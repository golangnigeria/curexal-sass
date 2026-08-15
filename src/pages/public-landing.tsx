import React from "react";
import { Activity, ShieldCheck, Zap, ArrowRight } from "lucide-react";

export const PublicLandingPage: React.FC = () => {
  return (
    <div className="min-h-screen bg-[#090d16] text-white flex flex-col justify-between p-8">
      {/* Navbar */}
      <header className="flex items-center justify-between border-b border-slate-800 pb-6">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-2xl bg-sky-500/10 border border-sky-500/20 flex items-center justify-center">
            <Activity className="w-6 h-6 text-sky-400" />
          </div>
          <span className="text-2xl font-black tracking-tight bg-gradient-to-r from-sky-400 to-indigo-400 bg-clip-text text-transparent">
            CUREXAL
          </span>
        </div>
        <nav className="flex items-center gap-6 text-sm font-medium text-slate-300">
          <a href="#features" className="hover:text-white transition-colors">Features</a>
          <a href="#pricing" className="hover:text-white transition-colors">Pricing</a>
          <a href="#docs" className="hover:text-white transition-colors">Documentation</a>
        </nav>
      </header>

      {/* Hero Section */}
      <main className="max-w-4xl mx-auto text-center space-y-8 my-16">
        <h1 className="text-5xl md:text-6xl font-black tracking-tight leading-tight">
          Next-Generation Diagnostic & Clinical Management Architecture
        </h1>
        <p className="text-lg text-slate-400 max-w-2xl mx-auto leading-relaxed">
          Unified multi-tenant solution for healthcare networks, clinical diagnostic laboratories, and hospital facility systems.
        </p>
        <div className="flex items-center justify-center gap-4 pt-4">
          <a
            href="/developers"
            className="inline-flex items-center gap-2 px-6 py-3.5 bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold rounded-2xl transition-all shadow-lg shadow-sky-500/20"
          >
            Explore API Specs <ArrowRight className="w-4 h-4" />
          </a>
        </div>
      </main>

      {/* Footer */}
      <footer className="border-t border-slate-800 pt-6 text-center text-xs text-slate-500">
        © 2026 Curexal Inc. All rights reserved. Enterprise Healthcare SaaS Architecture.
      </footer>
    </div>
  );
};
