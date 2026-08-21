import { useState } from "react";
import { cn } from "@/lib/utils";
import {
  Globe,
  Building2,
  FlaskConical,
  HeartPulse,
  Stethoscope,
  Scale,
  Pill,
  ArrowRight,
  CheckCircle2,
  Zap,
  Network,
  Shield,
} from "lucide-react";

const nodes = [
  {
    label: "Diagnostic Laboratories",
    icon: FlaskConical,
    active: true,
    color: "#0F766E",
    tagline: "The engine of diagnostics",
    description: "Laboratories on the Curexal network process specimens, validate results through pathologist sign-off, and dispatch verified reports to referring clinicians and patients, all digitally, all tracked, and all audit-ready.",
    capabilities: [
      "Specimen accessioning & barcode tracking",
      "Analyzer integration & result auto-mapping",
      "Pathologist electronic sign-off with audit trail",
      "Automatic result dispatch to referring providers",
    ],
    metric: "LIS/LIMS",
    metricLabel: "Workflows supported",
  },
  {
    label: "Hospitals & Clinics",
    icon: Stethoscope,
    active: true,
    color: "#0D9488",
    tagline: "Where patient journeys begin",
    description: "Clinics and hospitals create digital referrals that reach connected laboratories instantly. Results flow back automatically the moment a pathologist signs off, with no phone calls, paper, or chasing.",
    capabilities: [
      "Digital referral creation in seconds",
      "Real-time result receipt from partner labs",
      "Patient chart integration & clinical notes",
      "Referral commission tracking & settlement",
    ],
    metric: "EMR",
    metricLabel: "Integrations supported",
  },
  {
    label: "Patients",
    icon: HeartPulse,
    active: true,
    color: "#14B8A6",
    tagline: "Healthcare, finally accessible",
    description: "Patients on the Curexal network access their verified lab results through a secure Patient Vault, accessible from any device, anywhere. They find labs through the marketplace, book tests, compare prices, and own their complete health history.",
    capabilities: [
      "Secure Patient Vault with full test history",
      "Verified PDF reports with digital signatures",
      "Marketplace: find labs, compare prices, book tests",
      "Notifications via email, SMS, or WhatsApp",
    ],
    metric: "256-bit",
    metricLabel: "Vault Encryption",
  },
  {
    label: "Diagnostic Centers",
    icon: Building2,
    active: true,
    color: "#0F766E",
    tagline: "Specialized testing, networked",
    description: "Imaging centers, specialized pathology labs, and reference laboratories join the network to become discoverable through the marketplace and receive referrals from connected clinics across the country.",
    capabilities: [
      "Discoverable through the diagnostic marketplace",
      "Receive referrals from connected clinics",
      "Specialized test catalog management",
      "Cross-organizational result sharing",
    ],
    metric: "DICOM",
    metricLabel: "Imaging standard",
  },
  {
    label: "Pharmacies",
    icon: Pill,
    active: false,
    color: "#6B7280",
    tagline: "Coming to the network",
    description: "Pharmacies will join the Curexal network to receive verified prescriptions, coordinate dispensing, and close the loop on the patient journey, from diagnosis to treatment.",
    capabilities: [
      "Digital prescription verification",
      "Connected dispensing workflows",
      "Patient medication history access",
      "Order coordination across the network",
    ],
    metric: "e-Rx",
    metricLabel: "Coming soon",
  },
  {
    label: "Regulatory Bodies",
    icon: Scale,
    active: false,
    color: "#6B7280",
    tagline: "Compliance by design",
    description: "Regulatory bodies and accreditation agencies will access anonymized, aggregated quality metrics and audit data from the network, enabling oversight without compromising individual patient data or provider independence.",
    capabilities: [
      "Aggregated quality metrics & TAT reports",
      "Immutable audit trail access",
      "Accreditation compliance dashboards",
      "Population health intelligence",
    ],
    metric: "ISO 15189",
    metricLabel: "Compliance standard",
  },
];

