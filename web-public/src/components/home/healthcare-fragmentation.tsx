import { motion } from "framer-motion";
import { User, Stethoscope, Phone, FlaskConical, FileText, Pill, CheckCircle2, AlertTriangle, ArrowRight, ArrowDown } from "lucide-react";

const steps = [
  {
    step: "01",
    actor: "Patient",
    icon: User,
    action: "Visits clinic with symptoms",
    friction: "Manually shares medical history verbally",
  },
  {
    step: "02",
    actor: "Clinic",
    icon: Stethoscope,
    action: "Issues paper lab requisition",
    friction: "No digital link to diagnostic center",
  },
  {
    step: "03",
    actor: "Paper / Phone",
    icon: Phone,
    action: "Patient carries paper slip",
    friction: "High risk of lost slip or misread handwriting",
  },
  {
    step: "04",
    actor: "Laboratory",
    icon: FlaskConical,
    action: "Performs diagnostic tests",
    friction: "Manual result entry and delayed printouts",
  },
  {
    step: "05",
    actor: "Result Dispatch",
    icon: FileText,
    action: "Patient returns for paper result",
    friction: "Extra travel costs and repeated clinic trips",
  },
  {
    step: "06",
    actor: "Pharmacy",
    icon: Pill,
    action: "Dispenses prescribed medication",
    friction: "No access to lab values or doctor notes",
  },
  {
    step: "07",
    actor: "Patient Care",
    icon: CheckCircle2,
    action: "Treatment initiated after delays",
    friction: "Fragmented health record left behind",
  },
];

export function HealthcareFragmentationSection() {
  return (
    <section id="fragmentation" className="py-10 sm:py-16 bg-white dark:bg-[#0B1120] text-slate-900 dark:text-white relative overflow-hidden">
      <div className="max-w-[1280px] mx-auto px-4 sm:px-6 relative z-10">
        
        {/* Editorial Header */}
        <div className="max-w-2xl mb-8 sm:mb-12">
          <h2 className="text-2xl sm:text-4xl font-black tracking-tight text-slate-900 dark:text-white leading-tight mb-3">
            Healthcare is connected in real life. <br />
            <span className="bg-clip-text text-transparent bg-gradient-to-r from-rose-500 via-amber-500 to-[#0F766E]">
              But disconnected digitally.
            </span>
          </h2>

          <p className="text-xs sm:text-sm text-slate-600 dark:text-slate-300 leading-relaxed">
            When clinics, laboratories, pharmacies, and patients operate in isolated silos, healthcare relies on paper slips, phone calls, and manual footwork. Curexal bridges these gaps with real-time digital coordination.
          </p>
        </div>

        {/* Mobile Vertical Scrollable / Desktop Grid (No Scrollbar) */}
        <div className="max-h-[380px] lg:max-h-none overflow-y-auto lg:overflow-visible [scrollbar-width:none] [-ms-overflow-style:none] [&::-webkit-scrollbar]:hidden pr-1">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-7 gap-3">
            {steps.map((item, idx) => {
              const IconComponent = item.icon;
              return (
                <motion.div
                  key={item.step}
                  initial={{ opacity: 0, y: 15 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  viewport={{ once: true }}
                  transition={{ duration: 0.4, delay: idx * 0.05 }}
                  className="p-3 sm:p-3.5 rounded-xl bg-slate-50 dark:bg-slate-900/80 border border-slate-200/80 dark:border-slate-800 flex flex-col justify-between space-y-2 hover:border-[#0F766E]/40 transition-colors"
                >
                  <div className="flex items-center justify-between">
                    <span className="text-[10px] font-mono font-bold text-slate-400">{item.step}</span>
                    <div className="w-7 h-7 rounded-lg bg-teal-50 dark:bg-teal-950 flex items-center justify-center text-[#0F766E] dark:text-teal-400">
                      <IconComponent className="w-3.5 h-3.5" />
                    </div>
                  </div>

                  <div>
                    <h3 className="text-xs font-bold text-slate-900 dark:text-white mb-0.5">{item.actor}</h3>
                    <p className="text-[11px] font-medium text-slate-600 dark:text-slate-300 mb-1.5">{item.action}</p>
                    <div className="p-1.5 rounded-md bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-900/40 text-[9px] sm:text-[10px] text-rose-700 dark:text-rose-300 font-medium">
                      ⚠️ {item.friction}
                    </div>
                  </div>

                  {idx < steps.length - 1 && (
                    <div className="flex lg:hidden justify-center text-slate-400 pt-0.5">
                      <ArrowDown className="w-3.5 h-3.5" />
                    </div>
                  )}
                </motion.div>
              );
            })}
          </div>
        </div>

        {/* Bottom Callout */}
        <div className="mt-8 p-4 rounded-xl bg-[#F0FDFA] dark:bg-slate-900/80 border border-teal-200 dark:border-teal-800/80 flex flex-col sm:flex-row items-center justify-between gap-3 text-center sm:text-left">
          <div>
            <h4 className="text-xs sm:text-sm font-extrabold text-slate-900 dark:text-white mb-0.5">
              Curexal standardizes healthcare coordination
            </h4>
            <p className="text-[11px] text-slate-600 dark:text-slate-400">
              No repeated hospital trips. Direct digital test requests, result dispatch, and prescription tracking.
            </p>
          </div>
          <div className="inline-flex items-center gap-1.5 text-[10px] sm:text-xs font-bold text-[#0F766E] dark:text-teal-400 bg-white dark:bg-teal-950 px-3 py-1.5 rounded-lg border border-teal-200 dark:border-teal-800 flex-shrink-0">
            <span>Seamless Coordination</span>
            <ArrowRight className="w-3 h-3" />
          </div>
        </div>

      </div>
    </section>
  );
}
