import { motion } from "framer-motion";
import { SignalHigh, Smartphone, MessageSquare, CreditCard, Globe2 } from "lucide-react";

const pillars = [
  {
    icon: SignalHigh,
    title: "Low-Bandwidth Optimization",
    description: "Designed for variable cellular connectivity with light payload transmission and instant cached state.",
  },
  {
    icon: Smartphone,
    title: "Offline-Capable Workflows",
    description: "Accession specimens and log patient notes offline. Records queue safely and sync when connection restores.",
  },
  {
    icon: MessageSquare,
    title: "WhatsApp & SMS Delivery",
    description: "Deliver test status updates and secure PDF result links directly to channels African patients already use daily.",
  },
  {
    icon: CreditCard,
    title: "African Payment Rails",
    description: "Native support for local card settlement, mobile transfers, and patient wallet options to eliminate payment friction.",
  },
];

export function AfricaRealityInfrastructureSection() {
  return (
    <section id="africa-infrastructure" className="py-10 sm:py-16 bg-[#F8FAFC] dark:bg-[#0B1120] text-slate-900 dark:text-white relative overflow-hidden">
      <div className="max-w-[1280px] mx-auto px-4 sm:px-6">
        
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-12 items-center">
          
          {/* Left Column: Editorial Philosophy */}
          <div className="lg:col-span-5 space-y-4">
            <h2 className="text-2xl sm:text-4xl font-black text-slate-900 dark:text-white tracking-tight leading-tight">
              Healthcare infrastructure isn't always predictable. <br />
              <span className="bg-clip-text text-transparent bg-gradient-to-r from-[#0F766E] via-[#0D9488] to-[#14B8A6]">
                Curexal is built for reality.
              </span>
            </h2>

            <p className="text-xs sm:text-sm text-slate-600 dark:text-slate-300 leading-relaxed">
              Western healthcare software assumes zero-latency fiber internet, steady grid power, and centralized insurance systems. Curexal is architected specifically around the operational realities of healthcare in African markets.
            </p>
          </div>

          {/* Right Column: 4 Core Infrastructure Pillars */}
          <div className="lg:col-span-7 grid grid-cols-1 sm:grid-cols-2 gap-3 max-h-[360px] sm:max-h-none overflow-y-auto [scrollbar-width:none] [-ms-overflow-style:none] [&::-webkit-scrollbar]:hidden pr-1">
            {pillars.map((item, idx) => {
              const IconComponent = item.icon;
              return (
                <motion.div
                  key={item.title}
                  initial={{ opacity: 0, y: 15 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  viewport={{ once: true }}
                  transition={{ duration: 0.4, delay: idx * 0.08 }}
                  className="p-4 rounded-xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 flex flex-col justify-between space-y-2 hover:border-[#0F766E]/40 transition-colors shadow-xs"
                >
                  <div className="w-8 h-8 rounded-lg bg-teal-50 dark:bg-teal-950 flex items-center justify-center text-[#0F766E] dark:text-teal-400">
                    <IconComponent className="w-4 h-4" />
                  </div>

                  <div>
                    <h3 className="text-xs font-bold text-slate-900 dark:text-white mb-1">
                      {item.title}
                    </h3>
                    <p className="text-[11px] text-slate-600 dark:text-slate-400 leading-relaxed">
                      {item.description}
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
