#!/usr/bin/env python3
"""Migrate explicit workspaces: docs/workspace-* → docs/workspaces/workspace-*.

One-shot hygiene migration (protocol mount-point only):
  1. git mv each docs/workspace-* into docs/workspaces/
  2. Rewrite absolute path strings docs/workspace- → docs/workspaces/workspace-
  3. Fix relative markdown links that climb out of a workspace past docs/

Does NOT change goal id shape, intra-workspace parent links, or governance semantics.

Usage (repo root, clean tree recommended):
  python scripts/migrate_workspaces_container.py
  python scripts/migrate_workspaces_container.py --dry-run
"""

from __future__ import annotations

import argparse
import re
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DOCS = ROOT / "docs"
CONTAINER = DOCS / "workspaces"

# Text files that may carry path strings or relative links.
TEXT_SUFFIXES = {
    ".md",
    ".txt",
    ".py",
    ".ps1",
    ".sh",
    ".mjs",
    ".js",
    ".ts",
    ".tsx",
    ".json",
    ".yml",
    ".yaml",
    ".toml",
    ".go",
}

# Skip bulky / generated / VCS trees.
SKIP_DIR_NAMES = {
    ".git",
    "node_modules",
    "dist",
    "build",
    ".next",
    "coverage",
    "__pycache__",
    ".pytest_cache",
    "vendor",
    ".goal-governance-updates",
}

# Absolute path token (safe: docs/workspaces/ does not match docs/workspace-).
ABS_OLD = "docs/workspace-"
ABS_NEW = "docs/workspaces/workspace-"

# Windows-style absolute path in install.ps1 messages etc.
ABS_OLD_WIN = "docs\\workspace-"
ABS_NEW_WIN = "docs\\workspaces\\workspace-"

MD_LINK_RE = re.compile(r"(!?\[[^\]]*\])\(([^)]+)\)")
WS_NAME_RE = re.compile(r"^workspace-\d{3}-")


def list_old_workspaces() -> list[Path]:
    if not DOCS.is_dir():
        return []
    return sorted(
        p
        for p in DOCS.iterdir()
        if p.is_dir() and WS_NAME_RE.match(p.name) is not None
    )


def git_mv(src: Path, dst: Path, dry_run: bool) -> None:
    dst.parent.mkdir(parents=True, exist_ok=True)
    rel_src = src.relative_to(ROOT).as_posix()
    rel_dst = dst.relative_to(ROOT).as_posix()
    if dry_run:
        print(f"DRY git mv {rel_src} -> {rel_dst}")
        return
    subprocess.run(
        ["git", "mv", str(src), str(dst)],
        cwd=ROOT,
        check=True,
    )


def rewrite_absolute_paths(text: str) -> tuple[str, int]:
    n = 0
    if ABS_OLD in text:
        c = text.count(ABS_OLD)
        text = text.replace(ABS_OLD, ABS_NEW)
        n += c
    if ABS_OLD_WIN in text:
        c = text.count(ABS_OLD_WIN)
        text = text.replace(ABS_OLD_WIN, ABS_NEW_WIN)
        n += c
    return text, n


def _split_link_target(raw: str) -> tuple[str, str]:
    """Split destination from optional title; keep destination only for rewrite."""
    s = raw.strip()
    if not s:
        return raw, ""
    # angle-bracket destinations: <path>
    if s.startswith("<") and ">" in s:
        end = s.index(">")
        dest = s[1:end]
        rest = s[end + 1 :]
        return dest, rest
    # dest "title" or dest 'title'
    parts = s.split(None, 1)
    if len(parts) == 1:
        return parts[0], ""
    return parts[0], " " + parts[1]


def fix_relative_link_in_workspace(dest: str, file_parent_depth_in_ws: int) -> str | None:
    """Return new dest if a workspace-escaping relative link needs one more ../."""
    if not dest or dest.startswith(("#", "http://", "https://", "mailto:", "data:")):
        return None
    # Only pure relative ups
    norm = dest.replace("\\", "/")
    anchor = ""
    if "#" in norm:
        norm, frag = norm.split("#", 1)
        anchor = "#" + frag
    if not norm.startswith("../"):
        return None

    up = 0
    rest = norm
    while rest.startswith("../"):
        up += 1
        rest = rest[3:]
    if rest.startswith("./"):
        rest = rest[2:]

    parent_depth = file_parent_depth_in_ws
    if up <= parent_depth:
        return None  # stays inside workspace

    out = up - parent_depth
    first = rest.split("/")[0] if rest else ""
    # Sibling workspace under the new container: still one level out of ws root.
    if out == 1 and WS_NAME_RE.match(first):
        return None

    # Everything else that left the old docs/ flat layout needs +1 ../
    new_dest = "../" * (up + 1) + rest + anchor
    return new_dest


def fix_relative_links_outside_to_workspaces(dest: str, file_path: Path) -> str | None:
    """Fix ../workspace-NNN relative links from files outside the container."""
    if not dest:
        return None
    norm = dest.replace("\\", "/")
    anchor = ""
    if "#" in norm:
        norm, frag = norm.split("#", 1)
        anchor = "#" + frag

    # Match .../workspace-NNN-slug... that is NOT already .../workspaces/workspace-
    # Common forms: ../workspace-001-..., ../../workspace-001-...
    m = re.match(r"^((?:../)+)(workspace-\d{3}-[^#]*)$", norm)
    if not m:
        return None
    prefix, tail = m.group(1), m.group(2)
    # If this file already lives under docs/workspaces, skip (handled elsewhere).
    try:
        file_path.relative_to(CONTAINER)
        return None
    except ValueError:
        pass
    # Avoid double-prefix if somehow already workspaces/
    if "/workspaces/" in norm:
        return None
    return f"{prefix}workspaces/{tail}{anchor}"


