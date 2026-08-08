import { motion } from "framer-motion";
import { ShoppingBag, FlaskConical, Stethoscope, Pill, PackageCheck, ArrowRight, ShieldCheck } from "lucide-react";
import { Link } from "react-router-dom";

const categories = [
  {
    title: "Diagnostic Testing & Labs",
    icon: FlaskConical,
    badge: "Public Marketplace",
    description: "Search tests by accession category, compare turnaround times across accredited laboratories, and book appointments directly.",
    tags: ["Full Blood Count", "Lipid Profile", "HbA1c", "LFT & KFT", "ISO 15189 Labs"],
    href: "/marketplace",
  },
  {
    title: "Outpatient Clinics & EMR",
    icon: Stethoscope,
    badge: "Provider Discovery",
    description: "Discover verified general practice clinics and specialist consults with integrated digital E-referral routing.",
    tags: ["General Consults", "Pediatrics", "Cardiology", "Tele-consult", "E-Referrals"],
    href: "/marketplace",
  },
  {
    title: "Pharmacies & Prescriptions",
    icon: Pill,
    badge: "Fulfillment Network",
    description: "Send electronic prescriptions directly to partner dispensing pharmacies with automated stock verification.",
    tags: ["E-Prescription", "Refill Alerts", "Chronic Meds", "Doorstep Dispatch"],
    href: "/marketplace",
  },
  {
    title: "Medical Products & Reagents",
    icon: PackageCheck,
    badge: "B2B Storefront",
    description: "Healthcare facilities order analyzer reagents, diagnostic rapid kits, and clinic consumables directly from verified suppliers.",
    tags: ["Reagents", "Rapid Test Strips", "Specimen Tubes", "B2B Orders"],
    href: "/marketplace",
  },
];

export function MarketplacePreviewCarouselSection() {
  return (
    <section id="marketplace-preview" className="py-10 sm:py-16 bg-white dark:bg-[#0B1120] text-slate-900 dark:text-white relative overflow-hidden">
      <div className="max-w-[1280px] mx-auto px-4 sm:px-6 relative z-10">
        
        {/* Section Header */}
        <div className="flex flex-col sm:flex-row items-start sm:items-end justify-between gap-4 mb-8 sm:mb-12">
          <div className="max-w-2xl">
            <div className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-teal-50 dark:bg-teal-950/80 border border-teal-200 dark:border-teal-800 text-[#0F766E] dark:text-teal-400 text-[10px] sm:text-xs font-bold uppercase tracking-wider mb-3">
              <ShoppingBag className="w-3.5 h-3.5" />
              <span>DISCOVERY & REVENUE ENGINE</span>
            </div>

            <h2 className="text-2xl sm:text-4xl font-black text-slate-900 dark:text-white tracking-tight leading-tight mb-2">
              The Connected Healthcare <br />
              <span className="bg-clip-text text-transparent bg-gradient-to-r from-[#0F766E] via-[#0D9488] to-[#14B8A6]">
                Marketplace Network.
              </span>
            </h2>

            <p className="text-xs sm:text-sm text-slate-600 dark:text-slate-300">
              Patients discover and book accredited healthcare services. Organizations purchase verified medical supplies and expand patient reach.
            </p>
          </div>

          <Link
            to="/marketplace"
            className="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-[#0F766E] hover:bg-[#115E59] text-white text-xs font-bold transition-all shadow-md flex-shrink-0 cursor-pointer"
          >
            <span>Explore Marketplace</span>
            <ArrowRight className="w-3.5 h-3.5" />
          </Link>
        </div>

        {/* Horizontal/Vertical Scrollable Cards Container with Hidden Scrollbars */}
        <div className="flex flex-col sm:flex-row items-stretch gap-3.5 max-h-[380px] sm:max-h-none overflow-y-auto sm:overflow-x-auto [scrollbar-width:none] [-ms-overflow-style:none] [&::-webkit-scrollbar]:hidden pb-2 pr-1">
          {categories.map((cat, idx) => {
            const IconComponent = cat.icon;
            return (
              <motion.div
                key={cat.title}
                initial={{ opacity: 0, x: 20 }}
                whileInView={{ opacity: 1, x: 0 }}
                viewport={{ once: true }}
                transition={{ duration: 0.4, delay: idx * 0.08 }}
                className="w-full sm:w-[300px] p-4 sm:p-5 rounded-2xl bg-slate-50 dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 flex flex-col justify-between flex-shrink-0 hover:border-[#0F766E]/40 transition-colors"
              >
                <div>
                  <div className="flex items-center justify-between mb-3">
                    <div className="w-8 h-8 rounded-xl bg-teal-50 dark:bg-teal-950 flex items-center justify-center text-[#0F766E] dark:text-teal-400">
                      <IconComponent className="w-4 h-4" />
                    </div>
                    <span className="text-[9px] font-bold text-[#0F766E] dark:text-teal-300 bg-teal-50 dark:bg-teal-950 px-2 py-0.5 rounded-full border border-teal-200 dark:border-teal-800">
                      {cat.badge}
                    </span>
                  </div>

                  <h3 className="text-sm font-bold text-slate-900 dark:text-white mb-1.5 leading-tight">
                    {cat.title}
                  </h3>

                  <p className="text-[11px] text-slate-600 dark:text-slate-400 leading-relaxed mb-4">
                    {cat.description}
                  </p>
                </div>

                <div>
                  <div className="flex flex-wrap gap-1 mb-4">
                    {cat.tags.map((tag) => (
                      <span key={tag} className="text-[9px] font-semibold text-slate-600 dark:text-slate-300 bg-white dark:bg-slate-800 px-2 py-0.5 rounded-md border border-slate-200 dark:border-slate-700">
                        {tag}
                      </span>
                    ))}
                  </div>

                  <Link
                    to={cat.href}
                    className="inline-flex items-center gap-1 text-[11px] font-bold text-[#0F766E] dark:text-teal-400 hover:text-teal-500 transition-colors"
                  >
                    <span>Search Categories</span>
                    <ArrowRight className="w-3 h-3" />
                  </Link>
                </div>
              </motion.div>
            );
          })}
        </div>

        {/* Verification Guarantee Footer */}
        <div className="mt-6 flex items-center gap-1.5 text-[10px] text-slate-500 dark:text-slate-400">
          <ShieldCheck className="w-3.5 h-3.5 text-[#0F766E] flex-shrink-0" />
          <span>Strict Provider Verification: All participating labs, clinics, and suppliers undergo accreditation verification before network listing.</span>
        </div>

      </div>
    </section>
  );
}
