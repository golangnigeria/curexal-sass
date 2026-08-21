import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPut } from "@/api/client";
import type {
  PricingRule,
  PaymentGatewayConfig,
  UpdateGatewayPayload,
} from "@/api/contracts";

export function usePricingRules() {
  return useQuery({
    queryKey: ["platform", "pricing"],
    queryFn: () => apiGet<PricingRule[]>("/platform/pricing"),
  });
}

export function useUpdatePricingRule() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: PricingRule) =>
      apiPut<PricingRule>("/platform/pricing", payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["platform", "pricing"] });
    },
  });
}

export function usePaymentGateways() {
  return useQuery({
    queryKey: ["platform", "payment-gateways"],
    queryFn: () => apiGet<PaymentGatewayConfig[]>("/platform/payment-gateways"),
  });
}

export function usePaymentGateway(provider: string) {
  return useQuery({
    queryKey: ["platform", "payment-gateways", provider],
    queryFn: () => apiGet<PaymentGatewayConfig>(`/platform/payment-gateways/${provider}`),
    enabled: !!provider,
  });
}

export function useUpdatePaymentGateway(provider: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: UpdateGatewayPayload) =>
      apiPut<PaymentGatewayConfig>(`/platform/payment-gateways/${provider}`, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["platform", "payment-gateways"] });
    },
  });
}
