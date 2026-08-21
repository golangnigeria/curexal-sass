import React, { useState } from "react";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { useSubscribeCapability } from "@/api/hooks/use-organization";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
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
} from "lucide-react";

const addOnCapabilitiesList = [
  {
    code: "laboratory.analyzer_integration",
    name: "LIMS Analyzer Middleware",
    desc: "ASTM / HL7 bidirectional interfacing for automated hematology & chemistry analyzers.",
    priceNGN: 25000,
    priceUSD: 20,
    module: "laboratory",
    popular: true,
  },
  {
    code: "radiology.pacs_dicom",
    name: "Radiology DICOM & PACS Suite",
    desc: "Direct modality worklists, PACS image routing, and browser DICOM viewer.",
    priceNGN: 35000,
    priceUSD: 30,
    module: "radiology",
    popular: true,
  },
  {
    code: "laboratory.advanced_qc",
    name: "Advanced QC & Levey-Jennings",
    desc: "Westgard rules, control lot tracking, and automatic delta error alerts.",
    priceNGN: 15000,
    priceUSD: 12,
    module: "laboratory",
    popular: false,
  },
  {
    code: "clinical.inpatient_wards",
    name: "Inpatient Bed & Ward Management",
    desc: "Bed allocation board, nurse handovers, and ward transfer workflows.",
    priceNGN: 30000,
    priceUSD: 25,
    module: "clinical",
    popular: false,
  },
  {
    code: "pharmacy.advanced_inventory",
    name: "FEFO Pharmacy Batch Tracking",
    desc: "First-Expired First-Out automation, supplier POs, and stock threshold alerts.",
    priceNGN: 20000,
    priceUSD: 18,
    module: "pharmacy",
    popular: false,
  },
  {
    code: "qms.iso15189",
    name: "ISO 15189 Quality & CAPA Suite",
    desc: "Full accreditation workflow, CAPA incident investigation, and audit readiness.",
    priceNGN: 40000,
    priceUSD: 35,
    module: "qms",
    popular: false,
  },
];

export default function OrganizationBillingPage() {
  const { data: bootstrap } = useBootstrap();
  const subscribeMutation = useSubscribeCapability();

  const orgPlan = bootstrap?.organization?.subscription || "smart";
  const currency = bootstrap?.workspace?.currency || "NGN";
  const activeCapabilities = bootstrap?.capabilities || [];
  const limits = bootstrap?.limits || { maxBranches: 1, maxMembers: 5, storageGb: 10 };

  const [loadingCode, setLoadingCode] = useState<string | null>(null);

  const handleSubscribe = async (capCode: string) => {
    setLoadingCode(capCode);
    try {
      await subscribeMutation.mutateAsync({ capabilityCode: capCode, currency });
      toast.success("Capability Add-On Activated!", {
        description: "Your organization entitlements and navigation have been updated immediately.",
      });
    } catch (err: any) {
      toast.error("Failed to activate add-on: " + (err.message || "Network error"));
    } finally {
      setLoadingCode(null);
    }
  };

  const formatPrice = (ngn: number, usd: number) => {
    if (currency === "USD") return `$${usd} / mo`;
    return `₦${ngn.toLocaleString()} / mo`;
  };

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
        {[
          {
            code: "smart",
            name: "Smart Starter",
            price: "Free / ₦0",
            desc: "Essential reception & patient registration.",
            features: ["1 Branch Facility", "5 Staff Member Seats", "10 GB Cloud Storage", "Basic Lab & Clinic"],
          },
          {
            code: "optimize",
            name: "Optimize Tier",
            price: "₦35,000 / mo",
            desc: "Growing diagnostic centers & multi-clinics.",
            features: ["3 Branch Facilities", "25 Staff Member Seats", "50 GB Cloud Storage", "Analyzer Interfacing"],
          },
          {
            code: "pro",
            name: "Pro Tier",
            price: "₦95,000 / mo",
            desc: "Full hospital & diagnostic laboratory networks.",
            features: ["10 Branch Facilities", "100 Staff Member Seats", "200 GB Cloud Storage", "DICOM PACS & LIS"],
          },
          {
            code: "enterprise",
            name: "Enterprise Custom",
            price: "Custom Contract",
            desc: "Tertiary hospital groups & health networks.",
            features: ["Unlimited Branches", "Unlimited Staff Seats", "5 TB Dedicated Storage", "White-label Custom Domain"],
          },
        ].map((plan) => {
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
                <div className="text-xl font-bold text-foreground mt-2">{plan.price}</div>
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
                  className={`w-full text-xs h-8 ${!isCurrent ? "bg-primary text-primary-foreground" : ""}`}
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
          {addOnCapabilitiesList.map((addon) => {
            const isOwned = activeCapabilities.includes(addon.code);
            const isLoadingThis = loadingCode === addon.code;

            return (
              <Card key={addon.code} className="border-border shadow-sm hover:shadow-md transition-all flex flex-col justify-between bg-card">
                <CardHeader className="pb-3">
                  <div className="flex items-center justify-between">
                    <Badge variant="outline" className="text-[9px] font-mono uppercase border-border">
                      {addon.module}
                    </Badge>
                    {isOwned ? (
                      <Badge className="text-[9px] bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/30">
                        Active License
                      </Badge>
                    ) : (
                      <span className="text-xs font-mono font-bold text-foreground">
                        {formatPrice(addon.priceNGN, addon.priceUSD)}
                      </span>
                    )}
                  </div>
                  <CardTitle className="text-sm font-bold text-foreground mt-2">{addon.name}</CardTitle>
                  <CardDescription className="text-xs min-h-[36px] text-muted-foreground">
                    {addon.desc}
                  </CardDescription>
                </CardHeader>
                <CardContent className="pt-0">
                  <Button
                    size="sm"
                    variant={isOwned ? "outline" : "default"}
                    disabled={isOwned || isLoadingThis}
                    onClick={() => handleSubscribe(addon.code)}
                    className={`w-full text-xs h-8 gap-1.5 ${
                      !isOwned ? "bg-primary text-primary-foreground shadow" : "text-muted-foreground"
                    }`}
                  >
                    {isOwned ? (
                      <>
                        <CheckCircle2 className="w-3.5 h-3.5 text-emerald-500" />
                        Entitled
                      </>
                    ) : isLoadingThis ? (
                      "Activating..."
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
    </div>
  );
}
