import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import {
  useOrgBranches,
  useCreateBranch,
  useUpdateBranch,
  useDeactivateBranch,
  useSetHeadquarters,
} from "@/api/hooks/use-organization";
import { BranchPayload } from "@/api/services/organization.service";
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
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Textarea } from "@/components/ui/textarea";
import { toast } from "sonner";
import {
  Building2,
  Plus,
  Search,
  CheckCircle2,
  Layers,
  MapPin,
  Phone,
  Mail,
  ArrowRight,
  Sparkles,
  Shield,
  Activity,
  Stethoscope,
  Pill,
  Microscope,
  MoreVertical,
  Pencil,
  Star,
  PowerOff,
  Clock,
  Check,
  AlertTriangle,
} from "lucide-react";

function resolveDefaultWorkspaceModule(facilityType?: string): string {
  const norm = (facilityType || "").toLowerCase();
  if (norm.includes("clinic") || norm.includes("outpatient") || norm.includes("emr")) {
    return "clinical";
  }
  if (norm.includes("pharmacy")) {
    return "pharmacy";
  }
  if (norm.includes("radiology") || norm.includes("imaging")) {
    return "radiology";
  }
  return "clinical";
}

const facilityTypeOptions = [
  {
    code: "clinic",
    name: "Outpatient Clinic & Practice",
    icon: Stethoscope,
    desc: "Primary healthcare, doctor consultations, nursing triage, and cashier billing",
    caps: ["clinical.basic", "core.billing", "core.patient"],
  },
  {
    code: "specialist_clinic",
    name: "Specialist Outpatient Clinic",
    icon: Activity,
    desc: "Consultant specialty care, electronic prescriptions, and outpatient procedural services",
    caps: ["clinical.basic", "core.billing", "core.patient"],
  },
  {
    code: "telehealth_center",
    name: "Telehealth & Virtual Care Center",
    icon: Stethoscope,
    desc: "Remote consultations, digital care desk, and online patient check-in",
    caps: ["clinical.basic", "telehealth.addon", "core.billing"],
  },
];

function getFacilityTypeIcon(codeOrName?: string) {
  const normalized = (codeOrName || "").toLowerCase();
  if (normalized.includes("lab")) return Microscope;
  if (normalized.includes("clinic")) return Stethoscope;
  if (normalized.includes("pharmacy")) return Pill;
  if (normalized.includes("hospital")) return Building2;
  return Activity;
}

function generateBranchCode(branchName: string): string {
  const cleaned = branchName.trim().toUpperCase();
  if (!cleaned) return "";
  const words = cleaned.split(/\s+/).filter(Boolean);
  if (words.length === 1) {
    return words[0].slice(0, 3) + "-01";
  }
  return words.map((w) => w[0]).join("").slice(0, 4) + "-01";
}

