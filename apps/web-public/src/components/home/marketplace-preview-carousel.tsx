import { motion } from "framer-motion";
import { FlaskConical, Stethoscope, Radio, ArrowRight, ShieldCheck, Share2, Send, CheckCircle2, FileCheck } from "lucide-react";
import { Link } from "react-router-dom";

const referralPillars = [
  {
    title: "Electronic Doctor-to-Lab Requisitions",
    icon: Send,
    badge: "Paperless Orders",
    description: "Doctors order diagnostic tests directly from their consultation canvas. Tests arrive instantly at partner laboratory accession benches with zero handwriting errors.",
    tags: ["Instant Requisition", "No Missing Forms", "Direct LIS Ingest", "Bi-directional Sync"],
    href: "/marketplace",
  },
  {
    title: "Pathology & Laboratory Network",
    icon: FlaskConical,
    badge: "Accredited Labs",
    description: "Partner medical laboratories receive electronic test orders from regional clinics, process specimens on interfaced analyzers, and deliver verified results automatically.",
    tags: ["Analyzer Telemetry", "Pathologist Sign-off", "Automated TAT", "ISO 15189 Aligned"],
    href: "/marketplace",
  },
  {
    title: "Radiology & DICOM Image Sharing",
    icon: Radio,
    badge: "Digital Modalities",
    description: "Refer patients for specialized Ultrasounds, X-Rays, CT, and MRI scans. Radiologists sign structured reports and share high-res DICOM links with zero physical film costs.",
    tags: ["Web DICOM PACS", "Zero Film Waste", "Structured Reports", "Direct Doctor Chart Sync"],
    href: "/marketplace",
  },
  {
    title: "Zero-Leakage Patient Result Vault",
    icon: FileCheck,
    badge: "Encrypted Delivery",
    description: "Verified results flow directly back into the referring clinic's patient chart while simultaneously dispatching encrypted PDF reports to the patient via WhatsApp & SMS.",
    tags: ["WhatsApp PDF", "EMR Chart Update", "Tamper-Proof Signatures", "HIPAA & NDPR"],
    href: "/marketplace",
  },
];

export function MarketplacePreviewCarouselSection() {
  return (
    <section id="referral-highway" className="py-12 sm:py-20 bg-white dark:bg-[#0B1120] text-slate-900 dark:text-white relative overflow-hidden">
      <div className="max-w-[1280px] mx-auto px-4 sm:px-6 relative z-10">
        
        {/* Section Header */}
        <div className="flex flex-col sm:flex-row items-start sm:items-end justify-between gap-4 mb-8 sm:mb-12">
          <div className="max-w-2xl">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-teal-500/10 border border-teal-500/20 text-[#0F766E] dark:text-teal-400 text-xs font-bold uppercase tracking-wider mb-3">
              <Share2 className="w-3.5 h-3.5" />
              <span>B2B Interoperability Layer</span>
            </div>

            <h2 className="text-2xl sm:text-4xl font-black text-slate-900 dark:text-white tracking-tight leading-tight mb-2">
              The Connected Diagnostic & <br />
              <span className="bg-clip-text text-transparent bg-gradient-to-r from-[#0F766E] via-[#0D9488] to-[#14B8A6]">
                Referral Highway.
              </span>
            </h2>

            <p className="text-xs sm:text-sm text-slate-600 dark:text-slate-300">
              Stop losing patient referrals to illegible paper slips and lost phone calls. Connect your clinic, laboratory, and imaging center on a unified digital referral highway.
            </p>
          </div>

          <Link
            to="/marketplace"
            className="inline-flex items-center gap-1.5 px-4 py-2.5 rounded-xl bg-[#0F766E] hover:bg-[#115E59] text-white text-xs font-bold transition-all shadow-md flex-shrink-0 cursor-pointer"
          >
            <span>Explore Referral Network</span>
            <ArrowRight className="w-3.5 h-3.5" />
          </Link>
        </div>

        {/* 4 Connected Cards Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {referralPillars.map((cat, idx) => {
            const IconComponent = cat.icon;
            return (
              <motion.div
                key={cat.title}
                initial={{ opacity: 0, y: 15 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ duration: 0.35, delay: idx * 0.08 }}
                className="p-5 rounded-2xl bg-slate-50 dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 flex flex-col justify-between hover:border-[#0F766E]/50 transition-colors shadow-xs"
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

                  <h3 className="text-sm font-bold text-slate-900 dark:text-white mb-2 leading-tight">
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
                    <span>View Network Workflow</span>
                    <ArrowRight className="w-3 h-3" />
                  </Link>
                </div>
              </motion.div>
            );
          })}
        </div>

        {/* Verification Guarantee Footer */}
        <div className="mt-8 flex items-center gap-2 text-xs text-slate-500 dark:text-slate-400 bg-teal-50/50 dark:bg-teal-950/20 p-3.5 rounded-xl border border-teal-200/60 dark:border-teal-800/40">
          <ShieldCheck className="w-4 h-4 text-[#0F766E] flex-shrink-0" />
          <span>Accredited Provider Interoperability: All participating clinics, diagnostic laboratories, and imaging centers undergo regulatory credential verification before joining the live exchange.</span>
        </div>

      </div>
    </section>
  );
}
