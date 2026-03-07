---
description: Generate, score, and rank solution variants for better code decisions
agent: solution-ranker
subtask: true
---

# Rank Solutions Command

Generate, score, and rank solution variants for: $ARGUMENTS

## Your Task

1. **Analyze** the task — read relevant source files, understand constraints, identify edge cases
2. **Generate 3 solution variants** using different strategies (minimal, robust, optimal)
3. **Score all variants** in a comparison matrix (1-10 across 8 dimensions)
4. **Recommend the best variant** with implementation plan and testing strategy
5. **Wait for confirmation** — MUST receive user approval before handing off to implementation

## Modes

If the user specifies a mode, follow it:

- `analyze <task>` — Full analysis + 3 variants + comparison + recommendation
- `quick <task>` — Abbreviated one-paragraph variants + matrix + recommendation
- `compare <approach1> vs <approach2>` — Side-by-side scoring of two approaches
- `refactor <file-or-function>` — 3 refactoring strategies scored and ranked
- `debug <bug>` — 3 fix strategies scored by correctness, risk, and speed
- `migrate <from> to <to>` — 3 migration strategies scored by risk and effort

If no mode is specified, default to `analyze`.

## Output Format

```
TASK ANALYSIS
=============
[Problem, scope, constraints, edge cases]

SOLUTION VARIANTS
=================
[3 variants with approach, files, code, trade-offs, risk]

COMPARISON MATRIX
=================
[Scoring table across 8 dimensions, total /80]

RECOMMENDATION
==============
[Winner + reasoning + alternatives]

IMPLEMENTATION PLAN
===================
[Phased steps for the recommended variant]

TESTING STRATEGY
================
[What to test and in what order]
```

**CRITICAL**: Do NOT write any code until the user explicitly confirms with "yes", "proceed", or similar affirmative response.
