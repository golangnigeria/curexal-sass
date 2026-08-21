import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { cn } from "@/lib/utils";
import {
  FlaskConical,
  Stethoscope,
  Pill,
  Radio,
  HeartPulse,
  TrendingUp,
  ArrowRight,
  CheckCircle2,
  AlertTriangle,
  Zap,
  DollarSign,
  Users,
  BarChart3,
  ShoppingBag,
  Target,
  Rocket,
  ChevronRight,
} from "lucide-react";
import { Link } from "react-router-dom";

const providers = [
  {
    id: "labs",
    label: "Laboratories",
    icon: FlaskConical,
    color: "#0F766E",
    challenges: [
      "Low patient volume and poor visibility",
      "Manual reporting and delayed result delivery",
      "Revenue leakage from uncollected payments",
      "Losing referrals to bigger competitors",
    ],
    benefits: [
      "Marketplace exposure to new clinics and patients",
      "Digital referral network to receive referrals automatically",
      "Faster turnaround time with automated workflows",
      "Online payment collection with zero leakage",
      "Multi-branch management from one dashboard",
    ],
    outcomes: [
      { metric: "3×", label: "More referrals" },
      { metric: "40%", label: "Faster TAT" },
      { metric: "↑", label: "Revenue growth" },
    ],
    futureRevenue: ["Marketplace referrals", "Home sample collection", "Corporate wellness", "Subscription diagnostics"],
  },
  {
    id: "clinics",
    label: "Clinics",
    icon: Stethoscope,
    color: "#0D9488",
    challenges: [
      "Patients lost during paper referrals",
      "No visibility after sending patients to labs",
      "Manual follow-up for every result",
      "Poor patient retention and low repeat visits",
    ],
    benefits: [
      "Digital referrals with real-time status tracking",
      "Automatic result receipt from partner labs",
      "Referral commission earnings from marketplace",
      "Better patient experience drives retention",
      "Online consultation and appointment scheduling",
    ],
    outcomes: [
      { metric: "↑", label: "Patient retention" },
      { metric: "₦", label: "Commission income" },
      { metric: "★", label: "Better outcomes" },
    ],
    futureRevenue: ["Consultation bookings", "Follow-up appointments", "Referral commissions", "Corporate contracts"],
  },
  {
    id: "pharmacies",
    label: "Pharmacies",
    icon: Pill,
    color: "#14B8A6",
    challenges: [
      "Prescription errors from manual processing",
      "Stock wastage and expired inventory",
      "Limited repeat customers",
      "No digital presence for ordering",
    ],
    benefits: [
      "Digital prescription verification",
      "Integrated inventory with expiry alerts",
      "Automatic refill reminders for patients",
      "Patient loyalty and retention tools",
      "Connected dispensing workflows",
    ],
    outcomes: [
      { metric: "↑", label: "Repeat business" },
      { metric: "↓", label: "Inventory loss" },
      { metric: "⚡", label: "Faster dispensing" },
    ],
    futureRevenue: ["Online refills", "Prescription subscriptions", "Medication delivery", "Vaccination programs"],
  },
  {
    id: "imaging",
    label: "Imaging Centers",
    icon: Radio,
    color: "#0F766E",
    challenges: [
      "Delayed report delivery to referring doctors",
      "Manual referral processing",
      "Underutilized equipment and idle capacity",
    ],
    benefits: [
      "Digital referrals from connected clinics",
      "Online bookings through the marketplace",
      "Structured reporting and faster turnaround",
      "Marketplace visibility to thousands of patients",
    ],
    outcomes: [
      { metric: "↑", label: "Equipment utilization" },
      { metric: "3×", label: "More referrals" },
      { metric: "₦", label: "Revenue per unit" },
    ],
    futureRevenue: ["Online scheduling", "Priority reporting tiers", "Corporate imaging contracts"],
  },
  {
    id: "patients",
    label: "Patients",
    icon: HeartPulse,
    color: "#2DD4BF",
    challenges: [
      "Repeating tests because reports are lost",
      "Long waiting times at laboratories",
      "No way to compare lab prices or quality",
    ],
    benefits: [
      "Digital health vault with lifetime test history",
      "Online appointment booking and scheduling",
      "Marketplace: find labs, compare prices, and reviews",
      "Digital payments with no cash and no queues",
      "WhatsApp & SMS notifications for results",
    ],
    outcomes: [
      { metric: "⚡", label: "Faster care" },
      { metric: "↓", label: "Lower costs" },
      { metric: "★", label: "Better experience" },
    ],
    futureRevenue: [],
  },
];

