import React, { useState } from "react";
import { CapabilityGate } from "@/components/design-system/capability-gate";
import { DataTable } from "@/components/design-system/data-table";
import { StatusPill } from "@/components/design-system/status-pill";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import {
  Microscope,
  Barcode,
  CheckCircle2,
  AlertTriangle,
  Play,
  FileCheck,
  Cpu,
  Plus,
  Search,
  Filter,
} from "lucide-react";

interface SampleItem {
  id: string;
  barcode: string;
  patientName: string;
  patientMrn: string;
  testName: string;
  specimenType: string;
  status: "collected" | "accessioned" | "in_analysis" | "results_pending" | "authorized";
  priority: "routine" | "stat" | "urgent";
  receivedAt: string;
}

const mockSamples: SampleItem[] = [
  { id: "s-1", barcode: "BC-2026-8401", patientName: "Amina Yusuf", patientMrn: "PAT-0012", testName: "Complete Blood Count (CBC)", specimenType: "Whole Blood (EDTA)", status: "results_pending", priority: "stat", receivedAt: "10 mins ago" },
  { id: "s-2", barcode: "BC-2026-8402", patientName: "Chinedu Okafor", patientMrn: "PAT-0034", testName: "Fasting Blood Sugar (FBS)", specimenType: "Fluoride Plasma", status: "in_analysis", priority: "routine", receivedAt: "25 mins ago" },
  { id: "s-3", barcode: "BC-2026-8403", patientName: "Babatunde Lawal", patientMrn: "PAT-0078", testName: "Lipid Profile Panel", specimenType: "Serum", status: "accessioned", priority: "routine", receivedAt: "40 mins ago" },
  { id: "s-4", barcode: "BC-2026-8404", patientName: "Fatima Bello", patientMrn: "PAT-0091", testName: "Liver Function Test (LFT)", specimenType: "Serum", status: "authorized", priority: "routine", receivedAt: "1 hour ago" },
  { id: "s-5", barcode: "BC-2026-8405", patientName: "Emeka Eze", patientMrn: "PAT-0105", testName: "Hepatitis B Surface Antigen (HBsAg)", specimenType: "Serum", status: "collected", priority: "urgent", receivedAt: "2 hours ago" },
];

export default function WorkspaceLaboratoryPage() {
  const [samples, setSamples] = useState<SampleItem[]>(mockSamples);
  const [selectedSample, setSelectedSample] = useState<SampleItem | null>(null);

  const handleAuthorize = (sampleId: string) => {
    setSamples((prev) =>
      prev.map((s) => (s.id === sampleId ? { ...s, status: "authorized" } : s))
    );
    toast.success("Laboratory Results Electronically Authorized!", {
      description: "Tamper-evident verification stamp generated.",
    });
  };

  return (
    <CapabilityGate
      capability="laboratory.basic"
      moduleCode="laboratory"
      title="Medical Laboratory (LIMS) Module"
      description="Unlock full specimen accessioning, barcode label printing, automated analyzer interfacing, and two-step verification."
      requiredPlan="Optimize or Pro"
    >
      <div className="space-y-8 animate-fade-in">
        {/* Header */}
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-6">
          <div>
            <div className="flex items-center gap-2.5 mb-1">
              <h1 className="text-2xl font-bold tracking-tight text-foreground flex items-center gap-2">
                <Microscope className="w-6 h-6 text-teal-600 dark:text-teal-400" />
                Laboratory Information System (LIS)
              </h1>
              <Badge variant="outline" className="border-teal-500/40 text-teal-600 dark:text-teal-400 bg-teal-500/10 text-[10px] font-mono">
                Accessioning & QC
              </Badge>
            </div>
            <p className="text-xs text-muted-foreground">
              Specimen reception, analyzer worklists, delta checking, and consultant pathologist sign-off.
            </p>
          </div>

          <div className="flex items-center gap-2.5">
            <Button size="sm" variant="outline" className="text-xs h-9 gap-1.5">
              <Barcode className="w-3.5 h-3.5" />
              Scan Barcode (F2)
            </Button>
            <Button size="sm" className="text-xs h-9 gap-1.5 bg-teal-600 hover:bg-teal-700 text-white shadow">
              <Plus className="w-3.5 h-3.5" />
              Accession Specimen
            </Button>
          </div>
        </div>

        {/* LIS Worklist Table */}
        <DataTable
          data={samples}
          searchPlaceholder="Search barcode, patient name, or test..."
          searchKey="barcode"
          columns={[
            {
              header: "Sample Barcode",
              cell: (s) => (
                <div className="flex items-center gap-2">
                  <div className="p-1.5 rounded bg-secondary text-foreground font-mono font-bold text-xs">
                    {s.barcode}
                  </div>
                  {s.priority === "stat" && (
                    <Badge className="bg-rose-500 text-white text-[9px] font-mono">STAT</Badge>
                  )}
                </div>
              ),
            },
            {
              header: "Patient Details",
              cell: (s) => (
                <div>
                  <p className="font-semibold text-foreground">{s.patientName}</p>
                  <p className="text-[11px] font-mono text-muted-foreground">{s.patientMrn}</p>
                </div>
              ),
            },
            {
              header: "Investigation & Specimen",
              cell: (s) => (
                <div>
                  <p className="font-medium text-foreground">{s.testName}</p>
                  <p className="text-[10px] text-muted-foreground">{s.specimenType}</p>
                </div>
              ),
            },
            {
              header: "Stage / Status",
              cell: (s) => <StatusPill status={s.status} />,
            },
            {
              header: "Actions",
              className: "text-right",
              cell: (s) => (
                <div className="flex items-center justify-end gap-2">
                  {s.status === "results_pending" && (
                    <Button
                      size="sm"
                      onClick={() => handleAuthorize(s.id)}
                      className="text-xs h-7 gap-1 bg-teal-600 hover:bg-teal-700 text-white"
                    >
                      <CheckCircle2 className="w-3 h-3" /> Authorize
                    </Button>
                  )}
                  {s.status === "authorized" && (
                    <span className="text-[11px] text-teal-600 dark:text-teal-400 font-semibold flex items-center gap-1">
                      <FileCheck className="w-3.5 h-3.5" /> Verified
                    </span>
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
