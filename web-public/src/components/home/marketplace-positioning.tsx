import React from "react";
import { ShoppingBag, ArrowRight, Truck, ShieldCheck, RefreshCw, Layers } from "lucide-react";
import { Link } from "react-router-dom";

export function MarketplacePositioningSection() {
  return (
    <section className="py-20 bg-slate-50 dark:bg-slate-900/50 border-b border-slate-200/60 dark:border-slate-800">
      <div className="max-w-[1280px] mx-auto px-4 sm:px-6">
        
        {/* Header */}
        <div className="text-center max-w-3xl mx-auto mb-14">
          <h2 className="text-3xl sm:text-4xl font-extrabold text-slate-900 dark:text-white tracking-tight mb-4">
            A healthcare network, not just a marketplace.
          </h2>
          <p className="text-sm sm:text-base text-slate-600 dark:text-slate-300 leading-relaxed">
            The Curexal Marketplace is not an isolated e-commerce store. It is an integrated coordination layer connecting healthcare needs to accredited providers, diagnostic labs, pharmacies, and medical suppliers.
          </p>
        </div>

        {/* Network Layer Diagram */}
        <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 sm:p-10 shadow-sm mb-10">
          <div className="text-center max-w-xl mx-auto mb-8">
            <h3 className="text-sm font-bold text-slate-900 dark:text-white uppercase tracking-wider">
              Coordinated Service & Supply Movement
            </h3>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-5 gap-3 items-center text-center">
            
            <div className="p-4 rounded-2xl bg-slate-50 dark:bg-slate-800/60 border border-slate-200 dark:border-slate-700">
              <span className="text-xs font-bold text-slate-900 dark:text-white">1. Patient / Facility</span>
              <p className="text-[11px] text-slate-500 dark:text-slate-400 mt-1">Identifies healthcare or diagnostic need</p>
            </div>

            <ArrowRight className="hidden md:block mx-auto text-teal-600 w-5 h-5" />

            <div className="p-4 rounded-2xl bg-teal-50 dark:bg-teal-950/60 border border-teal-200 dark:border-teal-800">
              <span className="text-xs font-bold text-[#0F766E] dark:text-teal-300">2. Curexal Network</span>
              <p className="text-[11px] text-slate-600 dark:text-slate-400 mt-1">Matches requisition or B2B supply order</p>
            </div>

            <ArrowRight className="hidden md:block mx-auto text-teal-600 w-5 h-5" />

            <div className="p-4 rounded-2xl bg-slate-50 dark:bg-slate-800/60 border border-slate-200 dark:border-slate-700">
              <span className="text-xs font-bold text-slate-900 dark:text-white">3. Verified Fulfillment</span>
              <p className="text-[11px] text-slate-500 dark:text-slate-400 mt-1">Service delivered or product dispatched</p>
            </div>

          </div>
        </div>

        {/* Feature Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="p-6 rounded-3xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 space-y-3">
            <div className="w-10 h-10 rounded-xl bg-teal-50 dark:bg-teal-950 text-[#0F766E] flex items-center justify-center">
              <Layers className="w-5 h-5" />
            </div>
            <h3 className="text-base font-bold text-slate-900 dark:text-white">Diagnostic & Care Services</h3>
            <p className="text-xs text-slate-600 dark:text-slate-400 leading-relaxed">
              Patients and clinics can search accredited laboratories, compare diagnostic test pricing, filter by TAT turnaround time, and book verified services connected directly to LIMS systems.
            </p>
          </div>

          <div className="p-6 rounded-3xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 space-y-3">
            <div className="w-10 h-10 rounded-xl bg-teal-50 dark:bg-teal-950 text-[#0F766E] flex items-center justify-center">
              <Truck className="w-5 h-5" />
            </div>
            <h3 className="text-base font-bold text-slate-900 dark:text-white">B2B Medical Supply Procurement</h3>
            <p className="text-xs text-slate-600 dark:text-slate-400 leading-relaxed">
              Diagnostic laboratories, clinics, and pharmacies can order reagents, test kits, personal protective equipment, and pharmaceutical supplies directly from accredited medical suppliers.
            </p>
          </div>
        </div>

        <div className="mt-10 text-center">
          <Link to="/marketplace">
            <button className="px-6 py-3 rounded-xl bg-[#0F766E] hover:bg-[#115E59] text-white text-xs sm:text-sm font-bold transition-all shadow-md cursor-pointer border-0 inline-flex items-center gap-2">
              <span>Explore Marketplace Network</span>
              <ArrowRight className="w-4 h-4" />
            </button>
          </Link>
        </div>

      </div>
    </section>
  );
}