def process_workspace_file(path: Path, ws_root: Path, dry_run: bool) -> int:
    if path.suffix.lower() not in TEXT_SUFFIXES and path.name not in {
        "Dockerfile",
        "Makefile",
    }:
        return 0
    try:
        original = path.read_text(encoding="utf-8")
    except UnicodeDecodeError:
        return 0

    text, n_abs = rewrite_absolute_paths(original)

    rel = path.relative_to(ws_root)
    parent_depth = len(rel.parts) - 1  # depth of parent dir inside ws

    changes = n_abs

    def repl(match: re.Match[str]) -> str:
        nonlocal changes
        label, raw_target = match.group(1), match.group(2)
        dest, rest = _split_link_target(raw_target)
        new_dest = fix_relative_link_in_workspace(dest, parent_depth)
        if new_dest is None:
            return match.group(0)
        changes += 1
        # Preserve angle brackets if original used them
        stripped = raw_target.strip()
        if stripped.startswith("<"):
            return f"{label}(<{new_dest}>{rest})"
        return f"{label}({new_dest}{rest})"

    text2 = MD_LINK_RE.sub(repl, text)
    if text2 != original and not dry_run:
        path.write_text(text2, encoding="utf-8", newline="\n")
    elif text2 != original and dry_run:
        print(f"DRY rewrite workspace file {path.relative_to(ROOT).as_posix()} ({changes})")
    return changes


def process_non_workspace_file(path: Path, dry_run: bool) -> int:
    if path.suffix.lower() not in TEXT_SUFFIXES and path.name not in {
        "Dockerfile",
        "Makefile",
        "AGENTS.md",
        "README.md",
        "QUICKSTART.md",
    }:
        return 0
    # Never rewrite this migrator (contains the old token as search pattern).
    if path.resolve() == Path(__file__).resolve():
        return 0
    # Skip files inside old or new workspace trees (handled separately / already moved)
    try:
        path.relative_to(CONTAINER)
        return 0
    except ValueError:
        pass

    try:
        original = path.read_text(encoding="utf-8")
    except UnicodeDecodeError:
        return 0

    text, n_abs = rewrite_absolute_paths(original)
    changes = n_abs

    def repl(match: re.Match[str]) -> str:
        nonlocal changes
        label, raw_target = match.group(1), match.group(2)
        dest, rest = _split_link_target(raw_target)
        new_dest = fix_relative_links_outside_to_workspaces(dest, path)
        if new_dest is None:
            return match.group(0)
        changes += 1
        stripped = raw_target.strip()
        if stripped.startswith("<"):
            return f"{label}(<{new_dest}>{rest})"
        return f"{label}({new_dest}{rest})"

    text2 = MD_LINK_RE.sub(repl, text)
    if text2 != original and not dry_run:
        path.write_text(text2, encoding="utf-8", newline="\n")
    elif text2 != original and dry_run:
        print(f"DRY rewrite {path.relative_to(ROOT).as_posix()} ({changes})")
    return changes


def iter_repo_files() -> list[Path]:
    out: list[Path] = []
    for p in ROOT.rglob("*"):
        if not p.is_file():
            continue
        # skip dirs
        parts = set(p.relative_to(ROOT).parts)
        if parts & SKIP_DIR_NAMES:
            continue
        if any(part in SKIP_DIR_NAMES for part in p.relative_to(ROOT).parts):
            continue
        out.append(p)
    return out


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--dry-run", action="store_true")
    args = ap.parse_args()
    dry = args.dry_run

    old = list_old_workspaces()
    if not old:
        # Maybe already migrated
        existing = (
            sorted(CONTAINER.iterdir())
            if CONTAINER.is_dir()
            else []
        )
        ws_existing = [p for p in existing if p.is_dir() and WS_NAME_RE.match(p.name)]
        if ws_existing:
            print(f"No docs/workspace-* left; container has {len(ws_existing)} workspaces.")
            print("Running path rewrite pass only…")
        else:
            print("No workspaces found to migrate.", file=sys.stderr)
            return 1
    else:
        if not dry:
            CONTAINER.mkdir(parents=True, exist_ok=True)
        print(f"Moving {len(old)} workspaces → docs/workspaces/")
        for src in old:
            dst = CONTAINER / src.name
            if dst.exists():
                print(f"Refusing overwrite: {dst}", file=sys.stderr)
                return 1
            git_mv(src, dst, dry)

    total = 0
    # Process workspace trees (new location, or dry-run still old location)
    ws_roots: list[Path] = []
    if dry and old:
        ws_roots = old
    else:
        if CONTAINER.is_dir():
            ws_roots = sorted(
                p
                for p in CONTAINER.iterdir()
                if p.is_dir() and WS_NAME_RE.match(p.name)
            )

    for ws in ws_roots:
        for f in ws.rglob("*"):
            if f.is_file():
                total += process_workspace_file(f, ws, dry)

    for f in iter_repo_files():
        # workspace files already processed
        skip = False
        for ws in ws_roots:
            try:
                f.relative_to(ws)
                skip = True
                break
            except ValueError:
                pass
        if skip:
            continue
        # also skip old flat paths if dry and not moved
        total += process_non_workspace_file(f, dry)

    print(f"Done. rewrite events≈{total} dry_run={dry}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
