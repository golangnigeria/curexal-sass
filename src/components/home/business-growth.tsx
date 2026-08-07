import { useState } from "react";
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
      "Limited marketing reach beyond physical location",
    ],
    benefits: [
      "Marketplace exposure to new clinics and patients",
      "Digital referral network to receive referrals automatically",
      "Faster turnaround time with automated workflows",
      "Online payment collection with zero leakage",
      "Home sample collection scheduling",
      "Multi-branch management from one dashboard",
      "Business analytics and ISO quality monitoring",
    ],
    outcomes: [
      { metric: "3×", label: "More referrals from connected clinics" },
      { metric: "40%", label: "Faster turnaround time" },
      { metric: "↑", label: "Higher monthly revenue" },
      { metric: "★", label: "Better patient satisfaction" },
    ],
    futureRevenue: [
      "Marketplace referrals",
      "Home sample collection",
      "Corporate wellness packages",
      "Preventive health packages",
      "Subscription diagnostics",
    ],
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
      "Shared diagnostic history across the network",
      "Online consultation and appointment scheduling",
      "Better patient experience drives retention",
    ],
    outcomes: [
      { metric: "↑", label: "Higher patient retention" },
      { metric: "₦", label: "Additional referral commission income" },
      { metric: "★", label: "Better clinical outcomes" },
      { metric: "🤝", label: "Increased patient trust" },
    ],
    futureRevenue: [
      "Consultation bookings",
      "Follow-up appointments",
      "Referral commissions",
      "Chronic care programs",
      "Corporate healthcare contracts",
    ],
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
      { metric: "↑", label: "Increased repeat business" },
      { metric: "↓", label: "Reduced inventory loss" },
      { metric: "⚡", label: "Faster dispensing" },
      { metric: "★", label: "Patient loyalty" },
    ],
    futureRevenue: [
      "Online refills",
      "Prescription subscriptions",
      "Medication delivery",
      "Vaccination programs",
    ],
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
      "Limited discoverability for patients",
    ],
    benefits: [
      "Digital referrals from connected clinics",
      "Online bookings through the marketplace",
      "Structured reporting and faster turnaround",
      "Marketplace visibility to thousands of patients",
      "Automated scheduling and capacity management",
    ],
    outcomes: [
      { metric: "↑", label: "Increased equipment utilization" },
      { metric: "3×", label: "More referrals from network" },
      { metric: "₦", label: "Higher revenue per machine" },
      { metric: "⚡", label: "Faster report delivery" },
    ],
    futureRevenue: [
      "Online scheduling",
      "Priority reporting tiers",
      "Specialist reporting",
      "Corporate imaging contracts",
    ],
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
      "Multiple hospital visits for one diagnosis",
    ],
    benefits: [
      "Digital health vault with lifetime test history",
      "Online appointment booking and scheduling",
      "Marketplace: find labs, compare prices, and reviews",
      "Digital payments with no cash and no queues",
      "WhatsApp & SMS notifications for results",
      "Verified PDF reports accessible anywhere",
    ],
    outcomes: [
      { metric: "⚡", label: "Faster care, fewer visits" },
      { metric: "↓", label: "Lower costs through comparison" },
      { metric: "★", label: "Better healthcare experience" },
      { metric: "📱", label: "Results on any device" },
    ],
    futureRevenue: [],
  },
];

const dashboards = [
  {
    label: "Executive",
    icon: DollarSign,
    color: "#0F766E",
    metrics: ["Daily & monthly revenue", "Revenue by branch", "Revenue by department", "Revenue by doctor", "Revenue by test", "Revenue by referral source", "Revenue by payment method"],
  },
  {
    label: "Operations",
    icon: BarChart3,
    color: "#0D9488",
    metrics: ["Average turnaround time", "Pending specimens", "Analyzer utilization", "Staff productivity", "Appointment completion rate", "No-show rate", "Queue length"],
  },
  {
    label: "Patient Growth",
    icon: Users,
    color: "#14B8A6",
    metrics: ["New vs returning patients", "Referral conversion rate", "Marketplace bookings", "Repeat laboratory tests", "Patient lifetime value", "Patient satisfaction"],
  },
  {
    label: "Financial",
    icon: Target,
    color: "#2DD4BF",
    metrics: ["Outstanding invoices", "Marketplace commissions", "Settlement reports", "Wallet balances", "Profit margins", "Collection rate"],
  },
];

