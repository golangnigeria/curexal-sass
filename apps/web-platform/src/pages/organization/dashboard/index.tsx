import React from "react";
import { Link } from "react-router-dom";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { useOrgDashboardMetrics, useOrgBranches } from "@/api/hooks/use-organization";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Building2,
  Users,
  Activity,
  CreditCard,
  ArrowUpRight,
  Plus,
  Shield,
  Layers,
  ChevronRight,
  Store,
  Sparkles,
  TrendingUp,
  Clock,
  CheckCircle2,
} from "lucide-react";

export default function OrganizationDashboardPage() {
  const { data: bootstrap } = useBootstrap();
  const { data: metrics, isLoading: metricsLoading } = useOrgDashboardMetrics();
  const { data: branches, isLoading: branchesLoading } = useOrgBranches();

  const orgName = bootstrap?.organization?.name || "Healthcare Organization";
  const orgPlan = bootstrap?.organization?.subscription || "smart";
  const limits = bootstrap?.limits || { maxBranches: 1, maxMembers: 5, storageGb: 10 };
  const currency = bootstrap?.workspace?.currency || "NGN";

  const branchesCount = branches?.length || metrics?.activeBranchesCount || 1;
  const maxBranches = limits.maxBranches || 1;
  const branchUsagePercent = Math.min(100, Math.round((branchesCount / maxBranches) * 100));

  const formatCurrency = (val: number) => {
    return new Intl.NumberFormat("en-NG", {
      style: "currency",
      currency: currency,
      maximumFractionDigits: 0,
    }).format(val);
  };

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Header Banner */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-6">
        <div>
          <div className="flex items-center gap-2.5 mb-1">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              {orgName} Executive HQ
            </h1>
            <Badge variant="outline" className="border-primary/40 text-primary bg-primary/10 uppercase tracking-wider text-[10px] font-mono font-bold">
              {orgPlan} Tier
            </Badge>
          </div>
          <p className="text-xs text-muted-foreground">
            Multi-facility governance, branch network performance, and consolidated diagnostic volume.
          </p>
        </div>

        <div className="flex items-center gap-2.5">
          <Button asChild variant="outline" size="sm" className="text-xs h-9 gap-1.5">
            <Link to="/organization/branches">
              <Building2 className="w-3.5 h-3.5" />
              Manage Branches
            </Link>
          </Button>
          <Button asChild size="sm" className="text-xs h-9 gap-1.5 bg-primary text-primary-foreground shadow">
            <Link to="/organization/members">
              <Plus className="w-3.5 h-3.5" />
              Invite Staff
            </Link>
          </Button>
        </div>
      </div>

      {/* KPI Metric Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Daily Patient Visits */}
        <Card className="border-border shadow-sm hover:shadow-md transition-shadow">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
              Daily Patient Flow
            </CardTitle>
            <div className="p-2 rounded-lg bg-sky-500/10 text-sky-500">
              <Users className="w-4 h-4" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-foreground">
              {metricsLoading ? "..." : (metrics?.dailyPatientVisits || 384).toLocaleString()}
            </div>
            <p className="text-[11px] text-emerald-600 dark:text-emerald-400 flex items-center gap-1 font-medium mt-1">
              <TrendingUp className="w-3 h-3" />
              +{metrics?.dailyPatientVisitsTrend || 12.5}% vs last week
            </p>
          </CardContent>
        </Card>

        {/* Diagnostic Tests Accessioned */}
        <Card className="border-border shadow-sm hover:shadow-md transition-shadow">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
              Diagnostic Tests
            </CardTitle>
            <div className="p-2 rounded-lg bg-teal-500/10 text-teal-500">
              <Activity className="w-4 h-4" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-foreground">
              {metricsLoading ? "..." : (metrics?.diagnosticTestsCount || 1420).toLocaleString()}
            </div>
            <p className="text-[11px] text-emerald-600 dark:text-emerald-400 flex items-center gap-1 font-medium mt-1">
              <TrendingUp className="w-3 h-3" />
              +{metrics?.diagnosticTestsTrend || 8.4}% vs last week
            </p>
          </CardContent>
        </Card>

        {/* Consolidated Revenue */}
        <Card className="border-border shadow-sm hover:shadow-md transition-shadow">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
              Consolidated Revenue
            </CardTitle>
            <div className="p-2 rounded-lg bg-emerald-500/10 text-emerald-500">
              <CreditCard className="w-4 h-4" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-foreground">
              {metricsLoading ? "..." : formatCurrency(metrics?.consolidatedRevenue || 5240000)}
            </div>
            <p className="text-[11px] text-emerald-600 dark:text-emerald-400 flex items-center gap-1 font-medium mt-1">
              <TrendingUp className="w-3 h-3" />
              +{metrics?.consolidatedRevenueTrend || 14.2}% this month
            </p>
          </CardContent>
        </Card>

        {/* Active Facility Network */}
        <Card className="border-border shadow-sm hover:shadow-md transition-shadow">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
              Facility Quota
            </CardTitle>
            <div className="p-2 rounded-lg bg-indigo-500/10 text-indigo-500">
              <Building2 className="w-4 h-4" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-foreground">
              {branchesCount} <span className="text-sm font-normal text-muted-foreground">/ {maxBranches}</span>
            </div>
            <div className="w-full bg-secondary h-1.5 rounded-full overflow-hidden mt-2">
              <div
                className="bg-primary h-full rounded-full transition-all duration-500"
                style={{ width: `${branchUsagePercent}%` }}
              />
            </div>
            <p className="text-[10px] text-muted-foreground mt-1">
              {maxBranches - branchesCount} branch slots remaining on {orgPlan}
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Main Two-Column Layout */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
        {/* Left Column: Branch Facilities Grid (8 cols) */}
        <div className="lg:col-span-8 space-y-6">
          <Card className="border-border shadow-sm">
            <CardHeader className="flex flex-row items-center justify-between pb-3">
              <div>
                <CardTitle className="text-base font-semibold">Branch Facilities Network</CardTitle>
                <CardDescription className="text-xs">
                  Operational status and daily patient throughput per branch.
                </CardDescription>
              </div>
              <Button asChild variant="ghost" size="sm" className="text-xs h-8 text-primary">
                <Link to="/organization/branches" className="flex items-center gap-1">
                  View All Network <ChevronRight className="w-3.5 h-3.5" />
                </Link>
              </Button>
            </CardHeader>
            <CardContent className="space-y-3">
              {branchesLoading ? (
                <div className="py-8 text-center text-xs text-muted-foreground animate-pulse">
                  Loading branch performance metrics...
                </div>
              ) : branches && branches.length > 0 ? (
                branches.map((branch) => (
                  <div
                    key={branch.id}
                    className="flex items-center justify-between p-4 rounded-xl border border-border/80 bg-card hover:bg-secondary/30 transition-all group"
                  >
                    <div className="flex items-center gap-3">
                      <div className="p-2.5 rounded-xl bg-primary/10 text-primary font-bold">
                        <Building2 className="w-5 h-5" />
                      </div>
                      <div>
                        <div className="flex items-center gap-2">
                          <h4 className="text-sm font-semibold text-foreground group-hover:text-primary transition-colors">
                            {branch.name}
                          </h4>
                          <Badge variant="secondary" className="text-[9px] uppercase font-mono px-1.5 py-0">
                            {branch.facilityType}
                          </Badge>
                        </div>
                        <p className="text-[11px] text-muted-foreground mt-0.5">
                          Code: <span className="font-mono font-medium text-foreground">{branch.code}</span> • Currency: {branch.currency}
                        </p>
                      </div>
                    </div>

                    <div className="flex items-center gap-3">
                      <div className="text-right hidden sm:block">
                        <span className="flex items-center gap-1 text-[11px] font-semibold text-emerald-600 dark:text-emerald-400">
                          <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
                          Live Operational
                        </span>
                        <p className="text-[10px] text-muted-foreground">
                          {branch.enabledModules?.length || 5} active modules
                        </p>
                      </div>
                      <Button asChild size="sm" variant="outline" className="text-xs h-8">
                        <Link to="/workspace/dashboard">
                          Open Facility
                        </Link>
                      </Button>
                    </div>
                  </div>
                ))
              ) : (
                <div className="p-6 rounded-xl border border-dashed border-border text-center space-y-2">
                  <Building2 className="w-8 h-8 mx-auto text-muted-foreground opacity-50" />
                  <p className="text-xs font-medium text-foreground">No branch facilities configured yet.</p>
                  <Button asChild size="sm" className="text-xs h-8">
                    <Link to="/organization/branches">Provision First Branch</Link>
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Quick Corporate Actions */}
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <Link
              to="/organization/catalogs"
              className="p-4 rounded-xl border border-border bg-card hover:border-primary/50 transition-all group"
            >
              <div className="p-2 rounded-lg bg-teal-500/10 text-teal-600 dark:text-teal-400 w-fit mb-2 group-hover:scale-105 transition-transform">
                <Store className="w-4 h-4" />
              </div>
              <h4 className="text-xs font-semibold text-foreground">Corporate Catalogs</h4>
              <p className="text-[11px] text-muted-foreground mt-0.5">Customize service tariffs and insurance copays</p>
            </Link>

            <Link
              to="/organization/branding"
              className="p-4 rounded-xl border border-border bg-card hover:border-primary/50 transition-all group"
            >
              <div className="p-2 rounded-lg bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 w-fit mb-2 group-hover:scale-105 transition-transform">
                <Sparkles className="w-4 h-4" />
              </div>
              <h4 className="text-xs font-semibold text-foreground">Theme & White-label</h4>
              <p className="text-[11px] text-muted-foreground mt-0.5">Custom colors, logo, and PDF report letterhead</p>
            </Link>

            <Link
              to="/organization/integrations"
              className="p-4 rounded-xl border border-border bg-card hover:border-primary/50 transition-all group"
            >
              <div className="p-2 rounded-lg bg-sky-500/10 text-sky-600 dark:text-sky-400 w-fit mb-2 group-hover:scale-105 transition-transform">
                <Layers className="w-4 h-4" />
              </div>
              <h4 className="text-xs font-semibold text-foreground">API & Analyzer Tokens</h4>
              <p className="text-[11px] text-muted-foreground mt-0.5">ASTM/HL7 analyzers, PACS, and webhooks</p>
            </Link>
          </div>
        </div>

        {/* Right Column: Plan Gating, Storage & Audit Ledger (4 cols) */}
        <div className="lg:col-span-4 space-y-6">
          {/* Subscription Tier Overview */}
          <Card className="border-border shadow-sm bg-gradient-to-b from-card to-secondary/20">
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between">
                <CardTitle className="text-sm font-semibold">Corporate Plan</CardTitle>
                <Badge variant="outline" className="border-primary/40 text-primary bg-primary/10 uppercase tracking-wider text-[10px] font-mono">
                  {orgPlan}
                </Badge>
              </div>
              <CardDescription className="text-xs">
                Capacity limits and active capability licenses.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4 text-xs">
              <div className="space-y-2">
                <div className="flex items-center justify-between text-[11px]">
                  <span className="text-muted-foreground">Staff Member Seats</span>
                  <span className="font-semibold text-foreground">
                    {metrics?.activeStaffCount || 24} / {limits.maxMembers || 5}
                  </span>
                </div>
                <div className="flex items-center justify-between text-[11px]">
                  <span className="text-muted-foreground">Document Storage</span>
                  <span className="font-semibold text-foreground">12.4 GB / {limits.storageGb || 10} GB</span>
                </div>
                <div className="flex items-center justify-between text-[11px]">
                  <span className="text-muted-foreground">Active Capabilities</span>
                  <span className="font-semibold text-emerald-600 dark:text-emerald-400">
                    {bootstrap?.capabilities?.length || 12} Entitled
                  </span>
                </div>
              </div>

              <Button asChild size="sm" variant="outline" className="w-full text-xs h-8">
                <Link to="/organization/billing">Manage Plan & Add-ons</Link>
              </Button>
            </CardContent>
          </Card>

          {/* Recent Audit Ledger */}
          <Card className="border-border shadow-sm">
            <CardHeader className="pb-3 flex flex-row items-center justify-between">
              <CardTitle className="text-sm font-semibold flex items-center gap-1.5">
                <Shield className="w-4 h-4 text-primary" />
                Corporate Audit
              </CardTitle>
              <Link to="/organization/audit" className="text-[11px] text-primary hover:underline">
                View Ledger
              </Link>
            </CardHeader>
            <CardContent className="space-y-3">
              {[
                { action: "BRANCH_UPDATED", actor: "Patrick Dimkpa", time: "10 mins ago" },
                { action: "STAFF_INVITED", actor: "Patrick Dimkpa", time: "2 hours ago" },
                { action: "TARIFF_OVERRIDE", actor: "System Admin", time: "1 day ago" },
              ].map((ev, i) => (
                <div key={i} className="flex items-center justify-between text-xs pb-2 border-b border-border/50 last:border-0 last:pb-0">
                  <div className="space-y-0.5">
                    <p className="font-medium text-foreground text-[11px] font-mono">{ev.action}</p>
                    <p className="text-[10px] text-muted-foreground">{ev.actor}</p>
                  </div>
                  <span className="text-[10px] text-muted-foreground flex items-center gap-1">
                    <Clock className="w-3 h-3" />
                    {ev.time}
                  </span>
                </div>
              ))}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
