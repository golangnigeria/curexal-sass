import React, { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { billingService } from "@/api/services/billing.service";
import { formatCurrency } from "@/lib/utils";
import { toast } from "sonner";
import {
  CreditCard,
  CheckCircle2,
  ShieldCheck,
  ExternalLink,
  Loader2,
  Lock,
  Sparkles,
} from "lucide-react";

export interface CheckoutItem {
  code: string;
  name: string;
  price: number;
  currency: string;
  type: "plan" | "addon";
  billingCycle?: "monthly" | "annual";
}

interface CommercialCheckoutDialogProps {
  isOpen: boolean;
  onClose: () => void;
  orgId: string;
  item: CheckoutItem | null;
  onSuccess?: () => void;
}

export function CommercialCheckoutDialog({
  isOpen,
  onClose,
  orgId,
  item,
  onSuccess,
}: CommercialCheckoutDialogProps) {
  const [selectedProvider, setSelectedProvider] = useState<"paystack" | "flutterwave" | "stripe" | "mock">("paystack");
  const [isProcessing, setIsProcessing] = useState(false);
  const [completedOrder, setCompletedOrder] = useState<any | null>(null);

  if (!item) return null;

  const cycle = item.billingCycle || "monthly";
  const vat = Math.round(item.price * 0.075);
  const total = item.price + vat;

  const handleCheckout = async () => {
    setIsProcessing(true);
    try {
      if (selectedProvider === "mock") {
        // Direct subscription activation for mock / sandbox
        await billingService.subscribeCapability(orgId, item.code);
        toast.success(`Successfully activated ${item.name}!`, {
          description: "Entitlements and licenses have been provisioned to your organization.",
        });
        if (onSuccess) onSuccess();
        onClose();
      } else {
        const orderRes = await billingService.createCommercialOrder(
          orgId,
          {
            billingCycle: cycle,
            currency: item.currency || "NGN",
            items: [
              {
                capabilityCode: item.code,
                billingCycle: cycle,
                quantity: 1,
              },
            ],
          },
          selectedProvider
        );

        if (orderRes.paymentUrl) {
          // Open secure payment gateway
          window.open(orderRes.paymentUrl, "_blank");
          toast.success("Payment Gateway Initialized", {
            description: `Order ${orderRes.orderNumber} created. Complete checkout in the newly opened window.`,
          });
        } else {
          setCompletedOrder(orderRes);
          toast.success("Order Created Successfully!");
        }

        if (onSuccess) onSuccess();
      }
    } catch (err: any) {
      toast.error("Payment initialization failed: " + (err.message || "Network error"));
    } finally {
      setIsProcessing(false);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={(open) => { if (!open) { onClose(); setCompletedOrder(null); } }}>
      <DialogContent className="sm:max-w-[480px]">
        <DialogHeader>
          <div className="flex items-center gap-2 mb-1">
            <Badge variant="outline" className="text-[10px] uppercase font-mono bg-primary/10 text-primary border-primary/20">
              Commercial Checkout
            </Badge>
            <div className="flex items-center gap-1 text-[11px] text-muted-foreground ml-auto">
              <Lock className="w-3 h-3 text-emerald-500" />
              <span>TLS 256-Bit Encrypted</span>
            </div>
          </div>
          <DialogTitle className="text-lg font-bold tracking-tight text-foreground">
            Complete Subscription Order
          </DialogTitle>
          <DialogDescription className="text-xs">
            Review order summary and select your preferred payment gateway provider.
          </DialogDescription>
        </DialogHeader>

        {completedOrder ? (
          <div className="py-6 text-center space-y-3">
            <div className="w-12 h-12 rounded-full bg-emerald-500/10 text-emerald-500 flex items-center justify-center mx-auto">
              <CheckCircle2 className="w-6 h-6" />
            </div>
            <h3 className="text-base font-bold text-foreground">Order Initialized!</h3>
            <p className="text-xs text-muted-foreground">
              Order Reference: <span className="font-mono text-primary font-bold">{completedOrder.orderNumber}</span>
            </p>
            <Button onClick={() => { onClose(); setCompletedOrder(null); }} className="w-full text-xs">
              Done
            </Button>
          </div>
        ) : (
          <div className="space-y-5 py-2">
            {/* Order Item Summary */}
            <div className="p-3.5 rounded-xl border border-border bg-card/60 space-y-2.5">
              <div className="flex items-start justify-between">
                <div>
                  <h4 className="text-xs font-bold text-foreground">{item.name}</h4>
                  <p className="text-[11px] text-muted-foreground capitalize">
                    {item.type} • {cycle} Billing Cycle
                  </p>
                </div>
                <div className="text-right">
                  <span className="text-xs font-bold font-mono text-foreground">
                    {formatCurrency(item.price, item.currency)}
                  </span>
                </div>
              </div>

              <div className="border-t border-border/60 pt-2 space-y-1 text-xs">
                <div className="flex justify-between text-muted-foreground text-[11px]">
                  <span>Subtotal</span>
                  <span className="font-mono">{formatCurrency(item.price, item.currency)}</span>
                </div>
                <div className="flex justify-between text-muted-foreground text-[11px]">
                  <span>VAT (7.5%)</span>
                  <span className="font-mono">{formatCurrency(vat, item.currency)}</span>
                </div>
                <div className="flex justify-between font-bold text-foreground pt-1 border-t border-border/40">
                  <span>Total Due</span>
                  <span className="font-mono text-primary">{formatCurrency(total, item.currency)}</span>
                </div>
              </div>
            </div>

            {/* Provider Selection */}
            <div className="space-y-2">
              <label className="text-xs font-medium text-foreground">Select Payment Method</label>
              <div className="grid grid-cols-2 gap-2.5">
                {[
                  { id: "paystack", name: "Paystack", desc: "Card, Bank Transfer, USSD" },
                  { id: "flutterwave", name: "Flutterwave", desc: "Barter, M-Pesa, Card" },
                  { id: "stripe", name: "Stripe", desc: "International Visa / MC / Amex" },
                  { id: "mock", name: "Direct Sandbox", desc: "Instant Activation for Testing" },
                ].map((prov) => (
                  <button
                    key={prov.id}
                    type="button"
                    onClick={() => setSelectedProvider(prov.id as any)}
                    className={`p-2.5 rounded-xl border text-left transition-all ${
                      selectedProvider === prov.id
                        ? "border-primary bg-primary/5 shadow-sm"
                        : "border-border hover:border-border/80 bg-card/40"
                    }`}
                  >
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-xs font-bold text-foreground">{prov.name}</span>
                      {selectedProvider === prov.id && (
                        <CheckCircle2 className="w-3.5 h-3.5 text-primary" />
                      )}
                    </div>
                    <p className="text-[10px] text-muted-foreground leading-tight">{prov.desc}</p>
                  </button>
                ))}
              </div>
            </div>
          </div>
        )}

        {!completedOrder && (
          <DialogFooter className="gap-2 sm:gap-0">
            <Button variant="outline" size="sm" onClick={onClose} disabled={isProcessing} className="text-xs">
              Cancel
            </Button>
            <Button
              size="sm"
              onClick={handleCheckout}
              disabled={isProcessing}
              className="text-xs gap-1.5 bg-primary text-primary-foreground shadow"
            >
              {isProcessing ? (
                <Loader2 className="w-3.5 h-3.5 animate-spin" />
              ) : (
                <CreditCard className="w-3.5 h-3.5" />
              )}
              {isProcessing ? "Processing..." : `Pay ${formatCurrency(total, item.currency)}`}
            </Button>
          </DialogFooter>
        )}
      </DialogContent>
    </Dialog>
  );
}
