Yes. But I would **not put a full production LIS, RIS, or Clinic/HMS inside Smart**. Smart should contain the **shared operational foundation** and a useful entry-level slice of each domain, while advanced workflows become composable capabilities.

The goal should be: **a Smart customer can actually run a small healthcare organization, but upgrading/add-ons unlock depth rather than fixing missing fundamentals.**

## 1. Smart — Core platform

These should be available to essentially every organization.

### Organization & administration

* Organization profile
* Facility profile
* Branch/workspace management
* Staff/user management
* Organization memberships
* Workspace memberships
* Roles and permissions
* Department management
* Service catalog
* Operating hours
* Contact information
* Organization verification
* Regulatory document submission
* Document status tracking
* Audit trail
* Activity history

### Patient/customer foundation

* Patient registration
* Patient demographic profile
* Patient ID/MRN
* Contact information
* Next of kin
* Basic patient search
* Patient status
* Patient history/timeline
* Duplicate patient detection
* Patient merge workflow
* Patient communication preferences

### Customer care

* Reception/front desk
* Patient check-in
* Queue management
* Appointment basics
* Service requests
* Referral recording
* Provider/facility directory
* Basic follow-up
* Notifications
* SMS/email/WhatsApp integration architecture

### Billing

* Service pricing
* Invoice generation
* Payment recording
* Receipts
* Outstanding balances
* Discounts
* Basic refunds
* Payment history
* Basic financial reports

### Platform-wide

* Dashboard
* Notifications
* Search
* Audit logs
* Reports
* Export
* Data access controls
* API foundation
* Import/export
* Multi-branch support

---

# 2. Smart LIS — Entry-level laboratory

Smart should have enough LIS functionality for a **small diagnostic laboratory** to operate.

### Patient & request

* Laboratory patient registration
* Laboratory order/request
* Test catalog
* Test groups/panels
* Test pricing
* Sample requirements
* Specimen type
* Collection instructions
* Basic barcode generation
* Sample collection
* Sample reception
* Sample status tracking

### Laboratory workflow

A Smart laboratory should be able to do:

**Order → Sample → Processing → Result → Verification → Report**

Features:

* Pending samples
* Received samples
* Rejected samples
* Collected samples
* Processing status
* Result entry
* Result editing
* Result verification
* Result authorization
* Reference ranges
* Units
* Abnormal flags
* Critical result flagging
* Result comments
* Basic result history
* Result report generation
* Result printing
* PDF result
* Patient result access

### Basic quality

* Sample rejection reason
* Result amendment
* Result audit trail
* Result version history
* Critical result notification
* Basic QC recording

### Smart LIS should **not** necessarily include:

* Advanced analyzer integration
* Instrument middleware
* Auto-validation rules
* Complex QC analytics
* Advanced microbiology
* Molecular workflows
* Blood bank
* Histopathology
* Cytology
* Advanced inventory
* Advanced quality management

Those become LIS capabilities/add-ons.

---

# 3. Production LIS — Advanced capabilities

When an organization needs a serious LIS, unlock the advanced capabilities.

### Laboratory management

* Advanced test catalog
* Panels/profiles
* Reflex testing
* Add-on tests
* Test cancellation
* Test rejection
* Repeat testing
* Dilution workflows
* Aliquot management
* Batch processing
* Worklists
* Department work queues

### Specimen management

* Barcode lifecycle
* Multiple specimens per order
* Aliquots
* Parent/child specimens
* Specimen routing
* Specimen storage
* Retention rules
* Chain of custody
* Transport tracking
* Rejection rules

### Result management

* Structured result entry
* Numeric results
* Text results
* Coded results
* Microbiology results
* Culture/sensitivity
* Susceptibility
* Molecular results
* Pathology results
* Histology
* Cytology
* Result templates
* Auto-calculation
* Derived results
* Delta checks
* Auto-validation
* Critical-value workflows
* Result correction
* Result amendment
* Result authorization

### Analyzer integration

* Instrument interfaces
* ASTM
* HL7
* Device connectivity
* Analyzer worklists
* Result import
* Bidirectional communication
* Instrument mapping
* QC data import
* Interface monitoring
* Interface error queues

### Advanced QC

* QC rules
* Levey-Jennings charts
* Westgard rules
* Control materials
* Lot tracking
* QC approval
* QC failures
* Corrective actions
* Calibration tracking

### LIS inventory

