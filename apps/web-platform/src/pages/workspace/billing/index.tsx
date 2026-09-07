import React, { useState } from "react";
import { CapabilityGate } from "@/components/design-system/capability-gate";
import { DataTable } from "@/components/data-display/data-table";
import { StatusBadge } from "@/components/feedback/status-badge";
import { DocumentViewerModal } from "@/features/documents/document-viewer-modal";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
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
  Sparkles,
  Search,
  User,
} from "lucide-react";
import { usePatients } from "@/api/hooks/use-patients";

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

const initialInvoices: InvoiceItem[] = [
  {
    id: "inv-1",
    invoiceNo: "INV-2026-0841",
    patientName: "Amina Yusuf",
    patientMrn: "PAT-0012",
    itemsSummary: "Doctor Consultation + CBC + Malaria Smear",
    amount: 17500,
    paymentMethod: "POS",
    status: "paid",
    createdAt: "15 mins ago",
  },
  {
    id: "inv-2",
    invoiceNo: "INV-2026-0842",
    patientName: "Chinedu Okafor",
    patientMrn: "PAT-0034",
    itemsSummary: "General Outpatient Consultation Fee",
    amount: 10000,
    paymentMethod: "CASH",
    status: "paid",
    createdAt: "45 mins ago",
  },
  {
    id: "inv-3",
    invoiceNo: "INV-2026-0843",
    patientName: "Babatunde Lawal",
    patientMrn: "PAT-0078",
    itemsSummary: "Specialist Consultation + Follow-up",
    amount: 15000,
    paymentMethod: "HMO_INSURANCE",
    status: "pending",
    createdAt: "1 hour ago",
  },
];

