#!/usr/bin/env python3
"""Bootstrap project name into template placeholders. Does not run git."""

from __future__ import annotations

import re
import sys
from datetime import date
from pathlib import Path


def normalize_name(raw: str) -> str:
    """Accept 'My Project' or accidental 'name=My Project' from shell/just."""
    name = raw.strip()
    if name.startswith("name="):
        name = name[5:].strip()
    # Strip surrounding quotes if present
    if len(name) >= 2 and name[0] == name[-1] and name[0] in {'"', "'"}:
        name = name[1:-1].strip()
    return name


def should_process(path: Path, root: Path) -> bool:
    rel = path.relative_to(root)
    parts = rel.parts
    if not parts:
        return False
    # Never rewrite methodology, skills, or tooling
    if parts[0] in {".git", "program", "scripts", ".agents"}:
        return False
    if path.suffix not in {".md", ".toml"}:
        return False
    return True


def main() -> int:
    if len(sys.argv) < 2:
        print(
            'error: name is required, e.g. just init name="my-project"',
            file=sys.stderr,
        )
        return 1

    name = normalize_name(sys.argv[1])
    if not name:
        print(
            'error: name is required, e.g. just init name="my-project"',
            file=sys.stderr,
        )
        return 1

    slug = re.sub(r"[^a-z0-9]+", "-", name.lower()).strip("-")
    if not slug:
        print("error: could not derive program_id from name", file=sys.stderr)
        return 1

    today = date.today().isoformat()
    root = Path.cwd()
    if not (root / "research-program.toml").is_file():
        print(
            "error: research-program.toml not found; run from template root",
            file=sys.stderr,
        )
        return 1

    print("Initializing project:")
    print(f"  program_name: {name}")
    print(f"  program_id:   {slug}")
    print(f"  date:         {today}")

    token_project = "{" + "{PROJECT_NAME}}"
    token_id = "{" + "{PROGRAM_ID}}"
    token_date = "{" + "{CREATED_DATE}}"

    updated: list[str] = []
    for path in root.rglob("*"):
        if not path.is_file() or not should_process(path, root):
            continue
        try:
            text = path.read_text(encoding="utf-8")
        except (UnicodeDecodeError, OSError):
            continue
        if (
            token_project not in text
            and token_id not in text
            and token_date not in text
        ):
            continue
        new = (
            text.replace(token_project, name)
            .replace(token_id, slug)
            .replace(token_date, today)
        )
        path.write_text(new, encoding="utf-8")
        updated.append(path.relative_to(root).as_posix())

    for u in sorted(updated):
        print(f"  updated: {u}")
    if not updated:
        print("  (no placeholders found — already initialized?)")
    print("Done. Next: discovery interview (see program/operator/getting-started.md).")
    print("Rigor tier remains 'standard' (proposed) until Blueprint approval.")
    print("This command did not run git.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
