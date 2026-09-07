import React, { useState } from "react";
import {
  Dialog,
  DialogContent,
  Badge,
  Button,
} from "@curexal/ui";
import {
  ShieldAlert,
  Activity,
  Clock,
  X,
} from "lucide-react";
import { StatusBadge } from "@curexal/design-system";
import type { CanonicalPatient } from "@curexal/contracts";

export interface PatientDrawerProps {
  isOpen: boolean;
  onClose: () => void;
  patient: Partial<CanonicalPatient> | null;
}

export const PatientDrawer: React.FC<PatientDrawerProps> = ({
  isOpen,
  onClose,
  patient,
}) => {
  const [activeTab, setActiveTab] = useState<"summary" | "timeline" | "labs" | "meds" | "billing">("summary");

  if (!isOpen || !patient) return null;

  const fullName = `${patient.firstName || "Patient"} ${patient.lastName || ""}`;

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-2xl p-0 overflow-hidden border-border bg-card shadow-2xl rounded-2xl">
        {/* Header */}
        <div className="p-5 border-b border-border bg-secondary/30 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-primary/10 text-primary font-bold flex items-center justify-center border border-primary/20">
              {patient.firstName?.[0]}
              {patient.lastName?.[0]}
            </div>
            <div>
              <div className="flex items-center gap-2">
                <h3 className="font-bold text-base text-foreground">{fullName}</h3>
                <Badge variant="outline" className="font-mono text-[10px] text-primary border-primary/30">
                  {patient.mrn || "PAT-0000"}
                </Badge>
              </div>
              <p className="text-xs text-muted-foreground">
                {patient.gender || "Patient"} • DOB: {patient.dateOfBirth ? new Date(patient.dateOfBirth).toLocaleDateString() : "N/A"}
              </p>
            </div>
          </div>

          <button
            type="button"
            onClick={onClose}
            className="p-1.5 rounded-lg text-muted-foreground hover:text-foreground hover:bg-secondary transition-colors"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Tab Navigation */}
        <div className="flex items-center gap-2 px-5 pt-3 border-b border-border text-xs">
          {[
            { id: "summary", label: "Clinical Summary" },
            { id: "timeline", label: "Care Timeline" },
            { id: "labs", label: "Lab Orders & Results" },
            { id: "meds", label: "Active Rx" },
            { id: "billing", label: "Ledger & Claims" },
          ].map((tab) => (
            <button
              key={tab.id}
              type="button"
              onClick={() => setActiveTab(tab.id as any)}
              className={`pb-2.5 px-1 font-semibold border-b-2 transition-all cursor-pointer ${
                activeTab === tab.id
                  ? "border-primary text-primary"
                  : "border-transparent text-muted-foreground hover:text-foreground"
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {/* Tab Content Body */}
        <div className="p-5 max-h-[60vh] overflow-y-auto space-y-4 text-xs">
          {activeTab === "summary" && (
            <div className="space-y-4">
              {/* Vitals Snapshot */}
              <div className="p-3.5 rounded-xl border border-border bg-secondary/20 space-y-2">
                <p className="font-bold text-foreground flex items-center gap-1.5">
                  <Activity className="w-4 h-4 text-teal-600" />
                  Latest Clinical Vitals
                </p>
                <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 font-mono text-xs">
                  <div><span className="text-muted-foreground">BP:</span> <strong>120/80 mmHg</strong></div>
                  <div><span className="text-muted-foreground">Pulse:</span> <strong>72 bpm</strong></div>
                  <div><span className="text-muted-foreground">Temp:</span> <strong>36.8 °C</strong></div>
                  <div><span className="text-muted-foreground">SpO2:</span> <strong>98%</strong></div>
                </div>
              </div>

              {/* Biological Markers & Known Allergies */}
              <div className="grid grid-cols-2 gap-3">
                <div className="p-3 rounded-xl border border-border bg-card space-y-1">
                  <p className="text-muted-foreground font-medium text-[11px]">Blood Group & Genotype</p>
                  <p className="font-bold text-foreground text-sm">
                    {patient.bloodGroup || "O+"} • {patient.genotype || "AA"}
                  </p>
                </div>
                <div className="p-3 rounded-xl border border-rose-500/20 bg-rose-500/5 space-y-1">
                  <p className="text-rose-600 dark:text-rose-400 font-medium text-[11px] flex items-center gap-1">
                    <ShieldAlert className="w-3.5 h-3.5" />
                    Allergies & Contraindications
                  </p>
                  <p className="font-bold text-rose-700 dark:text-rose-300 text-xs">
                    Penicillin (Severe) • NSAIDs
                  </p>
                </div>
              </div>

              {/* Active Diagnoses */}
              <div className="p-3.5 rounded-xl border border-border bg-card space-y-2">
                <p className="font-bold text-foreground">Active Problem List / Diagnoses</p>
                <ul className="space-y-1 text-muted-foreground list-disc list-inside">
                  <li>Essential Hypertension (Stage 1) - Managed on Amlodipine</li>
                  <li>Acute Upper Respiratory Tract Infection (Resolving)</li>
                </ul>
              </div>
            </div>
          )}

          {activeTab === "timeline" && (
            <div className="space-y-3">
              {[
                { time: "Today, 10:15 AM", title: "Outpatient Consultation", by: "Dr. Amina Yusuf", badge: "In Progress" },
                { time: "Yesterday, 04:30 PM", title: "Complete Blood Count (CBC)", by: "Main Pathology Lab", badge: "Authorized" },
                { time: "12 Aug 2026", title: "Prescription Dispensed", by: "Pharmacy Dispensary", badge: "Dispensed" },
              ].map((t, idx) => (
                <div key={idx} className="flex items-start gap-3 p-3 rounded-xl border border-border bg-secondary/10">
                  <Clock className="w-4 h-4 text-primary shrink-0 mt-0.5" />
                  <div className="space-y-0.5 flex-1">
                    <div className="flex items-center justify-between">
                      <p className="font-bold text-foreground">{t.title}</p>
                      <StatusBadge status={t.badge} size="sm" />
                    </div>
                    <p className="text-muted-foreground text-[11px]">{t.by} • {t.time}</p>
                  </div>
                </div>
              ))}
            </div>
          )}

          {activeTab === "labs" && (
            <div className="space-y-2.5">
              {[
                { test: "Complete Blood Count (CBC)", status: "authorized", date: "Yesterday", ref: "LIS-8401" },
                { test: "Fasting Blood Sugar (FBS)", status: "in_analysis", date: "Today", ref: "LIS-8402" },
                { test: "Lipid Profile Panel", status: "collected", date: "Today", ref: "LIS-8403" },
              ].map((lab, i) => (
                <div key={i} className="p-3 rounded-xl border border-border flex items-center justify-between">
                  <div>
                    <p className="font-semibold text-foreground">{lab.test}</p>
                    <p className="text-[11px] font-mono text-muted-foreground">Accession: {lab.ref} • {lab.date}</p>
                  </div>
                  <StatusBadge status={lab.status} />
                </div>
              ))}
            </div>
          )}

          {activeTab === "meds" && (
            <div className="space-y-2.5">
              {[
                { drug: "Amlodipine 10mg", dose: "1 tablet daily PO", duration: "30 days", status: "active" },
                { drug: "Artemether/Lumefantrine 80/480mg", dose: "1 tab BD x 3 days", duration: "3 days", status: "completed" },
              ].map((med, i) => (
                <div key={i} className="p-3 rounded-xl border border-border flex items-center justify-between">
                  <div>
                    <p className="font-semibold text-foreground">{med.drug}</p>
                    <p className="text-[11px] text-muted-foreground">{med.dose} • {med.duration}</p>
                  </div>
                  <StatusBadge status={med.status} />
                </div>
              ))}
            </div>
          )}

          {activeTab === "billing" && (
            <div className="space-y-3">
              <div className="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-700 dark:text-emerald-400 flex items-center justify-between">
                <span>Current Account Balance:</span>
                <span className="font-mono font-bold text-sm">₦0.00 (Settled)</span>
              </div>
              <div className="p-3 rounded-xl border border-border space-y-1">
                <div className="flex justify-between font-semibold text-foreground">
                  <span>INV-2026-0841</span>
                  <span>₦17,500</span>
                </div>
                <p className="text-[11px] text-muted-foreground">Doctor Consultation + CBC • Paid via POS</p>
              </div>
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="p-4 bg-secondary/30 border-t border-border flex items-center justify-between">
          <span className="text-[11px] text-muted-foreground">Curexal Patient 360 Record</span>
          <Button size="sm" onClick={onClose} variant="outline" className="text-xs h-8">
            Close Inspector
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
};
