import { useState } from "react";
import { MarketingNavbar } from "@/components/layouts/marketing-navbar";
import { MarketingFooter } from "@/components/layouts/marketing-footer";
import { SEOHead } from "@/components/seo/seo-head";
import { Link } from "react-router-dom";
import {
  Search,
  MapPin,
  Clock,
  Star,
  Building2,
  Stethoscope,
  FlaskConical,
  Radio,
  Pill,
  ShieldCheck,
  CheckCircle2,
  ArrowRight,
  Sparkles,
  Layers,
  Filter,
  UserPlus,
  Send,
  Network,
  Share2,
  Activity,
  FileCheck,
} from "lucide-react";
import { WaitlistModal } from "@/components/waitlist-modal";

interface AccreditedPartnerFacility {
  id: string;
  name: string;
  category: "clinic" | "laboratory" | "radiology";
  categoryLabel: string;
  location: string;
  accreditation: string;
  telemetryCapabilities: string[];
  sampleProcedures: { name: string; tat: string; routingType: string }[];
}

const accreditedFacilities: AccreditedPartnerFacility[] = [
  {
    id: "fac-1",
    name: "Everight Diagnostics & Pathology",
    category: "laboratory",
    categoryLabel: "Pathology & Medical Laboratory",
    location: "Ikeja, Lagos",
    accreditation: "ISO 15189 Aligned • MLSCN Certified",
    telemetryCapabilities: ["Mindray BS-800", "Sysmex XN-550", "Roche Cobas"],
    sampleProcedures: [
      { name: "Comprehensive Metabolic Panel (CMP)", tat: "2 Hours", routingType: "Direct Auto-Ingest" },
      { name: "Full Blood Count (5-Part Diff)", tat: "45 mins", routingType: "Automated LIS" },
      { name: "Lipid Profile & Atherogenic Index", tat: "3 Hours", routingType: "Bi-directional Link" },
    ],
  },
  {
    id: "fac-2",
    name: "Apex Precision Imaging & Ultrasound",
    category: "radiology",
    categoryLabel: "Radiology & Imaging Center",
    location: "Maitama, Abuja",
    accreditation: "NNRA Compliant • Board Radiologists",
    telemetryCapabilities: ["GE Voluson E10", "Digital X-Ray PA", "Multi-Slice CT"],
    sampleProcedures: [
      { name: "High-Resolution Pelvic Ultrasound", tat: "30 mins", routingType: "Web DICOM Link" },
      { name: "Digital Chest X-Ray (PA View)", tat: "1 Hour", routingType: "PACS Viewing" },
      { name: "Multi-Slice CT Scan (Head/Neck)", tat: "Same Day", routingType: "Structured Report" },
    ],
  },
  {
    id: "fac-3",
    name: "Curexal Outpatient Medical Center",
    category: "clinic",
    categoryLabel: "Outpatient Clinic & Practice",
    location: "Victoria Island, Lagos",
    accreditation: "HEFAMAA Certified • MDCN Licensed",
    telemetryCapabilities: ["Curexal Clinic OS", "Live SOAP Canvas", "Electronic Rx Engine"],
    sampleProcedures: [
      { name: "General Practice & Family Medicine", tat: "Immediate", routingType: "Clinical Encounter" },
      { name: "Endocrinology Specialist Review", tat: "By Schedule", routingType: "Specialist Queue" },
      { name: "Rapid Nurse Triage & Vitals Intake", tat: "5 mins", routingType: "MPI & Triage" },
    ],
  },
  {
    id: "fac-4",
    name: "Union Diagnostic & Clinical Services",
    category: "laboratory",
    categoryLabel: "Pathology & Molecular Lab",
    location: "Surulere, Lagos",
    accreditation: "MLSCN Registered • ISO Aligned",
    telemetryCapabilities: ["Automated PCR", "Immunoassay Telemetry"],
    sampleProcedures: [
      { name: "Thyroid Panel (TSH, Free T3, Free T4)", tat: "4 Hours", routingType: "Automated LIS" },
      { name: "HbA1c Glycated Hemoglobin", tat: "1 Hour", routingType: "Direct Telemetry" },
      { name: "Liver Function Tests (Full Panel)", tat: "2 Hours", routingType: "Pathologist Sign-off" },
    ],
  },
];

