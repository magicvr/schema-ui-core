#!/usr/bin/env python3
"""Transactional updater for an installed Goal Governance Skills package."""

from __future__ import annotations

import argparse
from datetime import datetime, timezone
from hashlib import sha256
import json
import os
from pathlib import Path, PurePosixPath
import re
import shutil
import stat
import subprocess
import sys
import tempfile
from urllib.request import Request, urlopen
import zipfile


SEMVER_RE = re.compile(
    r"^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)"
    r"(?:-[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?"
    r"(?:\+[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?$"
)


class UpdateError(RuntimeError):
    pass


def normalize_version(raw: str) -> str:
    value = (raw or "").strip()
    if value[:1].lower() == "v":
        value = value[1:]
    if not SEMVER_RE.fullmatch(value):
        raise UpdateError(f"invalid SemVer version: {raw!r}")
    return value


def file_sha256(path: Path) -> str:
    digest = sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def expected_digest(sidecar: Path, zip_name: str) -> str:
    match = re.fullmatch(
        r"([0-9a-fA-F]{64})[ \t]+(\S+)\s*",
        sidecar.read_text(encoding="utf-8").strip(),
    )
    if not match:
        raise UpdateError("invalid SHA-256 sidecar; expected '<hex>  <filename>'")
    if Path(match.group(2)).name != zip_name:
        raise UpdateError("SHA-256 sidecar filename does not match archive")
    return match.group(1).lower()


def verify_archive(zip_path: Path, sidecar: Path) -> str:
    expected = expected_digest(sidecar, zip_path.name)
    actual = file_sha256(zip_path)
    if actual != expected:
        raise UpdateError(f"SHA-256 mismatch: expected {expected}, got {actual}")
    return actual


def safe_extract(zip_path: Path, destination: Path) -> Path:
    roots: set[str] = set()
    with zipfile.ZipFile(zip_path) as archive:
        for info in archive.infolist():
            pure = PurePosixPath(info.filename)
            if pure.is_absolute() or not pure.parts or ".." in pure.parts:
                raise UpdateError(f"unsafe archive member: {info.filename}")
            mode = info.external_attr >> 16
            if stat.S_ISLNK(mode):
                raise UpdateError(f"archive symlink is not allowed: {info.filename}")
            roots.add(pure.parts[0])
        if len(roots) != 1:
            raise UpdateError("archive must contain exactly one top-level package directory")
        root_name = next(iter(roots))
        for info in archive.infolist():
            target = (destination / PurePosixPath(info.filename)).resolve()
            try:
                target.relative_to(destination.resolve())
            except ValueError as error:
                raise UpdateError(f"archive member escapes extraction root: {info.filename}") from error
        archive.extractall(destination)
    package = destination / root_name
    required = (
        package / "install.ps1",
        package / "install.sh",
        package / "contracts" / "skills-consumer-contract.json",
        package / "contracts" / "skills-consumer-contract.schema.json",
        package / "core" / "docs" / "architecture" / "principles.md",
    )
    missing = [str(path.relative_to(package)) for path in required if not path.is_file()]
    if missing:
        raise UpdateError("incomplete Skills package: " + ", ".join(missing))
    return package


