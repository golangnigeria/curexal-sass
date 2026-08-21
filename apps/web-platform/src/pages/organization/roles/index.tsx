import React, { useState } from "react";
import { useOrgRoles, useCreateRole } from "@/api/hooks/use-organization";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { toast } from "sonner";
import { Shield, Plus, CheckCircle2, Lock, Users, Sparkles } from "lucide-react";

const availablePermissionGroups = [
  {
    group: "Patient Care & Reception",
    permissions: [
      { code: "workspace:patient:read", label: "View Patient Demographics & Records" },
      { code: "workspace:patient:create", label: "Register New Patients & Queue" },
      { code: "workspace:patient:update", label: "Edit Patient Files & Merge Records" },
    ],
  },
  {
    group: "Laboratory Information System (LIS)",
    permissions: [
      { code: "workspace:sample:receive", label: "Accession Specimen & Barcode Printing" },
      { code: "workspace:worksheet:update", label: "Enter & Batch Worklist Test Results" },
      { code: "workspace:result:authorize", label: "Two-Step Pathologist Verification" },
    ],
  },
  {
    group: "Clinical & Outpatient EMR",
    permissions: [
      { code: "workspace:clinical:consult", label: "Conduct Consultations & SOAP Notes" },
      { code: "workspace:prescription:create", label: "Electronic Prescribing" },
      { code: "workspace:vitals:record", label: "Record Nursing Vitals & Triage" },
    ],
  },
  {
    group: "Billing, Cashier & Financials",
    permissions: [
      { code: "workspace:billing:create", label: "Generate Invoices & Process Payments" },
      { code: "workspace:refund:approve", label: "Authorize Bill Adjustments & Refunds" },
      { code: "organization:billing:read", label: "View Financial & Revenue Telemetry" },
    ],
  },
];

