import { MarketingNavbar } from "@/components/layouts/marketing-navbar";
import { MarketingFooter } from "@/components/layouts/marketing-footer";
import { SEOHead } from "@/components/seo/seo-head";

// Locked Hero Component (Untouched)
import { Hero } from "@/components/home/hero";

// Telemetry & Modern Sections below Hero
import { WaitlistStats } from "@/components/home/waitlist-stats";
import { HealthcareFragmentationSection } from "@/components/home/healthcare-fragmentation";
import { ConnectedStackedJourneySection } from "@/components/home/connected-stacked-journey";
import { PatientJourneyFlowSection } from "@/components/home/patient-journey-flow";
import { OrganizationNetworkNodeSection } from "@/components/home/organization-network-node";
import { MarketplacePreviewCarouselSection } from "@/components/home/marketplace-preview-carousel";
import { AfricaRealityInfrastructureSection } from "@/components/home/africa-reality-infrastructure";
import { HealthcareOperatingSystemSection } from "@/components/home/healthcare-operating-system";
import { BusinessGrowth } from "@/components/home/business-growth";
import { WaitlistResearchSection } from "@/components/home/waitlist-research-section";

export function LandingPage() {
  return (
    <div className="min-h-screen bg-white dark:bg-[#0B1120] text-slate-900 dark:text-white font-inter overflow-x-hidden">
      <SEOHead
        title="Curexal: Connected Healthcare Operating Network for Africa"
        description="Connecting patients, clinics, diagnostic laboratories, pharmacies and partners so referrals, diagnostics and healthcare transactions move together without friction."
      />
      
      <MarketingNavbar />
      
      {/* ── LOCKED HERO SECTION (UNTOUCHED) ── */}
      <Hero />
      
      {/* ── SECTIONS BELOW HERO ── */}
      
      {/* Live Telemetry Banner */}
      <div className="max-w-[1280px] mx-auto px-4 sm:px-6">
        <WaitlistStats />
      </div>

      {/* Section 1: Healthcare is Fragmented */}
      <HealthcareFragmentationSection />

      {/* Section 2: Curexal Connects the Journey (Stacked Cards) */}
      <ConnectedStackedJourneySection />

      {/* Section 3: One Patient Journey */}
      <PatientJourneyFlowSection />

      {/* Section 4: For Healthcare Organizations (Central Facility Node) */}
      <OrganizationNetworkNodeSection />

      {/* Section 5: Curexal Marketplace (Horizontal Asymmetric Preview) */}
      <MarketplacePreviewCarouselSection />

      {/* Section 6: Built for Africa */}
      <AfricaRealityInfrastructureSection />

      {/* Section 7: Healthcare Operating System */}
      <HealthcareOperatingSystemSection />

      {/* Business Growth Platform */}
      <BusinessGrowth />

      {/* Section 8: Help Us Build It */}
      <WaitlistResearchSection />

    

      <MarketingFooter />
    </div>
  );
}
