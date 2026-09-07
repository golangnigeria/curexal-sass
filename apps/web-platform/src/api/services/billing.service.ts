import { authClient } from "@/lib/auth-client";
import { apiGet, apiPost } from "@/api/client";

export interface CapabilityPricePayload {
  id: string;
  code: string;
  name: string;
  description: string;
  category: string;
  monthlyPrice: number;
  annualPrice: number;
  currency: string;
  isAddOn: boolean;
  isActive: boolean;
}

export interface InvoicePayload {
  id: string;
  invoiceNumber?: string;
  orderNumber?: string;
  organizationId: string;
  planCode?: string;
  total?: number;
  amount?: number;
  currency: string;
  status: "paid" | "pending" | "overdue" | "failed";
  createdAt?: string;
  billingDate?: string;
  dueDate?: string;
  paidAt?: string;
  pdfUrl?: string;
}

export interface CreateOrderPayload {
  billingCycle: "monthly" | "annual";
  currency: string;
  items: Array<{
    capabilityId?: string;
    capabilityCode: string;
    billingCycle: "monthly" | "annual";
    quantity?: number;
  }>;
}

export interface CommercialOrderResponse {
  id: string;
  orderNumber: string;
  status: string;
  currency: string;
  subtotal: number;
  tax: number;
  total: number;
  paymentUrl?: string;
  providerReference?: string;
}

class BillingService {
  async getCapabilityPrices(currency = "NGN"): Promise<CapabilityPricePayload[]> {
    try {
      const res = await apiGet<any[]>("/marketplace/capabilities");
      return (res || []).map((c) => ({
        id: c.id,
        code: c.code,
        name: c.name,
        description: c.description,
        category: c.category || "addon",
        monthlyPrice: c.monthlyPrice || c.basePrice || 0,
        annualPrice: c.annualPrice || (c.monthlyPrice ? c.monthlyPrice * 10 : 0),
        currency: c.currency || currency,
        isAddOn: c.isAddOn ?? true,
        isActive: c.isActive ?? true,
      }));
    } catch {
      return [];
    }
  }

  async getInvoices(orgId: string): Promise<InvoicePayload[]> {
    try {
      const res = await apiGet<any[]>(`/organizations/${orgId}/marketplace/orders`);
      return (res || []).map((inv) => ({
        id: inv.id,
        invoiceNumber: inv.invoiceNumber || inv.orderNumber || `INV-${inv.id.substring(0, 8).toUpperCase()}`,
        organizationId: inv.organizationId || orgId,
        planCode: inv.planCode || "Add-On Capability",
        amount: inv.total || inv.amount || 0,
        currency: inv.currency || "NGN",
        status: inv.status || "paid",
        billingDate: inv.createdAt || new Date().toISOString(),
        dueDate: inv.dueAt || inv.createdAt || new Date().toISOString(),
        paidAt: inv.paidAt,
        pdfUrl: inv.pdfUrl,
      }));
    } catch {
      return [];
    }
  }

  async subscribeCapability(orgId: string, capabilityCode: string): Promise<void> {
    await apiPost<void>(`/organizations/${orgId}/marketplace/subscribe`, {
      capabilityCode,
    });
  }

  async createCommercialOrder(
    orgId: string,
    payload: CreateOrderPayload,
    provider = "mock"
  ): Promise<CommercialOrderResponse> {
    return apiPost<CommercialOrderResponse>(
      `/organizations/${orgId}/marketplace/orders?provider=${encodeURIComponent(provider)}`,
      payload
    );
  }
}

export const billingService = new BillingService();
