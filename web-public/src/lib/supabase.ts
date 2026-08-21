import { createClient } from "@supabase/supabase-js";

const supabaseUrl = import.meta.env.VITE_SUPABASE_URL || "https://xrqwupliqiuotzbkcrja.supabase.co";
const supabaseKey =
  import.meta.env.VITE_SUPABASE_PUBLISHABLE_KEY ||
  import.meta.env.VITE_SUPABASE_ANON_KEY ||
  "sb_publishable_-yewa_hvqkBJQ8fuJ9nSCQ_sc0J-_iC";

export const supabase = createClient(supabaseUrl, supabaseKey);

export interface WaitlistSubmission {
  fullName: string;
  email: string;
  phone?: string;
  country?: string;
  state?: string;
  city?: string;
  persona?: string;
  organization?: string;
  biggestPainPoint?: string;
  desiredFeatures?: string[];
  timeline?: string;
  shapingPreference?: string;
  // Dynamic Customer Discovery Research Fields
  orgType?: string;
  externalPartners?: string[];
  referralBreakdownStory?: string;
  visibilityLossPoints?: string[];
  communicationChannels?: string[];
  unresolvedImpact?: string;
  topWorkflowChoice?: string;
  trustFactors?: string[];
  disconnectedArea?: string;
  placesInvolved?: string;
  selfCoordinatedTasks?: string;
  automaticCoordinationDesire?: string;
  testPilotInterest?: string;
}

export type WaitlistSubmissionStatus =
  | "SUCCESS"
  | "DUPLICATE"
  | "VALIDATION_ERROR"
  | "DATABASE_ERROR"
  | "NETWORK_ERROR";

export interface WaitlistSubmissionResult {
  status: WaitlistSubmissionStatus;
  message: string;
  isDuplicate?: boolean;
}

export interface DemoRequestSubmission {
  fullName: string;
  workEmail: string;
  facilityName: string;
  specimenVolume?: string;
  notes?: string;
}

export interface ContactFormSubmission {
  name: string;
  email: string;
  subject?: string;
  message: string;
}

export interface WaitlistAggregateStats {
  totalMembers: number;
  patientsCount: number;
  organizationsCount: number;
  suppliersCount: number;
  loading: boolean;
  error: boolean;
}

/**
 * Production-grade waitlist submission function with email normalization and deduplication.
 */
