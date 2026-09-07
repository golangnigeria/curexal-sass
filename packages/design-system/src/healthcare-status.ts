/**
 * Curexal Semantic Healthcare Status System
 * Unified taxonomy across Clinical Acuity, Workflow Lifecycle, Verification, and Product Lifecycles (LIS, RIS, Pharmacy, Billing).
 */

export type ClinicalAcuity = "routine" | "urgent" | "emergency" | "critical" | "resolved";

export type WorkflowLifecycle =
  | "pending"
  | "in_progress"
  | "completed"
  | "cancelled"
  | "rejected"
  | "failed";

export type VerificationState =
  | "unverified"
  | "under_review"
  | "verified"
  | "expired"
  | "approved"
  | "rejected";

export type LISLifecycle =
  | "collected"
  | "accessioned"
  | "in_analysis"
  | "results_pending"
  | "authorized"
  | "rejected";

export type PharmacyLifecycle =
  | "pending_verification"
  | "verified"
  | "dispensing"
  | "dispensed"
  | "cancelled";

export type BillingLifecycle =
  | "draft"
  | "pending"
  | "partially_paid"
  | "paid"
  | "overdue"
  | "voided";

export interface StatusMeta {
  label: string;
  variant: "routine" | "urgent" | "emergency" | "critical" | "pending" | "in_progress" | "completed" | "cancelled" | "rejected";
  dotColorClass?: string;
}

export function resolveStatusMeta(status: string | undefined): StatusMeta {
  const norm = (status || "").toLowerCase().replace(/[\s-]/g, "_");

  switch (norm) {
    // Clinical Acuity
    case "emergency":
    case "red":
    case "critical":
      return { label: "Emergency (Immediate)", variant: "emergency", dotColorClass: "bg-rose-500" };
    case "urgent":
    case "yellow":
      return { label: "Urgent Priority", variant: "urgent", dotColorClass: "bg-amber-500" };
    case "routine":
    case "green":
    case "standard":
      return { label: "Routine / Standard", variant: "routine", dotColorClass: "bg-teal-500" };
    case "resolved":
      return { label: "Resolved", variant: "completed", dotColorClass: "bg-emerald-500" };

    // LIS Lifecycle
    case "collected":
      return { label: "Sample Collected", variant: "pending", dotColorClass: "bg-slate-400" };
    case "accessioned":
      return { label: "Accessioned (LIMS)", variant: "in_progress", dotColorClass: "bg-blue-500" };
    case "in_analysis":
      return { label: "In Analyzer Queue", variant: "in_progress", dotColorClass: "bg-indigo-500" };
    case "results_pending":
      return { label: "Results Ready for Signoff", variant: "urgent", dotColorClass: "bg-amber-500" };
    case "authorized":
    case "verified":
    case "signed_off":
      return { label: "Electronically Authorized", variant: "completed", dotColorClass: "bg-emerald-500" };

    // Pharmacy Lifecycle
    case "pending_verification":
    case "pending":
    case "waiting":
    case "submitted":
      return { label: "Pending Verification", variant: "pending", dotColorClass: "bg-slate-400" };
    case "dispensing":
    case "in_consultation":
    case "in_progress":
      return { label: "In Progress", variant: "in_progress", dotColorClass: "bg-sky-500" };
    case "dispensed":
    case "completed":
    case "paid":
    case "active":
      return { label: "Completed", variant: "completed", dotColorClass: "bg-emerald-500" };
    case "partially_paid":
      return { label: "Partially Paid", variant: "urgent", dotColorClass: "bg-amber-500" };
    case "overdue":
    case "rejected":
    case "failed":
      return { label: "Action Required / Overdue", variant: "rejected", dotColorClass: "bg-rose-500" };
    case "cancelled":
    case "voided":
      return { label: "Cancelled", variant: "cancelled", dotColorClass: "bg-slate-400" };

    default:
      return {
        label: (status || "Unknown").replace(/_/g, " "),
        variant: "pending",
        dotColorClass: "bg-slate-400",
      };
  }
}
