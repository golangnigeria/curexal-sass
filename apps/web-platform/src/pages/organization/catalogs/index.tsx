import React, { useState } from "react";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { useOrgCatalogs, useUpdateCatalogPrice } from "@/api/hooks/use-organization";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import {
  BookOpen,
  Search,
  Check,
  Edit2,
  DollarSign,
  Activity,
  Stethoscope,
  Microscope,
  Pill,
} from "lucide-react";

// Default diagnostic and clinical catalog items if API query pending
const defaultCatalogItems = [
  { id: "cat-1", code: "LAB-CBC", name: "Complete Blood Count (Automated 5-Part Diff)", category: "Hematology", moduleCode: "laboratory", standardPrice: 7500, customPrice: 7000, currency: "NGN", taxRate: 0, isActive: true },
  { id: "cat-2", code: "LAB-FBS", name: "Fasting Blood Sugar (Glucose Hexokinase)", category: "Clinical Biochemistry", moduleCode: "laboratory", standardPrice: 3500, customPrice: 3500, currency: "NGN", taxRate: 0, isActive: true },
  { id: "cat-3", code: "LAB-LIPID", name: "Lipid Profile Panel (Cholesterol, HDL, LDL, Trig)", category: "Clinical Biochemistry", moduleCode: "laboratory", standardPrice: 12000, customPrice: 10500, currency: "NGN", taxRate: 0, isActive: true },
  { id: "cat-4", code: "LAB-LFT", name: "Liver Function Test (ALT, AST, ALP, Bilirubin)", category: "Clinical Biochemistry", moduleCode: "laboratory", standardPrice: 15000, customPrice: 14000, currency: "NGN", taxRate: 0, isActive: true },
  { id: "cat-5", code: "RAD-XR-CHEST", name: "Chest X-Ray (PA & Lateral View)", category: "Radiology", moduleCode: "radiology", standardPrice: 18000, customPrice: 16500, currency: "NGN", taxRate: 0, isActive: true },
  { id: "cat-6", code: "RAD-US-ABD", name: "Abdominal Ultrasound (Full Scan)", category: "Radiology", moduleCode: "radiology", standardPrice: 20000, customPrice: 18000, currency: "NGN", taxRate: 0, isActive: true },
  { id: "cat-7", code: "CLN-CONS-GP", name: "General Practice Consultation", category: "Outpatient", moduleCode: "clinical", standardPrice: 10000, customPrice: 10000, currency: "NGN", taxRate: 0, isActive: true },
  { id: "cat-8", code: "CLN-CONS-SPEC", name: "Specialist Consultant Review", category: "Outpatient", moduleCode: "clinical", standardPrice: 25000, customPrice: 25000, currency: "NGN", taxRate: 0, isActive: true },
];