export function BusinessGrowth() {
  const [activeProvider, setActiveProvider] = useState(0);
  const [activeDashboard, setActiveDashboard] = useState(0);
  const [activeTab, setActiveTab] = useState<"growth" | "intelligence">("growth");
  const current = providers[activeProvider];
  const Icon = current.icon;
  const currentDash = dashboards[activeDashboard];

  return (
    <section
      id="business-growth"
      className="section-padding bg-white dark:bg-[#0B1120]"
    >
      <div className="max-w-[1280px] mx-auto px-6">

        {/* Header */}
        <div className="text-center max-w-2xl mx-auto mb-10">
          <div className="inline-flex items-center gap-2 px-3 py-1.5 mb-5 rounded-full border border-[#0F766E]/20 bg-[#F0FDFA] dark:bg-[#0F766E]/10">
            <Rocket className="h-3.5 w-3.5 text-[#0F766E]" />
            <span className="text-xs font-semibold text-[#0F766E] tracking-wide">
              Business Growth Platform
            </span>
          </div>
          <h2 className="text-section text-gray-900 dark:text-white mb-4">
            Grow Your Healthcare Business<br />with Curexal.
          </h2>
          <p className="text-body text-gray-500 dark:text-gray-400 max-w-xl mx-auto">
            Curexal isn't just software that digitizes operations, it's a business growth platform that empowers healthcare providers to increase revenue, reduce costs, and expand into new markets.
          </p>
        </div>

        {/* Tab switcher */}
        <div className="flex items-center justify-center gap-2 mb-10">
          <button
            onClick={() => setActiveTab("growth")}
            className={cn(
              "px-5 py-2.5 rounded-[10px] text-sm font-semibold transition-all cursor-pointer border-0",
              activeTab === "growth"
                ? "bg-[#0F766E] text-white shadow-sm"
                : "bg-gray-100 dark:bg-[#1F2937] text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-[#374151]"
            )}
          >
            <span className="flex items-center gap-2">
              <TrendingUp className="w-4 h-4" />
              Growth by Provider
            </span>
          </button>
          <button
            onClick={() => setActiveTab("intelligence")}
            className={cn(
              "px-5 py-2.5 rounded-[10px] text-sm font-semibold transition-all cursor-pointer border-0",
              activeTab === "intelligence"
                ? "bg-[#0F766E] text-white shadow-sm"
                : "bg-gray-100 dark:bg-[#1F2937] text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-[#374151]"
            )}
          >
            <span className="flex items-center gap-2">
              <BarChart3 className="w-4 h-4" />
              Business Intelligence
            </span>
          </button>
        </div>

        {/* ── Tab: Growth by Provider ── */}
        {activeTab === "growth" && (
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-10 animate-fade-in">

            {/* Left: Provider selector */}
            <div className="lg:col-span-4">
              <div className="flex flex-col gap-1.5">
                {providers.map((p, idx) => {
                  const PIcon = p.icon;
                  const isActive = idx === activeProvider;
                  return (
                    <button
                      key={p.id}
                      onClick={() => setActiveProvider(idx)}
                      className={cn(
                        "w-full text-left px-4 py-3 rounded-[12px] border transition-all duration-300 cursor-pointer bg-transparent flex items-center gap-3",
                        isActive
                          ? "border-gray-200 dark:border-[#374151] bg-[#F8FAFC] dark:bg-[#111827] shadow-sm"
                          : "border-transparent hover:border-gray-200 dark:hover:border-[#374151] hover:bg-gray-50 dark:hover:bg-[#111827]/40"
                      )}
                    >
                      <div
                        className="w-9 h-9 rounded-[10px] flex items-center justify-center flex-shrink-0"
                        style={{ backgroundColor: isActive ? `${p.color}15` : `${p.color}08` }}
                      >
                        <PIcon className="w-4.5 h-4.5" style={{ color: p.color }} />
                      </div>
                      <span className={cn(
                        "text-sm font-semibold transition-colors",
                        isActive ? "text-gray-900 dark:text-white" : "text-gray-500 dark:text-gray-400"
                      )}>
                        {p.label}
                      </span>
                      <ArrowRight className={cn(
                        "w-3.5 h-3.5 ml-auto transition-all",
                        isActive ? "text-[#0F766E] opacity-100" : "opacity-0"
                      )} />
                    </button>
                  );
                })}
              </div>
            </div>

            {/* Right: Provider detail */}
            <div className="lg:col-span-8">
              <div
                key={activeProvider}
                className="rounded-[16px] border border-gray-200 dark:border-[#1F2937] bg-[#F8FAFC] dark:bg-[#111827] overflow-hidden animate-fade-in"
                style={{ boxShadow: "0 1px 2px rgba(0,0,0,0.04), 0 8px 32px rgba(0,0,0,0.06)" }}
              >
                {/* Header */}
                <div className="px-6 py-4 border-b border-gray-100 dark:border-[#1F2937] flex items-center gap-3">
                  <div
                    className="w-10 h-10 rounded-[10px] flex items-center justify-center"
                    style={{ backgroundColor: `${current.color}12` }}
                  >
                    <Icon className="w-5 h-5" style={{ color: current.color }} />
                  </div>
                  <div>
                    <h3 className="text-lg font-bold text-gray-900 dark:text-white">{current.label}</h3>
                    <p className="text-xs text-gray-400 dark:text-gray-500">Business challenges → Curexal benefits → Outcomes</p>
                  </div>
                </div>

                <div className="p-6 grid grid-cols-1 md:grid-cols-2 gap-6">

                  {/* Challenges */}
                  <div>
                    <div className="flex items-center gap-2 mb-3">
                      <AlertTriangle className="w-4 h-4 text-amber-500" />
                      <span className="text-xs font-bold text-gray-400 dark:text-gray-500 uppercase tracking-widest">
                        Business Challenges
                      </span>
                    </div>
                    <div className="space-y-2">
                      {current.challenges.map((c, i) => (
                        <div
                          key={c}
                          className="flex items-start gap-2.5 animate-fade-up"
                          style={{ animationDelay: `${i * 60}ms` }}
                        >
                          <div className="w-1.5 h-1.5 rounded-full bg-amber-400 mt-1.5 flex-shrink-0" />
                          <span className="text-[13px] text-gray-600 dark:text-gray-300">{c}</span>
                        </div>
                      ))}
                    </div>
                  </div>

                  {/* Curexal Benefits */}
                  <div>
                    <div className="flex items-center gap-2 mb-3">
                      <Zap className="w-4 h-4 text-[#0F766E]" />
                      <span className="text-xs font-bold text-gray-400 dark:text-gray-500 uppercase tracking-widest">
                        Curexal Benefits
                      </span>
                    </div>
                    <div className="space-y-2">
                      {current.benefits.map((b, i) => (
                        <div
                          key={b}
                          className="flex items-start gap-2.5 animate-fade-up"
                          style={{ animationDelay: `${i * 60}ms` }}
                        >
                          <CheckCircle2 className="w-3.5 h-3.5 mt-0.5 flex-shrink-0" style={{ color: current.color }} />
                          <span className="text-[13px] text-gray-700 dark:text-gray-300">{b}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>

                {/* Business Outcomes */}
                <div className="px-6 pb-5">
                  <div className="flex items-center gap-2 mb-3">
                    <TrendingUp className="w-4 h-4 text-emerald-500" />
                    <span className="text-xs font-bold text-gray-400 dark:text-gray-500 uppercase tracking-widest">
                      Business Outcomes
                    </span>
                  </div>
                  <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
                    {current.outcomes.map((o, i) => (
                      <div
                        key={o.label}
                        className="rounded-[10px] border border-gray-200 dark:border-[#374151] bg-white dark:bg-[#1F2937] p-3 text-center animate-fade-up"
                        style={{ animationDelay: `${(i + 1) * 100}ms` }}
                      >
                        <p className="text-xl font-bold text-gray-900 dark:text-white leading-none mb-1">{o.metric}</p>
                        <p className="text-[11px] text-gray-500 dark:text-gray-400 leading-snug">{o.label}</p>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Future Revenue */}
                {current.futureRevenue.length > 0 && (
                  <div className="px-6 pb-6 pt-2 border-t border-gray-100 dark:border-[#1F2937]">
                    <div className="flex items-center gap-2 mb-3">
                      <ShoppingBag className="w-4 h-4 text-violet-500" />
                      <span className="text-xs font-bold text-gray-400 dark:text-gray-500 uppercase tracking-widest">
                        New Revenue Streams
                      </span>
                    </div>
                    <div className="flex flex-wrap gap-2">
                      {current.futureRevenue.map((r) => (
                        <span
                          key={r}
                          className="px-2.5 py-1 rounded-[8px] text-xs font-medium text-gray-600 dark:text-gray-300 bg-gray-100 dark:bg-[#1F2937] border border-gray-200 dark:border-[#374151]"
                        >
                          {r}
                        </span>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        )}

        {/* ── Tab: Business Intelligence ── */}
        {activeTab === "intelligence" && (
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-10 animate-fade-in">

            {/* Left: Dashboard selector */}
            <div className="lg:col-span-4">
              <div className="flex flex-col gap-1.5">
                {dashboards.map((d, idx) => {
                  const DIcon = d.icon;
                  const isActive = idx === activeDashboard;
                  return (
                    <button
                      key={d.label}
                      onClick={() => setActiveDashboard(idx)}
                      className={cn(
                        "w-full text-left px-4 py-3.5 rounded-[12px] border transition-all duration-300 cursor-pointer bg-transparent flex items-center gap-3",
                        isActive
                          ? "border-gray-200 dark:border-[#374151] bg-[#F8FAFC] dark:bg-[#111827] shadow-sm"
                          : "border-transparent hover:border-gray-200 dark:hover:border-[#374151] hover:bg-gray-50 dark:hover:bg-[#111827]/40"
                      )}
                    >
                      <div
                        className="w-9 h-9 rounded-[10px] flex items-center justify-center flex-shrink-0"
                        style={{ backgroundColor: isActive ? `${d.color}15` : `${d.color}08` }}
                      >
                        <DIcon className="w-4.5 h-4.5" style={{ color: d.color }} />
                      </div>
                      <div>
                        <span className={cn(
                          "text-sm font-semibold transition-colors block",
                          isActive ? "text-gray-900 dark:text-white" : "text-gray-500 dark:text-gray-400"
                        )}>
                          {d.label} Dashboard
                        </span>
                        <span className="text-xs text-gray-400 dark:text-gray-500">{d.metrics.length} metrics</span>
                      </div>
                      <ArrowRight className={cn(
                        "w-3.5 h-3.5 ml-auto transition-all",
                        isActive ? "text-[#0F766E] opacity-100" : "opacity-0"
                      )} />
                    </button>
                  );
                })}
              </div>
            </div>

            {/* Right: Dashboard preview */}
            <div className="lg:col-span-8">
              <div
                key={activeDashboard}
                className="rounded-[16px] border border-gray-200 dark:border-[#1F2937] bg-white dark:bg-[#111827] overflow-hidden animate-fade-in"
                style={{ boxShadow: "0 1px 2px rgba(0,0,0,0.04), 0 8px 32px rgba(0,0,0,0.06)" }}
              >
                {/* Dashboard header */}
                <div
                  className="px-6 py-5 flex items-center gap-3"
                  style={{ backgroundColor: `${currentDash.color}06`, borderBottom: `1px solid ${currentDash.color}12` }}
                >
                  <div
                    className="w-11 h-11 rounded-[10px] flex items-center justify-center"
                    style={{ backgroundColor: `${currentDash.color}12` }}
                  >
                    {(() => { const DI = currentDash.icon; return <DI className="w-5 h-5" style={{ color: currentDash.color }} />; })()}
                  </div>
                  <div>
                    <h3 className="text-lg font-bold text-gray-900 dark:text-white">{currentDash.label} Dashboard</h3>
                    <p className="text-xs text-gray-400 dark:text-gray-500">Real-time business intelligence for healthcare organizations</p>
                  </div>
                </div>

                {/* Metric grid */}
                <div className="p-6">
                  <p className="text-[10px] font-bold text-gray-400 dark:text-gray-500 uppercase tracking-widest mb-4">
                    Key Metrics Available
                  </p>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                    {currentDash.metrics.map((m, i) => (
                      <div
                        key={m}
                        className="flex items-center gap-3 p-3 rounded-[10px] border border-gray-100 dark:border-[#1F2937] bg-[#F8FAFC] dark:bg-[#1F2937]/50 hover:border-gray-200 dark:hover:border-[#374151] transition-colors animate-fade-up"
                        style={{ animationDelay: `${i * 60}ms` }}
                      >
                        <div
                          className="w-8 h-8 rounded-[8px] flex items-center justify-center flex-shrink-0"
                          style={{ backgroundColor: `${currentDash.color}10` }}
                        >
                          <BarChart3 className="w-3.5 h-3.5" style={{ color: currentDash.color }} />
                        </div>
                        <span className="text-sm font-medium text-gray-700 dark:text-gray-300">{m}</span>
                      </div>
                    ))}
                  </div>
                </div>

                {/* CTA */}
                <div className="px-6 pb-6 pt-2 border-t border-gray-100 dark:border-[#1F2937] flex items-center justify-between">
                  <p className="text-xs text-gray-400 dark:text-gray-500">
                    Make better business decisions with real-time data.
                  </p>
                  <Link to="/book-demo">
                    <button className="flex items-center gap-1.5 text-sm font-semibold text-[#0F766E] hover:text-[#115E59] transition-colors bg-transparent border-0 cursor-pointer">
                      See it live
                      <ArrowRight className="w-4 h-4" />
                    </button>
                  </Link>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Philosophy callout */}
        <div className="mt-12 p-8 rounded-[16px] border border-[#0F766E]/15 bg-gradient-to-r from-[#F0FDFA] to-[#F0FDFA]/50 dark:from-[#0F766E]/5 dark:to-transparent text-center">
          <p className="text-[17px] leading-relaxed font-medium text-gray-700 dark:text-gray-300 max-w-3xl mx-auto">
            "Curexal is not just software that digitizes healthcare operations, it is a <strong className="text-[#0F766E]">business growth platform</strong> that empowers healthcare providers to operate more efficiently, collaborate more effectively, and generate sustainable revenue."
          </p>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-3">
            Grow your healthcare business. That's the Curexal promise.
          </p>
        </div>

      </div>
    </section>
  );
}
