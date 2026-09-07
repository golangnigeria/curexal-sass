/**
 * Currency Formatting Utilities
 */

export function formatCurrency(amount: number, currency: string = "NGN"): string {
  const symbolMap: Record<string, string> = {
    NGN: "₦",
    USD: "$",
    EUR: "€",
    GBP: "£",
  };

  const symbol = symbolMap[currency.toUpperCase()] || `${currency} `;
  return `${symbol}${amount.toLocaleString(undefined, {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  })}`;
}

export function formatPrice(amount?: number, currency: string = "NGN"): string {
  if (amount == null) return "—";
  return formatCurrency(amount, currency);
}
