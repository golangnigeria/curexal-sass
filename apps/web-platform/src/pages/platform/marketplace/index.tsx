import React, { useState } from "react";
import {
  Store,
  CheckCircle2,
  Plus,
  Layers,
  Sparkles,
  Shield,
  Clock,
  Search,
  ArrowUpRight,
  Filter,
} from "lucide-react";
import {
  useCapabilityCatalog,
  useGrantCapability,
  useStartTrialCapability,
} from "@/api/hooks/use-marketplace";
import { usePlatformOrganizations } from "@/api/hooks/use-organizations";
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
import { formatCurrency } from "@/lib/utils";
import { toast } from "sonner";

export default function MarketplacePage() {
  const { data: catalog, isLoading } = useCapabilityCatalog();
  const { data: orgs } = usePlatformOrganizations();

  const [search, setSearch] = useState("");
  const [categoryFilter, setCategoryFilter] = useState("all");

  // Grant Modal State
  const [openGrant, setOpenGrant] = useState(false);
  const [selectedOrgId, setSelectedOrgId] = useState("");
  const [selectedCapCode, setSelectedCapCode] = useState("");
  const [grantMode, setGrantMode] = useState<"permanent" | "trial">("permanent");
  const [trialDays, setTrialDays] = useState(14);

  const grantMutation = useGrantCapability(selectedOrgId);
  const trialMutation = useStartTrialCapability(selectedOrgId);

  const filteredCatalog = (catalog || []).filter((item) => {
    const name = item.name || item.capability?.name || "";
    const code = item.code || item.capability?.code || "";
    const desc = item.description || item.capability?.description || "";
    const cat = item.category || item.capability?.category || "";

    const matchesSearch =
      name.toLowerCase().includes(search.toLowerCase()) ||
      code.toLowerCase().includes(search.toLowerCase()) ||
      desc.toLowerCase().includes(search.toLowerCase());
    const matchesCat = categoryFilter === "all" || cat.toLowerCase() === categoryFilter.toLowerCase();
    return matchesSearch && matchesCat;
  });

  const handleGrantSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedOrgId || !selectedCapCode) {
      toast.error("Please select both an organization and a capability.");
      return;
    }

    try {
      if (grantMode === "trial") {
        await trialMutation.mutateAsync({
          capabilityCode: selectedCapCode,
          durationDays: trialDays,
        });
        toast.success(`Activated ${trialDays}-day trial for '${selectedCapCode}'!`);
      } else {
        await grantMutation.mutateAsync({
          capabilityCode: selectedCapCode,
          source: "addon",
        });
        toast.success(`Granted capability '${selectedCapCode}' to organization!`);
      }
      setOpenGrant(false);
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Operation failed.");
    }
  };

  const openGrantForCap = (capCode: string) => {
    setSelectedCapCode(capCode);
    setOpenGrant(true);
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              B2B Capability Marketplace
            </h1>
            <Badge variant="outline" className="border-primary/30 text-primary bg-primary/5 text-xs">
              Commercial Add-Ons
            </Badge>
          </div>
          <p className="text-sm text-muted-foreground mt-1">
            Global catalog of extensible diagnostic, clinical, and enterprise billing modules available for subscription.
          </p>
        </div>

        <Dialog open={openGrant} onOpenChange={setOpenGrant}>
          <DialogTrigger asChild>
            <Button className="bg-primary hover:bg-primary/90 text-primary-foreground text-xs h-9 gap-2 shadow-sm">
              <Plus className="h-4 w-4" />
              Grant Add-On to Org
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle className="text-base">Grant Capability Add-On</DialogTitle>
              <DialogDescription className="text-xs">
                Activate a modular capability for a registered healthcare network.
              </DialogDescription>
            </DialogHeader>

            <form onSubmit={handleGrantSubmit} className="space-y-3 py-2 text-xs">
              <div className="space-y-1">
                <Label htmlFor="targetOrg">Select Organization</Label>
                <Select value={selectedOrgId} onValueChange={setSelectedOrgId}>
                  <SelectTrigger id="targetOrg" className="text-xs">
                    <SelectValue placeholder="Choose organization" />
                  </SelectTrigger>
                  <SelectContent>
                    {(orgs || []).map((o) => (
                      <SelectItem key={o.id} value={o.id}>
                        {o.name} ({o.slug})
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-1">
                <Label htmlFor="targetCap">Select Capability</Label>
                <Select value={selectedCapCode} onValueChange={setSelectedCapCode}>
                  <SelectTrigger id="targetCap" className="text-xs">
                    <SelectValue placeholder="Choose capability code" />
                  </SelectTrigger>
                  <SelectContent>
                    {(catalog || []).map((c) => {
                      const capCode = c.code || c.capability?.code || "";
                      const capName = c.name || c.capability?.name || capCode;
                      return (
                        <SelectItem key={capCode} value={capCode}>
                          {capName} ({capCode})
                        </SelectItem>
                      );
                    })}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-1">
                <Label htmlFor="modeSelect">Grant Type</Label>
                <Select
                  value={grantMode}
                  onValueChange={(val: any) => setGrantMode(val)}
                >
                  <SelectTrigger id="modeSelect" className="text-xs">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="permanent">Standard Add-On Entitlement</SelectItem>
                    <SelectItem value="trial">Time-Limited Promotional Trial</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {grantMode === "trial" && (
                <div className="space-y-1">
                  <Label htmlFor="trialDays">Trial Duration (Days)</Label>
                  <Input
                    id="trialDays"
                    type="number"
                    min={1}
                    max={90}
                    value={trialDays}
                    onChange={(e) => setTrialDays(parseInt(e.target.value, 10))}
                    className="text-xs"
                  />
                </div>
              )}

              <DialogFooter className="pt-3">
                <Button type="button" variant="outline" size="sm" onClick={() => setOpenGrant(false)} className="text-xs">
                  Cancel
                </Button>
                <Button type="submit" size="sm" className="bg-primary text-primary-foreground text-xs">
                  Confirm Activation
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
              placeholder="Search capabilities by name, code, description..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-9 text-xs h-9"
            />
          </div>

          <Select value={categoryFilter} onValueChange={setCategoryFilter}>
            <SelectTrigger className="w-[160px] text-xs h-9">
              <SelectValue placeholder="All Categories" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Categories</SelectItem>
              <SelectItem value="clinical">Clinical</SelectItem>
              <SelectItem value="laboratory">Laboratory</SelectItem>
              <SelectItem value="radiology">Radiology</SelectItem>
              <SelectItem value="pharmacy">Pharmacy</SelectItem>
              <SelectItem value="billing">Billing & POS</SelectItem>
              <SelectItem value="platform">Platform Extensions</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </Card>

      {/* Capabilities Grid */}
      {isLoading ? (
        <div className="py-16 text-center text-sm text-muted-foreground">
          Loading marketplace catalog...
        </div>
      ) : filteredCatalog.length === 0 ? (
        <div className="py-16 text-center text-sm text-muted-foreground">
          No capabilities found in catalog.
        </div>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {filteredCatalog.map((item) => {
            const capCode = item.code || item.capability?.code || "";
            const capName = item.name || item.capability?.name || capCode;
            const capCategory = item.category || item.capability?.category || "general";
            const capDesc = item.description || item.capability?.description || "Modular platform extension";
            const price = item.monthlyPrice ?? item.prices?.[0]?.monthlyPrice ?? 0;
            const currency = item.currency || item.prices?.[0]?.currency || "NGN";

            return (
              <Card key={capCode || capName} className="card-enterprise hover-lift flex flex-col justify-between">
                <CardHeader className="pb-3">
                  <div className="flex items-start justify-between gap-2">
                    <div className="flex flex-col">
                      <CardTitle className="text-sm font-semibold text-foreground">
                        {capName}
                      </CardTitle>
                      <span className="font-mono text-[11px] text-primary">
                        {capCode}
                      </span>
                    </div>
                    <Badge variant="outline" className="text-[10px] capitalize font-medium">
                      {capCategory}
                    </Badge>
                  </div>
                  <CardDescription className="text-xs line-clamp-2 mt-2">
                    {capDesc}
                  </CardDescription>
                </CardHeader>

                <CardContent className="space-y-3 pt-0">
                  <div className="rounded-lg border border-border bg-secondary/30 p-2.5 flex items-center justify-between text-xs">
                    <span className="text-muted-foreground">Pricing:</span>
                    <span className="font-semibold text-foreground">
                      {price > 0 ? (
                        `${formatCurrency(price, currency)}/mo`
                      ) : (
                        <span className="text-emerald-600 font-medium">Core Plan Included</span>
                      )}
                    </span>
                  </div>

                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => openGrantForCap(capCode)}
                    className="w-full text-xs h-8 gap-1.5 hover:bg-primary hover:text-primary-foreground transition-colors"
                  >
                    <Plus className="h-3 w-3" /> Grant to Organization
                  </Button>
                </CardContent>
              </Card>
            );
          })}
        </div>
      )}
    </div>
  );
}
