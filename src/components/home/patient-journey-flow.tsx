import { motion } from "framer-motion";
import { Search, Calendar, FileSpreadsheet, Activity, FileCheck2, HeartHandshake, RefreshCw } from "lucide-react";

const journeyStages = [
  {
    stage: "01",
    label: "DISCOVER",
    title: "Find accredited labs & clinics",
    description: "Search tests by location, turnaround time, and price transparency.",
    icon: Search,
  },
  {
    stage: "02",
    label: "BOOK",
    title: "Schedule appointment or home sample",
    description: "Reserve your slot online with zero waiting time or phone calls.",
    icon: Calendar,
  },
  {
    stage: "03",
    label: "REFER",
    title: "Clinic places digital order",
    description: "Doctor sends e-referral directly to diagnostic partner.",
    icon: FileSpreadsheet,
  },
  {
    stage: "04",
    label: "DIAGNOSE",
    title: "Specimen testing & analyzer processing",
    description: "Automated LIMS barcode accessioning and result ingestion.",
    icon: Activity,
  },
  {
    stage: "05",
    label: "RESULT",
    title: "Pathologist signed PDF delivery",
    description: "Verified result dispatches directly to doctor EMR & patient mobile.",
    icon: FileCheck2,
  },
  {
    stage: "06",
    label: "TREAT",
    title: "E-Prescription & pharmacy dispensing",
    description: "Pharmacy receives verified prescription without handwriting errors.",
    icon: HeartHandshake,
  },
  {
    stage: "07",
    label: "FOLLOW UP",
    title: "Longitudinal health record history",
    description: "All past diagnostics securely stored in one digital health vault.",
    icon: RefreshCw,
  },
];

export function PatientJourneyFlowSection() {
  return (
    <section id="patient-journey" className="py-10 sm:py-16 bg-white dark:bg-[#0B1120] text-slate-900 dark:text-white relative overflow-hidden">
      <div className="max-w-[1280px] mx-auto px-4 sm:px-6">
        
        {/* Editorial Section Header */}
        <div className="max-w-2xl mb-8 sm:mb-12">
          <h2 className="text-2xl sm:text-4xl font-black tracking-tight text-slate-900 dark:text-white leading-tight mb-3">
            One patient. <br />
            <span className="bg-clip-text text-transparent bg-gradient-to-r from-[#0F766E] via-[#0D9488] to-[#14B8A6]">
              One connected journey.
            </span>
          </h2>

          <p className="text-xs sm:text-sm text-slate-600 dark:text-slate-300">
            From discovering diagnostic tests to receiving e-prescriptions, Curexal connects every step into a single continuous workflow.
          </p>
        </div>

        {/* Mobile Vertical Scrollable / Desktop 7-Stage Horizontal Container */}
        <div className="max-h-[360px] lg:max-h-none overflow-y-auto lg:overflow-visible [scrollbar-width:none] [-ms-overflow-style:none] [&::-webkit-scrollbar]:hidden pr-1">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-7 gap-3">
            {journeyStages.map((stg, idx) => {
              const IconComponent = stg.icon;
              return (
                <motion.div
                  key={stg.stage}
                  initial={{ opacity: 0, y: 15 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  viewport={{ once: true }}
                  transition={{ duration: 0.4, delay: idx * 0.05 }}
                  className="p-3 sm:p-3.5 rounded-xl bg-slate-50 dark:bg-slate-900/60 border border-slate-200/80 dark:border-slate-800 flex flex-col justify-between space-y-2 hover:border-[#0F766E]/40 transition-colors"
                >
                  <div className="flex items-center justify-between">
                    <div className="w-8 h-8 rounded-lg bg-[#0F766E] text-white flex items-center justify-center font-black text-xs shadow-xs">
                      <IconComponent className="w-4 h-4 text-white" />
                    </div>

                    <span className="text-[9px] font-mono font-bold text-[#0F766E] dark:text-teal-400 bg-teal-50 dark:bg-teal-950/60 px-2 py-0.5 rounded-full border border-teal-200 dark:border-teal-800/80">
                      {stg.stage}
                    </span>
                  </div>

                  <div>
                    <span className="text-[9px] font-black text-[#0F766E] dark:text-teal-400 uppercase tracking-wider block mb-0.5">
                      {stg.label}
                    </span>
                    <h3 className="text-xs font-bold text-slate-900 dark:text-white leading-snug mb-1">
                      {stg.title}
                    </h3>
                    <p className="text-[10px] text-slate-600 dark:text-slate-400 leading-normal">
                      {stg.description}
                    </p>
                  </div>
                </motion.div>
              );
            })}
          </div>
        </div>

      </div>
    </section>
  );
}
