import { MarketingNavbar } from "@/components/layouts/marketing-navbar";
import { MarketingFooter } from "@/components/layouts/marketing-footer";
import { SEOHead } from "@/components/seo/seo-head";
import { Link } from "react-router-dom";
import { ArrowRight, ShieldCheck, Network, Lock, Sparkles, Building2, CheckCircle2 } from "lucide-react";

export function AboutPage() {
  return (
    <div className="min-h-screen bg-white dark:bg-[#0B1120] font-inter">
      <SEOHead
        title="About Curexal: Connected Healthcare Operating System"
        description="Learn why Curexal is building the connected digital operating network for African healthcare providers and patients."
      />
      <MarketingNavbar />

      {/* Header */}
      <div className="pt-28 pb-16 bg-slate-50 dark:bg-[#0B1120] border-b border-slate-200/60 dark:border-slate-800">
        <div className="max-w-[1280px] mx-auto px-6">
          <div className="max-w-3xl">
            <div className="inline-flex items-center gap-1.5 px-3 py-1 mb-4 rounded-full border border-teal-500/30 bg-teal-50 dark:bg-teal-950/40 text-[#0F766E] dark:text-teal-400 text-xs font-bold uppercase tracking-wider">
              Product Philosophy & Mission
            </div>
            <h1 className="text-3xl sm:text-5xl font-black text-slate-900 dark:text-white tracking-tight mb-5">
              We believe healthcare is too fragmented.
            </h1>
            <p className="text-base sm:text-lg text-slate-600 dark:text-slate-300 leading-relaxed">
              A patient can visit one organization for consultation, another for diagnostics, another for imaging, and another for medication. Yet these healthcare organizations often operate using completely disconnected software systems.
            </p>
          </div>
        </div>
      </div>

      <main className="max-w-[1280px] mx-auto px-6 py-16 space-y-20">

        {/* Mission Statement Section */}
        <div className="p-8 sm:p-12 rounded-3xl bg-gradient-to-br from-[#0F766E] to-[#115E59] text-white shadow-xl space-y-4">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-white/10 text-teal-200 text-xs font-bold uppercase tracking-wider">
            <Sparkles className="w-3.5 h-3.5" />
            <span>Our Mission</span>
          </div>
          <blockquote className="text-xl sm:text-3xl font-extrabold leading-tight tracking-tight text-white max-w-4xl">
            "To build the digital backbone of healthcare delivery in Africa by connecting providers, patients, and partners on a single operating network."
          </blockquote>
        </div>

        {/* What Curexal Is: Healthcare OS */}
        <div className="space-y-8">
          <div className="max-w-3xl">
            <h2 className="text-2xl sm:text-4xl font-extrabold text-slate-900 dark:text-white tracking-tight mb-3">
              Curexal is a Healthcare Operating System.
            </h2>
            <p className="text-sm sm:text-base text-slate-600 dark:text-slate-300 leading-relaxed">
              Traditional healthcare software often focuses strictly on the boundary of one organization. Curexal is designed to connect multiple independent healthcare organizations while preserving their data boundaries and operational independence.
            </p>
          </div>

          {/* Model Comparison Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div className="p-6 sm:p-8 rounded-3xl bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 space-y-4">
              <h3 className="text-sm font-bold uppercase tracking-wider text-rose-500">Traditional Software Model</h3>
              <div className="space-y-2.5 text-xs text-slate-700 dark:text-slate-300">
                <p className="p-2.5 rounded-xl bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800">Clinic → Isolated System</p>
                <p className="p-2.5 rounded-xl bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800">Laboratory → Isolated System</p>
                <p className="p-2.5 rounded-xl bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800">Pharmacy → Isolated System</p>
                <p className="p-2.5 rounded-xl bg-[#FFF1F2] dark:bg-rose-950/40 text-rose-700 dark:text-rose-300 font-semibold border border-rose-200 dark:border-rose-900">Patient → Disconnected Experience</p>
              </div>
            </div>

            <div className="p-6 sm:p-8 rounded-3xl bg-teal-50/60 dark:bg-teal-950/30 border border-teal-200 dark:border-teal-800 space-y-4">
              <h3 className="text-sm font-bold uppercase tracking-wider text-[#0F766E] dark:text-teal-400">The Curexal Network Model</h3>
              <div className="space-y-2.5 text-xs text-slate-800 dark:text-slate-200 font-medium">
                <p className="p-2.5 rounded-xl bg-white dark:bg-slate-900 border border-teal-200 dark:border-teal-800">Clinic ──┐</p>
                <p className="p-2.5 rounded-xl bg-white dark:bg-slate-900 border border-teal-200 dark:border-teal-800">Laboratory ──┼──→ CUREXAL NETWORK</p>
                <p className="p-2.5 rounded-xl bg-white dark:bg-slate-900 border border-teal-200 dark:border-teal-800">Pharmacy ──┤</p>
                <p className="p-2.5 rounded-xl bg-[#F0FDFA] dark:bg-teal-950 text-[#0F766E] dark:text-teal-300 font-bold border border-teal-300 dark:border-teal-700">Patient ───┘ (Coordinated Care Journey)</p>
              </div>
            </div>
          </div>
        </div>

        {/* Collaboration By Default */}
        <div className="space-y-4 max-w-3xl">
          <h2 className="text-2xl sm:text-3xl font-extrabold text-slate-900 dark:text-white tracking-tight">
            Healthcare collaboration should be the default.
          </h2>
          <p className="text-sm sm:text-base text-slate-600 dark:text-slate-300 leading-relaxed">
            When a doctor creates a diagnostic referral, the receiving laboratory should know immediately. When a pathologist authorizes a test report, the referring doctor should receive the result without making the patient act as a courier. Cross-organizational coordination is built into our core architecture.
          </p>
        </div>

        {/* Tenant Privacy & Data Isolation */}
        <div className="p-6 sm:p-8 rounded-3xl bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 space-y-4">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-teal-50 dark:bg-teal-950 text-[#0F766E] flex items-center justify-center">
              <Lock className="w-5 h-5" />
            </div>
            <div>
              <h3 className="text-lg font-bold text-slate-900 dark:text-white">Tenant Privacy & Data Isolation Architecture</h3>
              <p className="text-xs text-slate-500 dark:text-slate-400">Cross-organizational coordination with strict tenant boundaries</p>
            </div>
          </div>

          <p className="text-xs sm:text-sm text-slate-600 dark:text-slate-300 leading-relaxed">
            Curexal is built on a schema-per-tenant architecture. Every healthcare organization operates within its own dedicated database schema, ensuring complete data isolation. Data exchange between facilities occurs strictly through encrypted, audit-logged cross-tenant APIs only when authorized by clinical workflow context.
          </p>
        </div>

        {/* Contact Footer */}
        <div id="contact" className="grid grid-cols-1 md:grid-cols-3 gap-5">
          {[
            { label: "General Enquiries", email: "hello@curexal.com" },
            { label: "Early Access & Pilots", email: "pilots@curexal.com" },
            { label: "Engineering & Architecture", email: "tech@curexal.com" },
          ].map((c) => (
            <div key={c.label} className="p-6 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800">
              <p className="text-xs font-semibold text-slate-400 dark:text-slate-500 uppercase tracking-widest mb-2">{c.label}</p>
              <a
                href={`mailto:${c.email}`}
                className="text-sm font-bold text-[#0F766E] dark:text-teal-400 hover:underline"
              >
                {c.email}
              </a>
            </div>
          ))}
        </div>

      </main>

      <MarketingFooter />
    </div>
  );
}