export default function OrganizationBranchesPage() {
  const navigate = useNavigate();
  const { data: bootstrap } = useBootstrap();
  const { data: branches, isLoading } = useOrgBranches();

  const createBranchMutation = useCreateBranch();
  const updateBranchMutation = useUpdateBranch();
  const deactivateBranchMutation = useDeactivateBranch();
  const setHeadquartersMutation = useSetHeadquarters();

  const limits = bootstrap?.limits || { maxBranches: 3, maxMembers: 25, storageGb: 50 };
  const orgPlan = bootstrap?.organization?.subscription || "pro";
  const defaultCurrency = bootstrap?.workspace?.currency || "NGN";

  const [searchQuery, setSearchQuery] = useState("");
  const [isCreateOpen, setIsCreateOpen] = useState(false);

  // Create Form State
  const [name, setName] = useState("");
  const [code, setCode] = useState("");
  const [isAutoCode, setIsAutoCode] = useState(true);
  const [facilityType, setFacilityType] = useState("diagnostic_center");
  const [isHeadquarters, setIsHeadquarters] = useState(false);
  const [currency, setCurrency] = useState(defaultCurrency);
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [address, setAddress] = useState("");
  const [city, setCity] = useState("");
  const [state, setState] = useState("");

  // Edit Modal State
  const [editModalOpen, setEditModalOpen] = useState(false);
  const [selectedBranch, setSelectedBranch] = useState<BranchPayload | null>(null);
  const [editName, setEditName] = useState("");
  const [editEmail, setEditEmail] = useState("");
  const [editPhone, setEditPhone] = useState("");
  const [editAddress, setEditAddress] = useState("");
  const [editCity, setEditCity] = useState("");
  const [editState, setEditState] = useState("");
  const [editStatus, setEditStatus] = useState<"ACTIVE" | "INACTIVE">("ACTIVE");

  // Confirmation Modal State
  const [confirmHqBranch, setConfirmHqBranch] = useState<BranchPayload | null>(null);
  const [confirmDeactivateBranch, setConfirmDeactivateBranch] = useState<BranchPayload | null>(null);

  const branchesCount = branches?.length || 0;
  const maxBranches = limits.maxBranches || 3;
  const isLimitReached = branchesCount >= maxBranches;

  const handleNameChange = (val: string) => {
    setName(val);
    if (isAutoCode) {
      setCode(generateBranchCode(val));
    }
  };

  const handleCreateBranch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !code.trim()) {
      toast.error("Please provide both branch name and identifier code.");
      return;
    }

    try {
      await createBranchMutation.mutateAsync({
        name: name.trim(),
        code: code.trim().toUpperCase(),
        facilityType,
        isHeadquarters,
        currency,
        email: email.trim() || undefined,
        phone: phone.trim() || undefined,
        address: address.trim() || undefined,
        city: city.trim() || undefined,
        state: state.trim() || undefined,
      });

      toast.success("Branch Facility Provisioned Successfully!", {
        description: `${name} (${code.toUpperCase()}) is now live in your network.`,
      });

      setIsCreateOpen(false);
      setName("");
      setCode("");
      setIsHeadquarters(false);
      setEmail("");
      setPhone("");
      setAddress("");
      setCity("");
      setState("");
    } catch (err: any) {
      toast.error("Failed to provision branch: " + (err.response?.data?.message || err.message || "Network error"));
    }
  };

  const handleOpenEdit = (branch: BranchPayload) => {
    setSelectedBranch(branch);
    setEditName(branch.name || "");
    setEditEmail(branch.email || "");
    setEditPhone(branch.phone || "");
    setEditAddress(branch.address || "");
    setEditCity(branch.city || "");
    setEditState(branch.state || "");
    setEditStatus((branch.status as "ACTIVE" | "INACTIVE") || "ACTIVE");
    setEditModalOpen(true);
  };

  const handleSaveEdit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedBranch) return;

    try {
      await updateBranchMutation.mutateAsync({
        branchId: selectedBranch.id,
        data: {
          name: editName.trim(),
          email: editEmail.trim() || undefined,
          phone: editPhone.trim() || undefined,
          address: editAddress.trim() || undefined,
          city: editCity.trim() || undefined,
          state: editState.trim() || undefined,
          status: editStatus,
          version: selectedBranch.version || 1,
        },
      });

      toast.success("Branch Facility Updated!", {
        description: `${editName} configuration has been saved.`,
      });
      setEditModalOpen(false);
    } catch (err: any) {
      if (err.response?.status === 409) {
        toast.error("Concurrency Conflict", {
          description: "This facility was modified by another administrator. Please refresh before saving.",
        });
      } else {
        toast.error("Failed to update branch: " + (err.response?.data?.message || err.message));
      }
    }
  };

  const handleSetHeadquarters = async () => {
    if (!confirmHqBranch) return;
    try {
      await setHeadquartersMutation.mutateAsync(confirmHqBranch.id);
      toast.success("Headquarters Reassigned Successfully!", {
        description: `${confirmHqBranch.name} is now the primary headquarters.`,
      });
      setConfirmHqBranch(null);
    } catch (err: any) {
      toast.error("Failed to set headquarters: " + (err.response?.data?.message || err.message));
    }
  };

  const handleDeactivate = async () => {
    if (!confirmDeactivateBranch) return;
    try {
      await deactivateBranchMutation.mutateAsync(confirmDeactivateBranch.id);
      toast.success("Branch Deactivated", {
        description: `${confirmDeactivateBranch.name} is now deactivated.`,
      });
      setConfirmDeactivateBranch(null);
    } catch (err: any) {
      toast.error("Failed to deactivate branch: " + (err.response?.data?.message || err.message));
    }
  };

  const filteredBranches = (branches || []).filter((b: any) => {
    const branchName = (b.name || "").toLowerCase();
    const branchCode = (b.code || "").toLowerCase();
    const branchFacilityType = (b.facilityTypeName || b.facilityTypeCode || b.facilityType || "").toLowerCase();
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
            Manage multi-facility clinical blueprints, headquarters governance, and operational quotas.
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
            <DialogContent className="max-w-xl bg-card border-border max-h-[90vh] overflow-y-auto">
              <DialogHeader>
                <DialogTitle className="text-base font-bold flex items-center gap-2">
                  <Building2 className="w-4 h-4 text-primary" />
                  Provision Branch Facility
                </DialogTitle>
                <DialogDescription className="text-xs">
                  Create a new physical or virtual clinical facility within your multi-tenant network.
                </DialogDescription>
              </DialogHeader>

              <form onSubmit={handleCreateBranch} className="space-y-4 py-2">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
                  <div className="space-y-1.5 md:col-span-2">
                    <Label htmlFor="branchName" className="text-xs font-semibold">
                      Facility Name <span className="text-destructive">*</span>
                    </Label>
                    <Input
                      id="branchName"
                      placeholder="e.g. Victoria Island Diagnostic Center"
                      value={name}
                      onChange={(e) => handleNameChange(e.target.value)}
                      required
                      className="text-xs h-9"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="branchCode" className="text-xs font-semibold">
                      Identifier Code <span className="text-destructive">*</span>
                    </Label>
                    <Input
                      id="branchCode"
                      placeholder="e.g. VI-01"
                      value={code}
                      onChange={(e) => {
                        setCode(e.target.value.toUpperCase());
                        setIsAutoCode(false);
                      }}
                      required
                      className="text-xs h-9 font-mono uppercase"
                    />
                  </div>
                </div>

                <div className="space-y-1.5">
                  <Label className="text-xs font-semibold">Facility Operational Blueprint</Label>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 max-h-52 overflow-y-auto pr-1">
                    {facilityTypeOptions.map((opt) => {
                      const Icon = opt.icon;
                      const isSelected = facilityType === opt.code;
                      return (
                        <div
                          key={opt.code}
                          onClick={() => setFacilityType(opt.code)}
                          className={`p-2.5 rounded-lg border text-left cursor-pointer transition-all flex flex-col justify-between ${
                            isSelected
                              ? "border-primary bg-primary/5 shadow-xs"
                              : "border-border hover:border-border/80 bg-card"
                          }`}
                        >
                          <div className="flex items-center gap-2.5 mb-1.5">
                            <div className={`p-1.5 rounded-md ${isSelected ? "bg-primary text-white" : "bg-secondary text-muted-foreground"}`}>
                              <Icon className="w-3.5 h-3.5" />
                            </div>
                            <p className="text-xs font-bold text-foreground truncate">{opt.name}</p>
                            {isSelected && <CheckCircle2 className="w-3.5 h-3.5 text-primary ml-auto shrink-0" />}
                          </div>
                          <p className="text-[10px] text-muted-foreground line-clamp-1">{opt.desc}</p>
                        </div>
                      );
                    })}
                  </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  <div className="space-y-1.5">
                    <Label htmlFor="branchPhone" className="text-xs font-semibold">Official Contact Phone</Label>
                    <Input
                      id="branchPhone"
                      placeholder="+234 800 000 0000"
                      value={phone}
                      onChange={(e) => setPhone(e.target.value)}
                      className="text-xs h-9"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="branchEmail" className="text-xs font-semibold">Official Contact Email</Label>
                    <Input
                      id="branchEmail"
                      type="email"
                      placeholder="branch@organization.com"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      className="text-xs h-9"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                  <div className="space-y-1.5 sm:col-span-2">
                    <Label htmlFor="address" className="text-xs font-semibold">Physical Street Address</Label>
                    <Input
                      id="address"
                      placeholder="e.g. 12 Adeola Odeku St, Victoria Island"
                      value={address}
                      onChange={(e) => setAddress(e.target.value)}
                      className="text-xs h-9"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="city" className="text-xs font-semibold">City / State</Label>
                    <Input
                      id="city"
                      placeholder="Lagos, Nigeria"
                      value={city}
                      onChange={(e) => setCity(e.target.value)}
                      className="text-xs h-9"
                    />
                  </div>
                </div>

                <div className="flex items-center gap-2 p-3 rounded-lg bg-secondary/30 border border-border">
                  <input
                    type="checkbox"
                    id="isHq"
                    checked={isHeadquarters}
                    onChange={(e) => setIsHeadquarters(e.target.checked)}
                    className="rounded border-border text-primary focus:ring-primary h-4 w-4"
                  />
                  <Label htmlFor="isHq" className="text-xs cursor-pointer select-none">
                    <span className="font-semibold text-foreground">Set as Primary Corporate Headquarters</span>
                    <p className="text-[10px] text-muted-foreground">
                      Designates this facility as the central legal and administrative hub.
                    </p>
                  </Label>
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
          {filteredBranches.map((branch) => {
            const rawType = branch.facilityTypeName || branch.facilityTypeCode || branch.facilityType || "diagnostic_center";
            const displayType = rawType.replace(/_/g, " ");
            const FacilityIcon = getFacilityTypeIcon(rawType);
            const branchSlug = branch.slug || branch.code?.toLowerCase() || "main";
            const isActive = branch.status === "ACTIVE";

            return (
              <Card key={branch.id} className="border-border shadow-xs hover:shadow-md transition-all flex flex-col justify-between group bg-card">
                <CardHeader className="pb-3">
                  <div className="flex items-start justify-between">
                    <div className="flex items-center gap-2">
                      <div className="p-2 rounded-xl bg-primary/10 text-primary font-bold">
                        <FacilityIcon className="w-4 h-4" />
                      </div>
                      <div>
                        <Badge variant="outline" className="text-[9px] uppercase font-mono px-2 py-0.5 border-primary/30 text-primary">
                          {displayType}
                        </Badge>
                      </div>
                    </div>

                    <div className="flex items-center gap-1">
                      {branch.isHeadquarters && (
                        <Badge className="bg-amber-500/15 text-amber-600 dark:text-amber-400 border-amber-500/30 text-[9px] gap-1 font-bold">
                          <Star className="w-2.5 h-2.5 fill-amber-500" />
                          HQ
                        </Badge>
                      )}

                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon" className="h-7 w-7 text-muted-foreground hover:text-foreground">
                            <MoreVertical className="w-3.5 h-3.5" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end" className="text-xs w-48">
                          <DropdownMenuLabel>Facility Actions</DropdownMenuLabel>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem onClick={() => handleOpenEdit(branch)} className="gap-2 cursor-pointer">
                            <Pencil className="w-3.5 h-3.5" />
                            Edit Facility
                          </DropdownMenuItem>

                          {!branch.isHeadquarters && isActive && (
                            <DropdownMenuItem
                              onClick={() => setConfirmHqBranch(branch)}
                              className="gap-2 cursor-pointer text-amber-600 focus:text-amber-600"
                            >
                              <Star className="w-3.5 h-3.5" />
                              Set as Headquarters
                            </DropdownMenuItem>
                          )}

                          <DropdownMenuSeparator />
                          <DropdownMenuItem
                            onClick={() => setConfirmDeactivateBranch(branch)}
                            className="gap-2 cursor-pointer text-destructive focus:text-destructive"
                          >
                            <PowerOff className="w-3.5 h-3.5" />
                            Deactivate Facility
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </div>

                  <CardTitle className="text-base font-bold text-foreground mt-2.5 group-hover:text-primary transition-colors">
                    {branch.name}
                  </CardTitle>
                  <CardDescription className="text-xs flex items-center gap-1.5 font-mono">
                    Code: <span className="font-semibold text-foreground">{branch.code}</span> • Slug: <span className="text-muted-foreground">/{branchSlug}</span>
                  </CardDescription>
                </CardHeader>

                <CardContent className="space-y-4 pt-0">
                  <div className="space-y-1.5 text-xs text-muted-foreground">
                    {branch.address && (
                      <p className="flex items-center gap-1.5 truncate">
                        <MapPin className="w-3.5 h-3.5 text-muted-foreground shrink-0" />
                        {branch.address}{branch.city ? `, ${branch.city}` : ""}
                      </p>
                    )}
                    {branch.phone && (
                      <p className="flex items-center gap-1.5 truncate">
                        <Phone className="w-3.5 h-3.5 text-muted-foreground shrink-0" />
                        {branch.phone}
                      </p>
                    )}
                    {branch.email && (
                      <p className="flex items-center gap-1.5 truncate">
                        <Mail className="w-3.5 h-3.5 text-muted-foreground shrink-0" />
                        {branch.email}
                      </p>
                    )}
                  </div>

                  <div className="pt-3 border-t border-border flex items-center justify-between">
                    <span className={`flex items-center gap-1 text-[11px] font-semibold ${
                      isActive ? "text-emerald-600 dark:text-emerald-400" : "text-destructive"
                    }`}>
                      <span className={`w-1.5 h-1.5 rounded-full ${isActive ? "bg-emerald-500 animate-pulse" : "bg-destructive"}`} />
                      {isActive ? "Live Active" : "Deactivated"}
                    </span>

                    <Button asChild size="sm" variant="outline" className="text-xs h-8 gap-1 group-hover:bg-primary group-hover:text-primary-foreground transition-all">
                      <Link to={`/${branchSlug}/${resolveDefaultWorkspaceModule(branch.facilityTypeName || branch.facilityTypeCode || branch.facilityType)}`}>
                        Open Workspace <ArrowRight className="w-3 h-3" />
                      </Link>
                    </Button>
                  </div>
                </CardContent>
              </Card>
            );
          })}
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

      {/* Edit Branch Modal */}
      <Dialog open={editModalOpen} onOpenChange={setEditModalOpen}>
        <DialogContent className="max-w-md bg-card border-border">
          <DialogHeader>
            <DialogTitle className="text-base font-bold flex items-center gap-2">
              <Pencil className="w-4 h-4 text-primary" />
              Edit Facility Configuration
            </DialogTitle>
            <DialogDescription className="text-xs">
              Update operational parameters and contact channels for {selectedBranch?.name}.
            </DialogDescription>
          </DialogHeader>

          <form onSubmit={handleSaveEdit} className="space-y-3.5 py-2 text-xs">
            <div className="space-y-1">
              <Label htmlFor="editName" className="font-semibold">Facility Name</Label>
              <Input
                id="editName"
                value={editName}
                onChange={(e) => setEditName(e.target.value)}
                required
                className="text-xs"
              />
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="space-y-1">
                <Label htmlFor="editPhone" className="font-semibold">Contact Phone</Label>
                <Input
                  id="editPhone"
                  value={editPhone}
                  onChange={(e) => setEditPhone(e.target.value)}
                  className="text-xs"
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="editEmail" className="font-semibold">Contact Email</Label>
                <Input
                  id="editEmail"
                  type="email"
                  value={editEmail}
                  onChange={(e) => setEditEmail(e.target.value)}
                  className="text-xs"
                />
              </div>
            </div>

            <div className="space-y-1">
              <Label htmlFor="editAddress" className="font-semibold">Physical Street Address</Label>
              <Textarea
                id="editAddress"
                value={editAddress}
                onChange={(e) => setEditAddress(e.target.value)}
                rows={2}
                className="text-xs resize-none"
              />
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="space-y-1">
                <Label htmlFor="editCity" className="font-semibold">City</Label>
                <Input
                  id="editCity"
                  value={editCity}
                  onChange={(e) => setEditCity(e.target.value)}
                  className="text-xs"
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="editState" className="font-semibold">State / Region</Label>
                <Input
                  id="editState"
                  value={editState}
                  onChange={(e) => setEditState(e.target.value)}
                  className="text-xs"
                />
              </div>
            </div>

            <DialogFooter className="pt-3">
              <Button type="button" variant="outline" size="sm" onClick={() => setEditModalOpen(false)} className="text-xs">
                Cancel
              </Button>
              <Button
                type="submit"
                size="sm"
                disabled={updateBranchMutation.isPending}
                className="text-xs bg-primary text-primary-foreground"
              >
                {updateBranchMutation.isPending ? "Saving..." : "Save Changes"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Set Headquarters Confirmation Dialog */}
      <Dialog open={!!confirmHqBranch} onOpenChange={(open) => !open && setConfirmHqBranch(null)}>
        <DialogContent className="max-w-md bg-card border-border">
          <DialogHeader>
            <DialogTitle className="text-base font-bold flex items-center gap-2 text-amber-600">
              <Star className="w-4 h-4" />
              Reassign Corporate Headquarters
            </DialogTitle>
            <DialogDescription className="text-xs pt-1">
              Are you sure you want to designate <strong>{confirmHqBranch?.name}</strong> as the primary headquarters?
              This will update corporate governance, billing entity records, and regulatory filings.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="pt-3">
            <Button variant="outline" size="sm" onClick={() => setConfirmHqBranch(null)} className="text-xs">
              Cancel
            </Button>
            <Button
              size="sm"
              onClick={handleSetHeadquarters}
              disabled={setHeadquartersMutation.isPending}
              className="text-xs bg-amber-600 hover:bg-amber-700 text-white"
            >
              {setHeadquartersMutation.isPending ? "Reassigning..." : "Confirm HQ Reassignment"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Deactivate Branch Confirmation Dialog */}
      <Dialog open={!!confirmDeactivateBranch} onOpenChange={(open) => !open && setConfirmDeactivateBranch(null)}>
        <DialogContent className="max-w-md bg-card border-border">
          <DialogHeader>
            <DialogTitle className="text-base font-bold flex items-center gap-2 text-destructive">
              <AlertTriangle className="w-4 h-4" />
              Deactivate Facility Branch
            </DialogTitle>
            <DialogDescription className="text-xs pt-1">
              Are you sure you want to deactivate <strong>{confirmDeactivateBranch?.name}</strong>?
              Clinical staff and workstations will lose operational access to this branch workspace until reactivated.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="pt-3">
            <Button variant="outline" size="sm" onClick={() => setConfirmDeactivateBranch(null)} className="text-xs">
              Cancel
            </Button>
            <Button
              size="sm"
              variant="destructive"
              onClick={handleDeactivate}
              disabled={deactivateBranchMutation.isPending}
              className="text-xs"
            >
              {deactivateBranchMutation.isPending ? "Deactivating..." : "Deactivate Facility"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
