import React, { useState } from "react";
import { CapabilityGate } from "@/components/design-system/capability-gate";
import { DataTable } from "@/components/data-display/data-table";
import { StatusBadge } from "@/components/feedback/status-badge";
import { PatientHeader } from "@/features/patients/patient-header";
import { PatientDrawer } from "@/features/patients/patient-drawer";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import { usePermissions } from "@/api/hooks/use-permissions";
import {
  Stethoscope,
  Users,
  Clock,
  HeartPulse,
  FileText,
  Plus,
  Play,
  CheckCircle2,
  Send,
  Microscope,
  Pill,
  Layers,
  ChevronRight,
  Sparkles,
  ShieldCheck,
  Activity,
  UserCheck,
  TrendingUp,
  AlertCircle,
  Building2,
} from "lucide-react";

interface ConsultationPatient {
  id: string;
  queueNo: number;
  firstName: string;
  lastName: string;
  mrn: string;
  gender: "MALE" | "FEMALE" | string;
  dateOfBirth: string;
  bloodGroup: string;
  genotype: string;
  triageCategory: "emergency" | "urgent" | "routine";
  vitals: { bp: string; pulse: string; temp: string; weight: string; spo2: string };
  chiefComplaint: string;
  status: "waiting" | "in_consultation" | "completed";
  waitingSince: string;
}

const mockQueue: ConsultationPatient[] = [
  {
    id: "pat-1",
    queueNo: 1,
    firstName: "Amina",
    lastName: "Yusuf",
    mrn: "PAT-0012",
    gender: "FEMALE",
    dateOfBirth: "1994-05-14",
    bloodGroup: "O+",
    genotype: "AA",
    triageCategory: "urgent",
    vitals: { bp: "140/90 mmHg", pulse: "88 bpm", temp: "38.2°C", weight: "68 kg", spo2: "98%" },
    chiefComplaint: "High fever, persistent chills and severe myalgia for 3 days.",
    status: "in_consultation",
    waitingSince: "12 mins ago",
  },
  {
    id: "pat-2",
    queueNo: 2,
    firstName: "Chinedu",
    lastName: "Okafor",
    mrn: "PAT-0034",
    gender: "MALE",
    dateOfBirth: "1981-11-20",
    bloodGroup: "A+",
    genotype: "AS",
    triageCategory: "routine",
    vitals: { bp: "120/80 mmHg", pulse: "72 bpm", temp: "36.8°C", weight: "82 kg", spo2: "99%" },
    chiefComplaint: "Routine hypertension follow-up and prescription renewal.",
    status: "waiting",
    waitingSince: "20 mins ago",
  },
  {
    id: "pat-3",
    queueNo: 3,
    firstName: "Babatunde",
    lastName: "Lawal",
    mrn: "PAT-0078",
    gender: "MALE",
    dateOfBirth: "1968-03-09",
    bloodGroup: "B+",
    genotype: "AA",
    triageCategory: "emergency",
    vitals: { bp: "175/110 mmHg", pulse: "104 bpm", temp: "37.1°C", weight: "90 kg", spo2: "94%" },
    chiefComplaint: "Sudden onset severe chest tightness and exertional dyspnea.",
    status: "waiting",
    waitingSince: "5 mins ago",
  },
  {
    id: "pat-4",
    queueNo: 4,
    firstName: "Fatima",
    lastName: "Bello",
    mrn: "PAT-0091",
    gender: "FEMALE",
    dateOfBirth: "1998-08-25",
    bloodGroup: "O+",
    genotype: "AA",
    triageCategory: "routine",
    vitals: { bp: "115/75 mmHg", pulse: "68 bpm", temp: "36.6°C", weight: "55 kg", spo2: "100%" },
    chiefComplaint: "Persistent dry cough for 2 weeks with nocturnal wheezing.",
    status: "completed",
    waitingSince: "1 hour ago",
  },
];

