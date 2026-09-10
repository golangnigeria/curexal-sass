import React, { useState } from "react";
import { CapabilityGate } from "@/components/design-system/capability-gate";
import { StatusBadge } from "@/components/feedback/status-badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { toast } from "sonner";
import {
  CreditCard,
  Receipt,
  Plus,
  DollarSign,
  CheckCircle2,
  FileText,
  Printer,
  Sparkles,
  Search,
  User,
  ArrowRight,
  ShieldCheck,
  TrendingUp,
  Clock,
  Banknote,
  Smartphone,
} from "lucide-react";
import { useInvoices, useProcessPayment } from "@/api/hooks/use-billing";
import { usePatients } from "@/api/hooks/use-patients";
import type { PatientInvoice, PaymentReceipt } from "@/api/contracts";

export default function WorkspaceBillingPage() {
  const { data: liveInvoices, isLoading: isLoadingInvoices } = useInvoices();
  const processPaymentMutation = useProcessPayment();

  // Payment Settlement Modal State
  const [selectedInvoice, setSelectedInvoice] = useState<PatientInvoice | null>(null);
  const [tenderType, setTenderType] = useState<"CASH" | "POS" | "BANK_TRANSFER" | "SPLIT">("POS");
  const [cashAmount, setCashAmount] = useState<number>(0);
  const [posAmount, setPosAmount] = useState<number>(0);
  const [transferAmount, setTransferAmount] = useState<number>(0);
  const [singleAmount, setSingleAmount] = useState<number>(0);

  // Active Receipt Modal
  const [activeReceipt, setActiveReceipt] = useState<PaymentReceipt | null>(null);

  // Open checkout modal for an invoice
  const handleOpenCheckout = (inv: PatientInvoice) => {
    setSelectedInvoice(inv);
    setSingleAmount(inv.balanceDue);
    setPosAmount(inv.balanceDue);
    setCashAmount(0);
    setTransferAmount(0);
    setTenderType("POS");
  };

  const handleSettlePayment = async () => {
    if (!selectedInvoice) return;

    let payAmount = 0;
    let breakdown: Record<string, number> | undefined;

    if (tenderType === "SPLIT") {
      payAmount = cashAmount + posAmount + transferAmount;
      breakdown = {
        CASH: cashAmount,
        POS: posAmount,
        BANK_TRANSFER: transferAmount,
      };
      if (payAmount <= 0) {
        toast.error("Split tender breakdown must be greater than zero");
        return;
      }
    } else {
      payAmount = singleAmount;
      if (payAmount <= 0) {
        toast.error("Payment amount must be greater than zero");
        return;
      }
    }

    try {
      const receipt = await processPaymentMutation.mutateAsync({
        invoiceId: selectedInvoice.id,
        payload: {
          amount: payAmount,
          tenderType,
          tenderBreakdown: breakdown,
        },
      });

      setSelectedInvoice(null);
      setActiveReceipt(receipt);
      toast.success(`Payment Settled: Receipt ${receipt.receiptNumber}`, {
        description: `Collected ₦${payAmount.toLocaleString()} via ${tenderType}. Care Journey updated.`,
      });
    } catch (err: any) {
      toast.error(err?.response?.data?.message || err?.message || "Failed to process payment");
    }
  };

  const formatCurrency = (val: number) => `₦${(val || 0).toLocaleString()}`;

  // Summary Metrics
  const invoiceList = liveInvoices || [];
  const totalReceivables = invoiceList.reduce((acc, i) => acc + i.totalAmount, 0);
  const totalCollected = invoiceList.reduce((acc, i) => acc + i.amountPaid, 0);
  const totalOutstanding = invoiceList.reduce((acc, i) => acc + i.balanceDue, 0);
  const paidCount = invoiceList.filter((i) => i.status === "PAID").length;

  return (
    <CapabilityGate
      capability="core.billing"
      moduleCode="billing"
      title="Point of Sale Billing & Cashier Register"
      description="Cashier point of sale, split multi-tender payments, and instant thermal receipt generation."
    >
      <div className="space-y-6 animate-in fade-in duration-300">
        {/* Header */}
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-5">
          <div>
            <div className="flex items-center gap-2.5 mb-1">
              <h1 className="text-2xl font-bold tracking-tight text-foreground flex items-center gap-2">
                <CreditCard className="w-6 h-6 text-teal-600 dark:text-teal-400" />
                Cashier Point of Sale & Billing Register
              </h1>
              <Badge variant="outline" className="border-teal-500/40 text-teal-600 dark:text-teal-400 bg-teal-500/10 text-[10px] font-mono">
                POS Terminal #1 • Split Tender Active
              </Badge>
            </div>
            <p className="text-xs text-muted-foreground">
              Automated encounter invoice settlement, diagnostic test collection, and instant receipt generation.
            </p>
          </div>

          <div className="flex items-center gap-2.5">
            <Badge variant="secondary" className="gap-1.5 text-xs py-1 px-3">
              <ShieldCheck className="w-3.5 h-3.5 text-emerald-500" />
              Audit Log Enforced
            </Badge>
          </div>
        </div>

        {/* Operational Statistics */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <Card className="border-border shadow-sm bg-card">
            <CardContent className="p-4 flex items-center gap-3.5">
              <div className="w-10 h-10 rounded-xl bg-teal-500/10 text-teal-600 flex items-center justify-center shrink-0">
                <Banknote className="w-5 h-5" />
              </div>
              <div>
                <p className="text-xs text-muted-foreground font-medium">Total Billed</p>
                <h3 className="text-lg font-bold text-foreground font-mono">
                  {formatCurrency(totalReceivables)}
                </h3>
              </div>
            </CardContent>
          </Card>

          <Card className="border-border shadow-sm bg-card">
            <CardContent className="p-4 flex items-center gap-3.5">
              <div className="w-10 h-10 rounded-xl bg-emerald-500/10 text-emerald-600 flex items-center justify-center shrink-0">
                <CheckCircle2 className="w-5 h-5" />
              </div>
              <div>
                <p className="text-xs text-muted-foreground font-medium">Settled Revenue</p>
                <h3 className="text-lg font-bold text-emerald-600 dark:text-emerald-400 font-mono">
                  {formatCurrency(totalCollected)}
                </h3>
              </div>
            </CardContent>
          </Card>

          <Card className="border-border shadow-sm bg-card">
            <CardContent className="p-4 flex items-center gap-3.5">
              <div className="w-10 h-10 rounded-xl bg-amber-500/10 text-amber-600 flex items-center justify-center shrink-0">
                <Clock className="w-5 h-5" />
              </div>
              <div>
                <p className="text-xs text-muted-foreground font-medium">Outstanding Balance</p>
                <h3 className="text-lg font-bold text-amber-600 dark:text-amber-400 font-mono">
                  {formatCurrency(totalOutstanding)}
                </h3>
              </div>
            </CardContent>
          </Card>

          <Card className="border-border shadow-sm bg-card">
            <CardContent className="p-4 flex items-center gap-3.5">
              <div className="w-10 h-10 rounded-xl bg-sky-500/10 text-sky-600 flex items-center justify-center shrink-0">
                <TrendingUp className="w-5 h-5" />
              </div>
              <div>
                <p className="text-xs text-muted-foreground font-medium">Clearance Rate</p>
                <h3 className="text-lg font-bold text-foreground font-mono">
                  {invoiceList.length > 0 ? Math.round((paidCount / invoiceList.length) * 100) : 100}%
                </h3>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Live Patient Invoices Register */}
        <Card className="border-border shadow-sm bg-card">
          <CardHeader className="pb-3 border-b border-border flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-sm font-bold text-foreground flex items-center gap-2">
                <Receipt className="w-4 h-4 text-primary" />
                Active Cashier Invoices & Clinical Encounter Bills
              </CardTitle>
              <CardDescription className="text-xs">
                Invoices auto-generated from concluded doctor consultations, diagnostic labs, and dispensary items.
              </CardDescription>
            </div>
            <Badge variant="outline" className="text-[10px] font-mono">
              {invoiceList.length} Total Invoices
            </Badge>
          </CardHeader>
          <CardContent className="p-0">
            {isLoadingInvoices ? (
              <div className="p-8 text-center text-xs text-muted-foreground">
                Loading live billing invoices from POS register...
              </div>
            ) : invoiceList.length === 0 ? (
              <div className="p-12 text-center space-y-2">
                <Receipt className="w-8 h-8 text-muted-foreground mx-auto opacity-50" />
                <p className="text-xs font-semibold text-foreground">No Invoices Pending</p>
                <p className="text-[11px] text-muted-foreground">
                  Invoices are automatically created when physicians conclude clinical encounters.
                </p>
              </div>
            ) : (
              <div className="divide-y divide-border">
                {invoiceList.map((inv) => (
                  <div
                    key={inv.id}
                    className="p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3 hover:bg-secondary/10 transition-colors"
                  >
                    <div className="flex items-start sm:items-center gap-3">
                      <div className="w-9 h-9 rounded-xl bg-primary/10 text-primary flex items-center justify-center font-bold text-xs font-mono shrink-0">
                        <Receipt className="w-4 h-4" />
                      </div>
                      <div>
                        <div className="flex items-center gap-2">
                          <span className="text-xs font-bold text-foreground">
                            {inv.patientName || "Patient"}
                          </span>
                          <span className="text-[10px] font-mono text-muted-foreground">
                            ({inv.mrn || "PAT-UNKNOWN"})
                          </span>
                          <span className="text-[10px] font-mono font-semibold px-1.5 py-0.5 rounded bg-secondary text-secondary-foreground">
                            {inv.invoiceNumber}
                          </span>
                        </div>
                        <p className="text-[11px] text-muted-foreground mt-0.5">
                          {inv.encounterId ? "Physician Clinical Consultation & Orders" : "General Outpatient Service"} • {new Date(inv.createdAt).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}
                        </p>
                      </div>
                    </div>

                    <div className="flex items-center gap-4 justify-between sm:justify-end">
                      <div className="text-right">
                        <p className="text-xs font-bold text-foreground font-mono">
                          {formatCurrency(inv.totalAmount)}
                        </p>
                        {inv.balanceDue > 0 ? (
                          <p className="text-[10px] font-mono text-amber-600 font-semibold">
                            Due: {formatCurrency(inv.balanceDue)}
                          </p>
                        ) : (
                          <p className="text-[10px] font-mono text-emerald-600 font-semibold">
                            Fully Settled
                          </p>
                        )}
                      </div>

                      <div className="flex items-center gap-2">
                        <Badge
                          variant="outline"
                          className={`text-[10px] font-mono uppercase ${
                            inv.status === "PAID"
                              ? "border-emerald-500/40 text-emerald-600 bg-emerald-500/10"
                              : inv.status === "PARTIALLY_PAID"
                              ? "border-amber-500/40 text-amber-600 bg-amber-500/10"
                              : "border-rose-500/40 text-rose-600 bg-rose-500/10"
                          }`}
                        >
                          {inv.status}
                        </Badge>

                        {inv.balanceDue > 0 ? (
                          <Button
                            size="sm"
                            onClick={() => handleOpenCheckout(inv)}
                            className="text-xs h-7 gap-1 bg-emerald-600 hover:bg-emerald-700 text-white font-semibold"
                          >
                            <CreditCard className="w-3 h-3" />
                            Collect POS
                          </Button>
                        ) : (
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => {
                              setActiveReceipt({
                                receiptNumber: `RCP-${inv.invoiceNumber.replace("INV-", "")}`,
                                invoiceNumber: inv.invoiceNumber,
                                patientName: inv.patientName || "Patient",
                                mrn: inv.mrn || "PAT-0000",
                                amountPaid: inv.amountPaid,
                                previousBalance: inv.totalAmount,
                                newBalance: 0,
                                tenderType: "SETTLED",
                                status: "SETTLED",
                                paidAt: inv.updatedAt,
                                cashierName: "Cashier Desk",
                              });
                            }}
                            className="text-xs h-7 gap-1"
                          >
                            <Printer className="w-3 h-3" />
                            Receipt
                          </Button>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Payment Checkout Modal (Split Tender POS) */}
        <Dialog open={Boolean(selectedInvoice)} onOpenChange={(open) => !open && setSelectedInvoice(null)}>
          <DialogContent className="max-w-md">
            <DialogHeader>
              <DialogTitle className="text-base font-bold flex items-center gap-2">
                <CreditCard className="w-4 h-4 text-emerald-600" />
                POS Cashier Settlement
              </DialogTitle>
              <DialogDescription className="text-xs">
                Settle invoice <span className="font-mono font-semibold">{selectedInvoice?.invoiceNumber}</span> for{" "}
                <span className="font-semibold text-foreground">{selectedInvoice?.patientName}</span>
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-2">
              {/* Outstanding Amount Banner */}
              <div className="p-3 rounded-xl bg-secondary/30 border border-border flex items-center justify-between">
                <div>
                  <p className="text-[11px] text-muted-foreground">Total Bill / Balance Due</p>
                  <p className="text-lg font-bold text-foreground font-mono">
                    {formatCurrency(selectedInvoice?.balanceDue || 0)}
                  </p>
                </div>
                <Badge variant="outline" className="border-primary/30 text-primary font-mono text-xs">
                  Split Tender Ready
                </Badge>
              </div>

              {/* Tender Type Selection */}
              <div>
                <label className="text-xs font-semibold text-foreground mb-1.5 block">Payment Tender Type</label>
                <div className="grid grid-cols-4 gap-2">
                  {(["POS", "CASH", "BANK_TRANSFER", "SPLIT"] as const).map((t) => (
                    <Button
                      key={t}
                      type="button"
                      size="sm"
                      variant={tenderType === t ? "default" : "outline"}
                      onClick={() => setTenderType(t)}
                      className="text-xs h-8 font-mono uppercase"
                    >
                      {t === "BANK_TRANSFER" ? "Transfer" : t}
                    </Button>
                  ))}
                </div>
              </div>

              {/* Tender Amounts */}
              {tenderType === "SPLIT" ? (
                <div className="space-y-2.5 p-3 rounded-xl border border-border bg-secondary/10">
                  <p className="text-xs font-semibold text-foreground">Split Multi-Tender Breakdown</p>
                  <div className="grid grid-cols-3 gap-2">
                    <div>
                      <label className="text-[10px] text-muted-foreground block mb-1">Cash (NGN)</label>
                      <Input
                        type="number"
                        value={cashAmount}
                        onChange={(e) => setCashAmount(Number(e.target.value))}
                        className="text-xs h-8 font-mono"
                      />
                    </div>
                    <div>
                      <label className="text-[10px] text-muted-foreground block mb-1">POS Card (NGN)</label>
                      <Input
                        type="number"
                        value={posAmount}
                        onChange={(e) => setPosAmount(Number(e.target.value))}
                        className="text-xs h-8 font-mono"
                      />
                    </div>
                    <div>
                      <label className="text-[10px] text-muted-foreground block mb-1">Transfer (NGN)</label>
                      <Input
                        type="number"
                        value={transferAmount}
                        onChange={(e) => setTransferAmount(Number(e.target.value))}
                        className="text-xs h-8 font-mono"
                      />
                    </div>
                  </div>
                  <div className="flex justify-between text-xs pt-1 border-t border-border font-mono">
                    <span className="text-muted-foreground">Tendered Total:</span>
                    <span className="font-bold text-foreground">
                      {formatCurrency(cashAmount + posAmount + transferAmount)}
                    </span>
                  </div>
                </div>
              ) : (
                <div>
                  <label className="text-xs font-semibold text-foreground mb-1 block">Tender Amount (NGN)</label>
                  <Input
                    type="number"
                    value={singleAmount}
                    onChange={(e) => setSingleAmount(Number(e.target.value))}
                    className="text-xs h-8 font-mono"
                  />
                </div>
              )}
            </div>

            <DialogFooter className="gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setSelectedInvoice(null)}
                className="text-xs h-8"
              >
                Cancel
              </Button>
              <Button
                size="sm"
                onClick={handleSettlePayment}
                disabled={processPaymentMutation.isPending}
                className="text-xs h-8 gap-1.5 bg-emerald-600 hover:bg-emerald-700 text-white font-semibold"
              >
                <CheckCircle2 className="w-3.5 h-3.5" />
                {processPaymentMutation.isPending ? "Processing..." : "Authorize Settlement"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* Thermal Payment Receipt Modal */}
        <Dialog open={Boolean(activeReceipt)} onOpenChange={(open) => !open && setActiveReceipt(null)}>
          <DialogContent className="max-w-sm font-mono text-xs">
            <DialogHeader className="text-center pb-2 border-b border-dashed border-border">
              <DialogTitle className="text-sm font-bold uppercase tracking-wider text-center">
                Curexal Clinic OS
              </DialogTitle>
              <DialogDescription className="text-[11px] text-center text-muted-foreground">
                Official Cashier Payment Clearance Receipt
              </DialogDescription>
            </DialogHeader>

            {activeReceipt && (
              <div className="space-y-3 py-2 text-[11px]">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Receipt No:</span>
                  <span className="font-bold text-foreground">{activeReceipt.receiptNumber}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Invoice No:</span>
                  <span className="font-semibold text-foreground">{activeReceipt.invoiceNumber}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Patient:</span>
                  <span className="font-semibold text-foreground">{activeReceipt.patientName}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">MRN:</span>
                  <span>{activeReceipt.mrn}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Date/Time:</span>
                  <span>{new Date(activeReceipt.paidAt).toLocaleString()}</span>
                </div>

                <div className="border-t border-b border-dashed border-border py-2 space-y-1">
                  <div className="flex justify-between text-xs font-bold">
                    <span>Amount Paid:</span>
                    <span>{formatCurrency(activeReceipt.amountPaid)}</span>
                  </div>
                  <div className="flex justify-between text-muted-foreground">
                    <span>Tender Type:</span>
                    <span>{activeReceipt.tenderType}</span>
                  </div>
                  {activeReceipt.tenderBreakdown && (
                    <div className="text-[10px] text-muted-foreground pl-2">
                      {Object.entries(activeReceipt.tenderBreakdown).map(([k, v]) => (
                        <div key={k} className="flex justify-between">
                          <span>• {k}:</span>
                          <span>{formatCurrency(v as number)}</span>
                        </div>
                      ))}
                    </div>
                  )}
                  <div className="flex justify-between font-semibold pt-1">
                    <span>Remaining Balance:</span>
                    <span>{formatCurrency(activeReceipt.newBalance)}</span>
                  </div>
                </div>

                <div className="text-center text-[10px] text-muted-foreground pt-1">
                  <p>Cashier: {activeReceipt.cashierName}</p>
                  <p className="mt-1">Thank you for visiting Curexal Health.</p>
                  <p>Retain receipt for dispensary & pharmacy clearance.</p>
                </div>
              </div>
            )}

            <DialogFooter>
              <Button
                variant="outline"
                size="sm"
                onClick={() => {
                  window.print();
                }}
                className="w-full text-xs h-8 gap-1.5"
              >
                <Printer className="w-3.5 h-3.5" />
                Print Thermal Receipt
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </CapabilityGate>
  );
}
