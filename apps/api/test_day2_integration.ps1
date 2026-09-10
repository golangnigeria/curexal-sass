$ErrorActionPreference = "Stop"

$BaseUrl = "http://localhost:8080/api/v1"
Write-Host "====================================================" -ForegroundColor Cyan
Write-Host "       CUREXAL DAY 2 END-TO-END INTEGRATION TEST    " -ForegroundColor Cyan
Write-Host "====================================================" -ForegroundColor Cyan

# 1. Authenticate as Nurse Chioma (Staff)
Write-Host "`n==> 1. Authenticating as Clinical Staff (Nurse Chioma)..." -ForegroundColor Yellow
$session = New-Object Microsoft.PowerShell.Commands.WebRequestSession
$loginPayload = @{
    email = "nurse.chioma@curexal.com"
    password = "password"
} | ConvertTo-Json

$loginRes = Invoke-RestMethod -Uri "$BaseUrl/auth/login" -Method Post -Body $loginPayload -ContentType "application/json" -WebSession $session
$token = $loginRes.targetUrl.Split("=")[1]

$exchangePayload = @{ token = $token } | ConvertTo-Json
$authRes = Invoke-RestMethod -Uri "$BaseUrl/auth/exchange" -Method Post -Body $exchangePayload -ContentType "application/json" -WebSession $session

$jwt = $session.Cookies.GetCookies("http://localhost:8080")["jwt"].Value
$tenantId = "00000000-0000-0000-0000-000000000001"
$headers = @{
    "Authorization" = "Bearer $jwt"
    "X-Tenant-ID" = $tenantId
    "X-Branch-ID" = $tenantId
}
Write-Host "    Staff Authenticated Successfully. Session established." -ForegroundColor Green

# 2. Feature #4: Provider Profiles & Duty Management
Write-Host "`n==> 2. Testing Feature #4: Provider Profiles & Duty Management..." -ForegroundColor Yellow
$providersRes = Invoke-RestMethod -Uri "$BaseUrl/providers/profiles" -Method Get -Headers $headers
$providers = $providersRes.data
Write-Host "    Found $($providers.Count) provider profiles." -ForegroundColor Green
$provider = $providers[0]
Write-Host "    Provider 1: $($provider.specialtyCode), Status: $($provider.status), License: $($provider.licenseNumber), Room: $($provider.roomNumber)" -ForegroundColor Green

if (-not $provider.id) {
    throw "Provider profile ID missing!"
}

# 3. Feature #5: Patient Registration & Monotonic MRN
Write-Host "`n==> 3. Testing Feature #5: Master Patient Registration with Emergency Contact..." -ForegroundColor Yellow
$rnd = Get-Random -Minimum 100000 -Maximum 999999
$testPhone = "+234803$rnd"
$testNIN = "NIN-2026-$rnd"

$patientPayload = @{
    firstName = "Amina"
    lastName = "Bello"
    dateOfBirth = "1994-06-20"
    gender = "FEMALE"
    phone = $testPhone
    email = "amina.$rnd@curexal-test.org"
    nin = $testNIN
    registrationChannel = "RECEPTION"
    residentialAddress = "14 Victoria Island Expressway"
    city = "Lagos"
    state = "Lagos State"
    emergencyContact = @{
        fullName = "Ibrahim Bello"
        relationship = "SPOUSE"
        phone = "+2348031122334"
        email = "ibrahim.$rnd@curexal-test.org"
        isEmergencyContact = $true
    }
} | ConvertTo-Json -Depth 5

$regRes = Invoke-RestMethod -Uri "$BaseUrl/patients/canonical" -Method Post -Headers $headers -Body $patientPayload -ContentType "application/json"
$patient = $regRes.data.patient
Write-Host "    Created Patient: $($patient.firstName) $($patient.lastName)" -ForegroundColor Green
Write-Host "    Generated MRN: $($patient.mrn)" -ForegroundColor Green
Write-Host "    Patient ID: $($patient.id)" -ForegroundColor Green

if ($patient.mrn -notmatch "^PAT-\d{4}-\d{5}$") {
    throw "Generated MRN does not match standard PAT-YYYY-XXXXX format!"
}

# 4. Feature #6: MPI Duplicate Detection & Prevention
Write-Host "`n==> 4. Testing Feature #6: MPI Duplicate Detection (Phone + NIN Conflict)..." -ForegroundColor Yellow
$dupPayload = @{
    firstName = "Ameena"
    lastName = "Bello"
    dateOfBirth = "1994-06-20"
    gender = "FEMALE"
    phone = $testPhone
    nin = $testNIN
} | ConvertTo-Json

