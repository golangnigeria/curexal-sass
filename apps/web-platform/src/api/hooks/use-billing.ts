import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  billingService,
  type InvoiceFilter,
} from "../services/billing.service";
import type { ProcessPaymentPayload } from "../contracts";

export const billingKeys = {
  all: ["billing"] as const,
  invoices: (filter?: InvoiceFilter) => [...billingKeys.all, "invoices", filter] as const,
  invoice: (id: string) => [...billingKeys.all, "invoice", id] as const,
  receipt: (number: string) => [...billingKeys.all, "receipt", number] as const,
};

/**
 * Hook to list clinic invoices
 */
export function useInvoices(filter?: InvoiceFilter) {
  return useQuery({
    queryKey: billingKeys.invoices(filter),
    queryFn: () => billingService.listInvoices(filter),
  });
}

/**
 * Hook to get detailed invoice
 */
export function useInvoice(invoiceId: string) {
  return useQuery({
    queryKey: billingKeys.invoice(invoiceId),
    queryFn: () => billingService.getInvoiceById(invoiceId),
    enabled: Boolean(invoiceId),
  });
}

/**
 * Hook to process POS cashier payment
 */
export function useProcessPayment() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      invoiceId,
      payload,
    }: {
      invoiceId: string;
      payload: ProcessPaymentPayload;
    }) => billingService.processPayment(invoiceId, payload),
    onSuccess: (data, vars) => {
      queryClient.invalidateQueries({ queryKey: billingKeys.all });
      queryClient.invalidateQueries({ queryKey: ["care-orchestration"] });
    },
  });
}

/**
 * Hook to retrieve payment receipt
 */
export function usePaymentReceipt(receiptNumber: string) {
  return useQuery({
    queryKey: billingKeys.receipt(receiptNumber),
    queryFn: () => billingService.getPaymentReceipt(receiptNumber),
    enabled: Boolean(receiptNumber),
  });
}
