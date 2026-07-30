#!/usr/bin/env python3
"""Sanity-check research program tree and acceptance consistency."""

from __future__ import annotations

import re
import sys
from pathlib import Path


REQ_FILES = [
    "README.md",
    "AGENTS.md",
    "research-program.toml",
    "Justfile",
    "docs/00-program-blueprint.md",
    "docs/01-research-charter.md",
    "docs/specifications/01-definitive-specification.md",
    "docs/specifications/02-definitive-specification-revised.md",
    "docs/plans/01-implementation-plan.md",
    "docs/plans/02-implementation-plan-revised.md",
    "docs/reviews/01-specification-adversarial-review.md",
    "docs/reviews/02-implementation-plan-adversarial-review.md",
    "program/README.md",
    "decisions/README.md",
    ".agents/skills/research-program/SKILL.md",
    ".agents/skills/research-stage/SKILL.md",
    ".agents/skills/research-validate/SKILL.md",
]

REQ_DIRS = [
    "docs/prompts",
    "docs/reports",
    "docs/reconciliations",
    "docs/evidence",
    "docs/handoffs",
    "docs/validations",
    "program/operator",
    "program/contracts",
    "program/templates",
    "program/reference",
    ".agents/skills",
    "scripts",
]


def field(block: str, key: str, default: str = "") -> str:
    m = re.search(rf'^{key}\s*=\s*"(.*)"', block, re.M)
    if m:
        return m.group(1)
    m = re.search(rf'^{key}\s*=\s*\[(.*?)\]', block, re.M | re.S)
    if m:
        return m.group(1)
    return default


def main() -> int:
    print("=== just check ===")
    fail = 0
    root = Path.cwd()

    for rel in REQ_FILES:
        if not (root / rel).is_file():
            print(f"MISSING: {rel}")
            fail = 1

    for rel in REQ_DIRS:
        if not (root / rel).is_dir():
            print(f"MISSING DIR: {rel}")
            fail = 1

    toml = root / "research-program.toml"
    if toml.is_file():
        text = toml.read_text(encoding="utf-8")
        token_project = "{" + "{PROJECT_NAME}}"
        token_id = "{" + "{PROGRAM_ID}}"
        token_date = "{" + "{CREATED_DATE}}"
        if token_project in text or token_id in text or token_date in text:
            print(
                'WARN: research-program.toml still has init placeholders — '
                'run: just init name="..."'
            )

        stages = re.split(r"\n\[\[stages\]\]\n", text)[1:]
        for block in stages:
            sid = field(block, "id")
            status = field(block, "status")
            if status != "accepted":
                continue
            outs = re.findall(r'"([^"]+)"', field(block, "outputs", ""))
            for path in outs:
                p = root / path
                if not p.is_file():
                    print(f"FAIL: stage {sid} accepted but missing output: {path}")
                    fail = 1
                    continue
                body = p.read_text(encoding="utf-8")
                if "Placeholder — not accepted" in body:
                    print(
                        f"FAIL: stage {sid} accepted but output still placeholder: {path}"
                    )
                    fail = 1
                if token_project in body or token_id in body:
                    print(
                        f"FAIL: stage {sid} accepted but output has init placeholders: {path}"
                    )
                    fail = 1

    if fail:
        print("check: FAILED")
        return 1

    print("check: OK")
    print(
        "Note: OK means tree shape and acceptance consistency; "
        "not that research is done."
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
