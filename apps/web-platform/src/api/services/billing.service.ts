import { authClient } from "@/lib/auth-client";

export interface CapabilityPricePayload {
  id: string;
  capabilityId: string;
  capabilityCode: string;
  name: string;
  module: string;
  currency: string;
  billingPeriod: string;
  price: number;
  isActive: boolean;
}

export interface InvoicePayload {
  id: string;
  invoiceNumber: string;
  organizationId: string;
  planCode: string;
  amount: number;
  currency: string;
  status: "paid" | "pending" | "overdue";
  billingDate: string;
  dueDate: string;
  paidAt?: string;
  pdfUrl?: string;
}

class BillingService {
  private async getCsrfHeader(): Promise<Record<string, string>> {
    const csrfToken = authClient.getCsrfToken?.() || "";
    return csrfToken ? { "X-CSRF-Token": csrfToken } : {};
  }

  async getCapabilityPrices(currency = "NGN"): Promise<CapabilityPricePayload[]> {
    const res = await fetch(`/api/v1/subscription/capabilities/prices?currency=${currency}`, {
      credentials: "include",
    });
    if (!res.ok) return [];
    return res.json();
  }

  async getInvoices(orgId: string): Promise<InvoicePayload[]> {
    const res = await fetch(`/api/v1/organizations/${orgId}/invoices`, {
      credentials: "include",
    });
    if (!res.ok) return [];
    return res.json();
  }

  async subscribeCapability(orgId: string, capabilityCode: string, currency = "NGN"): Promise<void> {
    const csrf = await this.getCsrfHeader();
    const res = await fetch(`/api/v1/subscription/capabilities/subscribe`, {
      method: "POST",
      headers: { "Content-Type": "application/json", ...csrf },
      credentials: "include",
      body: JSON.stringify({ organizationId: orgId, capabilityCode, currency }),
    });
    if (!res.ok) {
      const err = await res.json().catch(() => ({}));
      throw new Error(err.message || "Failed to activate capability subscription");
    }
  }
}

export const billingService = new BillingService();
