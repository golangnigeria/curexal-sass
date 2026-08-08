import React, { useState } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  Sparkles,
  CheckCircle2,
  User,
  ShoppingBag,
  Stethoscope,
  ArrowRight,
  ArrowLeft,
  Check,
  MessageSquare,
  FlaskRound,
  Rocket,
  Mail,
} from "lucide-react";
import { toast } from "sonner";
import { getApiUrl } from "@/api";
import { saveWaitlistToSupabase } from "@/lib/supabase";

interface WaitlistModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export type PersonaCategory = "Patient" | "Organization" | "Supplier";

export function WaitlistModal({ open, onOpenChange }: WaitlistModalProps) {
  const [step, setStep] = useState<number>(1);
  const [personaCategory, setPersonaCategory] = useState<PersonaCategory>("Patient");
  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [country, setCountry] = useState("Nigeria");
  const [state, setState] = useState("");
  const [city, setCity] = useState("");
  const [organization, setOrganization] = useState("");
  const [orgType, setOrgType] = useState("Clinic");

  // Research State
  const [disconnectedArea, setDisconnectedArea] = useState("Getting results back to my doctor");
  const [breakdownStory, setBreakdownStory] = useState("");
  const [selfCoordinatedTasks, setSelfCoordinatedTasks] = useState("");
  const [externalPartners, setExternalPartners] = useState<string[]>(["Clinics", "Laboratories"]);
  const [communicationChannels, setCommunicationChannels] = useState<string[]>(["WhatsApp messages", "Phone calls"]);

  const [timeline, setTimeline] = useState("Within 3 months");
  const [shapingPreference, setShapingPreference] = useState("Yes, Early Beta Tester");

  const [loading, setLoading] = useState(false);
  const [submitted, setSubmitted] = useState(false);
  const [isDuplicate, setIsDuplicate] = useState(false);
  const [modalFeedback, setModalFeedback] = useState("");

  const validateRequired = () => {
    const missing: string[] = [];
    if (!fullName.trim()) missing.push("Full Name");
    if (!email.trim()) missing.push("Email Address");

    if (missing.length > 0) {
      toast.error(`Please fill in required fields: ${missing.join(", ")}`);
      return false;
    }
    return true;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validateRequired()) return;

    setLoading(true);

    const payload = {
      fullName,
      email,
      phone,
      country,
      state,
      city,
      persona: personaCategory === "Organization" ? orgType : personaCategory,
      organization: organization || orgType,
      biggestPainPoint: breakdownStory || disconnectedArea,
      desiredFeatures: [
        ...externalPartners.map((p) => `Partner: ${p}`),
        ...communicationChannels.map((c) => `Channel: ${c}`),
      ],
      timeline,
      shapingPreference,
      orgType,
      externalPartners,
      communicationChannels,
      disconnectedArea,
      selfCoordinatedTasks,
    };

    const result = await saveWaitlistToSupabase(payload);

    fetch(getApiUrl("/waitlist"), {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    }).catch(() => null);

    setLoading(false);

    if (result.status === "SUCCESS") {
      setSubmitted(true);
      setIsDuplicate(false);
      setModalFeedback(result.message);
      toast.success(result.message);
    } else if (result.status === "DUPLICATE") {
      setSubmitted(true);
      setIsDuplicate(true);
      setModalFeedback("You're already on the Curexal early access list. We'll keep you updated on progress.");
      toast.info("You're already registered on our priority waitlist!");
    } else if (result.status === "VALIDATION_ERROR") {
      toast.error(result.message);
    } else {
      toast.error(result.message || "Unable to save your request right now. Please try again.");
    }
  };

  const resetAndClose = () => {
    setSubmitted(false);
    setStep(1);
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={resetAndClose}>
      <DialogContent className="max-w-lg w-[95vw] sm:w-full bg-white dark:bg-[#0B1120] border border-slate-200 dark:border-slate-800 rounded-2xl sm:rounded-3xl p-3.5 sm:p-6 shadow-2xl z-[999999] max-h-[85dvh] sm:max-h-[90vh] overflow-y-auto font-inter">
        {!submitted && (
          <DialogHeader className="space-y-1 text-left pb-2 border-b border-slate-100 dark:border-slate-800 pr-8">
            <div className="flex items-center gap-2">
              {step > 1 ? (
                <button
                  type="button"
                  onClick={() => setStep(step - 1)}
                  className="inline-flex items-center gap-1 text-xs font-bold text-slate-600 dark:text-slate-300 border-0 bg-slate-100 dark:bg-slate-800/80 px-2.5 py-1 rounded-lg transition-colors cursor-pointer"
                >
                  <ArrowLeft className="w-3.5 h-3.5" />
                  <span>Back</span>
                </button>
              ) : null}

              <div className="flex items-center gap-1.5">
                <div className="w-5 h-5 rounded-lg bg-teal-50 dark:bg-teal-950/60 border border-teal-200 dark:border-teal-800 flex items-center justify-center text-[#0F766E] dark:text-teal-400 font-bold text-[10px]">
                  {step}/5
                </div>
                <span className="text-[10px] sm:text-xs font-bold text-[#0F766E] dark:text-teal-400 uppercase tracking-wider">
                  Early Access Research
                </span>
              </div>
            </div>

            <DialogTitle className="text-sm sm:text-base font-extrabold text-slate-900 dark:text-white tracking-tight pt-1">
              {step === 1 && "Where do you experience healthcare from?"}
              {step === 2 && "Tell us about yourself"}
              {step === 3 && "Where does coordination break down?"}
              {step === 4 && "How soon do you need Curexal?"}
              {step === 5 && "How would you like to participate?"}
            </DialogTitle>
          </DialogHeader>
        )}

        {submitted ? (
          <div className="py-3 space-y-4 animate-in fade-in zoom-in-95 duration-200">
            <div className="p-4 sm:p-6 rounded-2xl bg-teal-50/90 dark:bg-teal-950/40 border border-teal-200/80 dark:border-teal-800/80 text-center space-y-2.5">
              <div className="w-12 h-12 rounded-full bg-[#0F766E] text-white flex items-center justify-center mx-auto shadow-md">
                <CheckCircle2 className="w-6 h-6" />
              </div>
              <div className="space-y-1">
                <h3 className="text-lg sm:text-xl font-black text-slate-900 dark:text-white tracking-tight">
                  {isDuplicate ? `Welcome back, ${fullName}!` : `🎉 You're on the list, ${fullName}!`}
                </h3>
                <p className="text-xs text-slate-600 dark:text-slate-300 max-w-xs mx-auto leading-relaxed">
                  {modalFeedback}
                </p>
              </div>
            </div>

            <div className="pt-2 flex items-center gap-2">
              <Button
                onClick={resetAndClose}
                className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-bold h-10 text-xs rounded-xl cursor-pointer"
              >
                Done
              </Button>
            </div>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="py-2 space-y-3" noValidate>
            {/* STEP 1: Persona */}
            {step === 1 && (
              <div className="space-y-2">
                {[
                  { id: "Patient" as PersonaCategory, label: "Patient", desc: "I receive healthcare and want a connected experience.", icon: User },
                  { id: "Organization" as PersonaCategory, label: "Healthcare Organization", desc: "Clinic, lab, pharmacy, hospital, or diagnostic center.", icon: Stethoscope },
                  { id: "Supplier" as PersonaCategory, label: "Supplier / Partner", desc: "I supply healthcare products, reagents, or services.", icon: ShoppingBag },
                ].map((item) => {
                  const active = personaCategory === item.id;
                  const Icon = item.icon;
                  return (
                    <button
                      key={item.id}
                      type="button"
                      onClick={() => setPersonaCategory(item.id)}
                      className={`w-full p-2.5 rounded-xl border text-left transition-all cursor-pointer flex items-start gap-2.5 ${
                        active
                          ? "bg-teal-50 dark:bg-teal-950 text-[#0F766E] dark:text-teal-300 border-[#0F766E]"
                          : "bg-slate-50 dark:bg-slate-900 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800"
                      }`}
                    >
                      <div className={`w-6 h-6 rounded-lg flex items-center justify-center flex-shrink-0 mt-0.5 ${active ? "bg-[#0F766E] text-white" : "bg-slate-200 dark:bg-slate-800 text-slate-500"}`}>
                        <Icon className="w-3.5 h-3.5" />
                      </div>
                      <div className="flex-1 min-w-0">
                        <h4 className="text-xs font-bold leading-tight">{item.label}</h4>
                        <p className="text-[10px] text-slate-500 dark:text-slate-400 leading-tight mt-0.5">{item.desc}</p>
                      </div>
                    </button>
                  );
                })}

                <Button
                  type="button"
                  onClick={() => setStep(2)}
                  className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-bold h-9 text-xs rounded-xl mt-2 flex items-center justify-center gap-1.5 cursor-pointer border-0"
                >
                  <span>Continue</span>
                  <ArrowRight className="w-3.5 h-3.5" />
                </Button>
              </div>
            )}

            {/* STEP 2: Contact Details */}
            {step === 2 && (
              <div className="space-y-2">
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                  <div>
                    <label className="text-[10px] font-semibold text-slate-700 dark:text-slate-300 block mb-0.5">
                      Full Name <span className="text-rose-500">*</span>
                    </label>
                    <Input
                      type="text"
                      placeholder="Dr. Sarah Johnson"
                      value={fullName}
                      onChange={(e) => setFullName(e.target.value)}
                      required
                      className="h-8 text-xs bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800 rounded-lg"
                    />
                  </div>
                  <div>
                    <label className="text-[10px] font-semibold text-slate-700 dark:text-slate-300 block mb-0.5">
                      Email Address <span className="text-rose-500">*</span>
                    </label>
                    <Input
                      type="email"
                      placeholder="sarah@facility.com"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      required
                      className="h-8 text-xs bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800 rounded-lg"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                  <div>
                    <label className="text-[10px] font-semibold text-slate-700 dark:text-slate-300 block mb-0.5">
                      Phone Number
                    </label>
                    <Input
                      type="tel"
                      placeholder="+234 800 000 0000"
                      value={phone}
                      onChange={(e) => setPhone(e.target.value)}
                      className="h-8 text-xs bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800 rounded-lg"
                    />
                  </div>
                  <div>
                    <label className="text-[10px] font-semibold text-slate-700 dark:text-slate-300 block mb-0.5">
                      Facility / Business (Optional)
                    </label>
                    <Input
                      type="text"
                      placeholder="Apex Diagnostics"
                      value={organization}
                      onChange={(e) => setOrganization(e.target.value)}
                      className="h-8 text-xs bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800 rounded-lg"
                    />
                  </div>
                </div>

                <Button
                  type="button"
                  onClick={() => {
                    if (validateRequired()) setStep(3);
                  }}
                  className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-bold h-9 text-xs rounded-xl mt-2 flex items-center justify-center gap-1.5 cursor-pointer border-0"
                >
                  <span>Continue to Research</span>
                  <ArrowRight className="w-3.5 h-3.5" />
                </Button>
              </div>
            )}

            {/* STEP 3: Research Questions */}
            {step === 3 && (
              <div className="space-y-2.5">
                {personaCategory === "Patient" ? (
                  <>
                    <label className="text-[11px] font-bold text-slate-800 dark:text-slate-200 block">
                      Q: What happens when healthcare feels disconnected?
                    </label>
                    <Textarea
                      placeholder="Tell us what happened last time your care required multiple labs or clinics..."
                      value={breakdownStory}
                      onChange={(e) => setBreakdownStory(e.target.value)}
                      className="text-xs bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800 rounded-lg h-20"
                    />
                  </>
                ) : (
                  <>
                    <label className="text-[11px] font-bold text-slate-800 dark:text-slate-200 block">
                      Q: What happens when you refer a patient or order from partners?
                    </label>
                    <Textarea
                      placeholder="Tell us where you lose visibility or experience manual delays..."
                      value={breakdownStory}
                      onChange={(e) => setBreakdownStory(e.target.value)}
                      className="text-xs bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800 rounded-lg h-20"
                    />
                  </>
                )}

                <Button
                  type="button"
                  onClick={() => setStep(4)}
                  className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-bold h-9 text-xs rounded-xl mt-2 flex items-center justify-center gap-1.5 cursor-pointer border-0"
                >
                  <span>Continue</span>
                  <ArrowRight className="w-3.5 h-3.5" />
                </Button>
              </div>
            )}

            {/* STEP 4: Timeline */}
            {step === 4 && (
              <div className="space-y-2.5">
                <p className="text-[11px] text-slate-500 dark:text-slate-400">
                  How soon do you need Curexal?
                </p>

                <div className="grid grid-cols-2 gap-2">
                  {["Immediately", "Within 3 months", "Within 6 months", "Just exploring"].map((t) => {
                    const active = timeline === t;
                    return (
                      <button
                        key={t}
                        type="button"
                        onClick={() => setTimeline(t)}
                        className={`p-2.5 rounded-xl border text-xs font-bold text-center transition-all cursor-pointer ${
                          active
                            ? "bg-teal-50 dark:bg-teal-950 text-[#0F766E] dark:text-teal-300 border-[#0F766E]"
                            : "bg-slate-50 dark:bg-slate-900 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800"
                        }`}
                      >
                        {t}
                      </button>
                    );
                  })}
                </div>

                <Button
                  type="button"
                  onClick={() => setStep(5)}
                  className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-bold h-9 text-xs rounded-xl mt-2 flex items-center justify-center gap-1.5 cursor-pointer border-0"
                >
                  <span>Final Step</span>
                  <ArrowRight className="w-3.5 h-3.5" />
                </Button>
              </div>
            )}

            {/* STEP 5: Participation */}
            {step === 5 && (
              <div className="space-y-2.5">
                <p className="text-[11px] text-slate-600 dark:text-slate-400">
                  How would you like to participate with our product team?
                </p>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                  {[
                    { id: "Yes, 15-minute customer interview", label: "15-Min User Interview", desc: "Speak directly with system architects", icon: MessageSquare },
                    { id: "Yes, Early Beta Tester", label: "Early Beta Tester", desc: "Get priority access to unreleased features", icon: FlaskRound },
                    { id: "Notify me on Launch Day", label: "Launch Day Alert", desc: "Instant notification when we go live", icon: Rocket },
                    { id: "Keep me updated via Email", label: "Email Updates", desc: "Weekly product & engineering progress", icon: Mail },
                  ].map((opt) => {
                    const active = shapingPreference === opt.id;
                    const Icon = opt.icon;
                    return (
                      <button
                        key={opt.id}
                        type="button"
                        onClick={() => setShapingPreference(opt.id)}
                        className={`p-2.5 rounded-xl border text-left transition-all cursor-pointer flex items-start gap-2 ${
                          active
                            ? "bg-teal-50 dark:bg-teal-950 text-[#0F766E] dark:text-teal-300 border-[#0F766E]"
                            : "bg-slate-50 dark:bg-slate-900 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800"
                        }`}
                      >
                        <div className={`w-5 h-5 rounded-lg flex items-center justify-center flex-shrink-0 mt-0.5 ${active ? "bg-[#0F766E] text-white" : "bg-slate-200 dark:bg-slate-800 text-slate-500"}`}>
                          <Icon className="w-3 h-3" />
                        </div>
                        <div className="flex-1 min-w-0">
                          <h4 className="text-xs font-bold leading-tight">{opt.label}</h4>
                          <p className="text-[9px] text-slate-500 dark:text-slate-400 leading-tight mt-0.5">{opt.desc}</p>
                        </div>
                      </button>
                    );
                  })}
                </div>

                <Button
                  type="submit"
                  disabled={loading}
                  className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-extrabold h-10 text-xs sm:text-sm rounded-xl shadow-md transition-all flex items-center justify-center gap-1.5 mt-2 cursor-pointer border-0"
                >
                  <Sparkles className="w-3.5 h-3.5" />
                  <span>{loading ? "Submitting..." : "Complete Registration"}</span>
                </Button>
              </div>
            )}
          </form>
        )}
      </DialogContent>
    </Dialog>
  );
}