export default function WorkspaceBillingPage() {
  const [invoices, setInvoices] = useState<InvoiceItem[]>(initialInvoices);
  const [activeReceiptDoc, setActiveReceiptDoc] = useState<any>(null);
  const [isCreatingInvoice, setIsCreatingInvoice] = useState(false);

  // New Invoice State
  const [patientSearch, setPatientSearch] = useState("");
  const [selectedPatientMrn, setSelectedPatientMrn] = useState("");
  const [patientName, setPatientName] = useState("");
  const [serviceCategory, setServiceCategory] = useState("Doctor Consultation");
  const [itemDescription, setItemDescription] = useState("Standard Outpatient Consultation");
  const [amount, setAmount] = useState<number>(10000);
  const [tenderMethod, setTenderMethod] = useState<"CASH" | "POS" | "TRANSFER" | "HMO_INSURANCE">("POS");

  // Query real patients from MPI
  const { data: patientList } = usePatients({
    query: patientSearch || undefined,
    limit: 5,
  });

  const handleSelectPatient = (p: any) => {
    setPatientName(`${p.firstName} ${p.lastName}`);
    setSelectedPatientMrn(p.mrn);
    setPatientSearch("");
  };

  const handleCreateInvoice = (e: React.FormEvent) => {
    e.preventDefault();
    if (!patientName.trim()) {
      toast.error("Please enter or select a patient from MPI");
      return;
    }

    const newInv: InvoiceItem = {
      id: `inv-${Date.now()}`,
      invoiceNo: `INV-2026-${Math.floor(1000 + Math.random() * 9000)}`,
      patientName,
      patientMrn: selectedPatientMrn || "PAT-0012",
      itemsSummary: itemDescription || serviceCategory,
      amount,
      paymentMethod: tenderMethod,
      status: "paid",
      createdAt: "Just now",
    };

    setInvoices([newInv, ...invoices]);
    setIsCreatingInvoice(false);
    setPatientName("");
    setSelectedPatientMrn("");
    setItemDescription("Standard Outpatient Consultation");
    toast.success("Payment Received & Cashier Receipt Issued!");
  };

  const handlePrintReceipt = (inv: InvoiceItem) => {
    setActiveReceiptDoc({
      title: `Cashier POS Receipt - ${inv.invoiceNo}`,
      documentType: "REPORT",
      uploadedAt: inv.createdAt,
      verifiedBy: "Facility Cashier Desk (Register #1)",
    });
  };

  const formatCurrency = (val: number) => `₦${val.toLocaleString()}`;

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
                POS Terminal #1 Active
              </Badge>
            </div>
            <p className="text-xs text-muted-foreground">
              Patient invoicing, diagnostic test payment collection, multi-tender settlement, and HMO claims.
            </p>
          </div>

          <div className="flex items-center gap-2.5">
            <Button
              size="sm"
              onClick={() => setIsCreatingInvoice(!isCreatingInvoice)}
              className="text-xs h-8 gap-1.5 bg-primary text-primary-foreground shadow"
            >
              <Plus className="w-3.5 h-3.5" />
              {isCreatingInvoice ? "Close Register" : "New Patient Invoice (F1)"}
            </Button>
          </div>
        </div>

        {/* New POS Invoice Form Drawer / Card */}
        {isCreatingInvoice && (
          <Card className="border-primary/40 bg-card shadow-md">
            <CardHeader className="pb-3 border-b border-border">
              <CardTitle className="text-sm font-bold text-foreground flex items-center gap-2">
                <Receipt className="w-4 h-4 text-primary" />
                New Cashier Transaction Register
              </CardTitle>
              <CardDescription className="text-xs">
                Select patient, add billable diagnostic investigations or clinical consult fees.
              </CardDescription>
            </CardHeader>
            <CardContent className="pt-4">
              <form onSubmit={handleCreateInvoice} className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4">
                <div className="relative">
                  <label className="text-xs font-semibold text-foreground mb-1 block">Patient Search (MPI)</label>
                  <div className="relative">
                    <Input
                      placeholder="Type name or MRN..."
                      value={patientName || patientSearch}
                      onChange={(e) => {
                        setPatientName(e.target.value);
                        setPatientSearch(e.target.value);
                      }}
                      className="text-xs h-8 pl-8"
                      required
                    />
                    <Search className="w-3.5 h-3.5 absolute left-2.5 top-2.5 text-muted-foreground" />
                  </div>
                  {patientSearch && patientList?.items && patientList.items.length > 0 && (
                    <div className="absolute top-full left-0 right-0 z-50 mt-1 bg-popover border border-border rounded-lg shadow-lg p-1 max-h-48 overflow-y-auto">
                      {patientList.items.map((p: any) => (
                        <button
                          key={p.id}
                          type="button"
                          onClick={() => handleSelectPatient(p)}
                          className="w-full text-left px-3 py-1.5 rounded text-xs hover:bg-accent hover:text-accent-foreground flex items-center justify-between"
                        >
                          <span className="font-semibold text-foreground">{p.firstName} {p.lastName}</span>
                          <span className="text-[10px] font-mono text-muted-foreground">{p.mrn}</span>
                        </button>
                      ))}
                    </div>
                  )}
                </div>
                <div>
                  <label className="text-xs font-semibold text-foreground mb-1 block">Service Item</label>
                  <select
                    value={itemDescription}
                    onChange={(e) => {
                      setItemDescription(e.target.value);
                      if (e.target.value.includes("Doctor Consultation")) setAmount(10000);
                      else if (e.target.value.includes("Specialist")) setAmount(20000);
                      else if (e.target.value.includes("Follow-up")) setAmount(5000);
                      else if (e.target.value.includes("Emergency")) setAmount(25000);
                    }}
                    className="w-full text-xs h-8 px-2 rounded-md border border-border bg-background"
                  >
                    <option value="Doctor Consultation (General Outpatient)">Doctor Consultation (General Outpatient) - ₦10,000</option>
                    <option value="Specialist Clinical Consultation">Specialist Clinical Consultation - ₦20,000</option>
                    <option value="Clinical Review & Follow-up">Clinical Review & Follow-up - ₦5,000</option>
                    <option value="Emergency Care Consultation & Triage">Emergency Care Consultation & Triage - ₦25,000</option>
                    <option value="Clinical Encounter Medication Plan">Clinical Encounter Medication Plan - ₦12,500</option>
                  </select>
                </div>
                <div>
                  <label className="text-xs font-semibold text-foreground mb-1 block">Tender Amount (NGN)</label>
                  <Input
                    type="number"
                    value={amount}
                    onChange={(e) => setAmount(Number(e.target.value))}
                    className="text-xs h-8 font-mono"
                    required
                  />
                </div>
                <div>
                  <label className="text-xs font-semibold text-foreground mb-1 block">Tender Method</label>
                  <select
                    value={tenderMethod}
                    onChange={(e) => setTenderMethod(e.target.value as any)}
                    className="w-full text-xs h-8 px-2 rounded-md border border-border bg-background"
                  >
                    <option value="POS">Card POS Terminal</option>
                    <option value="CASH">Cash</option>
                    <option value="TRANSFER">Direct Bank Transfer</option>
                    <option value="HMO_INSURANCE">HMO Insurance Co-Pay</option>
                  </select>
                </div>
                <div className="sm:col-span-2 md:col-span-4 flex justify-end gap-2 pt-2">
                  <Button
                    type="submit"
                    className="text-xs h-8 gap-1.5 bg-primary text-primary-foreground font-semibold shadow"
                  >
                    <CheckCircle2 className="w-3.5 h-3.5" />
                    Complete Tender & Print Receipt
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        )}

        {/* Invoice Register Table */}
        <DataTable
          data={invoices}
          searchPlaceholder="Search invoice number, patient name..."
          searchKey="invoiceNo"
          statusFilterKey="status"
          statusOptions={[
            { label: "Paid", value: "paid" },
            { label: "Pending Claims", value: "pending" },
          ]}
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
              header: "Total Amount",
              cell: (inv) => <span className="font-mono font-bold text-xs text-foreground">{formatCurrency(inv.amount)}</span>,
            },
            {
              header: "Payment Tender",
              cell: (inv) => (
                <Badge variant="outline" className="text-[10px] font-mono border-border">
                  {inv.paymentMethod.replace(/_/g, " ")}
                </Badge>
              ),
            },
            {
              header: "Status",
              cell: (inv) => <StatusBadge status={inv.status} />,
            },
            {
              header: "Actions",
              className: "text-right",
              cell: (inv) => (
                <div className="flex items-center justify-end gap-2">
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => handlePrintReceipt(inv)}
                    className="text-xs h-7 gap-1"
                  >
                    <FileText className="w-3 h-3" /> Receipt
                  </Button>
                </div>
              ),
            },
          ]}
        />

        {/* Universal Cashier Receipt Slip Preview */}
        <DocumentViewerModal
          isOpen={Boolean(activeReceiptDoc)}
          onClose={() => setActiveReceiptDoc(null)}
          document={activeReceiptDoc}
        />
      </div>
    </CapabilityGate>
  );
}
