#!/usr/bin/env python3
"""Minimal skills package packer for monorepo smoke / local zip.

Runs stage_skills_mirrors then zips the skills/ tree (excluding tests cache).
Producer release evidence / matrix gates live in the goal-governance producer
repo; this consumer monorepo only needs a stage-aware pack entrypoint.

Usage (repo root):
  python scripts/pack_skills_release.py --version X.Y.Z --output-dir dist/
"""

from __future__ import annotations

import argparse
import importlib.util
import sys
import zipfile
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SKILLS = ROOT / "skills"
STAGE_SCRIPT = ROOT / "scripts" / "stage_skills_mirrors.py"

SKIP_DIR_NAMES = {".git", "__pycache__", ".pytest_cache", "node_modules"}


def _load_stage():
    spec = importlib.util.spec_from_file_location("stage_skills_mirrors", STAGE_SCRIPT)
    if spec is None or spec.loader is None:
        raise RuntimeError(f"cannot load {STAGE_SCRIPT}")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def pack(version: str, output_dir: Path, skip_stage: bool = False) -> Path:
    if not skip_stage:
        stage = _load_stage()
        stage.stage()

    output_dir.mkdir(parents=True, exist_ok=True)
    zip_name = f"goal-governance-skills-v{version}.zip"
    zip_path = output_dir / zip_name
    root_prefix = f"goal-governance-skills-v{version}"

    with zipfile.ZipFile(zip_path, "w", compression=zipfile.ZIP_DEFLATED) as zf:
        for path in SKILLS.rglob("*"):
            if not path.is_file():
                continue
            if any(part in SKIP_DIR_NAMES for part in path.parts):
                continue
            arc = f"{root_prefix}/{path.relative_to(SKILLS).as_posix()}"
            zf.write(path, arcname=arc)

    return zip_path


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--version", required=True, help="package version X.Y.Z")
    parser.add_argument(
        "--output-dir",
        default="dist",
        help="directory for the zip (default: dist/)",
    )
    parser.add_argument(
        "--skip-stage",
        action="store_true",
        help="do not refresh mirrors before packing",
    )
    args = parser.parse_args(argv)

    zip_path = pack(
        version=args.version,
        output_dir=Path(args.output_dir),
        skip_stage=args.skip_stage,
    )
    print(f"wrote {zip_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
