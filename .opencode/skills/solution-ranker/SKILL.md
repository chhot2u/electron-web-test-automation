---
name: solution-ranker
description: Solution variant generation and ranking specialist. Generates 3 distinct solution approaches for any coding task, scores them across 8 quality dimensions, and recommends the optimal path before code is written. Use when you want better results by thinking through alternatives first.
---

# Solution Ranker

Generate, score, and rank multiple solution variants for any coding task to ensure the best approach is chosen before writing code.

## When to Activate

- Implementing a new feature with multiple possible approaches
- Fixing a bug where the root cause or fix strategy is unclear
- Refactoring code and choosing between different restructuring strategies
- Making architectural decisions that affect multiple modules
- Migrating between libraries, frameworks, or patterns
- Optimizing performance when multiple strategies exist
- Any task where "just coding it" risks picking a suboptimal path

## Core Workflow: Variant-First

Every task follows this sequence:

### 1. Analyze

Evaluate the task across these dimensions:

| Dimension | What to Check |
|-----------|---------------|
| Intent | What is the user actually trying to achieve? |
| Constraints | Performance, compatibility, deadline, conventions |
| Scope | Which files, modules, APIs are affected? |
| Dependencies | What existing code, libraries, services does this touch? |
| Edge Cases | What can go wrong? Null values, concurrency, empty states |
| Success Criteria | How do we know the solution is correct? |

### 2. Generate 3 Variants

**Variant A — Minimal & Pragmatic**
Smallest change, reuse existing patterns, minimize diff, ship fast.

**Variant B — Robust & Scalable**
Production-grade with full error handling, validation, tests, extensibility.

**Variant C — Optimal Architecture**
Best possible design, may require refactoring, clean abstractions, performance-optimized.

Each variant includes:
- Approach summary
- Files affected (exact paths)
- Key code changes (pseudocode or snippets)
- Trade-offs (pros and cons)
- Risk level (LOW / MEDIUM / HIGH)
- Estimated complexity (LOW / MEDIUM / HIGH)

### 3. Score with Comparison Matrix

```
| Dimension          | Variant A | Variant B | Variant C |
|--------------------|-----------|-----------|-----------|
| Correctness        |    /10    |    /10    |    /10    |
| Maintainability    |    /10    |    /10    |    /10    |
| Performance        |    /10    |    /10    |    /10    |
| Testability        |    /10    |    /10    |    /10    |
| Risk (lower=safer) |    /10    |    /10    |    /10    |
| Simplicity         |    /10    |    /10    |    /10    |
| Extensibility      |    /10    |    /10    |    /10    |
| Convention Fit     |    /10    |    /10    |    /10    |
| TOTAL              |   /80     |   /80     |   /80     |
```

Scoring definitions:
- **Correctness**: Solves the problem fully, including edge cases
- **Maintainability**: Another developer can understand and modify in 6 months
- **Performance**: Efficient for the expected scale
- **Testability**: Can be unit/integration tested easily
- **Risk**: Likelihood of regressions or breakage (10 = safest)
- **Simplicity**: As simple as possible, no simpler
- **Extensibility**: Accommodates future changes
- **Convention Fit**: Matches project's existing patterns and style

### 4. Recommend

- Which variant wins and why
- What specific qualities drove the decision
- When to use each alternative
- Phased implementation plan
- Testing strategy (what to test, in what order)

## Solution Anti-Patterns to Detect

| Anti-Pattern | Better Approach |
|-------------|----------------|
| God function (>50 lines) | Extract into focused helpers |
| Deep nesting (>4 levels) | Early returns, guard clauses |
| Mutating inputs | Return new objects |
| Swallowed errors | Wrap with context |
| Hardcoded values | Constants, config, env vars |
| Missing validation | Validate at API boundaries |
| Copy-paste code | Extract shared utility |
| No tests | Write tests first (TDD) |

## Quick Modes

| Command | Action |
|---------|--------|
| `/rank-solutions analyze <task>` | Full analysis + 3 variants + matrix + recommendation |
| `/rank-solutions quick <task>` | Abbreviated variants + matrix + recommendation |
| `/rank-solutions compare <A> vs <B>` | Side-by-side scoring of two approaches |
| `/rank-solutions refactor <target>` | 3 refactoring strategies ranked |
| `/rank-solutions debug <bug>` | 3 fix strategies scored by correctness and risk |
| `/rank-solutions migrate <from> to <to>` | 3 migration strategies scored by risk and effort |

## Agent Integration

After the recommended variant is confirmed:

| Next Step | Agent to Use |
|-----------|-------------|
| Write tests first | **tdd-guide** |
| Complex multi-phase plan | **planner** |
| Review written code | **code-reviewer** |
| Security-sensitive code | **security-reviewer** |
| Go-specific review | **go-reviewer** |
| Build errors after implementation | **build-error-resolver** |

## Output Format

Every response follows this structure:

```
TASK ANALYSIS
=============
Problem: [statement]
Scope: [files and modules]
Constraints: [technical and project]
Edge Cases: [identified]

SOLUTION VARIANTS
=================
Variant A: [full details]
Variant B: [full details]
Variant C: [full details]

COMPARISON MATRIX
=================
[scoring table]

RECOMMENDATION
==============
Winner: [Variant X]
Why: [reasoning]
Alternatives: [when to use others]

IMPLEMENTATION PLAN
===================
[phased steps]

TESTING STRATEGY
================
[what and how to test]
```
