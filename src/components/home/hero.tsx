import { useState } from "react";
import { Link } from "react-router-dom";
import {
  ArrowRight,
  ShieldCheck,
  Stethoscope,
  FlaskConical,
  Pill,
  User,
  Network,
  RefreshCw,
  ShoppingBag,
  CheckCircle2,
  Lock,
  Clock,
  FileText,
  Activity,
  Send,
  Building2,
  ChevronRight,
  Database
} from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import { WaitlistModal } from "@/components/waitlist-modal";

export function Hero() {
  const [waitlistOpen, setWaitlistOpen] = useState(false);
  const [activeTab, setActiveTab] = useState<"clinic" | "lab" | "patient">("lab");

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

            {/* Left: Coordination Content */}
            <div className="lg:col-span-6 flex flex-col items-start text-left">
              <p className="text-base sm:text-lg md:text-xl font-bold tracking-tight text-[#0F766E] dark:text-teal-400 mb-3 sm:mb-4">
                The Connection layer of healthcare
              </p>

              {/* Main Headline */}
              <h1 className="text-3xl xs:text-4xl sm:text-[48px] lg:text-[54px] font-black leading-[1.08] tracking-tight text-slate-900 dark:text-white mb-4 sm:mb-6">
                Healthcare shouldn't work in{" "}
                <span
                  className="bg-clip-text text-transparent bg-gradient-to-r from-[#0F766E] via-[#0D9488] to-[#14B8A6]"
                >
                  silos.
                </span>
              </h1>

              {/* Supporting Line */}
              <p className="text-base sm:text-lg leading-relaxed text-slate-600 dark:text-slate-300 max-w-xl mb-6 sm:mb-8 font-normal">
                Curexal connects patients, clinics, laboratories, pharmacies and healthcare partners so referrals, diagnostics, results and healthcare transactions can move together.
              </p>

              {/* CTAs */}
              <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 w-full sm:w-auto">
                <button
                  onClick={() => setWaitlistOpen(true)}
                  id="hero-primary-cta"
                  className="group flex items-center justify-center gap-2 px-6 py-3.5 rounded-xl bg-[#0F766E] hover:bg-[#115E59] text-white text-sm font-bold transition-all shadow-md cursor-pointer border-0"
                >
                  <span>Join Early Access</span>
                  <ArrowRight className="h-4 w-4 group-hover:translate-x-1 transition-transform" />
                </button>

                <a
                  href="#problem"
                  className="flex items-center justify-center gap-2 px-6 py-3.5 rounded-xl border border-slate-300 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/60 text-slate-800 dark:text-slate-200 text-sm font-semibold hover:bg-slate-100 dark:hover:bg-slate-800 transition-all cursor-pointer shadow-xs"
                >
                  <span>Help Shape Curexal</span>
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
                  <span>Cross-Org Coordination</span>
                </div>
                <div className="w-px h-3.5 bg-slate-200 dark:bg-slate-800 hidden sm:block" />
                <div className="flex items-center gap-1.5">
                  <Network className="w-4 h-4 text-[#0F766E]" />
                  <span>Single Connected Network</span>
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
                      curexal-control-plane // v2.4
                    </span>
                  </div>

                  <div className="flex items-center gap-1.5 px-2 py-0.5 rounded bg-emerald-500/10 border border-emerald-500/20 text-[10px] font-mono text-emerald-400 font-bold">
                    <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse" />
                    <span>SYNCED</span>
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
                    <span className="truncate">2. Lab LIMS</span>
                  </button>

                  <button
                    onClick={() => setActiveTab("patient")}
                    className={`flex-1 py-1.5 px-2 rounded-lg text-xs font-semibold flex items-center justify-center gap-1.5 transition-all cursor-pointer border-0 ${
                      activeTab === "patient"
                        ? "bg-[#0F766E] text-white shadow-sm font-bold"
                        : "bg-slate-800/60 text-slate-400 hover:text-slate-200 hover:bg-slate-800"
                    }`}
                  >
                    <User className="w-3.5 h-3.5" />
                    <span className="truncate">3. Patient Vault</span>
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
                          <span className="font-semibold text-white">St. Nicholas Outpatient Clinic</span>
                          <span className="font-mono text-[11px] text-teal-400">Dr. M. Adebayo, MD</span>
                        </div>

                        <div className="p-3 rounded-xl bg-slate-950/70 border border-slate-800 space-y-2">
                          <div className="flex items-center justify-between text-xs">
                            <span className="text-slate-400">Electronic Requisition</span>
                            <span className="font-mono text-teal-300 font-bold">REQ-2026-8891</span>
                          </div>
                          <p className="text-xs font-bold text-white">Patient: Amara Eze (34y, Female)</p>
                          <div className="flex flex-wrap gap-1.5 pt-1">
                            <span className="px-2 py-0.5 rounded bg-slate-800 text-[10px] text-slate-300 font-mono">Fasting Blood Sugar</span>
                            <span className="px-2 py-0.5 rounded bg-slate-800 text-[10px] text-slate-300 font-mono">Lipid Profile</span>
                            <span className="px-2 py-0.5 rounded bg-slate-800 text-[10px] text-slate-300 font-mono">HbA1c</span>
                          </div>
                        </div>

                        <div className="flex items-center justify-between px-3 py-2 rounded-lg bg-teal-950/40 border border-teal-800/50 text-[11px]">
                          <span className="text-teal-300 flex items-center gap-1.5">
                            <Send className="w-3.5 h-3.5 text-teal-400" />
                            Routed to Everight Pathology Lab
                          </span>
                          <span className="text-emerald-400 font-mono font-bold">Direct Dispatch</span>
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
                          <span className="font-semibold text-white">Everight Diagnostic & Pathology</span>
                          <span className="font-mono text-[11px] text-teal-400">LIS Auto-Ingest Node</span>
                        </div>

                        <div className="p-3 rounded-xl bg-slate-950/70 border border-slate-800 space-y-2.5">
                          <div className="flex items-center justify-between text-xs">
                            <span className="text-slate-400 font-mono">SPECIMEN #SP-9941</span>
                            <span className="flex items-center gap-1 text-emerald-400 font-bold text-[10px]">
                              <CheckCircle2 className="w-3.5 h-3.5" />
                              VERIFIED
                            </span>
                          </div>

                          <div className="grid grid-cols-2 gap-2 text-[11px]">
                            <div className="p-2 rounded-lg bg-slate-900 border border-slate-800">
                              <p className="text-[10px] text-slate-400">Analyzer Link</p>
                              <p className="font-mono font-bold text-white text-xs mt-0.5">Mindray BS-800</p>
                            </div>
                            <div className="p-2 rounded-lg bg-slate-900 border border-slate-800">
                              <p className="text-[10px] text-slate-400">Validation Signer</p>
                              <p className="font-mono font-bold text-teal-300 text-xs mt-0.5">Dr. C. Okonjo, FRCPath</p>
                            </div>
                          </div>
                        </div>

                        <div className="flex items-center justify-between px-3 py-2 rounded-lg bg-emerald-950/40 border border-emerald-800/50 text-[11px]">
                          <span className="text-emerald-300 flex items-center gap-1.5">
                            <FileText className="w-3.5 h-3.5 text-emerald-400" />
                            Digital PDF Signed & Ready
                          </span>
                          <span className="text-emerald-400 font-mono font-bold">Turnaround: 38 mins</span>
                        </div>
                      </motion.div>
                    )}

                    {activeTab === "patient" && (
                      <motion.div
                        key="patient"
                        initial={{ opacity: 0, y: 8 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -8 }}
                        transition={{ duration: 0.2 }}
                        className="space-y-3"
                      >
                        <div className="flex items-center justify-between text-xs text-slate-400 pb-2 border-b border-slate-800">
                          <span className="font-semibold text-white">Patient Unified Health Vault</span>
                          <span className="font-mono text-[11px] text-teal-400">Encrypted Delivery</span>
                        </div>

                        <div className="p-3 rounded-xl bg-slate-950/70 border border-slate-800 space-y-2">
                          <div className="flex items-center justify-between text-xs">
                            <span className="text-slate-400">Latest Diagnostic Record</span>
                            <span className="font-mono text-teal-300 font-bold">Today, 11:42 AM</span>
                          </div>
                          <p className="text-xs font-bold text-white">Complete Metabolic Panel (CMP)</p>
                          <div className="flex items-center gap-2 pt-1">
                            <span className="px-2 py-0.5 rounded bg-emerald-500/20 text-emerald-300 text-[10px] font-bold">WhatsApp PDF Sent</span>
                            <span className="px-2 py-0.5 rounded bg-teal-500/20 text-teal-300 text-[10px] font-bold">Clinic Chart Updated</span>
                          </div>
                        </div>

                        <div className="flex items-center justify-between px-3 py-2 rounded-lg bg-teal-950/40 border border-teal-800/50 text-[11px]">
                          <span className="text-teal-300 flex items-center gap-1.5">
                            <Pill className="w-3.5 h-3.5 text-teal-400" />
                            E-Prescription Routed to Pharmacy
                          </span>
                          <span className="text-teal-300 font-mono font-bold">Zero Paper</span>
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
                      <span>FHIR HL7 Endpoints Live</span>
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
