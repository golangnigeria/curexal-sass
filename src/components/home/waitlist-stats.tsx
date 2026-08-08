import { useEffect, useState } from "react";
import { fetchWaitlistStats, type WaitlistAggregateStats } from "@/lib/supabase";
import { Users, Stethoscope, User, ShoppingBag, RefreshCw, ShieldCheck, PieChart as PieIcon } from "lucide-react";
import { ResponsiveContainer, PieChart, Pie, Cell, Tooltip } from "recharts";
import { motion } from "framer-motion";

export function WaitlistStats() {
  const [stats, setStats] = useState<WaitlistAggregateStats>({
    totalMembers: 0,
    patientsCount: 0,
    organizationsCount: 0,
    suppliersCount: 0,
    loading: true,
    error: false,
  });

  const loadStats = async () => {
    setStats((prev) => ({ ...prev, loading: true }));
    const data = await fetchWaitlistStats();
    setStats(data);
  };

  useEffect(() => {
    loadStats();
  }, []);

  const total = stats.totalMembers || 1;
  const patientsPct = Math.round((stats.patientsCount / total) * 100) || 0;
  const orgsPct = Math.round((stats.organizationsCount / total) * 100) || 0;
  const suppliersPct = Math.round((stats.suppliersCount / total) * 100) || 0;

  const chartData = [
    { name: "Patients", value: stats.patientsCount || 1, color: "#14B8A6" },
    { name: "Orgs & Labs", value: stats.organizationsCount || 1, color: "#0F766E" },
    { name: "Suppliers", value: stats.suppliersCount || 1, color: "#38BDF8" },
  ];
  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5, ease: "easeOut" }}
      className="w-full bg-slate-900/95 text-white rounded-2xl p-3.5 sm:p-5 border border-slate-800 shadow-xl relative overflow-hidden my-4 box-border font-inter backdrop-blur-xl"
    >
      {/* Background Radial Glow */}
      <div className="absolute top-0 right-0 w-64 h-64 bg-teal-500/10 rounded-full blur-3xl pointer-events-none" />

      <div className="relative z-10 space-y-3.5">
        
        {/* Compact Header Bar */}
        <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-2 pb-1.5 border-b border-slate-800/80">
          <div className="flex items-center gap-2">
            <div className="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 text-[9px] font-extrabold uppercase tracking-wider">
              <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse" />
              <span>LIVE TELEMETRY</span>
            </div>
            <h3 className="text-xs sm:text-sm font-extrabold text-white tracking-tight">
              Curexal Network Early Access
            </h3>
          </div>

          <button
            onClick={loadStats}
            className="text-slate-400 hover:text-teal-300 text-[10px] font-semibold inline-flex items-center gap-1 cursor-pointer bg-slate-800/80 hover:bg-slate-800 border border-slate-700/80 px-2 py-0.5 rounded-lg transition-all self-end sm:self-auto"
          >
            <RefreshCw className={`w-3 h-3 ${stats.loading ? "animate-spin text-teal-400" : ""}`} />
            <span>Refresh</span>
          </button>
        </div>

        {/* Compact Dashboard Content */}
        {stats.loading ? (
          <div className="grid grid-cols-1 md:grid-cols-12 gap-3 animate-pulse">
            <div className="md:col-span-7 grid grid-cols-2 gap-2">
              {[1, 2, 3, 4].map((i) => (
                <div key={i} className="p-3 rounded-xl bg-slate-800/50 border border-slate-700/50 h-16" />
              ))}
            </div>
            <div className="md:col-span-5 rounded-xl bg-slate-800/50 border border-slate-700/50 h-36" />
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-12 gap-3.5 items-center">
            
            {/* Left 7 Cols: Compact Metric Cards */}
            <div className="md:col-span-7 grid grid-cols-2 gap-2.5">
              
              {/* Total Members */}
              <motion.div
                whileHover={{ y: -2 }}
                className="p-2.5 sm:p-3 rounded-xl bg-gradient-to-br from-teal-950/60 to-slate-900 border border-teal-500/30 space-y-1 min-w-0"
              >
                <div className="flex items-center justify-between text-teal-400">
                  <span className="text-[9px] font-extrabold uppercase tracking-wider truncate">Total Waitlist</span>
                  <Users className="w-3.5 h-3.5 flex-shrink-0" />
                </div>
                <div>
                  <span className="text-base sm:text-xl font-black text-white block truncate">
                    {stats.totalMembers.toLocaleString()}
                  </span>
                  <span className="text-[9px] font-bold text-teal-300">100% Total Growth</span>
                </div>
                <div className="w-full bg-teal-950 h-1 rounded-full overflow-hidden">
                  <motion.div
                    initial={{ width: 0 }}
                    animate={{ width: "100%" }}
                    transition={{ duration: 0.8 }}
                    className="bg-[#14B8A6] h-full"
                  />
                </div>
              </motion.div>

              {/* Patients */}
              <motion.div
                whileHover={{ y: -2 }}
                className="p-2.5 sm:p-3 rounded-xl bg-slate-800/40 border border-slate-700/60 space-y-1 min-w-0"
              >
                <div className="flex items-center justify-between text-slate-400">
                  <span className="text-[9px] font-extrabold uppercase tracking-wider truncate">Patients</span>
                  <User className="w-3.5 h-3.5 text-teal-400 flex-shrink-0" />
                </div>
                <div>
                  <span className="text-base sm:text-xl font-black text-white block truncate">
                    {stats.patientsCount.toLocaleString()}
                  </span>
                  <span className="text-[9px] font-bold text-slate-400">{patientsPct}% of total</span>
                </div>
                <div className="w-full bg-slate-800 h-1 rounded-full overflow-hidden">
                  <motion.div
                    initial={{ width: 0 }}
                    animate={{ width: `${patientsPct}%` }}
                    transition={{ duration: 0.8 }}
                    className="bg-teal-400 h-full"
                  />
                </div>
              </motion.div>

              {/* Orgs & Labs */}
              <motion.div
                whileHover={{ y: -2 }}
                className="p-2.5 sm:p-3 rounded-xl bg-slate-800/40 border border-slate-700/60 space-y-1 min-w-0"
              >
                <div className="flex items-center justify-between text-slate-400">
                  <span className="text-[9px] font-extrabold uppercase tracking-wider truncate">Orgs &amp; Labs</span>
                  <Stethoscope className="w-3.5 h-3.5 text-teal-400 flex-shrink-0" />
                </div>
                <div>
                  <span className="text-base sm:text-xl font-black text-white block truncate">
                    {stats.organizationsCount.toLocaleString()}
                  </span>
                  <span className="text-[9px] font-bold text-slate-400">{orgsPct}% of total</span>
                </div>
                <div className="w-full bg-slate-800 h-1 rounded-full overflow-hidden">
                  <motion.div
                    initial={{ width: 0 }}
                    animate={{ width: `${orgsPct}%` }}
                    transition={{ duration: 0.8 }}
                    className="bg-[#0F766E] h-full"
                  />
                </div>
              </motion.div>

              {/* Suppliers */}
              <motion.div
                whileHover={{ y: -2 }}
                className="p-2.5 sm:p-3 rounded-xl bg-slate-800/40 border border-slate-700/60 space-y-1 min-w-0"
              >
                <div className="flex items-center justify-between text-slate-400">
                  <span className="text-[9px] font-extrabold uppercase tracking-wider truncate">Suppliers</span>
                  <ShoppingBag className="w-3.5 h-3.5 text-[#38BDF8] flex-shrink-0" />
                </div>
                <div>
                  <span className="text-base sm:text-xl font-black text-white block truncate">
                    {stats.suppliersCount.toLocaleString()}
                  </span>
                  <span className="text-[9px] font-bold text-slate-400">{suppliersPct}% of total</span>
                </div>
                <div className="w-full bg-slate-800 h-1 rounded-full overflow-hidden">
                  <motion.div
                    initial={{ width: 0 }}
                    animate={{ width: `${suppliersPct}%` }}
                    transition={{ duration: 0.8 }}
                    className="bg-[#38BDF8] h-full"
                  />
                </div>
              </motion.div>

            </div>

            {/* Right 5 Cols: Compact Recharts Donut */}
            <div className="md:col-span-5 bg-slate-800/40 border border-slate-700/60 rounded-xl p-3 flex flex-col items-center justify-center relative min-h-[140px]">
              <div className="flex items-center gap-1 text-[11px] font-bold text-slate-300 self-start mb-0.5">
                <PieIcon className="w-3.5 h-3.5 text-teal-400" />
                <span>Persona Composition</span>
              </div>

              <div className="w-full h-32 relative flex items-center justify-center">
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={chartData}
                      cx="50%"
                      cy="50%"
                      innerRadius={32}
                      outerRadius={48}
                      paddingAngle={3}
                      dataKey="value"
                    >
                      {chartData.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={entry.color} stroke="none" />
                      ))}
                    </Pie>
                    <Tooltip
                      contentStyle={{
                        backgroundColor: "#0B1120",
                        borderColor: "#334155",
                        borderRadius: "10px",
                        color: "#fff",
                        fontSize: "11px",
                        fontWeight: "600",
                        padding: "4px 8px",
                      }}
                    />
                  </PieChart>
                </ResponsiveContainer>

                {/* Donut Center Count */}
                <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
                  <span className="text-base font-black text-white leading-none">
                    {stats.totalMembers}
                  </span>
                  <span className="text-[8px] font-bold text-slate-400 uppercase tracking-tight mt-0.5">
                    Members
                  </span>
                </div>
              </div>

              {/* Chart Legend */}
              <div className="flex items-center justify-center gap-2 text-[9px] font-bold text-slate-300 pt-0.5 w-full flex-wrap">
                <div className="flex items-center gap-1">
                  <span className="w-2 h-2 rounded-full bg-[#14B8A6]" />
                  <span>Patients ({patientsPct}%)</span>
                </div>
                <div className="flex items-center gap-1">
                  <span className="w-2 h-2 rounded-full bg-[#0F766E]" />
                  <span>Orgs ({orgsPct}%)</span>
                </div>
                <div className="flex items-center gap-1">
                  <span className="w-2 h-2 rounded-full bg-[#38BDF8]" />
                  <span>Suppliers ({suppliersPct}%)</span>
                </div>
              </div>
            </div>

          </div>
        )}

        {/* Footer Note */}
        <div className="flex items-center justify-between gap-2 text-[9px] text-slate-400 pt-1 border-t border-slate-800/80">
          <span className="flex items-center gap-1">
            <ShieldCheck className="w-3 h-3 text-[#0F766E] flex-shrink-0" />
            <span>Verified Supabase Telemetry. Aggregate counts only. Zero personal data exposed.</span>
          </span>
        </div>

      </div>
    </motion.div>
  );
}