export default function OrganizationCatalogsPage() {
  const { data: bootstrap } = useBootstrap();
  const { data: serverCatalogs, isLoading } = useOrgCatalogs();
  const updatePriceMutation = useUpdateCatalogPrice();

  const currency = bootstrap?.workspace?.currency || "NGN";
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedCategory, setSelectedCategory] = useState("ALL");
  const [editingItemId, setEditingItemId] = useState<string | null>(null);
  const [editPriceValue, setEditPriceValue] = useState<number>(0);

  const catalogItems = (serverCatalogs && serverCatalogs.length > 0) ? serverCatalogs : defaultCatalogItems;

  const categories = ["ALL", ...Array.from(new Set(catalogItems.map((c) => c.category)))];

  const handleStartEdit = (item: any) => {
    setEditingItemId(item.id);
    setEditPriceValue(item.customPrice || item.standardPrice);
  };

  const handleSavePrice = async (itemId: string) => {
    try {
      await updatePriceMutation.mutateAsync({ itemId, customPrice: Number(editPriceValue) });
      toast.success("Catalog Tariff Updated!", {
        description: "Custom price override applied to patient billing POS.",
      });
      setEditingItemId(null);
    } catch (err: any) {
      toast.error("Failed to update price: " + (err.message || "Network error"));
    }
  };

  const filteredItems = catalogItems.filter((item: any) => {
    const itemName = (item.name || "").toLowerCase();
    const itemCode = (item.code || "").toLowerCase();
    const itemCategory = (item.category || "").toLowerCase();
    const query = searchQuery.toLowerCase().trim();

    const matchesSearch = !query ||
      itemName.includes(query) ||
      itemCode.includes(query) ||
      itemCategory.includes(query);
    const matchesCategory = selectedCategory === "ALL" || item.category === selectedCategory;
    return matchesSearch && matchesCategory;
  });

  const formatPrice = (val: number) => {
    return new Intl.NumberFormat("en-NG", {
      style: "currency",
      currency: currency,
      maximumFractionDigits: 0,
    }).format(val);
  };

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-6">
        <div>
          <div className="flex items-center gap-2.5 mb-1">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              Corporate Catalogs & Custom Pricing
            </h1>
            <Badge variant="outline" className="border-primary/40 text-primary bg-primary/10 text-[10px] font-mono">
              {catalogItems.length} Master Services
            </Badge>
          </div>
          <p className="text-xs text-muted-foreground">
            Set custom corporate tariffs, test investigation fees, and HMO insurance copayments across your branch network.
          </p>
        </div>
      </div>

      {/* Filter Tabs & Search */}
      <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
        {/* Category Pills */}
        <div className="flex items-center gap-1.5 overflow-x-auto w-full sm:w-auto pb-1">
          {categories.map((cat) => (
            <button
              key={cat}
              type="button"
              onClick={() => setSelectedCategory(cat)}
              className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-all ${
                selectedCategory === cat
                  ? "bg-primary text-primary-foreground shadow-sm font-semibold"
                  : "bg-secondary/40 text-muted-foreground hover:text-foreground hover:bg-secondary"
              }`}
            >
              {cat}
            </button>
          ))}
        </div>

        {/* Search */}
        <div className="w-full sm:max-w-xs relative">
          <Search className="w-3.5 h-3.5 absolute left-3 top-2.5 text-muted-foreground" />
          <Input
            placeholder="Search catalog code, test name..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="text-xs h-8 pl-9 bg-secondary/30"
          />
        </div>
      </div>

      {/* Catalog Table */}
      <Card className="border-border shadow-sm overflow-hidden">
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full text-xs text-left">
              <thead className="bg-secondary/40 border-b border-border text-[11px] font-bold text-muted-foreground uppercase tracking-wider">
                <tr>
                  <th className="py-3 px-4">Service Code & Investigation</th>
                  <th className="py-3 px-4">Category</th>
                  <th className="py-3 px-4">Standard Base Price</th>
                  <th className="py-3 px-4">Your Custom Corporate Tariff</th>
                  <th className="py-3 px-4 text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {filteredItems.map((item) => {
                  const isEditing = editingItemId === item.id;
                  const isDiscounted = item.customPrice < item.standardPrice;

                  return (
                    <tr key={item.id} className="hover:bg-secondary/20 transition-colors">
                      <td className="py-3 px-4">
                        <div>
                          <p className="font-semibold text-foreground">{item.name}</p>
                          <p className="text-[11px] font-mono text-muted-foreground">{item.code}</p>
                        </div>
                      </td>
                      <td className="py-3 px-4">
                        <Badge variant="outline" className="text-[10px] border-border font-medium">
                          {item.category}
                        </Badge>
                      </td>
                      <td className="py-3 px-4 font-mono text-muted-foreground">
                        {formatPrice(item.standardPrice)}
                      </td>
                      <td className="py-3 px-4">
                        {isEditing ? (
                          <div className="flex items-center gap-2 max-w-[140px]">
                            <Input
                              type="number"
                              value={editPriceValue}
                              onChange={(e) => setEditPriceValue(Number(e.target.value))}
                              className="text-xs h-8 font-mono"
                              autoFocus
                            />
                            <Button
                              size="icon"
                              className="h-8 w-8 bg-primary text-primary-foreground shrink-0"
                              onClick={() => handleSavePrice(item.id)}
                            >
                              <Check className="w-3.5 h-3.5" />
                            </Button>
                          </div>
                        ) : (
                          <div className="flex items-center gap-2">
                            <span className="font-mono font-bold text-foreground">
                              {formatPrice(item.customPrice || item.standardPrice)}
                            </span>
                            {isDiscounted && (
                              <Badge className="text-[9px] bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/30">
                                Override
                              </Badge>
                            )}
                          </div>
                        )}
                      </td>
                      <td className="py-3 px-4 text-right">
                        {!isEditing && (
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleStartEdit(item)}
                            className="text-xs h-8 gap-1 text-muted-foreground hover:text-foreground"
                          >
                            <Edit2 className="w-3.5 h-3.5" />
                            Edit Tariff
                          </Button>
                        )}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
