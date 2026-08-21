import React, { useState } from "react";
import {
  CreditCard,
  Shield,
  CheckCircle2,
  Lock,
  Eye,
  EyeOff,
  Edit2,
  RefreshCw,
  Plus,
} from "lucide-react";
import {
  usePricingRules,
  useUpdatePricingRule,
  usePaymentGateways,
  useUpdatePaymentGateway,
} from "@/api/hooks/use-pricing";
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
} from "@/components/ui/dialog";
import { Switch } from "@/components/ui/switch";
import { formatCurrency } from "@/lib/utils";
import type { PricingRule, PaymentGatewayConfig } from "@/api/contracts";
import { toast } from "sonner";

export default function PricingPage() {
  const { data: rules, isLoading: rulesLoading, refetch: refetchRules } = usePricingRules();
  const { data: gateways, isLoading: gwLoading, refetch: refetchGateways } = usePaymentGateways();

  const updateRuleMutation = useUpdatePricingRule();

  // Pricing Rule Edit State
  const [editingRule, setEditingRule] = useState<PricingRule | null>(null);
  const [ruleMonthly, setRuleMonthly] = useState(0);
  const [ruleAnnual, setRuleAnnual] = useState(0);
  const [ruleVAT, setRuleVAT] = useState(0);
  const [ruleActive, setRuleActive] = useState(true);

  // Gateway Edit State
  const [editingGateway, setEditingGateway] = useState<PaymentGatewayConfig | null>(null);
  const [gwName, setGwName] = useState("");
  const [gwEnabled, setGwEnabled] = useState(false);
  const [gwPriority, setGwPriority] = useState(1);
  const [gwSecretKey, setGwSecretKey] = useState("");
  const [gwPublicKey, setGwPublicKey] = useState("");
  const [gwWebhookSecret, setGwWebhookSecret] = useState("");

  const updateGwMutation = useUpdatePaymentGateway(editingGateway?.providerCode || "");

  const handleOpenEditRule = (rule: PricingRule) => {
    setEditingRule(rule);
    setRuleMonthly(rule.monthlyPrice);
    setRuleAnnual(rule.annualPrice);
    setRuleVAT(rule.vatPercentage);
    setRuleActive(rule.isActive);
  };

  const handleSaveRule = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingRule) return;

    try {
      await updateRuleMutation.mutateAsync({
        ...editingRule,
        monthlyPrice: ruleMonthly,
        annualPrice: ruleAnnual,
        vatPercentage: ruleVAT,
        isActive: ruleActive,
      });
      toast.success(`Pricing rule for '${editingRule.targetCode}' updated!`);
      setEditingRule(null);
      refetchRules();
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Failed to update pricing rule");
    }
  };

  const handleOpenEditGateway = (gw: PaymentGatewayConfig) => {
    setEditingGateway(gw);
    setGwName(gw.name);
    setGwEnabled(gw.isEnabled);
    setGwPriority(gw.priority);
    setGwSecretKey("");
    setGwPublicKey(gw.publicKey || "");
    setGwWebhookSecret(gw.webhookSecret || "");
  };

  const handleSaveGateway = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingGateway) return;

    try {
      await updateGwMutation.mutateAsync({
        name: gwName,
        isEnabled: gwEnabled,
        priority: gwPriority,
        publicKey: gwPublicKey || undefined,
        secretKey: gwSecretKey || undefined,
        webhookSecret: gwWebhookSecret || undefined,
        version: editingGateway.version,
      });
      toast.success(`Payment gateway '${editingGateway.name}' updated!`);
      setEditingGateway(null);
      refetchGateways();
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Failed to update gateway");
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-foreground">
          Pricing Rules & Payment Gateways
        </h1>
        <p className="text-sm text-muted-foreground mt-1">
          Manage commercial subscription pricing tiers, capability fees, and encrypted payment gateway vaults.
        </p>
      </div>

      <Tabs defaultValue="pricing" className="space-y-4">
        <TabsList className="bg-muted/60 p-1 border border-border">
          <TabsTrigger value="pricing" className="text-xs">Subscription & Capability Pricing</TabsTrigger>
          <TabsTrigger value="gateways" className="text-xs">Payment Gateway Vault</TabsTrigger>
        </TabsList>

        {/* TAB 1: PRICING RULES */}
        <TabsContent value="pricing" className="space-y-4">
          <Card className="card-enterprise overflow-hidden">
            <CardHeader className="flex flex-row items-center justify-between pb-3">
              <div>
                <CardTitle className="text-sm font-semibold">Configured Pricing Rules</CardTitle>
                <CardDescription className="text-xs">
                  Base rates with optimistic locking across currencies and subscription cycles.
                </CardDescription>
              </div>
            </CardHeader>
            <CardContent className="p-0">
              <table className="w-full text-left text-xs">
                <thead className="border-y border-border bg-muted/40 font-semibold uppercase tracking-wider text-muted-foreground">
                  <tr>
                    <th className="py-2.5 px-4">Target Type</th>
                    <th className="py-2.5 px-4">Target Code</th>
                    <th className="py-2.5 px-4">Monthly Price</th>
                    <th className="py-2.5 px-4">Annual Price</th>
                    <th className="py-2.5 px-4">VAT (%)</th>
                    <th className="py-2.5 px-4">Status</th>
                    <th className="py-2.5 px-4 text-right">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {rulesLoading ? (
                    <tr>
                      <td colSpan={7} className="py-8 text-center text-muted-foreground">
                        Loading pricing rules...
                      </td>
                    </tr>
                  ) : !rules || rules.length === 0 ? (
                    <tr>
                      <td colSpan={7} className="py-8 text-center text-muted-foreground">
                        No pricing rules configured.
                      </td>
                    </tr>
                  ) : (
                    rules.map((rule) => (
                      <tr key={rule.targetCode + rule.currency} className="hover:bg-muted/20">
                        <td className="py-3 px-4 capitalize font-medium text-foreground">
                          {rule.targetType}
                        </td>
                        <td className="py-3 px-4 font-mono font-semibold text-primary">
                          {rule.targetCode}
                        </td>
                        <td className="py-3 px-4 font-medium">
                          {formatCurrency(rule.monthlyPrice, rule.currency)}
                        </td>
                        <td className="py-3 px-4 font-medium">
                          {formatCurrency(rule.annualPrice, rule.currency)}
                        </td>
                        <td className="py-3 px-4 font-mono">
                          {rule.vatPercentage}%
                        </td>
                        <td className="py-3 px-4">
                          {rule.isActive ? (
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
                            onClick={() => handleOpenEditRule(rule)}
                            className="h-7 text-xs gap-1"
                          >
                            <Edit2 className="h-3 w-3" /> Edit
                          </Button>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </CardContent>
          </Card>

          {/* Edit Pricing Rule Modal */}
          <Dialog open={!!editingRule} onOpenChange={(open) => !open && setEditingRule(null)}>
            <DialogContent className="sm:max-w-md">
              <DialogHeader>
                <DialogTitle className="text-base">Edit Pricing Rule</DialogTitle>
                <DialogDescription className="text-xs">
                  Update commercial pricing for <code className="font-mono text-primary">{editingRule?.targetCode}</code>.
                </DialogDescription>
              </DialogHeader>

              <form onSubmit={handleSaveRule} className="space-y-3 py-2 text-xs">
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1">
                    <Label htmlFor="monthlyPrice">Monthly Price ({editingRule?.currency})</Label>
                    <Input
                      id="monthlyPrice"
                      type="number"
                      min={0}
                      value={ruleMonthly}
                      onChange={(e) => setRuleMonthly(parseFloat(e.target.value))}
                      className="text-xs"
                      required
                    />
                  </div>
                  <div className="space-y-1">
                    <Label htmlFor="annualPrice">Annual Price ({editingRule?.currency})</Label>
                    <Input
                      id="annualPrice"
                      type="number"
                      min={0}
                      value={ruleAnnual}
                      onChange={(e) => setRuleAnnual(parseFloat(e.target.value))}
                      className="text-xs"
                      required
                    />
                  </div>
                </div>

                <div className="space-y-1">
                  <Label htmlFor="vatPercentage">VAT Percentage (%)</Label>
                  <Input
                    id="vatPercentage"
                    type="number"
                    min={0}
                    max={100}
                    step={0.1}
                    value={ruleVAT}
                    onChange={(e) => setRuleVAT(parseFloat(e.target.value))}
                    className="text-xs"
                    required
                  />
                </div>

                <div className="flex items-center justify-between pt-2">
                  <Label htmlFor="ruleActive" className="cursor-pointer">Active in Marketplace</Label>
                  <Switch
                    id="ruleActive"
                    checked={ruleActive}
                    onCheckedChange={setRuleActive}
                  />
                </div>

                <DialogFooter className="pt-3">
                  <Button type="button" variant="outline" size="sm" onClick={() => setEditingRule(null)} className="text-xs">
                    Cancel
                  </Button>
                  <Button type="submit" size="sm" className="bg-primary text-primary-foreground text-xs">
                    Save Pricing Rule
                  </Button>
                </DialogFooter>
              </form>
            </DialogContent>
          </Dialog>
        </TabsContent>

        {/* TAB 2: PAYMENT GATEWAY VAULT */}
        <TabsContent value="gateways" className="space-y-4">
          <div className="grid gap-4 sm:grid-cols-2">
            {gwLoading ? (
              <div className="col-span-2 py-8 text-center text-xs text-muted-foreground">
                Loading payment gateway integrations...
              </div>
            ) : !gateways || gateways.length === 0 ? (
              <div className="col-span-2 py-8 text-center text-xs text-muted-foreground">
                No payment gateways found in database vault.
              </div>
            ) : (
              gateways.map((gw) => (
                <Card key={gw.providerCode} className="card-enterprise flex flex-col justify-between">
                  <CardHeader className="pb-3">
                    <div className="flex items-center justify-between">
                      <CardTitle className="text-sm font-semibold flex items-center gap-2">
                        <Lock className="h-4 w-4 text-primary" />
                        {gw.name}
                      </CardTitle>
                      {gw.isEnabled ? (
                        <Badge variant="outline" className="border-emerald-500/30 text-emerald-600 bg-emerald-500/5 text-[10px]">
                          Enabled
                        </Badge>
                      ) : (
                        <Badge variant="outline" className="border-border text-muted-foreground text-[10px]">
                          Disabled
                        </Badge>
                      )}
                    </div>
                    <CardDescription className="text-xs font-mono">
                      provider: {gw.providerCode} • priority: {gw.priority}
                    </CardDescription>
                  </CardHeader>

                  <CardContent className="space-y-3 pt-0 text-xs">
                    <div className="rounded-lg border border-border bg-secondary/30 p-2.5 space-y-1.5 font-mono text-[11px]">
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Public Key:</span>
                        <span className="truncate max-w-[180px]">{gw.publicKey || "—"}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Secret Key:</span>
                        <span>••••••••••••••••</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Currencies:</span>
                        <span>{gw.supportedCurrencies?.join(", ") || "NGN, USD"}</span>
                      </div>
                    </div>

                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleOpenEditGateway(gw)}
                      className="w-full text-xs h-8 gap-1.5"
                    >
                      <Edit2 className="h-3 w-3" /> Configure Vault Secrets
                    </Button>
                  </CardContent>
                </Card>
              ))
            )}
          </div>

          {/* Edit Gateway Vault Modal */}
          <Dialog open={!!editingGateway} onOpenChange={(open) => !open && setEditingGateway(null)}>
            <DialogContent className="sm:max-w-md">
              <DialogHeader>
                <DialogTitle className="text-base font-semibold">
                  Configure {editingGateway?.name} Vault
                </DialogTitle>
                <DialogDescription className="text-xs">
                  Encrypted credentials for commercial transaction processing.
                </DialogDescription>
              </DialogHeader>

              <form onSubmit={handleSaveGateway} className="space-y-3 py-2 text-xs">
                <div className="space-y-1">
                  <Label htmlFor="gwName">Display Name</Label>
                  <Input
                    id="gwName"
                    value={gwName}
                    onChange={(e) => setGwName(e.target.value)}
                    className="text-xs"
                    required
                  />
                </div>

                <div className="space-y-1">
                  <Label htmlFor="gwPubKey">Public Key</Label>
                  <Input
                    id="gwPubKey"
                    placeholder="pk_live_..."
                    value={gwPublicKey}
                    onChange={(e) => setGwPublicKey(e.target.value)}
                    className="text-xs font-mono"
                  />
                </div>

                <div className="space-y-1">
                  <Label htmlFor="gwSecKey">Secret Key (Encrypted at rest)</Label>
                  <Input
                    id="gwSecKey"
                    type="password"
                    placeholder="Leave empty to retain existing secret"
                    value={gwSecretKey}
                    onChange={(e) => setGwSecretKey(e.target.value)}
                    className="text-xs font-mono"
                  />
                </div>

                <div className="space-y-1">
                  <Label htmlFor="gwWebhook">Webhook Signing Secret</Label>
                  <Input
                    id="gwWebhook"
                    type="password"
                    placeholder="whsec_..."
                    value={gwWebhookSecret}
                    onChange={(e) => setGwWebhookSecret(e.target.value)}
                    className="text-xs font-mono"
                  />
                </div>

                <div className="flex items-center justify-between pt-2">
                  <Label htmlFor="gwEnabled" className="cursor-pointer">Enable Payment Provider</Label>
                  <Switch
                    id="gwEnabled"
                    checked={gwEnabled}
                    onCheckedChange={setGwEnabled}
                  />
                </div>

                <DialogFooter className="pt-3">
                  <Button type="button" variant="outline" size="sm" onClick={() => setEditingGateway(null)} className="text-xs">
                    Cancel
                  </Button>
                  <Button type="submit" size="sm" className="bg-primary text-primary-foreground text-xs">
                    Save Gateway
                  </Button>
                </DialogFooter>
              </form>
            </DialogContent>
          </Dialog>
        </TabsContent>
      </Tabs>
    </div>
  );
}
