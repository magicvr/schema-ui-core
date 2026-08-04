#!/usr/bin/env python3
"""Stage methodology mirrors into the Skills package (GOAL-022).

Canonical → package distribution:

  docs/templates/**              → skills/core/docs/templates/**
  docs/architecture/principles.md
  docs/architecture/workspace-protocol.md
                                 → skills/core/docs/architecture/
  docs/vision/alignment.md       → skills/core/docs/vision/alignment.md
  docs/contracts/consumer files  → skills/contracts/

Product-local architecture pages (overview, directory-layout, monorepo-layout,
module-architecture) are intentionally NOT staged: monorepo product docs may
diverge from the reusable methodology package.

Usage (repo root):
  python scripts/stage_skills_mirrors.py
  python scripts/stage_skills_mirrors.py --check   # exit 1 if drift
"""

from __future__ import annotations

import argparse
import shutil
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DOCS = ROOT / "docs"
CORE = ROOT / "skills" / "core" / "docs"
CONTRACTS_SRC = DOCS / "contracts"
CONTRACTS_DST = ROOT / "skills" / "contracts"

# Consumer contract profile only (0.11.0+). Producer matrix/runtime stay release-owned.
CONSUMER_CONTRACT_FILES = (
    "skills-consumer-contract.json",
    "skills-consumer-contract.schema.json",
)

ARCHITECTURE_MIRROR = (
    "principles.md",
    "workspace-protocol.md",
)

FIXTURE_SUBDIRS = ("fixtures/valid", "fixtures/invalid")


def _copy_file(src: Path, dst: Path) -> None:
    dst.parent.mkdir(parents=True, exist_ok=True)
    shutil.copy2(src, dst)


def _copy_tree(src: Path, dst: Path) -> None:
    if not src.is_dir():
        raise FileNotFoundError(f"missing source tree: {src}")
    if dst.exists():
        shutil.rmtree(dst)
    shutil.copytree(src, dst)


def stage() -> list[str]:
    changed: list[str] = []

    templates_src = DOCS / "templates"
    templates_dst = CORE / "templates"
    if not templates_src.is_dir():
        raise FileNotFoundError(f"missing {templates_src}")
    _copy_tree(templates_src, templates_dst)
    changed.append("skills/core/docs/templates/")

    for name in ARCHITECTURE_MIRROR:
        src = DOCS / "architecture" / name
        dst = CORE / "architecture" / name
        if not src.is_file():
            raise FileNotFoundError(f"missing {src}")
        _copy_file(src, dst)
        changed.append(f"skills/core/docs/architecture/{name}")

    alignment_src = DOCS / "vision" / "alignment.md"
    alignment_dst = CORE / "vision" / "alignment.md"
    if alignment_src.is_file():
        _copy_file(alignment_src, alignment_dst)
        changed.append("skills/core/docs/vision/alignment.md")

    CONTRACTS_DST.mkdir(parents=True, exist_ok=True)
    for name in CONSUMER_CONTRACT_FILES:
        src = CONTRACTS_SRC / name
        dst = CONTRACTS_DST / name
        if not src.is_file():
            raise FileNotFoundError(f"missing consumer contract: {src}")
        _copy_file(src, dst)
        changed.append(f"skills/contracts/{name}")

    for sub in FIXTURE_SUBDIRS:
        src_dir = CONTRACTS_SRC / Path(sub)
        dst_dir = CONTRACTS_DST / Path(sub)
        if src_dir.is_dir():
            _copy_tree(src_dir, dst_dir)
            changed.append(f"skills/contracts/{sub}/")

    return changed


def check_drift() -> list[str]:
    """Return list of drifted relative paths (empty if clean)."""
    drifts: list[str] = []

    def cmp_file(src: Path, dst: Path, rel: str) -> None:
        if not src.is_file():
            drifts.append(f"missing canonical: {rel}")
            return
        if not dst.is_file():
            drifts.append(f"missing mirror: {rel}")
            return
        if src.read_bytes() != dst.read_bytes():
            drifts.append(rel)

    for name in ARCHITECTURE_MIRROR:
        cmp_file(
            DOCS / "architecture" / name,
            CORE / "architecture" / name,
            f"architecture/{name}",
        )

    alignment = DOCS / "vision" / "alignment.md"
    if alignment.is_file():
        cmp_file(alignment, CORE / "vision" / "alignment.md", "vision/alignment.md")

    for path in (DOCS / "templates").rglob("*"):
        if path.is_file():
            rel = path.relative_to(DOCS / "templates").as_posix()
            cmp_file(path, CORE / "templates" / rel, f"templates/{rel}")

    for name in CONSUMER_CONTRACT_FILES:
        cmp_file(
            CONTRACTS_SRC / name,
            CONTRACTS_DST / name,
            f"contracts/{name}",
        )

    for sub in FIXTURE_SUBDIRS:
        src_dir = CONTRACTS_SRC / Path(sub)
        if not src_dir.is_dir():
            continue
        for path in src_dir.rglob("*.json"):
            rel = path.relative_to(CONTRACTS_SRC).as_posix()
            cmp_file(path, CONTRACTS_DST / rel, f"contracts/{rel}")

    return drifts


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--check",
        action="store_true",
        help="exit 1 if staged mirrors drift from canonical docs",
    )
    args = parser.parse_args(argv)

    if args.check:
        drifts = check_drift()
        if drifts:
            print("stage mirror drift:", file=sys.stderr)
            for item in drifts:
                print(f"  - {item}", file=sys.stderr)
            return 1
        print("stage mirrors OK (consumer methodology set)")
        return 0

    changed = stage()
    print("staged:")
    for item in changed:
        print(f"  {item}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
