import { useState } from "react";
import { Sparkles, ArrowRight, ShieldCheck } from "lucide-react";
import { WaitlistModal } from "@/components/waitlist-modal";

export function FinalEditorialCTASection() {
  const [waitlistOpen, setWaitlistOpen] = useState(false);

  return (
    <>
      <section id="final-cta" className="py-14 sm:py-20 bg-slate-900 dark:bg-slate-950 text-white relative overflow-hidden">
        {/* Radial Glow */}
        <div
          className="absolute top-0 left-1/2 -translate-x-1/2 w-[700px] h-[250px] pointer-events-none"
          style={{
            background: "radial-gradient(ellipse at center, rgba(15,118,110,0.18) 0%, transparent 70%)",
          }}
        />

        <div className="max-w-[1280px] mx-auto px-4 sm:px-6 relative z-10 text-center">
          
          <div className="max-w-3xl mx-auto space-y-6">
            
            {/* Main Headline */}
            <h2 className="text-3xl sm:text-5xl font-black text-white tracking-tight leading-tight">
              Healthcare works better <br />
              <span className="bg-clip-text text-transparent bg-gradient-to-r from-teal-300 via-emerald-300 to-sky-400">
                when it works together.
              </span>
            </h2>

            {/* Supporting Paragraph */}
            <p className="text-sm sm:text-lg text-slate-300 font-normal max-w-xl mx-auto leading-relaxed">
              Join clinics, laboratories, pharmacies, and patients pioneering Africa's connected healthcare operating network.
            </p>

            {/* Primary Action Button */}
            <div className="pt-2 flex flex-col sm:flex-row items-center justify-center gap-3">
              <button
                onClick={() => setWaitlistOpen(true)}
                className="group flex items-center justify-center gap-2.5 px-6 py-3.5 rounded-xl bg-[#0F766E] hover:bg-[#115E59] text-white text-sm font-bold transition-all shadow-xl cursor-pointer border-0 w-full sm:w-auto"
              >
                <Sparkles className="w-4 h-4 text-teal-200" />
                <span>Join Early Access</span>
                <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
              </button>
            </div>

            {/* Security Guarantee */}
            <div className="flex items-center justify-center gap-1.5 text-[10px] sm:text-xs font-medium text-slate-400 pt-3">
              <ShieldCheck className="w-3.5 h-3.5 text-teal-400" />
              <span>Multi-tenant Isolation • Zero personal data exposed • Verified Supabase Database</span>
            </div>

          </div>

        </div>
      </section>

      {/* Reused Waitlist Modal */}
      <WaitlistModal open={waitlistOpen} onOpenChange={setWaitlistOpen} />
    </>
  );
}
