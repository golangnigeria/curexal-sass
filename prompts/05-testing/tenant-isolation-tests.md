# TENANT ISOLATION TESTING PROMPT
**File**: `prompts/05-testing/tenant-isolation-tests.md`  

Execute automated security testing for cross-tenant and cross-branch attack vectors:
1. Organization A user $\to$ Organization B patient records (Must return HTTP 403).
2. Branch A user $\to$ Branch B operational routes (Must redirect or return HTTP 403).
3. Clinician $\to$ Organization Executive HQ governance routes (Must return HTTP 403).
4. Receptionist/Scientist $\to$ Doctor consultation and prescription authoring (Must return HTTP 403).
5. Unassigned/Empty roles $\to$ All protected resources (Must fail closed with HTTP 401/403).