export function MarketplacePage() {
  const [searchQuery, setSearchQuery] = useState("");
  const [categoryFilter, setCategoryFilter] = useState<"all" | "clinic" | "laboratory" | "radiology">("all");
  const [waitlistOpen, setWaitlistOpen] = useState(false);

  const filteredFacilities = accreditedFacilities.filter((f) => {
    const matchesCategory = categoryFilter === "all" || f.category === categoryFilter;
    const matchesSearch =
      f.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      f.location.toLowerCase().includes(searchQuery.toLowerCase()) ||
      f.sampleProcedures.some((s) => s.name.toLowerCase().includes(searchQuery.toLowerCase()));
    return matchesCategory && matchesSearch;
  });

  return (
    <div className="min-h-screen bg-white dark:bg-[#0B1120] font-inter">
      <SEOHead
        title="Diagnostic & Referral Network: Connected Labs, Clinics & Imaging Centers"
        description="B2B healthcare referral highway. Connect outpatient clinics with accredited diagnostic laboratories and radiology centers for paperless requisitions and automated result callbacks."
      />
      <MarketingNavbar />

      {/* Hero Section */}
      <div className="pt-28 pb-16 bg-[#F8FAFC] dark:bg-[#0B1120] border-b border-gray-100 dark:border-[#1F2937]">
        <div className="max-w-[1280px] mx-auto px-6">
          <div className="max-w-3xl">
            <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-teal-500/10 border border-teal-500/30 text-[#0F766E] dark:text-teal-400 text-xs font-bold uppercase tracking-wider mb-4">
              <Share2 className="w-3.5 h-3.5" />
              <span>B2B Diagnostic & Referral Highway</span>
            </div>

            <h1 className="text-3xl sm:text-4xl lg:text-5xl font-black tracking-tight text-gray-900 dark:text-white mb-5 leading-tight">
              The Connected Referral Highway for Clinics, Labs & Imaging.
            </h1>

            <p className="text-base sm:text-lg text-gray-600 dark:text-gray-300 leading-relaxed mb-8">
              Eliminate paper referral slips and delayed phone calls. Enable your doctors to dispatch electronic test requisitions directly to partner diagnostic centers, with verified reports returning automatically to the patient chart.
            </p>

            {/* Network Node Search */}
            <div className="flex flex-col sm:flex-row items-center gap-3 p-2 bg-white dark:bg-[#131C31] rounded-2xl border border-gray-200 dark:border-gray-700 shadow-xl max-w-2xl">
              <div className="flex items-center gap-2.5 px-3 flex-1 w-full">
                <Search className="w-4 h-4 text-muted-foreground" />
                <input
                  type="text"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  placeholder="Search partner facility, blood test, ultrasound, or city..."
                  className="w-full bg-transparent text-sm text-foreground focus:outline-none placeholder:text-muted-foreground"
                />
              </div>
              <Link
                to="/book-demo"
                className="w-full sm:w-auto px-6 py-2.5 bg-[#0F766E] hover:bg-[#115E59] text-white text-xs font-bold rounded-xl transition-colors cursor-pointer shadow flex items-center justify-center gap-1.5"
              >
                <Sparkles className="w-3.5 h-3.5" />
                Connect Facility
              </Link>
            </div>

            {/* Category Filter Pills */}
            <div className="flex flex-wrap items-center gap-2 mt-6">
              {[
                { id: "all", label: "All Partner Facilities", icon: Layers },
                { id: "laboratory", label: "Diagnostic Labs (LIS)", icon: FlaskConical },
                { id: "radiology", label: "Imaging Centers (RIS)", icon: Radio },
                { id: "clinic", label: "Outpatient Clinics (EMR)", icon: Stethoscope },
              ].map(({ id, label, icon: Icon }) => (
                <button
                  key={id}
                  onClick={() => setCategoryFilter(id as any)}
                  className={`flex items-center gap-1.5 px-3.5 py-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer ${
                    categoryFilter === id
                      ? "bg-[#0F766E] text-white shadow-sm"
                      : "bg-white dark:bg-gray-800 text-slate-700 dark:text-slate-300 border border-slate-200 dark:border-slate-700 hover:border-slate-300"
                  }`}
                >
                  <Icon className="w-3.5 h-3.5" />
                  {label}
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Directory Preview Grid */}
      <main className="max-w-[1280px] mx-auto px-6 py-14">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8">
          <div>
            <h2 className="text-xl font-bold text-gray-900 dark:text-white">
              Accredited Diagnostic Network Directory
            </h2>
            <p className="text-xs sm:text-sm text-gray-500 dark:text-gray-400">
              Verified healthcare organizations connected to the Curexal electronic requisition exchange.
            </p>
          </div>
          <div className="flex items-center gap-2">
            <span className="text-xs font-semibold px-3 py-1 bg-emerald-500/10 border border-emerald-500/30 text-emerald-600 dark:text-emerald-400 rounded-full flex items-center gap-1.5">
              <CheckCircle2 className="w-3 h-3" />
              Active B2B Network
            </span>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {filteredFacilities.map((fac) => (
            <div
              key={fac.id}
              className="p-6 rounded-2xl bg-white dark:bg-[#131C31] border border-gray-200/80 dark:border-gray-800 hover:border-[#0F766E]/50 hover:shadow-lg transition-all flex flex-col justify-between group"
            >
              <div>
                {/* Header */}
                <div className="flex items-start justify-between gap-3 mb-3">
                  <div>
                    <span className="text-[10px] font-extrabold uppercase tracking-wider text-[#0F766E] dark:text-teal-400 bg-teal-50 dark:bg-teal-950/60 px-2 py-0.5 rounded-md border border-teal-200/60 dark:border-teal-800/60 inline-block mb-1.5">
                      {fac.categoryLabel}
                    </span>
                    <h3 className="text-base font-bold text-gray-900 dark:text-white group-hover:text-[#0F766E] transition-colors">
                      {fac.name}
                    </h3>
                  </div>
                  <div className="flex items-center gap-1 text-emerald-600 dark:text-emerald-400 text-xs font-bold bg-emerald-500/10 px-2 py-1 rounded-lg">
                    <Activity className="w-3.5 h-3.5" />
                    <span>Live Node</span>
                  </div>
                </div>

                {/* Location & Accreditation */}
                <div className="flex flex-wrap items-center gap-3 text-xs text-gray-500 dark:text-gray-400 mb-4">
                  <span className="flex items-center gap-1">
                    <MapPin className="w-3.5 h-3.5 text-muted-foreground" />
                    {fac.location}
                  </span>
                  <span>•</span>
                  <span className="flex items-center gap-1 text-emerald-600 dark:text-emerald-400 font-medium">
                    <ShieldCheck className="w-3.5 h-3.5" />
                    {fac.accreditation}
                  </span>
                </div>

                {/* Telemetry Links */}
                <div className="mb-4">
                  <p className="text-[10px] font-bold uppercase tracking-wider text-slate-400 mb-1.5">
                    Integrated Diagnostic Equipment
                  </p>
                  <div className="flex flex-wrap gap-1.5">
                    {fac.telemetryCapabilities.map((tel) => (
                      <span key={tel} className="text-[10px] font-mono font-medium px-2 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300">
                        {tel}
                      </span>
                    ))}
                  </div>
                </div>

                {/* Sample Procedures Table */}
                <div className="space-y-2 mb-6">
                  <p className="text-[10px] font-bold uppercase tracking-wider text-slate-400">
                    Electronic Requisition Categories
                  </p>
                  <div className="divide-y divide-gray-100 dark:divide-gray-800 rounded-xl bg-gray-50 dark:bg-gray-900/50 p-3 text-xs">
                    {fac.sampleProcedures.map((srv) => (
                      <div key={srv.name} className="py-2 first:pt-0 last:pb-0 flex items-center justify-between gap-2">
                        <span className="text-gray-700 dark:text-gray-300 font-medium">{srv.name}</span>
                        <div className="flex items-center gap-2 flex-shrink-0">
                          <span className="text-[10px] font-bold text-teal-600 dark:text-teal-400 bg-teal-50 dark:bg-teal-950/60 px-1.5 py-0.5 rounded border border-teal-200 dark:border-teal-800">
                            {srv.routingType}
                          </span>
                          <span className="text-[10px] text-muted-foreground bg-white dark:bg-gray-800 px-1.5 py-0.5 rounded border border-gray-200 dark:border-gray-700">
                            {srv.tat}
                          </span>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>

              {/* Action Button */}
              <Link
                to="/book-demo"
                className="w-full py-2.5 rounded-xl bg-teal-50 dark:bg-teal-950/40 border border-teal-200 dark:border-teal-800 text-[#0F766E] dark:text-teal-300 text-xs font-bold hover:bg-[#0F766E] hover:text-white dark:hover:bg-[#0F766E] dark:hover:text-white transition-all cursor-pointer flex items-center justify-center gap-1.5"
              >
                <span>Connect with this Facility</span>
                <ArrowRight className="w-3.5 h-3.5" />
              </Link>
            </div>
          ))}
        </div>

        {/* Facility Onboarding Banner */}
        <div className="mt-16 p-8 sm:p-12 rounded-3xl bg-gradient-to-br from-slate-900 via-[#0B1A28] to-teal-950 border border-teal-500/30 text-white shadow-2xl relative overflow-hidden">
          <div className="relative z-10 max-w-2xl">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-teal-500/20 text-teal-300 text-xs font-bold uppercase tracking-wider mb-4 border border-teal-500/30">
              <UserPlus className="w-3.5 h-3.5" />
              <span>For Diagnostic Labs, Imaging Centers & Clinics</span>
            </div>
            <h2 className="text-2xl sm:text-3xl font-extrabold tracking-tight mb-4">
              Join the Curexal Referral Highway.
            </h2>
            <p className="text-sm sm:text-base text-slate-300 leading-relaxed mb-6">
              Connect your laboratory or radiology center to partner clinics in your area. Receive verified electronic requisitions, speed up turnaround times, and expand your diagnostic testing volume with zero paperwork.
            </p>
            <div className="flex flex-wrap items-center gap-4">
              <Link
                to="/book-demo"
                className="px-6 py-3 bg-[#0F766E] hover:bg-[#115E59] text-white text-sm font-bold rounded-xl transition-all cursor-pointer shadow-lg shadow-teal-900/50 flex items-center gap-2"
              >
                <span>Register Facility Node</span>
                <ArrowRight className="w-4 h-4" />
              </Link>
              <button
                onClick={() => setWaitlistOpen(true)}
                className="text-sm font-semibold text-slate-300 hover:text-white underline underline-offset-4 cursor-pointer"
              >
                Request Interoperability Whitepaper
              </button>
            </div>
          </div>
        </div>
      </main>

      <WaitlistModal open={waitlistOpen} onOpenChange={setWaitlistOpen} />
      <MarketingFooter />
    </div>
  );
}
