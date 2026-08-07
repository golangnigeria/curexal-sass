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

/**
 * Persists 6-step customer research waitlist registrations to Supabase 'waitlist' table
 */
export async function saveWaitlistToSupabase(entry: WaitlistSubmission) {
  console.log("Persisting waitlist registration to Supabase database:", entry);

  const payload = {
    full_name: entry.fullName,
    email: entry.email,
    phone: entry.phone || null,
    country: entry.country || null,
    state: entry.state || null,
    city: entry.city || null,
    persona: entry.persona || "Patient",
    organization: entry.organization || null,
    biggest_pain_point: entry.biggestPainPoint || null,
    desired_features: entry.desiredFeatures || [],
    timeline: entry.timeline || null,
    shaping_preference: entry.shapingPreference || null,
  };

  const { data, error } = await supabase.from("waitlist").insert([payload]).select();

  if (error) {
    console.error("Supabase Database Insert Error (waitlist):", error);
    const fallbackPayload = {
      full_name: entry.fullName,
      email: entry.email,
      phone: entry.phone || null,
      country: entry.country || null,
      state: entry.state || null,
      city: entry.city || null,
      persona: entry.persona || "Patient",
      organization: entry.organization || null,
    };
    const fallbackRes = await supabase.from("waitlist").insert([fallbackPayload]);
    if (!fallbackRes.error) {
      console.log("Successfully persisted entry using fallback payload to Supabase 'waitlist':", fallbackRes.data);
      return true;
    }
    return false;
  }

  console.log("Successfully persisted waitlist entry to Supabase 'waitlist' table:", data);
  return true;
}

/**
 * Persists demo requests to Supabase 'demo_requests' table (or 'waitlist' fallback)
 */
export async function saveDemoRequestToSupabase(entry: DemoRequestSubmission) {
  console.log("Persisting demo request to Supabase:", entry);

  const payload = {
    full_name: entry.fullName,
    email: entry.workEmail,
    organization: entry.facilityName,
    notes: `Specimen volume: ${entry.specimenVolume || "N/A"} | ${entry.notes || ""}`,
    persona: "Laboratory",
    created_at: new Date().toISOString(),
  };

  const { data, error } = await supabase.from("demo_requests").insert([payload]).select();
  if (!error) {
    console.log("Persisted to 'demo_requests' table:", data);
    return true;
  }

  // Fallback to waitlist table
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
 * Persists contact form inquiries to Supabase 'contact_messages' table (or 'waitlist' fallback)
 */
export async function saveContactFormToSupabase(entry: ContactFormSubmission) {
  console.log("Persisting contact inquiry to Supabase:", entry);

  const payload = {
    name: entry.name,
    email: entry.email,
    subject: entry.subject || "General Inquiry",
    message: entry.message,
    created_at: new Date().toISOString(),
  };

  const { data, error } = await supabase.from("contact_messages").insert([payload]).select();
  if (!error) {
    console.log("Persisted to 'contact_messages' table:", data);
    return true;
  }

  return saveWaitlistToSupabase({
    fullName: entry.name,
    email: entry.email,
    persona: "Other",
    biggestPainPoint: `[${entry.subject || "Contact"}] ${entry.message}`,
  });
}

/**
 * Persists newsletter email subscriptions to Supabase 'newsletter_subscriptions' table (or 'waitlist' fallback)
 */
export async function saveNewsletterToSupabase(email: string) {
  console.log("Persisting newsletter email to Supabase:", email);

  const { data, error } = await supabase.from("newsletter_subscriptions").insert([{ email, created_at: new Date().toISOString() }]).select();
  if (!error) {
    console.log("Persisted to 'newsletter_subscriptions' table:", data);
    return true;
  }

  return saveWaitlistToSupabase({
    fullName: email.split("@")[0] || "Subscriber",
    email,
    persona: "Patient",
    shapingPreference: "Keep me updated via Email",
  });
}
