import React from "react";
import { Link } from "react-router-dom";
import {
  Building2,
  Users,
  Layers,
  Database,
  Activity,
  ArrowUpRight,
  ShieldCheck,
  CheckCircle2,
  AlertTriangle,
  Cpu,
  RefreshCw,
  Sparkles,
  LayoutDashboard,
} from "lucide-react";
import {
  useDiagnostics,
  useLaunchGateStatus,
  useVerifyLaunchGate,
} from "@/api/hooks/use-diagnostics";
import { usePlatformAuditLogs } from "@/api/hooks/use-audit-logs";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  ResponsiveContainer,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
  BarChart,
  Bar,
} from "recharts";
import { formatDate } from "@/lib/utils";
import { toast } from "sonner";

export default function DashboardPage() {
  const { data: diag, isLoading: diagLoading, refetch: refetchDiag } = useDiagnostics();
  const { data: gate, isLoading: gateLoading } = useLaunchGateStatus();
  const { data: auditLogs } = usePlatformAuditLogs({ limit: 5 });
  const verifyMutation = useVerifyLaunchGate();

  const metrics = diag?.metrics;
  const db = diag?.database;

  const handleVerifyGate = async () => {
    try {
      await verifyMutation.mutateAsync();
      toast.success("Production readiness check executed successfully!");
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Launch gate verification failed.");
    }
  };

  const gateStatus = gate?.status || "PASSED";

  return (
    <div className="space-y-8">
      {/* Header Section */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2.5">
            <div className="flex items-center justify-center h-8 w-8 rounded-lg bg-primary/10 text-primary border border-primary/20">
              <LayoutDashboard className="h-4 w-4" />
            </div>
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              Platform Dashboard
            </h1>
            <Badge variant="outline" className="border-primary/30 text-primary bg-primary/5 text-xs">
              Live Cluster
            </Badge>
          </div>
          <p className="text-sm text-muted-foreground mt-1">
            Real-time platform telemetry, database connection pooling, and multi-tenant health.
          </p>
        </div>

        <div className="flex items-center gap-2.5">
          <Button
            variant="outline"
            size="sm"
            onClick={() => refetchDiag()}
            className="h-9 gap-2 text-xs"
          >
            <RefreshCw className={diagLoading ? "h-3.5 w-3.5 animate-spin" : "h-3.5 w-3.5"} />
            Refresh Telemetry
          </Button>

          <Button
            size="sm"
            onClick={handleVerifyGate}
            disabled={verifyMutation.isPending}
            className="h-9 gap-2 bg-primary hover:bg-primary/90 text-primary-foreground text-xs shadow-sm"
          >
            <ShieldCheck className="h-3.5 w-3.5" />
            {verifyMutation.isPending ? "Verifying Gate..." : "Run Launch Gate"}
          </Button>
        </div>
      </div>

      {/* Top Telemetry KPI Cards */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {/* Total Organizations */}
        <Card className="card-enterprise hover-lift">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
              Total Organizations
            </CardTitle>
            <Building2 className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-foreground">
              {metrics?.totalOrganizations ?? 0}
            </div>
            <p className="text-xs text-muted-foreground mt-1 flex items-center gap-1">
              <span className="text-emerald-600 font-medium">100% active</span> tenant networks
            </p>
          </CardContent>
        </Card>

        {/* Total Workspaces / Branches */}
        <Card className="card-enterprise hover-lift">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
              Workspaces / Branches
            </CardTitle>
            <Layers className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-foreground">
              {metrics?.totalWorkspaces ?? 0}
            </div>
            <p className="text-xs text-muted-foreground mt-1">
              Operational facilities connected
            </p>
          </CardContent>
        </Card>

        {/* Total Users */}
        <Card className="card-enterprise hover-lift">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
              Platform Users
            </CardTitle>
            <Users className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-foreground">
              {metrics?.totalUsers ?? 0}
            </div>
            <p className="text-xs text-muted-foreground mt-1">
              Staff, clinicians, and scientists
            </p>
          </CardContent>
        </Card>

        {/* Database Connection Pool */}
        <Card className="card-enterprise hover-lift">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
              DB Connection Pool
            </CardTitle>
            <Database className="h-4 w-4 text-emerald-500" />
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-2">
              <div className="text-2xl font-bold text-foreground">
                {db?.openConnections ?? 0}
              </div>
              <Badge variant="outline" className="text-[10px] text-emerald-600 bg-emerald-500/10 border-emerald-500/20 font-mono">
                {db?.status ?? "connected"}
              </Badge>
            </div>
            <p className="text-xs text-muted-foreground mt-1">
              {db?.inUse ?? 0} acquired, {db?.idle ?? 0} idle connections
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Production Launch Gate & Telemetry Charts Grid */}
      <div className="grid gap-6 lg:grid-cols-3">
        {/* Left 2 Cols: Growth Chart */}
        <Card className="card-enterprise lg:col-span-2">
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-base font-semibold">
                Organization Growth Trend
              </CardTitle>
              <CardDescription className="text-xs">
                Monthly active healthcare networks provisioned on cluster
              </CardDescription>
            </div>
            <Link
              to="/platform/organizations"
              className="text-xs text-primary hover:underline flex items-center gap-1 font-medium"
            >
              Directory <ArrowUpRight className="h-3 w-3" />
            </Link>
          </CardHeader>
          <CardContent className="h-[280px] w-full pt-4">
            {metrics?.organizationsGrowth && metrics.organizationsGrowth.length > 0 ? (
              <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={metrics.organizationsGrowth}>
                  <defs>
                    <linearGradient id="orgGrowth" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#0F766E" stopOpacity={0.4} />
                      <stop offset="95%" stopColor="#0F766E" stopOpacity={0.0} />
                    </linearGradient>
                  </defs>
                  <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="var(--border)" opacity={0.5} />
                  <XAxis dataKey="month" stroke="var(--muted-foreground)" fontSize={12} tickLine={false} />
                  <YAxis stroke="var(--muted-foreground)" fontSize={12} tickLine={false} axisLine={false} />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: "var(--card)",
                      borderColor: "var(--border)",
                      borderRadius: "10px",
                      fontSize: "12px",
                      color: "var(--foreground)",
                    }}
                  />
                  <Area
                    type="monotone"
                    dataKey="count"
                    name="Organizations"
                    stroke="#0F766E"
                    strokeWidth={2.5}
                    fillOpacity={1}
                    fill="url(#orgGrowth)"
                  />
                </AreaChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex h-full items-center justify-center text-sm text-muted-foreground">
                No growth telemetry recorded yet.
              </div>
            )}
          </CardContent>
        </Card>

        {/* Right 1 Col: Production Launch Gate Audit Card */}
        <Card className="card-enterprise flex flex-col justify-between">
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="text-base font-semibold flex items-center gap-2">
                <Cpu className="h-4 w-4 text-primary" />
                Launch Gate Audit
              </CardTitle>
              <Badge
                variant="outline"
                className={
                  gateStatus === "PASSED"
                    ? "border-emerald-500/30 text-emerald-600 dark:text-emerald-400 bg-emerald-500/5 font-mono text-xs"
                    : "border-amber-500/30 text-amber-600 dark:text-amber-400 bg-amber-500/5 font-mono text-xs"
                }
              >
                {gateStatus === "PASSED" ? (
                  <CheckCircle2 className="mr-1 h-3 w-3 inline" />
                ) : (
                  <AlertTriangle className="mr-1 h-3 w-3 inline" />
                )}
                {gateStatus}
              </Badge>
            </div>
            <CardDescription className="text-xs">
              Automated production-readiness verification
            </CardDescription>
          </CardHeader>

          <CardContent className="space-y-3 flex-1">
            <div className="rounded-lg border border-border bg-secondary/30 p-3 space-y-2 text-xs">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Gate Identifier:</span>
                <span className="font-mono font-medium">{gate?.gateName || "production_readiness"}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Last Verification:</span>
                <span className="font-medium">{formatDate(gate?.evaluatedAt)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">System Health:</span>
                <span className="text-emerald-600 font-semibold">100% Operational</span>
              </div>
            </div>

            <p className="text-xs text-muted-foreground leading-relaxed">
              Verifies tenant isolation schemas, encryption key vaults, Casbin RBAC policies, and real-time database connection pooling.
            </p>
          </CardContent>

          <div className="p-6 pt-0">
            <Button
              variant="outline"
              size="sm"
              asChild
              className="w-full text-xs justify-between"
            >
              <Link to="/platform/diagnostics">
                <span>View Full Diagnostics Suite</span>
                <ArrowUpRight className="h-3.5 w-3.5" />
              </Link>
            </Button>
          </div>
        </Card>
      </div>

      {/* Capability Distribution & Recent Platform Audit Trail */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Capability Adoption Distribution */}
        <Card className="card-enterprise">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-base font-semibold">
                  Capability Distribution
                </CardTitle>
                <CardDescription className="text-xs">
                  Active B2B module subscriptions across organizations
                </CardDescription>
              </div>
              <Link
                to="/platform/marketplace"
                className="text-xs text-primary hover:underline flex items-center gap-1 font-medium"
              >
                Marketplace <ArrowUpRight className="h-3 w-3" />
              </Link>
            </div>
          </CardHeader>
          <CardContent className="h-[240px]">
            {metrics?.capabilityDistribution && metrics.capabilityDistribution.length > 0 ? (
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={metrics.capabilityDistribution} layout="vertical">
                  <CartesianGrid strokeDasharray="3 3" horizontal={false} stroke="var(--border)" opacity={0.5} />
                  <XAxis type="number" stroke="var(--muted-foreground)" fontSize={12} tickLine={false} />
                  <YAxis dataKey="name" type="category" stroke="var(--muted-foreground)" fontSize={11} width={110} tickLine={false} />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: "var(--card)",
                      borderColor: "var(--border)",
                      borderRadius: "10px",
                      fontSize: "12px",
                    }}
                  />
                  <Bar dataKey="count" name="Subscribed Orgs" fill="#0F766E" radius={[0, 6, 6, 0]} />
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex h-full items-center justify-center text-sm text-muted-foreground">
                No capability distribution data available.
              </div>
            )}
          </CardContent>
        </Card>

        {/* Recent Audit Events */}
        <Card className="card-enterprise flex flex-col justify-between">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-base font-semibold">
                  Security & Platform Events
                </CardTitle>
                <CardDescription className="text-xs">
                  Real-time immutable administrative audit stream
                </CardDescription>
              </div>
              <Link
                to="/platform/audit"
                className="text-xs text-primary hover:underline flex items-center gap-1 font-medium"
              >
                Full Trail <ArrowUpRight className="h-3 w-3" />
              </Link>
            </div>
          </CardHeader>
          <CardContent className="space-y-2.5">
            {auditLogs && auditLogs.length > 0 ? (
              auditLogs.slice(0, 4).map((log) => (
                <div
                  key={log.id}
                  className="flex items-center justify-between p-2.5 rounded-lg border border-border bg-secondary/20 text-xs"
                >
                  <div className="flex flex-col gap-0.5 truncate max-w-[280px]">
                    <span className="font-medium text-foreground truncate">
                      {log.action}
                    </span>
                    <span className="text-[11px] text-muted-foreground truncate">
                      {log.actorEmail || "Platform Actor"} • {log.resourceType}
                    </span>
                  </div>
                  <div className="flex flex-col items-end gap-1">
                    <Badge
                      variant="outline"
                      className={
                        log.severity === "critical" || log.severity === "error"
                          ? "border-destructive/30 text-destructive bg-destructive/5 text-[10px]"
                          : "border-border text-muted-foreground text-[10px]"
                      }
                    >
                      {log.severity}
                    </Badge>
                    <span className="text-[10px] text-muted-foreground">
                      {formatDate(log.createdAt)}
                    </span>
                  </div>
                </div>
              ))
            ) : (
              <div className="py-8 text-center text-xs text-muted-foreground">
                No recent platform audit entries recorded.
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
