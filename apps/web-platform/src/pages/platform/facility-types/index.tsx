import React, { useState } from "react";
import {
  Layers,
  Plus,
  Edit2,
  FileCode,
  CheckCircle2,
  ExternalLink,
  Code,
  Layout,
  ListOrdered,
  Eye,
} from "lucide-react";
import {
  useFacilityTypes,
  useCreateFacilityType,
  useUpdateFacilityType,
  useFacilityTypeRegistrationForm,
  useFacilityTypeNavigation,
  useFacilityTypeSetupSteps,
  useFacilityTypeDashboard,
} from "@/api/hooks/use-facility-types";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import type { FacilityTypeEntity } from "@/api/contracts";
import { toast } from "sonner";

export default function FacilityTypesPage() {
  const { data: types, isLoading, refetch } = useFacilityTypes();
  const createMutation = useCreateFacilityType();

  // Create Modal State
  const [openCreate, setOpenCreate] = useState(false);
  const [code, setCode] = useState("");
  const [name, setName] = useState("");
  const [category, setCategory] = useState("clinical");
  const [iconKey, setIconKey] = useState("Activity");
  const [description, setDescription] = useState("");
  const [isActive, setIsActive] = useState(true);

  // Edit Modal State
  const [editingType, setEditingType] = useState<FacilityTypeEntity | null>(null);
  const updateMutation = useUpdateFacilityType(editingType?.id || "");

  // Schema Viewer State
  const [viewerTypeId, setViewerTypeId] = useState<string | null>(null);
  const [viewerMode, setViewerMode] = useState<"form" | "nav" | "steps" | "dash">("form");

  const { data: formSchema } = useFacilityTypeRegistrationForm(viewerTypeId || "");
  const { data: navSchema } = useFacilityTypeNavigation(viewerTypeId || "");
  const { data: stepsSchema } = useFacilityTypeSetupSteps(viewerTypeId || "");
  const { data: dashSchema } = useFacilityTypeDashboard(viewerTypeId || "");

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createMutation.mutateAsync({
        code,
        name,
        category,
        iconKey,
        description,
        isActive,
        version: 1,
      });
      toast.success(`Facility type '${name}' created!`);
      setOpenCreate(false);
      setCode("");
      setName("");
      setDescription("");
      refetch();
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Failed to create facility type");
    }
  };

  const handleOpenEdit = (t: FacilityTypeEntity) => {
    setEditingType(t);
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingType) return;
    try {
      await updateMutation.mutateAsync(editingType);
      toast.success(`Facility type '${editingType.name}' updated!`);
      setEditingType(null);
      refetch();
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Failed to update facility type");
    }
  };

  const getViewerContent = () => {
    switch (viewerMode) {
      case "form":
        return formSchema ? JSON.stringify(formSchema, null, 2) : "Loading form schema...";
      case "nav":
        return navSchema ? JSON.stringify(navSchema, null, 2) : "Loading navigation schema...";
      case "steps":
        return stepsSchema ? JSON.stringify(stepsSchema, null, 2) : "Loading setup steps schema...";
      case "dash":
        return dashSchema ? JSON.stringify(dashSchema, null, 2) : "Loading dashboard schema...";
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            Facility Types & Schema Definitions
          </h1>
          <p className="text-sm text-muted-foreground mt-1">
            Dynamic facility models, registration forms, onboarding wizards, and customized navigation blueprints.
          </p>
        </div>

        <Dialog open={openCreate} onOpenChange={setOpenCreate}>
          <DialogTrigger asChild>
            <Button className="bg-primary hover:bg-primary/90 text-primary-foreground text-xs h-9 gap-2 shadow-sm">
              <Plus className="h-4 w-4" />
              Add Facility Type
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle className="text-base">Create Facility Type</DialogTitle>
              <DialogDescription className="text-xs">
                Define a new healthcare operating facility archetype on the cluster.
              </DialogDescription>
            </DialogHeader>

            <form onSubmit={handleCreate} className="space-y-3 py-2 text-xs">
              <div className="space-y-1">
                <Label htmlFor="typeName">Facility Type Name</Label>
                <Input
                  id="typeName"
                  placeholder="e.g. Diagnostic Center"
                  value={name}
                  onChange={(e) => {
                    setName(e.target.value);
                    if (!code) {
                      setCode(e.target.value.toLowerCase().replace(/[^a-z0-9]+/g, "_"));
                    }
                  }}
                  required
                  className="text-xs"
                />
              </div>

              <div className="space-y-1">
                <Label htmlFor="typeCode">Code Identifier</Label>
                <Input
                  id="typeCode"
                  placeholder="diagnostic_center"
                  value={code}
                  onChange={(e) => setCode(e.target.value)}
                  required
                  className="text-xs font-mono"
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-1">
                  <Label htmlFor="category">Category Tier</Label>
                  <Select value={category} onValueChange={setCategory}>
                    <SelectTrigger id="category" className="text-xs">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="clinical">Clinical Tier</SelectItem>
                      <SelectItem value="diagnostic">Diagnostic Tier</SelectItem>
                      <SelectItem value="retail">Retail Tier</SelectItem>
                      <SelectItem value="research">Research Tier</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-1">
                  <Label htmlFor="iconKey">Icon Key</Label>
                  <Input
                    id="iconKey"
                    value={iconKey}
                    onChange={(e) => setIconKey(e.target.value)}
                    className="text-xs font-mono"
                  />
                </div>
              </div>

              <div className="space-y-1">
                <Label htmlFor="typeDesc">Description</Label>
                <Input
                  id="typeDesc"
                  placeholder="Clinical & imaging diagnosis workflows"
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  className="text-xs"
                />
              </div>

              <div className="flex items-center justify-between pt-2">
                <Label htmlFor="activeToggle" className="cursor-pointer">Active State</Label>
                <Switch
                  id="activeToggle"
                  checked={isActive}
                  onCheckedChange={setIsActive}
                />
              </div>

              <DialogFooter className="pt-3">
                <Button type="button" variant="outline" size="sm" onClick={() => setOpenCreate(false)} className="text-xs">
                  Cancel
                </Button>
                <Button type="submit" size="sm" className="bg-primary text-primary-foreground text-xs">
                  Create Type
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {/* Facility Types List Card */}
      <Card className="card-enterprise overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="border-b border-border bg-muted/40 font-semibold uppercase tracking-wider text-muted-foreground">
              <tr>
                <th className="py-3 px-4">Facility Type & Code</th>
                <th className="py-3 px-4">Tier Category</th>
                <th className="py-3 px-4">Description</th>
                <th className="py-3 px-4">Status</th>
                <th className="py-3 px-4">Version</th>
                <th className="py-3 px-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {isLoading ? (
                <tr>
                  <td colSpan={6} className="py-8 text-center text-muted-foreground">
                    Loading facility types...
                  </td>
                </tr>
              ) : !types || types.length === 0 ? (
                <tr>
                  <td colSpan={6} className="py-8 text-center text-muted-foreground">
                    No facility types configured.
                  </td>
                </tr>
              ) : (
                types.map((t) => (
                  <tr key={t.id || t.code} className="hover:bg-muted/20">
                    <td className="py-3 px-4">
                      <div className="flex flex-col">
                        <span className="font-semibold text-foreground">{t.name}</span>
                        <span className="font-mono text-[11px] text-primary">{t.code}</span>
                      </div>
                    </td>
                    <td className="py-3 px-4 capitalize">
                      <Badge variant="secondary" className="text-[10px] capitalize">
                        {t.category}
                      </Badge>
                    </td>
                    <td className="py-3 px-4 text-muted-foreground max-w-[240px] truncate">
                      {t.description || "—"}
                    </td>
                    <td className="py-3 px-4">
                      {t.isActive ? (
                        <Badge variant="outline" className="border-emerald-500/30 text-emerald-600 bg-emerald-500/5 text-[10px]">
                          Active
                        </Badge>
                      ) : (
                        <Badge variant="outline" className="border-border text-muted-foreground text-[10px]">
                          Disabled
                        </Badge>
                      )}
                    </td>
                    <td className="py-3 px-4 font-mono">v{t.version}</td>
                    <td className="py-3 px-4 text-right">
                      <div className="flex items-center justify-end gap-1.5">
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => {
                            setViewerTypeId(t.code || t.id || "");
                            setViewerMode("form");
                          }}
                          className="h-7 text-xs gap-1"
                        >
                          <FileCode className="h-3 w-3" /> Schemas
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleOpenEdit(t)}
                          className="h-7 text-xs"
                        >
                          <Edit2 className="h-3 w-3" />
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </Card>

      {/* Schema Viewer Dialog */}
      <Dialog open={!!viewerTypeId} onOpenChange={(open) => !open && setViewerTypeId(null)}>
        <DialogContent className="sm:max-w-2xl">
          <DialogHeader>
            <DialogTitle className="text-base font-semibold flex items-center gap-2">
              <Code className="h-4 w-4 text-primary" />
              Dynamic Schema Blueprint: <span className="font-mono text-primary">{viewerTypeId}</span>
            </DialogTitle>
            <DialogDescription className="text-xs">
              Inspect backend-generated dynamic schemas for this facility tier.
            </DialogDescription>
          </DialogHeader>

          {/* Schema Selector Tabs */}
          <div className="flex items-center gap-2 border-b border-border pb-2 pt-1 text-xs">
            <Button
              variant={viewerMode === "form" ? "default" : "outline"}
              size="sm"
              onClick={() => setViewerMode("form")}
              className="h-7 text-xs"
            >
              Registration Form
            </Button>
            <Button
              variant={viewerMode === "nav" ? "default" : "outline"}
              size="sm"
              onClick={() => setViewerMode("nav")}
              className="h-7 text-xs"
            >
              Navigation Menu
            </Button>
            <Button
              variant={viewerMode === "steps" ? "default" : "outline"}
              size="sm"
              onClick={() => setViewerMode("steps")}
              className="h-7 text-xs"
            >
              Setup Steps
            </Button>
            <Button
              variant={viewerMode === "dash" ? "default" : "outline"}
              size="sm"
              onClick={() => setViewerMode("dash")}
              className="h-7 text-xs"
            >
              Dashboard Widgets
            </Button>
          </div>

          <div className="rounded-lg border border-border bg-muted/40 p-3 max-h-[350px] overflow-y-auto">
            <pre className="text-[11px] font-mono text-foreground leading-relaxed whitespace-pre-wrap">
              {getViewerContent()}
            </pre>
          </div>

          <DialogFooter>
            <Button size="sm" onClick={() => setViewerTypeId(null)} className="text-xs">
              Close Blueprint
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
