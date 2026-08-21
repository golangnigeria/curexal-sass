import React, { useState } from "react";
import {
  History,
  Shield,
  AlertTriangle,
  Users,
  Search,
  Filter,
  CheckCircle2,
  XCircle,
  FileCode,
  Calendar,
  Lock,
} from "lucide-react";
import {
  usePlatformAuditLogs,
  useAuditStats,
} from "@/api/hooks/use-audit-logs";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { formatDate } from "@/lib/utils";
import type { AuditLog } from "@/api/contracts";

export default function AuditLogsPage() {
  const [search, setSearch] = useState("");
  const [category, setCategory] = useState<string | undefined>(undefined);
  const [severity, setSeverity] = useState<string | undefined>(undefined);
  const [status, setStatus] = useState<string | undefined>(undefined);

  const [inspectLog, setInspectLog] = useState<AuditLog | null>(null);

  const { data: stats } = useAuditStats();
  const { data: logs, isLoading } = usePlatformAuditLogs({
    search: search || undefined,
    category: category === "all" ? undefined : category,
    severity: severity === "all" ? undefined : severity,
    status: status === "all" ? undefined : status,
    limit: 50,
  });

  const getSeverityBadge = (sev: string) => {
    switch (sev?.toLowerCase()) {
      case "critical":
      case "error":
        return (
          <Badge variant="outline" className="border-destructive/30 text-destructive bg-destructive/5 text-[10px]">
            {sev}
          </Badge>
        );
      case "warn":
      case "warning":
        return (
          <Badge variant="outline" className="border-amber-500/30 text-amber-600 bg-amber-500/5 text-[10px]">
            {sev}
          </Badge>
        );
      default:
        return (
          <Badge variant="outline" className="border-border text-muted-foreground text-[10px]">
            {sev || "info"}
          </Badge>
        );
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <div className="flex items-center gap-2">
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            Platform Audit Trail & Telemetry
          </h1>
          <Badge variant="outline" className="border-primary/30 text-primary bg-primary/5 text-xs font-mono">
            Immutable
          </Badge>
        </div>
        <p className="text-sm text-muted-foreground mt-1">
          Cryptographically auditable record of all administrative operations, logins, capability grants, and policy modifications.
        </p>
      </div>

      {/* KPI Stats */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <Card className="card-enterprise">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-xs font-medium text-muted-foreground uppercase">
              Total Logged Events
            </CardTitle>
            <History className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-foreground">
              {stats?.totalEvents ?? logs?.length ?? 0}
            </div>
            <p className="text-xs text-muted-foreground mt-1">Platform-wide events</p>
          </CardContent>
        </Card>

        <Card className="card-enterprise">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-xs font-medium text-muted-foreground uppercase">
              Security Incidents
            </CardTitle>
            <AlertTriangle className="h-4 w-4 text-amber-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-amber-600">
              {stats?.recentSecurityIncidents ?? 0}
            </div>
            <p className="text-xs text-muted-foreground mt-1">Rate limits & lockouts</p>
          </CardContent>
        </Card>

        <Card className="card-enterprise">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-xs font-medium text-muted-foreground uppercase">
              Active Actors
            </CardTitle>
            <Users className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-foreground">
              {stats?.activeActors ?? 1}
            </div>
            <p className="text-xs text-muted-foreground mt-1">Administrators & staff</p>
          </CardContent>
        </Card>

        <Card className="card-enterprise">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-xs font-medium text-muted-foreground uppercase">
              Error / Critical Events
            </CardTitle>
            <Shield className="h-4 w-4 text-emerald-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-foreground">
              {stats?.errorEvents ?? 0}
            </div>
            <p className="text-xs text-muted-foreground mt-1">System faults recorded</p>
          </CardContent>
        </Card>
      </div>

      {/* Filter & Search Bar */}
      <Card className="card-enterprise p-4">
        <div className="flex flex-col sm:flex-row items-center gap-3">
          <div className="relative flex-1 w-full">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search audit logs by action, actor, resource..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-9 text-xs h-9"
            />
          </div>

          <div className="flex items-center gap-2 w-full sm:w-auto">
            <Select value={severity || "all"} onValueChange={setSeverity}>
              <SelectTrigger className="w-[130px] text-xs h-9">
                <SelectValue placeholder="All Severities" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Severities</SelectItem>
                <SelectItem value="info">Info</SelectItem>
                <SelectItem value="warn">Warn</SelectItem>
                <SelectItem value="error">Error</SelectItem>
                <SelectItem value="critical">Critical</SelectItem>
              </SelectContent>
            </Select>

            <Select value={status || "all"} onValueChange={setStatus}>
              <SelectTrigger className="w-[120px] text-xs h-9">
                <SelectValue placeholder="All Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="success">Success</SelectItem>
                <SelectItem value="failure">Failure</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </Card>

      {/* Audit Logs Table */}
      <Card className="card-enterprise overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="border-b border-border bg-muted/40 font-semibold uppercase tracking-wider text-muted-foreground">
              <tr>
                <th className="py-3 px-4">Timestamp</th>
                <th className="py-3 px-4">Action</th>
                <th className="py-3 px-4">Actor</th>
                <th className="py-3 px-4">Category</th>
                <th className="py-3 px-4">Severity</th>
                <th className="py-3 px-4">Status</th>
                <th className="py-3 px-4 text-right">Details</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {isLoading ? (
                <tr>
                  <td colSpan={7} className="py-8 text-center text-muted-foreground">
                    Querying audit stream...
                  </td>
                </tr>
              ) : !logs || logs.length === 0 ? (
                <tr>
                  <td colSpan={7} className="py-8 text-center text-muted-foreground">
                    No platform audit events matching query.
                  </td>
                </tr>
              ) : (
                logs.map((log) => (
                  <tr key={log.id} className="hover:bg-muted/20">
                    <td className="py-3 px-4 font-mono text-muted-foreground">
                      {formatDate(log.createdAt)}
                    </td>
                    <td className="py-3 px-4 font-semibold text-foreground">
                      {log.action}
                    </td>
                    <td className="py-3 px-4">
                      <div className="flex flex-col">
                        <span className="font-medium text-foreground">
                          {log.actorEmail || "Platform Service"}
                        </span>
                        <span className="text-[10px] text-muted-foreground font-mono">
                          {log.actorRole || "system"}
                        </span>
                      </div>
                    </td>
                    <td className="py-3 px-4 font-mono text-[11px] text-muted-foreground">
                      {log.category || "security"}
                    </td>
                    <td className="py-3 px-4">
                      {getSeverityBadge(log.severity)}
                    </td>
                    <td className="py-3 px-4">
                      {log.status === "failure" ? (
                        <span className="text-destructive font-medium flex items-center gap-1">
                          <XCircle className="h-3 w-3" /> Failed
                        </span>
                      ) : (
                        <span className="text-emerald-600 font-medium flex items-center gap-1">
                          <CheckCircle2 className="h-3 w-3" /> Success
                        </span>
                      )}
                    </td>
                    <td className="py-3 px-4 text-right">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => setInspectLog(log)}
                        className="h-7 text-xs gap-1"
                      >
                        <FileCode className="h-3 w-3" /> Inspect
                      </Button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </Card>

      {/* Inspect Log Payload Modal */}
      <Dialog open={!!inspectLog} onOpenChange={(open) => !open && setInspectLog(null)}>
        <DialogContent className="sm:max-w-xl">
          <DialogHeader>
            <DialogTitle className="text-base font-semibold flex items-center gap-2">
              <History className="h-4 w-4 text-primary" />
              Audit Log Event Details
            </DialogTitle>
            <DialogDescription className="text-xs font-mono">
              Event ID: {inspectLog?.id}
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-3 py-2 text-xs">
            <div className="grid grid-cols-2 gap-2 rounded-lg border border-border bg-secondary/30 p-3">
              <div>
                <span className="text-muted-foreground">Action:</span>{" "}
                <span className="font-semibold text-foreground">{inspectLog?.action}</span>
              </div>
              <div>
                <span className="text-muted-foreground">Category:</span>{" "}
                <span className="font-mono">{inspectLog?.category}</span>
              </div>
              <div>
                <span className="text-muted-foreground">Actor:</span>{" "}
                <span>{inspectLog?.actorEmail || "Platform"}</span>
              </div>
              <div>
                <span className="text-muted-foreground">IP Address:</span>{" "}
                <span className="font-mono">{inspectLog?.ipAddress || "—"}</span>
              </div>
            </div>

            <div className="space-y-1">
              <Label className="text-xs font-semibold">Structured Audit Payload</Label>
              <div className="rounded-lg border border-border bg-muted/40 p-3 max-h-[250px] overflow-y-auto">
                <pre className="text-[11px] font-mono text-foreground leading-relaxed whitespace-pre-wrap">
                  {JSON.stringify(inspectLog?.details || inspectLog, null, 2)}
                </pre>
              </div>
            </div>
          </div>

          <DialogFooter>
            <Button size="sm" onClick={() => setInspectLog(null)} className="text-xs">
              Close
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
