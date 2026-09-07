import { useState } from "react";
import { Link } from "react-router-dom";
import { motion, AnimatePresence } from "framer-motion";
import {
  Stethoscope,
  FlaskConical,
  Radio,
  CheckCircle2,
  ArrowRight,
  ShieldCheck,
  Zap,
  Layers,
  Sparkles,
} from "lucide-react";

interface FacilityPillar {
  id: "clinic" | "lab" | "radiology";
  icon: any;
  title: string;
  badge: string;
  tagline: string;
  description: string;
  highlightMetric: { value: string; label: string };
  features: string[];
  workflowStages: { step: string; detail: string }[];
  ctaText: string;
  ctaHref: string;
}

const facilityPillars: FacilityPillar[] = [
  {
    id: "clinic",
    icon: Stethoscope,
    title: "Curexal Clinic OS",
    badge: "Outpatient Clinics & Practices",
    tagline: "The Fast, Paperless Practice Operating System",
    description:
      "Replace lost paper files, eliminate waiting room chaos, and stop revenue leakage with a unified clinical workspace designed for doctors, nurses, and reception staff.",
    highlightMetric: { value: "60%", label: "Reduction in Patient Intake Wait Times" },
    features: [
      "Master Patient Index (MPI) with instant duplicate search",
      "Nurse triage desk with vital signs, BMI & acuity triage",
      "Distraction-free doctor SOAP notes & ICD-10 coding",
      "Digital e-prescriptions with direct pharmacy routing",
      "Multi-tender cashier POS (Cash, Card, Transfer reconciliation)",
      "Automated e-referrals to partner diagnostic labs & imaging",
    ],
    workflowStages: [
      { step: "01. Intake", detail: "Instant MRN & queue dispatch" },
      { step: "02. Triage", detail: "Vitals & allergy capture" },
      { step: "03. Consult", detail: "SOAP canvas & ICD-10" },
      { step: "04. Settle", detail: "Reconciled POS invoice" },
    ],
    ctaText: "Explore Clinic OS",
    ctaHref: "/solutions#clinic-os",
  },
  {
    id: "lab",
    icon: FlaskConical,
    title: "Curexal Lab LIS",
    badge: "Diagnostic & Pathology Labs",
    tagline: "Precision Specimen Tracking & Automated Analyzer Telemetry",
    description:
      "A cloud-native Laboratory Information System for standalone pathology and medical labs. Automate specimen barcoding, ingest analyzer data, and deliver verified PDF reports to doctors & patients.",
    highlightMetric: { value: "35 mins", label: "Average Routine Turnaround Time" },
    features: [
      "Phlebotomy accessioning with unique barcode label printing",
      "Bidirectional analyzer telemetry (Sysmex, Mindray, Roche)",
      "Two-step quality control: scientist validation to pathologist sign-off",
      "Direct electronic test orders from partner clinics (zero paper)",
      "Automated report delivery via encrypted WhatsApp, SMS & Email",
      "ISO 15189 aligned audit trail & specimen chain-of-custody",
    ],
    workflowStages: [
      { step: "01. Accession", detail: "Barcode tag & rack sort" },
      { step: "02. Analysis", detail: "Direct analyzer link" },
      { step: "03. Validate", detail: "Pathologist digital sign" },
      { step: "04. Dispatch", detail: "WhatsApp & clinic sync" },
    ],
    ctaText: "Explore Lab LIS",
    ctaHref: "/solutions#lims",
  },
  {
    id: "radiology",
    icon: Radio,
    title: "Curexal Radiology RIS",
    badge: "Imaging & Diagnostic Centers",
    tagline: "Accelerate Imaging Modalities from Scan to Signed Report",
    description:
      "Manage Ultrasound, Digital X-Ray, CT, and MRI modalities with high-speed web DICOM viewing, structured radiologist templates, and instant digital sharing that eliminates physical film costs.",
    highlightMetric: { value: "100%", label: "Digital Imaging (Zero Film Waste)" },
    features: [
      "Modality appointment scheduler & live patient queue",
      "Web-based DICOM PACS viewer accessible from any device",
      "Standardized structured radiology reporting templates",
      "Digital electronic signature with timestamp verification",
      "Zero-film digital report links sent to referrers and patients",
      "Integrated billing for complex diagnostic imaging procedures",
    ],
    workflowStages: [
      { step: "01. Schedule", detail: "Modality intake & booking" },
      { step: "02. Acquire", detail: "DICOM PACS ingestion" },
      { step: "03. Report", detail: "Structured radiologist note" },
      { step: "04. Deliver", detail: "Digital image & PDF link" },
    ],
    ctaText: "Explore Radiology RIS",
    ctaHref: "/solutions#ris-pacs",
  },
];