* Reagents
* Consumables
* Lots
* Expiry
* Minimum stock
* Reorder levels
* Supplier management
* Consumption tracking

### Advanced laboratory reporting

* Configurable report templates
* Multi-signatory reports
* Department-specific reports
* Digital signatures
* Result delivery
* Referral laboratory workflows
* External laboratory results
* Advanced laboratory analytics

---

# 4. Smart RIS — Entry-level radiology

Smart should allow a small imaging center to start.

### Core RIS

* Radiology service catalog
* Imaging procedure catalog
* Pricing
* Radiology request
* Appointment
* Scheduling
* Modality selection
* Patient queue
* Examination status
* Radiologist assignment
* Basic reporting
* Report approval
* Report generation
* PDF report
* Patient result access

### Imaging workflow

**Request → Schedule → Perform → Report → Approve → Deliver**

### Smart RIS should include

* X-ray
* Ultrasound
* CT request
* MRI request
* Basic modality management
* Basic radiologist workflow
* Report templates
* Findings
* Impression
* Report status
* Critical finding flag
* Basic report history

---

# 5. Production RIS — Advanced

For serious imaging centers:

### Scheduling

* Modality calendars
* Resource scheduling
* Radiologist scheduling
* Room scheduling
* Technician scheduling
* Appointment conflicts
* Waiting lists
* Rescheduling
* Cancellation
* No-show management

### Radiology workflow

* Worklists
* Modality worklists
* Exam protocols
* Contrast tracking
* Procedure preparation
* Technician workflow
* Radiologist work queue
* Preliminary reports
* Final reports
* Report amendments
* Second opinions
* Peer review

### PACS/DICOM

This should be an advanced capability.

* DICOM
* PACS integration
* DICOM worklist
* Image routing
* Image storage integration
* Study metadata
* Series management
* Viewer integration
* External image sharing
* Imaging archive
* Study retrieval

### Advanced reporting

* Structured reporting
* Templates
* Voice dictation
* Report macros
* Digital signatures
* Critical findings
* Report distribution
* Referrer notifications

---

# 6. Smart Clinic/HMS — Entry-level

For a small clinic, Smart should provide a lightweight clinical workflow.

### Patient flow

**Registration → Appointment → Consultation → Service → Billing → Follow-up**

### Clinical basics

* Patient profile
* Appointment
* Queue
* Encounter
* Chief complaint
* Basic history
* Vital signs
* Clinical notes
* Diagnosis
* Treatment plan
* Prescription record
* Follow-up
* Referral
* Basic clinical history

### Basic consultation

* Encounter notes
* Diagnosis
* Symptoms
* Vital signs
* Allergies
* Basic medication history
* Basic procedure recording
* Basic clinical summary

### Basic reporting

* Patient visit history
* Encounter reports
* Provider activity
* Basic revenue reports
* Patient statistics

---

# 7. Production Clinic/HMS — Advanced

For larger clinics/hospitals:

### Clinical

* Advanced encounters
* SOAP notes
* Problem lists
* Diagnoses
* Allergies
* Medication history
* Care plans
* Clinical pathways
* Orders
* Referrals
* Procedures
* Discharge summaries
* Clinical documentation
* Electronic medical records

### Appointment

* Multi-provider scheduling
* Department scheduling
* Resource scheduling
* Recurring appointments
* Waitlists
* No-show management
* Referral appointments

### Nursing

* Nursing notes
* Vital signs
* Nursing assessments
* Medication administration
* Care plans
* Observation charts
* Handover

### Inpatient

* Admission
* Bed management
* Ward management
* Transfers
* Discharge
* Discharge summary
* Inpatient billing
* Medication administration

### Advanced clinical

* Clinical decision support
* Care pathways
* Alerts
* Order sets
* Clinical protocols
* Advanced analytics

---

# 8. Pharmacy

I would give Smart **basic pharmacy integration**, not a full enterprise pharmacy system.

### Smart

* Drug catalog
* Prescription
* Dispensing
* Basic stock
* Stock adjustment
* Basic expiry tracking
* Basic sales
* Prescription history

### Advanced Pharmacy

* Purchase orders
* Suppliers
* Batch management
* FEFO
* Multi-location inventory
* Stock transfers
* Controlled-drug workflows
* Reconciliation
* Procurement
* Advanced dispensing
* Pharmacy analytics

---

# 9. Inventory

### Smart

* Item catalog
* Stock levels
* Basic stock movement
* Low-stock alerts
* Basic adjustments

