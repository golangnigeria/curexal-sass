import React, { useState } from "react";
import { CapabilityGate } from "@/components/design-system/capability-gate";
import { DataTable } from "@/components/design-system/data-table";
import { StatusPill } from "@/components/design-system/status-pill";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import {
  CreditCard,
  Receipt,
  Plus,
  DollarSign,
  CheckCircle2,
  FileText,
  Printer,
} from "lucide-react";

interface InvoiceItem {
  id: string;
  invoiceNo: string;
  patientName: string;
  patientMrn: string;
  itemsSummary: string;
  amount: number;
  paymentMethod: "CASH" | "POS" | "TRANSFER" | "HMO_INSURANCE";
  status: "paid" | "pending" | "overdue";
  createdAt: string;
}

const mockInvoices: InvoiceItem[] = [
  { id: "inv-1", invoiceNo: "INV-2026-0841", patientName: "Amina Yusuf", patientMrn: "PAT-0012", itemsSummary: "Complete Blood Count + Doctor Consult", amount: 17500, paymentMethod: "POS", status: "paid", createdAt: "15 mins ago" },
  { id: "inv-2", invoiceNo: "INV-2026-0842", patientName: "Chinedu Okafor", patientMrn: "PAT-0034", itemsSummary: "Lipid Profile Panel + Fasting Blood Sugar", amount: 14000, paymentMethod: "CASH", status: "paid", createdAt: "45 mins ago" },
  { id: "inv-3", invoiceNo: "INV-2026-0843", patientName: "Babatunde Lawal", patientMrn: "PAT-0078", itemsSummary: "Chest X-Ray PA + Emergency Consult", amount: 26500, paymentMethod: "HMO_INSURANCE", status: "pending", createdAt: "1 hour ago" },
];

export default function WorkspaceBillingPage() {
  const [invoices, setInvoices] = useState<InvoiceItem[]>(mockInvoices);

  const handlePrintReceipt = (invNo: string) => {
    toast.info(`Generating Thermal Receipt for ${invNo}...`);
  };

  return (
    <CapabilityGate
      capability="core.billing"
      moduleCode="billing"
      title="Point of Sale Billing & Cashier Register"
      description="Unlock cashier point of sale, invoice receipts, split payments, and insurance claim adjudication."
    >
      <div className="space-y-8 animate-fade-in">
        {/* Header */}
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-6">
          <div>
            <div className="flex items-center gap-2.5 mb-1">
              <h1 className="text-2xl font-bold tracking-tight text-foreground flex items-center gap-2">
                <CreditCard className="w-6 h-6 text-primary" />
                Cashier Point of Sale & Billing
              </h1>
              <Badge variant="outline" className="border-primary/40 text-primary bg-primary/10 text-[10px] font-mono">
                POS Terminal
              </Badge>
            </div>
            <p className="text-xs text-muted-foreground">
              Patient invoicing, diagnostic test payment collection, and HMO insurance claims.
            </p>
          </div>

          <div className="flex items-center gap-2.5">
            <Button size="sm" className="text-xs h-9 gap-1.5 bg-primary text-primary-foreground shadow">
              <Plus className="w-3.5 h-3.5" />
              New Patient Invoice (F1)
            </Button>
          </div>
        </div>

        {/* Invoice Register Table */}
        <DataTable
          data={invoices}
          searchPlaceholder="Search invoice number, patient name..."
          searchKey="invoiceNo"
          columns={[
            {
              header: "Invoice #",
              cell: (inv) => <span className="font-mono font-bold text-xs">{inv.invoiceNo}</span>,
            },
            {
              header: "Patient",
              cell: (inv) => (
                <div>
                  <p className="font-semibold text-foreground">{inv.patientName}</p>
                  <p className="text-[11px] font-mono text-muted-foreground">{inv.patientMrn}</p>
                </div>
              ),
            },
            {
              header: "Services Billed",
              cell: (inv) => <span className="text-xs text-muted-foreground truncate max-w-xs">{inv.itemsSummary}</span>,
            },
            {
              header: "Amount (NGN)",
              cell: (inv) => <span className="font-mono font-bold text-xs text-foreground">₦{inv.amount.toLocaleString()}</span>,
            },
            {
              header: "Payment Method",
              cell: (inv) => (
                <Badge variant="outline" className="text-[9px] font-mono">
                  {inv.paymentMethod.replace("_", " ")}
                </Badge>
              ),
            },
            {
              header: "Status",
              cell: (inv) => <StatusPill status={inv.status} />,
            },
            {
              header: "Actions",
              className: "text-right",
              cell: (inv) => (
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => handlePrintReceipt(inv.invoiceNo)}
                  className="text-xs h-7 gap-1"
                >
                  <Printer className="w-3 h-3" /> Print Receipt
                </Button>
              ),
            },
          ]}
        />
      </div>
    </CapabilityGate>
  );
}
