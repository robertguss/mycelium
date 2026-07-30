#!/usr/bin/env python3
"""Print research-program.toml stage status and eligible next stages."""

from __future__ import annotations

import re
import sys
from pathlib import Path


def field(block: str, key: str, default: str = "") -> str:
    m = re.search(rf'^{key}\s*=\s*"(.*)"', block, re.M)
    if m:
        return m.group(1)
    m = re.search(rf'^{key}\s*=\s*\[(.*?)\]', block, re.M | re.S)
    if m:
        return m.group(1).strip()
    return default


def main() -> int:
    path = Path("research-program.toml")
    if not path.is_file():
        print("error: research-program.toml not found", file=sys.stderr)
        return 1

    text = path.read_text(encoding="utf-8")
    print("=== Research program status ===\n")
    for key in (
        "program_id",
        "program_name",
        "rigor_tier",
        "status",
        "created_date",
        "last_updated_date",
    ):
        m = re.search(rf"^{key}\s*=\s*(.+)$", text, re.M)
        if m:
            print(f"{key} = {m.group(1).strip()}")

    stages = re.split(r"\n\[\[stages\]\]\n", text)[1:]
    print("\n=== Stages ===\n")
    if not stages:
        print("(no stages found)")
        return 0

    rows = [
        (
            field(block, "id"),
            field(block, "status"),
            field(block, "name"),
            field(block, "depends_on", ""),
        )
        for block in stages
    ]
    accepted = {r[0] for r in rows if r[1] == "accepted"}

    print(f"{'ID':<22} {'STATUS':<22} {'NAME':<42} DEPENDS_ON")
    print("-" * 110)
    for sid, status, name, deps in rows:
        print(f"{sid:<22} {status:<22} {name:<42} {deps}")

    print("\n=== Eligible next stages (prerequisites all accepted) ===")
    eligible = []
    for block in stages:
        sid = field(block, "id")
        status = field(block, "status")
        if status in ("accepted", "superseded", "cancelled"):
            continue
        deps = re.findall(r'"([^"]+)"', field(block, "depends_on", ""))
        if all(d in accepted for d in deps):
            eligible.append((sid, status, field(block, "name")))

    if not eligible:
        print("(none — check blocked stages, incomplete prereqs, or finish discovery)")
    else:
        for sid, status, name in eligible:
            print(f"  - {sid} ({status}): {name}")

    print()
    print(
        "Placeholders do not count as accepted. "
        "See program/operator/resume-protocol.md"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
