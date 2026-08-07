import { useState } from "react";
import {
  Building2,
  FlaskConical,
  Users,
  Activity,
  Globe,
  CheckCircle2,
  ArrowRight,
  Zap,
  Network,
  ShieldCheck,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { Link } from "react-router-dom";

const steps = [
  {
    step: "01",
    title: "Join the Network",
    subtitle: "Create your organization",
    desc: "Register as a laboratory, clinic, or diagnostic center. Verify your facility credentials and get your own secure, isolated workspace, provisioned in under 60 seconds.",
    icon: Building2,
    color: "#0F766E",
    highlights: [
      "Dedicated database per organization",
      "Schema-isolated data sovereignty",
      "Instant workspace provisioning",
    ],
  },
  {
    step: "02",
    title: "Configure Your Workflows",
    subtitle: "Set up your clinical directory",
    desc: "Define your test panels, analyte reference ranges, pricing, and notification thresholds. Import your existing catalog or start from Curexal's built-in clinical directory templates.",
    icon: FlaskConical,
    color: "#0D9488",
    highlights: [
      "Clinical test panels & reference ranges",
      "Custom pricing per test or package",
      "Automated critical value alerts",
    ],
  },
  {
    step: "03",
    title: "Onboard Your Team",
    subtitle: "Invite & assign roles",
    desc: "Add pathologists, technologists, phlebotomists, and administrators. Fine-grained role-based access ensures everyone sees exactly what they need, with no security risks.",
    icon: Users,
    color: "#14B8A6",
    highlights: [
      "Role-based access control (RBAC)",
      "Pathologist, technologist, admin roles",
      "Audit trail on every action",
    ],
  },
  {
    step: "04",
    title: "Go Live on the Network",
    subtitle: "Start receiving referrals",
    desc: "Publish your lab on the diagnostic marketplace. Accept referrals from connected clinics. Deliver results digitally to clinicians and patients. Commission settlements tracked automatically.",
    icon: Globe,
    color: "#0F766E",
    highlights: [
      "Discoverable on the marketplace",
      "Digital referrals from partner clinics",
      "Automated commission reconciliation",
    ],
  },
];

export function HowItWorks() {
  const [activeStep, setActiveStep] = useState(0);
  const current = steps[activeStep];
  const Icon = current.icon;

  return (
    <section id="how-it-works" className="section-padding bg-white dark:bg-[#0B1120]">
      <div className="max-w-[1280px] mx-auto px-6">

        {/* Header */}
        <div className="text-center max-w-2xl mx-auto mb-14">
          <div className="inline-flex items-center gap-2 px-3 py-1.5 mb-5 rounded-full border border-[#0F766E]/20 bg-[#F0FDFA] dark:bg-[#0F766E]/10">
            <Zap className="h-3.5 w-3.5 text-[#0F766E]" />
            <span className="text-xs font-semibold text-[#0F766E] tracking-wide">
              Go Live in Days, Not Months
            </span>
          </div>
          <h2 className="text-section text-gray-900 dark:text-white mb-4">
            Four Steps to Join<br />the Healthcare Network.
          </h2>
          <p className="text-body text-gray-500 dark:text-gray-400">
            From signup to your first digital referral, Curexal gets your organization connected and operational faster than any traditional system.
          </p>
        </div>

        {/* Interactive Stepper */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-12">

          {/* Left: Step Selector */}
          <div className="lg:col-span-5 flex flex-col gap-2">
            {steps.map((step, idx) => {
              const StepIcon = step.icon;
              const isActive = idx === activeStep;
              const isCompleted = idx < activeStep;

              return (
                <button
                  key={step.step}
                  onClick={() => setActiveStep(idx)}
                  className={cn(
                    "relative w-full text-left p-4 rounded-[14px] border transition-all duration-300 cursor-pointer bg-transparent",
                    isActive
                      ? "border-[#0F766E]/30 bg-[#F0FDFA] dark:bg-[#0F766E]/10 shadow-sm"
                      : "border-gray-200 dark:border-[#1F2937] hover:border-gray-300 dark:hover:border-[#374151] hover:bg-gray-50 dark:hover:bg-[#1F2937]/50"
                  )}
                >
                  <div className="flex items-center gap-4">
                    {/* Step number / check */}
                    <div
                      className={cn(
                        "w-10 h-10 rounded-[10px] flex items-center justify-center flex-shrink-0 transition-all duration-300",
                        isActive
                          ? "bg-[#0F766E] shadow-sm"
                          : isCompleted
                            ? "bg-[#0F766E]/10"
                            : "bg-gray-100 dark:bg-[#1F2937]"
                      )}
                    >
                      {isCompleted ? (
                        <CheckCircle2 className="w-5 h-5 text-[#0F766E]" />
                      ) : (
                        <StepIcon
                          className={cn(
                            "w-5 h-5 transition-colors",
                            isActive ? "text-white" : "text-gray-400 dark:text-gray-500"
                          )}
                        />
                      )}
                    </div>

                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2">
                        <span className={cn(
                          "text-[10px] font-bold uppercase tracking-widest",
                          isActive ? "text-[#0F766E]" : "text-gray-400 dark:text-gray-500"
                        )}>
                          Step {step.step}
                        </span>
                        {isActive && (
                          <span className="w-1.5 h-1.5 rounded-full bg-[#0F766E] animate-pulse" />
                        )}
                      </div>
                      <h3 className={cn(
                        "text-[15px] font-semibold mt-0.5 transition-colors",
                        isActive ? "text-gray-900 dark:text-white" : "text-gray-600 dark:text-gray-400"
                      )}>
                        {step.title}
                      </h3>
                    </div>

                    {/* Arrow indicator */}
                    <ArrowRight className={cn(
                      "w-4 h-4 flex-shrink-0 transition-all duration-300",
                      isActive ? "text-[#0F766E] translate-x-0 opacity-100" : "text-gray-300 -translate-x-1 opacity-0"
                    )} />
                  </div>

                  {/* Active progress bar at bottom */}
                  {isActive && (
                    <div className="absolute bottom-0 left-4 right-4 h-[2px] bg-gray-200 dark:bg-[#374151] rounded-full overflow-hidden">
                      <div
                        className="h-full bg-[#0F766E] rounded-full"
                        style={{
                          width: `${((activeStep + 1) / steps.length) * 100}%`,
                          transition: "width 0.5s ease",
                        }}
                      />
                    </div>
                  )}
                </button>
              );
            })}

            {/* Progress indicator */}
            <div className="mt-3 flex items-center justify-between px-2">
              <span className="text-[11px] font-semibold text-gray-400 dark:text-gray-500">
                {activeStep + 1} of {steps.length} steps
              </span>
              <div className="flex gap-1.5">
                {steps.map((_, idx) => (
                  <div
                    key={idx}
                    className={cn(
                      "h-1.5 rounded-full transition-all duration-300",
                      idx === activeStep
                        ? "w-6 bg-[#0F766E]"
                        : idx < activeStep
                          ? "w-1.5 bg-[#0F766E]/40"
                          : "w-1.5 bg-gray-200 dark:bg-[#374151]"
                    )}
                  />
                ))}
              </div>
            </div>
          </div>

          {/* Right: Active Step Detail Card */}
          <div className="lg:col-span-7">
            <div
              key={activeStep}
              className="h-full rounded-[16px] border border-gray-200 dark:border-[#1F2937] bg-[#F8FAFC] dark:bg-[#111827] p-8 lg:p-10 animate-fade-in flex flex-col justify-between"
              style={{ boxShadow: "0 1px 2px rgba(0,0,0,0.04), 0 8px 32px rgba(0,0,0,0.06)" }}
            >
              <div>
                {/* Step badge */}
                <div className="flex items-center gap-3 mb-6">
                  <div
                    className="w-12 h-12 rounded-[12px] flex items-center justify-center"
                    style={{ backgroundColor: `${current.color}15` }}
                  >
                    <Icon className="w-6 h-6" style={{ color: current.color }} />
                  </div>
                  <div>
                    <span className="text-[10px] font-bold text-gray-400 dark:text-gray-500 uppercase tracking-widest">
                      Step {current.step}
                    </span>
                    <h3 className="text-xl font-bold text-gray-900 dark:text-white">
                      {current.title}
                    </h3>
                  </div>
                </div>

                {/* Description */}
                <p className="text-[15px] leading-relaxed text-gray-600 dark:text-gray-300 mb-6">
                  {current.desc}
                </p>

                {/* Highlights */}
                <div className="space-y-3">
                  {current.highlights.map((h, i) => (
                    <div
                      key={h}
                      className="flex items-center gap-3 animate-fade-up"
                      style={{ animationDelay: `${i * 100}ms` }}
                    >
                      <div
                        className="w-6 h-6 rounded-full flex items-center justify-center flex-shrink-0"
                        style={{ backgroundColor: `${current.color}15` }}
                      >
                        <CheckCircle2 className="w-3.5 h-3.5" style={{ color: current.color }} />
                      </div>
                      <span className="text-sm text-gray-700 dark:text-gray-300 font-medium">{h}</span>
                    </div>
                  ))}
                </div>
              </div>

              {/* Bottom action */}
              <div className="mt-8 pt-6 border-t border-gray-200 dark:border-[#1F2937] flex items-center justify-between">
                <div className="flex items-center gap-2 text-xs text-gray-400 dark:text-gray-500">
                  <ShieldCheck className="w-3.5 h-3.5 text-[#0F766E]" />
                  <span>Your data stays yours. Always.</span>
                </div>
                {activeStep < steps.length - 1 ? (
                  <button
                    onClick={() => setActiveStep(activeStep + 1)}
                    className="flex items-center gap-1.5 text-sm font-semibold text-[#0F766E] hover:text-[#115E59] transition-colors bg-transparent border-0 cursor-pointer"
                  >
                    Next step
                    <ArrowRight className="w-4 h-4" />
                  </button>
                ) : (
                  <Link to="/book-demo">
                    <button className="flex items-center gap-2 px-4 py-2 rounded-[10px] bg-[#0F766E] hover:bg-[#115E59] text-white text-sm font-semibold transition-colors cursor-pointer border-0">
                      Get Started
                      <ArrowRight className="w-4 h-4" />
                    </button>
                  </Link>
                )}
              </div>
            </div>
          </div>

        </div>
      </div>
    </section>
  );
}
