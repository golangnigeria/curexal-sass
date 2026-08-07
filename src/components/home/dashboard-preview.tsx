import { useState, useEffect } from "react";
import { cn } from "@/lib/utils";
import {
  Building2,
  FlaskConical,
  Stethoscope,
  HeartPulse,
  Send,
  ArrowRight,
  CheckCircle2,
  Clock,
  FileText,
  CreditCard,
  Zap,
} from "lucide-react";

const networkFlow = [
  {
    step: "01",
    icon: Stethoscope,
    actor: "Clinic",
    location: "Port Harcourt",
    action: "Refers a patient for Full Blood Count",
    detail: "A clinician at a connected clinic creates a digital referral for a patient who needs a Full Blood Count (FBC). The referral is sent instantly to the nearest accredited laboratory on the network, with no paper, no phone calls, and no envelopes.",
    color: "#0F766E",
    status: "Referral Sent",
    statusIcon: Send,
    highlights: [
      "Patient details securely transmitted",
      "Laboratory auto-selected based on proximity & capability",
      "Referring clinician tracked for commission settlement",
    ],
  },
  {
    step: "02",
    icon: FlaskConical,
    actor: "Laboratory",
    location: "Independent Lab, Lagos",
    action: "Receives, accessions, processes, validates",
    detail: "The laboratory receives the referral automatically. The specimen is accessioned, assigned a barcode, processed through the analyzer, and the technologist logs the results. The pathologist validates and electronically signs the report.",
    color: "#0D9488",
    status: "Processing",
    statusIcon: Clock,
    highlights: [
      "Specimen chain of custody tracked end-to-end",
      "Analyzer results auto-mapped to patient record",
      "Pathologist e-signs with full audit trail",
    ],
  },
  {
    step: "03",
    icon: Send,
    actor: "Clinician",
    location: "Referring Doctor",
    action: "Automatically receives the validated result",
    detail: "The moment the pathologist signs off, the referring clinician receives the validated result in real-time. No phone call. No waiting. No manual follow-up. The result appears in their connected dashboard with the patient's full referral context.",
    color: "#14B8A6",
    status: "Result Delivered",
    statusIcon: FileText,
    highlights: [
      "Real-time push notification to clinician",
      "Full referral context preserved",
      "Zero manual coordination required",
    ],
  },
  {
    step: "04",
    icon: HeartPulse,
    actor: "Patient",
    location: "Anywhere",
    action: "Gets notified and accesses verified report",
    detail: "The patient receives a notification (email, SMS, or WhatsApp) that their result is ready. They log into their Patient Vault and access the verified PDF report, digitally signed by the pathologist, accredited by the laboratory, and accessible from any device.",
    color: "#2DD4BF",
    status: "Patient Notified",
    statusIcon: CheckCircle2,
    highlights: [
      "Verified PDF with digital pathologist signature",
      "Accessible on any device, anywhere",
      "Full test history in Patient Vault",
    ],
  },
  {
    step: "05",
    icon: CreditCard,
    actor: "Settlement",
    location: "Automatic",
    action: "Commission automatically reconciled",
    detail: "The referral commission between the clinic and the laboratory is automatically calculated, tracked, and reconciled. Every naira is auditable. No spreadsheets, no disputes, no manual accounting.",
    color: "#0F766E",
    status: "Settled",
    statusIcon: CheckCircle2,
    highlights: [
      "Automatic commission calculation",
      "Full transaction audit trail",
      "Reconciliation without spreadsheets",
    ],
  },
];

