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

      {/* Interactive Patient Flow & Throughput Funnel */}
      <Card className="border-border shadow-sm bg-card">
        <CardHeader className="pb-3 border-b border-border flex flex-row items-center justify-between">
          <div>
            <CardTitle className="text-sm font-bold text-foreground flex items-center gap-2">
              <Activity className="w-4 h-4 text-teal-600" />
              Live Patient Care Pipeline (Today's Census)
            </CardTitle>
            <CardDescription className="text-xs">
              Real-time patient distribution across clinic service stations.
            </CardDescription>
          </div>
          <Badge variant="outline" className="text-[10px] font-mono text-teal-600 border-teal-500/30">
            48 Total Registered
          </Badge>
        </CardHeader>
        <CardContent className="p-4">
          <div className="grid grid-cols-2 sm:grid-cols-5 gap-3 text-center">
            <Link
              to={`/${activeBranchSlug}/reception`}
              className="p-3 rounded-xl bg-sky-500/10 border border-sky-500/20 hover:border-sky-500/40 transition-colors group cursor-pointer block"
            >
              <p className="text-[11px] font-medium text-sky-700 dark:text-sky-300">1. Reception Intake</p>
              <p className="text-xl font-bold text-sky-900 dark:text-sky-100 mt-0.5">8</p>
              <p className="text-[10px] text-muted-foreground mt-0.5 group-hover:text-foreground transition-colors">Awaiting MPI →</p>
            </Link>
            <Link
              to={`/${activeBranchSlug}/care-desk`}
              className="p-3 rounded-xl bg-teal-500/10 border border-teal-500/20 hover:border-teal-500/40 transition-colors group cursor-pointer block"
            >
              <p className="text-[11px] font-medium text-teal-700 dark:text-teal-300">2. Nursing Triage</p>
              <p className="text-xl font-bold text-teal-900 dark:text-teal-100 mt-0.5">5</p>
              <p className="text-[10px] text-muted-foreground mt-0.5 group-hover:text-foreground transition-colors">Taking Vitals →</p>
            </Link>
            <Link
              to={`/${activeBranchSlug}/clinical`}
              className="p-3 rounded-xl bg-indigo-500/10 border border-indigo-500/20 hover:border-indigo-500/40 transition-colors group cursor-pointer block"
            >
              <p className="text-[11px] font-medium text-indigo-700 dark:text-indigo-300">3. In Consultation</p>
              <p className="text-xl font-bold text-indigo-900 dark:text-indigo-100 mt-0.5">4</p>
              <p className="text-[10px] text-muted-foreground mt-0.5 group-hover:text-foreground transition-colors">With Physicians →</p>
            </Link>
            <Link
              to={`/${activeBranchSlug}/billing`}
              className="p-3 rounded-xl bg-amber-500/10 border border-amber-500/20 hover:border-amber-500/40 transition-colors group cursor-pointer block"
            >
              <p className="text-[11px] font-medium text-amber-700 dark:text-amber-300">4. Cashier Billing</p>
              <p className="text-xl font-bold text-amber-900 dark:text-amber-100 mt-0.5">3</p>
              <p className="text-[10px] text-muted-foreground mt-0.5 group-hover:text-foreground transition-colors">Awaiting POS →</p>
            </Link>
            <div className="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 col-span-2 sm:col-span-1">
              <p className="text-[11px] font-medium text-emerald-700 dark:text-emerald-300">5. Completed</p>
              <p className="text-xl font-bold text-emerald-900 dark:text-emerald-100 mt-0.5">28</p>
              <p className="text-[10px] text-muted-foreground mt-0.5">Discharged</p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Consultation Room Status & Critical Triage Alerts */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Active Consultation Rooms Board */}
        <Card className="lg:col-span-2 border-border shadow-sm bg-card">
          <CardHeader className="pb-3 border-b border-border flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-sm font-bold text-foreground flex items-center gap-2">
                <Stethoscope className="w-4 h-4 text-indigo-600" />
                Physician Consultation Rooms
              </CardTitle>
              <CardDescription className="text-xs">
                Attending doctors on duty and active room assignments.
              </CardDescription>
            </div>
            <Button asChild size="sm" variant="outline" className="text-xs h-7">
              <Link to={`/${activeBranchSlug}/clinical`}>Open EMR Canvas</Link>
            </Button>
          </CardHeader>
          <CardContent className="p-4 space-y-3">
            <div className="p-3 rounded-xl border border-border bg-secondary/15 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-lg bg-indigo-500/10 text-indigo-600 font-bold flex items-center justify-center text-xs">
                  R1
                </div>
                <div>
                  <p className="text-xs font-bold text-foreground">Dr. Emeka Nwosu • Consultation Room 1</p>
                  <p className="text-[11px] text-muted-foreground">In Consult with Amina Yusuf (14 mins elapsed)</p>
                </div>
              </div>
              <Badge className="bg-emerald-500/10 text-emerald-600 border-emerald-500/30 text-[10px]">
                In Progress
              </Badge>
            </div>

            <div className="p-3 rounded-xl border border-border bg-secondary/15 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-lg bg-indigo-500/10 text-indigo-600 font-bold flex items-center justify-center text-xs">
                  R2
                </div>
                <div>
                  <p className="text-xs font-bold text-foreground">Dr. Sarah Adebayo • Consultation Room 2</p>
                  <p className="text-[11px] text-muted-foreground">In Consult with Babatunde Lawal (6 mins elapsed)</p>
                </div>
              </div>
              <Badge className="bg-emerald-500/10 text-emerald-600 border-emerald-500/30 text-[10px]">
                In Progress
              </Badge>
            </div>

            <div className="p-3 rounded-xl border border-emerald-500/30 bg-emerald-500/5 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-lg bg-emerald-500/10 text-emerald-600 font-bold flex items-center justify-center text-xs">
                  R3
                </div>
                <div>
                  <p className="text-xs font-bold text-foreground">Dr. Chinedu Eke • Consultation Room 3</p>
                  <p className="text-[11px] text-emerald-600 font-medium">Ready for next patient (Next in Queue: Chinedu Okafor)</p>
                </div>
              </div>
              <Button asChild size="sm" className="text-xs h-7 gap-1 bg-emerald-600 hover:bg-emerald-700 text-white shadow-sm">
                <Link to={`/${activeBranchSlug}/clinical`}>Call Patient</Link>
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Critical Triage Alerts Widget */}
        <Card className="border-rose-500/30 bg-rose-500/5 shadow-sm">
          <CardHeader className="pb-3 border-b border-rose-500/20">
            <CardTitle className="text-sm font-bold text-rose-700 dark:text-rose-400 flex items-center gap-1.5">
              <Activity className="w-4 h-4 text-rose-600" />
              High-Acuity Triage Alerts
            </CardTitle>
            <CardDescription className="text-xs text-rose-600/80">
              Vitals exceeding clinical thresholds requiring immediate attention.
            </CardDescription>
          </CardHeader>
          <CardContent className="p-4 space-y-3">
            <div className="p-3 rounded-xl border border-rose-500/20 bg-background/80 space-y-1.5 shadow-sm">
              <div className="flex items-center justify-between">
                <span className="text-xs font-bold text-foreground">Babatunde Lawal (PAT-0078)</span>
                <Badge variant="outline" className="text-[10px] uppercase font-mono border-rose-500 text-rose-600 bg-rose-500/10 font-bold">
                  Emergency
                </Badge>
              </div>
              <p className="text-[11px] text-rose-600 font-medium">
                BP 175/110 mmHg • SpO2 94% • HR 104 bpm
              </p>
              <p className="text-[11px] text-muted-foreground">
                Chest tightness & exertional dyspnea (Waiting: 5 mins)
              </p>
              <Button asChild size="sm" variant="destructive" className="w-full text-xs h-7 mt-2 shadow-sm font-semibold">
                <Link to={`/${activeBranchSlug}/clinical`}>Expedite to Doctor</Link>
              </Button>
            </div>
          </CardContent>
        </Card>
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