def read_contract(package: Path) -> dict[str, object]:
    path = package / "contracts" / "skills-consumer-contract.json"
    try:
        payload = json.loads(path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as error:
        raise UpdateError(f"cannot read consumer contract: {path}") from error
    if not isinstance(payload, dict):
        raise UpdateError("consumer contract root must be an object")
    return payload


def protocol_version(contract: dict[str, object]) -> str:
    protocol = contract.get("protocol")
    if not isinstance(protocol, dict):
        raise UpdateError("consumer contract is missing protocol")
    return normalize_version(str(protocol.get("version") or ""))


def assert_protocol_compatible(current: str, incoming: str, allow_upgrade: bool) -> None:
    current_parts = tuple(int(part) for part in normalize_version(current).split(".")[:2])
    incoming_parts = tuple(int(part) for part in normalize_version(incoming).split(".")[:2])
    if current_parts != incoming_parts and not allow_upgrade:
        raise UpdateError(
            f"protocol boundary change {current} -> {incoming}; pass --allow-protocol-upgrade after review"
        )


def latest_version(repo: str) -> str:
    request = Request(
        f"https://api.github.com/repos/{repo}/releases/latest",
        headers={"Accept": "application/vnd.github+json", "User-Agent": "goal-governance-updater"},
    )
    try:
        with urlopen(request, timeout=30) as response:
            payload = json.load(response)
    except Exception as error:  # urllib surfaces several transport-specific subclasses
        raise UpdateError(f"latest release discovery failed: {error}") from error
    if not isinstance(payload, dict) or not payload.get("tag_name"):
        raise UpdateError("latest release response is missing tag_name")
    return normalize_version(str(payload["tag_name"]))


def download(url: str, destination: Path) -> None:
    request = Request(url, headers={"User-Agent": "goal-governance-updater"})
    try:
        with urlopen(request, timeout=60) as response, destination.open("wb") as output:
            shutil.copyfileobj(response, output)
    except Exception as error:
        raise UpdateError(f"download failed: {url}: {error}") from error


def managed_file_pairs(package: Path, target: Path) -> list[tuple[Path, Path]]:
    pairs: list[tuple[Path, Path]] = []
    fixed = {
        "install/claude/AGENTS.md": "AGENTS.md",
        "install/copilot/copilot-instructions.md": ".github/copilot-instructions.md",
        "core/docs/README.md": "docs/README.md",
    }
    for source, destination in fixed.items():
        pairs.append((package / source, target / destination))
    host_roots = {
        "install/claude/skills": ".claude/skills",
        "install/grok/skills": ".grok/skills",
        "install/codex/skills": ".agents/skills",
        "core/docs/architecture": "docs/architecture",
        "core/docs/templates": "docs/templates",
        "core/docs/vision": "docs/vision",
    }
    for source_root, destination_root in host_roots.items():
        root = package / source_root
        if not root.is_dir():
            continue
        for source in sorted(root.rglob("*")):
            if source.is_file() and not source.name.startswith("."):
                pairs.append((source, target / destination_root / source.relative_to(root)))
    copilot_prompts = package / "install" / "copilot" / "prompts"
    if copilot_prompts.is_dir():
        for source in sorted(copilot_prompts.glob("*.md")):
            pairs.append((source, target / ".github" / "prompts" / f"{source.stem}.prompt.md"))
    return pairs


def modified_managed_files(
    package: Path,
    target: Path,
    incoming_package: Path | None = None,
) -> list[Path]:
    current = {destination: source for source, destination in managed_file_pairs(package, target)}
    incoming = (
        {destination: source for source, destination in managed_file_pairs(incoming_package, target)}
        if incoming_package is not None
        else {}
    )
    modified: list[Path] = []
    for destination, source in current.items():
        if source.is_file() and destination.is_file() and file_sha256(source) != file_sha256(destination):
            modified.append(destination)
    for destination, source in incoming.items():
        if destination in current or not source.is_file() or not destination.is_file():
            continue
        if file_sha256(source) != file_sha256(destination):
            modified.append(destination)
    return sorted(set(modified))


def backup_external_files(
    package: Path,
    target: Path,
    backup: Path,
    incoming_package: Path | None = None,
) -> list[str]:
    absent: list[str] = []
    project_backup = backup / "project"
    destinations = {destination for _source, destination in managed_file_pairs(package, target)}
    if incoming_package is not None:
        destinations.update(
            destination for _source, destination in managed_file_pairs(incoming_package, target)
        )
    for destination in sorted(destinations):
        relative = destination.relative_to(target)
        if destination.is_file():
            saved = project_backup / relative
            saved.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy2(destination, saved)
        else:
            absent.append(relative.as_posix())
    return absent


def restore_external_files(target: Path, backup: Path, absent: list[str]) -> None:
    for relative in absent:
        path = target / PurePosixPath(relative)
        if path.is_file():
            path.unlink()
    project_backup = backup / "project"
    if project_backup.is_dir():
        for source in sorted(project_backup.rglob("*")):
            if source.is_file():
                destination = target / source.relative_to(project_backup)
                destination.parent.mkdir(parents=True, exist_ok=True)
                shutil.copy2(source, destination)


def run_installer(package: Path, target: Path) -> None:
    if os.name == "nt":
        executable = shutil.which("powershell") or shutil.which("pwsh")
        if not executable:
            raise UpdateError("PowerShell is required to apply the update on Windows")
        command = [
            executable,
            "-NoProfile",
            "-ExecutionPolicy",
            "Bypass",
            "-File",
            str(package / "install.ps1"),
            "-All",
            "-NonInteractive",
            "-Force",
            "-SkillsDir",
            str(package),
        ]
    else:
        bash = shutil.which("bash")
        if not bash:
            raise UpdateError("bash is required to apply the update on this platform")
        command = [
            bash,
            str(package / "install.sh"),
            "--all",
            "--non-interactive",
            "--force",
            "--skills-dir",
            str(package),
        ]
    result = subprocess.run(command, cwd=target, text=True, capture_output=True, check=False)
    if result.returncode != 0:
        raise UpdateError(
            f"package install failed ({result.returncode})\nstdout:\n{result.stdout}\nstderr:\n{result.stderr}"
        )


def update_package(args: argparse.Namespace) -> dict[str, object]:
    target = Path(args.target_dir).resolve()
    skills = (target / args.skills_dir).resolve() if not Path(args.skills_dir).is_absolute() else Path(args.skills_dir).resolve()
    try:
        skills.relative_to(target)
    except ValueError as error:
        raise UpdateError("skills directory must stay within target directory") from error
    if not skills.is_dir():
        raise UpdateError(f"installed Skills directory not found: {skills}")

    version = latest_version(args.repo) if args.latest else normalize_version(args.version)
    current_protocol = protocol_version(read_contract(skills))
    archive_name = f"goal-governance-skills-v{version}.zip"
    tag = args.release_tag or f"v{version}"

    with tempfile.TemporaryDirectory(prefix="goal-governance-update-") as temp_name:
        temp = Path(temp_name)
        if args.zip_path:
            archive = Path(args.zip_path).resolve()
            sidecar = Path(args.sha256_path).resolve() if args.sha256_path else Path(str(archive) + ".sha256")
        else:
            archive = temp / archive_name
            sidecar = temp / f"{archive_name}.sha256"
            base = f"https://github.com/{args.repo}/releases/download/{tag}"
            download(f"{base}/{archive_name}", archive)
            download(f"{base}/{archive_name}.sha256", sidecar)
        if not archive.is_file() or not sidecar.is_file():
            raise UpdateError("archive or SHA-256 sidecar does not exist")
        digest = verify_archive(archive, sidecar)
        incoming = safe_extract(archive, temp / "extract")
        incoming_protocol = protocol_version(read_contract(incoming))
        assert_protocol_compatible(current_protocol, incoming_protocol, args.allow_protocol_upgrade)

        modified = modified_managed_files(skills, target, incoming)
        if modified and not args.force_managed:
            shown = ", ".join(str(path.relative_to(target)) for path in modified[:8])
            raise UpdateError(f"managed files have local changes; review or pass --force-managed: {shown}")

        plan = {
            "version": version,
            "current_protocol": current_protocol,
            "incoming_protocol": incoming_protocol,
            "archive_sha256": digest,
            "managed_conflicts": [str(path.relative_to(target)) for path in modified],
        }
        if args.dry_run:
            return {**plan, "result": "dry-run"}

        stamp = datetime.now(timezone.utc).strftime("%Y%m%dT%H%M%SZ")
        backup = target / ".goal-governance-updates" / f"{stamp}-v{version}"
        backup.mkdir(parents=True, exist_ok=False)
        absent = backup_external_files(skills, target, backup, incoming)
        old_skills = backup / "skills"
        shutil.move(str(skills), str(old_skills))
        try:
            shutil.move(str(incoming), str(skills))
            if not args.skip_install:
                run_installer(skills, target)
            state = {
                "format": "goal-governance.skills-install-state/v1",
                "version": version,
                "protocol": incoming_protocol,
                "archiveSha256": digest,
                "source": str(archive) if args.zip_path else f"{args.repo}@{tag}",
                "installedAt": datetime.now(timezone.utc).isoformat(),
                "rollbackPath": str(backup),
            }
            (skills / ".goal-governance-install.json").write_text(
                json.dumps(state, ensure_ascii=False, indent=2) + "\n",
                encoding="utf-8",
            )
        except Exception:
            failed = backup / "failed-new-skills"
            if skills.exists():
                shutil.move(str(skills), str(failed))
            shutil.move(str(old_skills), str(skills))
            restore_external_files(target, backup, absent)
            raise
        return {**plan, "result": "updated", "rollback_path": str(backup)}


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description=__doc__)
    version = parser.add_mutually_exclusive_group(required=True)
    version.add_argument("--version", help="target SemVer (optional leading v)")
    version.add_argument("--latest", action="store_true", help="discover the latest GitHub Release")
    parser.add_argument("--target-dir", default=".", help="consumer project root")
    parser.add_argument("--skills-dir", default="skills", help="Skills path under target")
    parser.add_argument("--zip-path", help="offline release zip")
    parser.add_argument("--sha256-path", help="offline SHA-256 sidecar")
    parser.add_argument("--repo", default="magicvr/goal-governance")
    parser.add_argument("--release-tag", default="")
    parser.add_argument("--allow-protocol-upgrade", action="store_true")
    parser.add_argument("--force-managed", action="store_true")
    parser.add_argument("--dry-run", action="store_true")
    parser.add_argument("--skip-install", action="store_true", help=argparse.SUPPRESS)
    return parser


def main(argv: list[str] | None = None) -> int:
    try:
        result = update_package(build_parser().parse_args(argv))
    except UpdateError as error:
        print(f"Error: {error}", file=sys.stderr)
        return 1
    print(json.dumps(result, ensure_ascii=False, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