const dashboards = [
  {
    label: "Executive",
    icon: DollarSign,
    color: "#0F766E",
    metrics: ["Daily & monthly revenue", "Revenue by branch", "Revenue by department", "Revenue by doctor", "Revenue by test", "Revenue by referral source"],
  },
  {
    label: "Operations",
    icon: BarChart3,
    color: "#0D9488",
    metrics: ["Average turnaround time", "Pending specimens", "Analyzer utilization", "Staff productivity", "Appointment completion rate", "Queue length"],
  },
  {
    label: "Patient Growth",
    icon: Users,
    color: "#14B8A6",
    metrics: ["New vs returning patients", "Referral conversion rate", "Marketplace bookings", "Patient lifetime value", "Patient satisfaction"],
  },
  {
    label: "Financial",
    icon: Target,
    color: "#2DD4BF",
    metrics: ["Outstanding invoices", "Marketplace commissions", "Settlement reports", "Wallet balances", "Collection rate"],
  },
];

const contentVariants = {
  initial: { opacity: 0, y: 12 },
  animate: { opacity: 1, y: 0, transition: { duration: 0.35, ease: "easeOut" as const } },
  exit: { opacity: 0, y: -8, transition: { duration: 0.2 } },
};

const staggerContainer = {
  animate: { transition: { staggerChildren: 0.04 } },
};

const staggerItem = {
  initial: { opacity: 0, x: -8 },
  animate: { opacity: 1, x: 0, transition: { duration: 0.3 } },
};

