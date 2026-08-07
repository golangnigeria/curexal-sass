import { MarketingNavbar } from "@/components/layouts/marketing-navbar";
import { MarketingFooter } from "@/components/layouts/marketing-footer";
import { SEOHead } from "@/components/seo/seo-head";
import { Link } from "react-router-dom";
import { BookOpen, ArrowRight, FileText, ShieldCheck, Rss } from "lucide-react";

const sections = [
  {
    icon: BookOpen,
    label: "Documentation",
    desc: "Complete reference for the Curexal LIS platform, administration guides, and integration documentation.",
    href: "/developers",
    cta: "Explore APIs",
  },
  {
    icon: FileText,
    id: "releases",
    label: "Release Notes",
    desc: "Changelog for every platform update, including new features, fixes, and deprecations.",
    href: "#",
    cta: "View Releases",
  },
  {
    icon: ShieldCheck,
    id: "compliance",
    label: "Compliance Center",
    desc: "ISO 15189 readiness guides, HIPAA alignment documentation, and data processing agreements.",
    href: "#",
    cta: "View Compliance Docs",
  },
  {
    icon: Rss,
    id: "blog",
    label: "Blog",
    desc: "Insights on laboratory management, African healthcare infrastructure, and diagnostic technology.",
    href: "#",
    cta: "Read Articles",
  },
];

export function ResourcesPage() {
  return (
    <div className="min-h-screen bg-white dark:bg-[#0B1120] font-sans">
      <MarketingNavbar />

      <div className="pt-28 pb-16 bg-[#F8FAFC] dark:bg-[#0B1120] border-b border-gray-100 dark:border-[#1F2937]">
        <div className="max-w-[1280px] mx-auto px-6">
          <div className="max-w-xl">
            <div className="accent-line mb-4" />
            <h1 className="text-hero text-gray-900 dark:text-white mb-5">Resources</h1>
            <p className="text-body text-gray-500 dark:text-gray-400">
              Documentation, compliance guides, release notes, and insights for laboratory teams and platform administrators.
            </p>
          </div>
        </div>
      </div>

      <main className="max-w-[1280px] mx-auto px-6 py-16">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-5">
          {sections.map((s) => {
            const Icon = s.icon;
            return (
              <Link key={s.label} to={s.href} className="block group">
                <div className="card-enterprise p-7 hover-lift h-full">
                  <div className="w-10 h-10 rounded-[10px] bg-[#F0FDFA] dark:bg-[#0F766E]/10 flex items-center justify-center mb-5 group-hover:bg-[#0F766E]/15 transition-colors">
                    <Icon className="h-5 w-5 text-[#0F766E]" strokeWidth={1.75} />
                  </div>
                  <h3 className="text-[16px] font-semibold text-gray-900 dark:text-white mb-2">{s.label}</h3>
                  <p className="text-sm text-gray-500 dark:text-gray-400 leading-relaxed mb-5">{s.desc}</p>
                  <span className="flex items-center gap-1.5 text-sm font-semibold text-[#0F766E] group-hover:gap-2.5 transition-all">
                    {s.cta}
                    <ArrowRight className="h-4 w-4" />
                  </span>
                </div>
              </Link>
            );
          })}
        </div>
      </main>

      <MarketingFooter />
    </div>
  );
}
