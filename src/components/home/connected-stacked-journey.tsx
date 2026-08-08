import { motion } from "framer-motion";
import { User, Stethoscope, FlaskConical, Pill, Building2, CheckCircle2 } from "lucide-react";

const stackedCards = [
  {
    role: "PATIENT",
    title: "Patient Health Vault & Access",
    icon: User,
    color: "bg-teal-50 dark:bg-teal-950/60",
    borderColor: "border-teal-200 dark:border-teal-800",
    accentColor: "text-[#0F766E] dark:text-teal-400",
    badge: "Care Recipient",
    features: [
      "Find accredited diagnostic labs & clinics",
      "Direct digital result delivery on mobile",
      "Lifetime verified electronic test history",
      "Zero repeated physical slip trips",
    ],
  },
  {
    role: "CLINIC",
    title: "Clinical EMR & E-Referrals",
    icon: Stethoscope,
    color: "bg-slate-50 dark:bg-slate-900",
    borderColor: "border-slate-200 dark:border-slate-800",
    accentColor: "text-teal-600 dark:text-teal-300",
    badge: "Outpatient Provider",
    features: [
      "Instant digital lab order placement",
      "Real-time referral status tracking",
      "Direct result ingestion to patient chart",
      "Multi-branch consultation management",
    ],
  },
  {
    role: "LABORATORY",
    title: "Laboratory LIMS & Result Dispatch",
    icon: FlaskConical,
    color: "bg-[#F0FDFA] dark:bg-slate-900",
    borderColor: "border-teal-200 dark:border-teal-800",
    accentColor: "text-emerald-600 dark:text-emerald-400",
    badge: "Diagnostic Facility",
    features: [
      "Automated specimen accessioning & barcoding",
      "Sysmex, Mindray & Roche analyzer ingestion",
      "Pathologist digital stamp & e-signatures",
      "ISO 15189 TAT compliance monitoring",
    ],
  },
  {
    role: "PHARMACY",
    title: "Digital Prescription & Inventory POS",
    icon: Pill,
    color: "bg-sky-50 dark:bg-slate-900",
    borderColor: "border-sky-200 dark:border-sky-900",
    accentColor: "text-sky-600 dark:text-sky-400",
    badge: "Dispensing Partner",
    features: [
      "Direct e-prescription receiving from clinics",
      "Automated stock reordering & expiry tracking",
      "Refill subscription delivery management",
      "Zero prescription fulfillment errors",
    ],
  },
  {
    role: "HEALTHCARE PARTNER",
    title: "Suppliers & B2B Procurement Marketplace",
    icon: Building2,
    color: "bg-teal-50/70 dark:bg-slate-900",
    borderColor: "border-teal-200 dark:border-teal-800",
    accentColor: "text-teal-700 dark:text-teal-300",
    badge: "Reagent & Equipment Supplier",
    features: [
      "B2B storefront ordering for labs & clinics",
      "Automated reagent replenishment orders",
      "Transparent settlement reconciliation",
      "Unified B2B healthcare supply chain",
    ],
  },
];

export function ConnectedStackedJourneySection() {
  return (
    <section id="stacked-journey" className="py-10 sm:py-16 bg-[#F8FAFC] dark:bg-[#0B1120] text-slate-900 dark:text-white relative overflow-hidden">
      <div className="max-w-[1280px] mx-auto px-4 sm:px-6 relative z-10">
        
        {/* Section Header */}
        <div className="text-center max-w-2xl mx-auto mb-8 sm:mb-12">
          <div className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-teal-50 dark:bg-teal-950/80 border border-teal-200 dark:border-teal-800 text-[#0F766E] dark:text-teal-400 text-[10px] sm:text-xs font-bold uppercase tracking-wider mb-3">
            <span>CUREXAL OPERATING ARCHITECTURE</span>
          </div>

          <h2 className="text-2xl sm:text-4xl font-black text-slate-900 dark:text-white tracking-tight leading-tight mb-3">
            Curexal connects every node <br className="hidden sm:inline" />
            <span className="bg-clip-text text-transparent bg-gradient-to-r from-[#0F766E] via-[#0D9488] to-[#14B8A6]">
              in the healthcare continuum.
            </span>
          </h2>

          <p className="text-xs sm:text-sm text-slate-600 dark:text-slate-300">
            A single connected network enabling seamless coordination across patients, clinics, diagnostic centers, pharmacies, and equipment suppliers.
          </p>
        </div>

        {/* Vertical Stacked Cards with Hidden Scrollbar on Mobile */}
        <div className="max-w-2xl mx-auto space-y-4 max-h-[420px] lg:max-h-none overflow-y-auto lg:overflow-visible [scrollbar-width:none] [-ms-overflow-style:none] [&::-webkit-scrollbar]:hidden pr-1">
          {stackedCards.map((card, idx) => {
            const IconComponent = card.icon;
            return (
              <motion.div
                key={card.role}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ duration: 0.4, delay: idx * 0.08 }}
                className={`sticky top-20 sm:top-24 p-4 sm:p-5 rounded-2xl bg-white dark:bg-slate-900 border ${card.borderColor} shadow-md transition-all box-border`}
                style={{
                  zIndex: idx + 10,
                  marginTop: idx === 0 ? 0 : `-8px`,
                }}
              >
                <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-2 pb-3 border-b border-slate-100 dark:border-slate-800">
                  <div className="flex items-center gap-2.5">
                    <div className="w-8 h-8 sm:w-10 sm:h-10 rounded-xl bg-teal-50 dark:bg-teal-950 flex items-center justify-center">
                      <IconComponent className={`w-4 h-4 sm:w-5 sm:h-5 ${card.accentColor}`} />
                    </div>
                    <div>
                      <span className={`text-[10px] font-mono font-black uppercase tracking-wider ${card.accentColor}`}>
                        {card.role}
                      </span>
                      <h3 className="text-xs sm:text-sm font-bold text-slate-900 dark:text-white leading-tight">
                        {card.title}
                      </h3>
                    </div>
                  </div>

                  <span className="text-[9px] font-bold text-slate-600 dark:text-slate-300 bg-slate-100 dark:bg-slate-800 px-2.5 py-0.5 rounded-full border border-slate-200 dark:border-slate-700">
                    {card.badge}
                  </span>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 pt-3">
                  {card.features.map((feat) => (
                    <div key={feat} className="flex items-start gap-2 text-[11px] sm:text-xs text-slate-600 dark:text-slate-300">
                      <CheckCircle2 className={`w-3.5 h-3.5 mt-0.5 flex-shrink-0 ${card.accentColor}`} />
                      <span>{feat}</span>
                    </div>
                  ))}
                </div>
              </motion.div>
            );
          })}
        </div>

      </div>
    </section>
  );
}
