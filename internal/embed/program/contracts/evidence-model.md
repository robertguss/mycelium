# Evidence Model Contract

## Claim classification

Every material claim should be classifiable as:

| Class                     | Meaning                                                                |
| ------------------------- | ---------------------------------------------------------------------- |
| Verified fact             | Directly confirmed through primary evidence or reliable measurement    |
| Official claim            | Stated by the responsible vendor, maintainer, standard, or institution |
| Independent corroboration | Confirmed by a strong independent source                               |
| Community observation     | Reported by practitioners but not independently proven                 |
| Experimental result       | Observed in a documented evidence spike                                |
| Inference                 | Reasoned from cited evidence                                           |
| Architectural judgment    | Decision balancing evidence and constraints                            |
| User decision             | Explicitly selected by the owner                                       |
| Hypothesis                | Not yet sufficiently verified                                          |

## Evidence Ledger (minimum fields)

| Field                    | Meaning                                                 |
| ------------------------ | ------------------------------------------------------- |
| Evidence ID              | Stable `EVD-###` if allocated                           |
| Claim                    | Proposition supported                                   |
| Classification           | One of the classes above                                |
| Source or spike          | Citation or `SPK-###`                                   |
| Source tier              | Charter-defined tier                                    |
| Date                     | Publication, release, or experiment date                |
| Access or execution date | When verified                                           |
| Confidence               | High, Medium, or Low                                    |
| Limitations              | What the evidence does not prove                        |
| Contradictory evidence   | Related conflict, if any                                |
| Downstream use           | `REC-###`, `REQ-###`, `DEC-###`, risk, or open question |
| Revalidation trigger     | When to check again                                     |

## Recommendation evidence threshold

A major recommendation must include: problem solved; requirements and
constraints; credible alternatives; supporting evidence; tradeoffs; confidence;
failure modes; revisit triggers. **Popularity alone is not sufficient.**

## Evidence before confidence

Confident prose is not evidence. Distinguish verified fact, official claim,
corroboration, community report, experiment, inference, judgment, user
preference, and unverified hypothesis.
