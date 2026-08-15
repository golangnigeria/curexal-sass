import { Link } from "react-router-dom";
import {
  Check,
  ArrowRight,
  Network,
  Clock,
} from "lucide-react";

export function ProvidersCta() {
  return (
    <section className="section-padding bg-[#F8FAFC] dark:bg-[#0B1120] border-y border-gray-100 dark:border-[#1F2937]">
      <div className="max-w-[1280px] mx-auto px-6">
        <div className="grid md:grid-cols-2 gap-12 items-center">
          <div>
            <h2 className="text-section text-gray-900 dark:text-white mb-6">
              Built for Africa's Healthcare Reality.
            </h2>
            <p className="text-body text-gray-500 dark:text-gray-400 mb-8 max-w-lg">
              Most healthcare software is built for countries with reliable electricity and always-on internet. Curexal is intentionally engineered for unstable connectivity, power interruptions, local payment rails, and African healthcare workflows.
            </p>
            <ul className="space-y-3 mb-8">
              {[
                "Schema-per-tenant data isolation, ensuring your data never mixes",
                "Offline-resilient architecture for unreliable connectivity",
                "Local payment integration (Paystack, Flutterwave)",
                "WhatsApp-first patient notifications",
                "Immutable audit trail on every clinical action",
              ].map((item) => (
                <li key={item} className="flex items-center gap-3 text-sm text-gray-700 dark:text-gray-300">
                  <div className="w-5 h-5 rounded-full bg-[#F0FDFA] dark:bg-[#0F766E]/10 flex items-center justify-center flex-shrink-0">
                    <Check className="h-3 w-3 text-[#0F766E]" />
                  </div>
                  {item}
                </li>
              ))}
            </ul>
            <Link to="/book-demo">
              <button
                className="flex items-center gap-2 px-6 py-3 rounded-[12px] bg-[#0F766E] hover:bg-[#115E59] text-white text-sm font-semibold transition-colors cursor-pointer border-0"
                style={{ boxShadow: "0 1px 2px rgba(0,0,0,0.12), 0 4px 16px rgba(15,118,110,0.25)" }}
              >
                <span>Book Demo</span>
                <span className="text-[10px] bg-amber-500/20 text-amber-300 font-extrabold px-1.5 py-0.5 rounded flex items-center gap-1">
                  <Clock className="w-3 h-3" />
                  Soon
                </span>
                <ArrowRight className="h-4 w-4 ml-1" />
              </button>
            </Link>
          </div>

          <div className="space-y-4">
            {[
              {
                title: "LIMS Module",
                desc: "Specimen reception, barcoding, bench worksheets, auto-flagging, pathology sign-off.",
              },
              {
                title: "EMR Module",
                desc: "Patient registration, consultation notes, digital lab/pharmacy orders, billing integration.",
              },
              {
                title: "Diagnostic Marketplace",
                desc: "Public lab directory, patient test discovery, automated digital report delivery.",
              },
            ].map((module, i) => (
              <div key={i} className="p-6 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm">
                <h3 className="text-base font-bold text-slate-900 dark:text-white mb-1">{module.title}</h3>
                <p className="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">{module.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
