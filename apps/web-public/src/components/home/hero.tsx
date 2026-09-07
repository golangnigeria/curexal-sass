import { useState } from "react";
import { Link } from "react-router-dom";
import {
  ArrowRight,
  ShieldCheck,
  Stethoscope,
  FlaskConical,
  Radio,
  Pill,
  User,
  Network,
  RefreshCw,
  CheckCircle2,
  Lock,
  Clock,
  FileText,
  Activity,
  Send,
  Building2,
  ChevronRight,
  Database,
  Sparkles,
  Barcode,
  Layers,
} from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import { WaitlistModal } from "@/components/waitlist-modal";

export function Hero() {
  const [waitlistOpen, setWaitlistOpen] = useState(false);
  const [activeTab, setActiveTab] = useState<"clinic" | "lab" | "radiology">("clinic");

  return (
    <>
      <section
        id="hero"
        className="relative z-0 isolate pt-14 sm:pt-20 lg:pt-24 min-h-[90vh] flex items-center overflow-hidden bg-white dark:bg-[#0B1120]"
      >
        {/* Ambient background effects */}
        <div className="absolute inset-0 dot-grid opacity-30 pointer-events-none" />

        {/* Top-center teal glow */}
        <div
          className="absolute top-0 left-1/2 -translate-x-1/2 w-[900px] h-[350px] pointer-events-none"
          style={{
            background: "radial-gradient(ellipse at center, rgba(15,118,110,0.08) 0%, transparent 70%)",
          }}
        />

        <div className="relative z-10 max-w-[1280px] mx-auto px-4 sm:px-6 py-8 sm:py-16 lg:py-20 w-full">
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-12 items-center">

            {/* Left: Positioning & Value Proposition */}
            <div className="lg:col-span-6 flex flex-col items-start text-left">
              <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-teal-500/10 border border-teal-500/30 text-[#0F766E] dark:text-teal-400 text-xs font-bold uppercase tracking-wider mb-4">
                <Layers className="w-3.5 h-3.5" />
                <span>Unified Cloud Platform for Clinics, Labs & Imaging</span>
              </div>

              {/* Main Headline */}
              <h1 className="text-3xl xs:text-4xl sm:text-[46px] lg:text-[52px] font-black leading-[1.1] tracking-tight text-slate-900 dark:text-white mb-4 sm:mb-6">
                The Operating System for{" "}
                <span className="bg-clip-text text-transparent bg-gradient-to-r from-[#0F766E] via-[#0D9488] to-[#14B8A6]">
                  Modern Healthcare Facilities.
                </span>
              </h1>

              {/* Supporting Line */}
              <p className="text-base sm:text-lg leading-relaxed text-slate-600 dark:text-slate-300 max-w-xl mb-6 sm:mb-8 font-normal">
                Purpose-built workspaces for outpatient clinics, diagnostic laboratories, and radiology centers. Eliminate lost paper records, automate cashier billing, and connect referrals without friction.
              </p>

              {/* Facility Quick Switcher Pills */}
              <div className="flex flex-wrap items-center gap-2 mb-6 w-full">
                <span className="text-xs font-bold text-slate-400 uppercase tracking-wider block w-full sm:w-auto mr-1">
                  Preview Facility:
                </span>
                {[
                  { id: "clinic", label: "Clinic EMR", icon: Stethoscope },
                  { id: "lab", label: "Diagnostic LIS", icon: FlaskConical },
                  { id: "radiology", label: "Radiology RIS", icon: Radio },
                ].map(({ id, label, icon: Icon }) => (
                  <button
                    key={id}
                    onClick={() => setActiveTab(id as any)}
                    className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer border ${
                      activeTab === id
                        ? "bg-[#0F766E] text-white border-[#0F766E] shadow-sm"
                        : "bg-slate-100 dark:bg-slate-800/80 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-700 hover:border-teal-500/50"
                    }`}
                  >
                    <Icon className="w-3.5 h-3.5" />
                    <span>{label}</span>
                  </button>
                ))}
              </div>

              {/* CTAs */}
              <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 w-full sm:w-auto">
                <Link
                  to="/book-demo"
                  id="hero-primary-cta"
                  className="group flex items-center justify-center gap-2 px-6 py-3.5 rounded-xl bg-[#0F766E] hover:bg-[#115E59] text-white text-sm font-bold transition-all shadow-md cursor-pointer border-0"
                >
                  <span>Book Facility Demo</span>
                  <ArrowRight className="h-4 w-4 group-hover:translate-x-1 transition-transform" />
                </Link>

                <a
                  href="#facility-workspaces"
                  className="flex items-center justify-center gap-2 px-6 py-3.5 rounded-xl border border-slate-300 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/60 text-slate-800 dark:text-slate-200 text-sm font-semibold hover:bg-slate-100 dark:hover:bg-slate-800 transition-all cursor-pointer shadow-xs"
                >
                  <span>Explore Workspaces</span>
                </a>
              </div>

              {/* Key Trust & Architecture Pillars */}
              <div className="flex flex-wrap items-center gap-3 sm:gap-6 mt-8 sm:mt-10 text-slate-500 dark:text-slate-400 text-xs font-medium border-t border-slate-100 dark:border-slate-800/80 pt-6 w-full">
                <div className="flex items-center gap-1.5">
                  <ShieldCheck className="w-4 h-4 text-[#0F766E]" />
                  <span>Tenant Data Isolation</span>
                </div>
                <div className="w-px h-3.5 bg-slate-200 dark:bg-slate-800 hidden sm:block" />
                <div className="flex items-center gap-1.5">
                  <RefreshCw className="w-4 h-4 text-[#0F766E]" />
                  <span>Cross-Facility Referrals</span>
                </div>
                <div className="w-px h-3.5 bg-slate-200 dark:bg-slate-800 hidden sm:block" />
                <div className="flex items-center gap-1.5">
                  <Network className="w-4 h-4 text-[#0F766E]" />
                  <span>Zero-Paper Diagnostic Loop</span>
                </div>
              </div>

            </div>

            {/* Right: Realistic Multi-Facility Healthcare Operating Console */}
            <div className="lg:col-span-6 flex items-center justify-center mt-4 lg:mt-0 w-full">
              <div className="w-full max-w-[500px] bg-slate-900/95 dark:bg-[#0c1322]/95 border border-slate-700/80 dark:border-slate-800/90 rounded-2xl shadow-2xl overflow-hidden backdrop-blur-xl">
                
                {/* Console Header Bar */}
                <div className="px-4 py-3 bg-slate-950/80 border-b border-slate-800 flex items-center justify-between gap-2">
                  <div className="flex items-center gap-2">
                    <div className="flex items-center gap-1.5">
                      <div className="w-2.5 h-2.5 rounded-full bg-rose-500/80" />
                      <div className="w-2.5 h-2.5 rounded-full bg-amber-500/80" />
                      <div className="w-2.5 h-2.5 rounded-full bg-emerald-500/80" />
                    </div>
                    <span className="text-[11px] font-mono text-slate-400 pl-1">
                      curexal-facility-os // v2.6
                    </span>
                  </div>

                  <div className="flex items-center gap-1.5 px-2 py-0.5 rounded bg-emerald-500/10 border border-emerald-500/20 text-[10px] font-mono text-emerald-400 font-bold">
                    <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse" />
                    <span>TENANT ISOLATED</span>
                  </div>
                </div>

                {/* Facility Node Switcher Tabs */}
                <div className="p-2.5 bg-slate-900 border-b border-slate-800/80 flex items-center gap-1.5">
                  <button
                    onClick={() => setActiveTab("clinic")}
                    className={`flex-1 py-1.5 px-2 rounded-lg text-xs font-semibold flex items-center justify-center gap-1.5 transition-all cursor-pointer border-0 ${
                      activeTab === "clinic"
                        ? "bg-[#0F766E] text-white shadow-sm font-bold"
                        : "bg-slate-800/60 text-slate-400 hover:text-slate-200 hover:bg-slate-800"
                    }`}
                  >
                    <Stethoscope className="w-3.5 h-3.5" />
                    <span className="truncate">1. Clinic EMR</span>
                  </button>

                  <button
                    onClick={() => setActiveTab("lab")}
                    className={`flex-1 py-1.5 px-2 rounded-lg text-xs font-semibold flex items-center justify-center gap-1.5 transition-all cursor-pointer border-0 ${
                      activeTab === "lab"
                        ? "bg-[#0F766E] text-white shadow-sm font-bold"
                        : "bg-slate-800/60 text-slate-400 hover:text-slate-200 hover:bg-slate-800"
                    }`}
                  >
                    <FlaskConical className="w-3.5 h-3.5" />
                    <span className="truncate">2. Lab LIS</span>
                  </button>

                  <button
                    onClick={() => setActiveTab("radiology")}
                    className={`flex-1 py-1.5 px-2 rounded-lg text-xs font-semibold flex items-center justify-center gap-1.5 transition-all cursor-pointer border-0 ${
                      activeTab === "radiology"
                        ? "bg-[#0F766E] text-white shadow-sm font-bold"
                        : "bg-slate-800/60 text-slate-400 hover:text-slate-200 hover:bg-slate-800"
                    }`}
                  >
                    <Radio className="w-3.5 h-3.5" />
                    <span className="truncate">3. Radiology RIS</span>
                  </button>
                </div>

                {/* Console Content Window */}
                <div className="p-4 sm:p-5 space-y-4">
                  <AnimatePresence mode="wait">
                    {activeTab === "clinic" && (
                      <motion.div
                        key="clinic"
                        initial={{ opacity: 0, y: 8 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -8 }}
                        transition={{ duration: 0.2 }}
                        className="space-y-3"
                      >
                        <div className="flex items-center justify-between text-xs text-slate-400 pb-2 border-b border-slate-800">
                          <span className="font-semibold text-white">Curexal Outpatient Clinic EMR</span>
                          <span className="font-mono text-[11px] text-teal-400">Dr. M. Adebayo, MD</span>
                        </div>

                        <div className="p-3 rounded-xl bg-slate-950/70 border border-slate-800 space-y-2">
                          <div className="flex items-center justify-between text-xs">
                            <span className="text-slate-400">Encounter SOAP Canvas</span>
                            <span className="font-mono text-teal-300 font-bold">MRN-2026-8891</span>
                          </div>
                          <p className="text-xs font-bold text-white">Patient: Amara Eze (34y, Female) • BP 124/82</p>
                          <div className="flex flex-wrap gap-1.5 pt-1">
                            <span className="px-2 py-0.5 rounded bg-slate-800 text-[10px] text-slate-300 font-mono">ICD-10: E11.9 (Type 2 Diabetes)</span>
                            <span className="px-2 py-0.5 rounded bg-teal-900/60 text-teal-300 text-[10px] font-mono">Metformin 500mg PO</span>
                          </div>
                        </div>

                        <div className="flex items-center justify-between px-3 py-2 rounded-lg bg-teal-950/40 border border-teal-800/50 text-[11px]">
                          <span className="text-teal-300 flex items-center gap-1.5">
                            <Send className="w-3.5 h-3.5 text-teal-400" />
                            Electronic Lab Order Routed to Partner LIS
                          </span>
                          <span className="text-emerald-400 font-mono font-bold">Instant Dispatch</span>
                        </div>
                      </motion.div>
                    )}

                    {activeTab === "lab" && (
                      <motion.div
                        key="lab"
                        initial={{ opacity: 0, y: 8 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -8 }}
                        transition={{ duration: 0.2 }}
                        className="space-y-3"
                      >
                        <div className="flex items-center justify-between text-xs text-slate-400 pb-2 border-b border-slate-800">
                          <span className="font-semibold text-white">Curexal Laboratory LIS</span>
                          <span className="font-mono text-[11px] text-teal-400">Accessioning Desk</span>
                        </div>

                        <div className="p-3 rounded-xl bg-slate-950/70 border border-slate-800 space-y-2.5">
                          <div className="flex items-center justify-between text-xs">
                            <span className="text-slate-400 font-mono">BARCODE #LAB-9941</span>
                            <span className="flex items-center gap-1 text-emerald-400 font-bold text-[10px]">
                              <CheckCircle2 className="w-3.5 h-3.5" />
                              VERIFIED & SIGNED
                            </span>
                          </div>

                          <div className="grid grid-cols-2 gap-2 text-[11px]">
                            <div className="p-2 rounded-lg bg-slate-900 border border-slate-800">
                              <p className="text-[10px] text-slate-400">Analyzer Link</p>
                              <p className="font-mono font-bold text-white text-xs mt-0.5">Mindray BS-800 Telemetry</p>
                            </div>
                            <div className="p-2 rounded-lg bg-slate-900 border border-slate-800">
                              <p className="text-[10px] text-slate-400">Pathologist Sign-off</p>
                              <p className="font-mono font-bold text-teal-300 text-xs mt-0.5">Dr. C. Okonjo, FRCPath</p>
                            </div>
                          </div>
                        </div>

                        <div className="flex items-center justify-between px-3 py-2 rounded-lg bg-emerald-950/40 border border-emerald-800/50 text-[11px]">
                          <span className="text-emerald-300 flex items-center gap-1.5">
                            <FileText className="w-3.5 h-3.5 text-emerald-400" />
                            Auto-dispatched via WhatsApp & Chart Sync
                          </span>
                          <span className="text-emerald-400 font-mono font-bold">TAT: 32 mins</span>
                        </div>
                      </motion.div>
                    )}

                    {activeTab === "radiology" && (
                      <motion.div
                        key="radiology"
                        initial={{ opacity: 0, y: 8 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -8 }}
                        transition={{ duration: 0.2 }}
                        className="space-y-3"
                      >
                        <div className="flex items-center justify-between text-xs text-slate-400 pb-2 border-b border-slate-800">
                          <span className="font-semibold text-white">Curexal Radiology RIS & PACS</span>
                          <span className="font-mono text-[11px] text-teal-400">Modality Review</span>
                        </div>

                        <div className="p-3 rounded-xl bg-slate-950/70 border border-slate-800 space-y-2.5">
                          <div className="flex items-center justify-between text-xs">
                            <span className="text-slate-400 font-mono">SCAN #RAD-3042</span>
                            <span className="flex items-center gap-1 text-teal-400 font-bold text-[10px]">
                              <Activity className="w-3.5 h-3.5" />
                              DICOM INGESTED
                            </span>
                          </div>

                          <div className="grid grid-cols-2 gap-2 text-[11px]">
                            <div className="p-2 rounded-lg bg-slate-900 border border-slate-800">
                              <p className="text-[10px] text-slate-400">Modality Type</p>
                              <p className="font-mono font-bold text-white text-xs mt-0.5">High-Res Digital X-Ray</p>
                            </div>
                            <div className="p-2 rounded-lg bg-slate-900 border border-slate-800">
                              <p className="text-[10px] text-slate-400">Reporting Radiologist</p>
                              <p className="font-mono font-bold text-teal-300 text-xs mt-0.5">Dr. T. Danladi, FWACS</p>
                            </div>
                          </div>
                        </div>

                        <div className="flex items-center justify-between px-3 py-2 rounded-lg bg-teal-950/40 border border-teal-800/50 text-[11px]">
                          <span className="text-teal-300 flex items-center gap-1.5">
                            <Radio className="w-3.5 h-3.5 text-teal-400" />
                            Zero-Film Web DICOM Link Ready
                          </span>
                          <span className="text-teal-300 font-mono font-bold">Signed in 24m</span>
                        </div>
                      </motion.div>
                    )}
                  </AnimatePresence>

                  {/* Multi-Tenant Security & Network Status Footer */}
                  <div className="pt-2 border-t border-slate-800/80 flex items-center justify-between text-[10px] font-mono text-slate-400">
                    <div className="flex items-center gap-1.5 text-slate-400">
                      <Lock className="w-3 h-3 text-[#0F766E]" />
                      <span>Schema Isolation: 100%</span>
                    </div>
                    <div className="flex items-center gap-1.5 text-teal-400">
                      <Activity className="w-3 h-3" />
                      <span>FHIR & DICOM Ready</span>
                    </div>
                  </div>
                </div>

              </div>
            </div>

          </div>
        </div>
      </section>

      {/* Waitlist Modal */}
      <WaitlistModal open={waitlistOpen} onOpenChange={setWaitlistOpen} />
    </>
  );
}