export async function saveWaitlistToSupabase(entry: WaitlistSubmission): Promise<WaitlistSubmissionResult> {
  if (!entry.fullName || !entry.fullName.trim() || !entry.email || !entry.email.trim()) {
    return {
      status: "VALIDATION_ERROR",
      message: "Please fill in all required fields (Full Name and Email Address).",
    };
  }

  // Normalize email: trim surrounding whitespace and convert to lowercase
  const normalizedEmail = entry.email.trim().toLowerCase();
  const normalizedName = entry.fullName.trim();

  // 1. Pre-insert duplicate check (case-insensitive email search)
  try {
    const { data: existing, error: checkError } = await supabase
      .from("waitlist")
      .select("id")
      .ilike("email", normalizedEmail)
      .limit(1);

    if (checkError) {
      console.warn("Pre-insert duplicate check query warning:", checkError.message);
    }

    if (existing && existing.length > 0) {
      console.info("Duplicate waitlist registration detected for email:", normalizedEmail);
      return {
        status: "DUPLICATE",
        isDuplicate: true,
        message: "You're already on the Curexal early access list. We'll keep you updated on progress.",
      };
    }
  } catch (preCheckErr) {
    console.warn("Pre-insert duplicate check exception:", preCheckErr);
  }

  // Compile customer research answers into desired_features and biggest_pain_point
  const researchSummary: string[] = [
    ...(entry.desiredFeatures || []),
    ...(entry.externalPartners ? [`Partners: ${entry.externalPartners.join(", ")}`] : []),
    ...(entry.communicationChannels ? [`Channels: ${entry.communicationChannels.join(", ")}`] : []),
    ...(entry.trustFactors ? [`Trust Factors: ${entry.trustFactors.join(", ")}`] : []),
    ...(entry.placesInvolved ? [`Facilities Involved: ${entry.placesInvolved}`] : []),
    ...(entry.testPilotInterest ? [`Pilot Interest: ${entry.testPilotInterest}`] : []),
  ];

  const primaryPainPoint = [
    entry.biggestPainPoint,
    entry.referralBreakdownStory ? `Referral Breakdown: ${entry.referralBreakdownStory}` : null,
    entry.selfCoordinatedTasks ? `Self-Coordinated Tasks: ${entry.selfCoordinatedTasks}` : null,
    entry.unresolvedImpact ? `Unresolved Impact: ${entry.unresolvedImpact}` : null,
    entry.automaticCoordinationDesire ? `Desired Outcome: ${entry.automaticCoordinationDesire}` : null,
  ]
    .filter(Boolean)
    .join(" | ");

  const payload = {
    full_name: normalizedName,
    email: normalizedEmail,
    phone: entry.phone ? entry.phone.trim() : null,
    country: entry.country || null,
    state: entry.state || null,
    city: entry.city || null,
    persona: entry.persona || "Patient",
    organization: entry.organization || entry.orgType || null,
    biggest_pain_point: primaryPainPoint || entry.biggestPainPoint || null,
    desired_features: researchSummary,
    timeline: entry.timeline || null,
    shaping_preference: entry.shapingPreference || null,
  };

  try {
    const { error } = await supabase.from("waitlist").insert([payload]);

    if (error) {
      console.error("Supabase Database Insert Error (waitlist):", error);
      const errCode = error.code || "";
      const errMsg = (error.message || "").toLowerCase();

      // Postgres 23505 unique violation check
      if (
        errCode === "23505" ||
        errMsg.includes("unique") ||
        errMsg.includes("duplicate") ||
        errMsg.includes("already exists")
      ) {
        return {
          status: "DUPLICATE",
          isDuplicate: true,
          message: "You're already on the Curexal early access list. We'll keep you updated on progress.",
        };
      }

      // Fallback payload for strict schema databases
      const fallbackPayload = {
        full_name: normalizedName,
        email: normalizedEmail,
        phone: entry.phone ? entry.phone.trim() : null,
        country: entry.country || null,
        state: entry.state || null,
        city: entry.city || null,
        persona: entry.persona || "Patient",
        organization: entry.organization || null,
      };

      const { error: fallbackErr } = await supabase.from("waitlist").insert([fallbackPayload]);

      if (fallbackErr) {
        const fbCode = fallbackErr.code || "";
        const fbMsg = (fallbackErr.message || "").toLowerCase();
        if (
          fbCode === "23505" ||
          fbMsg.includes("unique") ||
          fbMsg.includes("duplicate") ||
          fbMsg.includes("already exists")
        ) {
          return {
            status: "DUPLICATE",
            isDuplicate: true,
            message: "You're already on the Curexal early access list. We'll keep you updated on progress.",
          };
        }

        return {
          status: "DATABASE_ERROR",
          message: "Unable to complete registration due to a temporary database issue. Please try again.",
        };
      }
    }

    return {
      status: "SUCCESS",
      message: "Spot reserved on our priority waitlist!",
    };
  } catch (netErr) {
    console.error("Network / unexpected error saving to Supabase:", netErr);
    return {
      status: "NETWORK_ERROR",
      message: "Unable to reach database services. Please check your internet connection and try again.",
    };
  }
}

/**
 * Fetches live aggregate statistics from public.waitlist without exposing private user rows.
 */
