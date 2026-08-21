import React from "react";
import { WifiOff, CreditCard, MessageSquare, RefreshCw, CheckCircle2, Clock } from "lucide-react";

export function AfricaInfrastructureSection() {
  const pillars = [
    {
      icon: WifiOff,
      title: "Connectivity Resilience",
      desc: "Healthcare operations cannot pause during internet downtime. Curexal is designed to queue local operations and continue essential workflows regardless of temporary network failures.",
      status: "Designed For",
      statusType: "building",
    },
    {
      icon: RefreshCw,
      title: "Offline-First Synchronization",
      desc: "Local data is stored securely on device and automatically synchronizes with central servers as soon as internet connectivity is restored, preserving data integrity.",
      status: "Designed For",
      statusType: "building",
    },
    {
      icon: CreditCard,
      title: "African Payment Infrastructure",
      desc: "Built to accommodate local payment methods, mobile money, B2B partner settlements, and split payment workflows tailored to African health facilities.",
      status: "Currently Available",
      statusType: "available",
    },
    {
      icon: MessageSquare,
      title: "WhatsApp & SMS Notifications",
      desc: "Utilizes WhatsApp and SMS messaging for patient test notifications, sample status alerts, and doctor result dispatches where appropriate.",
      status: "Currently Available",
      statusType: "available",
    },
  ];

  return (
    <section className="py-20 bg-white dark:bg-[#0B1120] border-b border-slate-100 dark:border-slate-800">
      <div className="max-w-[1280px] mx-auto px-4 sm:px-6">
        
        {/* Header */}
        <div className="text-center max-w-3xl mx-auto mb-14">
          <h2 className="text-3xl sm:text-4xl font-extrabold text-slate-900 dark:text-white tracking-tight mb-4">
            Built for the realities of healthcare in Africa.
          </h2>
          <p className="text-sm sm:text-base text-slate-600 dark:text-slate-400 leading-relaxed">
            Healthcare software built for Western infrastructure often fails in African clinical environments due to network instability, fragmented payment methods, and communication habits. Curexal is engineered around these realities.
          </p>
        </div>

        {/* Pillars Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {pillars.map((item) => {
            const Icon = item.icon;
            const isAvailable = item.statusType === "available";
            return (
              <div
                key={item.title}
                className="p-6 sm:p-8 rounded-3xl bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 flex flex-col justify-between space-y-4 shadow-xs"
              >
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <div className="w-10 h-10 rounded-xl bg-teal-50 dark:bg-teal-950 text-[#0F766E] flex items-center justify-center">
                      <Icon className="w-5 h-5" />
                    </div>

                    <span
                      className={`text-[10px] font-extrabold uppercase px-2.5 py-1 rounded-full border flex items-center gap-1 ${
                        isAvailable
                          ? "bg-emerald-50 dark:bg-emerald-950/60 text-emerald-700 dark:text-emerald-300 border-emerald-300 dark:border-emerald-800"
                          : "bg-amber-50 dark:bg-amber-950/60 text-amber-700 dark:text-amber-300 border-amber-300 dark:border-amber-800"
                      }`}
                    >
                      {isAvailable ? <CheckCircle2 className="w-3 h-3" /> : <Clock className="w-3 h-3" />}
                      <span>{item.status}</span>
                    </span>
                  </div>

                  <h3 className="text-lg font-bold text-slate-900 dark:text-white">{item.title}</h3>
                  <p className="text-xs sm:text-sm text-slate-600 dark:text-slate-400 leading-relaxed">
                    {item.desc}
                  </p>
                </div>

                <div className="pt-2 border-t border-slate-200/60 dark:border-slate-800 text-[11px] font-semibold text-slate-500 dark:text-slate-500">
                  {isAvailable ? "✓ Active in Curexal Production Architecture" : "⚡ In Active Engineering Roadmap (V2 Core)"}
                </div>
              </div>
            );
          })}
        </div>

      </div>
    </section>
  );
}
