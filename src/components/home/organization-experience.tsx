import React, { useState } from "react";
import { Stethoscope, FlaskConical, Pill, ArrowRight, Building2 } from "lucide-react";

export function OrganizationExperienceSection() {
  const [activeTab, setActiveTab] = useState<"clinic" | "lab" | "pharmacy">("clinic");

  const orgWorkflows = {
    clinic: {
      title: "Clinic Workflow",
      subtitle: "Connect your outpatient consultation rooms directly with diagnostic laboratories.",
      steps: [
        { label: "1. Patient Visit", desc: "Clinician conducts patient examination" },
        { label: "2. Digital Referral", desc: "Lab requisition created electronically in chart" },
        { label: "3. Lab Dispatch", desc: "Requisition transmitted instantly to partner lab" },
        { label: "4. Result Sync", desc: "Pathologist authorized PDF delivered back to chart" },
        { label: "5. Care Decision", desc: "Physician initiates targeted treatment without delay" },
      ],
    },
    lab: {
      title: "Laboratory Workflow",
      subtitle: "Receive electronic requisitions, manage chain of custody, and auto-dispatch results.",
      steps: [
        { label: "1. Electronic Requisition", desc: "Lab queue populates before patient sample arrives" },
        { label: "2. Sample Accessioning", desc: "Barcoded accession number assigned to specimen" },
        { label: "3. Instrument Testing", desc: "Analyzer auto-transmits raw test results to LIMS" },
        { label: "4. Pathologist Sign-Off", desc: "Multi-tier verification & digital stamp authorization" },
        { label: "5. Automated Dispatch", desc: "Result delivered to referring doctor & patient" },
      ],
    },
    pharmacy: {
      title: "Pharmacy Workflow",
      subtitle: "Fulfill electronic prescriptions and coordinate medical supply reorders efficiently.",
      steps: [
        { label: "1. Prescription / Product", desc: "Electronic prescription received from provider or patient" },
        { label: "2. Inventory Verification", desc: "Live stock level & expiry check" },
        { label: "3. Fulfillment & Dispensing", desc: "Safety check & barcode verification" },
        { label: "4. Delivery / Pickup", desc: "Dispensed to patient with digital record" },
        { label: "5. Supplier Reorder", desc: "Automated replenishment order to B2B suppliers" },
      ],
    },
  };

  const activeData = orgWorkflows[activeTab];

  return (
    <section className="py-20 bg-white dark:bg-[#0B1120] border-b border-slate-100 dark:border-slate-800">
      <div className="max-w-[1280px] mx-auto px-4 sm:px-6">
        
        {/* Header */}
        <div className="text-center max-w-3xl mx-auto mb-12">
          <div className="inline-flex items-center gap-1.5 px-3 py-1 mb-3 rounded-full border border-teal-500/30 bg-teal-50 dark:bg-teal-950/40 text-[#0F766E] dark:text-teal-400 text-xs font-bold uppercase tracking-wider">
            For Healthcare Organizations
          </div>
          <h2 className="text-3xl sm:text-4xl font-extrabold text-slate-900 dark:text-white tracking-tight mb-4">
            Your organization doesn't have to operate alone.
          </h2>
          <p className="text-sm sm:text-base text-slate-600 dark:text-slate-400">
            Curexal connects your facility's internal software capabilities to a broader healthcare network without compromising your operational autonomy.
          </p>
        </div>

        {/* Tab Selection */}
        <div className="flex items-center justify-center gap-2 mb-10">
          <button
            onClick={() => setActiveTab("clinic")}
            className={`px-5 py-2.5 rounded-xl text-xs font-bold transition-all cursor-pointer flex items-center gap-2 ${
              activeTab === "clinic"
                ? "bg-[#0F766E] text-white shadow-md"
                : "bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 hover:bg-slate-200"
            }`}
          >
            <Stethoscope className="w-4 h-4" />
            <span>Clinics & Hospitals</span>
          </button>

          <button
            onClick={() => setActiveTab("lab")}
            className={`px-5 py-2.5 rounded-xl text-xs font-bold transition-all cursor-pointer flex items-center gap-2 ${
              activeTab === "lab"
                ? "bg-[#0F766E] text-white shadow-md"
                : "bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 hover:bg-slate-200"
            }`}
          >
            <FlaskConical className="w-4 h-4" />
            <span>Diagnostic Laboratories</span>
          </button>

          <button
            onClick={() => setActiveTab("pharmacy")}
            className={`px-5 py-2.5 rounded-xl text-xs font-bold transition-all cursor-pointer flex items-center gap-2 ${
              activeTab === "pharmacy"
                ? "bg-[#0F766E] text-white shadow-md"
                : "bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 hover:bg-slate-200"
            }`}
          >
            <Pill className="w-4 h-4" />
            <span>Pharmacies & Suppliers</span>
          </button>
        </div>

        {/* Tab Content Display */}
        <div className="bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 sm:p-10">
          <div className="max-w-2xl mb-8">
            <h3 className="text-xl font-extrabold text-slate-900 dark:text-white mb-2">{activeData.title}</h3>
            <p className="text-xs sm:text-sm text-slate-600 dark:text-slate-400">{activeData.subtitle}</p>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-3">
            {activeData.steps.map((step, idx) => (
              <div
                key={step.label}
                className="p-4 rounded-2xl bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 space-y-2 relative"
              >
                <span className="text-[10px] font-bold text-[#0F766E] uppercase tracking-wider block">
                  Stage 0{idx + 1}
                </span>
                <h4 className="text-xs font-bold text-slate-900 dark:text-white leading-tight">{step.label}</h4>
                <p className="text-[11px] text-slate-500 dark:text-slate-400 leading-snug">{step.desc}</p>
              </div>
            ))}
          </div>
        </div>

      </div>
    </section>
  );
}
