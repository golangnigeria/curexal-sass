import React, { useState } from "react";
import { CapabilityGate } from "@/components/design-system/capability-gate";
import { DataTable } from "@/components/design-system/data-table";
import { StatusPill } from "@/components/design-system/status-pill";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import {
  Pill,
  CheckCircle2,
  AlertCircle,
  Plus,
  Package,
} from "lucide-react";

interface PrescriptionOrder {
  id: string;
  rxNumber: string;
  patientName: string;
  prescribingDoctor: string;
  items: Array<{ name: string; dose: string; qty: number; batch: string }>;
  status: "pending" | "completed";
  orderedAt: string;
}

const mockPrescriptions: PrescriptionOrder[] = [
  { id: "rx-1", rxNumber: "RX-2026-091", patientName: "Amina Yusuf", prescribingDoctor: "Dr. Amina Yusuf", items: [{ name: "Artemether/Lumefantrine 80/480mg", dose: "1 tab BD x 3 days", qty: 6, batch: "BAT-8921" }, { name: "Paracetamol 500mg", dose: "2 tabs TDS x 3 days", qty: 18, batch: "BAT-7740" }], status: "pending", orderedAt: "15 mins ago" },
  { id: "rx-2", rxNumber: "RX-2026-092", patientName: "Chinedu Okafor", prescribingDoctor: "Dr. A. Adebayo", items: [{ name: "Amlodipine 10mg", dose: "1 tab daily", qty: 30, batch: "BAT-6651" }], status: "completed", orderedAt: "45 mins ago" },
];

export default function WorkspacePharmacyPage() {
  const [prescriptions, setPrescriptions] = useState<PrescriptionOrder[]>(mockPrescriptions);

  const handleDispense = (id: string) => {
    setPrescriptions((prev) =>
      prev.map((p) => (p.id === id ? { ...p, status: "completed" } : p))
    );
    toast.success("Prescription Dispensed Successfully!", {
      description: "Inventory stock depleted via FEFO batch allocation.",
    });
  };

  return (
    <CapabilityGate
      capability="pharmacy.basic"
      moduleCode="pharmacy"
      title="Dispensary & Pharmacy Suite"
      description="Unlock prescription validation, FEFO batch allocation, stock threshold alerts, and drug interaction checking."
      requiredPlan="Optimize or Pro"
    >
      <div className="space-y-8 animate-fade-in">
        {/* Header */}
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-6">
          <div>
            <div className="flex items-center gap-2.5 mb-1">
              <h1 className="text-2xl font-bold tracking-tight text-foreground flex items-center gap-2">
                <Pill className="w-6 h-6 text-emerald-600 dark:text-emerald-400" />
                Pharmacy & Prescription Dispensary
              </h1>
              <Badge variant="outline" className="border-emerald-500/40 text-emerald-600 dark:text-emerald-400 bg-emerald-500/10 text-[10px] font-mono">
                FEFO Batch Tracking
              </Badge>
            </div>
            <p className="text-xs text-muted-foreground">
              Electronic prescription queue, batch/lot expiration enforcement, and inventory dispensing.
            </p>
          </div>

          <div className="flex items-center gap-2.5">
            <Button size="sm" className="text-xs h-9 gap-1.5 bg-emerald-600 hover:bg-emerald-700 text-white shadow">
              <Package className="w-3.5 h-3.5" />
              Stock Inventory (POS)
            </Button>
          </div>
        </div>

        {/* Prescription Table */}
        <DataTable
          data={prescriptions}
          searchPlaceholder="Search Rx number, patient..."
          searchKey="rxNumber"
          columns={[
            {
              header: "Prescription #",
              cell: (p) => <span className="font-mono font-bold text-xs">{p.rxNumber}</span>,
            },
            {
              header: "Patient",
              cell: (p) => <span className="font-semibold text-foreground">{p.patientName}</span>,
            },
            {
              header: "Prescribed Items & Dosage",
              cell: (p) => (
                <div className="space-y-1">
                  {p.items.map((item, idx) => (
                    <div key={idx} className="text-xs flex items-center gap-1.5">
                      <span className="font-medium text-foreground">{item.name}</span>
                      <span className="text-[10px] font-mono text-muted-foreground">({item.dose})</span>
                      <Badge variant="secondary" className="text-[9px] font-mono">Lot: {item.batch}</Badge>
                    </div>
                  ))}
                </div>
              ),
            },
            {
              header: "Prescribing Clinician",
              cell: (p) => <span className="text-xs text-muted-foreground">{p.prescribingDoctor}</span>,
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
                  {p.status === "pending" && (
                    <Button
                      size="sm"
                      onClick={() => handleDispense(p.id)}
                      className="text-xs h-7 gap-1 bg-emerald-600 hover:bg-emerald-700 text-white"
                    >
                      <CheckCircle2 className="w-3 h-3" /> Dispense Rx
                    </Button>
                  )}
                  {p.status === "completed" && (
                    <span className="text-[11px] text-emerald-600 dark:text-emerald-400 font-semibold flex items-center gap-1">
                      Dispensed
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
