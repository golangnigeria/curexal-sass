import React, { useState } from "react";
import { CapabilityGate } from "@/components/design-system/capability-gate";
import { DataTable } from "@/components/design-system/data-table";
import { StatusPill } from "@/components/design-system/status-pill";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import {
  Stethoscope,
  Users,
  Clock,
  HeartPulse,
  FileText,
  Plus,
  Play,
  CheckCircle2,
} from "lucide-react";

interface ConsultationPatient {
  id: string;
  queueNo: number;
  patientName: string;
  patientMrn: string;
  ageGender: string;
  triageCategory: "emergency" | "urgent" | "routine";
  vitals: { bp: string; pulse: string; temp: string; weight: string };
  chiefComplaint: string;
  status: "waiting" | "in_consultation" | "completed";
  waitingSince: string;
}

const mockQueue: ConsultationPatient[] = [
  { id: "q-1", queueNo: 1, patientName: "Amina Yusuf", patientMrn: "PAT-0012", ageGender: "32F", triageCategory: "urgent", vitals: { bp: "140/90", pulse: "88 bpm", temp: "38.2°C", weight: "68 kg" }, chiefComplaint: "High fever, chills and severe body ache for 3 days.", status: "in_consultation", waitingSince: "12 mins ago" },
  { id: "q-2", queueNo: 2, patientName: "Chinedu Okafor", patientMrn: "PAT-0034", ageGender: "45M", triageCategory: "routine", vitals: { bp: "120/80", pulse: "72 bpm", temp: "36.8°C", weight: "82 kg" }, chiefComplaint: "Routine hypertension follow-up & prescription renewal.", status: "waiting", waitingSince: "20 mins ago" },
  { id: "q-3", queueNo: 3, patientName: "Babatunde Lawal", patientMrn: "PAT-0078", ageGender: "58M", triageCategory: "emergency", vitals: { bp: "170/110", pulse: "102 bpm", temp: "37.1°C", weight: "90 kg" }, chiefComplaint: "Sudden onset chest tightness and shortness of breath.", status: "waiting", waitingSince: "5 mins ago" },
  { id: "q-4", queueNo: 4, patientName: "Fatima Bello", patientMrn: "PAT-0091", ageGender: "26F", triageCategory: "routine", vitals: { bp: "115/75", pulse: "68 bpm", temp: "36.6°C", weight: "55 kg" }, chiefComplaint: "Persistent dry cough for 2 weeks.", status: "completed", waitingSince: "1 hour ago" },
];

export default function WorkspaceClinicalPage() {
  const [queue, setQueue] = useState<ConsultationPatient[]>(mockQueue);
  const [activePatient, setActivePatient] = useState<ConsultationPatient | null>(mockQueue[0]);

  const handleStartConsult = (p: ConsultationPatient) => {
    setActivePatient(p);
    setQueue((prev) =>
      prev.map((item) => (item.id === p.id ? { ...item, status: "in_consultation" } : item))
    );
    toast.info(`Started Consultation for ${p.patientName}`);
  };

  return (
    <CapabilityGate
      capability="clinical.basic"
      moduleCode="clinical"
      title="Outpatient Clinical & EMR Suite"
      description="Unlock doctor consultations, OPD triage queues, electronic SOAP notes, and formulary prescribing."
      requiredPlan="Optimize or Pro"
    >
      <div className="space-y-8 animate-fade-in">
        {/* Header */}
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-6">
          <div>
            <div className="flex items-center gap-2.5 mb-1">
              <h1 className="text-2xl font-bold tracking-tight text-foreground flex items-center gap-2">
                <Stethoscope className="w-6 h-6 text-sky-600 dark:text-sky-400" />
                Outpatient Clinic & EMR
              </h1>
              <Badge variant="outline" className="border-sky-500/40 text-sky-600 dark:text-sky-400 bg-sky-500/10 text-[10px] font-mono">
                Live OPD Triage
              </Badge>
            </div>
            <p className="text-xs text-muted-foreground">
              Doctor consultation queue, clinical triage vitals, and electronic medical record management.
            </p>
          </div>

          <div className="flex items-center gap-2.5">
            <Button size="sm" className="text-xs h-9 gap-1.5 bg-sky-600 hover:bg-sky-700 text-white shadow">
              <Plus className="w-3.5 h-3.5" />
              Triage New Patient
            </Button>
          </div>
        </div>

        {/* Queue Table */}
        <DataTable
          data={queue}
          searchPlaceholder="Search queue, MRN, patient name..."
          searchKey="patientName"
          columns={[
            {
              header: "Queue #",
              cell: (p) => (
                <div className="font-mono font-bold text-xs bg-secondary w-7 h-7 rounded-full flex items-center justify-center">
                  {p.queueNo}
                </div>
              ),
            },
            {
              header: "Patient",
              cell: (p) => (
                <div>
                  <p className="font-semibold text-foreground">{p.patientName}</p>
                  <p className="text-[11px] font-mono text-muted-foreground">{p.patientMrn} • {p.ageGender}</p>
                </div>
              ),
            },
            {
              header: "Triage Vitals",
              cell: (p) => (
                <div className="text-[11px] font-mono space-y-0.5">
                  <p><span className="text-muted-foreground">BP:</span> <span className="font-bold">{p.vitals.bp}</span> | <span className="text-muted-foreground">Pulse:</span> {p.vitals.pulse}</p>
                  <p><span className="text-muted-foreground">Temp:</span> {p.vitals.temp} | <span className="text-muted-foreground">Wt:</span> {p.vitals.weight}</p>
                </div>
              ),
            },
            {
              header: "Chief Complaint",
              cell: (p) => (
                <p className="text-xs text-foreground max-w-xs truncate">{p.chiefComplaint}</p>
              ),
            },
            {
              header: "Status",
              cell: (p) => <StatusPill status={p.status} />,
            },
            {
              header: "Actions",
              className: "text-right",
              cell: (p) => (
                <div className="flex items-center justify-end gap-2">
                  {p.status === "waiting" && (
                    <Button
                      size="sm"
                      onClick={() => handleStartConsult(p)}
                      className="text-xs h-7 gap-1 bg-sky-600 hover:bg-sky-700 text-white"
                    >
                      <Play className="w-3 h-3" /> Call Patient
                    </Button>
                  )}
                  {p.status === "in_consultation" && (
                    <Badge className="bg-sky-500/10 text-sky-600 dark:text-sky-400 border-sky-500/30 text-[10px]">
                      In Room 2
                    </Badge>
                  )}
                </div>
              ),
            },
          ]}
        />
      </div>
    </CapabilityGate>
  );
}