### Advanced

* Warehouses
* Bins
* Lots
* Serial numbers
* Expiry
* Purchase orders
* Goods receiving
* Stock transfers
* Reorder automation
* Supplier management
* Consumption analytics

---

# 10. QMS — Very important for Africa

This can become one of Curexal's strongest differentiators, particularly for diagnostic laboratories.

### Smart QMS

* Basic document management
* SOP repository
* Audit trail
* Incident recording
* Corrective action recording

### Advanced QMS

* ISO 15189 workflows
* ISO 9001 workflows
* Document control
* Version control
* Non-conformities
* CAPA
* Internal audits
* External audits
* Risk management
* Quality indicators
* Equipment calibration
* Competency management
* Training records
* Proficiency testing
* EQA
* QC monitoring
* Management review

---

# 11. Interoperability — This should be a Curexal differentiator

Don't make interoperability an expensive afterthought.

The **foundation** should exist in Smart.

### Smart

* API
* Webhooks
* CSV import/export
* Basic external referrals
* Basic result sharing

### Advanced interoperability

* HL7
* FHIR
* DICOM
* ASTM
* LIS analyzer interfaces
* RIS/PACS integration
* Hospital integrations
* External laboratory integration
* External pharmacy integration
* Insurance integrations
* National health integrations

This is particularly important to your **"not fragmented / not siloed"** positioning.

---

# 12. Communication

### Smart

* In-app notifications
* Email
* Basic SMS
* Patient notifications
* Appointment reminders
* Result-ready notification

### Advanced

* WhatsApp
* Automated campaigns
* Two-way messaging
* Provider notifications
* Critical-result notifications
* Bulk communication
* Communication templates
* Notification workflows

---

# 13. Analytics

### Smart

* Patient count
* Visits
* Orders
* Revenue
* Payments
* Basic service statistics
* Basic operational dashboard

### Advanced

* LIS analytics
* RIS analytics
* Clinical analytics
* Revenue analytics
* Staff performance
* Branch comparison
* Quality indicators
* Turnaround time
* Test utilization
* Modality utilization
* Patient retention
* Operational forecasting

---

# 14. The architecture I recommend for your plans

Instead of making plans look like this:

> Smart = bad
> Optimize = better
> Pro = better
> Enterprise = everything

make them look like:

### SMART

**Core platform**

* Customer Care
* Patient foundation
* Billing
* Basic LIS
* Basic RIS
* Basic Clinic
* Basic Pharmacy
* Basic Inventory
* Basic QMS
* Interoperability foundation

---

### OPTIMIZE

Everything in Smart

*

**Operational depth**

* Advanced LIS
* Advanced laboratory workflow
* Advanced inventory
* Advanced QMS
* Advanced analytics
* More automation

---

### PRO

Everything in Optimize

*

**Multi-department healthcare**

* Advanced Clinic/HMS
* Advanced RIS
* Advanced Pharmacy
* Advanced interoperability
* Advanced reporting
* Advanced workflows

---

### ENTERPRISE

Everything in Pro

*

**Enterprise infrastructure**

* Full LIS
* Full RIS/PACS integration
* Full HMS
* Full Pharmacy
* Full Inventory
* Full QMS
* Advanced analytics
* Enterprise integrations
* Multi-organization structures
* Advanced API
* Dedicated infrastructure options
* Enterprise support

---

# 15. But the most important part: Add-ons

This is where your architecture becomes powerful.

A Smart organization should be able to say:

> "I don't need the entire LIS. I only need sample tracking and result management."

So they purchase:

```text
Smart
+
laboratory.samples
+
laboratory.results
```

Another says:

> "I need X-ray reporting."

```text
Smart
+
radiology.orders
+
radiology.reporting
```

Another:

> "I need both."

```text
Smart
+
laboratory.samples
+
laboratory.results
+
radiology.orders
+
radiology.reporting
```

And another:

```text
Smart
+
laboratory.*
+
radiology.*
+
clinical.*
```

without changing their organization, users, patients, workspace, billing, or identity.

### This is the model I'd build Curexal around:

**Plan = what you start with.**

**Capability = what your organization can use.**

**Add-on = what you can purchase when you need more.**

**Permission = what a particular employee can do.**

**Facility type = what kind of healthcare organization you are.**

**Organization status = whether the organization is permitted to operate.**

That gives you a much stronger African-market strategy than simply having four rigid plans.