export function EcosystemSection() {
  const [activeIdx, setActiveIdx] = useState(0);
  const current = nodes[activeIdx];
  const Icon = current.icon;

  return (
    <section
      id="ecosystem"
      className="section-padding bg-[#F8FAFC] dark:bg-[#0B1120] border-y border-gray-100 dark:border-[#1F2937]"
    >
      <div className="max-w-[1280px] mx-auto px-6">

        {/* Header */}
        <div className="text-center max-w-2xl mx-auto mb-14">
          <div className="accent-line mx-auto mb-4" />
          <h2 className="text-section text-gray-900 dark:text-white mb-4">
            One Network.<br />Every Healthcare Provider Connected.
          </h2>
          <p className="text-body text-gray-500 dark:text-gray-400">
            Independent healthcare providers across Africa collaborate as one connected system, without sacrificing their independence or ownership of data. Click each provider to explore.
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-10">

          {/* Left: Provider Selector */}
          <div className="lg:col-span-5">
            <div className="flex flex-col gap-1.5">
              {nodes.map((node, idx) => {
                const NodeIcon = node.icon;
                const isActive = idx === activeIdx;

                return (
                  <button
                    key={node.label}
                    onClick={() => setActiveIdx(idx)}
                    className={cn(
                      "relative w-full text-left px-4 py-3.5 rounded-[12px] border transition-all duration-300 cursor-pointer bg-transparent",
                      isActive
                        ? "border-gray-200 dark:border-[#374151] bg-white dark:bg-[#111827] shadow-md"
                        : "border-transparent hover:border-gray-200 dark:hover:border-[#374151] hover:bg-white/60 dark:hover:bg-[#111827]/40"
                    )}
                  >
                    <div className="flex items-center gap-3">
                      {/* Icon */}
                      <div
                        className={cn(
                          "w-10 h-10 rounded-[10px] flex items-center justify-center flex-shrink-0 transition-all duration-300",
                          isActive ? "scale-105" : ""
                        )}
                        style={{
                          backgroundColor: isActive ? `${node.color}12` : node.active ? `${node.color}08` : "#F3F4F6",
                          border: isActive ? `2px solid ${node.color}30` : "2px solid transparent",
                        }}
                      >
                        <NodeIcon
                          className="w-5 h-5 transition-colors"
                          style={{ color: isActive || node.active ? node.color : "#9CA3AF" }}
                        />
                      </div>

                      {/* Text */}
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          <span className={cn(
                            "text-[14px] font-semibold transition-colors",
                            isActive ? "text-gray-900 dark:text-white" : "text-gray-600 dark:text-gray-400"
                          )}>
                            {node.label}
                          </span>
                        </div>
                        <p className={cn(
                          "text-xs transition-colors",
                          isActive ? "text-gray-500 dark:text-gray-400" : "text-gray-400 dark:text-gray-500"
                        )}>
                          {node.tagline}
                        </p>
                      </div>

                      {/* Status badge */}
                      <span className={cn(
                        "flex-shrink-0 px-2 py-0.5 rounded-full text-[10px] font-bold transition-all",
                        node.active
                          ? isActive
                            ? "text-white bg-[#0F766E]"
                            : "text-[#0F766E] bg-[#F0FDFA] dark:bg-[#0F766E]/10"
                          : "text-gray-400 bg-gray-100 dark:bg-[#1F2937]"
                      )}>
                        {node.active ? "On Network" : "Coming Soon"}
                      </span>
                    </div>
                  </button>
                );
              })}
            </div>

            {/* Network count */}
            <div className="mt-5 px-4 flex items-center gap-2 text-xs text-gray-400 dark:text-gray-500">
              <Network className="w-3.5 h-3.5 text-[#0F766E]" />
              <span className="font-semibold">{nodes.filter(n => n.active).length} provider types</span>
              <span>active on the network</span>
            </div>
          </div>

          {/* Right: Active Provider Detail */}
          <div className="lg:col-span-7">
            <div
              key={activeIdx}
              className="h-full rounded-[16px] border border-gray-200 dark:border-[#1F2937] bg-white dark:bg-[#111827] overflow-hidden animate-fade-in flex flex-col"
              style={{ boxShadow: "0 1px 2px rgba(0,0,0,0.04), 0 8px 32px rgba(0,0,0,0.06)" }}
            >
              {/* Colored header bar */}
              <div
                className="px-8 py-5 flex items-center justify-between"
                style={{ backgroundColor: `${current.color}08`, borderBottom: `1px solid ${current.color}15` }}
              >
                <div className="flex items-center gap-3">
                  <div
                    className="w-11 h-11 rounded-[10px] flex items-center justify-center"
                    style={{ backgroundColor: `${current.color}15` }}
                  >
                    <Icon className="w-5 h-5" style={{ color: current.color }} />
                  </div>
                  <div>
                    <h3 className="text-lg font-bold text-gray-900 dark:text-white">{current.label}</h3>
                    <p className="text-xs font-medium" style={{ color: current.color }}>{current.tagline}</p>
                  </div>
                </div>
                {/* Metric */}
                <div className="text-right hidden sm:block">
                  <p className="text-2xl font-bold text-gray-900 dark:text-white leading-none">{current.metric}</p>
                  <p className="text-[11px] text-gray-400 dark:text-gray-500 mt-0.5">{current.metricLabel}</p>
                </div>
              </div>

              {/* Body */}
              <div className="px-8 py-6 flex-1 flex flex-col justify-between">
                <div>
                  {/* Description */}
                  <p className="text-[15px] leading-relaxed text-gray-600 dark:text-gray-300 mb-6 animate-fade-up">
                    {current.description}
                  </p>

                  {/* Capabilities */}
                  <div className="space-y-2.5">
                    <p className="text-[10px] font-bold text-gray-400 dark:text-gray-500 uppercase tracking-widest mb-3">
                      Network Capabilities
                    </p>
                    {current.capabilities.map((cap, i) => (
                      <div
                        key={cap}
                        className="flex items-center gap-3 animate-fade-up"
                        style={{ animationDelay: `${(i + 1) * 80}ms` }}
                      >
                        <div
                          className="w-5 h-5 rounded-full flex items-center justify-center flex-shrink-0"
                          style={{ backgroundColor: `${current.color}12` }}
                        >
                          <CheckCircle2 className="w-3 h-3" style={{ color: current.color }} />
                        </div>
                        <span className="text-sm text-gray-700 dark:text-gray-300">{cap}</span>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Footer */}
                <div className="mt-6 pt-5 border-t border-gray-100 dark:border-[#1F2937] flex items-center justify-between">
                  <div className="flex items-center gap-1.5 text-xs text-gray-400 dark:text-gray-500">
                    <Shield className="w-3.5 h-3.5 text-[#0F766E]" />
                    <span>Data sovereignty guaranteed</span>
                  </div>
                  {current.active ? (
                    <div className="flex items-center gap-1.5 text-xs font-semibold text-[#0F766E]">
                      <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
                      Active on network
                    </div>
                  ) : (
                    <span className="text-xs font-semibold text-gray-400 dark:text-gray-500">
                      Joining the network soon
                    </span>
                  )}
                </div>
              </div>
            </div>
          </div>

        </div>
      </div>
    </section>
  );
}
