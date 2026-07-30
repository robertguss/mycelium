# Artifact-Driven Research Program — operator commands
# Never runs git. Humans own commit/push.

set shell := ["bash", "-euo", "pipefail", "-c"]

# Default: list available recipes
default:
    @just --list

# Bootstrap project name into placeholders (edits files only; no git)
# Usage: just init "my-research-project"
#    or: just init name="my-research-project"
init name:
    python3 scripts/init_project.py {{quote(name)}}

# Show program status from the manifest
status:
    python3 scripts/status_program.py

# Sanity-check tree, placeholders vs accepted stages, required paths
check:
    python3 scripts/check_program.py
