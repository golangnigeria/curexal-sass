import { useState } from "react";
import { Link } from "react-router-dom";
import { ArrowRight, Sparkles, ShieldCheck, Stethoscope, FlaskConical, Pill, User, Network, RefreshCw, ShoppingBag, CheckCircle2 } from "lucide-react";
import { motion } from "framer-motion";
import { WaitlistModal } from "@/components/waitlist-modal";

export function Hero() {
  const [waitlistOpen, setWaitlistOpen] = useState(false);

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
                  <Sparkles className="h-4 w-4 text-teal-200" />
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

            {/* Right: Living Healthcare Network Visual (Framer Motion Animated) */}
            <div className="lg:col-span-6 flex items-center justify-center mt-4 lg:mt-0 w-full">
              <motion.div
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ duration: 0.7, ease: "easeOut" }}
                className="relative w-full max-w-[340px] xs:max-w-[380px] sm:max-w-[440px] aspect-square flex items-center justify-center mx-auto box-border"
              >
                
                {/* 1. Rotating SVG Network Dash Rings */}
                <motion.svg
                  animate={{ rotate: 360 }}
                  transition={{ duration: 40, repeat: Infinity, ease: "linear" }}
                  className="absolute inset-0 w-full h-full pointer-events-none opacity-30"
                  viewBox="0 0 400 400"
                  fill="none"
                >
                  <circle cx="200" cy="200" r="145" stroke="#0F766E" strokeWidth="1.5" strokeDasharray="8 8" />
                  <circle cx="200" cy="200" r="95" stroke="#0D9488" strokeWidth="1" strokeDasharray="4 4" />
                </motion.svg>

                {/* 2. Flowing Animated Connection Rays */}
                <svg className="absolute inset-0 w-full h-full pointer-events-none z-0" viewBox="0 0 400 400" fill="none">
                  {/* Connection lines */}
                  <line x1="200" y1="55" x2="200" y2="150" stroke="#0F766E" strokeWidth="2" strokeDasharray="4 4" className="opacity-40" />
                  <line x1="330" y1="325" x2="240" y2="240" stroke="#0D9488" strokeWidth="2" strokeDasharray="4 4" className="opacity-40" />
                  <line x1="70" y1="325" x2="160" y2="240" stroke="#14B8A6" strokeWidth="2" strokeDasharray="4 4" className="opacity-40" />
                  <line x1="55" y1="175" x2="150" y2="190" stroke="#2DD4BF" strokeWidth="2" strokeDasharray="4 4" className="opacity-40" />
                  <line x1="345" y1="175" x2="250" y2="190" stroke="#0F766E" strokeWidth="2" strokeDasharray="4 4" className="opacity-40" />

                  {/* Flowing Pulse Particles */}
                  <motion.circle
                    r="3.5"
                    fill="#14B8A6"
                    animate={{ cx: [200, 200], cy: [55, 150] }}
                    transition={{ duration: 2.2, repeat: Infinity, ease: "linear" }}
                  />
                  <motion.circle
                    r="3.5"
                    fill="#0D9488"
                    animate={{ cx: [330, 240], cy: [325, 240] }}
                    transition={{ duration: 2.6, repeat: Infinity, ease: "linear", delay: 0.4 }}
                  />
                  <motion.circle
                    r="3.5"
                    fill="#2DD4BF"
                    animate={{ cx: [70, 160], cy: [325, 240] }}
                    transition={{ duration: 2.4, repeat: Infinity, ease: "linear", delay: 0.8 }}
                  />
                  <motion.circle
                    r="3.5"
                    fill="#0F766E"
                    animate={{ cx: [55, 150], cy: [175, 190] }}
                    transition={{ duration: 2.1, repeat: Infinity, ease: "linear", delay: 0.2 }}
                  />
                  <motion.circle
                    r="3.5"
                    fill="#14B8A6"
                    animate={{ cx: [345, 250], cy: [175, 190] }}
                    transition={{ duration: 2.5, repeat: Infinity, ease: "linear", delay: 0.6 }}
                  />
                </svg>

                {/* 3. Central Hub: CUREXAL OPERATING NETWORK (Circular Hub) */}
                <motion.div
                  animate={{ y: [-3, 3, -3] }}
                  transition={{ duration: 4, repeat: Infinity, ease: "easeInOut" }}
                  className="relative z-20 flex flex-col items-center justify-center p-4 rounded-full bg-gradient-to-b from-[#0F766E] to-[#115E59] text-white border border-teal-400/40 text-center w-32 h-32 sm:w-36 sm:h-36 group cursor-pointer shadow-lg"
                >
                  {/* Glowing Radar Circle Ring Behind Hub */}
                  <motion.div
                    animate={{ scale: [1, 1.3, 1], opacity: [0.6, 0, 0.6] }}
                    transition={{ duration: 3, repeat: Infinity, ease: "easeInOut" }}
                    className="absolute inset-0 rounded-full bg-teal-500/25 pointer-events-none"
                  />

                  <div className="w-9 h-9 sm:w-10 sm:h-10 rounded-full bg-white/20 flex items-center justify-center mb-1">
                    <Network className="w-5 h-5 sm:w-6 sm:h-6 text-white" />
                  </div>
                  <span className="text-[11px] sm:text-xs font-black uppercase tracking-wider text-white">CUREXAL</span>
                  <span className="text-[8px] sm:text-[9px] text-teal-200 font-bold tracking-tight">OPERATING NETWORK</span>
                </motion.div>

                {/* 4. Floating Satellite Nodes */}

                {/* Node 1: Clinic (Top Center) */}
                <motion.div
                  animate={{ y: [-4, 4, -4] }}
                  transition={{ duration: 3.8, repeat: Infinity, ease: "easeInOut" }}
                  whileHover={{ scale: 1.08 }}
                  className="absolute top-[3%] left-1/2 -translate-x-1/2 z-10 flex flex-col items-center cursor-pointer"
                >
                  <div className="px-2.5 sm:px-3 py-1.5 rounded-xl bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 flex items-center gap-1.5">
                    <div className="w-6 h-6 rounded-lg bg-teal-50 dark:bg-teal-950 flex items-center justify-center text-[#0F766E] dark:text-teal-400">
                      <Stethoscope className="w-3.5 h-3.5" />
                    </div>
                    <span className="text-[11px] font-extrabold text-slate-800 dark:text-white">CLINIC</span>
                  </div>
                  <span className="mt-1 text-[8px] sm:text-[9px] font-bold text-[#0F766E] dark:text-teal-300 bg-teal-50 dark:bg-teal-950 px-2 py-0.5 rounded-full border border-teal-200/80 dark:border-teal-800/80">
                    Digital Referral
                  </span>
                </motion.div>

                {/* Node 2: Laboratory (Bottom Right) */}
                <motion.div
                  animate={{ y: [4, -4, 4] }}
                  transition={{ duration: 4.2, repeat: Infinity, ease: "easeInOut", delay: 0.5 }}
                  whileHover={{ scale: 1.08 }}
                  className="absolute bottom-[4%] right-[2%] z-10 flex flex-col items-end cursor-pointer"
                >
                  <div className="px-2.5 sm:px-3 py-1.5 rounded-xl bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 flex items-center gap-1.5">
                    <div className="w-6 h-6 rounded-lg bg-teal-50 dark:bg-teal-950 flex items-center justify-center text-[#0F766E] dark:text-teal-400">
                      <FlaskConical className="w-3.5 h-3.5" />
                    </div>
                    <span className="text-[11px] font-extrabold text-slate-800 dark:text-white">LABORATORY</span>
                  </div>
                  <span className="mt-1 text-[8px] sm:text-[9px] font-bold text-[#0F766E] dark:text-teal-300 bg-teal-50 dark:bg-teal-950 px-2 py-0.5 rounded-full border border-teal-200/80 dark:border-teal-800/80">
                    Verified Results
                  </span>
                </motion.div>

                {/* Node 3: Pharmacy (Bottom Left) */}
                <motion.div
                  animate={{ y: [-4, 4, -4] }}
                  transition={{ duration: 3.6, repeat: Infinity, ease: "easeInOut", delay: 1 }}
                  whileHover={{ scale: 1.08 }}
                  className="absolute bottom-[4%] left-[2%] z-10 flex flex-col items-start cursor-pointer"
                >
                  <div className="px-2.5 sm:px-3 py-1.5 rounded-xl bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 flex items-center gap-1.5">
                    <div className="w-6 h-6 rounded-lg bg-teal-50 dark:bg-teal-950 flex items-center justify-center text-[#0F766E] dark:text-teal-400">
                      <Pill className="w-3.5 h-3.5" />
                    </div>
                    <span className="text-[11px] font-extrabold text-slate-800 dark:text-white">PHARMACY</span>
                  </div>
                  <span className="mt-1 text-[8px] sm:text-[9px] font-bold text-[#0F766E] dark:text-teal-300 bg-teal-50 dark:bg-teal-950 px-2 py-0.5 rounded-full border border-teal-200/80 dark:border-teal-800/80">
                    Prescriptions
                  </span>
                </motion.div>

                {/* Node 4: Patient (Mid Left) */}
                <motion.div
                  animate={{ y: [4, -4, 4] }}
                  transition={{ duration: 4, repeat: Infinity, ease: "easeInOut", delay: 0.3 }}
                  whileHover={{ scale: 1.08 }}
                  className="absolute top-[34%] left-[1%] z-10 flex flex-col items-start cursor-pointer"
                >
                  <div className="px-2.5 sm:px-3 py-1.5 rounded-xl bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 flex items-center gap-1.5">
                    <div className="w-6 h-6 rounded-lg bg-teal-50 dark:bg-teal-950 flex items-center justify-center text-[#0F766E] dark:text-teal-400">
                      <User className="w-3.5 h-3.5" />
                    </div>
                    <span className="text-[11px] font-extrabold text-slate-800 dark:text-white">PATIENT</span>
                  </div>
                  <span className="mt-1 text-[8px] sm:text-[9px] font-bold text-[#0F766E] dark:text-teal-300 bg-teal-50 dark:bg-teal-950 px-2 py-0.5 rounded-full border border-teal-200/80 dark:border-teal-800/80">
                    Care Records
                  </span>
                </motion.div>

                {/* Node 5: Suppliers (Mid Right) */}
                <motion.div
                  animate={{ y: [-4, 4, -4] }}
                  transition={{ duration: 4.4, repeat: Infinity, ease: "easeInOut", delay: 0.7 }}
                  whileHover={{ scale: 1.08 }}
                  className="absolute top-[34%] right-[1%] z-10 flex flex-col items-end cursor-pointer"
                >
                  <div className="px-2.5 sm:px-3 py-1.5 rounded-xl bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 flex items-center gap-1.5">
                    <div className="w-6 h-6 rounded-lg bg-teal-50 dark:bg-teal-950 flex items-center justify-center text-[#0F766E] dark:text-teal-400">
                      <ShoppingBag className="w-3.5 h-3.5" />
                    </div>
                    <span className="text-[11px] font-extrabold text-slate-800 dark:text-white">SUPPLIERS</span>
                  </div>
                  <span className="mt-1 text-[8px] sm:text-[9px] font-bold text-[#0F766E] dark:text-teal-300 bg-teal-50 dark:bg-teal-950 px-2 py-0.5 rounded-full border border-teal-200/80 dark:border-teal-800/80">
                    B2B Orders
                  </span>
                </motion.div>

              </motion.div>
            </div>

          </div>
        </div>
      </section>

      {/* Waitlist Modal */}
      <WaitlistModal open={waitlistOpen} onOpenChange={setWaitlistOpen} />
    </>
  );
}
