import { useState } from "react";
import { MarketingNavbar } from "@/components/layouts/marketing-navbar";
import { MarketingFooter } from "@/components/layouts/marketing-footer";
import { SEOHead } from "@/components/seo/seo-head";
import { Link } from "react-router-dom";
import {
  Stethoscope,
  HeartPulse,
  CreditCard,
  FlaskConical,
  Radio,
  Pill,
  Building2,
  Video,
  ArrowRight,
  Sparkles,
  CheckCircle2,
  Clock,
  Layers,
  ShieldCheck,
} from "lucide-react";
import { BusinessGrowth } from "@/components/home/business-growth";

interface ProductSolution {
  id: string;
  icon: any;
  label: string;
  category: "active" | "coming_soon";
  statusBadge: string;
  badgeVariant: "live" | "coming_soon";
  headline: string;
  desc: string;
  features: string[];
  href: string;
  cta: string;
}

const allProducts: ProductSolution[] = [
  // ── 1. ACTIVE CLINIC OS PRODUCTS ─────────────────────────────
  {
    id: "clinic-os",
    icon: Stethoscope,
    label: "Clinic Outpatient OS",
    category: "active",
    statusBadge: "Available in Clinic MVP",
    badgeVariant: "live",
    headline: "Complete Practice Management & Outpatient EMR",
    desc: "Designed for modern healthcare clinics, outpatient centers, and specialist practices. Unifies patient intake, nursing triage, doctor SOAP notes, ICD-10 coding, and cashier billing into one secure operational workflow.",
    features: [
      "Walk-in patient registration & Master Patient Index (MPI)",
      "Nursing triage desk with vital signs & acuity scoring",
      "Electronic SOAP clinical consultation room canvas",
      "Standard ICD-10 diagnostic coding & e-prescriptions",
      "Point-of-Sale (POS) invoicing & multi-tender receipts",
      "Role-based access control (Doctor, Nurse, Receptionist, Cashier, Admin)",
    ],
    href: "/book-demo",
    cta: "Book Clinic OS Demo",
  },
  {
    id: "patient-portal",
    icon: HeartPulse,
    label: "Patient Health Portal",
    category: "active",
    statusBadge: "Available in Clinic MVP",
    badgeVariant: "live",
    headline: "Unified Care Continuum for Every Patient",
    desc: "Give patients instant self-service access to verified consultation summaries, digital e-prescriptions, cashier receipts, and upcoming appointment schedules in a zero-leakage encrypted vault.",
    features: [
      "Real-time care journey milestone tracker",
      "Digital prescription downloads and dosage guidance",
      "Point-of-Sale billing history & downloadable receipts",
      "Doctor appointment scheduling and check-in confirmation",
      "End-to-end encrypted health records vault",
      "Accessible from any mobile browser or desktop device",
    ],
    href: "/book-demo",
    cta: "Request Portal Demo",
  },
  {
    id: "billing-pos",
    icon: CreditCard,
    label: "Cashier POS & Billing Engine",
    category: "active",
    statusBadge: "Available in Clinic MVP",
    badgeVariant: "live",
    headline: "Real-Time Healthcare Invoicing & Multi-Tender POS",
    desc: "Eliminate revenue leakage with automated fee schedules, split payment processing (Cash, Card, Bank Transfer), and instant receipt generation.",
    features: [
      "Clinical service fee catalog & custom price overrides",
      "Multi-tender cashier register with split payment capabilities",
      "Automated receipt numbering & print templates",
      "Daily revenue reconciliation & cashier shift settlement",
      "Patient invoice history and outstanding balance ledger",
      "Audit-ready financial transaction log",
    ],
    href: "/book-demo",
    cta: "Request Billing Demo",
  },

  // ── 2. FUTURE VERTICAL PRODUCTS (TAGGED COMING SOON) ─────────
  {
    id: "lims",
    icon: FlaskConical,
    label: "Medical Laboratory LIS",
    category: "active",
    statusBadge: "Facility Operating Suite",
    badgeVariant: "live",
    headline: "Specimen Accessioning, Barcoding & Automated Analyzer Telemetry",
    desc: "Automate phlebotomy intake, specimen barcode label generation, bidirectional analyzer interfacing (Sysmex, Mindray, Roche), and two-step pathologist authorization. Built for ISO 15189 pathology environments.",
    features: [
      "Specimen accessioning with unique barcode tracking",
      "Automated laboratory analyzer integration worklists",
      "Two-step consultant pathologist review & digital signature",
      "Internal Quality Control (IQC) & Levey-Jennings charts",
      "Turnaround time (TAT) monitoring & breach alerts",
      "Direct PDF report dispatch via SMS and WhatsApp",
    ],
    href: "/book-demo",
    cta: "Book Lab LIS Demo",
  },
  {
    id: "ris-pacs",
    icon: Radio,
    label: "Radiology RIS & PACS",
    category: "active",
    statusBadge: "Facility Operating Suite",
    badgeVariant: "live",
    headline: "DICOM Imaging Modalities, Cloud PACS & Structured Reports",
    desc: "Connect Ultrasound, X-Ray, CT, and MRI modalities into a centralized PACS archive with browser-based DICOM viewing and structured radiologist reporting.",
    features: [
      "DICOM Modality Worklist (MWL) integration & queueing",
      "Zero-footprint web DICOM viewer accessible from any device",
      "Structured radiological reporting with signature stamps",
      "Zero-film digital report links dispatched to referrers",
      "Cloud-backed medical imaging archive and PACS backup",
      "Prior scan comparison and historical study timeline",
    ],
    href: "/book-demo",
    cta: "Book Radiology RIS Demo",
  },
  {
    id: "pharmacy",
    icon: Pill,
    label: "Pharmacy & Dispensary",
    category: "coming_soon",
    statusBadge: "Coming Soon",
    badgeVariant: "coming_soon",
    headline: "FEFO Drug Inventory, Batch Tracking & Dispensing",
    desc: "Seamlessly transition electronic prescriptions from clinic consultation rooms to the dispensary with automated stock depletion, batch/expiry alerts, and supplier reordering.",
    features: [
      "First-Expiry-First-Out (FEFO) batch and lot allocation",
      "Direct e-prescription reception from clinic consult rooms",
      "Drug-drug interaction and allergy safety warnings",
      "Automated low-stock threshold and reorder point alerts",
      "Controlled substance tracking and audit log",
      "Supplier purchase order and goods receipt management",
    ],
    href: "/waitlist",
    cta: "Join Pharmacy Waitlist",
  },
  {
    id: "hospital-his",
    icon: Building2,
    label: "Hospital Inpatient & Wards (HIS)",
    category: "coming_soon",
    statusBadge: "Coming Soon",
    badgeVariant: "coming_soon",
    headline: "Ward Bed Board, Admissions & Inpatient Care",
    desc: "Comprehensive inpatient hospital management featuring interactive ward bed boards, admission/transfer/discharge (ADT) workflows, nurse shift handovers, and surgery theater scheduling.",
    features: [
      "Interactive visual ward bed map and occupancy grid",
      "Admission, transfer, and discharge (ADT) workflows",
      "Nurse medication administration records (eMAR)",
      "Shift handover notes and doctor round chart updates",
      "Operating room (OR) and surgical theater scheduling",
      "Comprehensive inpatient billing and bed charges accumulation",
    ],
    href: "/waitlist",
    cta: "Join Hospital HIS Waitlist",
  },
  {
    id: "telehealth",
    icon: Video,
    label: "Telehealth & Virtual Care",
    category: "coming_soon",
    statusBadge: "Coming Soon",
    badgeVariant: "coming_soon",
    headline: "Remote Doctor Consultations & Digital Triage",
    desc: "High-definition, low-bandwidth video and audio consultations with real-time in-call clinical note taking, screen sharing, and instant e-prescription dispatch.",
    features: [
      "Secure peer-to-peer WebRTC video and voice consultations",
      "In-consultation real-time SOAP note documentation",
      "Patient waiting room and automated queue notifications",
      "Low-bandwidth adaptive streaming for 3G/4G connectivity",
      "Instant digital prescription and receipt sharing",
      "Session recording with patient consent and audit trail",
    ],
    href: "/waitlist",
    cta: "Join Telehealth Waitlist",
  },
];