export function NetworkFlowSection() {
  const [activeStep, setActiveStep] = useState(0);
  const [isAutoPlaying, setIsAutoPlaying] = useState(true);
  const current = networkFlow[activeStep];
  const Icon = current.icon;
  const StatusIcon = current.statusIcon;

  // Auto-advance every 5 seconds
  useEffect(() => {
    if (!isAutoPlaying) return;
    const timer = setInterval(() => {
      setActiveStep((prev) => (prev + 1) % networkFlow.length);
    }, 5000);
    return () => clearInterval(timer);
  }, [isAutoPlaying]);

  const handleStepClick = (idx: number) => {
    setActiveStep(idx);
    setIsAutoPlaying(false);
  };

  return (
    <section
      id="network-flow"
      className="section-padding bg-[#F8FAFC] dark:bg-[#0B1120] border-y border-gray-100 dark:border-[#1F2937]"
    >
      <div className="max-w-[1280px] mx-auto px-6">

        {/* Header */}
        <div className="text-center max-w-2xl mx-auto mb-14">
          <div className="accent-line mx-auto mb-4" />
          <h2 className="text-section text-gray-900 dark:text-white mb-4">
            This Is How Healthcare<br />Should Work.
          </h2>
          <p className="text-body text-gray-500 dark:text-gray-400">
            Follow one patient journey across the entire network, from clinic referral to commission settlement. Click each step to explore.
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-10">

          {/* Left: Visual Journey Timeline */}
          <div className="lg:col-span-5">

            {/* Horizontal progress bar */}
            <div className="mb-6 flex items-center gap-1">
              {networkFlow.map((_, idx) => (
                <div key={idx} className="flex-1 h-1 rounded-full overflow-hidden bg-gray-200 dark:bg-[#1F2937]">
                  <div
                    className={cn(
                      "h-full rounded-full transition-all duration-500",
                      idx < activeStep
                        ? "bg-[#0F766E] w-full"
                        : idx === activeStep
                          ? "bg-[#0F766E] animate-shimmer"
                          : "w-0"
                    )}
                    style={idx === activeStep ? { width: "100%" } : idx < activeStep ? { width: "100%" } : { width: "0%" }}
                  />
                </div>
              ))}
            </div>

            {/* Step cards */}
            <div className="flex flex-col gap-2">
              {networkFlow.map((node, idx) => {
                const NodeIcon = node.icon;
                const isActive = idx === activeStep;
                const isCompleted = idx < activeStep;

                return (
                  <button
                    key={node.step}
                    onClick={() => handleStepClick(idx)}
                    className={cn(
                      "relative w-full text-left px-4 py-3.5 rounded-[12px] border transition-all duration-300 cursor-pointer bg-transparent",
                      isActive
                        ? "border-[#0F766E]/30 bg-white dark:bg-[#111827] shadow-md"
                        : "border-transparent hover:border-gray-200 dark:hover:border-[#374151] hover:bg-white dark:hover:bg-[#111827]/50"
                    )}
                  >
                    <div className="flex items-center gap-3">
                      {/* Icon orb */}
                      <div
                        className={cn(
                          "w-9 h-9 rounded-full flex items-center justify-center flex-shrink-0 transition-all duration-300",
                          isActive
                            ? "scale-110"
                            : isCompleted
                              ? "opacity-60"
                              : "opacity-40"
                        )}
                        style={{
                          backgroundColor: isActive || isCompleted ? `${node.color}15` : "transparent",
                          border: `2px solid ${isActive ? node.color : isCompleted ? `${node.color}40` : "#E5E7EB"}`,
                        }}
                      >
                        {isCompleted ? (
                          <CheckCircle2 className="w-4 h-4" style={{ color: node.color }} />
                        ) : (
                          <NodeIcon className="w-4 h-4" style={{ color: isActive ? node.color : "#9CA3AF" }} />
                        )}
                      </div>

                      {/* Text */}
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          <span className={cn(
                            "text-[13px] font-semibold transition-colors",
                            isActive ? "text-gray-900 dark:text-white" : "text-gray-500 dark:text-gray-400"
                          )}>
                            {node.actor}
                          </span>
                          {isActive && (
                            <span
                              className="text-[10px] font-bold px-1.5 py-0.5 rounded-full"
                              style={{ color: node.color, backgroundColor: `${node.color}12` }}
                            >
                              {node.status}
                            </span>
                          )}
                        </div>
                        <p className={cn(
                          "text-xs truncate transition-colors",
                          isActive ? "text-gray-500 dark:text-gray-400" : "text-gray-400 dark:text-gray-500"
                        )}>
                          {node.action}
                        </p>
                      </div>

                      {/* Connector arrow */}
                      <ArrowRight className={cn(
                        "w-3.5 h-3.5 flex-shrink-0 transition-all duration-300",
                        isActive ? "text-[#0F766E] opacity-100" : "opacity-0"
                      )} />
                    </div>
                  </button>
                );
              })}
            </div>

            {/* Autoplay toggle */}
            <div className="mt-4 flex items-center justify-between px-2">
              <button
                onClick={() => setIsAutoPlaying(!isAutoPlaying)}
                className="text-[11px] font-semibold text-gray-400 dark:text-gray-500 hover:text-[#0F766E] transition-colors bg-transparent border-0 cursor-pointer flex items-center gap-1.5"
              >
                <div className={cn(
                  "w-2 h-2 rounded-full",
                  isAutoPlaying ? "bg-emerald-500 animate-pulse" : "bg-gray-300 dark:bg-gray-600"
                )} />
                {isAutoPlaying ? "Auto-playing" : "Paused"}
              </button>
              <span className="text-[11px] font-semibold text-gray-400 dark:text-gray-500">
                {activeStep + 1} of {networkFlow.length}
              </span>
            </div>
          </div>

          {/* Right: Active Step Detail */}
          <div className="lg:col-span-7">
            <div
              key={activeStep}
              className="h-full rounded-[16px] border border-gray-200 dark:border-[#1F2937] bg-white dark:bg-[#111827] p-8 lg:p-10 animate-fade-in flex flex-col justify-between"
              style={{ boxShadow: "0 1px 2px rgba(0,0,0,0.04), 0 8px 32px rgba(0,0,0,0.06)" }}
            >
              <div>
                {/* Header row */}
                <div className="flex items-center gap-4 mb-6">
                  <div
                    className="w-14 h-14 rounded-[14px] flex items-center justify-center animate-fade-up"
                    style={{ backgroundColor: `${current.color}12` }}
                  >
                    <Icon className="w-7 h-7" style={{ color: current.color }} />
                  </div>
                  <div>
                    <div className="flex items-center gap-2">
                      <span className="text-[10px] font-bold text-gray-400 dark:text-gray-500 uppercase tracking-widest">
                        Step {current.step}
                      </span>
                      <span className="text-[10px] font-semibold text-gray-400 dark:text-gray-500">
                        · {current.location}
                      </span>
                    </div>
                    <h3 className="text-xl font-bold text-gray-900 dark:text-white mt-0.5">
                      {current.action}
                    </h3>
                  </div>
                </div>

                {/* Status pill */}
                <div
                  className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full mb-6 animate-fade-up delay-100"
                  style={{ backgroundColor: `${current.color}10`, border: `1px solid ${current.color}20` }}
                >
                  <StatusIcon className="w-3.5 h-3.5" style={{ color: current.color }} />
                  <span className="text-xs font-bold" style={{ color: current.color }}>
                    {current.status}
                  </span>
                </div>

                {/* Description */}
                <p className="text-[15px] leading-relaxed text-gray-600 dark:text-gray-300 mb-6 animate-fade-up delay-100">
                  {current.detail}
                </p>

                {/* Highlights */}
                <div className="space-y-3">
                  {current.highlights.map((h, i) => (
                    <div
                      key={h}
                      className="flex items-center gap-3 animate-fade-up"
                      style={{ animationDelay: `${(i + 2) * 100}ms` }}
                    >
                      <div
                        className="w-6 h-6 rounded-full flex items-center justify-center flex-shrink-0"
                        style={{ backgroundColor: `${current.color}12` }}
                      >
                        <Zap className="w-3 h-3" style={{ color: current.color }} />
                      </div>
                      <span className="text-sm text-gray-700 dark:text-gray-300 font-medium">{h}</span>
                    </div>
                  ))}
                </div>
              </div>

              {/* Bottom navigation */}
              <div className="mt-8 pt-6 border-t border-gray-100 dark:border-[#1F2937] flex items-center justify-between">
                <button
                  onClick={() => handleStepClick(Math.max(0, activeStep - 1))}
                  disabled={activeStep === 0}
                  className={cn(
                    "text-sm font-semibold transition-colors bg-transparent border-0 cursor-pointer flex items-center gap-1.5",
                    activeStep === 0
                      ? "text-gray-300 dark:text-gray-600 cursor-default"
                      : "text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200"
                  )}
                >
                  <ArrowRight className="w-4 h-4 rotate-180" />
                  Previous
                </button>

                {activeStep < networkFlow.length - 1 ? (
                  <button
                    onClick={() => handleStepClick(activeStep + 1)}
                    className="flex items-center gap-1.5 text-sm font-semibold text-[#0F766E] hover:text-[#115E59] transition-colors bg-transparent border-0 cursor-pointer"
                  >
                    Next step
                    <ArrowRight className="w-4 h-4" />
                  </button>
                ) : (
                  <div className="flex items-center gap-2 text-sm font-semibold text-[#0F766E]">
                    <CheckCircle2 className="w-4 h-4" />
                    Journey complete
                  </div>
                )}
              </div>
            </div>
          </div>

        </div>

        {/* Summary callout */}
        <div className="mt-10 p-6 rounded-[14px] border border-[#0F766E]/20 bg-[#F0FDFA] dark:bg-[#0F766E]/5 text-center">
          <p className="text-[15px] font-semibold text-[#0F766E] dark:text-teal-400">
            Five organizations. One patient journey. Zero paper. Zero phone calls. Fully traceable.
          </p>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            That's the Curexal network.
          </p>
        </div>

      </div>
    </section>
  );
}
