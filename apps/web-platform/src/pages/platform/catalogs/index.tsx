import React, { useState } from "react";
import {
  BookOpen,
  Search,
  Plus,
  Edit2,
  Stethoscope,
  Activity,
  Radio,
  Pill,
  CheckCircle2,
  FileText,
} from "lucide-react";
import {
  useCatalogItems,
  useCreateCatalogItem,
  useUpdateCatalogItem,
  useSearchICD10,
} from "@/api/hooks/use-catalogs";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Switch } from "@/components/ui/switch";
import type { CatalogDomain, CatalogItem } from "@/api/contracts";
import { toast } from "sonner";

export default function MasterCatalogsPage() {
  const [activeDomain, setActiveDomain] = useState<CatalogDomain>("lab");
  const [search, setSearch] = useState("");

  const { data: items, isLoading, refetch } = useCatalogItems(activeDomain);
  const createMutation = useCreateCatalogItem(activeDomain);

  // ICD-10 Query State
  const [icdQuery, setIcdQuery] = useState("");
  const { data: icdResults, isLoading: icdLoading } = useSearchICD10(icdQuery);

  // Create/Edit Item Modal State
  const [openCreate, setOpenCreate] = useState(false);
  const [editingItem, setEditingItem] = useState<CatalogItem | null>(null);
  const updateMutation = useUpdateCatalogItem(activeDomain, editingItem?.id || "");

  const [code, setCode] = useState("");
  const [name, setName] = useState("");
  const [category, setCategory] = useState("");
  const [description, setDescription] = useState("");
  const [systemGroup, setSystemGroup] = useState("");
  const [basePrice, setBasePrice] = useState<number | string>("");
  const [isActive, setIsActive] = useState(true);

  const filteredItems = (items || []).filter(
    (item) =>
      (item.name || "").toLowerCase().includes(search.toLowerCase()) ||
      (item.code || "").toLowerCase().includes(search.toLowerCase()) ||
      (item.category || "").toLowerCase().includes(search.toLowerCase())
  );

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createMutation.mutateAsync({
        domain: activeDomain,
        code,
        name,
        category,
        description,
        systemGroup,
        basePrice: basePrice !== "" ? Number(basePrice) : 0,
        isActive,
        version: 1,
      });
      toast.success(`Catalog item '${name}' created!`);
      setOpenCreate(false);
      setCode("");
      setName("");
      setCategory("");
      setDescription("");
      setSystemGroup("");
      setBasePrice("");
      refetch();
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Failed to create catalog item");
    }
  };

  const handleOpenEdit = (item: CatalogItem) => {
    setEditingItem(item);
    setCode(item.code);
    setName(item.name);
    setCategory(item.category);
    setDescription(item.description || "");
    setSystemGroup(item.systemGroup || "");
    setBasePrice(item.basePrice !== undefined ? item.basePrice : "");
    setIsActive(item.isActive);
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingItem) return;
    try {
      await updateMutation.mutateAsync({
        ...editingItem,
        code,
        name,
        category,
        description,
        systemGroup,
        basePrice: basePrice !== "" ? Number(basePrice) : 0,
        isActive,
      });
      toast.success(`Catalog item '${name}' updated!`);
      setEditingItem(null);
      refetch();
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Failed to update catalog item");
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              Master Reference Catalogs
            </h1>
            <Badge variant="outline" className="border-primary/30 text-primary bg-primary/5 text-xs font-mono">
              EPIC-010
            </Badge>
          </div>
          <p className="text-sm text-muted-foreground mt-1">
            Global standardized clinical test dictionaries, laboratory panels, radiology modalities, and drug formularies.
          </p>
        </div>

        <Dialog open={openCreate} onOpenChange={setOpenCreate}>
          <DialogTrigger asChild>
            <Button className="bg-primary hover:bg-primary/90 text-primary-foreground text-xs h-9 gap-2 shadow-sm">
              <Plus className="h-4 w-4" />
              Add Master Item
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle className="text-base font-semibold">
                Add {activeDomain.toUpperCase()} Master Item
              </DialogTitle>
              <DialogDescription className="text-xs">
                Create a global reference catalog entry accessible to all health networks.
              </DialogDescription>
            </DialogHeader>

            <form onSubmit={handleCreate} className="space-y-3 py-2 text-xs">
              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-1">
                  <Label htmlFor="itemCode">Code Identifier</Label>
                  <Input
                    id="itemCode"
                    placeholder="e.g. FBC_001"
                    value={code}
                    onChange={(e) => setCode(e.target.value)}
                    required
                    className="text-xs font-mono"
                  />
                </div>
                <div className="space-y-1">
                  <Label htmlFor="itemCategory">Category</Label>
                  <Input
                    id="itemCategory"
                    placeholder="e.g. Hematology"
                    value={category}
                    onChange={(e) => setCategory(e.target.value)}
                    required
                    className="text-xs"
                  />
                </div>
              </div>

              <div className="space-y-1">
                <Label htmlFor="itemName">Item Name</Label>
                <Input
                  id="itemName"
                  placeholder="e.g. Full Blood Count"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  required
                  className="text-xs"
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-1">
                  <Label htmlFor="itemGroup">System / Specimen Group</Label>
                  <Input
                    id="itemGroup"
                    placeholder="e.g. Whole Blood"
                    value={systemGroup}
                    onChange={(e) => setSystemGroup(e.target.value)}
                    className="text-xs"
                  />
                </div>
                <div className="space-y-1">
                  <Label htmlFor="itemBasePrice">Base Benchmark Price (₦)</Label>
                  <Input
                    id="itemBasePrice"
                    type="number"
                    step="0.01"
                    placeholder="0.00"
                    value={basePrice}
                    onChange={(e) => setBasePrice(e.target.value)}
                    className="text-xs font-mono"
                  />
                </div>
              </div>

              <div className="space-y-1">
                <Label htmlFor="itemDesc">Description</Label>
                <Input
                  id="itemDesc"
                  placeholder="Standard diagnostic profile"
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  className="text-xs"
                />
              </div>

              <div className="flex items-center justify-between pt-2">
                <Label htmlFor="itemActive" className="cursor-pointer">Active State</Label>
                <Switch
                  id="itemActive"
                  checked={isActive}
                  onCheckedChange={setIsActive}
                />
              </div>

              <DialogFooter className="pt-3">
                <Button type="button" variant="outline" size="sm" onClick={() => setOpenCreate(false)} className="text-xs">
                  Cancel
                </Button>
                <Button type="submit" size="sm" className="bg-primary text-primary-foreground text-xs">
                  Create Item
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {/* Main Tabs for Domains */}
      <Tabs
        value={activeDomain}
        onValueChange={(val: any) => {
          setActiveDomain(val);
          setSearch("");
        }}
        className="space-y-4"
      >
        <TabsList className="bg-muted/60 p-1 border border-border">
          <TabsTrigger value="lab" className="text-xs gap-1.5">
            <Activity className="h-3.5 w-3.5 text-primary" /> Laboratory (LIS)
          </TabsTrigger>
          <TabsTrigger value="clinical" className="text-xs gap-1.5">
            <Stethoscope className="h-3.5 w-3.5 text-blue-500" /> Clinical Tests
          </TabsTrigger>
          <TabsTrigger value="radiology" className="text-xs gap-1.5">
            <Radio className="h-3.5 w-3.5 text-purple-500" /> Radiology & Imaging
          </TabsTrigger>
          <TabsTrigger value="pharmacy" className="text-xs gap-1.5">
            <Pill className="h-3.5 w-3.5 text-amber-500" /> Pharmacy & Formularies
          </TabsTrigger>
        </TabsList>

        {/* Filter Bar */}
        <Card className="card-enterprise p-4">
          <div className="relative">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder={`Search ${activeDomain} catalog by code, name, category...`}
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-9 text-xs h-9"
            />
          </div>
        </Card>

        {/* Items Table */}
        <Card className="card-enterprise overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs">
              <thead className="border-b border-border bg-muted/40 font-semibold uppercase tracking-wider text-muted-foreground">
                <tr>
                  <th className="py-3 px-4">Code Identifier</th>
                  <th className="py-3 px-4">Item Name</th>
                  <th className="py-3 px-4">Category</th>
                  <th className="py-3 px-4">Group / Specimen</th>
                  <th className="py-3 px-4">Base Price</th>
                  <th className="py-3 px-4">Status</th>
                  <th className="py-3 px-4 text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {isLoading ? (
                  <tr>
                    <td colSpan={7} className="py-8 text-center text-muted-foreground">
                      Loading catalog items...
                    </td>
                  </tr>
                ) : filteredItems.length === 0 ? (
                  <tr>
                    <td colSpan={7} className="py-8 text-center text-muted-foreground">
                      No master catalog items found for this domain.
                    </td>
                  </tr>
                ) : (
                  filteredItems.map((item) => (
                    <tr key={item.id || item.code} className="hover:bg-muted/20">
                      <td className="py-3 px-4 font-mono font-semibold text-primary">
                        {item.code}
                      </td>
                      <td className="py-3 px-4 font-medium text-foreground">
                        {item.name}
                      </td>
                      <td className="py-3 px-4 text-muted-foreground">
                        {item.category}
                      </td>
                      <td className="py-3 px-4 font-mono text-[11px] text-muted-foreground">
                        {item.systemGroup || "—"}
                      </td>
                      <td className="py-3 px-4 font-mono text-xs text-foreground font-semibold">
                        {item.basePrice && item.basePrice > 0
                          ? `₦${item.basePrice.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
                          : "—"}
                      </td>
                      <td className="py-3 px-4">
                        {item.isActive ? (
                          <Badge variant="outline" className="border-emerald-500/30 text-emerald-600 bg-emerald-500/5 text-[10px]">
                            Active
                          </Badge>
                        ) : (
                          <Badge variant="outline" className="border-border text-muted-foreground text-[10px]">
                            Inactive
                          </Badge>
                        )}
                      </td>
                      <td className="py-3 px-4 text-right">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleOpenEdit(item)}
                          className="h-7 text-xs"
                        >
                          <Edit2 className="h-3 w-3" />
                        </Button>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </Card>
      </Tabs>

      {/* ICD-10 Search Explorer Card */}
      <Card className="card-enterprise">
        <CardHeader>
          <CardTitle className="text-sm font-semibold flex items-center gap-2">
            <BookOpen className="h-4 w-4 text-primary" />
            Global ICD-10 Diagnostic Master Reference
          </CardTitle>
          <CardDescription className="text-xs">
            Query standardized WHO ICD-10 medical diagnostic terminology.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="relative">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search ICD-10 by code or diagnostic description (e.g. malaria, hypertension, E11)..."
              value={icdQuery}
              onChange={(e) => setIcdQuery(e.target.value)}
              className="pl-9 text-xs h-9"
            />
          </div>

          {icdLoading && (
            <p className="text-xs text-muted-foreground">Searching ICD-10 codes...</p>
          )}

          {icdResults && icdResults.length > 0 && (
            <div className="rounded-lg border border-border overflow-hidden">
              <table className="w-full text-left text-xs">
                <thead className="border-b border-border bg-muted/40 font-semibold uppercase text-muted-foreground">
                  <tr>
                    <th className="py-2 px-3">ICD-10 Code</th>
                    <th className="py-2 px-3">Diagnostic Description</th>
                    <th className="py-2 px-3">Category</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {icdResults.slice(0, 8).map((code) => (
                    <tr key={code.code} className="hover:bg-muted/20">
                      <td className="py-2 px-3 font-mono font-semibold text-primary">{code.code}</td>
                      <td className="py-2 px-3 text-foreground">{code.description}</td>
                      <td className="py-2 px-3 text-muted-foreground">{code.category || "General"}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Edit Item Modal */}
      <Dialog open={!!editingItem} onOpenChange={(open) => !open && setEditingItem(null)}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="text-base font-semibold">Edit Master Catalog Item</DialogTitle>
            <DialogDescription className="text-xs">
              Update global reference dictionary parameters.
            </DialogDescription>
          </DialogHeader>

          <form onSubmit={handleUpdate} className="space-y-3 py-2 text-xs">
            <div className="grid grid-cols-2 gap-3">
              <div className="space-y-1">
                <Label htmlFor="editItemCode">Code Identifier</Label>
                <Input
                  id="editItemCode"
                  value={code}
                  onChange={(e) => setCode(e.target.value)}
                  required
                  className="text-xs font-mono"
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="editItemCategory">Category</Label>
                <Input
                  id="editItemCategory"
                  value={category}
                  onChange={(e) => setCategory(e.target.value)}
                  required
                  className="text-xs"
                />
              </div>
            </div>

            <div className="space-y-1">
              <Label htmlFor="editItemName">Item Name</Label>
              <Input
                id="editItemName"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
                className="text-xs"
              />
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div className="space-y-1">
                <Label htmlFor="editItemGroup">System / Specimen Group</Label>
                <Input
                  id="editItemGroup"
                  value={systemGroup}
                  onChange={(e) => setSystemGroup(e.target.value)}
                  className="text-xs"
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="editItemBasePrice">Base Benchmark Price (₦)</Label>
                <Input
                  id="editItemBasePrice"
                  type="number"
                  step="0.01"
                  placeholder="0.00"
                  value={basePrice}
                  onChange={(e) => setBasePrice(e.target.value)}
                  className="text-xs font-mono"
                />
              </div>
            </div>

            <div className="space-y-1">
              <Label htmlFor="editItemDesc">Description</Label>
              <Input
                id="editItemDesc"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                className="text-xs"
              />
            </div>

            <div className="flex items-center justify-between pt-2">
              <Label htmlFor="editItemActive" className="cursor-pointer">Active State</Label>
              <Switch
                id="editItemActive"
                checked={isActive}
                onCheckedChange={setIsActive}
              />
            </div>

            <DialogFooter className="pt-3">
              <Button type="button" variant="outline" size="sm" onClick={() => setEditingItem(null)} className="text-xs">
                Cancel
              </Button>
              <Button type="submit" size="sm" className="bg-primary text-primary-foreground text-xs">
                Save Changes
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
