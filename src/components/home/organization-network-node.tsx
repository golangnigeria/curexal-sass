import { motion } from "framer-motion";
import { Building2, Stethoscope, FlaskConical, Pill, Radio, Handshake, Network } from "lucide-react";

const nodes = [
  { label: "CLINICS", desc: "Digital Referrals", icon: Stethoscope, position: "top-[2%] left-1/2 -translate-x-1/2" },
  { label: "LABORATORIES", desc: "LIMS Accessioning", icon: FlaskConical, position: "top-[26%] right-[1%]" },
  { label: "PHARMACY", desc: "E-Prescriptions POS", icon: Pill, position: "bottom-[6%] right-[4%]" },
  { label: "IMAGING", desc: "PACS Link Sharing", icon: Radio, position: "bottom-[6%] left-[4%]" },
  { label: "SUPPLIERS", desc: "B2B Reagent Orders", icon: Handshake, position: "top-[26%] left-[1%]" },
];

export function OrganizationNetworkNodeSection() {
  return (
    <section id="organization-network" className="py-10 sm:py-16 bg-[#F8FAFC] dark:bg-[#0B1120] text-slate-900 dark:text-white relative overflow-hidden">
      <div className="max-w-[1280px] mx-auto px-4 sm:px-6 relative z-10">
        
        {/* Editorial Header */}
        <div className="max-w-2xl mb-8 sm:mb-12">
          <h2 className="text-2xl sm:text-4xl font-black text-slate-900 dark:text-white tracking-tight leading-tight mb-3">
            Your facility doesn't operate alone. <br />
            <span className="bg-clip-text text-transparent bg-gradient-to-r from-[#0F766E] via-[#0D9488] to-[#14B8A6]">
              Coordinate across organizational boundaries.
            </span>
          </h2>

          <p className="text-xs sm:text-sm text-slate-600 dark:text-slate-300">
            Curexal connects standalone clinics, ISO laboratories, imaging facilities, pharmacies, and suppliers into a unified healthcare coordination web.
          </p>
        </div>

        {/* Central Organization Network Node Illustration (Compact Scale) */}
        <div className="relative w-full max-w-[340px] xs:max-w-[380px] sm:max-w-[440px] aspect-square mx-auto flex items-center justify-center">
          
          {/* SVG Animated Rays */}
          <svg className="absolute inset-0 w-full h-full pointer-events-none z-0" viewBox="0 0 400 400" fill="none">
            <circle cx="200" cy="200" r="140" stroke="#0F766E" strokeWidth="1.5" strokeDasharray="6 6" className="opacity-30" />
            <circle cx="200" cy="200" r="90" stroke="#0D9488" strokeWidth="1" strokeDasharray="4 4" className="opacity-25" />

            <line x1="200" y1="200" x2="200" y2="55" stroke="#0F766E" strokeWidth="2" strokeDasharray="4 4" className="opacity-40" />
            <line x1="200" y1="200" x2="330" y2="130" stroke="#0D9488" strokeWidth="2" strokeDasharray="4 4" className="opacity-40" />
            <line x1="200" y1="200" x2="300" y2="320" stroke="#14B8A6" strokeWidth="2" strokeDasharray="4 4" className="opacity-40" />
            <line x1="200" y1="200" x2="100" y2="320" stroke="#2DD4BF" strokeWidth="2" strokeDasharray="4 4" className="opacity-40" />
            <line x1="200" y1="200" x2="70" y2="130" stroke="#0F766E" strokeWidth="2" strokeDasharray="4 4" className="opacity-40" />

            <motion.circle r="3.5" fill="#14B8A6" animate={{ cx: [200, 200], cy: [200, 55] }} transition={{ duration: 2.2, repeat: Infinity }} />
            <motion.circle r="3.5" fill="#0D9488" animate={{ cx: [200, 330], cy: [200, 130] }} transition={{ duration: 2.6, repeat: Infinity, delay: 0.3 }} />
            <motion.circle r="3.5" fill="#2DD4BF" animate={{ cx: [200, 300], cy: [200, 320] }} transition={{ duration: 2.4, repeat: Infinity, delay: 0.6 }} />
            <motion.circle r="3.5" fill="#0F766E" animate={{ cx: [200, 100], cy: [200, 320] }} transition={{ duration: 2.1, repeat: Infinity, delay: 0.4 }} />
            <motion.circle r="3.5" fill="#14B8A6" animate={{ cx: [200, 70], cy: [200, 130] }} transition={{ duration: 2.5, repeat: Infinity, delay: 0.7 }} />
          </svg>

          {/* Central Organization Hub */}
          <motion.div
            animate={{ scale: [1, 1.03, 1] }}
            transition={{ duration: 4, repeat: Infinity, ease: "easeInOut" }}
            className="relative z-20 w-28 h-28 sm:w-32 sm:h-32 rounded-full bg-gradient-to-br from-[#0F766E] to-[#115E59] text-white flex flex-col items-center justify-center p-3 text-center border-2 border-teal-400/40 shadow-xl"
          >
            <div className="w-8 h-8 rounded-full bg-white/20 flex items-center justify-center mb-1">
              <Network className="w-4 h-4 text-white" />
            </div>
            <span className="text-[10px] font-black uppercase tracking-wider">YOUR FACILITY</span>
            <span className="text-[8px] text-teal-200 font-extrabold uppercase">CUREXAL NODE</span>
          </motion.div>

          {/* 5 Satellite Nodes */}
          {nodes.map((node, idx) => {
            const IconComponent = node.icon;
            return (
              <motion.div
                key={node.label}
                animate={{ y: [-3, 3, -3] }}
                transition={{ duration: 3.5 + idx * 0.4, repeat: Infinity, ease: "easeInOut" }}
                className={`absolute ${node.position} z-10 flex flex-col items-center cursor-pointer`}
              >
                <div className="px-2.5 py-1 rounded-xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-md flex items-center gap-1.5 backdrop-blur-md">
                  <div className="w-5 h-5 rounded-lg bg-teal-50 dark:bg-teal-950 flex items-center justify-center text-[#0F766E] dark:text-teal-400">
                    <IconComponent className="w-3 h-3" />
                  </div>
                  <span className="text-[10px] font-black text-slate-800 dark:text-white">{node.label}</span>
                </div>
                <span className="mt-0.5 text-[8px] text-[#0F766E] dark:text-teal-300 font-bold bg-teal-50 dark:bg-teal-950 px-1.5 py-0.5 rounded-full border border-teal-200 dark:border-teal-800">
                  {node.desc}
                </span>
              </motion.div>
            );
          })}

        </div>

      </div>
    </section>
  );
}