export default function OrganizationRolesPage() {
  const { data: roles, isLoading } = useOrgRoles();
  const createRoleMutation = useCreateRole();

  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [roleCode, setRoleCode] = useState("");
  const [roleName, setRoleName] = useState("");
  const [description, setDescription] = useState("");
  const [selectedPermissions, setSelectedPermissions] = useState<string[]>([
    "workspace:patient:read",
    "workspace:sample:receive",
  ]);

  const togglePermission = (code: string) => {
    setSelectedPermissions((prev) =>
      prev.includes(code) ? prev.filter((p) => p !== code) : [...prev, code]
    );
  };

  const handleCreateRole = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!roleName.trim() || !roleCode.trim()) {
      toast.error("Please provide both role name and unique identifier.");
      return;
    }

    try {
      await createRoleMutation.mutateAsync({
        code: roleCode.trim().toLowerCase().replace(/\s+/g, "_"),
        name: roleName.trim(),
        description: description.trim(),
        permissions: selectedPermissions,
      });

      toast.success("Role Created Successfully!", {
        description: `${roleName} is now available for staff assignment.`,
      });

      setIsCreateOpen(false);
      setRoleCode("");
      setRoleName("");
      setDescription("");
    } catch (err: any) {
      toast.error("Failed to create role: " + (err.message || "Network error"));
    }
  };

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-6">
        <div>
          <div className="flex items-center gap-2.5 mb-1">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              Roles & Granular RBAC Permissions
            </h1>
            <Badge variant="outline" className="border-primary/40 text-primary bg-primary/10 text-[10px] font-mono">
              Enterprise Security
            </Badge>
          </div>
          <p className="text-xs text-muted-foreground">
            Enforce role-based access control across clinical, laboratory, radiology, pharmacy, and billing workspaces.
          </p>
        </div>

        <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
          <DialogTrigger asChild>
            <Button size="sm" className="text-xs h-9 gap-1.5 bg-primary text-primary-foreground shadow">
              <Plus className="w-3.5 h-3.5" />
              Create Custom Role
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-xl bg-card border-border max-h-[85vh] overflow-y-auto">
            <DialogHeader>
              <DialogTitle className="text-base font-bold">Create Custom Organizational Role</DialogTitle>
              <DialogDescription className="text-xs">
                Define a tailored role with granular permission bindings.
              </DialogDescription>
            </DialogHeader>

            <form onSubmit={handleCreateRole} className="space-y-4 py-2">
              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-1.5">
                  <Label htmlFor="roleName" className="text-xs font-medium">Role Name</Label>
                  <Input
                    id="roleName"
                    placeholder="e.g. Senior Phlebotomist"
                    value={roleName}
                    onChange={(e) => setRoleName(e.target.value)}
                    required
                    className="text-xs h-9"
                  />
                </div>
                <div className="space-y-1.5">
                  <Label htmlFor="roleCode" className="text-xs font-medium">Role Code (ID)</Label>
                  <Input
                    id="roleCode"
                    placeholder="e.g. senior_phlebotomist"
                    value={roleCode}
                    onChange={(e) => setRoleCode(e.target.value)}
                    required
                    className="text-xs h-9 font-mono"
                  />
                </div>
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="desc" className="text-xs font-medium">Description</Label>
                <Input
                  id="desc"
                  placeholder="Primary duties and scope of authority..."
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  className="text-xs h-9"
                />
              </div>

              <div className="space-y-3 pt-2">
                <Label className="text-xs font-semibold text-foreground">Granular Module Permissions</Label>
                <div className="space-y-4">
                  {availablePermissionGroups.map((group) => (
                    <div key={group.group} className="space-y-2 p-3 rounded-lg border border-border bg-secondary/10">
                      <p className="text-[11px] font-bold text-foreground uppercase tracking-wider">{group.group}</p>
                      <div className="space-y-1.5">
                        {group.permissions.map((p) => {
                          const isChecked = selectedPermissions.includes(p.code);
                          return (
                            <label
                              key={p.code}
                              className="flex items-center gap-2 text-xs cursor-pointer hover:text-foreground"
                            >
                              <input
                                type="checkbox"
                                checked={isChecked}
                                onChange={() => togglePermission(p.code)}
                                className="rounded border-input text-primary focus:ring-primary w-3.5 h-3.5"
                              />
                              <span className="text-muted-foreground">{p.label}</span>
                              <span className="text-[10px] font-mono text-muted-foreground/60">({p.code})</span>
                            </label>
                          );
                        })}
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              <DialogFooter className="pt-3">
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() => setIsCreateOpen(false)}
                  className="text-xs h-9"
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  size="sm"
                  disabled={createRoleMutation.isPending}
                  className="text-xs h-9 bg-primary text-primary-foreground gap-1.5"
                >
                  {createRoleMutation.isPending ? "Creating..." : "Create Role"}
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {/* Role Cards Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {[
          { code: "owner", name: "Organization Owner", desc: "Full executive control over billing, branches, and staff.", count: 1, system: true },
          { code: "org_admin", name: "Organization Administrator", desc: "Administrative access across all operational branches.", count: 3, system: true },
          { code: "pathologist", name: "Chief / Consultant Pathologist", desc: "Two-step laboratory result validation and QC authorization.", count: 2, system: false },
          { code: "radiologist", name: "Consultant Radiologist", desc: "DICOM PACS review and diagnostic imaging reports.", count: 1, system: false },
          { code: "lab_scientist", name: "Medical Laboratory Scientist", desc: "Accessioning, worklist execution, and analyzer review.", count: 8, system: false },
          { code: "cashier", name: "Billing & Cashier Officer", desc: "Patient invoicing, POS receipts, and HMO claims.", count: 4, system: false },
        ].map((r) => (
          <Card key={r.code} className="border-border shadow-sm hover:shadow-md transition-all flex flex-col justify-between">
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between">
                <div className="p-2 rounded-lg bg-primary/10 text-primary">
                  <Shield className="w-4 h-4" />
                </div>
                {r.system && (
                  <Badge variant="secondary" className="text-[9px] uppercase font-mono">
                    System Core
                  </Badge>
                )}
              </div>
              <CardTitle className="text-sm font-bold text-foreground mt-2">{r.name}</CardTitle>
              <CardDescription className="text-xs font-mono text-muted-foreground">{r.code}</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4 pt-0">
              <p className="text-xs text-muted-foreground min-h-[32px]">{r.desc}</p>
              <div className="pt-3 border-t border-border flex items-center justify-between text-xs">
                <span className="flex items-center gap-1.5 text-muted-foreground">
                  <Users className="w-3.5 h-3.5" />
                  {r.count} Staff Assigned
                </span>
                <span className="text-[11px] font-semibold text-primary">Active</span>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
