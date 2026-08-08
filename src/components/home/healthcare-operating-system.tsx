import { motion } from "framer-motion";
import { Cpu, Users, Building2, Handshake, Network, ShieldCheck } from "lucide-react";

export function HealthcareOperatingSystemSection() {
  return (
    <section id="operating-system" className="py-10 sm:py-16 bg-white dark:bg-[#0B1120] text-slate-900 dark:text-white relative overflow-hidden">
      <div className="max-w-[1280px] mx-auto px-4 sm:px-6 relative z-10">
        
        {/* Editorial Section Header */}
        <div className="max-w-2xl mx-auto text-center mb-8 sm:mb-12">
          <div className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-teal-50 dark:bg-teal-950/80 border border-teal-200 dark:border-teal-800 text-[#0F766E] dark:text-teal-400 text-[10px] sm:text-xs font-bold uppercase tracking-wider mb-3">
            <Cpu className="w-3.5 h-3.5" />
            <span>ENTERPRISE NETWORK ARCHITECTURE</span>
          </div>

          <h2 className="text-2xl sm:text-4xl font-black text-slate-900 dark:text-white tracking-tight leading-tight mb-3">
            Curexal is not just another HMS. <br />
            <span className="bg-clip-text text-transparent bg-gradient-to-r from-[#0F766E] via-[#0D9488] to-[#14B8A6]">
              It is the healthcare operating system.
            </span>
          </h2>

          <p className="text-xs sm:text-sm text-slate-600 dark:text-slate-300">
            A multi-tenant coordination platform where isolated facility software modules connect seamlessly to a unified regional healthcare network.
          </p>
        </div>

        {/* Conceptual Architecture Diagram (Compact Scale) */}
        <div className="max-w-3xl mx-auto p-4 sm:p-7 rounded-2xl bg-slate-50 dark:bg-slate-900/90 border border-slate-200 dark:border-slate-800 shadow-md relative">
          
          {/* Top Layer: Central Operating Network Node */}
          <div className="flex flex-col items-center justify-center mb-6">
            <motion.div
              animate={{ y: [-2, 2, -2] }}
              transition={{ duration: 4, repeat: Infinity, ease: "easeInOut" }}
              className="px-4 py-2.5 rounded-xl bg-gradient-to-r from-[#0F766E] to-[#115E59] text-white border border-teal-400/40 flex items-center gap-2.5 shadow-md"
            >
              <div className="w-7 h-7 rounded-lg bg-white/20 flex items-center justify-center">
                <Network className="w-3.5 h-3.5 text-white" />
              </div>
              <div>
                <span className="text-[10px] sm:text-xs font-black uppercase tracking-wider block">CUREXAL OPERATING NETWORK</span>
                <span className="text-[8px] sm:text-[9px] text-teal-200 font-bold">Cross-Tenant Orchestration Layer</span>
              </div>
            </motion.div>

            {/* Connecting Vertical Rays */}
            <div className="w-0.5 h-6 bg-gradient-to-b from-[#0F766E] to-slate-300 dark:to-slate-700 my-1.5" />
          </div>

          {/* Middle Layer: 3 Core Stakeholder Nodes */}
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 mb-6">
            
            {/* Patients Node */}
            <motion.div
              initial={{ opacity: 0, y: 15 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.4 }}
              className="p-3 rounded-xl bg-white dark:bg-slate-950/80 border border-slate-200 dark:border-slate-800 text-center space-y-1.5"
            >
              <div className="w-7 h-7 rounded-lg bg-teal-50 dark:bg-teal-950 text-[#0F766E] dark:text-teal-400 mx-auto flex items-center justify-center">
                <Users className="w-3.5 h-3.5" />
              </div>
              <h3 className="text-xs font-bold text-slate-900 dark:text-white">PATIENTS</h3>
              <p className="text-[10px] text-slate-600 dark:text-slate-400">Unified health vault & marketplace booking</p>
            </motion.div>

            {/* Organizations Node */}
            <motion.div
              initial={{ opacity: 0, y: 15 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.4, delay: 0.08 }}
              className="p-3 rounded-xl bg-white dark:bg-slate-950/80 border border-slate-200 dark:border-slate-800 text-center space-y-1.5"
            >
              <div className="w-7 h-7 rounded-lg bg-teal-50 dark:bg-teal-950 text-[#0F766E] dark:text-teal-400 mx-auto flex items-center justify-center">
                <Building2 className="w-3.5 h-3.5" />
              </div>
              <h3 className="text-xs font-bold text-slate-900 dark:text-white">ORGANIZATIONS</h3>
              <p className="text-[10px] text-slate-600 dark:text-slate-400">Isolated LIMS & EMR workflows with analytics</p>
            </motion.div>

            {/* Partners Node */}
            <motion.div
              initial={{ opacity: 0, y: 15 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.4, delay: 0.16 }}
              className="p-3 rounded-xl bg-white dark:bg-slate-950/80 border border-slate-200 dark:border-slate-800 text-center space-y-1.5"
            >
              <div className="w-7 h-7 rounded-lg bg-teal-50 dark:bg-teal-950 text-[#0F766E] dark:text-teal-400 mx-auto flex items-center justify-center">
                <Handshake className="w-3.5 h-3.5" />
              </div>
              <h3 className="text-xs font-bold text-slate-900 dark:text-white">PARTNERS</h3>
              <p className="text-[10px] text-slate-600 dark:text-slate-400">B2B supplier storefronts & settlements</p>
            </motion.div>

          </div>

          {/* Bottom Layer: Unified Healthcare Network */}
          <div className="flex flex-col items-center justify-center pt-1">
            <div className="w-0.5 h-5 bg-gradient-to-b from-slate-300 dark:from-slate-700 to-[#0F766E] mb-1.5" />
            <div className="px-4 py-2 rounded-full bg-teal-50 dark:bg-teal-950/90 border border-teal-200 dark:border-teal-800 text-[#0F766E] dark:text-teal-300 text-[10px] font-black uppercase tracking-wider flex items-center gap-1.5">
              <ShieldCheck className="w-3.5 h-3.5 text-[#0F766E] dark:text-teal-400" />
              <span>SINGLE CONNECTED HEALTHCARE NETWORK</span>
            </div>
          </div>

        </div>

      </div>
    </section>
  );
}