export function SolutionsPage() {
  const [activeFilter, setActiveFilter] = useState<"all" | "active" | "coming_soon">("all");

  const filteredProducts = allProducts.filter((product) => {
    if (activeFilter === "all") return true;
    return product.category === activeFilter;
  });

  return (
    <div className="min-h-screen bg-white dark:bg-[#0B1120] font-inter">
      <SEOHead
        title="Healthcare Solutions: Clinic OS, EMR, LIMS, RIS & Hospital Products"
        description="Explore Curexal's modular healthcare operating system suite, from active Clinic OS and EMR to upcoming LIMS, Radiology RIS, and Hospital solutions."
      />
      <MarketingNavbar />

      {/* Hero Header */}
      <div className="pt-28 pb-16 bg-[#F8FAFC] dark:bg-[#0B1120] border-b border-gray-100 dark:border-[#1F2937]">
        <div className="max-w-[1280px] mx-auto px-6">
          <div className="max-w-3xl">
            <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-teal-50 dark:bg-teal-950/60 border border-teal-200/80 dark:border-teal-800/80 text-[#0F766E] dark:text-teal-300 text-xs font-semibold uppercase tracking-wider mb-4">
              <Sparkles className="w-3.5 h-3.5" />
              <span>Unified Healthcare Product Suite</span>
            </div>
            <h1 className="text-3xl sm:text-4xl lg:text-5xl font-extrabold tracking-tight text-gray-900 dark:text-white mb-5 leading-tight">
              Enterprise Solutions for Every Healthcare Sector.
            </h1>
            <p className="text-base sm:text-lg text-gray-500 dark:text-gray-400 leading-relaxed">
              Curexal provides a modular, multi-tenant operating system kernel. Explore our active <strong>Clinic OS MVP</strong> suite alongside specialized upcoming vertical platforms engineered for diagnostic and hospital networks.
            </p>

            {/* Filter Tabs */}
            <div className="flex flex-wrap items-center gap-2.5 mt-8">
              <button
                onClick={() => setActiveFilter("all")}
                className={`px-4 py-2 rounded-xl text-xs font-semibold transition-all cursor-pointer ${
                  activeFilter === "all"
                    ? "bg-[#0F766E] text-white shadow-sm"
                    : "bg-white dark:bg-gray-800 text-gray-600 dark:text-gray-300 border border-gray-200 dark:border-gray-700 hover:border-gray-300"
                }`}
              >
                All Products ({allProducts.length})
              </button>
              <button
                onClick={() => setActiveFilter("active")}
                className={`px-4 py-2 rounded-xl text-xs font-semibold transition-all cursor-pointer flex items-center gap-1.5 ${
                  activeFilter === "active"
                    ? "bg-[#0F766E] text-white shadow-sm"
                    : "bg-white dark:bg-gray-800 text-gray-600 dark:text-gray-300 border border-gray-200 dark:border-gray-700 hover:border-gray-300"
                }`}
              >
                <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
                Live in Clinic OS ({allProducts.filter((p) => p.category === "active").length})
              </button>
              <button
                onClick={() => setActiveFilter("coming_soon")}
                className={`px-4 py-2 rounded-xl text-xs font-semibold transition-all cursor-pointer flex items-center gap-1.5 ${
                  activeFilter === "coming_soon"
                    ? "bg-[#0F766E] text-white shadow-sm"
                    : "bg-white dark:bg-gray-800 text-gray-600 dark:text-gray-300 border border-gray-200 dark:border-gray-700 hover:border-gray-300"
                }`}
              >
                <Clock className="w-3.5 h-3.5 text-amber-500" />
                Coming Soon ({allProducts.filter((p) => p.category === "coming_soon").length})
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Solutions Product Grid */}
      <main className="max-w-[1280px] mx-auto px-6 py-16 space-y-16">
        {filteredProducts.map((sol, idx) => {
          const Icon = sol.icon;
          const isReversed = idx % 2 !== 0;
          const isLive = sol.badgeVariant === "live";

          return (
            <div
              key={sol.id}
              id={sol.id}
              className={`grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-20 items-start p-8 sm:p-10 rounded-3xl border transition-all ${
                isLive
                  ? "bg-white dark:bg-[#0F172A]/50 border-teal-500/20 shadow-lg shadow-teal-500/5 hover:border-teal-500/40"
                  : "bg-[#FAFAFA] dark:bg-[#0B1120] border-gray-200 dark:border-gray-800 hover:border-amber-500/30"
              }`}
            >
              <div>
                <div className="flex flex-wrap items-center gap-2.5 mb-4">
                  <div
                    className={`w-9 h-9 rounded-xl flex items-center justify-center ${
                      isLive
                        ? "bg-teal-500/10 text-teal-600 dark:text-teal-400 border border-teal-500/20"
                        : "bg-amber-500/10 text-amber-600 dark:text-amber-400 border border-amber-500/20"
                    }`}
                  >
                    <Icon className="h-5 w-5" strokeWidth={2} />
                  </div>
                  <span className="text-xs font-bold text-gray-700 dark:text-gray-300 uppercase tracking-wide">
                    {sol.label}
                  </span>

                  {/* Status Badge */}
                  {isLive ? (
                    <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-emerald-500/10 border border-emerald-500/30 text-emerald-600 dark:text-emerald-400 text-[11px] font-semibold">
                      <span className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
                      {sol.statusBadge}
                    </span>
                  ) : (
                    <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-amber-500/10 border border-amber-500/30 text-amber-600 dark:text-amber-400 text-[11px] font-semibold">
                      <Clock className="w-3 h-3" />
                      {sol.statusBadge}
                    </span>
                  )}
                </div>

                <h2 className="text-2xl sm:text-3xl font-extrabold text-gray-900 dark:text-white mb-4 tracking-tight">
                  {sol.headline}
                </h2>
                <p className="text-sm sm:text-base text-gray-600 dark:text-gray-400 mb-8 leading-relaxed">
                  {sol.desc}
                </p>

                <div className="flex items-center gap-4">
                  <Link to={sol.href}>
                    <button
                      className={`flex items-center gap-2 px-6 py-3 rounded-xl text-sm font-semibold transition-all cursor-pointer border-0 shadow-md ${
                        isLive
                          ? "bg-[#0F766E] hover:bg-[#115E59] text-white hover:scale-[1.02]"
                          : "bg-amber-600 hover:bg-amber-700 text-white hover:scale-[1.02]"
                      }`}
                    >
                      {sol.cta}
                      <ArrowRight className="h-4 w-4" />
                    </button>
                  </Link>

                  {isLive && (
                    <span className="text-xs text-muted-foreground hidden sm:inline-flex items-center gap-1">
                      <ShieldCheck className="w-4 h-4 text-emerald-500" /> Production Ready
                    </span>
                  )}
                </div>
              </div>

              {/* Key Capabilities Card */}
              <div className="bg-white dark:bg-[#131C31] p-6 sm:p-8 rounded-2xl border border-gray-200/80 dark:border-gray-800 shadow-sm">
                <div className="flex items-center justify-between pb-4 border-b border-gray-100 dark:border-gray-800 mb-5">
                  <h3 className="text-xs font-bold text-gray-500 dark:text-gray-400 uppercase tracking-widest flex items-center gap-2">
                    <Layers className="w-3.5 h-3.5" />
                    Core Platform Capabilities
                  </h3>
                  {isLive ? (
                    <span className="text-[10px] font-mono font-bold text-emerald-600 dark:text-emerald-400 uppercase">
                      Live MVP Feature
                    </span>
                  ) : (
                    <span className="text-[10px] font-mono font-bold text-amber-600 dark:text-amber-400 uppercase">
                      Roadmap Release
                    </span>
                  )}
                </div>

                <ul className="space-y-3.5">
                  {sol.features.map((f) => (
                    <li key={f} className="flex items-start gap-3 text-sm text-gray-700 dark:text-gray-300">
                      <CheckCircle2
                        className={`h-4 w-4 mt-0.5 flex-shrink-0 ${
                          isLive ? "text-[#0F766E] dark:text-teal-400" : "text-amber-500"
                        }`}
                      />
                      <span>{f}</span>
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