export function BusinessGrowth() {
  const [activeProvider, setActiveProvider] = useState(0);
  const [activeDashboard, setActiveDashboard] = useState(0);
  const [activeTab, setActiveTab] = useState<"growth" | "intelligence">("growth");
  const current = providers[activeProvider];
  const Icon = current.icon;
  const currentDash = dashboards[activeDashboard];

  return (
    <section id="business-growth" className="py-10 sm:py-16 bg-white dark:bg-[#0B1120]">
      <div className="max-w-[1280px] mx-auto px-4 sm:px-6">

        {/* Header */}
        <motion.div
          initial={{ opacity: 0, y: 16 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.5 }}
          className="max-w-2xl mx-auto text-center mb-8"
        >
          <h2 className="text-2xl sm:text-4xl font-black tracking-tight text-slate-900 dark:text-white leading-tight mb-3">
            Grow Your Healthcare Business<br />
            <span className="bg-clip-text text-transparent bg-gradient-to-r from-[#0F766E] via-[#0D9488] to-[#14B8A6]">
              with Curexal.
            </span>
          </h2>
          <p className="text-xs sm:text-sm text-slate-600 dark:text-slate-300 max-w-xl mx-auto">
            Curexal isn't just software — it's a business growth platform that empowers healthcare providers to increase revenue, reduce costs, and expand into new markets.
          </p>
        </motion.div>

        {/* Tab Switcher */}
        <div className="flex items-center justify-center gap-1.5 mb-8">
          <button
            onClick={() => setActiveTab("growth")}
            className={cn(
              "px-4 py-2 rounded-lg text-xs font-bold transition-all cursor-pointer border-0",
              activeTab === "growth"
                ? "bg-[#0F766E] text-white shadow-sm"
                : "bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 hover:bg-slate-200 dark:hover:bg-slate-700"
            )}
          >
            <span className="flex items-center gap-1.5">
              <TrendingUp className="w-3.5 h-3.5" />
              Growth by Provider
            </span>
          </button>
          <button
            onClick={() => setActiveTab("intelligence")}
            className={cn(
              "px-4 py-2 rounded-lg text-xs font-bold transition-all cursor-pointer border-0",
              activeTab === "intelligence"
                ? "bg-[#0F766E] text-white shadow-sm"
                : "bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 hover:bg-slate-200 dark:hover:bg-slate-700"
            )}
          >
            <span className="flex items-center gap-1.5">
              <BarChart3 className="w-3.5 h-3.5" />
              Business Intelligence
            </span>
          </button>
        </div>

        <AnimatePresence mode="wait">
          {/* ── Tab: Growth by Provider ── */}
          {activeTab === "growth" && (
            <motion.div
              key="growth"
              initial={{ opacity: 0, y: 12 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -8 }}
              transition={{ duration: 0.3 }}
              className="grid grid-cols-1 lg:grid-cols-12 gap-5 lg:gap-8"
            >
              {/* Left: Provider Pills — Horizontal Scroll on Mobile */}
              <div className="lg:col-span-3">
                <div className="flex lg:flex-col gap-1.5 overflow-x-auto lg:overflow-visible [scrollbar-width:none] [-ms-overflow-style:none] [&::-webkit-scrollbar]:hidden pb-1 lg:pb-0">
                  {providers.map((p, idx) => {
                    const PIcon = p.icon;
                    const isActive = idx === activeProvider;
                    return (
                      <motion.button
                        key={p.id}
                        whileTap={{ scale: 0.97 }}
                        onClick={() => setActiveProvider(idx)}
                        className={cn(
                          "flex items-center gap-2 px-3 py-2 rounded-lg border transition-all cursor-pointer bg-transparent flex-shrink-0 text-left",
                          isActive
                            ? "border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 shadow-xs"
                            : "border-transparent hover:bg-slate-50 dark:hover:bg-slate-900/50"
                        )}
                      >
                        <div
                          className="w-7 h-7 rounded-lg flex items-center justify-center flex-shrink-0"
                          style={{ backgroundColor: isActive ? `${p.color}15` : `${p.color}08` }}
                        >
                          <PIcon className="w-3.5 h-3.5" style={{ color: p.color }} />
                        </div>
                        <span className={cn(
                          "text-xs font-bold transition-colors whitespace-nowrap",
                          isActive ? "text-slate-900 dark:text-white" : "text-slate-500 dark:text-slate-400"
                        )}>
                          {p.label}
                        </span>
                        <ChevronRight className={cn(
                          "w-3 h-3 ml-auto transition-all hidden lg:block",
                          isActive ? "text-[#0F766E] opacity-100" : "opacity-0"
                        )} />
                      </motion.button>
                    );
                  })}
                </div>
              </div>

              {/* Right: Provider Detail Panel */}
              <div className="lg:col-span-9">
                <AnimatePresence mode="wait">
                  <motion.div
                    key={activeProvider}
                    variants={contentVariants}
                    initial="initial"
                    animate="animate"
                    exit="exit"
                    className="rounded-2xl border border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900 overflow-hidden shadow-xs"
                  >
                    {/* Header Bar */}
                    <div className="px-4 py-3 border-b border-slate-100 dark:border-slate-800 flex items-center gap-2.5">
                      <div
                        className="w-8 h-8 rounded-lg flex items-center justify-center"
                        style={{ backgroundColor: `${current.color}12` }}
                      >
                        <Icon className="w-4 h-4" style={{ color: current.color }} />
                      </div>
                      <div>
                        <h3 className="text-sm font-bold text-slate-900 dark:text-white">{current.label}</h3>
                        <p className="text-[10px] text-slate-400">Challenges → Curexal Benefits → Outcomes</p>
                      </div>
                    </div>

                    {/* Content Grid */}
                    <div className="p-4 grid grid-cols-1 sm:grid-cols-2 gap-4">

                      {/* Challenges */}
                      <motion.div variants={staggerContainer} initial="initial" animate="animate">
                        <div className="flex items-center gap-1.5 mb-2.5">
                          <AlertTriangle className="w-3.5 h-3.5 text-amber-500" />
                          <span className="text-[10px] font-black text-slate-400 uppercase tracking-widest">
                            Business Challenges
                          </span>
                        </div>
                        <div className="space-y-1.5">
                          {current.challenges.map((c) => (
                            <motion.div key={c} variants={staggerItem} className="flex items-start gap-2">
                              <div className="w-1.5 h-1.5 rounded-full bg-amber-400 mt-1.5 flex-shrink-0" />
                              <span className="text-[11px] text-slate-600 dark:text-slate-300">{c}</span>
                            </motion.div>
                          ))}
                        </div>
                      </motion.div>

                      {/* Benefits */}
                      <motion.div variants={staggerContainer} initial="initial" animate="animate">
                        <div className="flex items-center gap-1.5 mb-2.5">
                          <Zap className="w-3.5 h-3.5 text-[#0F766E]" />
                          <span className="text-[10px] font-black text-slate-400 uppercase tracking-widest">
                            Curexal Benefits
                          </span>
                        </div>
                        <div className="space-y-1.5">
                          {current.benefits.map((b) => (
                            <motion.div key={b} variants={staggerItem} className="flex items-start gap-2">
                              <CheckCircle2 className="w-3.5 h-3.5 mt-0.5 flex-shrink-0" style={{ color: current.color }} />
                              <span className="text-[11px] text-slate-700 dark:text-slate-300">{b}</span>
                            </motion.div>
                          ))}
                        </div>
                      </motion.div>
                    </div>

                    {/* Outcomes Row */}
                    <div className="px-4 pb-4">
                      <div className="flex items-center gap-1.5 mb-2.5">
                        <TrendingUp className="w-3.5 h-3.5 text-emerald-500" />
                        <span className="text-[10px] font-black text-slate-400 uppercase tracking-widest">
                          Business Outcomes
                        </span>
                      </div>
                      <div className="grid grid-cols-3 gap-2">
                        {current.outcomes.map((o, i) => (
                          <motion.div
                            key={o.label}
                            initial={{ opacity: 0, scale: 0.9 }}
                            animate={{ opacity: 1, scale: 1 }}
                            transition={{ duration: 0.3, delay: i * 0.08 }}
                            className="rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 p-2.5 text-center"
                          >
                            <p className="text-lg font-black text-slate-900 dark:text-white leading-none mb-0.5">{o.metric}</p>
                            <p className="text-[9px] text-slate-500 dark:text-slate-400">{o.label}</p>
                          </motion.div>
                        ))}
                      </div>
                    </div>

                    {/* Future Revenue Tags */}
                    {current.futureRevenue.length > 0 && (
                      <div className="px-4 pb-4 pt-2 border-t border-slate-100 dark:border-slate-800">
                        <div className="flex items-center gap-1.5 mb-2">
                          <ShoppingBag className="w-3.5 h-3.5 text-violet-500" />
                          <span className="text-[10px] font-black text-slate-400 uppercase tracking-widest">
                            New Revenue Streams
                          </span>
                        </div>
                        <div className="flex flex-wrap gap-1.5">
                          {current.futureRevenue.map((r, i) => (
                            <motion.span
                              key={r}
                              initial={{ opacity: 0, scale: 0.9 }}
                              animate={{ opacity: 1, scale: 1 }}
                              transition={{ duration: 0.25, delay: i * 0.04 }}
                              className="px-2 py-0.5 rounded-md text-[10px] font-semibold text-slate-600 dark:text-slate-300 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700"
                            >
                              {r}
                            </motion.span>
                          ))}
                        </div>
                      </div>
                    )}
                  </motion.div>
                </AnimatePresence>
              </div>
            </motion.div>
          )}

          {/* ── Tab: Business Intelligence ── */}
          {activeTab === "intelligence" && (
            <motion.div
              key="intelligence"
              initial={{ opacity: 0, y: 12 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -8 }}
              transition={{ duration: 0.3 }}
              className="grid grid-cols-1 lg:grid-cols-12 gap-5 lg:gap-8"
            >
              {/* Left: Dashboard Selector — Horizontal on Mobile */}
              <div className="lg:col-span-3">
                <div className="flex lg:flex-col gap-1.5 overflow-x-auto lg:overflow-visible [scrollbar-width:none] [-ms-overflow-style:none] [&::-webkit-scrollbar]:hidden pb-1 lg:pb-0">
                  {dashboards.map((d, idx) => {
                    const DIcon = d.icon;
                    const isActive = idx === activeDashboard;
                    return (
                      <motion.button
                        key={d.label}
                        whileTap={{ scale: 0.97 }}
                        onClick={() => setActiveDashboard(idx)}
                        className={cn(
                          "flex items-center gap-2 px-3 py-2 rounded-lg border transition-all cursor-pointer bg-transparent flex-shrink-0 text-left",
                          isActive
                            ? "border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 shadow-xs"
                            : "border-transparent hover:bg-slate-50 dark:hover:bg-slate-900/50"
                        )}
                      >
                        <div
                          className="w-7 h-7 rounded-lg flex items-center justify-center flex-shrink-0"
                          style={{ backgroundColor: isActive ? `${d.color}15` : `${d.color}08` }}
                        >
                          <DIcon className="w-3.5 h-3.5" style={{ color: d.color }} />
                        </div>
                        <div className="whitespace-nowrap">
                          <span className={cn(
                            "text-xs font-bold transition-colors block",
                            isActive ? "text-slate-900 dark:text-white" : "text-slate-500 dark:text-slate-400"
                          )}>
                            {d.label}
                          </span>
                          <span className="text-[9px] text-slate-400">{d.metrics.length} metrics</span>
                        </div>
                        <ChevronRight className={cn(
                          "w-3 h-3 ml-auto transition-all hidden lg:block",
                          isActive ? "text-[#0F766E] opacity-100" : "opacity-0"
                        )} />
                      </motion.button>
                    );
                  })}
                </div>
              </div>

              {/* Right: Dashboard Detail Panel */}
              <div className="lg:col-span-9">
                <AnimatePresence mode="wait">
                  <motion.div
                    key={activeDashboard}
                    variants={contentVariants}
                    initial="initial"
                    animate="animate"
                    exit="exit"
                    className="rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 overflow-hidden shadow-xs"
                  >
                    {/* Dashboard Header */}
                    <div
                      className="px-4 py-3 flex items-center gap-2.5"
                      style={{ backgroundColor: `${currentDash.color}06`, borderBottom: `1px solid ${currentDash.color}12` }}
                    >
                      <div
                        className="w-8 h-8 rounded-lg flex items-center justify-center"
                        style={{ backgroundColor: `${currentDash.color}12` }}
                      >
                        {(() => { const DI = currentDash.icon; return <DI className="w-4 h-4" style={{ color: currentDash.color }} />; })()}
                      </div>
                      <div>
                        <h3 className="text-sm font-bold text-slate-900 dark:text-white">{currentDash.label} Dashboard</h3>
                        <p className="text-[10px] text-slate-400">Real-time business intelligence</p>
                      </div>
                    </div>

                    {/* Metric Grid */}
                    <motion.div className="p-4" variants={staggerContainer} initial="initial" animate="animate">
                      <p className="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-3">
                        Key Metrics Available
                      </p>
                      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                        {currentDash.metrics.map((m) => (
                          <motion.div
                            key={m}
                            variants={staggerItem}
                            className="flex items-center gap-2.5 p-2.5 rounded-lg border border-slate-100 dark:border-slate-800 bg-slate-50 dark:bg-slate-800/50 hover:border-slate-200 dark:hover:border-slate-700 transition-colors"
                          >
                            <div
                              className="w-6 h-6 rounded-md flex items-center justify-center flex-shrink-0"
                              style={{ backgroundColor: `${currentDash.color}10` }}
                            >
                              <BarChart3 className="w-3 h-3" style={{ color: currentDash.color }} />
                            </div>
                            <span className="text-[11px] font-medium text-slate-700 dark:text-slate-300">{m}</span>
                          </motion.div>
                        ))}
                      </div>
                    </motion.div>

                    {/* CTA Footer */}
                    <div className="px-4 pb-4 pt-2 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between">
                      <p className="text-[10px] text-slate-400">
                        Make better business decisions with real-time data.
                      </p>
                      <Link to="/book-demo">
                        <button className="flex items-center gap-1 text-[11px] font-bold text-[#0F766E] hover:text-[#115E59] transition-colors bg-transparent border-0 cursor-pointer">
                          See it live
                          <ArrowRight className="w-3 h-3" />
                        </button>
                      </Link>
                    </div>
                  </motion.div>
                </AnimatePresence>
              </div>
            </motion.div>
          )}
        </AnimatePresence>


      </div>
    </section>
  );
}
