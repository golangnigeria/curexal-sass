# LEGACY PURGE & REFACTORING PROMPT
**File**: `prompts/07-refactoring/legacy-purge.md`  

Execute the canonical 7-step legacy purge process:
1. **DISCOVER**: Map all references, imports, routes, controllers, services, database tables, and tests.
2. **CLASSIFY**: Classify into KEEP, REFACTOR, REPLACE, or DELETE.
3. **DETERMINE REPLACEMENT**: Identify the canonical replacement in `/specs`.
4. **MIGRATE**: Migrate consumers to the canonical replacement.
5. **TEST**: Run test suites to verify zero regression.
6. **DELETE**: Safely delete obsolete code and unreferenced files.
7. **VERIFY**: Search the repository again to ensure zero dangling references.
