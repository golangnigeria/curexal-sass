# ARCHITECTURE & SECURITY CODE REVIEW PROMPT
**File**: `prompts/04-review/code-review.md`  

Review proposed code changes against `/specs`:
1. Are all boundaries and tenant isolations strictly enforced on the server?
2. Are all database operations wrapped in explicit transactions with constraints?
3. Does the frontend act strictly as a projection of backend authorization?
4. Are all error responses predictable, safe, and free of sensitive PHI or stack traces?
5. Do tests validate business rules and security attack vectors?
