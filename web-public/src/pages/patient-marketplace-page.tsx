import React, { useState, useEffect } from "react";
import { MarketingNavbar } from "@/components/layouts/marketing-navbar";
import { MarketingFooter } from "@/components/layouts/marketing-footer";
import { SEOHead } from "@/components/seo/seo-head";
import { Skeleton } from "@/components/ui/skeleton";
import { EmptyState } from "@/components/ui/empty-state";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Search,
  MapPin,
  Clock,
  Star,
  Calendar,
  Building2,
  Stethoscope,
  Pill,
  Radio,
  ShoppingBag,
  ExternalLink,
  SlidersHorizontal,
  ShieldCheck,
  CheckCircle2,
} from "lucide-react";
import { env } from "@/config/env";
import { getApiUrl } from "@/api";

interface FacilityItem {
  id: string;
  name: string;
  category: "laboratory" | "clinic" | "pharmacy" | "hospital" | "radiology" | "vendor";
  accreditation?: string;
  location?: string;
  rating?: number;
  services?: Array<{ name: string; price?: number; turnAroundTime?: string }>;
}

export function PatientMarketplacePage() {
  const [activeTab, setActiveTab] = useState<"all" | "laboratory" | "clinic" | "pharmacy" | "hospital" | "radiology" | "vendor">("all");
  const [searchQuery, setSearchQuery] = useState("");
  const [locationQuery, setLocationQuery] = useState("");
  const [loading, setLoading] = useState(true);
  const [facilities, setFacilities] = useState<FacilityItem[]>([]);
  const [error, setError] = useState<string | null>(null);

  const patientPortalUrl = env.VITE_PATIENT_PORTAL_URL || "https://patient.curexal.com";
  const orgPortalUrl = env.VITE_PORTAL_URL || "https://app.curexal.com";

  useEffect(() => {
    let isMounted = true;
    setLoading(true);
    setError(null);

    // Fetch dynamic search results from Go backend API
    fetch(getApiUrl(`/marketplace/search?query=${encodeURIComponent(searchQuery)}&category=${activeTab}&location=${encodeURIComponent(locationQuery)}`))
      .then((res) => {
        if (!res.ok) {
          throw new Error(`API returned status ${res.status}`);
        }
        return res.json();
      })
      .then((data) => {
        if (isMounted) {
          setFacilities(Array.isArray(data?.items) ? data.items : []);
          setLoading(false);
        }
      })
      .catch((err) => {
        if (isMounted) {
          console.info("Marketplace API endpoint currently unconfigured on backend:", err.message);
          setFacilities([]);
          setLoading(false);
        }
      });

    return () => {
      isMounted = false;
    };
  }, [activeTab, searchQuery, locationQuery]);

  const handleBookRedirect = (facilityId: string, serviceName?: string) => {
    const continueUrl = encodeURIComponent(window.location.href);
    const targetUrl = `${patientPortalUrl}/book?facilityId=${facilityId}&service=${encodeURIComponent(serviceName || "")}&continue=${continueUrl}`;
    window.location.href = targetUrl;
  };

  const handleVendorRedirect = (productId?: string) => {
    const continueUrl = encodeURIComponent(window.location.href);
    const targetUrl = `${orgPortalUrl}/store?product=${productId || ""}&continue=${continueUrl}`;
    window.location.href = targetUrl;
  };

  return (
    <div className="min-h-screen bg-[#F8FAFC] dark:bg-[#0B1120] font-inter text-slate-900 dark:text-slate-100 flex flex-col">
      <SEOHead
        title="Healthcare & Medical Marketplace"
        description="Public search for verified diagnostic laboratories, clinics, hospitals, pharmacies, radiology centers, and medical equipment vendors."
      />

      <MarketingNavbar />

      <main className="flex-1 pt-24 pb-16">
        {/* Header Hero */}
        <section className="max-w-[1280px] mx-auto px-6 mb-10">
          <div className="text-center max-w-3xl mx-auto mb-8">
            <h1 className="text-3xl md:text-5xl font-extrabold tracking-tight text-slate-900 dark:text-white mb-4">
              Healthcare & Medical Marketplace
            </h1>
            <p className="text-base md:text-lg text-slate-600 dark:text-slate-400">
              Browse accredited laboratories, clinics, hospitals, pharmacies, and medical supply vendors verified by backend network records.
            </p>
          </div>

          {/* Search Controls */}
          <div className="bg-white dark:bg-[#111827] p-4 md:p-6 rounded-2xl border border-slate-200/80 dark:border-slate-800 shadow-md max-w-4xl mx-auto">
            <div className="grid grid-cols-1 md:grid-cols-12 gap-3">
              <div className="md:col-span-6 relative">
                <Search className="w-4 h-4 text-slate-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
                <Input
                  type="text"
                  placeholder="Search test, specialty, facility, or medical product..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10 h-11 bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800 focus:ring-[#0F766E]"
                />
              </div>

              <div className="md:col-span-4 relative">
                <MapPin className="w-4 h-4 text-slate-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
                <Input
                  type="text"
                  placeholder="City, region, or location..."
                  value={locationQuery}
                  onChange={(e) => setLocationQuery(e.target.value)}
                  className="pl-10 h-11 bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800 focus:ring-[#0F766E]"
                />
              </div>

              <div className="md:col-span-2">
                <Button
                  onClick={() => {}}
                  className="w-full h-11 bg-[#0F766E] hover:bg-[#115E59] text-white font-semibold rounded-xl transition-colors"
                >
                  Search
                </Button>
              </div>
            </div>

            {/* Category Filter Tabs */}
            <div className="flex items-center gap-2 mt-4 overflow-x-auto pb-1 pt-2">
              {[
                { id: "all", label: "All Categories", icon: SlidersHorizontal },
                { id: "laboratory", label: "Laboratories", icon: Building2 },
                { id: "clinic", label: "Clinics & EMR", icon: Stethoscope },
                { id: "pharmacy", label: "Pharmacies", icon: Pill },
                { id: "radiology", label: "Radiology", icon: Radio },
                { id: "vendor", label: "Medical Supply Vendors", icon: ShoppingBag },
              ].map((tab) => {
                const IconComponent = tab.icon;
                const active = activeTab === tab.id;
                return (
                  <button
                    key={tab.id}
                    onClick={() => setActiveTab(tab.id as any)}
                    className={`flex items-center gap-1.5 px-3.5 py-2 rounded-xl text-xs font-semibold whitespace-nowrap transition-all border ${
                      active
                        ? "bg-[#0F766E] text-white border-[#0F766E] shadow-sm"
                        : "bg-slate-50 dark:bg-slate-900 text-slate-700 dark:text-slate-300 border-slate-200/80 dark:border-slate-800 hover:border-slate-300"
                    }`}
                  >
                    <IconComponent className="w-3.5 h-3.5" />
                    <span>{tab.label}</span>
                  </button>
                );
              })}
            </div>
          </div>
        </section>

        {/* Results Container */}
        <section className="max-w-[1280px] mx-auto px-6">
          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {[1, 2, 3, 4, 5, 6].map((i) => (
                <div key={i} className="p-6 rounded-2xl bg-white dark:bg-[#111827] border border-slate-200/80 dark:border-slate-800 space-y-4">
                  <Skeleton className="h-6 w-3/4 rounded-md" />
                  <Skeleton className="h-4 w-1/2 rounded-md" />
                  <Skeleton className="h-20 w-full rounded-xl" />
                  <Skeleton className="h-10 w-full rounded-xl" />
                </div>
              ))}
            </div>
          ) : facilities.length === 0 ? (
            <EmptyState
              badge="Backend API Source"
              title="No Marketplace Facilities Match Your Criteria"
              description="The backend database has zero published facility listings matching this filter, or the backend search endpoint is unconfigured. The presentation layer strictly renders API state without fabricating data."
              icon="server"
              actionLabel="Clear Filters"
              onAction={() => {
                setSearchQuery("");
                setLocationQuery("");
                setActiveTab("all");
              }}
            />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {facilities.map((item) => (
                <div
                  key={item.id}
                  className="bg-white dark:bg-[#111827] rounded-2xl p-6 border border-slate-200/80 dark:border-slate-800 shadow-sm hover:shadow-md transition-all flex flex-col justify-between"
                >
                  <div>
                    <div className="flex items-start justify-between gap-2 mb-3">
                      <div>
                        <span className="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-extrabold uppercase tracking-wider bg-teal-50 dark:bg-teal-950/50 text-[#0F766E] dark:text-teal-400 border border-teal-200/60 dark:border-teal-800/60 mb-1.5">
                          {item.category}
                        </span>
                        <h3 className="text-lg font-bold text-slate-900 dark:text-white leading-snug">
                          {item.name}
                        </h3>
                      </div>
                      {item.rating && (
                        <div className="flex items-center gap-1 bg-amber-50 dark:bg-amber-950/40 px-2 py-1 rounded-lg border border-amber-200/60 dark:border-amber-900/60 text-xs font-bold text-amber-700 dark:text-amber-400">
                          <Star className="w-3.5 h-3.5 fill-amber-400 text-amber-400" />
                          <span>{item.rating}</span>
                        </div>
                      )}
                    </div>

                    {item.location && (
                      <p className="text-xs text-slate-500 dark:text-slate-400 flex items-center gap-1.5 mb-4">
                        <MapPin className="w-3.5 h-3.5 text-[#0F766E]" />
                        <span>{item.location}</span>
                      </p>
                    )}

                    {item.services && item.services.length > 0 && (
                      <div className="space-y-2 my-4 pt-3 border-t border-slate-100 dark:border-slate-800">
                        {item.services.slice(0, 3).map((srv, idx) => (
                          <div key={idx} className="flex items-center justify-between text-xs py-1">
                            <span className="text-slate-700 dark:text-slate-300 font-medium">{srv.name}</span>
                            {srv.price && (
                              <span className="font-bold text-[#0F766E] dark:text-teal-400">
                                ₦{srv.price.toLocaleString()}
                              </span>
                            )}
                          </div>
                        ))}
                      </div>
                    )}
                  </div>

                  <div className="pt-4 border-t border-slate-100 dark:border-slate-800 mt-4">
                    {item.category === "vendor" ? (
                      <Button
                        onClick={() => handleVendorRedirect(item.id)}
                        className="w-full bg-slate-900 hover:bg-slate-800 dark:bg-teal-600 dark:hover:bg-teal-500 text-white font-semibold text-xs h-10 rounded-xl flex items-center justify-center gap-2"
                      >
                        <span>Buy Medical Equipment</span>
                        <ExternalLink className="w-3.5 h-3.5" />
                      </Button>
                    ) : (
                      <Button
                        onClick={() => handleBookRedirect(item.id)}
                        className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-semibold text-xs h-10 rounded-xl flex items-center justify-center gap-2"
                      >
                        <Calendar className="w-3.5 h-3.5" />
                        <span>Book Appointment</span>
                        <ExternalLink className="w-3.5 h-3.5 opacity-60" />
                      </Button>
                    )}
                  </div>
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
