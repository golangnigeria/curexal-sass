import React, { useState } from "react";
import { CapabilityGate } from "@/components/design-system/capability-gate";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import {
  Building2,
  Users,
  Plus,
  HeartPulse,
  Clock,
  Bed,
} from "lucide-react";

interface BedItem {
  id: string;
  ward: string;
  bedNo: string;
  status: "occupied" | "available" | "cleaning";
  patientName?: string;
  patientMrn?: string;
  admittedDays?: number;
}

const mockBeds: BedItem[] = [
  { id: "b-1", ward: "Male Medical Ward", bedNo: "MMW-01", status: "occupied", patientName: "Chinedu Okafor", patientMrn: "PAT-0034", admittedDays: 3 },
  { id: "b-2", ward: "Male Medical Ward", bedNo: "MMW-02", status: "available" },
  { id: "b-3", ward: "Male Medical Ward", bedNo: "MMW-03", status: "occupied", patientName: "Babatunde Lawal", patientMrn: "PAT-0078", admittedDays: 1 },
  { id: "b-4", ward: "Female Surgical Ward", bedNo: "FSW-01", status: "occupied", patientName: "Amina Yusuf", patientMrn: "PAT-0012", admittedDays: 4 },
  { id: "b-5", ward: "Female Surgical Ward", bedNo: "FSW-02", status: "cleaning" },
  { id: "b-6", ward: "Female Surgical Ward", bedNo: "FSW-03", status: "available" },
];

export default function WorkspaceHospitalPage() {
  const [beds, setBeds] = useState<BedItem[]>(mockBeds);

  return (
    <CapabilityGate
      capability="clinical.inpatient_wards"
      moduleCode="clinical"
      title="Inpatient Hospital & Ward Management (HIS)"
      description="Unlock interactive ward bed boards, inpatient admissions, nursing shift handovers, and surgical scheduling."
      requiredPlan="Pro or Enterprise"
    >
      <div className="space-y-8 animate-fade-in">
        {/* Header */}
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-6">
          <div>
            <div className="flex items-center gap-2.5 mb-1">
              <h1 className="text-2xl font-bold tracking-tight text-foreground flex items-center gap-2">
                <Building2 className="w-6 h-6 text-violet-600 dark:text-violet-400" />
                Inpatient Wards & Bed Board (HIS)
              </h1>
              <Badge variant="outline" className="border-violet-500/40 text-violet-600 dark:text-violet-400 bg-violet-500/10 text-[10px] font-mono">
                Ward Management
              </Badge>
            </div>
            <p className="text-xs text-muted-foreground">
              Real-time inpatient bed allocation, admission/discharge workflows, and nursing shift handovers.
            </p>
          </div>

          <div className="flex items-center gap-2.5">
            <Button size="sm" className="text-xs h-9 gap-1.5 bg-violet-600 hover:bg-violet-700 text-white shadow">
              <Plus className="w-3.5 h-3.5" />
              Admit Inpatient
            </Button>
          </div>
        </div>

        {/* Visual Bed Grid */}
        <div className="space-y-6">
          <div className="flex items-center justify-between">
            <h3 className="text-sm font-bold text-foreground">Active Ward Beds</h3>
            <div className="flex items-center gap-3 text-xs">
              <span className="flex items-center gap-1.5"><span className="w-2 h-2 rounded-full bg-rose-500" /> Occupied (3)</span>
              <span className="flex items-center gap-1.5"><span className="w-2 h-2 rounded-full bg-emerald-500" /> Available (2)</span>
              <span className="flex items-center gap-1.5"><span className="w-2 h-2 rounded-full bg-amber-500" /> Cleaning (1)</span>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {beds.map((bed) => (
              <Card key={bed.id} className="border-border shadow-sm bg-card hover:shadow-md transition-all">
                <CardHeader className="pb-2">
                  <div className="flex items-center justify-between">
                    <span className="text-xs font-mono font-bold text-foreground">{bed.bedNo}</span>
                    <Badge
                      className={`text-[9px] uppercase font-mono ${
                        bed.status === "occupied"
                          ? "bg-rose-500/10 text-rose-600 border-rose-500/30"
                          : bed.status === "available"
                          ? "bg-emerald-500/10 text-emerald-600 border-emerald-500/30"
                          : "bg-amber-500/10 text-amber-600 border-amber-500/30"
                      }`}
                    >
                      {bed.status}
                    </Badge>
                  </div>
                  <CardDescription className="text-xs">{bed.ward}</CardDescription>
                </CardHeader>
                <CardContent className="space-y-3 pt-0">
                  {bed.status === "occupied" ? (
                    <div className="p-3 rounded-lg bg-secondary/30 border border-border/50 space-y-1 text-xs">
                      <p className="font-semibold text-foreground">{bed.patientName}</p>
                      <p className="text-[11px] font-mono text-muted-foreground">{bed.patientMrn}</p>
                      <p className="text-[10px] text-muted-foreground flex items-center gap-1 pt-1">
                        <Clock className="w-3 h-3" /> Day {bed.admittedDays} of Admission
                      </p>
                    </div>
                  ) : (
                    <div className="py-4 text-center text-xs text-muted-foreground">
                      {bed.status === "available" ? "Ready for patient admission" : "Housekeeping in progress"}
                    </div>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </div>
    </CapabilityGate>
  );
}
