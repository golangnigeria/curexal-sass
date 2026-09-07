import { z } from "zod";

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

/**
 * Canonical Zod Schemas
 */
export const signInSchema = z.object({
  email: z.string().email("Please enter a valid email address"),
  password: z.string().min(6, "Password must be at least 6 characters"),
});

export const patientRegistrationSchema = z.object({
  firstName: z.string().min(2, "First name is required"),
  middleName: z.string().optional(),
  lastName: z.string().min(2, "Last name is required"),
  gender: z.enum(["MALE", "FEMALE", "OTHER"]),
  dateOfBirth: z.string().min(4, "Date of birth is required"),
  phone: z.string().min(8, "Valid phone number is required"),
  email: z.string().email("Invalid email").optional().or(z.literal("")),
  nin: z.string().optional(),
  bloodGroup: z.string().optional(),
  genotype: z.string().optional(),
  address: z.string().optional(),
});
