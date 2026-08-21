import React, { useState } from "react";
import { CapabilityGate } from "@/components/design-system/capability-gate";
import { DataTable } from "@/components/design-system/data-table";
import { StatusPill } from "@/components/design-system/status-pill";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import {
  Activity,
  Radio,
  Eye,
  FileCheck,
  Plus,
  Tv,
  Sparkles,
} from "lucide-react";

interface ScanItem {
  id: string;
  accessionNo: string;
  patientName: string;
  patientMrn: string;
  modality: "X-RAY" | "ULTRASOUND" | "CT" | "MRI";
  procedure: string;
  status: "scheduled" | "in_progress" | "images_acquired" | "reported";
  scheduledTime: string;
  radiologist?: string;
}

const mockScans: ScanItem[] = [
  { id: "rad-1", accessionNo: "ACC-8401", patientName: "Babatunde Lawal", patientMrn: "PAT-0078", modality: "X-RAY", procedure: "Chest PA View", status: "images_acquired", scheduledTime: "10:15 AM", radiologist: "Dr. K. Obi" },
  { id: "rad-2", accessionNo: "ACC-8402", patientName: "Amina Yusuf", patientMrn: "PAT-0012", modality: "ULTRASOUND", procedure: "Pelvic Scan (Transabdominal)", status: "scheduled", scheduledTime: "11:00 AM" },
  { id: "rad-3", accessionNo: "ACC-8403", patientName: "Chinedu Okafor", patientMrn: "PAT-0034", modality: "CT", procedure: "CT Brain without Contrast", status: "reported", scheduledTime: "09:00 AM", radiologist: "Dr. K. Obi" },
];

export default function WorkspaceRadiologyPage() {
  const [scans, setScans] = useState<ScanItem[]>(mockScans);

  const handleOpenViewer = (acc: string) => {
    toast.info(`Launching DICOM PACS Viewer for Accession ${acc}...`);
  };

  return (
    <CapabilityGate
      capability="radiology.basic"
      moduleCode="radiology"
      title="Radiology Information System (RIS) & PACS"
      description="Unlock modality worklists (X-Ray, Ultrasound, CT, MRI), DICOM PACS viewer integration, and structured radiologist reporting."
      requiredPlan="Optimize or Pro"
    >
      <div className="space-y-8 animate-fade-in">
        {/* Header */}
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-6">
          <div>
            <div className="flex items-center gap-2.5 mb-1">
              <h1 className="text-2xl font-bold tracking-tight text-foreground flex items-center gap-2">
                <Activity className="w-6 h-6 text-amber-600 dark:text-amber-400" />
                Radiology & Diagnostic Imaging (RIS)
              </h1>
              <Badge variant="outline" className="border-amber-500/40 text-amber-600 dark:text-amber-400 bg-amber-500/10 text-[10px] font-mono">
                DICOM / PACS Ready
              </Badge>
            </div>
            <p className="text-xs text-muted-foreground">
              Modality scan worklists, image routing, radiologist dictation, and verified diagnostic reports.
            </p>
          </div>

          <div className="flex items-center gap-2.5">
            <Button size="sm" className="text-xs h-9 gap-1.5 bg-amber-600 hover:bg-amber-700 text-white shadow">
              <Plus className="w-3.5 h-3.5" />
              Schedule Modality Scan
            </Button>
          </div>
        </div>

        {/* Scans Table */}
        <DataTable
          data={scans}
          searchPlaceholder="Search accession number, procedure..."
          searchKey="accessionNo"
          columns={[
            {
              header: "Accession #",
              cell: (s) => <span className="font-mono font-bold text-xs">{s.accessionNo}</span>,
            },
            {
              header: "Patient",
              cell: (s) => (
                <div>
                  <p className="font-semibold text-foreground">{s.patientName}</p>
                  <p className="text-[11px] font-mono text-muted-foreground">{s.patientMrn}</p>
                </div>
              ),
            },
            {
              header: "Modality & Procedure",
              cell: (s) => (
                <div className="flex items-center gap-2">
                  <Badge variant="outline" className="text-[10px] font-mono border-amber-500/40 text-amber-600 dark:text-amber-400">
                    {s.modality}
                  </Badge>
                  <span className="font-medium text-foreground text-xs">{s.procedure}</span>
                </div>
              ),
            },
            {
              header: "Scheduled Time",
              cell: (s) => <span className="font-mono text-xs text-muted-foreground">{s.scheduledTime}</span>,
            },
            {
              header: "Status",
              cell: (s) => <StatusPill status={s.status} />,
            },
            {
              header: "Actions",
              className: "text-right",
              cell: (s) => (
                <div className="flex items-center justify-end gap-2">
                  {s.status === "images_acquired" && (
                    <Button
                      size="sm"
                      onClick={() => handleOpenViewer(s.accessionNo)}
                      className="text-xs h-7 gap-1 bg-amber-600 hover:bg-amber-700 text-white"
                    >
                      <Tv className="w-3 h-3" /> PACS Viewer
                    </Button>
                  )}
                  {s.status === "reported" && (
                    <span className="text-[11px] text-emerald-600 dark:text-emerald-400 font-semibold flex items-center gap-1">
                      <FileCheck className="w-3.5 h-3.5" /> Signed Off
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
