# Commit Boundaries

Each accepted substantive artifact should normally receive its own coherent
commit. **Humans run git**; Just recipes do not.

Recommended patterns:

```text
docs: bootstrap research program
docs: add program blueprint
docs: add research charter
docs: add <track> research prompt
docs: add <track> research report
docs: reconcile replicated <track> research
docs: add definitive specification
docs: add specification adversarial review
docs: publish revised definitive specification
docs: add implementation plan
docs: add implementation plan adversarial review
docs: publish final revised implementation plan
```

Prompt installation and report execution should not be mixed in one commit
unless repository policy explicitly requires it.
