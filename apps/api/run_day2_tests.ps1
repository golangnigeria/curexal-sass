$env:GOCACHE = "c:\Users\HomePC\Desktop\program\fullstack_Curexal\.cache\go-build"

Write-Host "==> 1. Testing Patient Module Services (MRN & MPI)..." -ForegroundColor Cyan
go test -v ./internal/modules/patient/service

Write-Host "`n==> 2. Testing Orchestration Module Services (Provider Profiles & Acuity)..." -ForegroundColor Cyan
go test -v ./internal/modules/orchestration/service

Write-Host "`n==> 3. Testing Operations Module Services (Appointments & Delivery Channels)..." -ForegroundColor Cyan
go test -v ./internal/modules/operations/service

Write-Host "`n==> 4. Testing Day 2 Integration Suites..." -ForegroundColor Cyan
go test -v -run "TestMPI|TestQueueWorkflow|TestDeliveryChannel|TestCriticalRoutes" ./internal/testing
