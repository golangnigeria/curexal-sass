import React from "react";
import { FileText, Paperclip, PhoneCall, FileSpreadsheet, CheckCircle2 } from "lucide-react";

export function HeroSupporting() {
  return (
    <section className="py-16 sm:py-20 bg-slate-50 dark:bg-slate-900/50 border-y border-slate-200/60 dark:border-slate-800/60">
      <div className="max-w-[1280px] mx-auto px-4 sm:px-6">
        <div className="text-center max-w-3xl mx-auto mb-12">
          <h2 className="text-2xl sm:text-4xl font-extrabold text-slate-900 dark:text-white tracking-tight mb-4">
            Healthcare is connected in real life. <br className="hidden sm:inline" />
            <span className="text-[#0F766E] dark:text-teal-400">Why isn't the software?</span>
          </h2>
          <p className="text-sm sm:text-base text-slate-600 dark:text-slate-300 leading-relaxed">
            Patients naturally move across independent providers, from consultations at a clinic to diagnostic testing at a lab, imaging at a diagnostic center, and medications at a pharmacy. Yet the software behind them remains trapped inside individual buildings.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6">
          <div className="bg-white dark:bg-slate-900 p-5 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm space-y-3">
            <div className="w-10 h-10 rounded-xl bg-rose-50 dark:bg-rose-950/60 border border-rose-200 dark:border-rose-800 text-rose-600 dark:text-rose-400 flex items-center justify-center">
              <FileText className="w-5 h-5" />
            </div>
            <h3 className="text-sm font-bold text-slate-900 dark:text-white">Referrals become paper</h3>
            <p className="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
              Patients carry handwritten paper requests from clinic to lab, risking lost context, unreadable notes, and delayed care.
            </p>
          </div>

          <div className="bg-white dark:bg-slate-900 p-5 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm space-y-3">
            <div className="w-10 h-10 rounded-xl bg-amber-50 dark:bg-amber-950/60 border border-amber-200 dark:border-amber-800 text-amber-600 dark:text-amber-400 flex items-center justify-center">
              <Paperclip className="w-5 h-5" />
            </div>
            <h3 className="text-sm font-bold text-slate-900 dark:text-white">Results become attachments</h3>
            <p className="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
              Diagnostic reports get emailed as detached PDFs or printed papers, requiring manual re-entry into patient charts.
            </p>
          </div>

          <div className="bg-white dark:bg-slate-900 p-5 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm space-y-3">
            <div className="w-10 h-10 rounded-xl bg-sky-50 dark:bg-sky-950/60 border border-sky-200 dark:border-sky-800 text-sky-600 dark:text-sky-400 flex items-center justify-center">
              <PhoneCall className="w-5 h-5" />
            </div>
            <h3 className="text-sm font-bold text-slate-900 dark:text-white">Comms become phone & WhatsApp</h3>
            <p className="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
              Doctors and laboratory directors rely on manual phone calls and unorganized messaging threads to follow up on urgent samples.
            </p>
          </div>

          <div className="bg-white dark:bg-slate-900 p-5 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm space-y-3">
            <div className="w-10 h-10 rounded-xl bg-indigo-50 dark:bg-indigo-950/60 border border-indigo-200 dark:border-indigo-800 text-indigo-600 dark:text-indigo-400 flex items-center justify-center">
              <FileSpreadsheet className="w-5 h-5" />
            </div>
            <h3 className="text-sm font-bold text-slate-900 dark:text-white">Settlements become spreadsheets</h3>
            <p className="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
              Inter-organization billing, partner commissions, and supplier payments require tedious end-of-month manual reconciliation.
            </p>
          </div>
        </div>

        <div className="mt-10 p-6 rounded-2xl bg-teal-900/10 border border-teal-500/30 text-center max-w-3xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-3 text-left">
            <div className="w-10 h-10 rounded-full bg-[#0F766E] text-white flex items-center justify-center flex-shrink-0">
              <CheckCircle2 className="w-5 h-5" />
            </div>
            <div>
              <h4 className="text-sm font-bold text-slate-900 dark:text-white">Curexal is designed to connect those workflows.</h4>
              <p className="text-xs text-slate-600 dark:text-slate-300">Preserving organizational independence while establishing automatic, secure data movement.</p>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
