import React from "react";
import { Link, useParams } from "react-router-dom";
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
  const { branchSlug } = useParams<{ branchSlug?: string }>();

  const activeBranchSlug = branchSlug || bootstrap?.branch?.slug || bootstrap?.workspace?.slug || "main";
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

        <div className="flex items-center gap-3">
          <Badge variant="secondary" className="px-3 py-1.5 text-xs font-mono">
            {currency} Currency Engine
          </Badge>
          <Button asChild size="sm" className="gap-1.5 text-xs">
            <Link to={`/${activeBranchSlug}/billing`}>
              <Sparkles className="w-3.5 h-3.5" />
              New Patient Transaction
            </Link>
          </Button>
        </div>
      </div>

      {/* KPI Highlights */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          title="Patients Registered Today"
          value="48"
          icon={Users}
          iconColorClass="text-sky-500 bg-sky-500/10"
          trendPercentage={12}
          trendLabel="vs last week"
        />
        <StatCard
          title="Clinical & Lab Orders"
          value="112"
          icon={Activity}
          iconColorClass="text-teal-500 bg-teal-500/10"
          trendPercentage={-3}
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
        <h3 className="text-base font-bold text-foreground">Facility Clinical & Operational Workspaces</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {/* Patient Reception */}
          <Link
            to={`/${activeBranchSlug}/reception`}
            className="p-5 rounded-2xl border border-border bg-card hover:border-sky-500/50 transition-all flex flex-col justify-between group shadow-sm"
          >
            <div>
              <div className="p-2.5 rounded-xl bg-sky-500/10 text-sky-600 dark:text-sky-400 w-fit mb-3 group-hover:scale-105 transition-transform">
                <Users className="w-5 h-5" />
              </div>
              <h4 className="text-sm font-bold text-foreground group-hover:text-sky-600 dark:group-hover:text-sky-400 transition-colors">
                Patient Reception & MPI
              </h4>
              <p className="text-xs text-muted-foreground mt-1">
                Walk-in intake, Master Patient Index identity verification, and appointment booking.
              </p>
            </div>
            <div className="pt-4 border-t border-border mt-4 flex items-center justify-between text-xs text-sky-600 dark:text-sky-400 font-semibold">
              <span>Open Reception</span>
              <ArrowRight className="w-3.5 h-3.5" />
            </div>
          </Link>

          {/* Triage & Care Desk */}
          <Link
            to={`/${activeBranchSlug}/care-desk`}
            className="p-5 rounded-2xl border border-border bg-card hover:border-teal-500/50 transition-all flex flex-col justify-between group shadow-sm"
          >
            <div>
              <div className="p-2.5 rounded-xl bg-teal-500/10 text-teal-600 dark:text-teal-400 w-fit mb-3 group-hover:scale-105 transition-transform">
                <Activity className="w-5 h-5" />
              </div>
              <h4 className="text-sm font-bold text-foreground group-hover:text-teal-600 dark:group-hover:text-teal-400 transition-colors">
                Nursing Triage Desk
              </h4>
              <p className="text-xs text-muted-foreground mt-1">
                Vital signs recording, automated acuity scoring, and nurse-to-doctor handoff queue.
              </p>
            </div>
            <div className="pt-4 border-t border-border mt-4 flex items-center justify-between text-xs text-teal-600 dark:text-teal-400 font-semibold">
              <span>Open Triage Desk</span>
              <ArrowRight className="w-3.5 h-3.5" />
            </div>
          </Link>

          {/* Doctor EMR Room */}
          <Link
            to={`/${activeBranchSlug}/clinical`}
            className="p-5 rounded-2xl border border-border bg-card hover:border-indigo-500/50 transition-all flex flex-col justify-between group shadow-sm"
          >
            <div>
              <div className="p-2.5 rounded-xl bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 w-fit mb-3 group-hover:scale-105 transition-transform">
                <Stethoscope className="w-5 h-5" />
              </div>
              <h4 className="text-sm font-bold text-foreground group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors">
                Outpatient Clinic & EMR
              </h4>
              <p className="text-xs text-muted-foreground mt-1">
                Doctor consultation queue, electronic SOAP notes, ICD-10 coding, and digital prescriptions.
              </p>
            </div>
            <div className="pt-4 border-t border-border mt-4 flex items-center justify-between text-xs text-indigo-600 dark:text-indigo-400 font-semibold">
              <span>Open Consultations</span>
              <ArrowRight className="w-3.5 h-3.5" />
            </div>
          </Link>

          {/* Billing POS */}
          <Link
            to={`/${activeBranchSlug}/billing`}
            className="p-5 rounded-2xl border border-border bg-card hover:border-emerald-500/50 transition-all flex flex-col justify-between group shadow-sm"
          >
            <div>
              <div className="p-2.5 rounded-xl bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 w-fit mb-3 group-hover:scale-105 transition-transform">
                <CreditCard className="w-5 h-5" />
              </div>
              <h4 className="text-sm font-bold text-foreground group-hover:text-emerald-600 dark:group-hover:text-emerald-400 transition-colors">
                Cashier Billing & POS
              </h4>
              <p className="text-xs text-muted-foreground mt-1">
                Clinical fee settlement, multi-tender POS receipts (Cash, Card, Transfer), and payment history.
              </p>
            </div>
            <div className="pt-4 border-t border-border mt-4 flex items-center justify-between text-xs text-emerald-600 dark:text-emerald-400 font-semibold">
              <span>Open Cashier POS</span>
              <ArrowRight className="w-3.5 h-3.5" />
            </div>
          </Link>
        </div>
      </div>
    </div>
  );
}
