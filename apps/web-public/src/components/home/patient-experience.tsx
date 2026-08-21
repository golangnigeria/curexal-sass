import React from "react";
import { UserCheck, Search, Calendar, HeartPulse, Shield, FileCheck, CheckCircle2, ArrowRight } from "lucide-react";
import { Link } from "react-router-dom";

export function PatientExperienceSection() {
  const patientSteps = [
    { title: "Find", desc: "Discover accredited clinics, labs, and pharmacies." },
    { title: "Book", desc: "Schedule clinic appointments & diagnostic tests." },
    { title: "Receive Care", desc: "Consult with verified healthcare providers." },
    { title: "Diagnostics", desc: "Sample collected without repeated trips or lost paper." },
    { title: "Results", desc: "Verified PDF reports delivered straight to your portal." },
    { title: "Continue Care", desc: "Direct follow-ups with referring doctor." },
  ];

  return (
    <section className="py-20 bg-slate-50 dark:bg-slate-900/40 border-b border-slate-200/60 dark:border-slate-800">
      <div className="max-w-[1280px] mx-auto px-4 sm:px-6">
        
        {/* Header */}
        <div className="text-center max-w-3xl mx-auto mb-14">
          <h2 className="text-3xl sm:text-4xl font-extrabold text-slate-900 dark:text-white tracking-tight mb-4">
            For patients: Stop being the paper courier between doctors and labs.
          </h2>
          <p className="text-sm sm:text-base text-slate-600 dark:text-slate-300 leading-relaxed">
            You shouldn't have to carry paper requisitions across town, make multiple trips to check if test results are ready, or re-explain your medical history at every facility.
          </p>
        </div>

        {/* 6-Step Journey Ribbon */}
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3 mb-12">
          {patientSteps.map((step, idx) => (
            <div
              key={step.title}
              className="p-4 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 text-center space-y-1.5 shadow-sm"
            >
              <span className="w-6 h-6 rounded-full bg-teal-50 dark:bg-teal-950 text-[#0F766E] dark:text-teal-400 font-bold text-xs inline-flex items-center justify-center">
                {idx + 1}
              </span>
              <h3 className="text-xs font-bold text-slate-900 dark:text-white">{step.title}</h3>
              <p className="text-[11px] text-slate-500 dark:text-slate-400 leading-tight">{step.desc}</p>
            </div>
          ))}
        </div>

        {/* Value Points Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          
          <div className="p-6 rounded-3xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm space-y-3">
            <div className="w-10 h-10 rounded-xl bg-teal-50 dark:bg-teal-950 text-[#0F766E] flex items-center justify-center">
              <Search className="w-5 h-5" />
            </div>
            <h3 className="text-base font-bold text-slate-900 dark:text-white">Discover & Book Online</h3>
            <p className="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
              Find accredited pathology laboratories, imaging centers, and clinics near you. Book diagnostic appointments or home sample collection directly.
            </p>
          </div>

          <div className="p-6 rounded-3xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm space-y-3">
            <div className="w-10 h-10 rounded-xl bg-teal-50 dark:bg-teal-950 text-[#0F766E] flex items-center justify-center">
              <FileCheck className="w-5 h-5" />
            </div>
            <h3 className="text-base font-bold text-slate-900 dark:text-white">Digital Test Reports</h3>
            <p className="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
              Receive electronically signed PDF diagnostic results directly on your phone as soon as the pathologist authorizes them.
            </p>
          </div>

          <div className="p-6 rounded-3xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm space-y-3">
            <div className="w-10 h-10 rounded-xl bg-teal-50 dark:bg-teal-950 text-[#0F766E] flex items-center justify-center">
              <HeartPulse className="w-5 h-5" />
            </div>
            <h3 className="text-base font-bold text-slate-900 dark:text-white">Reduced Unnecessary Trips</h3>
            <p className="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
              Your referring doctor automatically receives your results as soon as they are ready, enabling faster care decisions without extra laboratory visits.
            </p>
          </div>

        </div>

        <div className="mt-10 text-center">
          <Link to="/waitlist">
            <button className="px-6 py-3 rounded-xl bg-[#0F766E] hover:bg-[#115E59] text-white text-xs sm:text-sm font-bold transition-all shadow-md cursor-pointer border-0 inline-flex items-center gap-2">
              <span>Join Patient Early Access</span>
              <ArrowRight className="w-4 h-4" />
            </button>
          </Link>
        </div>

      </div>
    </section>
  );
}
