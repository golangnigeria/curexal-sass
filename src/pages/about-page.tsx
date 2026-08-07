import { MarketingNavbar } from "@/components/layouts/marketing-navbar";
import { MarketingFooter } from "@/components/layouts/marketing-footer";
import { SEOHead } from "@/components/seo/seo-head";
import { Link } from "react-router-dom";
import { ArrowRight } from "lucide-react";

export function AboutPage() {
  return (
    <div className="min-h-screen bg-white dark:bg-[#0B1120] font-inter">
      <SEOHead
        title="About Curexal — Connected Healthcare Infrastructure"
        description="Learn about Curexal's mission to build connected digital infrastructure for laboratories, clinics, and healthcare providers."
      />
      <MarketingNavbar />

      {/* Header */}
      <div className="pt-28 pb-16 bg-[#F8FAFC] dark:bg-[#0B1120] border-b border-gray-100 dark:border-[#1F2937]">
        <div className="max-w-[1280px] mx-auto px-6">
          <div className="max-w-2xl">
            <div className="accent-line mb-4" />
            <h1 className="text-hero text-gray-900 dark:text-white mb-5">
              Building the Operating Network<br />for African Healthcare.
            </h1>
            <p className="text-body text-gray-500 dark:text-gray-400">
              Curexal was built by a team frustrated with disconnected healthcare systems, paper referrals, and manual coordination across Africa. We set out to build the infrastructure that enables independent healthcare providers to collaborate as one connected system, without sacrificing their independence or ownership of data.
            </p>
          </div>
        </div>
      </div>

      <main className="max-w-[1280px] mx-auto px-6 py-16 space-y-20">

        {/* Mission */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-start">
          <div>
            <h2 className="text-section text-gray-900 dark:text-white mb-4">Our Mission</h2>
            <p className="text-body text-gray-500 dark:text-gray-400 mb-4">
              We believe every patient in Africa deserves access to accurate, timely, and affordable diagnostic care, and that every healthcare provider deserves infrastructure that actually works for how they operate.
            </p>
            <p className="text-body text-gray-500 dark:text-gray-400 mb-4">
              Most healthcare software is built for countries with reliable electricity and always-on internet. Curexal is intentionally engineered for unstable connectivity, power interruptions, WhatsApp-first communication, local payment rails, and African healthcare workflows. That isn't localization, it's product strategy.
            </p>
            <p className="text-body text-gray-500 dark:text-gray-400">
              Our competitors sell software. Curexal sells a new operating model for healthcare, where thousands of independent providers function like one coordinated healthcare network.
            </p>
          </div>
          <div className="grid grid-cols-2 gap-4">
            {[
              { label: "Founded", value: "2026" },
              { label: "Data Sovereignty", value: "100%" },
              { label: "ISO 15189", value: "Ready" },
              { label: "Uptime SLA", value: "99.9%" },
            ].map((s) => (
              <div key={s.label} className="card-enterprise p-6">
                <span className="text-[28px] font-bold text-gray-900 dark:text-white">{s.value}</span>
                <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">{s.label}</p>
              </div>
            ))}
          </div>
        </div>

        {/* Vision in one sentence */}
        <div className="card-enterprise p-8 border-[#0F766E]/20 bg-[#F0FDFA]/50 dark:bg-[#0F766E]/5">
          <h2 className="text-section text-gray-900 dark:text-white mb-4">The Vision</h2>
          <p className="text-[17px] leading-relaxed text-gray-700 dark:text-gray-300 max-w-3xl">
            Curexal is building the operating network that enables independent healthcare providers across Africa to collaborate as one connected healthcare system, without sacrificing their independence or ownership of data.
          </p>
        </div>

        {/* Careers */}
        <div id="careers" className="card-enterprise p-8">
          <h2 className="text-section text-gray-900 dark:text-white mb-4">Careers</h2>
          <p className="text-body text-gray-500 dark:text-gray-400 mb-6 max-w-2xl">
            We're building a small, focused team solving hard problems in African healthcare infrastructure. We hire for clarity of thinking, ownership mindset, and genuine care for the mission.
          </p>
          <a href="mailto:careers@curexal.com">
            <button className="flex items-center gap-2 px-5 py-2.5 rounded-[10px] bg-[#0F766E] hover:bg-[#115E59] text-white text-sm font-semibold transition-colors cursor-pointer border-0">
              View Open Roles
              <ArrowRight className="h-4 w-4" />
            </button>
          </a>
        </div>

        {/* Contact */}
        <div id="contact" className="grid grid-cols-1 md:grid-cols-3 gap-5">
          {[
            { label: "General Enquiries", email: "hello@curexal.com" },
            { label: "Sales & Partnerships", email: "sales@curexal.com" },
            { label: "Technical Support", email: "support@curexal.com" },
          ].map((c) => (
            <div key={c.label} className="card-enterprise p-6">
              <p className="text-xs font-semibold text-gray-400 dark:text-gray-500 uppercase tracking-widest mb-3">{c.label}</p>
              <a
                href={`mailto:${c.email}`}
                className="text-sm font-medium text-[#0F766E] hover:text-[#115E59] transition-colors"
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