export default function WorkspaceClinicalPage() {
  const { can, isClinician, isExecutive, isNurse, isManager, role } = usePermissions();

  const [queue, setQueue] = useState<ConsultationPatient[]>(mockQueue);
  const [activePatient, setActivePatient] = useState<ConsultationPatient | null>(mockQueue[0]);
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);

  // SOAP State for Doctors
  const [subjective, setSubjective] = useState(
    "Patient presents with a 3-day history of high intermittent fever, chills, generalized body weakness, and headache. Denies persistent vomiting or diarrhea."
  );
  const [objective, setObjective] = useState(
    "Febrile to touch (38.2°C). Mild pallor, no icterus. Chest is clear with vesicular breath sounds bilaterally. Abdomen soft, non-tender, no organomegaly."
  );
  const [assessment, setAssessment] = useState(
    "1. Acute Febrile Illness - Suspected Malaria (Severe/Complicated rule-out)\n2. Mild Dehydration secondary to pyrexia"
  );
  const [plan, setPlan] = useState(
    "1. Urgent Full Blood Count (FBC) + Malaria Parasite Giemsa Thick Film (LIS)\n2. Oral Rehydration Solution (ORS) 1L + Paracetamol 1g PO stat\n3. Prescribe Artemether-Lumefantrine pending smear confirmation"
  );

  const [orderLabs, setOrderLabs] = useState(true);
  const [orderMeds, setOrderMeds] = useState(true);
  const [orderRad, setOrderRad] = useState(false);

  const handleStartConsult = (p: ConsultationPatient) => {
    setActivePatient(p);
    setQueue((prev) =>
      prev.map((item) => (item.id === p.id ? { ...item, status: "in_consultation" } : item))
    );
    toast.info(`Called patient: ${p.firstName} ${p.lastName} into Consultation Room 1`);
  };

  const handleCompleteEncounter = () => {
    if (!activePatient) return;
    setQueue((prev) =>
      prev.map((item) => (item.id === activePatient.id ? { ...item, status: "completed" } : item))
    );
    toast.success(`Encounter Completed & Orders Dispatched!`, {
      description: `Dispatched Lab Orders to LIS queue and Prescriptions to Pharmacy Dispensary.`,
    });
    const next = queue.find((p) => p.status === "waiting");
    setActivePatient(next || null);
  };

  // Determine Persona View
  const showExecutiveView = isExecutive || (isManager && !isClinician);

  return (
    <CapabilityGate
      capability="clinical.basic"
      moduleCode="clinical"
      title="Outpatient Clinical & EMR Suite"
      description="Doctor consultation queue, electronic SOAP encounters, diagnostic work orders, and digital prescribing."
      requiredPlan="Optimize or Pro"
    >
      <div className="space-y-6 animate-in fade-in duration-300">
        {/* =========================================================================
            PERSONA VIEW 1: EXECUTIVE & OPERATIONS MANAGEMENT VIEW
           ========================================================================= */}
        {showExecutiveView ? (
          <div className="space-y-6">
            {/* Header */}
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-5">
              <div>
                <div className="flex items-center gap-2.5 mb-1">
                  <h1 className="text-2xl font-bold tracking-tight text-foreground flex items-center gap-2">
                    <Activity className="w-6 h-6 text-teal-600 dark:text-teal-400" />
                    Clinical Department Operations & Quality Overview
                  </h1>
                  <Badge variant="outline" className="border-primary/40 text-primary bg-primary/10 text-[10px] font-mono">
                    {isExecutive ? "Executive Supervision" : "Facility Management"}
                  </Badge>
                </div>
                <p className="text-xs text-muted-foreground">
                  Real-time departmental patient census, doctor utilization, clinical triage throughput, and diagnostic order turnaround.
                </p>
              </div>

              <div className="flex items-center gap-2.5">
                <Badge variant="secondary" className="gap-1 text-xs py-1 px-3">
                  <ShieldCheck className="w-3.5 h-3.5 text-emerald-500" />
                  EMR Documentation: 100% Compliant
                </Badge>
              </div>
            </div>

            {/* Department Operational Metrics */}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
              <Card className="border-border shadow-sm bg-card">
                <CardContent className="p-4 flex items-center gap-3.5">
                  <div className="w-10 h-10 rounded-xl bg-teal-500/10 text-teal-600 flex items-center justify-center shrink-0">
                    <Users className="w-5 h-5" />
                  </div>
                  <div>
                    <p className="text-[11px] font-medium text-muted-foreground">Total Patient Census</p>
                    <h3 className="text-xl font-bold text-foreground">24 Patients</h3>
                    <p className="text-[10px] text-emerald-600 font-medium mt-0.5">8 In Consult | 16 Waiting</p>
                  </div>
                </CardContent>
              </Card>

              <Card className="border-border shadow-sm bg-card">
                <CardContent className="p-4 flex items-center gap-3.5">
                  <div className="w-10 h-10 rounded-xl bg-sky-500/10 text-sky-600 flex items-center justify-center shrink-0">
                    <Clock className="w-5 h-5" />
                  </div>
                  <div>
                    <p className="text-[11px] font-medium text-muted-foreground">Avg Consultation Time</p>
                    <h3 className="text-xl font-bold text-foreground">14.2 mins</h3>
                    <p className="text-[10px] text-muted-foreground mt-0.5">Target: &lt; 20 mins</p>
                  </div>
                </CardContent>
              </Card>

              <Card className="border-border shadow-sm bg-card">
                <CardContent className="p-4 flex items-center gap-3.5">
                  <div className="w-10 h-10 rounded-xl bg-amber-500/10 text-amber-600 flex items-center justify-center shrink-0">
                    <UserCheck className="w-5 h-5" />
                  </div>
                  <div>
                    <p className="text-[11px] font-medium text-muted-foreground">Attending Doctors Active</p>
                    <h3 className="text-xl font-bold text-foreground">3 Physicians</h3>
                    <p className="text-[10px] text-muted-foreground mt-0.5">Rooms 1, 2, 3 Active</p>
                  </div>
                </CardContent>
              </Card>

              <Card className="border-border shadow-sm bg-card">
                <CardContent className="p-4 flex items-center gap-3.5">
                  <div className="w-10 h-10 rounded-xl bg-purple-500/10 text-purple-600 flex items-center justify-center shrink-0">
                    <TrendingUp className="w-5 h-5" />
                  </div>
                  <div>
                    <p className="text-[11px] font-medium text-muted-foreground">Diagnostic Dispatches</p>
                    <h3 className="text-xl font-bold text-foreground">42 Orders Today</h3>
                    <p className="text-[10px] text-muted-foreground mt-0.5">31 LIS | 11 Pharmacy</p>
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* Notice Banner */}
            <div className="p-4 rounded-xl border border-border bg-secondary/20 flex items-start gap-3">
              <AlertCircle className="w-5 h-5 text-primary shrink-0 mt-0.5" />
              <div className="text-xs space-y-1">
                <p className="font-semibold text-foreground">
                  Executive Supervision Mode Active
                </p>
                <p className="text-muted-foreground leading-relaxed">
                  Individual patient clinical documentation (SOAP notes, diagnoses, and medical prescribing) is restricted to credentialed Attending Physicians on duty. Executive metrics reflect aggregated department performance, bed census, and clinical turnaround times.
                </p>
              </div>
            </div>

            {/* Department Active Queue & Doctor Roster Grid */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              {/* Active Clinical Patient Flow */}
              <div className="lg:col-span-2 space-y-4">
                <Card className="border-border shadow-sm bg-card">
                  <CardHeader className="pb-3 border-b border-border flex flex-row items-center justify-between">
                    <div>
                      <CardTitle className="text-sm font-bold text-foreground">
                        Department Patient Queue & Triage Acuity
                      </CardTitle>
                      <CardDescription className="text-xs">
                        Real-time status of patients currently in triage, consultation, and diagnostic transit.
                      </CardDescription>
                    </div>
                    <Badge variant="outline" className="text-[10px] font-mono">
                      Live Queue ({queue.length})
                    </Badge>
                  </CardHeader>
                  <CardContent className="p-0">
                    <div className="divide-y divide-border">
                      {queue.map((p) => (
                        <div key={p.id} className="p-4 flex items-center justify-between hover:bg-secondary/10 transition-colors">
                          <div className="flex items-center gap-3">
                            <div className="w-8 h-8 rounded-full bg-secondary flex items-center justify-center text-xs font-bold font-mono">
                              #{p.queueNo}
                            </div>
                            <div>
                              <p className="text-xs font-bold text-foreground flex items-center gap-2">
                                {p.firstName} {p.lastName}
                                <span className="text-[10px] font-normal text-muted-foreground font-mono">
                                  ({p.mrn})
                                </span>
                              </p>
                              <p className="text-[11px] text-muted-foreground">
                                {p.gender} • {p.bloodGroup} • {p.chiefComplaint}
                              </p>
                            </div>
                          </div>

                          <div className="flex items-center gap-2.5">
                            <Badge
                              variant="outline"
                              className={`text-[10px] uppercase font-mono ${
                                p.triageCategory === "emergency"
                                  ? "border-rose-500/40 text-rose-600 bg-rose-500/10"
                                  : p.triageCategory === "urgent"
                                  ? "border-amber-500/40 text-amber-600 bg-amber-500/10"
                                  : "border-teal-500/40 text-teal-600 bg-teal-500/10"
                              }`}
                            >
                              {p.triageCategory}
                            </Badge>
                            <StatusBadge status={p.status} />
                          </div>
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              </div>

              {/* On-Duty Physician Roster */}
              <div className="space-y-4">
                <Card className="border-border shadow-sm bg-card">
                  <CardHeader className="pb-3 border-b border-border">
                    <CardTitle className="text-sm font-bold text-foreground flex items-center gap-2">
                      <Stethoscope className="w-4 h-4 text-teal-600" />
                      Physicians On Duty Roster
                    </CardTitle>
                    <CardDescription className="text-xs">
                      Consultation room allocation & live practitioner status.
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="pt-4 space-y-3.5">
                    <div className="p-3 rounded-lg border border-border bg-secondary/10 flex items-center justify-between">
                      <div>
                        <p className="text-xs font-bold text-foreground">Dr. Emeka Nwosu</p>
                        <p className="text-[11px] text-muted-foreground">Chief Medical Officer • Room 1</p>
                      </div>
                      <Badge className="bg-emerald-500/10 text-emerald-600 border-emerald-500/30 text-[10px]">
                        In Consult
                      </Badge>
                    </div>

                    <div className="p-3 rounded-lg border border-border bg-secondary/10 flex items-center justify-between">
                      <div>
                        <p className="text-xs font-bold text-foreground">Dr. Sarah Adebayo</p>
                        <p className="text-[11px] text-muted-foreground">Consultant Specialist • Room 2</p>
                      </div>
                      <Badge className="bg-emerald-500/10 text-emerald-600 border-emerald-500/30 text-[10px]">
                        In Consult
                      </Badge>
                    </div>

                    <div className="p-3 rounded-lg border border-border bg-secondary/10 flex items-center justify-between">
                      <div>
                        <p className="text-xs font-bold text-foreground">Dr. Chinedu Eke</p>
                        <p className="text-[11px] text-muted-foreground">Primary Care Physician • Room 3</p>
                      </div>
                      <Badge variant="outline" className="text-[10px] text-muted-foreground">
                        Available
                      </Badge>
                    </div>
                  </CardContent>
                </Card>
              </div>
            </div>
          </div>
        ) : (
          /* =========================================================================
             PERSONA VIEW 2: DOCTOR PRACTICE MODE (ATTENDING PHYSICIAN)
             ========================================================================= */
          <div className="space-y-6">
            {/* Top Header */}
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-5">
              <div>
                <div className="flex items-center gap-2.5 mb-1">
                  <h1 className="text-2xl font-bold tracking-tight text-foreground flex items-center gap-2">
                    <Stethoscope className="w-6 h-6 text-teal-600 dark:text-teal-400" />
                    Outpatient Clinic & EMR Practice Canvas
                  </h1>
                  <Badge variant="outline" className="border-teal-500/40 text-teal-600 dark:text-teal-400 bg-teal-500/10 text-[10px] font-mono">
                    Consultation Room 1
                  </Badge>
                </div>
                <p className="text-xs text-muted-foreground">
                  Doctor consultation queue, clinical triage vitals, and electronic medical record management.
                </p>
              </div>

              <div className="flex items-center gap-2.5">
                <Button
                  size="sm"
                  onClick={() => {
                    toast.info("Triage Intake modal opened");
                  }}
                  className="text-xs h-8 gap-1.5 bg-primary text-primary-foreground shadow"
                >
                  <Plus className="w-3.5 h-3.5" />
                  Triage Intake
                </Button>
              </div>
            </div>

            {/* Active Patient Clinical Encounter Workspace */}
            {activePatient && (
              <div className="space-y-4">
                <PatientHeader
                  patient={{
                    firstName: activePatient.firstName,
                    lastName: activePatient.lastName,
                    mrn: activePatient.mrn,
                    gender: activePatient.gender,
                    dateOfBirth: activePatient.dateOfBirth,
                    bloodGroup: activePatient.bloodGroup,
                    genotype: activePatient.genotype,
                    activeEncounterId: activePatient.id,
                  }}
                  onOpenDrawer={() => setIsDrawerOpen(true)}
                />

                {/* SOAP Note & Diagnostics Grid */}
                <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                  {/* SOAP Documentation Canvas */}
                  <div className="lg:col-span-2 space-y-4">
                    <Card className="border-border shadow-sm bg-card">
                      <CardHeader className="pb-3 border-b border-border">
                        <div className="flex items-center justify-between">
                          <CardTitle className="text-sm font-bold text-foreground flex items-center gap-2">
                            <FileText className="w-4 h-4 text-primary" />
                            Electronic Clinical Encounter Note (SOAP)
                          </CardTitle>
                          <span className="text-[10px] font-mono text-muted-foreground">
                            ICD-10 / SNOMED CT Auto-Coding Active
                          </span>
                        </div>
                      </CardHeader>
                      <CardContent className="pt-4 space-y-3.5">
                        {/* Subjective */}
                        <div>
                          <label className="text-xs font-semibold text-foreground flex items-center justify-between mb-1">
                            <span>S - Subjective (History of Presenting Complaint)</span>
                          </label>
                          <textarea
                            rows={2}
                            value={subjective}
                            onChange={(e) => setSubjective(e.target.value)}
                            className="w-full text-xs p-2.5 rounded-lg border border-border bg-secondary/20 focus:bg-background transition-colors focus:outline-none focus:ring-1 focus:ring-primary"
                          />
                        </div>

                        {/* Objective */}
                        <div>
                          <label className="text-xs font-semibold text-foreground flex items-center justify-between mb-1">
                            <span>O - Objective (Physical Findings & Vitals)</span>
                            <span className="text-[11px] font-mono text-muted-foreground">
                              {activePatient.vitals.bp} | {activePatient.vitals.temp}
                            </span>
                          </label>
                          <textarea
                            rows={2}
                            value={objective}
                            onChange={(e) => setObjective(e.target.value)}
                            className="w-full text-xs p-2.5 rounded-lg border border-border bg-secondary/20 focus:bg-background transition-colors focus:outline-none focus:ring-1 focus:ring-primary"
                          />
                        </div>

                        {/* Assessment */}
                        <div>
                          <label className="text-xs font-semibold text-foreground mb-1 block">
                            A - Clinical Assessment & Working Diagnosis
                          </label>
                          <textarea
                            rows={2}
                            value={assessment}
                            onChange={(e) => setAssessment(e.target.value)}
                            className="w-full text-xs p-2.5 rounded-lg border border-border bg-secondary/20 focus:bg-background transition-colors focus:outline-none focus:ring-1 focus:ring-primary font-medium"
                          />
                        </div>

                        {/* Plan */}
                        <div>
                          <label className="text-xs font-semibold text-foreground mb-1 block">
                            P - Management Plan, Interventions & Patient Counseling
                          </label>
                          <textarea
                            rows={2}
                            value={plan}
                            onChange={(e) => setPlan(e.target.value)}
                            className="w-full text-xs p-2.5 rounded-lg border border-border bg-secondary/20 focus:bg-background transition-colors focus:outline-none focus:ring-1 focus:ring-primary"
                          />
                        </div>
                      </CardContent>
                    </Card>
                  </div>

                  {/* One-Click Orders & Direct Dispatch */}
                  <div className="space-y-4">
                    <Card className="border-border shadow-sm bg-card">
                      <CardHeader className="pb-3 border-b border-border">
                        <CardTitle className="text-sm font-bold text-foreground">
                          Diagnostic Work Orders
                        </CardTitle>
                        <CardDescription className="text-xs">
                          Instantly queue orders into laboratory, radiology, and pharmacy.
                        </CardDescription>
                      </CardHeader>
                      <CardContent className="pt-4 space-y-3">
                        <label className="flex items-start gap-2.5 p-2.5 rounded-lg border border-border bg-secondary/20 hover:bg-secondary/40 cursor-pointer transition-colors">
                          <input
                            type="checkbox"
                            checked={orderLabs}
                            onChange={(e) => setOrderLabs(e.target.checked)}
                            className="mt-0.5 rounded text-primary"
                          />
                          <div className="text-xs">
                            <p className="font-semibold text-foreground flex items-center gap-1.5">
                              <Microscope className="w-3.5 h-3.5 text-sky-600" />
                              Lab: CBC + Malaria Smear
                            </p>
                            <p className="text-[11px] text-muted-foreground">Auto-generates LIS accession & barcode</p>
                          </div>
                        </label>

                        <label className="flex items-start gap-2.5 p-2.5 rounded-lg border border-border bg-secondary/20 hover:bg-secondary/40 cursor-pointer transition-colors">
                          <input
                            type="checkbox"
                            checked={orderMeds}
                            onChange={(e) => setOrderMeds(e.target.checked)}
                            className="mt-0.5 rounded text-primary"
                          />
                          <div className="text-xs">
                            <p className="font-semibold text-foreground flex items-center gap-1.5">
                              <Pill className="w-3.5 h-3.5 text-emerald-600" />
                              Rx: Artemether-Lumefantrine + Paracetamol
                            </p>
                            <p className="text-[11px] text-muted-foreground">Dispatches to Pharmacy dispensary queue</p>
                          </div>
                        </label>

                        <label className="flex items-start gap-2.5 p-2.5 rounded-lg border border-border bg-secondary/20 hover:bg-secondary/40 cursor-pointer transition-colors">
                          <input
                            type="checkbox"
                            checked={orderRad}
                            onChange={(e) => setOrderRad(e.target.checked)}
                            className="mt-0.5 rounded text-primary"
                          />
                          <div className="text-xs">
                            <p className="font-semibold text-foreground flex items-center gap-1.5">
                              <Layers className="w-3.5 h-3.5 text-indigo-600" />
                              Imaging: Chest X-Ray PA
                            </p>
                            <p className="text-[11px] text-muted-foreground">Dispatches to RIS / PACS modality worklist</p>
                          </div>
                        </label>

                        <Button
                          onClick={handleCompleteEncounter}
                          className="w-full text-xs h-9 gap-1.5 bg-emerald-600 hover:bg-emerald-700 text-white shadow font-semibold mt-2"
                        >
                          <Send className="w-3.5 h-3.5" />
                          Sign SOAP & Dispatch Orders
                        </Button>
                      </CardContent>
                    </Card>
                  </div>
                </div>
              </div>
            )}

            {/* Waiting Room & Queue List */}
            <Card className="border-border shadow-sm bg-card">
              <CardHeader className="pb-3 border-b border-border flex flex-row items-center justify-between">
                <div>
                  <CardTitle className="text-sm font-bold text-foreground">
                    Attending Physician Consultation Queue
                  </CardTitle>
                  <CardDescription className="text-xs">
                    Triaged patients queued for doctor review and medical assessment.
                  </CardDescription>
                </div>
                <Badge variant="outline" className="text-[10px] font-mono">
                  {queue.filter((p) => p.status === "waiting").length} Waiting
                </Badge>
              </CardHeader>
              <CardContent className="p-0">
                <div className="divide-y divide-border">
                  {queue.map((p) => (
                    <div
                      key={p.id}
                      className={`p-4 flex items-center justify-between transition-colors ${
                        activePatient?.id === p.id ? "bg-primary/5" : "hover:bg-secondary/10"
                      }`}
                    >
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 rounded-full bg-secondary flex items-center justify-center text-xs font-bold font-mono">
                          #{p.queueNo}
                        </div>
                        <div>
                          <p className="text-xs font-bold text-foreground flex items-center gap-2">
                            {p.firstName} {p.lastName}
                            <span className="text-[10px] font-normal text-muted-foreground font-mono">
                              ({p.mrn})
                            </span>
                          </p>
                          <p className="text-[11px] text-muted-foreground">
                            {p.gender} • {p.bloodGroup} • {p.chiefComplaint}
                          </p>
                        </div>
                      </div>

                      <div className="flex items-center gap-3">
                        <Badge
                          variant="outline"
                          className={`text-[10px] uppercase font-mono ${
                            p.triageCategory === "emergency"
                              ? "border-rose-500/40 text-rose-600 bg-rose-500/10"
                              : p.triageCategory === "urgent"
                              ? "border-amber-500/40 text-amber-600 bg-amber-500/10"
                              : "border-teal-500/40 text-teal-600 bg-teal-500/10"
                          }`}
                        >
                          {p.triageCategory}
                        </Badge>
                        <StatusBadge status={p.status} />

                        {p.status === "waiting" && (
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => handleStartConsult(p)}
                            className="text-xs h-7 gap-1 text-primary border-primary/30 hover:bg-primary/10"
                          >
                            <Play className="w-3 h-3" />
                            Call In
                          </Button>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Patient 360 Drawer */}
            <PatientDrawer
              isOpen={isDrawerOpen}
              onClose={() => setIsDrawerOpen(false)}
              patient={
                activePatient
                  ? {
                      id: activePatient.id,
                      firstName: activePatient.firstName,
                      lastName: activePatient.lastName,
                      mrn: activePatient.mrn,
                      gender: activePatient.gender as any,
                      dateOfBirth: activePatient.dateOfBirth,
                      bloodGroup: activePatient.bloodGroup,
                      genotype: activePatient.genotype,
                    }
                  : null
              }
            />
          </div>
        )}
      </div>
    </CapabilityGate>
  );
}
