import React from "react";
import { Stethoscope, FlaskConical, Pill, Building2, User, ArrowRight, ShieldCheck, Zap } from "lucide-react";

export function CurexalWaySection() {
  const participants = [
    { name: "Clinic A", type: "Outpatient Clinic", desc: "Creates digital referrals & lab requisitions" },
    { name: "Laboratory", type: "Diagnostic Facility", desc: "Processes specimens & verifies digital results" },
    { name: "Pharmacy", type: "Retail & Hospital Pharmacy", desc: "Fulfills prescriptions & diagnostic supplies" },
    { name: "Clinic B", type: "Specialist Center", desc: "Receives transferred results for care decisions" },
    { name: "Patient", type: "Individual Care Recipient", desc: "Tracks results & bookings without carrying paper" },
  ];

  return (
    <section className="py-20 bg-slate-900 text-white relative overflow-hidden">
      {/* Background Glow */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[700px] h-[400px] bg-teal-500/10 blur-[120px] pointer-events-none" />

      <div className="max-w-[1280px] mx-auto px-4 sm:px-6 relative z-10">
        
        {/* Header */}
        <div className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="text-3xl sm:text-5xl font-black tracking-tight text-white mb-4">
            One network. <br className="hidden sm:inline" />
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-teal-400 to-emerald-300">
              Many healthcare organizations.
            </span>
          </h2>
          <p className="text-slate-300 text-base leading-relaxed">
            Curexal does not require every healthcare provider to merge into one corporation. Instead, independent organizations retain their ownership, branding, and data boundaries while participating in a connected healthcare network.
          </p>
        </div>

        {/* Network Hub Visualization */}
        <div className="bg-slate-950/80 border border-slate-800 rounded-3xl p-6 sm:p-10 shadow-2xl">
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-center">
            
            {/* Left: Independent Participants List */}
            <div className="lg:col-span-5 space-y-3">
              <h3 className="text-xs font-bold uppercase tracking-wider text-slate-400 mb-2">
                Independent Healthcare Participants
              </h3>

              {participants.map((item) => (
                <div
                  key={item.name}
                  className="p-3.5 rounded-2xl bg-slate-900/90 border border-slate-800 flex items-center justify-between hover:border-teal-500/40 transition-colors"
                >
                  <div>
                    <h4 className="text-sm font-bold text-white">{item.name}</h4>
                    <p className="text-[11px] text-slate-400">{item.desc}</p>
                  </div>
                  <span className="text-[10px] font-bold text-teal-400 bg-teal-950/80 px-2 py-0.5 rounded-full border border-teal-800 flex-shrink-0">
                    Independent
                  </span>
                </div>
              ))}
            </div>

            {/* Middle: Flow Arrow */}
            <div className="lg:col-span-2 flex flex-col items-center justify-center py-4 lg:py-0">
              <div className="w-12 h-12 rounded-full bg-teal-500/20 border border-teal-500/40 flex items-center justify-center text-teal-400 animate-pulse">
                <ArrowRight className="w-6 h-6 rotate-90 lg:rotate-0" />
              </div>
              <span className="text-[11px] font-bold text-teal-300 mt-2 text-center">Controlled Data Sharing</span>
            </div>

            {/* Right: Curexal Operating Network Core */}
            <div className="lg:col-span-5">
              <div className="p-8 rounded-3xl bg-gradient-to-b from-[#0F766E] to-[#115E59] border border-teal-400/30 text-center space-y-4 shadow-xl">
                <div className="w-16 h-16 rounded-2xl bg-white/10 mx-auto flex items-center justify-center">
                  <Zap className="w-8 h-8 text-teal-200" />
                </div>
                <h3 className="text-xl font-black tracking-tight text-white uppercase">CUREXAL NETWORK</h3>
                <p className="text-xs text-teal-100 leading-relaxed">
                  Coordinates patient movements, referral workflows, diagnostic sample tracking, verified result delivery, and payment settlements across facility boundaries.
                </p>

                <div className="pt-2 flex items-center justify-center gap-2 text-xs font-bold text-teal-200 bg-white/10 py-2 rounded-xl border border-white/10">
                  <ShieldCheck className="w-4 h-4 text-teal-300" />
                  <span>Tenant Isolation Architecture</span>
                </div>
              </div>
            </div>

          </div>
        </div>

        {/* Emphasized Statement */}
        <div className="mt-12 text-center">
          <div className="inline-block px-6 py-3 rounded-full bg-teal-500/10 border border-teal-500/30 text-teal-300 font-extrabold text-sm sm:text-base tracking-wide">
            INDEPENDENT ORGANIZATIONS. COORDINATED HEALTHCARE.
          </div>
        </div>

      </div>
    </section>
  );
}
