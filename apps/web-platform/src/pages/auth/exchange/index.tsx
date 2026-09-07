import React, { useEffect, useState, useRef } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { ShieldCheck, AlertCircle, RefreshCw } from "lucide-react";

export default function ExchangePage() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [status, setStatus] = useState<"processing" | "success" | "error">("processing");
  const [errorMessage, setErrorMessage] = useState<string>("");
  const executedRef = useRef(false);

  useEffect(() => {
    if (executedRef.current) return;
    executedRef.current = true;

    const rawToken = searchParams.get("token");

    // Immediately remove token from browser address bar and history to prevent exposure/replay
    if (typeof window !== "undefined" && window.history) {
      window.history.replaceState({}, document.title, window.location.pathname);
    }

    if (!rawToken) {
      setStatus("error");
      setErrorMessage("No authorization token provided for workspace exchange.");
      return;
    }

    const exchange = async () => {
      try {
        const res = await fetch("/api/v1/auth/exchange", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include",
          body: JSON.stringify({ token: rawToken }),
        });

        if (!res.ok) {
          const data = await res.json().catch(() => ({}));
          throw new Error(data.message || data.error || "Workspace exchange token expired or invalid.");
        }

        const data = await res.json();
        setStatus("success");

        const targetPath = data.destinationPath || "/workspace/dashboard";
        // Force full page navigation to initialize clean host context & bootstrap
        window.location.replace(targetPath);
      } catch (err: any) {
        setStatus("error");
        setErrorMessage(err.message || "Failed to establish workspace session.");
      }
    };

    exchange();
  }, [searchParams]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-950 px-4">
      <div className="w-full max-w-md bg-slate-900 border border-slate-800 rounded-2xl p-8 shadow-2xl text-center">
        {status === "processing" && (
          <div className="flex flex-col items-center gap-4">
            <div className="relative flex items-center justify-center w-16 h-16 rounded-full bg-cyan-500/10 text-cyan-400 border border-cyan-500/20 animate-pulse">
              <RefreshCw className="w-8 h-8 animate-spin text-cyan-400" />
            </div>
            <h2 className="text-xl font-semibold text-white tracking-tight">Authenticating Workspace</h2>
            <p className="text-sm text-slate-400 max-w-xs">
              Securely exchanging credentials and establishing your clinic session...
            </p>
          </div>
        )}

        {status === "success" && (
          <div className="flex flex-col items-center gap-4">
            <div className="flex items-center justify-center w-16 h-16 rounded-full bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
              <ShieldCheck className="w-8 h-8 text-emerald-400" />
            </div>
            <h2 className="text-xl font-semibold text-white tracking-tight">Access Granted</h2>
            <p className="text-sm text-slate-400 max-w-xs">
              Redirecting to your authorized workspace...
            </p>
          </div>
        )}

        {status === "error" && (
          <div className="flex flex-col items-center gap-4">
            <div className="flex items-center justify-center w-16 h-16 rounded-full bg-rose-500/10 text-rose-400 border border-rose-500/20">
              <AlertCircle className="w-8 h-8 text-rose-400" />
            </div>
            <h2 className="text-xl font-semibold text-white tracking-tight">Handoff Expired</h2>
            <p className="text-sm text-rose-300 max-w-xs">{errorMessage}</p>
            <button
              onClick={() => navigate("/login")}
              className="mt-2 w-full py-2.5 px-4 rounded-xl bg-cyan-600 hover:bg-cyan-500 text-white font-medium text-sm transition-colors shadow-lg shadow-cyan-900/20"
            >
              Return to Sign In
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
