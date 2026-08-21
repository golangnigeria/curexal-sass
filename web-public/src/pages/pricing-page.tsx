import React, { useState, useEffect } from "react";
import { MarketingNavbar } from "@/components/layouts/marketing-navbar";
import { MarketingFooter } from "@/components/layouts/marketing-footer";
import { SEOHead } from "@/components/seo/seo-head";
import { Skeleton } from "@/components/ui/skeleton";
import { EmptyState } from "@/components/ui/empty-state";
import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";
import {
  Check,
  ArrowRight,
  Sparkles,
  HelpCircle,
  Building2,
  ExternalLink,
} from "lucide-react";
import { getApiUrl } from "@/api";
import { env } from "@/config/env";

interface PlanItem {
  id: string;
  name: string;
  description: string;
  priceMonthly: string;
  priceAnnual: string;
  features: string[];
  popular?: boolean;
}

export function PricingPage() {
  const [annual, setAnnual] = useState(true);
  const [loading, setLoading] = useState(true);
  const [plans, setPlans] = useState<PlanItem[]>([]);
  const orgPortalUrl = env.VITE_PORTAL_URL || "https://app.curexal.com";

  useEffect(() => {
    let isMounted = true;
    setLoading(true);

    fetch(getApiUrl("/plans"))
      .then((res) => {
        if (!res.ok) throw new Error(`Status ${res.status}`);
        return res.json();
      })
      .then((data) => {
        if (isMounted) {
          setPlans(Array.isArray(data?.plans) ? data.plans : []);
          setLoading(false);
        }
      })
      .catch((err) => {
        if (isMounted) {
          console.info("Plans API endpoint currently unconfigured on backend:", err.message);
          setPlans([]);
          setLoading(false);
        }
      });

    return () => {
      isMounted = false;
    };
  }, []);

  return (
    <div className="min-h-screen bg-[#F8FAFC] dark:bg-[#0B1120] font-inter text-slate-900 dark:text-slate-100 flex flex-col">
      <SEOHead
        title="Enterprise Pricing & Plans"
        description="Transparent subscription plans for laboratories, clinics, and healthcare organizations."
      />

      <MarketingNavbar />

      <main className="flex-1 pt-24 pb-16">
        <section className="max-w-[1280px] mx-auto px-6 mb-12 text-center">
          <h1 className="text-3xl md:text-5xl font-extrabold tracking-tight text-slate-900 dark:text-white mb-4">
            Plans Built for Healthcare Scale
          </h1>
          <p className="text-base md:text-lg text-slate-600 dark:text-slate-400 max-w-2xl mx-auto mb-8">
            Choose the workspace plan that fits your clinic, laboratory, or healthcare facility. All plans include PostgreSQL data isolation and RBAC security.
          </p>

          {/* Monthly / Annual Toggle */}
          <div className="inline-flex items-center p-1 rounded-2xl bg-slate-200/70 dark:bg-slate-800 border border-slate-300/60 dark:border-slate-700">
            <button
              onClick={() => setAnnual(false)}
              className={`px-4 py-2 rounded-xl text-xs font-bold transition-all ${
                !annual
                  ? "bg-white dark:bg-slate-900 text-slate-900 dark:text-white shadow-sm"
                  : "text-slate-600 dark:text-slate-400"
              }`}
            >
              Monthly Billing
            </button>
            <button
              onClick={() => setAnnual(true)}
              className={`px-4 py-2 rounded-xl text-xs font-bold transition-all flex items-center gap-1.5 ${
                annual
                  ? "bg-[#0F766E] text-white shadow-sm"
                  : "text-slate-600 dark:text-slate-400"
              }`}
            >
              <span>Annual Billing</span>
              <span className="px-1.5 py-0.5 rounded-full bg-teal-400/30 text-teal-100 text-[10px]">
                Save 20%
              </span>
            </button>
          </div>
        </section>

        {/* Pricing Cards / Empty State */}
        <section className="max-w-[1280px] mx-auto px-6">
          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              {[1, 2, 3].map((i) => (
                <div key={i} className="p-8 rounded-3xl bg-white dark:bg-[#111827] border border-slate-200/80 dark:border-slate-800 space-y-4">
                  <Skeleton className="h-6 w-1/2 rounded-md" />
                  <Skeleton className="h-10 w-3/4 rounded-md" />
                  <Skeleton className="h-32 w-full rounded-2xl" />
                  <Skeleton className="h-12 w-full rounded-xl" />
                </div>
              ))}
            </div>
          ) : plans.length === 0 ? (
            <EmptyState
              badge="Backend API Source"
              title="No Subscription Plans Configured in Backend"
              description="The backend database has zero active plans configured in database tables (or the `/api/v1/plans` endpoint is not provisioned). `web-public` strictly reflects backend API data without hardcoding pricing."
              icon="server"
              actionLabel="Contact Sales for Custom Quote"
              onAction={() => {
                window.location.href = "/contact";
              }}
            />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              {plans.map((plan) => (
                <div
                  key={plan.id}
                  className={`rounded-3xl p-8 bg-white dark:bg-[#111827] border flex flex-col justify-between transition-all ${
                    plan.popular
                      ? "border-[#0F766E] shadow-xl ring-2 ring-[#0F766E]/20"
                      : "border-slate-200/80 dark:border-slate-800 shadow-sm"
                  }`}
                >
                  <div>
                    {plan.popular && (
                      <span className="inline-block px-3 py-1 rounded-full text-[10px] font-extrabold uppercase tracking-wider bg-[#0F766E] text-white mb-4">
                        Most Popular
                      </span>
                    )}
                    <h3 className="text-xl font-bold text-slate-900 dark:text-white mb-1">
                      {plan.name}
                    </h3>
                    <p className="text-xs text-slate-500 dark:text-slate-400 mb-6 leading-relaxed">
                      {plan.description}
                    </p>

                    <div className="mb-6">
                      <span className="text-3xl font-extrabold text-slate-900 dark:text-white">
                        ₦{annual ? plan.priceAnnual : plan.priceMonthly}
                      </span>
                      <span className="text-xs text-slate-500 dark:text-slate-400 ml-1">
                        / month
                      </span>
                    </div>

                    <ul className="space-y-3 mb-8">
                      {plan.features.map((feat, idx) => (
                        <li key={idx} className="flex items-start gap-2 text-xs text-slate-700 dark:text-slate-300">
                          <Check className="w-4 h-4 text-[#0F766E] flex-shrink-0 mt-0.5" />
                          <span>{feat}</span>
                        </li>
                      ))}
                    </ul>
                  </div>

                  <a
                    href={`${orgPortalUrl}/register?plan=${plan.id}`}
                    className={`w-full h-11 rounded-xl text-xs font-bold flex items-center justify-center gap-2 transition-all ${
                      plan.popular
                        ? "bg-[#0F766E] hover:bg-[#115E59] text-white shadow-sm"
                        : "bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700 text-slate-900 dark:text-white"
                    }`}
                  >
                    <span>Subscribe Plan</span>
                    <ExternalLink className="w-3.5 h-3.5" />
                  </a>
                </div>
              ))}
            </div>
          )}
        </section>
      </main>

      <MarketingFooter />
    </div>
  );
}
