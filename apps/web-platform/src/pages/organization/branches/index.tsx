import React, { useState } from "react";
import { Link } from "react-router-dom";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { useOrgBranches, useCreateBranch } from "@/api/hooks/use-organization";
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
import {
  Building2,
  Plus,
  Search,
  CheckCircle2,
  Layers,
  MapPin,
  Phone,
  ArrowRight,
  Sparkles,
  Shield,
  Activity,
  Stethoscope,
  Pill,
  Microscope,
} from "lucide-react";

const facilityTypeOptions = [
  { code: "diagnostic_center", name: "Diagnostic Center (Hybrid)", icon: Activity, desc: "Combined pathology & radiology imaging" },
  { code: "laboratory", name: "Medical Laboratory (LIS)", icon: Microscope, desc: "Clinical biochemistry, hematology & microbiology" },
  { code: "clinic", name: "Outpatient Clinic (EMR)", icon: Stethoscope, desc: "Primary care, consultations & triage" },
  { code: "pharmacy", name: "Community Pharmacy", icon: Pill, desc: "Prescription dispensing & drug inventory" },
  { code: "hospital", name: "Hospital (HIS)", icon: Building2, desc: "Inpatient wards, emergency & admissions" },
];

export default function OrganizationBranchesPage() {
  const { data: bootstrap } = useBootstrap();
  const { data: branches, isLoading } = useOrgBranches();
  const createBranchMutation = useCreateBranch();

  const limits = bootstrap?.limits || { maxBranches: 1, maxMembers: 5, storageGb: 10 };
  const orgPlan = bootstrap?.organization?.subscription || "smart";
  const defaultCurrency = bootstrap?.workspace?.currency || "NGN";

  const [searchQuery, setSearchQuery] = useState("");
  const [isCreateOpen, setIsCreateOpen] = useState(false);

  // Form State
  const [name, setName] = useState("");
  const [code, setCode] = useState("");
  const [facilityType, setFacilityType] = useState("diagnostic_center");
  const [currency, setCurrency] = useState(defaultCurrency);
  const [address, setAddress] = useState("");
  const [phone, setPhone] = useState("");

  const branchesCount = branches?.length || 0;
  const maxBranches = limits.maxBranches || 1;
  const isLimitReached = branchesCount >= maxBranches;

  const handleCreateBranch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !code.trim()) {
      toast.error("Please provide both branch name and branch code.");
      return;
    }

    try {
      await createBranchMutation.mutateAsync({
        name: name.trim(),
        code: code.trim().toUpperCase(),
        facilityType,
        currency,
        address: address.trim(),
        phone: phone.trim(),
      });

      toast.success("Branch Facility Provisioned Successfully!", {
        description: `${name} has been configured in your network.`,
      });

      setIsCreateOpen(false);
      setName("");
      setCode("");
      setAddress("");
      setPhone("");
    } catch (err: any) {
      toast.error("Failed to provision branch: " + (err.message || "Network error"));
    }
  };

  const filteredBranches = (branches || []).filter((b: any) => {
    const branchName = (b.name || "").toLowerCase();
    const branchCode = (b.code || "").toLowerCase();
    const branchFacilityType = (b.facilityType || b.facility_type || "").toLowerCase();
    const query = searchQuery.toLowerCase().trim();

    if (!query) return true;
    return (
      branchName.includes(query) ||
      branchCode.includes(query) ||
      branchFacilityType.includes(query)
    );
  });

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Header & Controls */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-6">
        <div>
          <div className="flex items-center gap-2.5 mb-1">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              Branch Facilities Network
            </h1>
            <Badge variant="outline" className="border-primary/40 text-primary bg-primary/10 text-[10px] font-mono uppercase font-bold">
              {branchesCount} of {maxBranches} Provisioned
            </Badge>
          </div>
          <p className="text-xs text-muted-foreground">
            Manage your multi-facility network, clinical blueprints, and facility operating quotas.
          </p>
        </div>

        <div className="flex items-center gap-3">
          <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
            <DialogTrigger asChild>
              <Button
                size="sm"
                className="text-xs h-9 gap-1.5 bg-primary text-primary-foreground shadow"
                disabled={isLimitReached}
              >
                <Plus className="w-3.5 h-3.5" />
                Provision New Branch
              </Button>
            </DialogTrigger>
            <DialogContent className="max-w-md bg-card border-border">
              <DialogHeader>
                <DialogTitle className="text-base font-bold">Provision Branch Facility</DialogTitle>
                <DialogDescription className="text-xs">
                  Create a new physical or virtual clinical branch within your organization network.
                </DialogDescription>
              </DialogHeader>

              <form onSubmit={handleCreateBranch} className="space-y-4 py-2">
                <div className="space-y-1.5">
                  <Label htmlFor="branchName" className="text-xs font-medium">Branch Facility Name</Label>
                  <Input
                    id="branchName"
                    placeholder="e.g. Victoria Island Diagnostic Center"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    required
                    className="text-xs h-9"
                  />
                </div>

                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1.5">
                    <Label htmlFor="branchCode" className="text-xs font-medium">Branch Code (Identifier)</Label>
                    <Input
                      id="branchCode"
                      placeholder="e.g. VI-01"
                      value={code}
                      onChange={(e) => setCode(e.target.value)}
                      required
                      className="text-xs h-9 font-mono uppercase"
                    />
                  </div>
                  <div className="space-y-1.5">
                    <Label htmlFor="currency" className="text-xs font-medium">Operating Currency</Label>
                    <Input
                      id="currency"
                      value={currency}
                      onChange={(e) => setCurrency(e.target.value)}
                      required
                      className="text-xs h-9 font-mono"
                    />
                  </div>
                </div>

                <div className="space-y-1.5">
                  <Label className="text-xs font-medium">Facility Blueprint</Label>
                  <div className="grid grid-cols-1 gap-2 max-h-48 overflow-y-auto pr-1">
                    {facilityTypeOptions.map((opt) => {
                      const Icon = opt.icon;
                      const isSelected = facilityType === opt.code;
                      return (
                        <div
                          key={opt.code}
                          onClick={() => setFacilityType(opt.code)}
                          className={`p-2.5 rounded-lg border text-left cursor-pointer transition-all flex items-center gap-3 ${
                            isSelected
                              ? "border-primary bg-primary/5 shadow-sm"
                              : "border-border hover:border-border/80 bg-card"
                          }`}
                        >
                          <div className={`p-2 rounded-lg ${isSelected ? "bg-primary text-white" : "bg-secondary text-muted-foreground"}`}>
                            <Icon className="w-4 h-4" />
                          </div>
                          <div className="flex-1 min-w-0">
                            <p className="text-xs font-semibold text-foreground truncate">{opt.name}</p>
                            <p className="text-[10px] text-muted-foreground truncate">{opt.desc}</p>
                          </div>
                          {isSelected && <CheckCircle2 className="w-4 h-4 text-primary shrink-0" />}
                        </div>
                      );
                    })}
                  </div>
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="address" className="text-xs font-medium">Physical Address</Label>
                  <Input
                    id="address"
                    placeholder="e.g. 12 Adeola Odeku St, Victoria Island, Lagos"
                    value={address}
                    onChange={(e) => setAddress(e.target.value)}
                    className="text-xs h-9"
                  />
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
                    disabled={createBranchMutation.isPending}
                    className="text-xs h-9 bg-primary text-primary-foreground gap-1.5"
                  >
                    {createBranchMutation.isPending ? "Provisioning..." : "Provision Facility"}
                  </Button>
                </DialogFooter>
              </form>
            </DialogContent>
          </Dialog>
        </div>
      </div>

      {/* Plan Limit Warning Banner if Reached */}
      {isLimitReached && (
        <div className="p-4 rounded-xl border border-amber-500/30 bg-amber-500/10 text-amber-900 dark:text-amber-200 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Sparkles className="w-5 h-5 text-amber-500 shrink-0" />
            <div>
              <p className="text-xs font-bold">Branch Quota Limit Reached ({branchesCount} / {maxBranches})</p>
              <p className="text-[11px] text-amber-800/80 dark:text-amber-300/80">
                You have utilized all allocated branch facilities on your {orgPlan} plan. Upgrade to expand your network.
              </p>
            </div>
          </div>
          <Button asChild size="sm" className="text-xs h-8 bg-amber-600 hover:bg-amber-700 text-white">
            <Link to="/organization/billing">Upgrade Plan</Link>
          </Button>
        </div>
      )}

      {/* Search Bar */}
      <div className="max-w-sm relative">
        <Search className="w-3.5 h-3.5 absolute left-3 top-2.5 text-muted-foreground" />
        <Input
          placeholder="Filter branches by name, code, or type..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="text-xs h-8 pl-9 bg-secondary/30"
        />
      </div>

      {/* Branch Cards Grid */}
      {isLoading ? (
        <div className="py-16 text-center text-xs text-muted-foreground animate-pulse">
          Loading branch facilities network...
        </div>
      ) : filteredBranches.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredBranches.map((branch) => (
            <Card key={branch.id} className="border-border shadow-sm hover:shadow-md transition-all flex flex-col justify-between group">
              <CardHeader className="pb-3">
                <div className="flex items-start justify-between">
                  <div className="p-2.5 rounded-xl bg-primary/10 text-primary font-bold">
                    <Building2 className="w-5 h-5" />
                  </div>
                  <Badge variant="outline" className="text-[9px] uppercase font-mono px-2 py-0.5 border-primary/30 text-primary">
                    {branch.facilityType.replace("_", " ")}
                  </Badge>
                </div>
                <CardTitle className="text-base font-bold text-foreground mt-2 group-hover:text-primary transition-colors">
                  {branch.name}
                </CardTitle>
                <CardDescription className="text-xs flex items-center gap-1 font-mono">
                  Code: <span className="font-semibold text-foreground">{branch.code}</span> • Currency: {branch.currency}
                </CardDescription>
              </CardHeader>

              <CardContent className="space-y-4 pt-0">
                <div className="space-y-1.5 text-xs text-muted-foreground">
                  {branch.address && (
                    <p className="flex items-center gap-1.5 truncate">
                      <MapPin className="w-3.5 h-3.5 text-muted-foreground shrink-0" />
                      {branch.address}
                    </p>
                  )}
                  {branch.phone && (
                    <p className="flex items-center gap-1.5 truncate">
                      <Phone className="w-3.5 h-3.5 text-muted-foreground shrink-0" />
                      {branch.phone}
                    </p>
                  )}
                </div>

                <div className="pt-3 border-t border-border flex items-center justify-between">
                  <span className="flex items-center gap-1 text-[11px] font-semibold text-emerald-600 dark:text-emerald-400">
                    <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
                    Live Active
                  </span>

                  <Button asChild size="sm" variant="outline" className="text-xs h-8 gap-1 group-hover:bg-primary group-hover:text-primary-foreground transition-all">
                    <Link to="/workspace/dashboard">
                      Open Workspace <ArrowRight className="w-3 h-3" />
                    </Link>
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      ) : (
        <div className="p-12 rounded-2xl border border-dashed border-border text-center space-y-3 bg-card">
          <Building2 className="w-10 h-10 mx-auto text-muted-foreground opacity-50" />
          <h3 className="text-sm font-semibold text-foreground">No branch facilities match your filter</h3>
          <p className="text-xs text-muted-foreground max-w-sm mx-auto">
            Try adjusting your search criteria or provision a new branch facility.
          </p>
        </div>
      )}
    </div>
  );
}
