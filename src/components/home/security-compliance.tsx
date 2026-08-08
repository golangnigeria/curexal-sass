import { useState } from "react";
import {
  ShieldCheck,
  FileText,
  Users,
  Lock,
  Server,
  Layers,
  Database,
  FileCheck,
  CheckCircle2,
  Zap,
  Shield,
  Terminal,
  Key,
  RefreshCw,
  AlertOctagon,
  Play,
  Check,
} from "lucide-react";
import { cn } from "@/lib/utils";

export function SecurityCompliance() {
  const [activeTab, setActiveTab] = useState<"sovereignty" | "audit" | "zerotrust">("sovereignty");
  
  // Interactive Simulation State for Tab 1 (Data Sovereignty)
  const [selectedTenant, setSelectedTenant] = useState<"apollo" | "lifecare" | "meridian">("apollo");
  const [simulatedAttempt, setSimulatedAttempt] = useState<boolean>(false);

  // Interactive Simulation State for Tab 2 (Audit Logs)
  const [logItems, setLogItems] = useState([
    { id: "1", time: "16:42:01.002", details: "PATHOLOGIST_SIGN_OFF: FBC #8492 signed", hash: "8f9a2b7c" },
    { id: "2", time: "16:40:15.891", details: "SPECIMEN_ACCESSION: Blood #SP-9921", hash: "3c1d4e5f" },
  ]);
  const [isAppendingLog, setIsAppendingLog] = useState(false);

  // Interactive Simulation State for Tab 3 (Encryption)
  const [keyRotated, setKeyRotated] = useState(false);

  const handleSimulateCrossAccess = () => {
    setSimulatedAttempt(true);
    setTimeout(() => setSimulatedAttempt(false), 3000);
  };

  const handleSimulateLogEvent = () => {
    setIsAppendingLog(true);
    const newLog = {
      id: String(Date.now()),
      time: new Date().toISOString().substring(11, 23),
      details: "ANALYZER_INGESTION: CBC Run ingested",
      hash: Math.random().toString(16).substring(2, 10),
    };
    setLogItems((prev) => [newLog, ...prev.slice(0, 1)]);
    setTimeout(() => setIsAppendingLog(false), 500);
  };

  return (
    <section
      id="security"
      className="py-14 bg-[#0B1120] text-white border-y border-slate-800 relative overflow-hidden"
    >
      {/* Background radial gradient */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[300px] bg-teal-500/10 blur-[120px] rounded-full pointer-events-none" />

      <div className="max-w-[1280px] mx-auto px-6 relative z-10">

        {/* Compact Header */}
        <div className="text-center max-w-2xl mx-auto mb-8">
          <div className="inline-flex items-center gap-1.5 px-3 py-1 mb-3 rounded-full border border-teal-500/30 bg-teal-500/10 text-teal-300 text-xs font-semibold">
            <Shield className="h-3.5 w-3.5 text-teal-400" />
            <span>Enterprise Security Architecture</span>
          </div>
          <h2 className="text-2xl sm:text-3xl font-bold text-white mb-2">
            Network-Grade Security. Data Sovereignty by Design.
          </h2>
          <p className="text-xs sm:text-sm text-slate-400">
            Isolated PostgreSQL database schemas, cryptographic SHA-256 audit trails, and end-to-end TLS 1.3 encryption.
          </p>
        </div>

        {/* Compact Tab Bar */}
        <div className="flex items-center justify-center gap-2 mb-6">
          {[
            { id: "sovereignty", label: "Schema Isolation", icon: Database },
            { id: "audit", label: "Cryptographic Audit Trail", icon: FileCheck },
            { id: "zerotrust", label: "Zero-Trust & Encryption", icon: Lock },
          ].map((tab) => {
            const TabIcon = tab.icon;
            const isActive = activeTab === tab.id;
            return (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id as any)}
                className={cn(
                  "px-3.5 py-1.5 rounded-[10px] text-xs font-semibold transition-all duration-200 cursor-pointer border flex items-center gap-2",
                  isActive
                    ? "bg-[#0F766E] border-teal-400 text-white shadow-sm"
                    : "bg-slate-900/80 border-slate-800 text-slate-400 hover:text-white hover:bg-slate-800"
                )}
              >
                <TabIcon className={cn("w-3.5 h-3.5", isActive ? "text-white" : "text-slate-400")} />
                <span>{tab.label}</span>
              </button>
            );
          })}
        </div>

        {/* Main Compact Dual Console Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 items-stretch mb-8">

          {/* Left Interactive Terminal */}
          <div className="lg:col-span-7 rounded-[16px] bg-slate-900/90 border border-slate-800 p-5 shadow-xl flex flex-col justify-between">
            
            {/* Header */}
            <div className="flex items-center justify-between pb-3 mb-4 border-b border-slate-800">
              <div className="flex items-center gap-2 text-xs font-mono text-slate-400">
                <Terminal className="w-3.5 h-3.5 text-teal-400" />
                <span>curexal-security-node://v2.4</span>
              </div>
              <span className="text-[11px] font-mono text-emerald-400 font-semibold flex items-center gap-1">
                <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse" />
                Enforced
              </span>
            </div>

            {/* TAB 1: SCHEMA ISOLATION CONSOLE */}
            {activeTab === "sovereignty" && (
              <div className="space-y-4 text-xs animate-fade-in">
                <div className="flex items-center justify-between">
                  <span className="text-slate-300 font-semibold">Select Tenant Schema:</span>
                  <div className="flex gap-1.5">
                    {[
                      { id: "apollo", label: "Apollo Diag" },
                      { id: "lifecare", label: "LifeCare Lab" },
                      { id: "meridian", label: "Meridian Health" },
                    ].map((t) => (
                      <button
                        key={t.id}
                        onClick={() => setSelectedTenant(t.id as any)}
                        className={cn(
                          "px-2.5 py-1 rounded-[6px] text-[11px] font-mono border transition-all cursor-pointer",
                          selectedTenant === t.id
                            ? "bg-teal-950 border-teal-400 text-teal-300 font-bold"
                            : "bg-slate-950 border-slate-800 text-slate-400"
                        )}
                      >
                        {t.label}
                      </button>
                    ))}
                  </div>
                </div>

                <div className="p-3 rounded-[10px] bg-slate-950 font-mono text-[11px] border border-slate-800 space-y-1.5">
                  <div className="text-slate-400">$ psql -c "SHOW search_path;"</div>
                  <div className="text-emerald-400 font-semibold">
                    search_path = "tenant_{selectedTenant === "apollo" ? "apollo_diag" : selectedTenant === "lifecare" ? "lifecare_lab" : "meridian_health"}", public
                  </div>

                  {simulatedAttempt && (
                    <div className="mt-2 p-2 rounded bg-red-950/60 border border-red-500/40 text-red-300 text-[10.5px] animate-fade-in flex items-center gap-2">
                      <AlertOctagon className="w-3.5 h-3.5 text-red-400 flex-shrink-0" />
                      <span>[DENIED] 403 Forbidden: Cross-tenant schema boundary enforced.</span>
                    </div>
                  )}
                </div>

                <div className="flex items-center justify-between pt-1">
                  <button
                    onClick={handleSimulateCrossAccess}
                    disabled={simulatedAttempt}
                    className="px-3 py-1.5 rounded-[8px] bg-red-500/20 hover:bg-red-500/30 text-red-300 border border-red-500/30 text-[11px] font-bold transition-all cursor-pointer flex items-center gap-1.5"
                  >
                    <Play className="w-3 h-3 text-red-400" />
                    Test Breach Guard
                  </button>
                  <span className="text-[11px] text-slate-500 font-mono">Zero Cross-Talk</span>
                </div>
              </div>
            )}

            {/* TAB 2: AUDIT TRAIL CONSOLE */}
            {activeTab === "audit" && (
              <div className="space-y-3 text-xs animate-fade-in">
                <div className="p-3 rounded-[10px] bg-slate-950 font-mono text-[11px] border border-slate-800 space-y-1.5">
                  <div className="text-slate-400 border-b border-slate-800 pb-1 flex justify-between">
                    <span>UTC Time</span>
                    <span>Action &amp; Details</span>
                    <span>HMAC Hash</span>
                  </div>

                  {logItems.map((item) => (
                    <div
                      key={item.id}
                      className={cn(
                        "p-1.5 rounded flex items-center justify-between text-[11px]",
                        isAppendingLog && item.id === logItems[0].id ? "bg-teal-950 border border-teal-500/40" : "bg-slate-900/50"
                      )}
                    >
                      <span className="text-slate-500">{item.time}</span>
                      <span className="text-teal-300 truncate max-w-[180px]">{item.details}</span>
                      <span className="text-emerald-400 font-mono">sig:{item.hash}</span>
                    </div>
                  ))}
                </div>

                <div className="flex items-center justify-between pt-1">
                  <button
                    onClick={handleSimulateLogEvent}
                    className="px-3 py-1.5 rounded-[8px] bg-teal-500/20 hover:bg-teal-500/30 text-teal-300 border border-teal-500/30 text-[11px] font-bold transition-all cursor-pointer flex items-center gap-1.5"
                  >
                    <Zap className="w-3 h-3 text-teal-400" />
                    Append Log Event
                  </button>
                  <span className="text-[11px] text-emerald-400 font-mono">SHA-256 Verified</span>
                </div>
              </div>
            )}

            {/* TAB 3: ENCRYPTION CONSOLE */}
            {activeTab === "zerotrust" && (
              <div className="space-y-3 text-xs animate-fade-in">
                <div className="grid grid-cols-2 gap-2 text-[11px] font-mono">
                  <div className="p-2.5 rounded bg-slate-950 border border-slate-800">
                    <span className="text-slate-500 block text-[10px]">In Transit</span>
                    <span className="text-white font-bold block mt-0.5">TLS 1.3 (AES-256)</span>
                  </div>
                  <div className="p-2.5 rounded bg-slate-950 border border-slate-800">
                    <span className="text-slate-500 block text-[10px]">Storage Key</span>
                    <span className="text-white font-bold block mt-0.5">
                      {keyRotated ? "Key Rotated" : "AES-256 Active"}
                    </span>
                  </div>
                </div>

                <div className="flex items-center justify-between pt-1">
                  <button
                    onClick={() => setKeyRotated(!keyRotated)}
                    className="px-3 py-1.5 rounded-[8px] bg-emerald-500/20 hover:bg-emerald-500/30 text-emerald-300 border border-emerald-500/30 text-[11px] font-bold transition-all cursor-pointer flex items-center gap-1.5"
                  >
                    <RefreshCw className="w-3 h-3 text-emerald-400" />
                    {keyRotated ? "Reset Key" : "Rotate KMS Key"}
                  </button>
                  <span className="text-[11px] text-slate-400 font-mono">Zero Trust Enforced</span>
                </div>
              </div>
            )}
          </div>

          {/* Right Guarantees Box */}
          <div className="lg:col-span-5 rounded-[16px] bg-slate-900/60 border border-slate-800 p-5 flex flex-col justify-between">
            <div>
              <div className="flex items-center gap-2 mb-3">
                <ShieldCheck className="w-4 h-4 text-teal-400" />
                <h3 className="text-sm font-bold text-white">Security &amp; Privacy Guarantees</h3>
              </div>

              <div className="space-y-2 mb-4 text-xs">
                {[
                  "Schema-per-tenant isolation for 100% data sovereignty",
                  "Append-only SHA-256 HMAC cryptographic audit log",
                  "Pathologist digital e-signatures on every report",
                  "AES-256 storage volume encryption & TLS 1.3 transit",
                  "Role-Based Access Control (RBAC) across all roles",
                ].map((g, i) => (
                  <div key={i} className="flex items-start gap-2">
                    <Check className="w-3.5 h-3.5 text-teal-400 flex-shrink-0 mt-0.5" />
                    <span className="text-slate-300 font-medium leading-snug">{g}</span>
                  </div>
                ))}
              </div>
            </div>

            <div className="pt-3 border-t border-slate-800 flex flex-wrap gap-1.5">
              {["ISO 15189 Aligned", "HIPAA Privacy Controls", "NDPR Protection", "AES-256 Encrypted"].map((b) => (
                <span
                  key={b}
                  className="px-2 py-0.5 rounded text-[10.5px] font-semibold text-teal-300 bg-teal-500/10 border border-teal-500/20"
                >
                  {b}
                </span>
              ))}
            </div>
          </div>

        </div>

        {/* Compact 4-Badge Ribbon */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-3 text-center">
          {[
            { label: "ISO 15189 Aligned", sub: "Pathology Quality Workflows" },
            { label: "Immutable Audit Log", sub: "SHA-256 Signed History" },
            { label: "Fine-Grained RBAC", sub: "Role Isolation Enforced" },
            { label: "AES-256 & TLS 1.3", sub: "Data Encryption" },
          ].map((item) => (
            <div key={item.label} className="p-3 rounded-[12px] bg-slate-900/60 border border-slate-800">
              <span className="text-xs font-bold text-white block">{item.label}</span>
              <span className="text-[11px] text-slate-400 block mt-0.5">{item.sub}</span>
            </div>
          ))}
        </div>

      </div>
    </section>
  );
}
