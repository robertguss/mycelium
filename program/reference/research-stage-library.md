# Research Stage Library

The architect selects only tracks justified by the project.

| Track                                   | Answers                                                                 | Use when                                                                 |
| --------------------------------------- | ----------------------------------------------------------------------- | ------------------------------------------------------------------------ |
| Domain and problem                      | Real problem, domain model, vocabulary, workflow, institutional context | Domain complex, regulated, specialized, or poorly understood             |
| User and workflow                       | Users, jobs, friction, success                                          | User behavior shapes architecture or scope                               |
| Ecosystem, tooling, dependency          | Languages, frameworks, libraries, tools, standards, versions            | Ecosystem broad, volatile, or consequential                              |
| Architecture and system design          | Best architecture; component interaction                                | Multiple plausible architectures or high integration complexity          |
| Security and threat model               | Assets, trust boundaries, threats, controls, residual risk              | Credentials, sensitive data, untrusted input, privilege, public exposure |
| Data and integration                    | APIs, contracts, consistency, rate limits, failure modes, migration     | External systems or data quality architecture-defining                   |
| Testing and verification                | Methods that prove desired properties                                   | Correctness, AI code quality, reliability, state-space complexity        |
| Operations, deployment, reliability     | Deploy, observe, recover, operate                                       | Runtime ops affect design                                                |
| AI-native repository and agent workflow | Repos, instructions, boundaries, checks for coding agents               | Agents will implement or maintain substantially                          |
| Performance and scalability             | Workloads, budgets, bottlenecks                                         | Perf/scale could change architecture                                     |
| Migration and compatibility             | Safe move of users, data, APIs, workflows                               | Replacement, transition, compatibility, staged adoption                  |
| Legal, regulatory, privacy, compliance  | Obligations on behavior, data, records, distribution                    | Regulated data, contracts, licensing, jurisdiction                       |
| Financial, cost, feasibility            | Cost to build, operate, maintain, replace                               | Cost, ROI, budget central                                                |
| Market and competitive                  | Alternatives, user value, differentiation                               | Product viability or positioning                                         |
| Scientific or empirical validation      | Hypothesis under controlled observation                                 | Empirical claims over engineering convention                             |
| Risk and failure-mode                   | Technical, operational, organizational, economic failure                | Novel, consequential, hard to reverse                                    |

## Stage selection rule

For every selected track: why it exists; why another cannot absorb it; which
decision it informs; which artifact consumes it. For obvious omitted tracks:
brief why unnecessary.

**Avoid overbuilding.** Simplicity has positive weight.
