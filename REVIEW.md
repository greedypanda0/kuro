# Kuro Review Notes

## Summary
I reviewed the current `core` and `cli` codebase against the goals in the README and then applied the key improvements we discussed. The core now handles repo root discovery, ignore semantics, path normalization, transactional updates, schema migrations, and baseline invariant tests.

## Ratings
- Architecture: 8/10
- Code quality: 8/10
- Correctness and invariants: 7/10
- CLI ergonomics: 7/10
- Maintainability: 8/10
- Overall: 8/10

## What is working well
- Clear separation between `core` and `cli`
- SQLite schema maps cleanly to VCS concepts
- UI layer is centralized and consistent
- Explicit errors in `core/errors` improve clarity

## Applied improvements
1. Repo root discovery  
   Added a repo root resolver and updated CLI commands to use repo-relative paths.

2. Ignore semantics  
   Reworked ignore handling to support pattern matching and consistent repo-relative checks.

3. Path normalization  
   Normalized staged paths to repo-relative, forward-slash paths and prevented staging outside the repo.

4. Transactions for multi-step changes  
   Wrapped staging and branch operations in database transactions.

5. Schema migrations  
   Introduced a migration table and incremental migration execution.

6. Tests for invariants  
   Added baseline tests for default refs/HEAD and staging lifecycle.

7. Core cleanup  
   Removed the stray `main` in `core` to keep the package focused.

## Proposed next steps
- Expand ignore support with negation and anchored patterns
- Add tests for branch creation/deletion edge cases and HEAD invariants
- Add migration tests to ensure forward compatibility
- Add CLI-level tests for `add`, `remove`, and `status`
- Consider a pathspec utility to unify matching across features