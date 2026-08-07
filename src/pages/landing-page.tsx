import { MarketingNavbar } from "@/components/layouts/marketing-navbar";
import { MarketingFooter } from "@/components/layouts/marketing-footer";
import { SEOHead } from "@/components/seo/seo-head";
import { Hero } from "@/components/home/hero";
import { TrustedMarquee } from "@/components/home/trusted-marquee";
import { EcosystemSection } from "@/components/home/ecosystem";
import { PlatformModules } from "@/components/home/platform-modules";
import { BusinessGrowth } from "@/components/home/business-growth";
import { HowItWorks } from "@/components/home/how-it-works";
import { SecurityCompliance } from "@/components/home/security-compliance";
import { Testimonials } from "@/components/home/testimonials";
import { CtaDark } from "@/components/home/cta-dark";

export function LandingPage() {
  return (
    <div className="min-h-screen bg-white dark:bg-[#0B1120] text-slate-900 dark:text-white font-inter">
      <SEOHead
        title="Curexal — Enterprise Healthcare Platform & Marketplace"
        description="Public platform for diagnostic laboratories, clinics, hospitals, pharmacies, and medical supply vendors."
      />
      <MarketingNavbar />
      <Hero />
      <TrustedMarquee />
      <BusinessGrowth />
      <EcosystemSection />
      <PlatformModules />
      <HowItWorks />
      <SecurityCompliance />
      <Testimonials />
      <CtaDark />
      <MarketingFooter />
    </div>
  );
}
