import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  Stethoscope,
  Video,
  Phone,
  Building2,
  MessageSquare,
  Sparkles,
  Calendar,
  Clock,
  CheckCircle2,
  ArrowRight,
  ArrowLeft,
  Loader2,
  Tag,
} from "lucide-react";
import { useCreateCareRequest } from "@curexal/patient";
import { toast } from "sonner";

export function NewConsultationBookingPage() {
  const navigate = useNavigate();
  const createCareRequest = useCreateCareRequest();

  // Wizard state
  const [step, setStep] = useState<1 | 2 | 3>(1);
  const [serviceType, setServiceType] = useState("GENERAL_CONSULTATION");
  const [preferredMode, setPreferredMode] = useState<"VIDEO" | "AUDIO" | "IN_PERSON" | "ASYNC_CHAT">("VIDEO");
  const [urgency, setUrgency] = useState<"ROUTINE" | "URGENT" | "EMERGENCY">("ROUTINE");
  const [chiefComplaint, setChiefComplaint] = useState("");
  const [selectedSymptoms, setSelectedSymptoms] = useState<string[]>([]);
  const [preferredDate, setPreferredDate] = useState("");
  const [preferredTime, setPreferredTime] = useState("MORNING");

  const commonSymptoms = [
    "Fever",
    "Headache",
    "Dry Cough",
    "Sore Throat",
    "Chest Pain",
    "Fatigue & Weakness",
    "Nausea / Vomiting",
    "Abdominal Pain",
    "Joint Pain",
    "Skin Rash",
    "Dizziness",
    "Shortness of Breath",
  ];

  const toggleSymptom = (sym: string) => {
    setSelectedSymptoms((prev) =>
      prev.includes(sym) ? prev.filter((s) => s !== sym) : [...prev, sym]
    );
  };

  const handleBook = () => {
    if (!chiefComplaint) return;

    // Resolve patient identity from session/token
    const storedPatient = localStorage.getItem("curexal_portal_patient");
    let patientId = "";
    if (storedPatient) {
      try {
        const parsed = JSON.parse(storedPatient);
        if (parsed.id) patientId = parsed.id;
      } catch (e) {}
    }

    if (!patientId) {
      toast.error("Please sign in to complete your consultation booking.");
      navigate("/login");
      return;
    }

    createCareRequest.mutate(
      {
        patientId,
        serviceType,
        preferredMode,
        urgency,
        chiefComplaint,
        symptoms: selectedSymptoms,
        preferredTimeWindow: {
          date: preferredDate || "IMMEDIATE",
          timeSlot: preferredTime,
        },
      },
      {
        onSuccess: () => {
          toast.success("Care Request submitted successfully!");
          navigate("/dashboard");
        },
        onError: (err: any) => {
          toast.error(err.message || "Failed to submit care request.");
        },
      }
    );
  };

  return (
    <div className="max-w-4xl mx-auto space-y-6 animate-in fade-in duration-300">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900 dark:text-white flex items-center gap-2.5">
            <Stethoscope className="w-6 h-6 text-teal-600 dark:text-teal-400" />
            <span>Book Clinical Consultation</span>
          </h1>
          <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">
            Initiate a Care Request for in-person or telehealth medical attention
          </p>
        </div>

        {/* Step Indicator */}
        <div className="flex items-center gap-2">
          {[1, 2, 3].map((s) => (
            <div
              key={s}
              className={`w-8 h-8 rounded-xl flex items-center justify-center font-bold text-xs transition-all ${
                step === s
                  ? "bg-teal-600 text-white shadow-lg shadow-teal-600/30 ring-2 ring-teal-400/30"
                  : step > s
                  ? "bg-emerald-500/20 text-emerald-600 dark:text-emerald-400 border border-emerald-500/30"
                  : "bg-slate-100 dark:bg-slate-800 text-slate-400 dark:text-slate-500 border border-slate-200 dark:border-slate-700"
              }`}
            >
              {step > s ? <CheckCircle2 className="w-4 h-4" /> : s}
            </div>
          ))}
        </div>
      </div>

      {/* Wizard Content Cards */}
      <div className="p-6 sm:p-8 rounded-3xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-2xl space-y-6">
        {/* STEP 1 */}
        {step === 1 && (
          <div className="space-y-6 animate-in slide-in-from-right-4 duration-200">
            <div>
              <h2 className="text-base font-semibold text-slate-900 dark:text-white">Select Consultation Service</h2>
              <p className="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
                Choose the clinical discipline for your care demand
              </p>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              {[
                { id: "GENERAL_CONSULTATION", title: "General Outpatient Practice", desc: "Primary care, wellness exams, common illness intake", icon: Stethoscope },
                { id: "SPECIALIST", title: "Specialist Physician", desc: "Cardiology, Pediatrics, Gynecology, Dermatology", icon: Sparkles },
                { id: "LAB_TEST", title: "Diagnostic Lab Order", desc: "Blood panels, pathology, toxicology investigations", icon: Calendar },
                { id: "REFILL", title: "Medication Refill", desc: "Physician prescription review & pharmacy dispatch", icon: Clock },
              ].map((serv) => {
                const Icon = serv.icon;
                const isSelected = serviceType === serv.id;
                return (
                  <button
                    key={serv.id}
                    type="button"
                    onClick={() => setServiceType(serv.id)}
                    className={`p-4 rounded-2xl border text-left transition-all flex items-start gap-4 cursor-pointer ${
                      isSelected
                        ? "bg-teal-500/10 border-teal-500 ring-1 ring-teal-500 shadow-lg shadow-teal-500/10"
                        : "bg-slate-50 dark:bg-slate-950/60 border-slate-200 dark:border-slate-800 hover:border-slate-300 dark:hover:border-slate-700"
                    }`}
                  >
                    <div className={`p-2.5 rounded-xl ${isSelected ? "bg-teal-500/20 text-teal-600 dark:text-teal-400" : "bg-slate-200 dark:bg-slate-800 text-slate-500"}`}>
                      <Icon className="w-5 h-5" />
                    </div>
                    <div>
                      <div className="font-semibold text-sm text-slate-900 dark:text-white">{serv.title}</div>
                      <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">{serv.desc}</p>
                    </div>
                  </button>
                );
              })}
            </div>

            <div className="pt-4 border-t border-slate-200 dark:border-slate-800 space-y-3">
              <div>
                <h3 className="text-sm font-semibold text-slate-900 dark:text-white">Select Consultation Channel</h3>
                <p className="text-xs text-slate-500 dark:text-slate-400">Choose how you prefer to communicate with the doctor</p>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-4 gap-3">
                {[
                  { id: "VIDEO", title: "Live Video", desc: "HD WebRTC Call", icon: Video },
                  { id: "AUDIO", title: "Phone Audio", desc: "Voice Call", icon: Phone },
                  { id: "IN_PERSON", title: "In-Person", desc: "Clinic Visit", icon: Building2 },
                  { id: "ASYNC_CHAT", title: "Secure Chat", desc: "Text & Media", icon: MessageSquare },
                ].map((mode) => {
                  const Icon = mode.icon;
                  const isSelected = preferredMode === mode.id;
                  return (
                    <button
                      key={mode.id}
                      type="button"
                      onClick={() => setPreferredMode(mode.id as any)}
                      className={`p-3.5 rounded-xl border text-center transition-all flex flex-col items-center gap-2 cursor-pointer ${
                        isSelected
                          ? "bg-teal-500/20 border-teal-500 text-teal-700 dark:text-teal-300 ring-1 ring-teal-500 shadow-md shadow-teal-500/10"
                          : "bg-slate-50 dark:bg-slate-950/60 border-slate-200 dark:border-slate-800 text-slate-500 hover:text-slate-900 dark:hover:text-white"
                      }`}
                    >
                      <Icon className="w-5 h-5" />
                      <div>
                        <div className="font-semibold text-xs text-slate-900 dark:text-white">{mode.title}</div>
                        <div className="text-[10px] text-slate-400">{mode.desc}</div>
                      </div>
                    </button>
                  );
                })}
              </div>
            </div>

            <div className="flex justify-end pt-4">
              <button
                type="button"
                onClick={() => setStep(2)}
                className="flex items-center gap-2 px-6 py-2.5 rounded-xl bg-[#0F766E] hover:bg-[#115E59] text-white text-xs font-semibold shadow-lg shadow-teal-600/20 transition-all cursor-pointer"
              >
                <span>Continue to Symptoms</span>
                <ArrowRight className="w-4 h-4" />
              </button>
            </div>
          </div>
        )}

        {/* STEP 2 */}
        {step === 2 && (
          <div className="space-y-6 animate-in slide-in-from-right-4 duration-200">
            <div>
              <h2 className="text-base font-semibold text-slate-900 dark:text-white">Symptoms & Chief Complaint</h2>
              <p className="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
                Describe your condition so the matching engine routes to the most qualified provider
              </p>
            </div>

            <div>
              <label className="block text-xs font-medium text-slate-700 dark:text-slate-300 mb-1.5">
                Describe Your Primary Medical Reason / Chief Complaint <span className="text-rose-500">*</span>
              </label>
              <textarea
                required
                rows={3}
                value={chiefComplaint}
                onChange={(e) => setChiefComplaint(e.target.value)}
                placeholder="e.g. Experiencing persistent fever, dull headache, and mild throat irritation for the past 2 days..."
                className="w-full p-3 bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl text-sm text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:border-teal-500 focus:ring-1 focus:ring-teal-500"
              />
            </div>

            <div>
              <label className="block text-xs font-medium text-slate-700 dark:text-slate-300 mb-2">
                Select Common Symptoms (Click all that apply)
              </label>
              <div className="flex flex-wrap gap-2">
                {commonSymptoms.map((sym) => {
                  const isSelected = selectedSymptoms.includes(sym);
                  return (
                    <button
                      key={sym}
                      type="button"
                      onClick={() => toggleSymptom(sym)}
                      className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium border transition-all cursor-pointer ${
                        isSelected
                          ? "bg-teal-500/20 border-teal-500/50 text-teal-700 dark:text-teal-300 shadow-sm"
                          : "bg-slate-50 dark:bg-slate-950 border-slate-200 dark:border-slate-800 text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200"
                      }`}
                    >
                      <Tag className="w-3 h-3" />
                      <span>{sym}</span>
                    </button>
                  );
                })}
              </div>
            </div>

            <div className="pt-4 border-t border-slate-200 dark:border-slate-800 flex items-center justify-between">
              <button
                type="button"
                onClick={() => setStep(1)}
                className="flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-medium text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white bg-slate-100 dark:bg-slate-800 transition-colors cursor-pointer"
              >
                <ArrowLeft className="w-4 h-4" />
                <span>Back</span>
              </button>
              <button
                type="button"
                disabled={!chiefComplaint.trim()}
                onClick={() => setStep(3)}
                className="flex items-center gap-2 px-6 py-2.5 rounded-xl bg-[#0F766E] hover:bg-[#115E59] disabled:opacity-50 text-white text-xs font-semibold shadow-lg shadow-teal-600/20 transition-all cursor-pointer"
              >
                <span>Continue to Urgency & Time</span>
                <ArrowRight className="w-4 h-4" />
              </button>
            </div>
          </div>
        )}

        {/* STEP 3 */}
        {step === 3 && (
          <div className="space-y-6 animate-in slide-in-from-right-4 duration-200">
            <div>
              <h2 className="text-base font-semibold text-slate-900 dark:text-white">Urgency & Preferred Schedule</h2>
              <p className="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
                Set clinical triage priority and appointment availability
              </p>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              {[
                { id: "ROUTINE", title: "Routine (Standard)", desc: "Normal queuing and consultation scheduling" },
                { id: "URGENT", title: "Urgent (Priority)", desc: "Elevated acuity queue; priority doctor allocation" },
                { id: "EMERGENCY", title: "Emergency Alert", desc: "Immediate clinical intervention escalation" },
              ].map((urg) => {
                const isSelected = urgency === urg.id;
                return (
                  <button
                    key={urg.id}
                    type="button"
                    onClick={() => setUrgency(urg.id as any)}
                    className={`p-4 rounded-2xl border text-left transition-all cursor-pointer ${
                      isSelected
                        ? "bg-teal-500/10 border-teal-500 ring-1 ring-teal-500 shadow-lg shadow-teal-500/10"
                        : "bg-slate-50 dark:bg-slate-950/60 border-slate-200 dark:border-slate-800"
                    }`}
                  >
                    <div className="font-semibold text-sm text-slate-900 dark:text-white">{urg.title}</div>
                    <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">{urg.desc}</p>
                  </button>
                );
              })}
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 pt-2">
              <div>
                <label className="block text-xs font-medium text-slate-700 dark:text-slate-300 mb-1">
                  Preferred Date
                </label>
                <input
                  type="date"
                  value={preferredDate}
                  onChange={(e) => setPreferredDate(e.target.value)}
                  className="w-full px-3 py-2 bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl text-sm text-slate-900 dark:text-white focus:outline-none focus:border-teal-500"
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-slate-700 dark:text-slate-300 mb-1">
                  Time Slot Window
                </label>
                <select
                  value={preferredTime}
                  onChange={(e) => setPreferredTime(e.target.value)}
                  className="w-full px-3 py-2 bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl text-sm text-slate-900 dark:text-white focus:outline-none focus:border-teal-500"
                >
                  <option value="IMMEDIATE">Immediate Next Available Doctor</option>
                  <option value="MORNING">Morning (08:00 AM – 12:00 PM)</option>
                  <option value="AFTERNOON">Afternoon (12:00 PM – 04:00 PM)</option>
                  <option value="EVENING">Evening (04:00 PM – 08:00 PM)</option>
                </select>
              </div>
            </div>

            {/* Summary Review */}
            <div className="p-4 rounded-2xl bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800/80 text-xs space-y-2">
              <div className="font-semibold text-slate-900 dark:text-white">Summary of Request</div>
              <div className="text-slate-600 dark:text-slate-400 space-y-1">
                <div>Service: <strong className="text-slate-900 dark:text-slate-200">{serviceType}</strong> via <strong className="text-teal-600 dark:text-teal-400">{preferredMode}</strong></div>
                <div>Reason: <span className="text-slate-800 dark:text-slate-300">{chiefComplaint}</span></div>
                {selectedSymptoms.length > 0 && (
                  <div>Symptoms: <span className="text-slate-800 dark:text-slate-300">{selectedSymptoms.join(", ")}</span></div>
                )}
              </div>
            </div>

            <div className="pt-4 border-t border-slate-200 dark:border-slate-800 flex items-center justify-between">
              <button
                type="button"
                onClick={() => setStep(2)}
                className="flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-medium text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white bg-slate-100 dark:bg-slate-800 transition-colors cursor-pointer"
              >
                <ArrowLeft className="w-4 h-4" />
                <span>Back</span>
              </button>
              <button
                type="button"
                disabled={createCareRequest.isPending}
                onClick={handleBook}
                className="flex items-center gap-2 px-6 py-2.5 rounded-xl bg-[#0F766E] hover:bg-[#115E59] text-white text-xs font-semibold shadow-lg shadow-teal-600/25 transition-all cursor-pointer"
              >
                {createCareRequest.isPending ? (
                  <>
                    <Loader2 className="w-4 h-4 animate-spin" />
                    <span>Submitting Care Request...</span>
                  </>
                ) : (
                  <>
                    <CheckCircle2 className="w-4 h-4" />
                    <span>Confirm & Launch Care Journey</span>
                  </>
                )}
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
