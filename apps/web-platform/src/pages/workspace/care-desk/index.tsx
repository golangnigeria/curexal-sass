import React, { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import {
  Activity,
  HeartPulse,
  Users,
  Search,
  Filter,
  CheckCircle2,
  Clock,
  AlertTriangle,
  Flame,
  Stethoscope,
  Video,
  Building2,
  Phone,
  ChevronRight,
  ShieldCheck,
  Sparkles,
  ArrowRight,
  Loader2,
  X,
  Sliders,
} from "lucide-react";
import {
  useCareRequests,
  useSubmitTriage,
  useMatchProviders,
  useAssignProvider,
} from "../../../api/hooks/use-care-orchestration";
import type {
  CareRequest,
  MatchedProviderCandidate,
  SubmitTriagePayload,
} from "../../../api/contracts";

export const CareAgentDeskWorkspacePage: React.FC = () => {
  const { branchSlug } = useParams<{ branchSlug: string }>();
  const navigate = useNavigate();

  const [statusFilter, setStatusFilter] = useState("");
  const [urgencyFilter, setUrgencyFilter] = useState("");
  const [channelFilter, setChannelFilter] = useState("");
  const [activeRequest, setActiveRequest] = useState<CareRequest | null>(null);
  const [isTriageOpen, setIsTriageOpen] = useState(false);
  const [isMatchOpen, setIsMatchOpen] = useState(false);

  const { data: requestList, isLoading, refetch } = useCareRequests({
    status: statusFilter || undefined,
    urgency: urgencyFilter || undefined,
    limit: 50,
  });

  // Triage state
  const [systolicBp, setSystolicBp] = useState<number | undefined>();
  const [diastolicBp, setDiastolicBp] = useState<number | undefined>();
  const [pulseRate, setPulseRate] = useState<number | undefined>();
  const [temperature, setTemperature] = useState<number | undefined>(37.0);
  const [spo2, setSpo2] = useState<number | undefined>(98);
  const [respiratoryRate, setRespiratoryRate] = useState<number | undefined>(18);
  const [painScore, setPainScore] = useState<number | undefined>(0);
  const [triageNotes, setTriageNotes] = useState("");

  const submitTriage = useSubmitTriage();
  const assignProvider = useAssignProvider();

  // Provider matching query for active request
  const { data: matchingData, isLoading: isMatchingLoading } = useMatchProviders(
    activeRequest?.id || ""
  );

  // Live Auto-Acuity calculation
  const calculateLiveAcuity = () => {
    if (
      (temperature && temperature >= 39.5) ||
      (spo2 && spo2 <= 91) ||
      (painScore && painScore >= 8) ||
      (systolicBp && (systolicBp >= 180 || systolicBp <= 85))
    ) {
      return { level: "RED", label: "Emergency (Immediate)", color: "text-rose-400 bg-rose-500/10 border-rose-500/30" };
    }
    if (
      (temperature && temperature >= 38.2) ||
      (spo2 && spo2 <= 94) ||
      (painScore && painScore >= 5) ||
      (systolicBp && (systolicBp >= 140 || systolicBp <= 95))
    ) {
      return { level: "YELLOW", label: "Urgent Priority", color: "text-amber-400 bg-amber-500/10 border-amber-500/30" };
    }
    return { level: "GREEN", label: "Standard / Routine", color: "text-emerald-400 bg-emerald-500/10 border-emerald-500/30" };
  };

  const liveAcuity = calculateLiveAcuity();

  const handleSaveTriage = (e: React.FormEvent) => {
    e.preventDefault();
    if (!activeRequest) return;

    submitTriage.mutate(
      {
        requestId: activeRequest.id,
        payload: {
          systolicBp,
          diastolicBp,
          pulseRate,
          temperature,
          spo2,
          respiratoryRate,
          painScore,
          triageNotes,
        },
      },
      {
        onSuccess: () => {
          setIsTriageOpen(false);
          setIsMatchOpen(true);
          refetch();
        },
      }
    );
  };

  const handleAssignDoctor = (provider: MatchedProviderCandidate) => {
    if (!activeRequest) return;

    assignProvider.mutate(
      {
        requestId: activeRequest.id,
        providerId: provider.providerId,
      },
      {
        onSuccess: () => {
          setIsMatchOpen(false);
          setActiveRequest(null);
          refetch();
        },
      }
    );
  };

  const allItems: CareRequest[] = requestList?.items || [];
  const items = allItems.filter((i: CareRequest) => {
    if (!channelFilter) return true;
    const mode = (i.preferredMode || i.deliveryChannel || "").toLowerCase();
    if (channelFilter === "video") {
      return mode === "video";
    }
    if (channelFilter === "in_person") {
      return mode === "in_person" || mode === "";
    }
    return true;
  });
  const emergencyCount = allItems.filter((i: CareRequest) => i.urgency === "EMERGENCY").length;
  const pendingTriageCount = allItems.filter((i: CareRequest) => i.status === "SUBMITTED" || i.status === "WAITING_TRIAGE").length;

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 p-6 space-y-6 animate-in fade-in duration-300">
      {/* Header Banner */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 p-6 rounded-2xl bg-gradient-to-r from-slate-900 via-slate-900 to-indigo-950/40 border border-slate-800 shadow-xl">
        <div className="flex items-center gap-4">
          <div className="p-3.5 rounded-2xl bg-indigo-500/10 text-indigo-400 border border-indigo-500/20 shadow-inner">
            <Activity className="w-7 h-7" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h1 className="text-2xl font-bold tracking-tight text-white">
                Care Coordination & Triage Desk
              </h1>
              <span className="px-2.5 py-0.5 rounded-full text-[11px] font-semibold bg-indigo-500/10 text-indigo-400 border border-indigo-500/20 uppercase tracking-wider">
                Intake Engine
              </span>
            </div>
            <p className="text-xs text-slate-400 mt-1">
              Clinical triage vitals assessment, HMO intake verification, and intelligent doctor allocation for branch{" "}
              <span className="text-indigo-300 font-mono">/{branchSlug}</span>
            </p>
          </div>
        </div>
      </div>

      {/* Metrics Row */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div className="p-4 rounded-xl bg-slate-900 border border-slate-800 flex items-center justify-between">
          <div>
            <div className="text-xs text-slate-400 font-medium">Pending Triage / Intake</div>
            <div className="text-2xl font-bold text-white mt-1">{pendingTriageCount}</div>
          </div>
          <div className="p-2.5 rounded-lg bg-slate-800 text-cyan-400">
            <Clock className="w-5 h-5" />
          </div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900 border border-slate-800 flex items-center justify-between">
          <div>
            <div className="text-xs text-slate-400 font-medium">Emergency Escalations</div>
            <div className="text-2xl font-bold text-rose-400 mt-1">{emergencyCount}</div>
          </div>
          <div className="p-2.5 rounded-lg bg-rose-500/10 text-rose-400">
            <Flame className="w-5 h-5" />
          </div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900 border border-slate-800 flex items-center justify-between">
          <div>
            <div className="text-xs text-slate-400 font-medium">Available Doctors</div>
            <div className="text-2xl font-bold text-emerald-400 mt-1">On Duty</div>
          </div>
          <div className="p-2.5 rounded-lg bg-emerald-500/10 text-emerald-400">
            <Stethoscope className="w-5 h-5" />
          </div>
        </div>
      </div>

      {/* Filter Bar */}
      <div className="flex flex-col sm:flex-row items-center gap-3 p-3 bg-slate-900 border border-slate-800 rounded-xl">
        <div className="flex items-center gap-2 w-full sm:w-auto">
          <span className="text-xs text-slate-400 font-medium">Status:</span>
          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            className="px-3 py-1.5 bg-slate-950 border border-slate-800 rounded-lg text-xs text-white focus:outline-none focus:border-cyan-500"
          >
            <option value="">All Care Demands</option>
            <option value="SUBMITTED">Pending Triage (Submitted)</option>
            <option value="TRIAGED">Triaged (Ready for Match)</option>
            <option value="MATCHED">Doctor Matched</option>
          </select>
        </div>

        <div className="flex items-center gap-2 w-full sm:w-auto">
          <span className="text-xs text-slate-400 font-medium">Urgency:</span>
          <select
            value={urgencyFilter}
            onChange={(e) => setUrgencyFilter(e.target.value)}
            className="px-3 py-1.5 bg-slate-950 border border-slate-800 rounded-lg text-xs text-white focus:outline-none focus:border-cyan-500"
          >
            <option value="">All Acuity Levels</option>
            <option value="ROUTINE">Routine / Standard</option>
            <option value="URGENT">Urgent Priority</option>
            <option value="EMERGENCY">Emergency Alert</option>
          </select>
        </div>

        <div className="flex items-center gap-2 w-full sm:w-auto">
          <span className="text-xs text-slate-400 font-medium">Channel:</span>
          <select
            value={channelFilter}
            onChange={(e) => setChannelFilter(e.target.value)}
            className="px-3 py-1.5 bg-slate-950 border border-slate-800 rounded-lg text-xs text-white focus:outline-none focus:border-cyan-500"
          >
            <option value="">All Delivery Channels</option>
            <option value="in_person">In-Person (Physical Clinic)</option>
            <option value="video">Telehealth (Virtual WebRTC)</option>
          </select>
        </div>
      </div>

      {/* Care Requests List */}
      <div className="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden shadow-xl">
        <div className="px-6 py-4 border-b border-slate-800 flex items-center justify-between">
          <h2 className="text-sm font-semibold text-white">Inbound Care Demand Worklist</h2>
          <span className="text-xs text-slate-400">{items.length} requests in queue</span>
        </div>

        {isLoading ? (
          <div className="p-12 text-center text-slate-400 text-xs animate-pulse">
            Loading care requests worklist...
          </div>
        ) : items.length === 0 ? (
          <div className="p-12 text-center space-y-3">
            <div className="w-12 h-12 rounded-full bg-slate-800 text-slate-500 flex items-center justify-center mx-auto">
              <CheckCircle2 className="w-6 h-6" />
            </div>
            <div className="text-sm font-medium text-white">Worklist is clear</div>
            <p className="text-xs text-slate-400">All inbound care requests have been triaged and allocated.</p>
          </div>
        ) : (
          <div className="divide-y divide-slate-800">
            {items.map((req: CareRequest) => {
              const isRed = req.urgency === "EMERGENCY";
              const isYellow = req.urgency === "URGENT";
              return (
                <div
                  key={req.id}
                  className={`p-5 hover:bg-slate-800/40 transition-colors flex flex-col md:flex-row md:items-center justify-between gap-4 ${
                    isRed ? "bg-rose-950/20 border-l-4 border-l-rose-500" : ""
                  }`}
                >
                  <div className="space-y-1.5">
                    <div className="flex items-center gap-2">
                      <span className="font-semibold text-white text-sm">
                        {req.patientName || "Walk-In / Portal Patient"}
                      </span>
                      {req.mrn && (
                        <span className="px-2 py-0.5 rounded-md bg-slate-800 text-cyan-400 font-mono text-[11px] font-semibold border border-slate-700">
                          {req.mrn}
                        </span>
                      )}
                      <span className="px-2 py-0.5 rounded-md bg-slate-800 text-slate-400 font-mono text-[10px]">
                        {req.requestNumber}
                      </span>
                      <span
                        className={`px-2 py-0.5 rounded-md text-[10px] font-semibold border ${
                          isRed
                            ? "bg-rose-500/20 text-rose-400 border-rose-500/40 animate-pulse"
                            : isYellow
                            ? "bg-amber-500/10 text-amber-400 border-amber-500/30"
                            : "bg-emerald-500/10 text-emerald-400 border-emerald-500/20"
                        }`}
                      >
                        {req.urgency}
                      </span>
                      <span className="px-2 py-0.5 rounded-md bg-cyan-500/10 text-cyan-400 text-[10px] font-medium border border-cyan-500/20">
                        {req.status}
                      </span>
                      {(() => {
                        const elapsedMins = Math.max(0, Math.floor((Date.now() - new Date(req.createdAt).getTime()) / 60000));
                        const isSlaBreached = elapsedMins >= 60;
                        const isSlaWarning = elapsedMins >= 30 && elapsedMins < 60;
                        return (
                          <span
                            className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[11px] font-mono font-medium border ${
                              isSlaBreached
                                ? "bg-rose-500/20 text-rose-300 border-rose-500/40 animate-pulse"
                                : isSlaWarning
                                ? "bg-amber-500/10 text-amber-300 border-amber-500/30"
                                : "bg-slate-800 text-slate-300 border-slate-700"
                            }`}
                          >
                            <Clock className="w-3 h-3 text-slate-400" />
                            <span>{elapsedMins}m wait</span>
                          </span>
                        );
                      })()}
                    </div>

                    <p className="text-xs text-slate-300">
                      Chief Complaint: <strong>{req.chiefComplaint || "General medical consultation"}</strong>
                    </p>

                    <div className="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-slate-400">
                      <span>Service: <strong className="text-slate-200">{req.serviceType}</strong></span>
                      <span className="inline-flex items-center gap-1">
                        Channel:{" "}
                        {req.preferredMode === "VIDEO" || req.preferredMode === "video" ? (
                          <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded bg-purple-500/10 text-purple-300 border border-purple-500/20 text-[10px] font-semibold">
                            <Video className="w-3 h-3 text-purple-400" />
                            <span>TELEHEALTH (WebRTC Room)</span>
                          </span>
                        ) : (
                          <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded bg-blue-500/10 text-blue-300 border border-blue-500/20 text-[10px] font-semibold">
                            <Building2 className="w-3 h-3 text-blue-400" />
                            <span>IN-PERSON (Suite 104)</span>
                          </span>
                        )}
                      </span>
                      {req.symptomsJson && req.symptomsJson.length > 0 && (
                        <span>Symptoms: <strong className="text-slate-300">{req.symptomsJson.join(", ")}</strong></span>
                      )}
                    </div>
                  </div>

                  {/* Actions */}
                  <div className="flex items-center gap-2 self-end md:self-center">
                    <button
                      onClick={() => {
                        setActiveRequest(req);
                        setIsTriageOpen(true);
                      }}
                      className="flex items-center gap-1.5 px-3.5 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold border border-slate-700 transition-colors"
                    >
                      <HeartPulse className="w-3.5 h-3.5 text-rose-400" />
                      <span>Take Triage Vitals</span>
                    </button>
                    <button
                      onClick={() => {
                        setActiveRequest(req);
                        setIsMatchOpen(true);
                      }}
                      className="flex items-center gap-1.5 px-4 py-2 rounded-xl bg-gradient-to-r from-cyan-500 to-blue-600 hover:from-cyan-400 hover:to-blue-500 text-white text-xs font-semibold shadow-lg shadow-cyan-500/20 transition-all"
                    >
                      <Stethoscope className="w-3.5 h-3.5" />
                      <span>Match Doctor</span>
                    </button>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* 1. CLINICAL TRIAGE VITALS MODAL */}
      {isTriageOpen && activeRequest && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/70 backdrop-blur-sm animate-in fade-in duration-200">
          <div className="relative w-full max-w-2xl bg-slate-900 border border-slate-800 rounded-3xl shadow-2xl text-slate-100 overflow-hidden">
            <div className="flex items-center justify-between px-6 py-4 bg-slate-900/90 border-b border-slate-800">
              <div className="flex items-center gap-3">
                <div className="p-2.5 rounded-xl bg-rose-500/10 text-rose-400 border border-rose-500/20">
                  <HeartPulse className="w-5 h-5" />
                </div>
                <div>
                  <h3 className="text-base font-bold text-white">Clinical Triage & Vitals Intake</h3>
                  <p className="text-xs text-slate-400">
                    Recording vitals for {activeRequest.patientName || activeRequest.requestNumber}
                  </p>
                </div>
              </div>
              <button
                onClick={() => setIsTriageOpen(false)}
                className="p-2 rounded-lg text-slate-400 hover:text-white bg-slate-800"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <form onSubmit={handleSaveTriage} className="p-6 space-y-5">
              {/* Live Calculated Acuity Banner */}
              <div className={`p-4 rounded-2xl border flex items-center justify-between ${liveAcuity.color}`}>
                <div className="flex items-center gap-3">
                  <Flame className="w-5 h-5" />
                  <div>
                    <div className="text-xs font-bold uppercase tracking-wider">Computed Acuity Level</div>
                    <div className="text-sm font-semibold">{liveAcuity.label}</div>
                  </div>
                </div>
                <span className="text-xs font-mono font-bold px-3 py-1 rounded-full bg-slate-950/80 border border-current">
                  {liveAcuity.level}
                </span>
              </div>

              {/* Vitals Grid */}
              <div className="grid grid-cols-2 sm:grid-cols-3 gap-4">
                <div>
                  <label className="block text-xs font-medium text-slate-300 mb-1">
                    Systolic BP (mmHg)
                  </label>
                  <input
                    type="number"
                    value={systolicBp || ""}
                    onChange={(e) => setSystolicBp(e.target.value ? parseInt(e.target.value) : undefined)}
                    placeholder="120"
                    className="w-full px-3 py-2 bg-slate-950 border border-slate-800 rounded-xl text-sm text-white focus:outline-none focus:border-cyan-500"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-slate-300 mb-1">
                    Diastolic BP (mmHg)
                  </label>
                  <input
                    type="number"
                    value={diastolicBp || ""}
                    onChange={(e) => setDiastolicBp(e.target.value ? parseInt(e.target.value) : undefined)}
                    placeholder="80"
                    className="w-full px-3 py-2 bg-slate-950 border border-slate-800 rounded-xl text-sm text-white focus:outline-none focus:border-cyan-500"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-slate-300 mb-1">
                    Pulse Rate (bpm)
                  </label>
                  <input
                    type="number"
                    value={pulseRate || ""}
                    onChange={(e) => setPulseRate(e.target.value ? parseInt(e.target.value) : undefined)}
                    placeholder="72"
                    className="w-full px-3 py-2 bg-slate-950 border border-slate-800 rounded-xl text-sm text-white focus:outline-none focus:border-cyan-500"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-slate-300 mb-1">
                    Temperature (°C)
                  </label>
                  <input
                    type="number"
                    step="0.1"
                    value={temperature || ""}
                    onChange={(e) => setTemperature(e.target.value ? parseFloat(e.target.value) : undefined)}
                    placeholder="37.0"
                    className="w-full px-3 py-2 bg-slate-950 border border-slate-800 rounded-xl text-sm text-white focus:outline-none focus:border-cyan-500"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-slate-300 mb-1">
                    SpO2 Oxygen (%)
                  </label>
                  <input
                    type="number"
                    value={spo2 || ""}
                    onChange={(e) => setSpo2(e.target.value ? parseInt(e.target.value) : undefined)}
                    placeholder="98"
                    className="w-full px-3 py-2 bg-slate-950 border border-slate-800 rounded-xl text-sm text-white focus:outline-none focus:border-cyan-500"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-slate-300 mb-1">
                    Pain Score (0–10)
                  </label>
                  <input
                    type="number"
                    min={0}
                    max={10}
                    value={painScore || ""}
                    onChange={(e) => setPainScore(e.target.value ? parseInt(e.target.value) : undefined)}
                    placeholder="0"
                    className="w-full px-3 py-2 bg-slate-950 border border-slate-800 rounded-xl text-sm text-white focus:outline-none focus:border-cyan-500"
                  />
                </div>
              </div>

              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">
                  Triage Clinical Notes & Observations
                </label>
                <textarea
                  rows={2}
                  value={triageNotes}
                  onChange={(e) => setTriageNotes(e.target.value)}
                  placeholder="Patient alert, conscious, ambulatory, reported acute pain in left quadrant..."
                  className="w-full p-3 bg-slate-950 border border-slate-800 rounded-xl text-sm text-white focus:outline-none focus:border-cyan-500"
                />
              </div>

              <div className="pt-3 border-t border-slate-800 flex items-center justify-between">
                <button
                  type="button"
                  onClick={() => setIsTriageOpen(false)}
                  className="px-4 py-2 text-xs font-medium text-slate-400 hover:text-white bg-slate-800 rounded-xl"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={submitTriage.isPending}
                  className="flex items-center gap-2 px-6 py-2.5 rounded-xl bg-gradient-to-r from-cyan-500 to-blue-600 hover:from-cyan-400 text-white text-xs font-semibold shadow-lg shadow-cyan-500/20"
                >
                  {submitTriage.isPending ? (
                    <>
                      <Loader2 className="w-4 h-4 animate-spin" />
                      <span>Saving Vitals...</span>
                    </>
                  ) : (
                    <>
                      <CheckCircle2 className="w-4 h-4" />
                      <span>Confirm Triage & Match Doctor</span>
                    </>
                  )}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* 2. PROVIDER MATCHING & ALLOCATION DRAWER */}
      {isMatchOpen && activeRequest && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/70 backdrop-blur-sm animate-in fade-in duration-200">
          <div className="relative w-full max-w-2xl bg-slate-900 border border-slate-800 rounded-3xl shadow-2xl text-slate-100 overflow-hidden">
            <div className="flex items-center justify-between px-6 py-4 bg-slate-900/90 border-b border-slate-800">
              <div className="flex items-center gap-3">
                <div className="p-2.5 rounded-xl bg-cyan-500/10 text-cyan-400 border border-cyan-500/20">
                  <Stethoscope className="w-5 h-5" />
                </div>
                <div>
                  <h3 className="text-base font-bold text-white">Provider Matching & Doctor Allocation</h3>
                  <p className="text-xs text-slate-400">
                    Scored candidates for {activeRequest.serviceType} via {activeRequest.preferredMode}
                  </p>
                </div>
              </div>
              <button
                onClick={() => setIsMatchOpen(false)}
                className="p-2 rounded-lg text-slate-400 hover:text-white bg-slate-800"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <div className="p-6 space-y-4 max-h-[70vh] overflow-y-auto">
              {isMatchingLoading ? (
                <div className="p-12 text-center text-xs text-slate-500 animate-pulse">
                  Computing provider scores across specialty, channel, and queue capacity...
                </div>
              ) : !matchingData?.candidates || matchingData.candidates.length === 0 ? (
                <div className="p-8 text-center space-y-2">
                  <div className="text-sm font-medium text-white">No active providers found</div>
                  <p className="text-xs text-slate-400">
                    No doctor profiles found with specialty {activeRequest.serviceType}.
                  </p>
                </div>
              ) : (
                <div className="space-y-3">
                  {matchingData.candidates.map((cand: MatchedProviderCandidate) => (
                    <div
                      key={cand.providerId}
                      className="p-4 rounded-2xl bg-slate-950 border border-slate-800 hover:border-cyan-500/40 transition-all flex flex-col sm:flex-row sm:items-center justify-between gap-4"
                    >
                      <div className="space-y-1">
                        <div className="flex items-center gap-2">
                          <span className="font-semibold text-white text-sm">{cand.providerName}</span>
                          <span className="px-2 py-0.5 rounded-md bg-emerald-500/10 text-emerald-400 text-[10px] font-semibold border border-emerald-500/20">
                            {cand.status}
                          </span>
                          <span className="px-2 py-0.5 rounded-md bg-cyan-500/10 text-cyan-300 font-mono text-[10px] font-bold border border-cyan-500/30">
                            {cand.matchScore}% Match
                          </span>
                        </div>
                        <div className="text-xs text-slate-400">
                          Specialty: <strong className="text-slate-200">{cand.specialtyCode}</strong> • Active Queue:{" "}
                          <strong className="text-amber-400">{cand.currentActiveQueue}</strong> / {cand.maxActiveQueue}
                        </div>
                        <div className="text-[11px] text-slate-500">
                          {cand.matchingReasons.join(" • ")}
                        </div>
                      </div>

                      <button
                        type="button"
                        onClick={() => handleAssignDoctor(cand)}
                        disabled={assignProvider.isPending}
                        className="flex items-center gap-1.5 px-4 py-2 rounded-xl bg-gradient-to-r from-cyan-500 to-blue-600 hover:from-cyan-400 text-white text-xs font-semibold shadow-lg shadow-cyan-500/20 transition-all self-end sm:self-center"
                      >
                        <CheckCircle2 className="w-4 h-4" />
                        <span>Assign Doctor</span>
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
