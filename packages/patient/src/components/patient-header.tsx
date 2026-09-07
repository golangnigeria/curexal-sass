import React from "react";
import {
  Calendar,
  Droplet,
  FileText,
} from "lucide-react";
import { Badge } from "@curexal/ui";
import type { CanonicalPatient } from "@curexal/contracts";

export interface PatientHeaderProps {
  patient: Partial<CanonicalPatient> & {
    ageGender?: string;
    vitalsSummary?: string;
    activeEncounterId?: string;
  };
  onOpenDrawer?: () => void;
  className?: string;
}

export const PatientHeader: React.FC<PatientHeaderProps> = ({
  patient,
  onOpenDrawer,
  className = "",
}) => {
  const initials = `${patient.firstName?.[0] || "P"}${patient.lastName?.[0] || ""}`;
  const fullName = `${patient.firstName || "Unknown"} ${patient.middleName ? `${patient.middleName} ` : ""}${patient.lastName || "Patient"}`;

  return (
    <div className={`p-4 rounded-xl border border-border bg-card shadow-sm flex flex-col sm:flex-row sm:items-center justify-between gap-4 ${className}`}>
      <div className="flex items-center gap-3.5">
        {/* Patient Avatar Circle */}
        <div className="w-11 h-11 rounded-xl bg-primary/10 text-primary font-bold text-sm flex items-center justify-center border border-primary/20 shrink-0">
          {initials}
        </div>

        {/* Patient Core Info */}
        <div className="space-y-1">
          <div className="flex flex-wrap items-center gap-2">
            <h3 className="font-bold text-sm text-foreground">{fullName}</h3>
            <Badge variant="outline" className="text-[10px] font-mono border-primary/30 text-primary bg-primary/5">
              MRN: {patient.mrn || "PAT-0000"}
            </Badge>
            {patient.gender && (
              <Badge variant="secondary" className="text-[10px] font-mono uppercase">
                {patient.gender}
              </Badge>
            )}
            {patient.activeEncounterId && (
              <Badge className="text-[9px] bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/30">
                In Active Consult
              </Badge>
            )}
          </div>

          <div className="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-muted-foreground">
            {patient.dateOfBirth && (
              <span className="flex items-center gap-1">
                <Calendar className="w-3.5 h-3.5" />
                DOB: {new Date(patient.dateOfBirth).toLocaleDateString()}
              </span>
            )}
            {patient.bloodGroup && (
              <span className="flex items-center gap-1 text-rose-600 dark:text-rose-400 font-medium">
                <Droplet className="w-3.5 h-3.5" />
                Blood: {patient.bloodGroup}
              </span>
            )}
            {patient.genotype && (
              <span className="text-amber-600 dark:text-amber-400 font-medium">
                Genotype: {patient.genotype}
              </span>
            )}
          </div>
        </div>
      </div>

      {/* Action / Inspection Trigger */}
      {onOpenDrawer && (
        <button
          type="button"
          onClick={onOpenDrawer}
          className="self-end sm:self-center text-xs font-semibold text-primary hover:underline flex items-center gap-1 cursor-pointer"
        >
          <FileText className="w-3.5 h-3.5" />
          <span>Patient 360 Record</span>
        </button>
      )}
    </div>
  );
};
