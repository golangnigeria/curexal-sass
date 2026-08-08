import { useState } from "react";
import { motion } from "framer-motion";
import { Sparkles, MessageSquare, ArrowRight, User, Stethoscope, ShoppingBag } from "lucide-react";
import { WaitlistModal } from "@/components/waitlist-modal";

const researchFocus = [
  {
    role: "PATIENTS",
    question: "What currently frustrates your healthcare experience?",
    icon: User,
    points: ["Paper slip travel hassles", "Unclear test prices & turnaround times", "Fragmented medical history"],
  },
  {
    role: "CLINICS & LABS",
    question: "Where does healthcare coordination break down for your facility?",
    icon: Stethoscope,
    points: ["Manual paper referrals", "Delayed lab result reception", "Uncollected payment leakage"],
  },
  {
    role: "SUPPLIERS",
    question: "What do you struggle with when supplying healthcare products?",
    icon: ShoppingBag,
    points: ["Manual reagent order processing", "Unpredictable inventory demand", "Reconciliation delays"],
  },
];

export function WaitlistResearchSection() {
  const [modalOpen, setModalOpen] = useState(false);

  return (
    <>
      <section id="help-us-build" className="py-10 sm:py-16 bg-[#F8FAFC] dark:bg-[#0B1120] text-slate-900 dark:text-white relative overflow-hidden">
        <div className="max-w-[1280px] mx-auto px-4 sm:px-6">
          
          {/* Editorial Section Header */}
          <div className="max-w-2xl mx-auto text-center mb-8 sm:mb-12">
            <div className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-teal-50 dark:bg-teal-950/80 border border-teal-200 dark:border-teal-800 text-[#0F766E] dark:text-teal-400 text-[10px] sm:text-xs font-bold uppercase tracking-wider mb-3">
              <MessageSquare className="w-3.5 h-3.5" />
              <span>CUSTOMER DISCOVERY & CO-DESIGN</span>
            </div>

            <h2 className="text-2xl sm:text-4xl font-black tracking-tight text-slate-900 dark:text-white leading-tight mb-3">
              Help us build the healthcare system <br />
              <span className="bg-clip-text text-transparent bg-gradient-to-r from-[#0F766E] via-[#0D9488] to-[#14B8A6]">
                you actually need.
              </span>
            </h2>

            <p className="text-xs sm:text-sm text-slate-600 dark:text-slate-300">
              We are actively gathering operational realities from patients, clinics, laboratories, and suppliers to tailor the Curexal network.
            </p>
          </div>

          {/* Mobile Vertical Scrollable / Desktop 3-Column Grid */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8 max-h-[380px] md:max-h-none overflow-y-auto [scrollbar-width:none] [-ms-overflow-style:none] [&::-webkit-scrollbar]:hidden pr-1">
            {researchFocus.map((item, idx) => {
              const IconComponent = item.icon;
              return (
                <motion.div
                  key={item.role}
                  initial={{ opacity: 0, y: 15 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  viewport={{ once: true }}
                  transition={{ duration: 0.4, delay: idx * 0.08 }}
                  className="p-4 sm:p-5 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 space-y-3 flex flex-col justify-between shadow-xs"
                >
                  <div className="space-y-2.5">
                    <div className="flex items-center justify-between">
                      <span className="text-[10px] font-mono font-black text-[#0F766E] dark:text-teal-400 uppercase tracking-wider">
                        {item.role}
                      </span>
                      <div className="w-7 h-7 rounded-lg bg-teal-50 dark:bg-teal-950 flex items-center justify-center text-[#0F766E] dark:text-teal-400">
                        <IconComponent className="w-3.5 h-3.5" />
                      </div>
                    </div>

                    <h3 className="text-xs sm:text-sm font-bold text-slate-900 dark:text-white leading-snug">
                      "{item.question}"
                    </h3>

                    <ul className="space-y-1.5 pt-2 border-t border-slate-100 dark:border-slate-800">
                      {item.points.map((pt) => (
                        <li key={pt} className="flex items-center gap-1.5 text-[11px] text-slate-600 dark:text-slate-400">
                          <span className="w-1.5 h-1.5 rounded-full bg-[#0F766E] dark:bg-teal-400 flex-shrink-0" />
                          <span>{pt}</span>
                        </li>
                      ))}
                    </ul>
                  </div>
                </motion.div>
              );
            })}
          </div>

          {/* CTA Action */}
          <div className="text-center">
            <button
              onClick={() => setModalOpen(true)}
              className="inline-flex items-center gap-2 px-6 py-3 rounded-xl bg-[#0F766E] hover:bg-[#115E59] text-white text-xs sm:text-sm font-bold transition-all shadow-md cursor-pointer border-0"
            >
              <Sparkles className="w-4 h-4 text-teal-200" />
              <span>Join Early Access & Research</span>
              <ArrowRight className="w-4 h-4" />
            </button>
          </div>

        </div>
      </section>

      {/* Reused Waitlist Modal */}
      <WaitlistModal open={modalOpen} onOpenChange={setModalOpen} />
    </>
  );
}
