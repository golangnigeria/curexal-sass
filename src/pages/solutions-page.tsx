import { MarketingNavbar } from "@/components/layouts/marketing-navbar";
import { MarketingFooter } from "@/components/layouts/marketing-footer";
import { SEOHead } from "@/components/seo/seo-head";
import { Link } from "react-router-dom";
import { FlaskConical, Stethoscope, User, ShoppingBag, ArrowRight } from "lucide-react";
import { BusinessGrowth } from "@/components/home/business-growth";

const solutions = [
  {
    id: "lims",
    icon: FlaskConical,
    label: "Laboratory LIMS",
    headline: "Laboratory Requisition and Results Workflow",
    desc: "Automate specimen accessioning, barcode label generation, analyzer auto-result ingestion, and multi-tier pathologist validation. Built for ISO 15189 pathology environments.",
    features: [
      "Specimen chain of custody tracking",
      "Auto-interface with Sysmex, Mindray, Roche analyzers",
      "Pathologist digital stamp and e-signature",
      "Turnaround time analytics and TAT breach alerts",
      "Accession number generation and label printing",
      "Quality control charts (Levey-Jennings)",
    ],
    href: "/book-demo",
    cta: "Request LIMS Demo",
  },
  {
    id: "emr",
    icon: Stethoscope,
    label: "Clinic EMR",
    headline: "Electronic Medical Records for Outpatient Clinics",
    desc: "Enable physicians to place electronic lab orders, view results directly in the clinical chart, and manage patient visits, all connected to the laboratory workflow.",
    features: [
      "Electronic lab test ordering",
      "Direct result dispatch to clinical chart",
      "Patient visit notes and consultation records",
      "Prescription management",
      "Vital signs tracking",
      "Referral letter generation",
    ],
    href: "/book-demo",
    cta: "Request EMR Demo",
  },
  {
    id: "patient",
    icon: User,
    label: "Patient Portal",
    headline: "Secure Health Records for Every Patient",
    desc: "Give patients a private, verified health vault where they can view their complete test history, download electronically signed PDF lab reports, and track upcoming appointments.",
    features: [
      "Complete diagnostic test history",
      "Electronically signed PDF report downloads",
      "Appointment scheduling and tracking",
      "Home sample phlebotomy booking",
      "Shareable health summaries for referrals",
      "Trend analytics across test results over time",
    ],
    href: "/marketplace",
    cta: "Explore Patient Marketplace",
  },
  {
    id: "marketplace",
    icon: ShoppingBag,
    label: "Diagnostic Marketplace",
    headline: "A Public Network for Patients to Find Accredited Labs",
    desc: "Patients can search by test name, compare prices across accredited laboratories, filter by location and turnaround time, and book diagnostic appointments, all from one interface.",
    features: [
      "Search by test name, lab name, or location",
      "Price comparison across partner laboratories",
      "Home sample phlebotomy scheduling",
      "Turnaround time visibility",
      "Verified accreditation badges",
      "Instant digital report delivery post-testing",
    ],
    href: "/marketplace",
    cta: "Open Marketplace",
  },
];

export function SolutionsPage() {
  return (
    <div className="min-h-screen bg-white dark:bg-[#0B1120] font-inter">
      <SEOHead
        title="Healthcare Solutions: LIMS, EMR and Marketplace Network"
        description="Software capabilities for diagnostic laboratories, outpatient clinics, and public healthcare networks."
      />
      <MarketingNavbar />

      {/* Header */}
      <div className="pt-28 pb-16 bg-[#F8FAFC] dark:bg-[#0B1120] border-b border-gray-100 dark:border-[#1F2937]">
        <div className="max-w-[1280px] mx-auto px-6">
          <div className="max-w-2xl">
            <div className="accent-line mb-4" />
            <h1 className="text-hero text-gray-900 dark:text-white mb-5">
              Solutions for Every<br />Healthcare Stakeholder.
            </h1>
            <p className="text-body text-gray-500 dark:text-gray-400">
              From ISO-accredited pathology laboratories to outpatient clinics and individual patients, Curexal provides a module for every role in the diagnostic care continuum.
            </p>
          </div>
        </div>
      </div>

      {/* Solutions */}
      <main className="max-w-[1280px] mx-auto px-6 py-16 space-y-16">
        {solutions.map((sol, idx) => {
          const Icon = sol.icon;
          const isReversed = idx % 2 !== 0;
          return (
            <div
              key={sol.id}
              id={sol.id}
              className={`grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-20 items-start ${isReversed ? "lg:direction-rtl" : ""}`}
            >
              <div>
                <div className="inline-flex items-center gap-2 px-3 py-1.5 mb-5 rounded-full bg-[#F0FDFA] dark:bg-[#0F766E]/10 border border-[#0F766E]/20">
                  <Icon className="h-4 w-4 text-[#0F766E]" strokeWidth={1.75} />
                  <span className="text-xs font-semibold text-[#0F766E]">{sol.label}</span>
                </div>
                <h2 className="text-section text-gray-900 dark:text-white mb-4">{sol.headline}</h2>
                <p className="text-body text-gray-500 dark:text-gray-400 mb-8">{sol.desc}</p>
                <Link to={sol.href}>
                  <button
                    className="flex items-center gap-2 px-5 py-2.5 rounded-[10px] bg-[#0F766E] hover:bg-[#115E59] text-white text-sm font-semibold transition-colors cursor-pointer border-0"
                    style={{ boxShadow: "0 1px 2px rgba(0,0,0,0.1)" }}
                  >
                    {sol.cta}
                    <ArrowRight className="h-4 w-4" />
                  </button>
                </Link>
              </div>

              <div className="card-enterprise p-7">
                <h3 className="text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-widest mb-5">
                  Key Capabilities
                </h3>
                <ul className="space-y-3">
                  {sol.features.map((f) => (
                    <li key={f} className="flex items-start gap-3 text-sm text-gray-700 dark:text-gray-300">
                      <span className="mt-1.5 w-1.5 h-1.5 rounded-full bg-[#0F766E] flex-shrink-0" />
                      {f}
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          );
        })}
      </main>

      <BusinessGrowth />

      <MarketingFooter />
    </div>
  );
}