export async function fetchWaitlistStats(): Promise<WaitlistAggregateStats> {
  try {
    const [totalRes, patientsRes, orgsRes, suppliersRes] = await Promise.all([
      supabase.from("waitlist").select("*", { count: "exact", head: true }),
      supabase.from("waitlist").select("*", { count: "exact", head: true }).in("persona", ["Patient", "patient"]),
      supabase.from("waitlist").select("*", { count: "exact", head: true }).in("persona", [
        "Clinic",
        "clinic",
        "Laboratory",
        "laboratory",
        "Hospital",
        "hospital",
        "Pharmacy",
        "pharmacy",
        "Diagnostic Centre",
        "diagnostic centre",
        "Organization",
        "organization",
        "Doctor",
        "doctor",
      ]),
      supabase.from("waitlist").select("*", { count: "exact", head: true }).in("persona", [
        "Medical Supplier",
        "medical supplier",
        "Supplier",
        "supplier",
      ]),
    ]);

    if (totalRes.error) {
      console.warn("Supabase waitlist total count query returned error:", totalRes.error);
    }

    const hasCounts = totalRes.count !== null && totalRes.count !== undefined;

    if (hasCounts) {
      return {
        totalMembers: totalRes.count || 0,
        patientsCount: patientsRes.count || 0,
        organizationsCount: orgsRes.count || 0,
        suppliersCount: suppliersRes.count || 0,
        loading: false,
        error: false,
      };
    }

    // Try backend API fallback if Supabase returns null count
    try {
      const apiRes = await fetch(import.meta.env.VITE_API_URL ? `${import.meta.env.VITE_API_URL}/api/v1/waitlist/stats` : "/api/v1/waitlist/stats");
      if (apiRes.ok) {
        const data = await apiRes.json();
        return {
          totalMembers: data.totalMembers ?? data.total ?? 0,
          patientsCount: data.patientsCount ?? data.patients ?? 0,
          organizationsCount: data.organizationsCount ?? data.organizations ?? 0,
          suppliersCount: data.suppliersCount ?? data.suppliers ?? 0,
          loading: false,
          error: false,
        };
      }
    } catch (apiErr) {
      console.warn("Backend API waitlist stats fallback error:", apiErr);
    }

    return {
      totalMembers: 0,
      patientsCount: 0,
      organizationsCount: 0,
      suppliersCount: 0,
      loading: false,
      error: false,
    };
  } catch (err) {
    console.warn("Could not fetch waitlist aggregate statistics:", err);
    return {
      totalMembers: 0,
      patientsCount: 0,
      organizationsCount: 0,
      suppliersCount: 0,
      loading: false,
      error: true,
    };
  }
}

/**
 * Persists demo requests to Supabase
 */
export async function saveDemoRequestToSupabase(entry: DemoRequestSubmission): Promise<WaitlistSubmissionResult> {
  console.log("Persisting demo request to Supabase:", entry);
  const normalizedEmail = entry.workEmail.trim().toLowerCase();

  const payload = {
    full_name: entry.fullName.trim(),
    email: normalizedEmail,
    organization: entry.facilityName.trim(),
    notes: `Specimen volume: ${entry.specimenVolume || "N/A"} | ${entry.notes || ""}`,
    persona: "Laboratory",
    created_at: new Date().toISOString(),
  };

  const { error } = await supabase.from("demo_requests").insert([payload]).select();
  if (!error) {
    return {
      status: "SUCCESS",
      message: "Spot reserved on our priority demo waitlist!",
    };
  }

  // Fallback to canonical waitlist function
  return saveWaitlistToSupabase({
    fullName: entry.fullName,
    email: entry.workEmail,
    organization: entry.facilityName,
    persona: "Laboratory",
    biggestPainPoint: entry.notes || "Live demo walkthrough request",
    timeline: "Immediately",
  });
}

/**
 * Persists contact form inquiries to Supabase
 */
export async function saveContactFormToSupabase(entry: ContactFormSubmission): Promise<WaitlistSubmissionResult> {
  const normalizedEmail = entry.email.trim().toLowerCase();

  const payload = {
    name: entry.name.trim(),
    email: normalizedEmail,
    subject: entry.subject || "General Inquiry",
    message: entry.message,
    created_at: new Date().toISOString(),
  };

  const { error } = await supabase.from("contact_messages").insert([payload]).select();
  if (!error) {
    return {
      status: "SUCCESS",
      message: "Message sent successfully!",
    };
  }

  return saveWaitlistToSupabase({
    fullName: entry.name,
    email: entry.email,
    persona: "Other",
    biggestPainPoint: `[${entry.subject || "Contact"}] ${entry.message}`,
  });
}

/**
 * Persists newsletter email subscriptions to Supabase
 */
export async function saveNewsletterToSupabase(email: string): Promise<WaitlistSubmissionResult> {
  const normalizedEmail = email.trim().toLowerCase();

  const { error } = await supabase.from("newsletter_subscriptions").insert([{ email: normalizedEmail, created_at: new Date().toISOString() }]).select();
  if (!error) {
    return {
      status: "SUCCESS",
      message: "Subscribed to newsletter updates!",
    };
  }

  return saveWaitlistToSupabase({
    fullName: normalizedEmail.split("@")[0] || "Subscriber",
    email: normalizedEmail,
    persona: "Patient",
    shapingPreference: "Keep me updated via Email",
  });
}
