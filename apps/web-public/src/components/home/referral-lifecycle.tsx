import React, { useState } from "react";
import {
  FileCheck2,
  UserCheck,
  Inbox,
  TestTube2,
  Microscope,
  CheckCircle,
  Send,
  Heart,
  ChevronRight,
  ChevronLeft,
  Sparkles,
} from "lucide-react";

const lifecycleSteps = [
  {
    step: 1,
    title: "1. Doctor Creates Referral",
    desc: "Clinician places digital laboratory referral during consultation, specifying required tests and clinical urgency.",
    icon: FileCheck2,
    actor: "Clinic / Doctor",
  },
  {
    step: 2,
    title: "2. Patient Directed to Provider",
    desc: "Patient receives digital notification with instructions, location, and test preparation details.",
    icon: UserCheck,
    actor: "Patient Portal / SMS",
  },
  {
    step: 3,
    title: "3. Laboratory Receives Request",
    desc: "Diagnostic lab receives electronic requisition in worklist queue before patient arrives.",
    icon: Inbox,
    actor: "Laboratory LIMS",
  },
  {
    step: 4,
    title: "4. Sample Collected & Barcoded",
    desc: "Phlebotomist accession sample with unique barcode label, creating verifiable chain of custody.",
    icon: TestTube2,
    actor: "Phlebotomy / Intake",
  },
  {
    step: 5,
    title: "5. Laboratory Processes Test",
    desc: "Analyzer instrument processes specimen and transmits auto-results to validation queue.",
    icon: Microscope,
    actor: "Lab Analyzer & Tech",
  },
  {
    step: 6,
    title: "6. Result Authorized",
    desc: "Pathologist reviews flags, attaches digital signature, and authorizes report delivery.",
    icon: CheckCircle,
    actor: "Pathologist Sign-Off",
  },
  {
    step: 7,
    title: "7. Referring Provider Receives Result",
    desc: "Electronically signed result instantly drops into referring doctor's EMR clinical chart.",
    icon: Send,
    actor: "Doctor EMR Chart",
  },
  {
    step: 8,
    title: "8. Patient Continues Care",
    desc: "Patient accesses verified report while clinician initiates timely treatment plan.",
    icon: Heart,
    actor: "Coordinated Care",
  },
];

export function ReferralLifecycleSection() {
  const [activeStep, setActiveStep] = useState(0);

  return (
    <section className="py-20 bg-white dark:bg-[#0B1120] border-b border-slate-100 dark:border-slate-800">
      <div className="max-w-[1280px] mx-auto px-4 sm:px-6">
        
        {/* Header */}
        <div className="text-center max-w-3xl mx-auto mb-14">
          <h2 className="text-3xl sm:text-4xl font-extrabold text-slate-900 dark:text-white tracking-tight mb-4">
            From referral to result, without losing the thread.
          </h2>
          <p className="text-sm sm:text-base text-slate-600 dark:text-slate-400">
            Track every diagnostic request across facilities with live status updates and automated digital handoffs.
          </p>
        </div>

        {/* Interactive Lifecycle Stepper */}
        <div className="bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 sm:p-8 shadow-sm">
          
          {/* Step Badges Bar */}
          <div className="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-8 gap-2 mb-8">
            {lifecycleSteps.map((item, index) => {
              const Icon = item.icon;
              const isCurrent = activeStep === index;
              return (
                <button
                  key={item.step}
                  onClick={() => setActiveStep(index)}
                  className={`p-2.5 rounded-2xl border text-left transition-all cursor-pointer flex flex-col justify-between h-24 ${
                    isCurrent
                      ? "bg-[#0F766E] text-white border-[#0F766E] shadow-md scale-102"
                      : "bg-white dark:bg-slate-800/60 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-700 hover:border-teal-500/50"
                  }`}
                >
                  <div className="flex items-center justify-between w-full">
                    <span className={`text-[10px] font-bold ${isCurrent ? "text-teal-100" : "text-slate-400"}`}>
                      0{item.step}
                    </span>
                    <Icon className={`w-4 h-4 ${isCurrent ? "text-teal-200" : "text-[#0F766E]"}`} />
                  </div>
                  <span className="text-[11px] font-bold leading-tight line-clamp-2">
                    {item.title.split(". ")[1]}
                  </span>
                </button>
              );
            })}
          </div>

          {/* Active Step Detail Visual */}
          <div className="p-6 sm:p-8 rounded-2xl bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 flex flex-col md:flex-row items-center justify-between gap-6">
            
            <div className="space-y-3 max-w-xl">
              <div className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-teal-50 dark:bg-teal-950/80 border border-teal-200 dark:border-teal-800 text-[#0F766E] dark:text-teal-300 text-[10px] font-bold">
                <Sparkles className="w-3 h-3" />
                <span>Actor: {lifecycleSteps[activeStep].actor}</span>
              </div>
              <h3 className="text-xl sm:text-2xl font-bold text-slate-900 dark:text-white">
                {lifecycleSteps[activeStep].title}
              </h3>
              <p className="text-sm text-slate-600 dark:text-slate-300 leading-relaxed">
                {lifecycleSteps[activeStep].desc}
              </p>
            </div>

            {/* Stepper Navigation Buttons */}
            <div className="flex items-center gap-2 flex-shrink-0">
              <button
                disabled={activeStep === 0}
                onClick={() => setActiveStep((prev) => Math.max(0, prev - 1))}
                className="px-4 py-2 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-800 text-xs font-bold text-slate-700 dark:text-slate-200 disabled:opacity-40 cursor-pointer flex items-center gap-1"
              >
                <ChevronLeft className="w-4 h-4" />
                <span>Previous</span>
              </button>
              <button
                disabled={activeStep === lifecycleSteps.length - 1}
                onClick={() => setActiveStep((prev) => Math.min(lifecycleSteps.length - 1, prev + 1))}
                className="px-4 py-2 rounded-xl bg-[#0F766E] hover:bg-[#115E59] text-white text-xs font-bold disabled:opacity-40 cursor-pointer flex items-center gap-1 border-0 shadow-sm"
              >
                <span>Next Stage</span>
                <ChevronRight className="w-4 h-4" />
              </button>
            </div>

          </div>

        </div>

      </div>
    </section>
  );
}
