import React, { useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import {
  Stethoscope,
  Video,
  VideoOff,
  Mic,
  MicOff,
  ScreenShare,
  PhoneOff,
  Sparkles,
  FileText,
  Microscope,
  Layers,
  Pill,
  CheckCircle2,
  Clock,
  AlertCircle,
  Save,
  Send,
  Plus,
  Trash2,
  ChevronRight,
  ShieldCheck,
  Activity,
  HeartPulse,
  Flame,
  ArrowLeft,
  Loader2,
} from "lucide-react";
import {
  useEncounter,
  useSaveSOAPNotes,
  useDispatchOrders,
  useCompleteEncounter,
} from "../../../api/hooks/use-encounters";
import type {
  LabOrderItem,
  RadiologyOrderItem,
  PrescriptionOrderItem,
} from "../../../api/services/encounter.service";

export const ClinicalEncounterRoomPage: React.FC = () => {
  const { branchSlug, encounterId } = useParams<{ branchSlug: string; encounterId: string }>();
  const navigate = useNavigate();

  const { data: encounter, isLoading } = useEncounter(encounterId || "");
  const saveSOAPNotes = useSaveSOAPNotes();
  const dispatchOrders = useDispatchOrders();
  const completeEncounter = useCompleteEncounter();

  // Telehealth video states
  const [isVideoOn, setIsVideoOn] = useState(true);
  const [isMicOn, setIsMicOn] = useState(true);
  const [isScreenSharing, setIsScreenSharing] = useState(false);

  // Active documentation tab
  const [activeTab, setActiveTab] = useState<"SOAP" | "LABS" | "RADIOLOGY" | "PHARMACY">("SOAP");

  // SOAP State
  const [subjective, setSubjective] = useState(
    "Patient presents with a 3-day history of low-grade fever, dry irritating cough, and generalized fatigue. No reported chest pain or shortness of breath."
  );
  const [objective, setObjective] = useState(
    "General: Alert, oriented, mild pallor. Chest: Vesicular breath sounds bilaterally, no crepitations or wheezes. CVS: S1 S2 present, no murmurs. Abdomen: Soft, non-tender, no organomegaly."
  );
  const [assessment, setAssessment] = useState(
    "Acute viral upper respiratory tract infection with mild systemic symptoms."
  );
  const [plan, setPlan] = useState(
    "1. Symptomatic relief with antipyretics and hydration.\n2. Diagnostic blood panel (FBC + Malaria).\n3. Rest for 3 days and review if symptoms persist."
  );
  const [primaryDiagCode, setPrimaryDiagCode] = useState("J06.9");
  const [primaryDiagName, setPrimaryDiagName] = useState("Acute Upper Respiratory Infection");

  // Clinical Orders State
  const [labOrders, setLabOrders] = useState<LabOrderItem[]>([
    { testCode: "LAB_FBC", testName: "Full Blood Count (FBC)", urgency: "ROUTINE", notes: "Check for leukocytosis" },
    { testCode: "LAB_MP", testName: "Malaria Parasite Microscopy", urgency: "ROUTINE", notes: "Rule out Plasmodium falciparum" },
  ]);

  const [radiologyOrders, setRadiologyOrders] = useState<RadiologyOrderItem[]>([]);

  const [prescriptions, setPrescriptions] = useState<PrescriptionOrderItem[]>([
    { medicationName: "Paracetamol Tablet", dosage: "1000mg", frequency: "TDS (8-Hourly)", durationDays: 3, instructions: "Take after food" },
    { medicationName: "Vitamin C Chewable", dosage: "500mg", frequency: "Daily", durationDays: 7, instructions: "Morning daily" },
  ]);

  // Order item inputs
  const [newLabName, setNewLabName] = useState("");
  const [newRadModality, setNewRadModality] = useState("XRAY");
  const [newRadStudy, setNewRadStudy] = useState("");
  const [newDrugName, setNewDrugName] = useState("");
  const [newDrugDosage, setNewDrugDosage] = useState("500mg");
  const [newDrugFreq, setNewDrugFreq] = useState("BD (12-Hourly)");
  const [newDrugDays, setNewDrugDays] = useState(5);

  const commonDiagnoses = [
    { code: "J06.9", name: "Acute Upper Respiratory Infection" },
    { code: "B50.9", name: "Plasmodium Falciparum Malaria" },
    { code: "I10", name: "Essential (Primary) Hypertension" },
    { code: "E11.9", name: "Type 2 Diabetes Mellitus" },
    { code: "A09", name: "Infectious Gastroenteritis and Colitis" },
    { code: "K29.7", name: "Gastritis, Unspecified" },
  ];

  const handleSaveSOAP = () => {
    if (!encounterId) return;
    saveSOAPNotes.mutate({
      encounterId,
      payload: {
        subjective,
        objective,
        assessment,
        plan,
        primaryDiagnosisCode: primaryDiagCode,
        primaryDiagnosisName: primaryDiagName,
      },
    });
  };

  const handleDispatchAllOrders = () => {
    if (!encounterId) return;
    dispatchOrders.mutate({
      encounterId,
      payload: {
        labOrders,
        radiologyOrders,
        prescriptionList: prescriptions,
      },
    });
  };

  const handleFinalizeEncounter = () => {
    if (!encounterId) return;
    completeEncounter.mutate(encounterId, {
      onSuccess: () => {
        navigate(`/${branchSlug}/clinical`);
      },
    });
  };

  const addLabOrder = () => {
    if (!newLabName.trim()) return;
    setLabOrders((prev) => [
      ...prev,
      {
        testCode: "LAB_CUSTOM_" + Date.now(),
        testName: newLabName.trim(),
        urgency: "ROUTINE",
        notes: "Clinical doctor order",
      },
    ]);
    setNewLabName("");
  };

  const addRadiologyOrder = () => {
    if (!newRadStudy.trim()) return;
    setRadiologyOrders((prev) => [
      ...prev,
      {
        modality: newRadModality as any,
        studyName: newRadStudy.trim(),
        urgency: "ROUTINE",
      },
    ]);
    setNewRadStudy("");
  };

  const addPrescription = () => {
    if (!newDrugName.trim()) return;
    setPrescriptions((prev) => [
      ...prev,
      {
        medicationName: newDrugName.trim(),
        dosage: newDrugDosage,
        frequency: newDrugFreq,
        durationDays: newDrugDays,
      },
    ]);
    setNewDrugName("");
  };

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans selection:bg-cyan-500/30">
      {/* Top Clinical Header Bar */}
      <header className="sticky top-0 z-40 w-full bg-slate-900/90 backdrop-blur-md border-b border-slate-800 px-4 py-3 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Link
            to={`/${branchSlug}/clinical`}
            className="p-2 rounded-xl text-slate-400 hover:text-white bg-slate-800 transition-colors"
          >
            <ArrowLeft className="w-4 h-4" />
          </Link>

          <div className="flex items-center gap-3">
            <div className="p-2.5 rounded-xl bg-cyan-500/10 text-cyan-400 border border-cyan-500/20">
              <Stethoscope className="w-5 h-5" />
            </div>
            <div>
              <div className="flex items-center gap-2">
                <span className="font-bold text-white text-base">
                  {encounter?.patientName || "Clinical Consultation"}
                </span>
                <span className="px-2 py-0.5 rounded-md bg-slate-800 text-cyan-400 font-mono text-xs font-semibold border border-slate-700">
                  {encounter?.mrn || "PAT-2026-00412"}
                </span>
                <span className="px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-400 text-[10px] font-semibold border border-emerald-500/20">
                  IN PROGRESS
                </span>
              </div>
              <p className="text-[11px] text-slate-400">
                Mode: <strong className="text-cyan-400">{encounter?.mode || "VIDEO (WebRTC Telehealth)"}</strong> • Type: {encounter?.encounterType || "OUTPATIENT"}
              </p>
            </div>
          </div>
        </div>

        {/* Action Controls */}
        <div className="flex items-center gap-3">
          <button
            type="button"
            onClick={handleSaveSOAP}
            disabled={saveSOAPNotes.isPending}
            className="flex items-center gap-1.5 px-3.5 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-medium border border-slate-700 transition-all"
          >
            <Save className="w-3.5 h-3.5 text-cyan-400" />
            <span>{saveSOAPNotes.isPending ? "Saving..." : "Save Notes"}</span>
          </button>

          <button
            type="button"
            onClick={handleFinalizeEncounter}
            disabled={completeEncounter.isPending}
            className="flex items-center gap-1.5 px-4 py-2 rounded-xl bg-gradient-to-r from-emerald-500 to-teal-600 hover:from-emerald-400 text-white text-xs font-semibold shadow-lg shadow-emerald-500/20 transition-all"
          >
            {completeEncounter.isPending ? (
              <Loader2 className="w-3.5 h-3.5 animate-spin" />
            ) : (
              <CheckCircle2 className="w-3.5 h-3.5" />
            )}
            <span>Conclude Consultation</span>
          </button>
        </div>
      </header>

      {/* Main Split Grid */}
      <div className="flex-1 grid grid-cols-1 lg:grid-cols-12 gap-4 p-4 max-w-7xl w-full mx-auto">
        {/* LEFT: Telehealth Stage & Patient Vitals (5 cols) */}
        <div className="lg:col-span-5 space-y-4 flex flex-col">
          {/* WebRTC Video Room */}
          <div className="relative aspect-video rounded-3xl bg-slate-900 border border-slate-800 overflow-hidden shadow-2xl flex flex-col justify-between p-4">
            {/* Background Stream Simulation */}
            <div className="absolute inset-0 bg-gradient-to-br from-slate-950 via-slate-900 to-indigo-950 flex items-center justify-center">
              {isVideoOn ? (
                <div className="text-center space-y-2">
                  <div className="w-16 h-16 rounded-full bg-cyan-500/20 border border-cyan-500/40 text-cyan-300 flex items-center justify-center mx-auto animate-pulse">
                    <Video className="w-8 h-8" />
                  </div>
                  <div className="text-xs font-medium text-white">Live Telehealth Encrypted Feed</div>
                  <div className="text-[10px] text-cyan-400 font-mono">WebRTC 1080p HD • Latency 24ms</div>
                </div>
              ) : (
                <div className="text-center space-y-1 text-slate-500 text-xs">
                  <VideoOff className="w-8 h-8 mx-auto text-slate-600" />
                  <div>Camera Muted</div>
                </div>
              )}
            </div>

            {/* Top Video Status */}
            <div className="relative z-10 flex items-center justify-between">
              <span className="px-2.5 py-1 rounded-full bg-slate-900/80 backdrop-blur-md border border-slate-700 text-emerald-400 text-[10px] font-mono flex items-center gap-1.5">
                <span className="w-2 h-2 rounded-full bg-emerald-400 animate-ping" />
                Live Session (00:14:32)
              </span>
            </div>

            {/* Self Video PIP (Picture in Picture) */}
            <div className="absolute right-4 bottom-16 w-24 h-16 rounded-xl bg-slate-800/90 border border-slate-700 shadow-lg flex items-center justify-center text-[10px] text-slate-400 font-mono">
              Dr. (You)
            </div>

            {/* Bottom Stream Controls */}
            <div className="relative z-10 flex items-center justify-center gap-3 py-1">
              <button
                type="button"
                onClick={() => setIsMicOn(!isMicOn)}
                className={`p-3 rounded-2xl transition-all shadow-md ${
                  isMicOn ? "bg-slate-800 text-white hover:bg-slate-700" : "bg-rose-500 text-white"
                }`}
                title={isMicOn ? "Mute Microphone" : "Unmute Microphone"}
              >
                {isMicOn ? <Mic className="w-4 h-4" /> : <MicOff className="w-4 h-4" />}
              </button>

              <button
                type="button"
                onClick={() => setIsVideoOn(!isVideoOn)}
                className={`p-3 rounded-2xl transition-all shadow-md ${
                  isVideoOn ? "bg-slate-800 text-white hover:bg-slate-700" : "bg-rose-500 text-white"
                }`}
                title={isVideoOn ? "Turn Camera Off" : "Turn Camera On"}
              >
                {isVideoOn ? <Video className="w-4 h-4" /> : <VideoOff className="w-4 h-4" />}
              </button>

              <button
                type="button"
                onClick={() => setIsScreenSharing(!isScreenSharing)}
                className={`p-3 rounded-2xl transition-all shadow-md ${
                  isScreenSharing ? "bg-cyan-500 text-white" : "bg-slate-800 text-white hover:bg-slate-700"
                }`}
                title="Share Screen"
              >
                <ScreenShare className="w-4 h-4" />
              </button>
            </div>
          </div>

          {/* Vitals Summary Card */}
          <div className="p-5 rounded-3xl bg-slate-900 border border-slate-800 shadow-xl space-y-3">
            <div className="flex items-center justify-between border-b border-slate-800 pb-2.5">
              <h3 className="text-xs font-bold text-white uppercase tracking-wider flex items-center gap-2">
                <HeartPulse className="w-4 h-4 text-rose-400" />
                <span>Recorded Triage Vitals</span>
              </h3>
              <span className="text-[10px] font-semibold text-emerald-400 bg-emerald-500/10 px-2 py-0.5 rounded border border-emerald-500/20">
                Acuity: GREEN
              </span>
            </div>

            <div className="grid grid-cols-3 gap-2 text-center">
              <div className="p-2.5 rounded-xl bg-slate-950 border border-slate-800">
                <div className="text-[10px] text-slate-400">Blood Pressure</div>
                <div className="text-xs font-bold text-white mt-0.5">120 / 80</div>
              </div>
              <div className="p-2.5 rounded-xl bg-slate-950 border border-slate-800">
                <div className="text-[10px] text-slate-400">Pulse Rate</div>
                <div className="text-xs font-bold text-white mt-0.5">72 bpm</div>
              </div>
              <div className="p-2.5 rounded-xl bg-slate-950 border border-slate-800">
                <div className="text-[10px] text-slate-400">Temperature</div>
                <div className="text-xs font-bold text-white mt-0.5">36.8 °C</div>
              </div>
              <div className="p-2.5 rounded-xl bg-slate-950 border border-slate-800">
                <div className="text-[10px] text-slate-400">Oxygen SpO2</div>
                <div className="text-xs font-bold text-white mt-0.5">98%</div>
              </div>
              <div className="p-2.5 rounded-xl bg-slate-950 border border-slate-800">
                <div className="text-[10px] text-slate-400">Resp Rate</div>
                <div className="text-xs font-bold text-white mt-0.5">18 bpm</div>
              </div>
              <div className="p-2.5 rounded-xl bg-slate-950 border border-slate-800">
                <div className="text-[10px] text-slate-400">Pain Score</div>
                <div className="text-xs font-bold text-emerald-400 mt-0.5">0 / 10</div>
              </div>
            </div>
          </div>
        </div>

        {/* RIGHT: SOAP Documentation & Order Dispatch Tabs (7 cols) */}
        <div className="lg:col-span-7 flex flex-col space-y-4">
          {/* Sub Navigation Tabs */}
          <div className="flex items-center gap-1.5 p-1.5 rounded-2xl bg-slate-900 border border-slate-800">
            {[
              { id: "SOAP", label: "SOAP Notes", icon: FileText },
              { id: "LABS", label: `Labs (${labOrders.length})`, icon: Microscope },
              { id: "RADIOLOGY", label: `Radiology (${radiologyOrders.length})`, icon: Layers },
              { id: "PHARMACY", label: `Pharmacy (${prescriptions.length})`, icon: Pill },
            ].map((tab) => {
              const Icon = tab.icon;
              const isActive = activeTab === tab.id;
              return (
                <button
                  key={tab.id}
                  type="button"
                  onClick={() => setActiveTab(tab.id as any)}
                  className={`flex-1 flex items-center justify-center gap-1.5 py-2 rounded-xl text-xs font-medium transition-all ${
                    isActive
                      ? "bg-cyan-500/20 text-cyan-300 border border-cyan-500/40 shadow-sm"
                      : "text-slate-400 hover:text-slate-200 hover:bg-slate-800"
                  }`}
                >
                  <Icon className="w-3.5 h-3.5" />
                  <span>{tab.label}</span>
                </button>
              );
            })}
          </div>

          {/* TAB 1: SOAP DOCUMENTATION */}
          {activeTab === "SOAP" && (
            <div className="p-6 rounded-3xl bg-slate-900 border border-slate-800 shadow-xl space-y-4 flex-1">
              <div>
                <label className="block text-xs font-semibold text-cyan-400 uppercase tracking-wider mb-1">
                  S — Subjective History of Presenting Illness
                </label>
                <textarea
                  rows={2}
                  value={subjective}
                  onChange={(e) => setSubjective(e.target.value)}
                  className="w-full p-3 bg-slate-950 border border-slate-800 rounded-xl text-xs text-white focus:outline-none focus:border-cyan-500 transition-colors"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-cyan-400 uppercase tracking-wider mb-1">
                  O — Objective Examination & Clinical Findings
                </label>
                <textarea
                  rows={2}
                  value={objective}
                  onChange={(e) => setObjective(e.target.value)}
                  className="w-full p-3 bg-slate-950 border border-slate-800 rounded-xl text-xs text-white focus:outline-none focus:border-cyan-500 transition-colors"
                />
              </div>

              {/* Assessment & ICD-10 Diagnosis Selector */}
              <div>
                <label className="block text-xs font-semibold text-cyan-400 uppercase tracking-wider mb-1">
                  A — Assessment & Primary ICD-10 Diagnosis
                </label>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 mb-2">
                  {commonDiagnoses.map((diag) => {
                    const isSelected = primaryDiagCode === diag.code;
                    return (
                      <button
                        key={diag.code}
                        type="button"
                        onClick={() => {
                          setPrimaryDiagCode(diag.code);
                          setPrimaryDiagName(diag.name);
                        }}
                        className={`p-2 rounded-xl text-left text-xs border transition-all ${
                          isSelected
                            ? "bg-cyan-500/20 border-cyan-500 text-cyan-200"
                            : "bg-slate-950 border-slate-800 text-slate-400 hover:border-slate-700"
                        }`}
                      >
                        <strong className="font-mono text-cyan-400">{diag.code}</strong> — {diag.name}
                      </button>
                    );
                  })}
                </div>
                <textarea
                  rows={2}
                  value={assessment}
                  onChange={(e) => setAssessment(e.target.value)}
                  className="w-full p-3 bg-slate-950 border border-slate-800 rounded-xl text-xs text-white focus:outline-none focus:border-cyan-500 transition-colors"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-cyan-400 uppercase tracking-wider mb-1">
                  P — Plan of Management & Treatment
                </label>
                <textarea
                  rows={3}
                  value={plan}
                  onChange={(e) => setPlan(e.target.value)}
                  className="w-full p-3 bg-slate-950 border border-slate-800 rounded-xl text-xs text-white focus:outline-none focus:border-cyan-500 transition-colors"
                />
              </div>
            </div>
          )}

          {/* TAB 2: LAB DIAGNOSTIC ORDERS */}
          {activeTab === "LABS" && (
            <div className="p-6 rounded-3xl bg-slate-900 border border-slate-800 shadow-xl space-y-4 flex-1">
              <div className="flex items-center justify-between border-b border-slate-800 pb-3">
                <div>
                  <h3 className="text-sm font-bold text-white">Diagnostic Laboratory Orders (LIS)</h3>
                  <p className="text-xs text-slate-400">Orders will dispatch directly to the Laboratory bench queue</p>
                </div>
                <button
                  type="button"
                  onClick={handleDispatchAllOrders}
                  className="flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-cyan-500/20 text-cyan-300 border border-cyan-500/40 text-xs font-semibold hover:bg-cyan-500/30 transition-all"
                >
                  <Send className="w-3 h-3" />
                  <span>Dispatch to LIS</span>
                </button>
              </div>

              {/* Add Lab Form */}
              <div className="flex gap-2">
                <input
                  type="text"
                  value={newLabName}
                  onChange={(e) => setNewLabName(e.target.value)}
                  placeholder="e.g. Fasting Lipid Profile, Liver Function Test (LFT)..."
                  className="flex-1 px-3 py-2 bg-slate-950 border border-slate-800 rounded-xl text-xs text-white focus:outline-none focus:border-cyan-500"
                />
                <button
                  type="button"
                  onClick={addLabOrder}
                  className="px-4 py-2 bg-cyan-500 hover:bg-cyan-400 text-white text-xs font-semibold rounded-xl"
                >
                  Add Test
                </button>
              </div>

              {/* Lab List */}
              <div className="space-y-2">
                {labOrders.map((lab, idx) => (
                  <div
                    key={idx}
                    className="p-3 rounded-xl bg-slate-950 border border-slate-800 flex items-center justify-between text-xs"
                  >
                    <div>
                      <span className="font-semibold text-white">{lab.testName}</span>
                      <span className="ml-2 text-[10px] text-slate-400">({lab.urgency})</span>
                      {lab.notes && <div className="text-[11px] text-slate-500 mt-0.5">{lab.notes}</div>}
                    </div>
                    <button
                      onClick={() => setLabOrders(labOrders.filter((_, i) => i !== idx))}
                      className="text-slate-500 hover:text-rose-400 p-1"
                    >
                      <Trash2 className="w-3.5 h-3.5" />
                    </button>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* TAB 3: RADIOLOGY ORDERS */}
          {activeTab === "RADIOLOGY" && (
            <div className="p-6 rounded-3xl bg-slate-900 border border-slate-800 shadow-xl space-y-4 flex-1">
              <div className="flex items-center justify-between border-b border-slate-800 pb-3">
                <div>
                  <h3 className="text-sm font-bold text-white">Radiology & Imaging Orders (PACS)</h3>
                  <p className="text-xs text-slate-400">Orders route to the Radiology department DICOM worklist</p>
                </div>
                <button
                  type="button"
                  onClick={handleDispatchAllOrders}
                  className="flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-cyan-500/20 text-cyan-300 border border-cyan-500/40 text-xs font-semibold hover:bg-cyan-500/30 transition-all"
                >
                  <Send className="w-3 h-3" />
                  <span>Dispatch to PACS</span>
                </button>
              </div>

              <div className="grid grid-cols-3 gap-2">
                <select
                  value={newRadModality}
                  onChange={(e) => setNewRadModality(e.target.value)}
                  className="px-3 py-2 bg-slate-950 border border-slate-800 rounded-xl text-xs text-white"
                >
                  <option value="XRAY">X-Ray</option>
                  <option value="ULTRASOUND">Ultrasound</option>
                  <option value="CT">CT Scan</option>
                  <option value="MRI">MRI</option>
                </select>
                <input
                  type="text"
                  value={newRadStudy}
                  onChange={(e) => setNewRadStudy(e.target.value)}
                  placeholder="e.g. Chest X-Ray (PA View)"
                  className="col-span-2 px-3 py-2 bg-slate-950 border border-slate-800 rounded-xl text-xs text-white"
                />
              </div>
              <button
                type="button"
                onClick={addRadiologyOrder}
                className="w-full py-2 bg-cyan-500 hover:bg-cyan-400 text-white text-xs font-semibold rounded-xl"
              >
                Add Imaging Study
              </button>

              <div className="space-y-2">
                {radiologyOrders.length === 0 ? (
                  <div className="p-6 text-center text-xs text-slate-500">No imaging studies added</div>
                ) : (
                  radiologyOrders.map((rad, idx) => (
                    <div
                      key={idx}
                      className="p-3 rounded-xl bg-slate-950 border border-slate-800 flex items-center justify-between text-xs"
                    >
                      <div>
                        <strong className="text-cyan-400">[{rad.modality}]</strong> {rad.studyName}
                      </div>
                      <button
                        onClick={() => setRadiologyOrders(radiologyOrders.filter((_, i) => i !== idx))}
                        className="text-slate-500 hover:text-rose-400 p-1"
                      >
                        <Trash2 className="w-3.5 h-3.5" />
                      </button>
                    </div>
                  ))
                )}
              </div>
            </div>
          )}

          {/* TAB 4: e-PRESCRIPTION PHARMACY */}
          {activeTab === "PHARMACY" && (
            <div className="p-6 rounded-3xl bg-slate-900 border border-slate-800 shadow-xl space-y-4 flex-1">
              <div className="flex items-center justify-between border-b border-slate-800 pb-3">
                <div>
                  <h3 className="text-sm font-bold text-white">e-Prescription & Pharmacy Dispense</h3>
                  <p className="text-xs text-slate-400">Medications route directly to the Pharmacy Dispensary</p>
                </div>
                <button
                  type="button"
                  onClick={handleDispatchAllOrders}
                  className="flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-cyan-500/20 text-cyan-300 border border-cyan-500/40 text-xs font-semibold hover:bg-cyan-500/30 transition-all"
                >
                  <Send className="w-3 h-3" />
                  <span>Dispatch to Pharmacy</span>
                </button>
              </div>

              {/* Add Drug Form */}
              <div className="grid grid-cols-1 sm:grid-cols-4 gap-2">
                <input
                  type="text"
                  value={newDrugName}
                  onChange={(e) => setNewDrugName(e.target.value)}
                  placeholder="Drug Name (e.g. Amoxicillin)"
                  className="sm:col-span-2 px-3 py-2 bg-slate-950 border border-slate-800 rounded-xl text-xs text-white"
                />
                <input
                  type="text"
                  value={newDrugDosage}
                  onChange={(e) => setNewDrugDosage(e.target.value)}
                  placeholder="Dosage (500mg)"
                  className="px-3 py-2 bg-slate-950 border border-slate-800 rounded-xl text-xs text-white"
                />
                <select
                  value={newDrugFreq}
                  onChange={(e) => setNewDrugFreq(e.target.value)}
                  className="px-3 py-2 bg-slate-950 border border-slate-800 rounded-xl text-xs text-white"
                >
                  <option value="OD (Once Daily)">OD (Once Daily)</option>
                  <option value="BD (12-Hourly)">BD (12-Hourly)</option>
                  <option value="TDS (8-Hourly)">TDS (8-Hourly)</option>
                  <option value="QDS (6-Hourly)">QDS (6-Hourly)</option>
                  <option value="PRN (As Needed)">PRN (As Needed)</option>
                </select>
              </div>

              <button
                type="button"
                onClick={addPrescription}
                className="w-full py-2 bg-cyan-500 hover:bg-cyan-400 text-white text-xs font-semibold rounded-xl"
              >
                Add Medication to Prescription
              </button>

              <div className="space-y-2">
                {prescriptions.map((rx, idx) => (
                  <div
                    key={idx}
                    className="p-3 rounded-xl bg-slate-950 border border-slate-800 flex items-center justify-between text-xs"
                  >
                    <div>
                      <span className="font-semibold text-white">{rx.medicationName}</span> —{" "}
                      <span className="text-cyan-400">{rx.dosage}</span> • {rx.frequency} for {rx.durationDays} days
                      {rx.instructions && (
                        <div className="text-[11px] text-slate-400 mt-0.5">{rx.instructions}</div>
                      )}
                    </div>
                    <button
                      onClick={() => setPrescriptions(prescriptions.filter((_, i) => i !== idx))}
                      className="text-slate-500 hover:text-rose-400 p-1"
                    >
                      <Trash2 className="w-3.5 h-3.5" />
                    </button>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
