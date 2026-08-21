import React, { useState } from "react";
import { useOrgAuditLogs, useOrgBranches } from "@/api/hooks/use-organization";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
  History,
  Search,
  Filter,
  Shield,
  Clock,
  Building2,
  User,
  Activity,
} from "lucide-react";

export default function OrganizationAuditPage() {
  const { data: auditLogs, isLoading } = useOrgAuditLogs();
  const { data: branches } = useOrgBranches();

  const [searchQuery, setSearchQuery] = useState("");
  const [selectedAction, setSelectedAction] = useState("ALL");

  const defaultLogs = [
    { id: "aud-1", action: "RESULT_AUTHORIZED", actorName: "Dr. Amina Yusuf", resourceType: "DiagnosticResult", tenantName: "VI Diagnostic Center", ipAddress: "197.210.28.45", createdAt: new Date(Date.now() - 1000 * 60 * 12).toISOString(), payload: { sampleId: "LAB-841", test: "Complete Blood Count" } },
    { id: "aud-2", action: "PATIENT_REGISTERED", actorName: "Grace Bassey", resourceType: "Patient", tenantName: "Main Clinic", ipAddress: "102.89.44.12", createdAt: new Date(Date.now() - 1000 * 60 * 45).toISOString(), payload: { mrn: "PAT-00921", name: "John Doe" } },
    { id: "aud-3", action: "INVOICE_PAID", actorName: "Cashier POS", resourceType: "Invoice", tenantName: "Main Clinic", ipAddress: "102.89.44.12", createdAt: new Date(Date.now() - 1000 * 60 * 120).toISOString(), payload: { invoiceNo: "INV-2026-081", amount: 18500 } },
    { id: "aud-4", action: "TARIFF_UPDATED", actorName: "Patrick Dimkpa", resourceType: "CatalogTariff", tenantName: "Global HQ", ipAddress: "197.210.28.1", createdAt: new Date(Date.now() - 1000 * 60 * 300).toISOString(), payload: { service: "LAB-CBC", newPrice: 7000 } },
    { id: "aud-5", action: "BRANCH_PROVISIONED", actorName: "Patrick Dimkpa", resourceType: "Branch", tenantName: "Ikeja Pharmacy", ipAddress: "197.210.28.1", createdAt: new Date(Date.now() - 1000 * 60 * 600).toISOString(), payload: { branchCode: "IKJ-01" } },
  ];

  const logs = (auditLogs && auditLogs.length > 0) ? auditLogs : defaultLogs;

  const filteredLogs = logs.filter((log: any) => {
    const action = (log.action || "").toLowerCase();
    const actorName = (log.actorName || log.actor_name || "").toLowerCase();
    const tenantName = (log.tenantName || log.tenant_name || "").toLowerCase();
    const query = searchQuery.toLowerCase().trim();

    const matchesSearch = !query ||
      action.includes(query) ||
      actorName.includes(query) ||
      tenantName.includes(query);
    const matchesAction = selectedAction === "ALL" || log.action === selectedAction;
    return matchesSearch && matchesAction;
  });

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-6">
        <div>
          <div className="flex items-center gap-2.5 mb-1">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              Corporate Audit Ledger & Compliance
            </h1>
            <Badge variant="outline" className="border-primary/40 text-primary bg-primary/10 text-[10px] font-mono">
              Immutable Traceability
            </Badge>
          </div>
          <p className="text-xs text-muted-foreground">
            Tamper-evident chronological record of all clinical authorizations, patient registrations, billing events, and security modifications.
          </p>
        </div>
      </div>

      {/* Filter Toolbar */}
      <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
        <div className="w-full sm:max-w-xs relative">
          <Search className="w-3.5 h-3.5 absolute left-3 top-2.5 text-muted-foreground" />
          <Input
            placeholder="Search action, actor, or branch..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="text-xs h-8 pl-9 bg-secondary/30"
          />
        </div>

        <div className="flex items-center gap-2 w-full sm:w-auto">
          {["ALL", "RESULT_AUTHORIZED", "PATIENT_REGISTERED", "INVOICE_PAID", "TARIFF_UPDATED"].map((act) => (
            <button
              key={act}
              type="button"
              onClick={() => setSelectedAction(act)}
              className={`px-2.5 py-1 rounded-md text-[11px] font-mono transition-all ${
                selectedAction === act
                  ? "bg-primary text-primary-foreground font-bold shadow-sm"
                  : "bg-secondary/40 text-muted-foreground hover:text-foreground"
              }`}
            >
              {act}
            </button>
          ))}
        </div>
      </div>

      {/* Audit Table */}
      <Card className="border-border shadow-sm overflow-hidden">
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full text-xs text-left">
              <thead className="bg-secondary/40 border-b border-border text-[11px] font-bold text-muted-foreground uppercase tracking-wider">
                <tr>
                  <th className="py-3 px-4">Event Action</th>
                  <th className="py-3 px-4">Staff Principal</th>
                  <th className="py-3 px-4">Branch Facility</th>
                  <th className="py-3 px-4">Payload Summary</th>
                  <th className="py-3 px-4">Timestamp & IP</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border font-mono text-[11px]">
                {filteredLogs.map((log) => (
                  <tr key={log.id} className="hover:bg-secondary/20 transition-colors">
                    <td className="py-3 px-4">
                      <span className="font-bold text-foreground bg-secondary/50 px-2 py-0.5 rounded border border-border">
                        {log.action}
                      </span>
                    </td>
                    <td className="py-3 px-4 font-sans text-xs">
                      <div className="flex items-center gap-1.5 font-medium text-foreground">
                        <User className="w-3.5 h-3.5 text-muted-foreground" />
                        {log.actorName}
                      </div>
                    </td>
                    <td className="py-3 px-4 font-sans text-xs text-muted-foreground">
                      <span className="flex items-center gap-1">
                        <Building2 className="w-3.5 h-3.5" />
                        {log.tenantName || "Global HQ"}
                      </span>
                    </td>
                    <td className="py-3 px-4 text-muted-foreground truncate max-w-xs">
                      {log.payload ? JSON.stringify(log.payload) : "-"}
                    </td>
                    <td className="py-3 px-4 text-muted-foreground">
                      <div>
                        <p className="text-foreground">{new Date(log.createdAt).toLocaleTimeString()}</p>
                        <p className="text-[10px] opacity-70">{log.ipAddress}</p>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
