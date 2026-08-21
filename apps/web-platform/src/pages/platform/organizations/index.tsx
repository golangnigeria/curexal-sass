import React, { useState } from "react";
import { Link } from "react-router-dom";
import {
  Building2,
  Search,
  Plus,
  ArrowUpRight,
  ShieldCheck,
  CheckCircle2,
  XCircle,
  Clock,
  Filter,
  MoreVertical,
  ExternalLink,
} from "lucide-react";
import {
  usePlatformOrganizations,
  useCreateOrganization,
} from "@/api/hooks/use-organizations";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { formatDate } from "@/lib/utils";
import { toast } from "sonner";

export default function OrganizationsPage() {
  const { data: orgs, isLoading, refetch } = usePlatformOrganizations();
  const createMutation = useCreateOrganization();

  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState("all");
  const [planFilter, setPlanFilter] = useState("all");

  // Create Modal State
  const [openCreate, setOpenCreate] = useState(false);
  const [name, setName] = useState("");
  const [slug, setSlug] = useState("");
  const [ownerEmail, setOwnerEmail] = useState("");
  const [ownerName, setOwnerName] = useState("");
  const [plan, setPlan] = useState("smart");

  const filteredOrgs = (orgs || []).filter((org) => {
    const matchesSearch =
      (org.name || "").toLowerCase().includes(search.toLowerCase()) ||
      (org.slug || "").toLowerCase().includes(search.toLowerCase());
    const matchesStatus = statusFilter === "all" || (org.status || "").toLowerCase() === statusFilter.toLowerCase();
    const matchesPlan = planFilter === "all" || (org.plan || "").toLowerCase() === planFilter.toLowerCase();
    return matchesSearch && matchesStatus && matchesPlan;
  });

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createMutation.mutateAsync({
        name,
        slug,
        owner: {
          email: ownerEmail,
          name: ownerName,
        },
        plan,
      });
      toast.success(`Organization '${name}' created! An invitation email has been sent to the owner.`);
      setOpenCreate(false);
      setName("");
      setSlug("");
      setOwnerEmail("");
      setOwnerName("");
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Failed to create organization");
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status?.toLowerCase()) {
      case "active":
        return (
          <Badge variant="outline" className="border-emerald-500/30 text-emerald-600 dark:text-emerald-400 bg-emerald-500/5 text-xs">
            <CheckCircle2 className="mr-1 h-3 w-3 inline" /> Active
          </Badge>
        );
      case "pending_verification":
        return (
          <Badge variant="outline" className="border-amber-500/30 text-amber-600 dark:text-amber-400 bg-amber-500/5 text-xs">
            <Clock className="mr-1 h-3 w-3 inline" /> Pending Review
          </Badge>
        );
      case "suspended":
        return (
          <Badge variant="outline" className="border-destructive/30 text-destructive bg-destructive/5 text-xs">
            <XCircle className="mr-1 h-3 w-3 inline" /> Suspended
          </Badge>
        );
      default:
        return <Badge variant="outline">{status || "Unknown"}</Badge>;
    }
  };

  return (
    <div className="space-y-6">
      {/* Header with Title & Action */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            Healthcare Organizations
          </h1>
          <p className="text-sm text-muted-foreground mt-1">
            Directory of all provisioned healthcare enterprise networks, plans, and compliance states.
          </p>
        </div>

        <Dialog open={openCreate} onOpenChange={setOpenCreate}>
          <DialogTrigger asChild>
            <Button className="bg-primary hover:bg-primary/90 text-primary-foreground text-xs h-9 gap-2 shadow-sm">
              <Plus className="h-4 w-4" />
              Provision Organization
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle className="text-base font-semibold">
                Provision New Organization
              </DialogTitle>
              <DialogDescription className="text-xs">
                Creates a new top-level healthcare network tenant on the cluster.
              </DialogDescription>
            </DialogHeader>

            <form onSubmit={handleCreate} className="space-y-4 py-2">
              <div className="space-y-1.5">
                <Label htmlFor="orgName" className="text-xs">Organization Name</Label>
                <Input
                  id="orgName"
                  placeholder="e.g. Apex Diagnostics Network"
                  value={name}
                  onChange={(e) => {
                    setName(e.target.value);
                    if (!slug) {
                      setSlug(e.target.value.toLowerCase().replace(/[^a-z0-9]+/g, "-"));
                    }
                  }}
                  required
                  className="text-xs"
                />
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="slug" className="text-xs">Subdomain / Tenant Slug</Label>
                <Input
                  id="slug"
                  placeholder="apex-diagnostics"
                  value={slug}
                  onChange={(e) => setSlug(e.target.value)}
                  required
                  className="text-xs font-mono"
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-1.5">
                  <Label htmlFor="ownerName" className="text-xs">Owner Name</Label>
                  <Input
                    id="ownerName"
                    placeholder="Dr. John Doe"
                    value={ownerName}
                    onChange={(e) => setOwnerName(e.target.value)}
                    required
                    className="text-xs"
                  />
                </div>
                <div className="space-y-1.5">
                  <Label htmlFor="ownerEmail" className="text-xs">Owner Email</Label>
                  <Input
                    id="ownerEmail"
                    type="email"
                    placeholder="owner@apex.org"
                    value={ownerEmail}
                    onChange={(e) => setOwnerEmail(e.target.value)}
                    required
                    className="text-xs"
                  />
                </div>
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="plan" className="text-xs">Subscription Tier</Label>
                <Select value={plan} onValueChange={setPlan}>
                  <SelectTrigger id="plan" className="text-xs">
                    <SelectValue placeholder="Select plan" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="smart">Smart (Starter)</SelectItem>
                    <SelectItem value="optimize">Optimize</SelectItem>
                    <SelectItem value="pro">Pro</SelectItem>
                    <SelectItem value="enterprise">Enterprise</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="rounded-md border border-primary/20 bg-primary/5 p-3 text-xs text-muted-foreground">
                <p className="font-medium text-foreground">Secure Verification Code Setup</p>
                <p className="mt-0.5 text-[11px] leading-relaxed">
                  A 6-character cryptographically secure verification code will be sent to the owner's email to verify their identity and set their private password.
                </p>
              </div>

              <DialogFooter className="pt-3">
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() => setOpenCreate(false)}
                  className="text-xs"
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  size="sm"
                  disabled={createMutation.isPending}
                  className="bg-primary text-primary-foreground text-xs"
                >
                  {createMutation.isPending ? "Provisioning..." : "Create Organization"}
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {/* Filter & Search Bar */}
      <Card className="card-enterprise p-4">
        <div className="flex flex-col sm:flex-row items-center gap-3">
          <div className="relative flex-1 w-full">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search by organization name, slug..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-9 text-xs h-9"
            />
          </div>

          <div className="flex items-center gap-2 w-full sm:w-auto">
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-[140px] text-xs h-9">
                <SelectValue placeholder="All Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="active">Active</SelectItem>
                <SelectItem value="pending_verification">Pending Review</SelectItem>
                <SelectItem value="suspended">Suspended</SelectItem>
              </SelectContent>
            </Select>

            <Select value={planFilter} onValueChange={setPlanFilter}>
              <SelectTrigger className="w-[140px] text-xs h-9">
                <SelectValue placeholder="All Plans" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Plans</SelectItem>
                <SelectItem value="smart">Smart</SelectItem>
                <SelectItem value="optimize">Optimize</SelectItem>
                <SelectItem value="pro">Pro</SelectItem>
                <SelectItem value="enterprise">Enterprise</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </Card>

      {/* Organizations Directory Table / List */}
      <Card className="card-enterprise overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="border-b border-border bg-muted/40 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              <tr>
                <th className="py-3 px-4">Organization & Slug</th>
                <th className="py-3 px-4">Subscription Plan</th>
                <th className="py-3 px-4">Status</th>
                <th className="py-3 px-4">Custom Domain</th>
                <th className="py-3 px-4">Provisioned Date</th>
                <th className="py-3 px-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border text-xs">
              {isLoading ? (
                <tr>
                  <td colSpan={6} className="py-8 text-center text-muted-foreground">
                    Loading organizations from cluster...
                  </td>
                </tr>
              ) : filteredOrgs.length === 0 ? (
                <tr>
                  <td colSpan={6} className="py-8 text-center text-muted-foreground">
                    No organizations match your query.
                  </td>
                </tr>
              ) : (
                filteredOrgs.map((org) => (
                  <tr key={org.id} className="hover:bg-muted/30 transition-colors">
                    <td className="py-3 px-4">
                      <div className="flex flex-col">
                        <Link
                          to={`/platform/organizations/${org.id}`}
                          className="font-semibold text-foreground hover:text-primary hover:underline"
                        >
                          {org.name}
                        </Link>
                        <span className="font-mono text-[11px] text-muted-foreground">
                          {org.slug}
                        </span>
                      </div>
                    </td>
                    <td className="py-3 px-4">
                      <Badge variant="secondary" className="capitalize text-xs font-medium">
                        {org.plan || "Smart"}
                      </Badge>
                    </td>
                    <td className="py-3 px-4">
                      {getStatusBadge(org.status)}
                    </td>
                    <td className="py-3 px-4 font-mono text-[11px] text-muted-foreground">
                      {org.customDomain || "—"}
                    </td>
                    <td className="py-3 px-4 text-muted-foreground">
                      {formatDate(org.createdAt)}
                    </td>
                    <td className="py-3 px-4 text-right">
                      <Button
                        variant="ghost"
                        size="sm"
                        asChild
                        className="h-8 text-xs hover:text-primary gap-1"
                      >
                        <Link to={`/platform/organizations/${org.id}`}>
                          Manage <ArrowUpRight className="h-3 w-3" />
                        </Link>
                      </Button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  );
}
