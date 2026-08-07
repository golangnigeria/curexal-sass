import { useState } from "react";
import { Link } from "react-router-dom";
import { ArrowRight, Sparkles, FlaskConical, Activity, ShieldCheck, Microscope, HeartPulse, Zap } from "lucide-react";
import { WaitlistModal } from "@/components/waitlist-modal";

const stats = [
  { value: "ISO 15189", label: "Accreditation Ready" },
  { value: "100%", label: "Tenant Data Isolation" },
  { value: "99.98%", label: "Uptime SLA Target" },
  { value: "HL7 / FHIR", label: "Interoperability Ready" },
];

const orbitNodes = [
  { icon: FlaskConical, label: "Laboratories", color: "#0F766E", delay: "0s" },
  { icon: Microscope, label: "Diagnostics", color: "#0D9488", delay: "-5s" },
  { icon: HeartPulse, label: "Patients", color: "#14B8A6", delay: "-10s" },
  { icon: ShieldCheck, label: "Compliance", color: "#2DD4BF", delay: "-15s" },
];

export function Hero() {
  const [waitlistOpen, setWaitlistOpen] = useState(false);

  return (
    <>
      <section
        id="hero"
        className="relative z-0 isolate pt-14 sm:pt-20 lg:pt-24 min-h-screen flex items-center overflow-hidden bg-white dark:bg-[#0B1120]"
      >
        {/* Ambient background effects */}
        <div className="absolute inset-0 dot-grid opacity-30 pointer-events-none" />

        {/* Top-center teal glow */}
        <div
          className="absolute top-0 left-1/2 -translate-x-1/2 w-[900px] h-[350px] pointer-events-none"
          style={{
            background: "radial-gradient(ellipse at center, rgba(15,118,110,0.07) 0%, transparent 70%)",
          }}
        />

        {/* Bottom-right accent glow */}
        <div
          className="absolute bottom-0 right-0 w-[500px] h-[350px] pointer-events-none"
          style={{
            background: "radial-gradient(ellipse at 70% 100%, rgba(15,118,110,0.04) 0%, transparent 70%)",
          }}
        />

        <div className="relative z-10 max-w-[1280px] mx-auto px-4 sm:px-6 py-8 sm:py-16 lg:py-24 w-full">
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-16 items-center">

            {/* Left — Content */}
            <div className="lg:col-span-6 xl:col-span-6 flex flex-col items-start text-left">

              {/* Eyebrow */}
              <div className="inline-flex items-center gap-1.5 px-2.5 py-1 mb-4 sm:mb-6 rounded-full border border-[#0F766E]/20 bg-[#F0FDFA] dark:bg-[#0F766E]/10 animate-fade-up">
                <span className="w-1.5 h-1.5 rounded-full bg-[#0F766E] animate-pulse" />
                <span className="text-[11px] sm:text-xs font-semibold text-[#0F766E] tracking-wide">
                  Healthcare Network Infrastructure
                </span>
              </div>

              {/* Headline */}
              <h1 className="text-2xl xs:text-3xl sm:text-[44px] lg:text-[52px] font-extrabold leading-[1.12] sm:leading-[1.08] tracking-tight text-gray-900 dark:text-white mb-4 sm:mb-6 animate-fade-up delay-100">
                Building the Digital Backbone of Healthcare{" "}
                <span
                  className="bg-clip-text text-transparent animate-gradient"
                  style={{
                    backgroundImage: "linear-gradient(135deg, #0F766E 0%, #0D9488 30%, #14B8A6 60%, #0F766E 100%)",
                    backgroundSize: "200% 200%",
                  }}
                >
                  Across Africa
                </span>
              </h1>

              {/* Subheadline — network positioning */}
              <p className="text-sm sm:text-[16px] lg:text-[17px] leading-relaxed text-gray-500 dark:text-gray-400 max-w-lg mb-6 sm:mb-8 animate-fade-up delay-200">
                Curexal is a healthcare operating network that connects laboratories, clinics, pharmacies, diagnostic centers, and patients through one secure platform.
              </p>

              {/* CTAs: Book Demo + Join Waitlist */}
              <div className="flex flex-row flex-wrap items-center gap-2.5 sm:gap-3 w-full sm:w-auto animate-fade-up delay-300">
                <Link to="/book-demo" id="hero-book-demo" className="flex-1 sm:flex-initial">
                  <button
                    className="w-full sm:w-auto group flex items-center justify-center gap-1.5 sm:gap-2 px-4 sm:px-6 py-2.5 sm:py-3 rounded-xl bg-[#0F766E] hover:bg-[#115E59] text-white text-xs sm:text-sm font-bold transition-all cursor-pointer border-0 shadow-sm"
                  >
                    <span>Book a Demo</span>
                    <ArrowRight className="h-3.5 w-3.5 group-hover:translate-x-0.5 transition-transform" />
                  </button>
                </Link>

                <Link to="/waitlist" id="hero-waitlist" className="flex-1 sm:flex-initial">
                  <button
                    className="w-full sm:w-auto flex items-center justify-center gap-1.5 sm:gap-2 px-4 sm:px-6 py-2.5 sm:py-3 rounded-xl border border-teal-200 dark:border-teal-800/80 bg-teal-50 dark:bg-teal-950/40 text-[#0F766E] dark:text-teal-300 text-xs sm:text-sm font-bold hover:bg-teal-100 dark:hover:bg-teal-900/60 transition-all cursor-pointer shadow-xs"
                  >
                    <Sparkles className="h-3.5 w-3.5 text-[#0F766E]" />
                    <span>Join Waitlist</span>
                  </button>
                </Link>
              </div>

              {/* Trust signals */}
              <div className="flex flex-wrap items-center gap-2.5 sm:gap-4 mt-6 sm:mt-8 animate-fade-up delay-400">
                <div className="flex items-center gap-1 text-[11px] sm:text-xs text-gray-400 dark:text-gray-500">
                  <ShieldCheck className="w-3.5 h-3.5 text-[#0F766E]" />
                  <span>HIPAA Aligned</span>
                </div>
                <div className="w-px h-3 bg-gray-200 dark:bg-[#374151]" />
                <div className="flex items-center gap-1 text-[11px] sm:text-xs text-gray-400 dark:text-gray-500">
                  <ShieldCheck className="w-3.5 h-3.5 text-[#0F766E]" />
                  <span>ISO 15189</span>
                </div>
                <div className="w-px h-3 bg-gray-200 dark:bg-[#374151]" />
                <div className="flex items-center gap-1 text-[11px] sm:text-xs text-gray-400 dark:text-gray-500">
                  <ShieldCheck className="w-3.5 h-3.5 text-[#0F766E]" />
                  <span>NDPR Privacy</span>
                </div>
              </div>

            </div>

            {/* Right — Responsive Animated Network Visual */}
            <div className="lg:col-span-6 xl:col-span-6 flex items-center justify-center animate-fade-up delay-400 mt-4 lg:mt-0">
              <div className="relative w-[260px] h-[260px] xs:w-[300px] xs:h-[300px] sm:w-[380px] sm:h-[380px] lg:w-[440px] lg:h-[440px]">

                {/* Outer orbit ring */}
                <div className="absolute inset-0 rounded-full border border-dashed border-[#0F766E]/15 dark:border-[#0F766E]/20" />

                {/* Middle orbit ring */}
                <div className="absolute inset-6 xs:inset-8 sm:inset-10 rounded-full border border-[#0F766E]/10 dark:border-[#0F766E]/15" />

                {/* Inner orbit ring */}
                <div className="absolute inset-12 xs:inset-16 sm:inset-20 rounded-full border border-[#0F766E]/8 animate-pulse-ring" />

                {/* Orbiting nodes */}
                <div className="absolute inset-0 animate-orbit" style={{ animationDuration: "24s" }}>
                  {orbitNodes.map((node, i) => {
                    const angle = (i * 90) * (Math.PI / 180);
                    const radius = 44;
                    const x = 50 + radius * Math.cos(angle);
                    const y = 50 + radius * Math.sin(angle);
                    return (
                      <div
                        key={node.label}
                        className="absolute flex flex-col items-center gap-1"
                        style={{
                          left: `${x}%`,
                          top: `${y}%`,
                          transform: "translate(-50%, -50%)",
                        }}
                      >
                        <div
                          className="w-8 h-8 sm:w-11 sm:h-11 rounded-[10px] sm:rounded-[12px] bg-white dark:bg-[#1F2937] border border-gray-200 dark:border-[#374151] flex items-center justify-center shadow-sm"
                          style={{
                            animation: `orbit 24s linear infinite reverse`,
                          }}
                        >
                          <node.icon className="w-4 h-4 sm:w-5 sm:h-5" style={{ color: node.color }} />
                        </div>
                      </div>
                    );
                  })}
                </div>

                {/* Floating micro-cards */}
                <div
                  className="absolute top-2 right-1 sm:top-6 sm:right-4 px-2.5 py-1.5 sm:px-3 sm:py-2 rounded-[8px] sm:rounded-[10px] bg-white dark:bg-[#1F2937] border border-gray-200 dark:border-[#374151] shadow-sm z-10"
                >
                  <div className="flex items-center gap-1.5">
                    <div className="w-1.5 h-1.5 sm:w-2 sm:h-2 rounded-full bg-emerald-500" />
                    <span className="text-[10px] sm:text-[11px] font-semibold text-gray-700 dark:text-gray-200">Referral Sent</span>
                  </div>
                </div>

                <div
                  className="absolute bottom-2 left-0 sm:bottom-6 sm:left-2 px-2.5 py-1.5 sm:px-3 sm:py-2 rounded-[8px] sm:rounded-[10px] bg-white dark:bg-[#1F2937] border border-gray-200 dark:border-[#374151] shadow-sm z-10"
                >
                  <div className="flex items-center gap-1.5">
                    <Zap className="w-3 h-3 text-amber-500" />
                    <span className="text-[10px] sm:text-[11px] font-semibold text-gray-700 dark:text-gray-200">Result Delivered</span>
                  </div>
                </div>

                {/* Center core */}
                <div className="absolute inset-0 flex items-center justify-center">
                  <div className="relative">
                    <div
                      className="w-14 h-14 sm:w-20 sm:h-20 rounded-full flex items-center justify-center bg-gradient-to-tr from-[#0F766E] to-[#14B8A6] shadow-lg"
                    >
                      <Activity className="w-6 h-6 sm:w-9 sm:h-9 text-white" />
                    </div>
                  </div>
                </div>

              </div>
            </div>

          </div>

          {/* Stats bar */}
          <div className="mt-10 sm:mt-16 pt-6 sm:pt-10 border-t border-gray-100 dark:border-[#1F2937] grid grid-cols-2 md:grid-cols-4 gap-4 sm:gap-8 text-left">
            {stats.map((s, i) => (
              <div key={s.label} className={`flex flex-col gap-0.5 animate-count-up delay-${(i + 5) * 100}`}>
                <span className="text-xl sm:text-[28px] font-bold text-gray-900 dark:text-white tracking-tight">{s.value}</span>
                <span className="text-xs sm:text-sm text-gray-500 dark:text-gray-400">{s.label}</span>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Waitlist Modal */}
      <WaitlistModal open={waitlistOpen} onOpenChange={setWaitlistOpen} />
    </>
  );
}
