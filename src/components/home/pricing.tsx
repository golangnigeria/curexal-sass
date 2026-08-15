import { useState } from "react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Check, HelpCircle, Activity, Sparkles, Clock } from "lucide-react";

export function Pricing() {
  const [isAnnual, setIsAnnual] = useState(false);

  const plans = [
    {
      name: "Smart",
      desc: "For small community clinics and private practices.",
      priceMonthly: "87,500",
      priceAnnual: "70,000",
      features: [
        "Up to 5 active staff members",
        "1,000 specimens processed / mo",
        "Core clinical analyte catalog",
        "Standard email support (24h)",
        "Basic audit logging (30 days)",
        "Single laboratory workspace",
      ],
      buttonText: "Book Demo (Soon)",
      buttonLink: "/book-demo",
      popular: false,
    },
    {
      name: "Optimized",
      desc: "For expanding clinics and reference laboratories.",
      priceMonthly: "145,000",
      priceAnnual: "116,000",
      features: [
        "Up to 15 active staff members",
        "5,000 specimens processed / mo",
        "Custom reference ranges",
        "Priority email support (12h)",
        "Audit logs history (1 year)",
        "Multi-device logging",
      ],
      buttonText: "Book Demo (Soon)",
      buttonLink: "/book-demo",
      popular: false,
    },
    {
      name: "Pro",
      desc: "For regional diagnostic centers and busy laboratory groups.",
      priceMonthly: "285,000",
      priceAnnual: "228,000",
      features: [
        "Up to 50 active staff members",
        "25,000 specimens processed / mo",
        "Pathologist validation gateways",
        "Priority support (under 4h)",
        "Unlimited audit trail history",
        "Real-time CSV & PDF report exports",
        "Multi-device analyzer integration APIs",
      ],
      buttonText: "Book Demo (Soon)",
      buttonLink: "/book-demo",
      popular: true,
      popularLabel: "Most Adopted",
    },
    {
      name: "Power",
      desc: "For hospital systems and high-volume reference laboratories.",
      priceMonthly: "Custom",
      priceAnnual: "Custom",
      features: [
        "Unlimited staff & pathologists",
        "Unlimited specimens / mo",
        "Dedicated onboarding manager",
        "24/7 phone & screen-share support",
        "Audit compliance readiness",
        "Custom analyzer interface development",
        "99.99% Uptime SLA target",
        "On-premise deployment options",
      ],
      buttonText: "Book Demo (Soon)",
      buttonLink: "/book-demo",
      popular: false,
    },
  ];

  return (
    <section id="pricing" className="py-24 px-6 border-t border-border bg-muted/5 relative">
      <div className="max-w-[1280px] mx-auto">
        <div className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="text-3xl sm:text-4xl font-extrabold tracking-tight text-slate-900 dark:text-white mb-4">
            Predictable Healthcare Operating Tier Plans
          </h2>
          <p className="text-slate-600 dark:text-slate-400 text-base max-w-2xl mx-auto">
            Choose the subscription model that matches your facility diagnostic processing volume. All tiers include full HIPAA & NDPR compliance isolation.
          </p>

          {/* Billing Switch */}
          <div className="flex items-center justify-center gap-3 mt-8">
            <span className={`text-sm font-semibold ${!isAnnual ? "text-slate-900 dark:text-white" : "text-slate-500"}`}>
              Monthly Billing
            </span>
            <button
              onClick={() => setIsAnnual(!isAnnual)}
              className="w-12 h-6 rounded-full bg-slate-200 dark:bg-slate-800 p-1 relative transition-colors cursor-pointer"
            >
              <div className={`w-4 h-4 rounded-full bg-[#0F766E] transition-transform ${isAnnual ? "translate-x-6" : ""}`} />
            </button>
            <span className={`text-sm font-semibold flex items-center gap-1.5 ${isAnnual ? "text-slate-900 dark:text-white" : "text-slate-500"}`}>
              <span>Annual Billing</span>
              <span className="text-[10px] font-bold px-2 py-0.5 rounded-full bg-teal-100 dark:bg-teal-950 text-[#0F766E] dark:text-teal-400">
                Save 20%
              </span>
            </span>
          </div>
        </div>

        {/* Pricing Cards Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {plans.map((plan, i) => (
            <div
              key={i}
              className={`rounded-2xl p-6 flex flex-col justify-between transition-all ${
                plan.popular
                  ? "bg-white dark:bg-slate-900 border-2 border-[#0F766E] shadow-xl relative"
                  : "bg-white dark:bg-slate-900/60 border border-slate-200 dark:border-slate-800 shadow-sm"
              }`}
            >
              <div>
                {plan.popular && (
                  <span className="absolute -top-3.5 left-1/2 -translate-x-1/2 px-3 py-1 rounded-full bg-[#0F766E] text-white text-[10px] font-extrabold uppercase tracking-wider shadow-md">
                    {plan.popularLabel}
                  </span>
                )}

                <h3 className="text-xl font-bold text-slate-900 dark:text-white">{plan.name}</h3>
                <p className="text-xs text-slate-500 dark:text-slate-400 mt-1 min-h-[36px]">{plan.desc}</p>

                <div className="my-6">
                  <span className="text-3xl font-black text-slate-900 dark:text-white">
                    {plan.priceMonthly === "Custom" ? "Custom" : `₦${isAnnual ? plan.priceAnnual : plan.priceMonthly}`}
                  </span>
                  {plan.priceMonthly !== "Custom" && (
                    <span className="text-xs text-slate-500 dark:text-slate-400 font-medium"> / month</span>
                  )}
                </div>

                <div className="space-y-2.5 pt-4 border-t border-slate-100 dark:border-slate-800">
                  {plan.features.map((feat, fi) => (
                    <div key={fi} className="flex items-start gap-2 text-xs text-slate-600 dark:text-slate-300">
                      <Check className="w-3.5 h-3.5 text-[#0F766E] flex-shrink-0 mt-0.5" />
                      <span>{feat}</span>
                    </div>
                  ))}
                </div>
              </div>

              <div className="pt-6">
                <Link to={plan.buttonLink} className="w-full block">
                  <Button
                    className={`w-full font-bold h-11 rounded-xl flex items-center justify-center gap-1.5 ${
                      plan.popular
                        ? "bg-[#0F766E] hover:bg-[#115E59] text-white"
                        : "bg-slate-100 dark:bg-slate-800 hover:bg-slate-200 text-slate-900 dark:text-white border border-slate-200 dark:border-slate-700"
                    }`}
                  >
                    <span>{plan.buttonText}</span>
                    <Clock className="w-3.5 h-3.5 opacity-70 text-amber-500" />
                  </Button>
                </Link>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
