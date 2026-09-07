export type PaymentTenderMethod = "CASH" | "POS" | "TRANSFER" | "HMO_INSURANCE" | "SPLIT";

export interface BillingInvoice {
  id: string;
  invoiceNumber: string;
  patientId: string;
  totalAmount: number;
  taxAmount?: number;
  discountAmount?: number;
  netPayable: number;
  tenderMethod: PaymentTenderMethod;
  status: "draft" | "pending" | "paid" | "partially_paid" | "voided";
  lineItems: Array<{
    description: string;
    category: "CONSULTATION" | "LAB" | "RADIOLOGY" | "PHARMACY" | "PROCEDURE";
    unitPrice: number;
    quantity: number;
    subtotal: number;
  }>;
  paidAt?: string;
  cashierName?: string;
}