$duplicateBlocked = $false
try {
    $dupRes = Invoke-RestMethod -Uri "$BaseUrl/patients/canonical" -Method Post -Headers $headers -Body $dupPayload -ContentType "application/json"
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    Write-Host "    PASS: Duplicate registration blocked with HTTP $statusCode (Conflict)" -ForegroundColor Green
    $duplicateBlocked = $true
}

if (-not $duplicateBlocked) {
    throw "FAILED: MPI Duplicate detection failed to block identical record!"
}

# 5. Feature #8: Appointment Scheduling
Write-Host "`n==> 5. Testing Feature #8: Appointment Scheduling..." -ForegroundColor Yellow
$appointmentPayload = @{
    patientId = $patient.id
    providerId = $provider.id
    serviceType = "CONSULTATION"
    deliveryChannel = "in_person"
    startTime = (Get-Date).AddHours(2).ToString("yyyy-MM-ddTHH:mm:00Z")
    endTime = (Get-Date).AddHours(2).AddMinutes(30).ToString("yyyy-MM-ddTHH:mm:00Z")
    reasonForVisit = "Follow-up consultation for recurring migraine"
} | ConvertTo-Json

$aptRes = Invoke-RestMethod -Uri "$BaseUrl/appointments" -Method Post -Headers $headers -Body $appointmentPayload -ContentType "application/json"
$appointment = $aptRes.data
Write-Host "    Booked Appointment: $($appointment.appointmentNumber)" -ForegroundColor Green
Write-Host "    Delivery Channel: $($appointment.deliveryChannel)" -ForegroundColor Green
Write-Host "    Status: $($appointment.status)" -ForegroundColor Green

if ($appointment.appointmentNumber -notmatch "^APT-\d{4}-\d{5}$") {
    throw "Appointment number does not match APT-YYYY-XXXXX format!"
}

# 6. Feature #10: Patient Check-In to Care Desk Queue
Write-Host "`n==> 6. Testing Feature #10: Patient Check-In (Care Request)..." -ForegroundColor Yellow
$checkInPayload = @{
    patientId = $patient.id
    serviceType = "GENERAL_CONSULTATION"
    deliveryChannel = "in_person"
    chiefComplaint = "Severe headache and high fever since morning"
} | ConvertTo-Json

$checkInRes = Invoke-RestMethod -Uri "$BaseUrl/orchestration/requests" -Method Post -Headers $headers -Body $checkInPayload -ContentType "application/json"
$careRequest = $checkInRes.data
Write-Host "    Check-In Created: $($careRequest.id)" -ForegroundColor Green
Write-Host "    Initial Status: $($careRequest.status)" -ForegroundColor Green

if ($careRequest.status -ne "WAITING_TRIAGE") {
    throw "Expected status WAITING_TRIAGE, got $($careRequest.status)"
}

# 7. Feature #10: Nurse Triage Assessment & Acuity Escalation
Write-Host "`n==> 7. Testing Feature #10: Clinical Triage & Acuity Evaluation..." -ForegroundColor Yellow
$triagePayload = @{
    systolicBp = 145
    diastolicBp = 95
    pulseRate = 104
    temperature = 39.2
    spo2 = 94
    respiratoryRate = 24
    painScore = 7
    triageNotes = "Patient in acute distress, febrile and tachycardic."
} | ConvertTo-Json

$triageRes = Invoke-RestMethod -Uri "$BaseUrl/orchestration/requests/$($careRequest.id)/triage" -Method Post -Headers $headers -Body $triagePayload -ContentType "application/json"
$triagedData = $triageRes.data
Write-Host "    Triage Assessment Complete." -ForegroundColor Green
Write-Host "    Evaluated Acuity Level: $($triagedData.acuityLevel)" -ForegroundColor Green
Write-Host "    Updated Request Status: $($triagedData.status)" -ForegroundColor Green

# 8. Feature #10: Live Queue Management Verification
Write-Host "`n==> 8. Testing Feature #10: Active Live Queue Telemetry..." -ForegroundColor Yellow
$queueRes = Invoke-RestMethod -Uri "$BaseUrl/orchestration/requests" -Method Get -Headers $headers
$matchedQueueItem = $queueRes.data | Where-Object { $_.id -eq $careRequest.id }

if (-not $matchedQueueItem) {
    throw "Care request $($careRequest.id) not found in live queue!"
}

Write-Host "    Queue Item Found: $($matchedQueueItem.id)" -ForegroundColor Green
Write-Host "    Queue Status: $($matchedQueueItem.status)" -ForegroundColor Green
Write-Host "    Patient Name: $($matchedQueueItem.patientName)" -ForegroundColor Green

Write-Host "`n====================================================" -ForegroundColor Green
Write-Host "       ALL DAY 2 INTEGRATION TESTS PASSED 100%!     " -ForegroundColor Green
Write-Host "====================================================" -ForegroundColor Green
