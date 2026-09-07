import React, { useState, useEffect } from "react";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { useCapabilityCatalog } from "@/api/hooks/use-marketplace";
import { usePricingRules } from "@/api/hooks/use-pricing";
import { billingService, type InvoicePayload } from "@/api/services/billing.service";
import { CommercialCheckoutDialog, type CheckoutItem } from "@/components/billing/checkout-dialog";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { formatCurrency, formatDate } from "@/lib/utils";
import { toast } from "sonner";
import {
  CreditCard,
  CheckCircle2,
  Sparkles,
  Zap,
  Shield,
  Activity,
  Layers,
  ArrowRight,
  Download,
  Receipt,
  Clock,
  ExternalLink,
} from "lucide-react";

export default function OrganizationBillingPage() {
  const { data: bootstrap, refetch: refetchBootstrap } = useBootstrap();
  const { data: catalog, refetch: refetchCatalog } = useCapabilityCatalog();
  const { data: pricingRules } = usePricingRules();

  const orgId = bootstrap?.organization?.id || "";
  const orgPlan = bootstrap?.organization?.subscription || "smart";
  const currency = bootstrap?.workspace?.currency || "NGN";
  const activeCapabilities = bootstrap?.capabilities || [];
  const limits = bootstrap?.limits || { maxBranches: 1, maxMembers: 5, storageGb: 10 };

  const [invoices, setInvoices] = useState<InvoicePayload[]>([]);
  const [isLoadingInvoices, setIsLoadingInvoices] = useState(false);

  // Checkout modal state
  const [checkoutItem, setCheckoutItem] = useState<CheckoutItem | null>(null);
  const [isCheckoutOpen, setIsCheckoutOpen] = useState(false);

  useEffect(() => {
    if (orgId) {
      setIsLoadingInvoices(true);
      billingService
        .getInvoices(orgId)
        .then((data) => setInvoices(data))
        .finally(() => setIsLoadingInvoices(false));
    }
  }, [orgId]);

  const handleOpenPlanCheckout = (planCode: string, planName: string, price: number) => {
    setCheckoutItem({
      code: planCode,
      name: planName,
      price,
      currency,
      type: "plan",
      billingCycle: "monthly",
    });
    setIsCheckoutOpen(true);
  };

  const handleOpenAddOnCheckout = (capCode: string, capName: string, price: number) => {
    setCheckoutItem({
      code: capCode,
      name: capName,
      price,
      currency,
      type: "addon",
      billingCycle: "monthly",
    });
    setIsCheckoutOpen(true);
  };

  const handleCheckoutSuccess = () => {
    if (refetchBootstrap) refetchBootstrap();
    if (refetchCatalog) refetchCatalog();
    if (orgId) {
      billingService.getInvoices(orgId).then((data) => setInvoices(data));
    }
  };

  // Plan tiers configuration with dynamic pricing rule overlay
  const planTiers = [
    {
      code: "smart",
      name: "Smart Starter",
      priceMonthly: 0,
      desc: "Essential reception, patient registration & core clinics.",
      features: ["1 Branch Facility", "5 Staff Member Seats", "10 GB Cloud Storage", "Basic Lab & Clinic"],
    },
    {
      code: "optimize",
      name: "Optimize Tier",
      priceMonthly: 35000,
      desc: "Growing diagnostic centers & multi-specialty clinics.",
      features: ["3 Branch Facilities", "25 Staff Member Seats", "50 GB Cloud Storage", "Analyzer Interfacing"],
    },
    {
      code: "pro",
      name: "Pro Tier",
      priceMonthly: 95000,
      desc: "Full hospital & diagnostic laboratory networks.",
      features: ["10 Branch Facilities", "100 Staff Member Seats", "200 GB Cloud Storage", "DICOM PACS & LIS"],
    },
    {
      code: "enterprise",
      name: "Enterprise Custom",
      priceMonthly: 250000,
      desc: "Tertiary hospital groups & regional healthcare systems.",
      features: ["Unlimited Branches", "Unlimited Staff Seats", "5 TB Dedicated Storage", "White-label Custom Domain"],
    },
  ].map((p) => {
    const rule = (pricingRules || []).find((r) => r.targetCode === p.code);
    return {
      ...p,
      priceMonthly: rule?.monthlyPrice ?? p.priceMonthly,
    };
  });

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-6">
        <div>
          <div className="flex items-center gap-2.5 mb-1">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              Corporate Subscription & Billing
            </h1>
            <Badge variant="outline" className="border-primary/40 text-primary bg-primary/10 uppercase tracking-wider text-[10px] font-mono font-bold">
              Active: {orgPlan} Plan
            </Badge>
          </div>
          <p className="text-xs text-muted-foreground">
            Manage your organization base subscription, capacity quotas, and specialized diagnostic add-on packages.
          </p>
        </div>
      </div>

      {/* Plan Tiers Grid */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        {planTiers.map((plan) => {
          const isCurrent = orgPlan === plan.code;
          return (
            <Card
              key={plan.code}
              className={`border transition-all flex flex-col justify-between ${
                isCurrent
                  ? "border-primary shadow-md bg-primary/5"
                  : "border-border shadow-sm bg-card hover:border-border/80"
              }`}
            >
              <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                  <span className="text-xs font-bold text-foreground uppercase tracking-wider">{plan.name}</span>
                  {isCurrent && (
                    <Badge className="text-[9px] bg-primary text-primary-foreground font-mono">Current Plan</Badge>
                  )}
                </div>
                <div className="text-xl font-bold text-foreground mt-2 font-mono">
                  {plan.priceMonthly === 0 ? "Free / ₦0" : `${formatCurrency(plan.priceMonthly, currency)} / mo`}
                </div>
                <CardDescription className="text-[11px] min-h-[30px]">{plan.desc}</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4 pt-0">
                <div className="space-y-1.5 text-xs">
                  {plan.features.map((f, i) => (
                    <p key={i} className="flex items-center gap-1.5 text-muted-foreground text-[11px]">
                      <CheckCircle2 className="w-3.5 h-3.5 text-primary shrink-0" />
                      {f}
                    </p>
                  ))}
                </div>

                <Button
                  size="sm"
                  variant={isCurrent ? "outline" : "default"}
                  disabled={isCurrent}
                  onClick={() => handleOpenPlanCheckout(plan.code, plan.name, plan.priceMonthly)}
                  className={`w-full text-xs h-8 ${!isCurrent ? "bg-primary text-primary-foreground shadow" : ""}`}
                >
                  {isCurrent ? "Active Plan" : "Switch Plan"}
                </Button>
              </CardContent>
            </Card>
          );
        })}
      </div>

      {/* In-App Capability Add-On Marketplace */}
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-base font-bold text-foreground flex items-center gap-2">
              <Sparkles className="w-4 h-4 text-primary" />
              Specialized Diagnostic & Clinical Add-Ons
            </h3>
            <p className="text-xs text-muted-foreground">
              Expand your facility capabilities on-demand. Activated add-ons instantly unlock in your workspace navigation.
            </p>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {(catalog || []).map((addon) => {
            const addonCode = addon.code || addon.id || "";
            const isOwned = activeCapabilities.includes(addonCode);
            const price = addon.basePrice || addon.monthlyPrice || 25000;

            return (
              <Card key={addonCode} className="border-border shadow-sm hover:shadow-md transition-all flex flex-col justify-between bg-card">
                <CardHeader className="pb-3">
                  <div className="flex items-center justify-between">
                    <Badge variant="outline" className="text-[9px] font-mono uppercase border-border">
                      {addon.category || "Add-On"}
                    </Badge>
                    {isOwned ? (
                      <Badge className="text-[9px] bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/30">
                        Active License
                      </Badge>
                    ) : (
                      <span className="text-xs font-mono font-bold text-foreground">
                        {formatCurrency(price, currency)} / mo
                      </span>
                    )}
                  </div>
                  <CardTitle className="text-sm font-bold text-foreground mt-2">{addon.name || addonCode}</CardTitle>
                  <CardDescription className="text-xs min-h-[36px] text-muted-foreground">
                    {addon.description}
                  </CardDescription>
                </CardHeader>
                <CardContent className="pt-0">
                  <Button
                    size="sm"
                    variant={isOwned ? "outline" : "default"}
                    disabled={isOwned}
                    onClick={() => handleOpenAddOnCheckout(addonCode, addon.name || addonCode, price)}
                    className={`w-full text-xs h-8 gap-1.5 ${
                      !isOwned ? "bg-primary text-primary-foreground shadow" : "text-muted-foreground"
                    }`}
                  >
                    {isOwned ? (
                      <>
                        <CheckCircle2 className="w-3.5 h-3.5 text-emerald-500" />
                        Entitled
                      </>
                    ) : (
                      <>
                        <Zap className="w-3.5 h-3.5" />
                        Activate Add-On
                      </>
                    )}
                  </Button>
                </CardContent>
              </Card>
            );
          })}
        </div>
      </div>

      {/* Commercial Invoices History */}
      <div className="space-y-4 pt-4 border-t border-border">
        <div>
          <h3 className="text-base font-bold text-foreground flex items-center gap-2">
            <Receipt className="w-4 h-4 text-primary" />
            Billing Invoices & Commercial Orders
          </h3>
          <p className="text-xs text-muted-foreground">
            View transaction receipts, payment references, and download official PDF tax invoices.
          </p>
        </div>

        {invoices.length === 0 ? (
          <Card className="border-border shadow-sm p-6 text-center text-muted-foreground text-xs">
            No historical invoices found for this organization.
          </Card>
        ) : (
          <div className="border border-border rounded-xl overflow-hidden shadow-sm">
            <table className="w-full text-left text-xs">
              <thead className="bg-muted/50 border-b border-border text-muted-foreground">
                <tr>
                  <th className="p-3 font-medium">Invoice #</th>
                  <th className="p-3 font-medium">Description / Plan</th>
                  <th className="p-3 font-medium">Date Issued</th>
                  <th className="p-3 font-medium">Amount</th>
                  <th className="p-3 font-medium">Status</th>
                  <th className="p-3 font-medium text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {invoices.map((inv) => (
                  <tr key={inv.id} className="hover:bg-muted/20 transition-colors">
                    <td className="p-3 font-mono font-bold text-foreground">
                      {inv.invoiceNumber || inv.orderNumber || inv.id.substring(0, 8)}
                    </td>
                    <td className="p-3 text-muted-foreground">{inv.planCode || "Subscription Order"}</td>
                    <td className="p-3 text-muted-foreground">
                      {formatDate(inv.billingDate || inv.createdAt || "")}
                    </td>
                    <td className="p-3 font-mono font-bold text-foreground">
                      {formatCurrency(inv.amount || inv.total || 0, inv.currency)}
                    </td>
                    <td className="p-3">
                      <Badge
                        variant="outline"
                        className={`text-[10px] capitalize ${
                          inv.status === "paid"
                            ? "bg-emerald-500/10 text-emerald-500 border-emerald-500/20"
                            : "bg-amber-500/10 text-amber-500 border-amber-500/20"
                        }`}
                      >
                        {inv.status}
                      </Badge>
                    </td>
                    <td className="p-3 text-right">
                      {inv.pdfUrl ? (
                        <a
                          href={inv.pdfUrl}
                          target="_blank"
                          rel="noreferrer"
                          className="inline-flex items-center gap-1 text-primary hover:underline font-medium"
                        >
                          <Download className="w-3 h-3" />
                          PDF
                        </a>
                      ) : (
                        <span className="text-[11px] text-muted-foreground font-mono">Paid</span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Commercial Checkout Dialog */}
      <CommercialCheckoutDialog
        isOpen={isCheckoutOpen}
        onClose={() => setIsCheckoutOpen(false)}
        orgId={orgId}
        item={checkoutItem}
        onSuccess={handleCheckoutSuccess}
      />
    </div>
  );
}
