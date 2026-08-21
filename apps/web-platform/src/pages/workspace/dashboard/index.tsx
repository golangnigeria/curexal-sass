import React from "react";
import { Link } from "react-router-dom";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { useCapabilities } from "@/api/hooks/use-capabilities";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { StatCard } from "@/components/design-system/stat-card";
import {
  Activity,
  Users,
  CreditCard,
  Building2,
  Stethoscope,
  Microscope,
  Pill,
  ChevronRight,
  ArrowRight,
  Shield,
  Layers,
  Sparkles,
} from "lucide-react";

export default function WorkspaceDashboardPage() {
  const { data: bootstrap } = useBootstrap();
  const { hasCapability } = useCapabilities();

  const workspace = bootstrap?.workspace;
  const facilityName = workspace?.name || "Main Diagnostic Facility";
  const facilityType = workspace?.facilityType || "Diagnostic Center";
  const currency = workspace?.currency || "NGN";

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-6">
        <div>
          <div className="flex items-center gap-2.5 mb-1">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              {facilityName}
            </h1>
            <Badge variant="outline" className="border-primary/40 text-primary bg-primary/10 text-[10px] font-mono uppercase font-bold">
              {facilityType}
            </Badge>
          </div>
          <p className="text-xs text-muted-foreground">
            Operational workspace for daily patient reception, clinical triage, diagnostics, and cashier billing.
          </p>
        </div>

        <div className="flex items-center gap-2.5">
          <Button asChild size="sm" className="text-xs h-9 gap-1.5 bg-primary text-primary-foreground shadow">
            <Link to="/workspace/billing">
              <CreditCard className="w-3.5 h-3.5" />
              Cashier Register POS
            </Link>
          </Button>
        </div>
      </div>

      {/* KPI Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          title="Today's Patient Queue"
          value="48"
          icon={Users}
          iconColorClass="text-sky-500 bg-sky-500/10"
          trendPercentage={12}
          trendLabel="vs yesterday"
        />
        <StatCard
          title="Lab Samples in Worklist"
          value="112"
          icon={Microscope}
          iconColorClass="text-teal-500 bg-teal-500/10"
          trendPercentage={6}
          trendLabel="vs yesterday"
        />
        <StatCard
          title="Today's POS Collections"
          value={`₦680,000`}
          icon={CreditCard}
          iconColorClass="text-emerald-500 bg-emerald-500/10"
          trendPercentage={18}
          trendLabel="vs yesterday"
        />
        <StatCard
          title="Active Facility Modules"
          value={bootstrap?.modules?.filter((m) => m.enabled).length || 5}
          icon={Layers}
          iconColorClass="text-indigo-500 bg-indigo-500/10"
        />
      </div>

      {/* Operational Department Canvases */}
      <div className="space-y-4">
        <h3 className="text-base font-bold text-foreground">Facility Clinical & Diagnostic Workspaces</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {/* LIMS */}
          <Link
            to="/workspace/laboratory"
            className="p-5 rounded-2xl border border-border bg-card hover:border-teal-500/50 transition-all flex flex-col justify-between group shadow-sm"
          >
            <div>
              <div className="p-2.5 rounded-xl bg-teal-500/10 text-teal-600 dark:text-teal-400 w-fit mb-3 group-hover:scale-105 transition-transform">
                <Microscope className="w-5 h-5" />
              </div>
              <h4 className="text-sm font-bold text-foreground group-hover:text-teal-600 dark:group-hover:text-teal-400 transition-colors">
                Medical Laboratory (LIMS)
              </h4>
              <p className="text-xs text-muted-foreground mt-1">
                Specimen accessioning, automated analyzer worklists, and two-step verification.
              </p>
            </div>
            <div className="pt-4 border-t border-border mt-4 flex items-center justify-between text-xs text-teal-600 dark:text-teal-400 font-semibold">
              <span>Open Accessioning</span>
              <ArrowRight className="w-3.5 h-3.5" />
            </div>
          </Link>

          {/* EMR */}
          <Link
            to="/workspace/clinical"
            className="p-5 rounded-2xl border border-border bg-card hover:border-sky-500/50 transition-all flex flex-col justify-between group shadow-sm"
          >
            <div>
              <div className="p-2.5 rounded-xl bg-sky-500/10 text-sky-600 dark:text-sky-400 w-fit mb-3 group-hover:scale-105 transition-transform">
                <Stethoscope className="w-5 h-5" />
              </div>
              <h4 className="text-sm font-bold text-foreground group-hover:text-sky-600 dark:group-hover:text-sky-400 transition-colors">
                Outpatient Clinic (EMR)
              </h4>
              <p className="text-xs text-muted-foreground mt-1">
                Doctor consultation queue, SOAP notes, vitals recording, and electronic Rx pad.
              </p>
            </div>
            <div className="pt-4 border-t border-border mt-4 flex items-center justify-between text-xs text-sky-600 dark:text-sky-400 font-semibold">
              <span>Open Consultations</span>
              <ArrowRight className="w-3.5 h-3.5" />
            </div>
          </Link>

          {/* Radiology RIS */}
          <Link
            to="/workspace/radiology"
            className="p-5 rounded-2xl border border-border bg-card hover:border-amber-500/50 transition-all flex flex-col justify-between group shadow-sm"
          >
            <div>
              <div className="p-2.5 rounded-xl bg-amber-500/10 text-amber-600 dark:text-amber-400 w-fit mb-3 group-hover:scale-105 transition-transform">
                <Activity className="w-5 h-5" />
              </div>
              <h4 className="text-sm font-bold text-foreground group-hover:text-amber-600 dark:group-hover:text-amber-400 transition-colors">
                Radiology & Imaging (RIS)
              </h4>
              <p className="text-xs text-muted-foreground mt-1">
                DICOM modality queues (X-Ray, Ultrasound, CT), PACS viewer launcher, and scan reports.
              </p>
            </div>
            <div className="pt-4 border-t border-border mt-4 flex items-center justify-between text-xs text-amber-600 dark:text-amber-400 font-semibold">
              <span>Open Modality Queue</span>
              <ArrowRight className="w-3.5 h-3.5" />
            </div>
          </Link>

          {/* Pharmacy */}
          <Link
            to="/workspace/pharmacy"
            className="p-5 rounded-2xl border border-border bg-card hover:border-emerald-500/50 transition-all flex flex-col justify-between group shadow-sm"
          >
            <div>
              <div className="p-2.5 rounded-xl bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 w-fit mb-3 group-hover:scale-105 transition-transform">
                <Pill className="w-5 h-5" />
              </div>
              <h4 className="text-sm font-bold text-foreground group-hover:text-emerald-600 dark:group-hover:text-emerald-400 transition-colors">
                Pharmacy & Dispensary
              </h4>
              <p className="text-xs text-muted-foreground mt-1">
                Prescription fulfillment, FEFO batch/lot dispensing, and stock reorder warnings.
              </p>
            </div>
            <div className="pt-4 border-t border-border mt-4 flex items-center justify-between text-xs text-emerald-600 dark:text-emerald-400 font-semibold">
              <span>Open Dispensary</span>
              <ArrowRight className="w-3.5 h-3.5" />
            </div>
          </Link>

          {/* Inpatient Hospital HIS */}
          <Link
            to="/workspace/hospital"
            className="p-5 rounded-2xl border border-border bg-card hover:border-violet-500/50 transition-all flex flex-col justify-between group shadow-sm"
          >
            <div>
              <div className="p-2.5 rounded-xl bg-violet-500/10 text-violet-600 dark:text-violet-400 w-fit mb-3 group-hover:scale-105 transition-transform">
                <Building2 className="w-5 h-5" />
              </div>
              <h4 className="text-sm font-bold text-foreground group-hover:text-violet-600 dark:group-hover:text-violet-400 transition-colors">
                Inpatient Wards (HIS)
              </h4>
              <p className="text-xs text-muted-foreground mt-1">
                Visual Ward Bed grid, admission/discharge management, and nurse shift handovers.
              </p>
            </div>
            <div className="pt-4 border-t border-border mt-4 flex items-center justify-between text-xs text-violet-600 dark:text-violet-400 font-semibold">
              <span>Open Ward Bed Board</span>
              <ArrowRight className="w-3.5 h-3.5" />
            </div>
          </Link>

          {/* Billing POS */}
          <Link
            to="/workspace/billing"
            className="p-5 rounded-2xl border border-border bg-card hover:border-primary/50 transition-all flex flex-col justify-between group shadow-sm"
          >
            <div>
              <div className="p-2.5 rounded-xl bg-primary/10 text-primary w-fit mb-3 group-hover:scale-105 transition-transform">
                <CreditCard className="w-5 h-5" />
              </div>
              <h4 className="text-sm font-bold text-foreground group-hover:text-primary transition-colors">
                Billing & Cashier Register
              </h4>
              <p className="text-xs text-muted-foreground mt-1">
                Service billing POS, split cashier receipts, HMO insurance claims, and payment records.
              </p>
            </div>
            <div className="pt-4 border-t border-border mt-4 flex items-center justify-between text-xs text-primary font-semibold">
              <span>Open Cashier POS</span>
              <ArrowRight className="w-3.5 h-3.5" />
            </div>
          </Link>
        </div>
      </div>
    </div>
  );
}
