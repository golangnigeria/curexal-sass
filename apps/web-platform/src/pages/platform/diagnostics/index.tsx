import React from "react";
import {
  Cpu,
  ShieldCheck,
  CheckCircle2,
  AlertTriangle,
  Database,
  Activity,
  Server,
  RefreshCw,
  Zap,
  Lock,
} from "lucide-react";
import {
  useDiagnostics,
  useLaunchGateStatus,
  useHealthMetrics,
  useVerifyLaunchGate,
} from "@/api/hooks/use-diagnostics";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { formatDate } from "@/lib/utils";
import { toast } from "sonner";

export default function DiagnosticsPage() {
  const { data: diag, isLoading: diagLoading, refetch: refetchDiag } = useDiagnostics();
  const { data: gate, isLoading: gateLoading, refetch: refetchGate } = useLaunchGateStatus();
  const { data: health, isLoading: healthLoading } = useHealthMetrics();
  const verifyMutation = useVerifyLaunchGate();

  const handleVerify = async () => {
    try {
      await verifyMutation.mutateAsync();
      toast.success("Production Launch Gate verification suite executed successfully!");
      refetchGate();
      refetchDiag();
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Launch gate verification failed.");
    }
  };

  const gateStatus = gate?.status || "PASSED";
  const db = diag?.database;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              Diagnostics & Launch Gate
            </h1>
            <Badge variant="outline" className="border-primary/30 text-primary bg-primary/5 text-xs">
              System Telemetry
            </Badge>
          </div>
          <p className="text-sm text-muted-foreground mt-1">
            Production readiness compliance, database connection health, and distributed component telemetry.
          </p>
        </div>

        <div className="flex items-center gap-2.5">
          <Button
            variant="outline"
            size="sm"
            onClick={() => {
              refetchDiag();
              refetchGate();
            }}
            className="h-9 gap-2 text-xs"
          >
            <RefreshCw className={diagLoading ? "h-3.5 w-3.5 animate-spin" : "h-3.5 w-3.5"} />
            Refresh
          </Button>

          <Button
            size="sm"
            onClick={handleVerify}
            disabled={verifyMutation.isPending}
            className="bg-primary hover:bg-primary/90 text-primary-foreground text-xs h-9 gap-2 shadow-sm"
          >
            <Zap className="h-3.5 w-3.5" />
            {verifyMutation.isPending ? "Executing Verification Suite..." : "Execute Launch Gate"}
          </Button>
        </div>
      </div>

      {/* Production Readiness Launch Gate Banner */}
      <Card className="card-enterprise border-primary/30 bg-primary/5 p-6">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-6">
          <div className="space-y-2">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary text-white shadow-md shadow-primary/20">
                <ShieldCheck className="h-5 w-5" />
              </div>
              <div>
                <h2 className="text-lg font-bold text-foreground">
                  Production Launch Gate: {gate?.gateName || "production_readiness"}
                </h2>
                <p className="text-xs text-muted-foreground font-mono">
                  Last evaluated: {formatDate(gate?.evaluatedAt)}
                </p>
              </div>
            </div>
            <p className="text-xs text-muted-foreground leading-relaxed max-w-2xl">
              The automated launch gate executes end-to-end checks validating PostgreSQL connection pools, tenant schema isolation, symmetric token vault encryption, Casbin RBAC policies, and background notification workers.
            </p>
          </div>

          <div className="flex flex-col items-start md:items-end gap-2 shrink-0">
            <Badge
              variant="outline"
              className={
                gateStatus === "PASSED"
                  ? "border-emerald-500/40 text-emerald-600 dark:text-emerald-400 bg-emerald-500/10 text-sm font-mono px-3 py-1"
                  : "border-amber-500/40 text-amber-600 dark:text-amber-400 bg-amber-500/10 text-sm font-mono px-3 py-1"
              }
            >
              {gateStatus === "PASSED" ? (
                <CheckCircle2 className="mr-1.5 h-4 w-4 inline" />
              ) : (
                <AlertTriangle className="mr-1.5 h-4 w-4 inline" />
              )}
              STATUS: {gateStatus}
            </Badge>
            <span className="text-[11px] text-muted-foreground">
              Permission: <code className="font-mono text-primary">platform:launch_gate:execute</code>
            </span>
          </div>
        </div>
      </Card>

      {/* Cluster Health & DB Connection Pool */}
      <div className="grid gap-6 md:grid-cols-2">
        <Card className="card-enterprise">
          <CardHeader>
            <CardTitle className="text-sm font-semibold flex items-center gap-2">
              <Database className="h-4 w-4 text-emerald-500" />
              PostgreSQL Connection Pool Status
            </CardTitle>
            <CardDescription className="text-xs">
              Live pgxpool statistics from Go Echo kernel.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3 text-xs">
            <div className="flex justify-between border-b border-border/50 pb-2">
              <span className="text-muted-foreground">Cluster Status:</span>
              <Badge variant="outline" className="border-emerald-500/30 text-emerald-600 bg-emerald-500/5 font-mono text-[10px]">
                {db?.status || "connected"}
              </Badge>
            </div>
            <div className="flex justify-between border-b border-border/50 pb-2">
              <span className="text-muted-foreground">Total Open Connections:</span>
              <span className="font-mono font-bold text-foreground">{db?.openConnections ?? 0}</span>
            </div>
            <div className="flex justify-between border-b border-border/50 pb-2">
              <span className="text-muted-foreground">Acquired / In-Use:</span>
              <span className="font-mono font-bold text-primary">{db?.inUse ?? 0}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Idle in Pool:</span>
              <span className="font-mono font-bold text-foreground">{db?.idle ?? 0}</span>
            </div>
          </CardContent>
        </Card>

        <Card className="card-enterprise">
          <CardHeader>
            <CardTitle className="text-sm font-semibold flex items-center gap-2">
              <Server className="h-4 w-4 text-primary" />
              Server Instance & Cluster Telemetry
            </CardTitle>
            <CardDescription className="text-xs">
              Uptime and runtime environment parameters.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3 text-xs">
            <div className="flex justify-between border-b border-border/50 pb-2">
              <span className="text-muted-foreground">Server State:</span>
              <span className="font-semibold text-emerald-600 capitalize">{diag?.status || "operational"}</span>
            </div>
            <div className="flex justify-between border-b border-border/50 pb-2">
              <span className="text-muted-foreground">Cluster Uptime:</span>
              <span className="font-mono">{Math.floor((diag?.uptimeSeconds || 0) / 3600)}h {Math.floor(((diag?.uptimeSeconds || 0) % 3600) / 60)}m</span>
            </div>
            <div className="flex justify-between border-b border-border/50 pb-2">
              <span className="text-muted-foreground">Total Provisioned Organizations:</span>
              <span className="font-mono font-semibold">{diag?.metrics?.totalOrganizations ?? 0}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Active Workspaces / Facilities:</span>
              <span className="font-mono font-semibold">{diag?.metrics?.totalWorkspaces ?? 0}</span>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Component Health Breakdown Table */}
      <Card className="card-enterprise overflow-hidden">
        <CardHeader className="pb-3">
          <CardTitle className="text-sm font-semibold">Subsystem Health Status</CardTitle>
          <CardDescription className="text-xs">
            Monitored platform sub-services and background workers.
          </CardDescription>
        </CardHeader>
        <CardContent className="p-0">
          <table className="w-full text-left text-xs">
            <thead className="border-y border-border bg-muted/40 font-semibold uppercase tracking-wider text-muted-foreground">
              <tr>
                <th className="py-2.5 px-4">Component</th>
                <th className="py-2.5 px-4">Health Status</th>
                <th className="py-2.5 px-4">Telemetry Payload</th>
                <th className="py-2.5 px-4 text-right">Last Verified</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {healthLoading ? (
                <tr>
                  <td colSpan={4} className="py-8 text-center text-muted-foreground">
                    Querying component health telemetry...
                  </td>
                </tr>
              ) : !health || health.length === 0 ? (
                <tr>
                  <td colSpan={4} className="py-8 text-center text-muted-foreground">
                    No component health telemetry received from backend.
                  </td>
                </tr>
              ) : (
                health.map((h) => (
                  <tr key={h.id || h.componentName} className="hover:bg-muted/20">
                    <td className="py-3 px-4 font-medium text-foreground">{h.componentName}</td>
                    <td className="py-3 px-4">
                      <Badge
                        variant="outline"
                        className={
                          h.status === "HEALTHY"
                            ? "border-emerald-500/30 text-emerald-600 bg-emerald-500/5 text-[10px]"
                            : "border-destructive/30 text-destructive bg-destructive/5 text-[10px]"
                        }
                      >
                        {h.status}
                      </Badge>
                    </td>
                    <td className="py-3 px-4 font-mono text-[11px] text-muted-foreground">
                      {JSON.stringify(h.metrics || {})}
                    </td>
                    <td className="py-3 px-4 text-right text-muted-foreground">
                      {formatDate(h.checkedAt)}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </CardContent>
      </Card>
    </div>
  );
}
