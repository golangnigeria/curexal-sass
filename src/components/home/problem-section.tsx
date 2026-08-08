import React from "react";
import { User, Stethoscope, FlaskConical, Binary, Pill, RefreshCcw, XCircle, CheckCircle2, ArrowRight, Network } from "lucide-react";

export function ProblemSection() {
  const steps = [
    { icon: User, label: "1. Patient", desc: "Seeks care for symptoms" },
    { icon: Stethoscope, label: "2. Clinic", desc: "Consultation & orders" },
    { icon: FlaskConical, label: "3. Laboratory", desc: "Bloodwork & pathology" },
    { icon: Binary, label: "4. Imaging", desc: "Radiology & scans" },
    { icon: Pill, label: "5. Pharmacy", desc: "Medication fulfillment" },
    { icon: RefreshCcw, label: "6. Follow-up", desc: "Care plan review" },
  ];

  return (
    <section id="problem" className="py-20 bg-white dark:bg-[#0B1120] border-b border-slate-100 dark:border-slate-800">
      <div className="max-w-[1280px] mx-auto px-4 sm:px-6">
        
        {/* Section Title */}
        <div className="text-center max-w-3xl mx-auto mb-14">
          <div className="inline-flex items-center gap-1.5 px-3 py-1 mb-3 rounded-full border border-rose-500/20 bg-rose-50 dark:bg-rose-950/40 text-rose-600 dark:text-rose-400 text-xs font-bold uppercase tracking-wider">
            Healthcare Fragmentation
          </div>
          <h2 className="text-2xl sm:text-4xl font-extrabold text-slate-900 dark:text-white tracking-tight mb-4">
            Healthcare doesn't happen inside one building.
          </h2>
          <p className="text-sm sm:text-base text-slate-600 dark:text-slate-400">
            A single patient healthcare episode routinely crosses multiple independent organizations.
          </p>
        </div>

        {/* Journey Flow Ribbon */}
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3 mb-16">
          {steps.map((s, idx) => {
            const Icon = s.icon;
            return (
              <div
                key={s.label}
                className="relative p-4 rounded-2xl bg-slate-50 dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 flex flex-col items-center text-center space-y-2"
              >
                <div className="w-10 h-10 rounded-xl bg-teal-50 dark:bg-teal-950 text-[#0F766E] dark:text-teal-400 flex items-center justify-center font-bold">
                  <Icon className="w-5 h-5" />
                </div>
                <h3 className="text-xs font-bold text-slate-900 dark:text-white">{s.label}</h3>
                <p className="text-[11px] text-slate-500 dark:text-slate-400 leading-tight">{s.desc}</p>
                {idx < steps.length - 1 && (
                  <ArrowRight className="hidden lg:block absolute -right-3 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-300 dark:text-slate-700 z-10" />
                )}
              </div>
            );
          })}
        </div>

        {/* Contrast Matrix: Traditional Silos vs Curexal Coordination */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 items-stretch">
          
          {/* Traditional Fragmented Model */}
          <div className="p-6 sm:p-8 rounded-3xl bg-rose-50/50 dark:bg-rose-950/20 border border-rose-200/80 dark:border-rose-900/50 space-y-6">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-xl bg-rose-500 text-white flex items-center justify-center font-bold text-lg">
                ❌
              </div>
              <div>
                <h3 className="text-lg font-bold text-slate-900 dark:text-white">Traditional Fragmented Model</h3>
                <p className="text-xs text-rose-600 dark:text-rose-400 font-semibold">Disconnected systems create friction & delay</p>
              </div>
            </div>

            <div className="space-y-3 pt-2">
              {[
                "Paper referrals hand-carried by patients",
                "Repeated patient visits due to missing records",
                "Lost laboratory information and lost diagnostic history",
                "Disconnected software systems that cannot talk to each other",
                "Zero real-time visibility into referral status",
                "Manual monthly partner commission & billing settlements",
              ].map((item) => (
                <div key={item} className="flex items-start gap-2.5 text-xs text-slate-700 dark:text-slate-300">
                  <XCircle className="w-4 h-4 text-rose-500 flex-shrink-0 mt-0.5" />
                  <span>{item}</span>
                </div>
              ))}
            </div>
          </div>

          {/* Curexal Connected Network */}
          <div className="p-6 sm:p-8 rounded-3xl bg-teal-50/60 dark:bg-teal-950/30 border border-teal-200/80 dark:border-teal-800/80 space-y-6">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-xl bg-[#0F766E] text-white flex items-center justify-center font-bold">
                <Network className="w-6 h-6 text-white" />
              </div>
              <div>
                <h3 className="text-lg font-bold text-slate-900 dark:text-white">The Curexal Coordination Model</h3>
                <p className="text-xs text-[#0F766E] dark:text-teal-400 font-semibold">One connected operating network</p>
              </div>
            </div>

            <div className="space-y-3 pt-2">
              {[
                "Digital referrals sent directly from clinic chart to laboratory",
                "Instant electronic test result dispatch to referring doctor",
                "Patient medical history unified across participating facilities",
                "End-to-end referral status tracking from collection to sign-off",
                "Automated B2B transactions and partner order coordination",
                "Independent tenant data boundaries with secure cross-org sharing",
              ].map((item) => (
                <div key={item} className="flex items-start gap-2.5 text-xs text-slate-800 dark:text-slate-200 font-medium">
                  <CheckCircle2 className="w-4 h-4 text-[#0F766E] dark:text-teal-400 flex-shrink-0 mt-0.5" />
                  <span>{item}</span>
                </div>
              ))}
            </div>
          </div>

        </div>

      </div>
    </section>
  );
}
