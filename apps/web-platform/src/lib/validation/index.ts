/**
 * Clinical and Form Input Validation Helpers
 */

export function isValidEmail(email: string): boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.trim());
}

export function isValidPhone(phone: string): boolean {
  const clean = phone.replace(/[\s\-\(\)\+]/g, "");
  return clean.length >= 8 && clean.length <= 15;
}

export function isValidMRN(mrn: string): boolean {
  return mrn.trim().length >= 4;
}

export function isValidNIN(nin: string): boolean {
  const clean = nin.replace(/\D/g, "");
  return clean.length === 11;
}

export function isFutureDate(dateStr: string): boolean {
  const date = new Date(dateStr);
  return date.getTime() > Date.now();
}
