import React from "react";
import { Link, useNavigate } from "react-router-dom";
import {
  Sparkles,
  CalendarPlus,
  Video,
  CheckCircle2,
  Clock,
  Activity,
  ArrowRight,
  ShieldCheck,
  Stethoscope,
  Microscope,
  Pill,
  CreditCard,
  FileCheck,
  AlertCircle,
} from "lucide-react";
import { useMyCareJourney, useMyCareRequests } from "@/api/hooks/use-care-orchestration";

export const PatientPortalDashboardPage: React.FC = () => {
  const navigate = useNavigate();
  const { data: journeyMilestones, isLoading: isJourneyLoading } = useMyCareJourney();
  const { data: careRequests, isLoading: isRequestsLoading } = useMyCareRequests();

  const milestones = Array.isArray(journeyMilestones) ? journeyMilestones : [];
  const requests = Array.isArray(careRequests) ? careRequests : [];

  // Canonical stages in order
  const defaultStages = [
    { code: "INTAKE", label: "Intake & Triage", icon: FileCheck, description: "Demographics & symptoms recorded" },
    { code: "TRIAGE", label: "Provider Matching", icon: Activity, description: "Doctor allocation & urgency assessment" },
    { code: "CONSULTATION", label: "Consultation", icon: Stethoscope, description: "Physical or telehealth session" },
    { code: "LAB_WORKLIST", label: "Diagnostics & Labs", icon: Microscope, description: "Ordered tests and pathology" },
    { code: "PHARMACY_DISPENSE", label: "Medication & Refills", icon: Pill, description: "Prescription fulfillment" },
    { code: "SETTLEMENT", label: "Discharge & Summary", icon: CreditCard, description: "Clinical notes & billing receipt" },
  ];

  // Map milestones to stage statuses
  const getStageStatus = (stageCode: string) => {
    if (!milestones || milestones.length === 0) {
      if (stageCode === "INTAKE") return "COMPLETED";
      if (stageCode === "TRIAGE") return "IN_PROGRESS";
      return "PENDING";
    }
    const found = milestones.find((m) => m.stageCode === stageCode);
    return found?.status || "PENDING";
  };

  return (
    <div className="space-y-8 animate-in fade-in duration-300">
      {/* Welcome Banner matching Public Aesthetic */}
      <div className="relative overflow-hidden p-6 sm:p-8 rounded-3xl bg-gradient-to-r from-teal-900/40 via-slate-900 to-slate-950 border border-teal-500/30 shadow-2xl text-white">
        <div className="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-6">
          <div className="space-y-2">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-teal-500/20 border border-teal-500/30 text-teal-300 text-xs font-semibold uppercase tracking-wider">
              <Sparkles className="w-3.5 h-3.5" />
              <span>Unified Care Continuum</span>
            </div>
            <h1 className="text-2xl sm:text-3xl font-extrabold tracking-tight text-white">
              Welcome to Your Care Workspace
            </h1>
            <p className="text-sm text-slate-300 max-w-xl">
              Track your real-time healthcare progression, access clinical consultations, and review certified diagnostic reports across all connected facilities.
            </p>
          </div>

          <Link
            to="/portal/consultations/new"
            className="flex items-center justify-center gap-2 px-6 py-3 rounded-2xl bg-[#0F766E] hover:bg-[#115E59] text-white text-xs font-bold shadow-lg shadow-teal-900/30 transition-all hover:scale-[1.02] active:scale-[0.98] self-start md:self-center"
          >
            <CalendarPlus className="w-4 h-4" />
            <span>Book Consultation</span>
          </Link>
        </div>
      </div>

      {/* Curexal Care Journey Tracker */}
      <div className="p-6 sm:p-8 rounded-3xl bg-white dark:bg-slate-900/80 border border-slate-200/80 dark:border-slate-800 shadow-xl space-y-6">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 border-b border-slate-200 dark:border-slate-800 pb-4">
          <div>
            <h2 className="text-lg font-bold text-slate-900 dark:text-white flex items-center gap-2">
              <Activity className="w-5 h-5 text-teal-600 dark:text-teal-400" />
              <span>Your Active Care Journey</span>
            </h2>
            <p className="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
              Live progression from initial intake through doctor consultation, diagnostic tests, and pharmacy fulfillment
            </p>
          </div>
          <span className="text-[11px] font-mono text-teal-600 dark:text-teal-400 bg-teal-500/10 px-3 py-1 rounded-full border border-teal-500/20 self-start sm:self-center font-bold">
            Orchestration Active
          </span>
        </div>

        {/* Milestone Steps Timeline */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-6 gap-4 relative">
          {defaultStages.map((stage, idx) => {
            const status = getStageStatus(stage.code);
            const Icon = stage.icon;

            const isCompleted = status === "COMPLETED";
            const isInProgress = status === "IN_PROGRESS";

            return (
              <div
                key={stage.code}
                className={`relative flex flex-col p-4 rounded-2xl border transition-all ${
                  isCompleted
                    ? "bg-emerald-500/5 dark:bg-emerald-950/20 border-emerald-500/30 text-emerald-700 dark:text-emerald-300"
                    : isInProgress
                    ? "bg-teal-500/10 dark:bg-teal-950/40 border-teal-500/50 text-teal-800 dark:text-teal-200 ring-1 ring-teal-500/30 shadow-lg shadow-teal-500/10"
                    : "bg-slate-50 dark:bg-slate-950/40 border-slate-200 dark:border-slate-800 text-slate-400 dark:text-slate-500"
                }`}
              >
                {/* Stage Header */}
                <div className="flex items-center justify-between mb-3">
                  <div
                    className={`w-8 h-8 rounded-xl flex items-center justify-center font-bold text-xs ${
                      isCompleted
                        ? "bg-emerald-500/20 text-emerald-600 dark:text-emerald-400 border border-emerald-500/30"
                        : isInProgress
                        ? "bg-teal-500/20 text-teal-600 dark:text-teal-300 border border-teal-500/40 animate-pulse"
                        : "bg-slate-200 dark:bg-slate-800 text-slate-500 border border-slate-300 dark:border-slate-700/50"
                    }`}
                  >
                    {isCompleted ? <CheckCircle2 className="w-4 h-4" /> : <Icon className="w-4 h-4" />}
                  </div>
                  <span className="text-[10px] font-mono text-slate-400 dark:text-slate-500 font-bold">0{idx + 1}</span>
                </div>

                {/* Title & Description */}
                <div className="font-bold text-xs text-slate-900 dark:text-white">{stage.label}</div>
                <p className="text-[11px] text-slate-500 dark:text-slate-400 mt-1 leading-relaxed">{stage.description}</p>

                {/* Status Badge */}
                <div className="mt-4 pt-2 border-t border-slate-200 dark:border-slate-800/60">
                  {isCompleted ? (
                    <span className="inline-flex items-center gap-1 text-[10px] font-bold text-emerald-600 dark:text-emerald-400">
                      <CheckCircle2 className="w-3 h-3" /> Completed
                    </span>
                  ) : isInProgress ? (
                    <span className="inline-flex items-center gap-1 text-[10px] font-bold text-teal-600 dark:text-teal-400 animate-pulse">
                      <Clock className="w-3 h-3" /> Active Now
                    </span>
                  ) : (
                    <span className="text-[10px] font-medium text-slate-400 dark:text-slate-500">Upcoming</span>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* Active Care Requests & Quick Cards */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left: Active Care Demands (2 cols) */}
        <div className="lg:col-span-2 p-6 rounded-3xl bg-white dark:bg-slate-900/80 border border-slate-200/80 dark:border-slate-800 shadow-xl space-y-4">
          <div className="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-3">
            <h3 className="text-sm font-bold text-slate-900 dark:text-white flex items-center gap-2">
              <Clock className="w-4 h-4 text-teal-600 dark:text-teal-400" />
              <span>Active Consultations & Requests</span>
            </h3>
            <Link
              to="/portal/consultations/new"
              className="text-xs text-teal-600 dark:text-teal-400 hover:underline font-bold flex items-center gap-1"
            >
              <span>Request Care</span>
              <ArrowRight className="w-3 h-3" />
            </Link>
          </div>

          {isRequestsLoading ? (
            <div className="p-8 text-center text-xs text-slate-400 animate-pulse">
              Loading your care requests...
            </div>
          ) : requests.length === 0 ? (
            <div className="p-8 text-center space-y-3">
              <div className="w-12 h-12 rounded-2xl bg-teal-50 dark:bg-slate-800 text-teal-600 dark:text-slate-400 flex items-center justify-center mx-auto">
                <Stethoscope className="w-6 h-6" />
              </div>
              <div className="text-sm font-bold text-slate-900 dark:text-white">No active care requests</div>
              <p className="text-xs text-slate-500 dark:text-slate-400 max-w-sm mx-auto">
                You currently have no pending consultation or diagnostic orders. Click below to book an appointment.
              </p>
              <Link
                to="/portal/consultations/new"
                className="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-teal-500/10 text-teal-700 dark:text-teal-300 border border-teal-500/20 text-xs font-bold hover:bg-teal-500/20 transition-all"
              >
                <CalendarPlus className="w-3.5 h-3.5" />
                <span>Start New Care Request</span>
              </Link>
            </div>
          ) : (
            <div className="space-y-3">
              {requests.map((req) => (
                <div
                  key={req.id}
                  className="p-4 rounded-2xl bg-slate-50 dark:bg-slate-950/60 border border-slate-200 dark:border-slate-800 hover:border-teal-500/40 transition-all flex flex-col sm:flex-row sm:items-center justify-between gap-4"
                >
                  <div className="space-y-1">
                    <div className="flex items-center gap-2">
                      <span className="font-bold text-slate-900 dark:text-white text-xs">{req.serviceType}</span>
                      <span className="px-2 py-0.5 rounded-md bg-slate-200 dark:bg-slate-800 text-slate-700 dark:text-teal-300 font-mono text-[10px] font-bold">
                        {req.requestNumber}
                      </span>
                      <span className="px-2 py-0.5 rounded-md bg-teal-500/10 text-teal-700 dark:text-teal-300 text-[10px] font-semibold">
                        {req.status}
                      </span>
                    </div>
                    <p className="text-xs text-slate-500 dark:text-slate-400">
                      Reason: <span className="text-slate-800 dark:text-slate-200 font-medium">{req.chiefComplaint || "General checkup"}</span>
                    </p>
                    <div className="flex items-center gap-3 text-[11px] text-slate-400 pt-1">
                      <span>Mode: <strong className="text-slate-700 dark:text-slate-300">{req.preferredMode}</strong></span>
                      <span>•</span>
                      <span>Urgency: <strong className="text-slate-700 dark:text-slate-300">{req.urgency}</strong></span>
                    </div>
                  </div>

                  <div className="flex items-center gap-2 self-end sm:self-center">
                    {req.preferredMode === "VIDEO" && (
                      <button
                        onClick={() => navigate(`/portal/consultations/${req.id}/room`)}
                        className="flex items-center gap-1.5 px-4 py-2 rounded-xl bg-[#0F766E] hover:bg-[#115E59] text-white text-xs font-bold shadow-md shadow-teal-900/20 transition-all"
                      >
                        <Video className="w-3.5 h-3.5" />
                        <span>Join Telehealth</span>
                      </button>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Right: Health Security & Privacy Card */}
        <div className="p-6 rounded-3xl bg-white dark:bg-slate-900/80 border border-slate-200/80 dark:border-slate-800 shadow-xl space-y-4">
          <h3 className="text-sm font-bold text-slate-900 dark:text-white flex items-center gap-2">
            <ShieldCheck className="w-4 h-4 text-emerald-600 dark:text-emerald-400" />
            <span>Health Privacy & Consent</span>
          </h3>
          <p className="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
            Your clinical data and diagnostic history are strictly isolated within accredited healthcare networks and protected by tenant encryption.
          </p>

          <div className="space-y-2.5 pt-2">
            <div className="p-3 rounded-xl bg-slate-50 dark:bg-slate-950/60 border border-slate-200 dark:border-slate-800 text-xs flex items-center justify-between">
              <span className="text-slate-700 dark:text-slate-300 font-medium">Telehealth Virtual Consent</span>
              <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20">
                ACTIVE
              </span>
            </div>
            <div className="p-3 rounded-xl bg-slate-50 dark:bg-slate-950/60 border border-slate-200 dark:border-slate-800 text-xs flex items-center justify-between">
              <span className="text-slate-700 dark:text-slate-300 font-medium">Diagnostic Lab Sharing</span>
              <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20">
                ACTIVE
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