export function FacilitySolutionsGrid() {
  const [activeTab, setActiveTab] = useState<"clinic" | "lab" | "radiology">("clinic");
  const activePillar = facilityPillars.find((p) => p.id === activeTab)!;
  const IconComponent = activePillar.icon;

  return (
    <section id="facility-workspaces" className="py-16 sm:py-24 bg-[#F8FAFC] dark:bg-[#080D18] text-slate-900 dark:text-white border-y border-slate-200/80 dark:border-slate-800 relative">
      <div className="max-w-[1280px] mx-auto px-4 sm:px-6">
        
        {/* Section Header */}
        <div className="max-w-3xl mx-auto text-center mb-12 sm:mb-16">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-teal-500/10 border border-teal-500/20 text-[#0F766E] dark:text-teal-400 text-xs font-bold uppercase tracking-wider mb-4">
            <Layers className="w-3.5 h-3.5" />
            <span>Dedicated Workspaces by Facility Type</span>
          </div>

          <h2 className="text-2xl sm:text-4xl lg:text-5xl font-black text-slate-900 dark:text-white tracking-tight leading-tight mb-4">
            One Operating Core. <br className="hidden sm:inline" />
            <span className="bg-clip-text text-transparent bg-gradient-to-r from-[#0F766E] via-[#0D9488] to-[#14B8A6]">
              Specialized Software for Your Exact Facility.
            </span>
          </h2>

          <p className="text-sm sm:text-base text-slate-600 dark:text-slate-300 leading-relaxed">
            Whether you operate a standalone medical lab, an imaging center, or a busy outpatient clinic, Curexal provides native workspaces tailored to your daily operations.
          </p>
        </div>

        {/* Interactive Facility Switcher Tabs */}
        <div className="flex justify-center mb-10">
          <div className="inline-flex p-1.5 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-md">
            {facilityPillars.map((p) => {
              const PIcon = p.icon;
              const isSelected = activeTab === p.id;
              return (
                <button
                  key={p.id}
                  onClick={() => setActiveTab(p.id)}
                  className={`flex items-center gap-2 px-4 sm:px-6 py-2.5 rounded-xl text-xs sm:text-sm font-bold transition-all cursor-pointer border-0 ${
                    isSelected
                      ? "bg-[#0F766E] text-white shadow-md"
                      : "text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-100 dark:hover:bg-slate-800"
                  }`}
                >
                  <PIcon className="w-4 h-4" />
                  <span>{p.title}</span>
                </button>
              );
            })}
          </div>
        </div>

        {/* Dynamic Facility Card Display */}
        <div className="max-w-5xl mx-auto">
          <AnimatePresence mode="wait">
            <motion.div
              key={activePillar.id}
              initial={{ opacity: 0, y: 16 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -16 }}
              transition={{ duration: 0.25 }}
              className="p-6 sm:p-10 rounded-3xl bg-white dark:bg-[#0F172A] border border-slate-200 dark:border-slate-800 shadow-xl"
            >
              <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
                
                {/* Left Column: Context & Features */}
                <div className="lg:col-span-7 space-y-6">
                  <div className="flex flex-wrap items-center gap-2.5">
                    <span className="px-2.5 py-1 rounded-lg bg-teal-50 dark:bg-teal-950/60 border border-teal-200 dark:border-teal-800 text-[#0F766E] dark:text-teal-300 text-xs font-extrabold uppercase tracking-wider">
                      {activePillar.badge}
                    </span>
                    <span className="text-xs text-slate-400 font-medium">•</span>
                    <span className="text-xs font-semibold text-emerald-600 dark:text-emerald-400 flex items-center gap-1">
                      <ShieldCheck className="w-3.5 h-3.5" />
                      Isolated Tenant Cloud Partition
                    </span>
                  </div>

                  <div>
                    <h3 className="text-2xl sm:text-3xl font-black text-slate-900 dark:text-white tracking-tight mb-2">
                      {activePillar.tagline}
                    </h3>
                    <p className="text-sm text-slate-600 dark:text-slate-300 leading-relaxed">
                      {activePillar.description}
                    </p>
                  </div>

                  {/* Operational Metric Callout */}
                  <div className="p-4 rounded-2xl bg-teal-50/50 dark:bg-teal-950/20 border border-teal-200/60 dark:border-teal-800/40 flex items-center gap-4">
                    <div className="w-12 h-12 rounded-xl bg-[#0F766E] text-white flex items-center justify-center flex-shrink-0 font-black text-lg shadow-md">
                      <Zap className="w-5 h-5" />
                    </div>
                    <div>
                      <p className="text-xl sm:text-2xl font-black text-slate-900 dark:text-white">
                        {activePillar.highlightMetric.value}
                      </p>
                      <p className="text-xs text-slate-500 dark:text-slate-400 font-medium">
                        {activePillar.highlightMetric.label}
                      </p>
                    </div>
                  </div>

                  {/* Feature Checklist */}
                  <div className="space-y-2.5 pt-2">
                    <p className="text-xs font-bold uppercase tracking-wider text-slate-400 dark:text-slate-500">
                      Standard Capabilities Included
                    </p>
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5 text-xs text-slate-700 dark:text-slate-200">
                      {activePillar.features.map((feat, i) => (
                        <div key={i} className="flex items-start gap-2">
                          <CheckCircle2 className="w-4 h-4 text-[#0F766E] dark:text-teal-400 flex-shrink-0 mt-0.5" />
                          <span className="leading-snug">{feat}</span>
                        </div>
                      ))}
                    </div>
                  </div>

                  {/* Action CTAs */}
                  <div className="flex flex-wrap items-center gap-3 pt-4">
                    <Link
                      to="/book-demo"
                      className="px-6 py-3 rounded-xl bg-[#0F766E] hover:bg-[#115E59] text-white text-xs sm:text-sm font-bold transition-all shadow-md flex items-center gap-2 cursor-pointer"
                    >
                      <span>Book Facility Walkthrough</span>
                      <ArrowRight className="w-4 h-4" />
                    </Link>

                    <Link
                      to={activePillar.ctaHref}
                      className="px-5 py-3 rounded-xl bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700 text-slate-800 dark:text-slate-200 text-xs sm:text-sm font-semibold transition-all cursor-pointer"
                    >
                      {activePillar.ctaText}
                    </Link>
                  </div>
                </div>

                {/* Right Column: Workflow Steps Card */}
                <div className="lg:col-span-5 p-6 rounded-2xl bg-slate-50 dark:bg-slate-950/80 border border-slate-200 dark:border-slate-800 space-y-4">
                  <div className="flex items-center justify-between pb-3 border-b border-slate-200 dark:border-slate-800">
                    <div className="flex items-center gap-2">
                      <div className="w-7 h-7 rounded-lg bg-teal-50 dark:bg-teal-950 flex items-center justify-center text-[#0F766E] dark:text-teal-400">
                        <IconComponent className="w-4 h-4" />
                      </div>
                      <span className="text-xs font-bold text-slate-900 dark:text-white">
                        Standard Operating Flow
                      </span>
                    </div>
                    <span className="text-[10px] font-mono font-bold text-emerald-600 dark:text-emerald-400">
                      LIVE WORKSPACE
                    </span>
                  </div>

                  <div className="space-y-3">
                    {activePillar.workflowStages.map((stg, idx) => (
                      <div
                        key={idx}
                        className="p-3 rounded-xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 flex items-center justify-between gap-3 shadow-xs"
                      >
                        <span className="text-xs font-mono font-bold text-[#0F766E] dark:text-teal-400">
                          {stg.step}
                        </span>
                        <span className="text-xs font-medium text-slate-700 dark:text-slate-300 text-right">
                          {stg.detail}
                        </span>
                      </div>
                    ))}
                  </div>

                  <div className="p-3.5 rounded-xl bg-teal-950/30 border border-teal-800/40 text-[11px] text-teal-300 flex items-start gap-2">
                    <Sparkles className="w-4 h-4 text-teal-400 flex-shrink-0 mt-0.5" />
                    <span className="leading-relaxed">
                      Cross-facility interoperability built-in: Requisitions, specimens, and digital records can move securely between partner facilities.
                    </span>
                  </div>
                </div>

              </div>
            </motion.div>
          </AnimatePresence>
        </div>

      </div>
    </section>
  );
}
