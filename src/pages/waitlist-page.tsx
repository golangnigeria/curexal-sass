import { useState } from "react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  CheckCircle2,
  User,
  ShoppingBag,
  ShieldCheck,
  Stethoscope,
  ArrowRight,
  ArrowLeft,
  Share2,
  Check,
  Sparkles,
  MessageSquare,
  FlaskRound,
  Rocket,
  Mail,
  Network,
  Building2,
  FlaskConical,
  Activity,
  ChevronRight,
  Sliders,
} from "lucide-react";
import { toast } from "sonner";
import { getApiUrl } from "@/api";
import { SEOHead } from "@/components/seo/seo-head";
import { MarketingNavbar } from "@/components/layouts/marketing-navbar";
import { MarketingFooter } from "@/components/layouts/marketing-footer";
import { saveWaitlistToSupabase } from "@/lib/supabase";
import { WaitlistStats } from "@/components/home/waitlist-stats";

export type PersonaCategory = "Patient" | "Organization" | "Supplier";

const STEP_TITLES = [
  "Select Your Role",
  "Your Information",
  "Coordination Research",
  "Launch Urgency",
  "Co-Design Participation",
];

export function WaitlistPage() {
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

  // Research Survey State
  const [disconnectedArea, setDisconnectedArea] = useState("Getting results back to my doctor");
  const [breakdownStory, setBreakdownStory] = useState("");
  const [placesInvolved, setPlacesInvolved] = useState("3-5 places");
  const [selfCoordinatedTasks, setSelfCoordinatedTasks] = useState("");
  const [automaticCoordinationDesire, setAutomaticCoordinationDesire] = useState("Auto-deliver lab reports directly to my doctor's chart");
  const [selectedTrustFactors, setSelectedTrustFactors] = useState<string[]>(["Data privacy guarantees"]);

  const [externalPartners, setExternalPartners] = useState<string[]>(["Clinics", "Laboratories"]);
  const [referralBreakdownStory, setReferralBreakdownStory] = useState("");
  const [visibilityLossPoints, setVisibilityLossPoints] = useState<string[]>(["Sample collection status", "Final report delivery"]);
  const [communicationChannels, setCommunicationChannels] = useState<string[]>(["WhatsApp messages", "Phone calls"]);
  const [unresolvedImpact, setUnresolvedImpact] = useState("");
  const [topWorkflowChoice, setTopWorkflowChoice] = useState("Digital Lab Referrals & Results");

  const [timeline, setTimeline] = useState("Within 3 months");
  const [shapingPreference, setShapingPreference] = useState("Yes, Early Beta Tester");

  const [loading, setLoading] = useState(false);
  const [submitted, setSubmitted] = useState(false);
  const [isDuplicateUser, setIsDuplicateUser] = useState(false);
  const [feedbackMessage, setFeedbackMessage] = useState("");

  const toggleArrayItem = (list: string[], setList: (val: string[]) => void, item: string) => {
    if (list.includes(item)) {
      setList(list.filter((i) => i !== item));
    } else {
      setList([...list, item]);
    }
  };

  const validateContactInfo = () => {
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
    if (!validateContactInfo()) return;

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
      biggestPainPoint: breakdownStory || referralBreakdownStory || disconnectedArea,
      desiredFeatures: [
        ...visibilityLossPoints,
        ...externalPartners.map((p) => `Partner: ${p}`),
        ...communicationChannels.map((c) => `Channel: ${c}`),
      ],
      timeline,
      shapingPreference,
      orgType,
      externalPartners,
      referralBreakdownStory,
      visibilityLossPoints,
      communicationChannels,
      unresolvedImpact,
      topWorkflowChoice,
      trustFactors: selectedTrustFactors,
      disconnectedArea,
      placesInvolved,
      selfCoordinatedTasks,
      automaticCoordinationDesire,
      testPilotInterest: shapingPreference,
    };

    const result = await saveWaitlistToSupabase(payload);

    setLoading(false);

    if (result.status === "SUCCESS") {
      setSubmitted(true);
      setIsDuplicateUser(false);
      setFeedbackMessage(result.message);
      toast.success(result.message);
      window.dispatchEvent(new CustomEvent("waitlist-updated"));
    } else if (result.status === "DUPLICATE") {
      setSubmitted(true);
      setIsDuplicateUser(true);
      setFeedbackMessage("You're already on the Curexal early access list. We'll keep you updated on progress.");
      toast.info("You're already registered on our priority waitlist!");
      window.dispatchEvent(new CustomEvent("waitlist-updated"));
    } else if (result.status === "VALIDATION_ERROR") {
      toast.error(result.message);
    } else {
      toast.error(result.message || "Unable to save your request right now. Please try again.");
    }
  };

  return (
    <div className="min-h-screen bg-white dark:bg-[#0B1120] text-slate-900 dark:text-white font-inter overflow-x-hidden relative">
      <SEOHead
        title="Help Us Build Curexal: Early Access & Customer Research"
        description="Join Curexal's customer discovery waitlist and tell us where healthcare coordination breaks down for you."
      />

      <MarketingNavbar />

      {/* Hero Ambient Backdrop */}
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[1000px] h-[450px] bg-gradient-to-b from-teal-500/10 via-teal-500/5 to-transparent blur-3xl pointer-events-none z-0" />

      <main className="relative z-10 pt-24 sm:pt-28 pb-20 px-4 sm:px-6 max-w-4xl mx-auto w-full">
        
        {/* Premium Page Header & Co-Design Positioning */}
        <div className="text-center max-w-2xl mx-auto mb-8 space-y-3">
          <h1 className="text-3xl sm:text-5xl font-black tracking-tight text-slate-900 dark:text-white leading-[1.1]">
            Help Us Build The Healthcare Operating Network.
          </h1>

          <p className="text-slate-600 dark:text-slate-300 text-sm sm:text-base leading-relaxed">
            What would healthcare look like if independent clinics, laboratories, pharmacies, suppliers, and patients could seamlessly coordinate? Tell us where care currently breaks down for you.
          </p>
        </div>

        {/* Live Marketing Statistics Component */}
        <WaitlistStats />

        {/* Co-Design Questionnaire Container */}
        <div className="bg-white dark:bg-[#0B1120]/90 border border-slate-200/80 dark:border-slate-800 rounded-3xl p-5 sm:p-10 shadow-2xl backdrop-blur-xl w-full">
          {submitted ? (
            <div className="space-y-6 py-4 animate-in fade-in zoom-in-95 duration-200">
              <div className="p-8 sm:p-10 rounded-3xl bg-teal-50/90 dark:bg-teal-950/40 border border-teal-200/80 dark:border-teal-800/80 text-center space-y-4">
                <div className="w-14 h-14 rounded-2xl bg-[#0F766E] text-white flex items-center justify-center mx-auto shadow-lg">
                  <CheckCircle2 className="w-8 h-8" />
                </div>
                <div className="space-y-1.5">
                  <h2 className="text-2xl sm:text-3xl font-black text-slate-900 dark:text-white tracking-tight">
                    {isDuplicateUser ? `Welcome Back, ${fullName}!` : `🎉 Thank You, ${fullName}!`}
                  </h2>
                  <p className="text-sm text-slate-600 dark:text-slate-300 max-w-lg mx-auto leading-relaxed font-medium">
                    {feedbackMessage}
                  </p>
                </div>

                <div className="pt-2 flex flex-wrap justify-center gap-2 max-w-md mx-auto text-xs font-semibold text-[#0F766E] dark:text-teal-300">
                  <span className="px-3 py-1.5 rounded-lg bg-white dark:bg-slate-900 border border-teal-200 dark:border-teal-800">
                    ✓ Priority Access Reserved
                  </span>
                  <span className="px-3 py-1.5 rounded-lg bg-white dark:bg-slate-900 border border-teal-200 dark:border-teal-800">
                    ✓ Research Input Recorded
                  </span>
                </div>
              </div>

              <div className="flex flex-col sm:flex-row items-center gap-3">
                <Button
                  onClick={() => {
                    navigator.clipboard?.writeText("https://curexal.com/waitlist");
                    toast.success("Waitlist research link copied!");
                  }}
                  variant="outline"
                  className="w-full sm:w-1/2 h-11 text-xs font-bold rounded-xl flex items-center justify-center gap-2 border-slate-200 dark:border-slate-800 cursor-pointer"
                >
                  <Share2 className="w-4 h-4 text-[#0F766E]" />
                  <span>Share Co-Design Link</span>
                </Button>
                <Link to="/" className="w-full sm:w-1/2">
                  <Button className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-bold h-11 text-xs rounded-xl cursor-pointer border-0 shadow-md">
                    Return to Homepage
                  </Button>
                </Link>
              </div>
            </div>
          ) : (
            <div className="w-full space-y-6">
              
              {/* Stepper Header Bar */}
              <div className="space-y-3 pb-4 border-b border-slate-100 dark:border-slate-800">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2.5">
                    <div className="w-7 h-7 rounded-xl bg-teal-50 dark:bg-teal-950 border border-teal-200 dark:border-teal-800 flex items-center justify-center text-[#0F766E] dark:text-teal-400 font-black text-xs">
                      0{step}
                    </div>
                    <div>
                      <span className="text-[10px] font-extrabold uppercase tracking-widest text-[#0F766E] dark:text-teal-400 block">
                        Stage 0{step} of 05
                      </span>
                      <h2 className="text-sm sm:text-base font-extrabold text-slate-900 dark:text-white tracking-tight">
                        {STEP_TITLES[step - 1]}
                      </h2>
                    </div>
                  </div>

                  {/* Progress Bar */}
                  <div className="hidden sm:flex items-center gap-1.5">
                    {[1, 2, 3, 4, 5].map((s) => (
                      <div
                        key={s}
                        className={`h-2 rounded-full transition-all duration-300 ${
                          s === step
                            ? "w-8 bg-[#0F766E]"
                            : s < step
                            ? "w-4 bg-teal-200 dark:bg-teal-900"
                            : "w-4 bg-slate-200 dark:bg-slate-800"
                        }`}
                      />
                    ))}
                  </div>
                </div>

                {/* Mobile Line Progress */}
                <div className="sm:hidden w-full bg-slate-100 dark:bg-slate-800 h-1.5 rounded-full overflow-hidden">
                  <div
                    className="bg-[#0F766E] h-full transition-all duration-300"
                    style={{ width: `${(step / 5) * 100}%` }}
                  />
                </div>
              </div>

              <form onSubmit={handleSubmit} className="space-y-6" noValidate>
                
                {/* ── STEP 1: PERSONA / ROLE SELECTION ── */}
                {step === 1 && (
                  <div className="space-y-5 animate-in fade-in duration-200">
                    <div className="space-y-1">
                      <h3 className="text-lg sm:text-xl font-black text-slate-900 dark:text-white tracking-tight">
                        Where do you experience healthcare from?
                      </h3>
                      <p className="text-xs text-slate-500 dark:text-slate-400">
                        Select your role so we can tailor our customer discovery research to your exact operational reality.
                      </p>
                    </div>

                    <div className="space-y-2.5">
                      {[
                        {
                          id: "Patient" as PersonaCategory,
                          title: "Patient / Care Recipient",
                          desc: "I receive healthcare and want connected medical records, referrals, and lab results.",
                          icon: User,
                        },
                        {
                          id: "Organization" as PersonaCategory,
                          title: "Healthcare Provider",
                          desc: "Clinic, laboratory, pharmacy, hospital, or diagnostic center operator.",
                          icon: Stethoscope,
                        },
                        {
                          id: "Supplier" as PersonaCategory,
                          title: "Supplier / Partner",
                          desc: "Medical equipment, reagents, pharmaceuticals, or healthcare support services.",
                          icon: ShoppingBag,
                        },
                      ].map((item) => {
                        const Icon = item.icon;
                        const active = personaCategory === item.id;
                        return (
                          <button
                            key={item.id}
                            type="button"
                            onClick={() => setPersonaCategory(item.id)}
                            className={`w-full p-3.5 sm:p-4 rounded-xl sm:rounded-2xl border text-left transition-all cursor-pointer flex items-center justify-between gap-3 ${
                              active
                                ? "bg-teal-50/90 dark:bg-teal-950/60 border-[#0F766E] shadow-sm ring-1 ring-[#0F766E]/30"
                                : "bg-slate-50/60 dark:bg-slate-900/60 border-slate-200 dark:border-slate-800 hover:border-slate-300 dark:hover:border-slate-700"
                            }`}
                          >
                            <div className="flex items-center gap-3 min-w-0 flex-1">
                              <div className={`w-9 h-9 rounded-xl flex items-center justify-center flex-shrink-0 transition-colors ${active ? "bg-[#0F766E] text-white" : "bg-white dark:bg-slate-800 text-slate-500 border border-slate-200 dark:border-slate-700"}`}>
                                <Icon className="w-4.5 h-4.5" />
                              </div>
                              <div className="space-y-0.5 min-w-0 flex-1">
                                <h4 className="text-xs sm:text-sm font-extrabold text-slate-900 dark:text-white leading-tight">
                                  {item.title}
                                </h4>
                                <p className="text-[11px] text-slate-500 dark:text-slate-400 leading-normal truncate">
                                  {item.desc}
                                </p>
                              </div>
                            </div>

                            <div className={`w-5 h-5 rounded-full border flex items-center justify-center flex-shrink-0 transition-colors ${active ? "border-[#0F766E] bg-[#0F766E] text-white" : "border-slate-300 dark:border-slate-700"}`}>
                              {active && <Check className="w-3 h-3" />}
                            </div>
                          </button>
                        );
                      })}
                    </div>

                    <div className="pt-3">
                      <Button
                        type="button"
                        onClick={() => setStep(2)}
                        className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-extrabold h-12 text-sm rounded-xl flex items-center justify-center gap-2 cursor-pointer shadow-md border-0"
                      >
                        <span>Continue to Your Information</span>
                        <ArrowRight className="w-4 h-4" />
                      </Button>
                    </div>
                  </div>
                )}

                {/* ── STEP 2: CONTACT INFORMATION ── */}
                {step === 2 && (
                  <div className="space-y-4 animate-in fade-in duration-200">
                    <div className="space-y-1">
                      <h3 className="text-lg sm:text-xl font-black text-slate-900 dark:text-white tracking-tight">
                        Tell us about yourself
                      </h3>
                      <p className="text-xs text-slate-500 dark:text-slate-400">
                        We'll use your contact details to reserve your priority early access slot.
                      </p>
                    </div>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                      <div>
                        <label className="text-xs font-bold text-slate-800 dark:text-slate-200 block mb-1">
                          Full Name <span className="text-rose-500">*</span>
                        </label>
                        <Input
                          type="text"
                          placeholder="Dr. Sarah Johnson"
                          value={fullName}
                          onChange={(e) => setFullName(e.target.value)}
                          required
                          className="h-10 text-xs bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800 rounded-xl"
                        />
                      </div>
                      <div>
                        <label className="text-xs font-bold text-slate-800 dark:text-slate-200 block mb-1">
                          Email Address <span className="text-rose-500">*</span>
                        </label>
                        <Input
                          type="email"
                          placeholder="sarah@facility.com"
                          value={email}
                          onChange={(e) => setEmail(e.target.value)}
                          required
                          className="h-10 text-xs bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800 rounded-xl"
                        />
                      </div>
                    </div>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                      <div>
                        <label className="text-xs font-bold text-slate-800 dark:text-slate-200 block mb-1">
                          Phone Number
                        </label>
                        <Input
                          type="tel"
                          placeholder="+234 800 000 0000"
                          value={phone}
                          onChange={(e) => setPhone(e.target.value)}
                          className="h-10 text-xs bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800 rounded-xl"
                        />
                      </div>
                      <div>
                        <label className="text-xs font-bold text-slate-800 dark:text-slate-200 block mb-1">
                          Facility / Business Name (Optional)
                        </label>
                        <Input
                          type="text"
                          placeholder="Apex Diagnostics"
                          value={organization}
                          onChange={(e) => setOrganization(e.target.value)}
                          className="h-10 text-xs bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800 rounded-xl"
                        />
                      </div>
                    </div>

                    <div className="grid grid-cols-3 gap-3">
                      <div>
                        <label className="text-[11px] font-bold text-slate-700 dark:text-slate-300 block mb-1">Country</label>
                        <Input
                          type="text"
                          value={country}
                          onChange={(e) => setCountry(e.target.value)}
                          className="h-9 text-xs bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800 rounded-xl"
                        />
                      </div>
                      <div>
                        <label className="text-[11px] font-bold text-slate-700 dark:text-slate-300 block mb-1">State</label>
                        <Input
                          type="text"
                          placeholder="Lagos"
                          value={state}
                          onChange={(e) => setState(e.target.value)}
                          className="h-9 text-xs bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800 rounded-xl"
                        />
                      </div>
                      <div>
                        <label className="text-[11px] font-bold text-slate-700 dark:text-slate-300 block mb-1">City</label>
                        <Input
                          type="text"
                          placeholder="Ikeja"
                          value={city}
                          onChange={(e) => setCity(e.target.value)}
                          className="h-9 text-xs bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800 rounded-xl"
                        />
                      </div>
                    </div>

                    <div className="flex items-center gap-3 pt-4 border-t border-slate-100 dark:border-slate-800">
                      <Button
                        type="button"
                        variant="outline"
                        onClick={() => setStep(1)}
                        className="px-5 h-11 text-xs font-bold rounded-xl flex items-center justify-center gap-1 border-slate-200 dark:border-slate-800 cursor-pointer"
                      >
                        <ArrowLeft className="w-4 h-4" />
                        <span>Back</span>
                      </Button>
                      <Button
                        type="button"
                        onClick={() => {
                          if (validateContactInfo()) setStep(3);
                        }}
                        className="flex-1 bg-[#0F766E] hover:bg-[#115E59] text-white font-extrabold h-11 text-xs sm:text-sm rounded-xl flex items-center justify-center gap-2 cursor-pointer shadow-md border-0"
                      >
                        <span>Continue to Research Questions</span>
                        <ArrowRight className="w-4 h-4" />
                      </Button>
                    </div>
                  </div>
                )}

                {/* ── STEP 3: DEEP COORDINATION DISCOVERY QUESTIONS ── */}
                {step === 3 && (
                  <div className="space-y-5 animate-in fade-in duration-200">
                    {personaCategory === "Patient" && (
                      <div className="space-y-4">
                        <div className="space-y-1">
                          <h3 className="text-lg sm:text-xl font-black text-slate-900 dark:text-white tracking-tight">
                            Patient Care Friction Discovery
                          </h3>
                          <p className="text-xs text-slate-500 dark:text-slate-400">
                            Help us identify where your care journey breaks down between doctors and labs.
                          </p>
                        </div>

                        <div>
                          <label className="text-xs font-bold text-slate-800 dark:text-slate-200 block mb-2">
                            Q1: Where does healthcare feel most disconnected for you?
                          </label>
                          <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                            {[
                              "Finding the right provider",
                              "Moving between providers",
                              "Getting diagnostic tests",
                              "Getting results back to my doctor",
                              "Booking appointments/services",
                              "Getting medicines/products",
                              "Keeping my health records together",
                              "Communicating with doctors",
                            ].map((opt) => (
                              <button
                                key={opt}
                                type="button"
                                onClick={() => setDisconnectedArea(opt)}
                                className={`p-2.5 rounded-xl border text-xs font-semibold text-left transition-all cursor-pointer ${
                                  disconnectedArea === opt
                                    ? "bg-teal-50 dark:bg-teal-950 text-[#0F766E] dark:text-teal-300 border-[#0F766E] font-bold"
                                    : "bg-slate-50 dark:bg-slate-900 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800"
                                }`}
                              >
                                {opt}
                              </button>
                            ))}
                          </div>
                        </div>

                        <div>
                          <label className="text-xs font-bold text-slate-800 dark:text-slate-200 block mb-1.5">
                            Q2: Tell us what happened the last time you experienced this.
                          </label>
                          <Textarea
                            placeholder="e.g. I had to travel across town twice to pick up paper lab results because the lab couldn't send them directly to my doctor..."
                            value={breakdownStory}
                            onChange={(e) => setBreakdownStory(e.target.value)}
                            className="text-xs bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800 rounded-xl h-24"
                          />
                        </div>

                        <div>
                          <label className="text-xs font-bold text-slate-800 dark:text-slate-200 block mb-1.5">
                            Q3: What did you have to do yourself that you wish providers coordinated automatically?
                          </label>
                          <Input
                            placeholder="e.g. Hand-deliver paper test requests and explain my history again..."
                            value={selfCoordinatedTasks}
                            onChange={(e) => setSelfCoordinatedTasks(e.target.value)}
                            className="h-10 text-xs bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800 rounded-xl"
                          />
                        </div>
                      </div>
                    )}

                    {personaCategory === "Organization" && (
                      <div className="space-y-4">
                        <div className="space-y-1">
                          <h3 className="text-lg sm:text-xl font-black text-slate-900 dark:text-white tracking-tight">
                            Healthcare Provider Operations Discovery
                          </h3>
                          <p className="text-xs text-slate-500 dark:text-slate-400">
                            Help us identify where external referrals, result dispatches, or partner communications break down.
                          </p>
                        </div>

                        <div>
                          <label className="text-xs font-bold text-slate-800 dark:text-slate-200 block mb-2">
                            Q1: What type of organization do you operate?
                          </label>
                          <div className="grid grid-cols-3 gap-2">
                            {["Clinic", "Laboratory", "Imaging Center", "Pharmacy", "Hospital", "Diagnostic Centre"].map((t) => (
                              <button
                                key={t}
                                type="button"
                                onClick={() => setOrgType(t)}
                                className={`p-2.5 rounded-xl border text-xs font-bold transition-all cursor-pointer ${
                                  orgType === t
                                    ? "bg-teal-50 dark:bg-teal-950 text-[#0F766E] dark:text-teal-300 border-[#0F766E]"
                                    : "bg-slate-50 dark:bg-slate-900 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800"
                                }`}
                              >
                                {t}
                              </button>
                            ))}
                          </div>
                        </div>

                        <div>
                          <label className="text-xs font-bold text-slate-800 dark:text-slate-200 block mb-2">
                            Q2: Which external partners do you regularly coordinate with? (Select all that apply)
                          </label>
                          <div className="grid grid-cols-2 gap-2">
                            {["Clinics", "Laboratories", "Imaging centers", "Pharmacies", "Hospitals", "Suppliers"].map((partner) => {
                              const selected = externalPartners.includes(partner);
                              return (
                                <button
                                  key={partner}
                                  type="button"
                                  onClick={() => toggleArrayItem(externalPartners, setExternalPartners, partner)}
                                  className={`p-2.5 rounded-xl border text-xs font-semibold text-left transition-all flex items-center justify-between cursor-pointer ${
                                    selected
                                      ? "bg-teal-50 dark:bg-teal-950 text-[#0F766E] dark:text-teal-300 border-[#0F766E]"
                                      : "bg-slate-50 dark:bg-slate-900 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800"
                                  }`}
                                >
                                  <span>{partner}</span>
                                  {selected && <Check className="w-3.5 h-3.5 text-[#0F766E]" />}
                                </button>
                              );
                            })}
                          </div>
                        </div>

                        <div>
                          <label className="text-xs font-bold text-slate-800 dark:text-slate-200 block mb-1.5">
                            Q3: What currently happens when you refer a patient to another provider?
                          </label>
                          <Textarea
                            placeholder="e.g. We give them a paper slip and have no visibility into whether they completed the test..."
                            value={referralBreakdownStory}
                            onChange={(e) => setReferralBreakdownStory(e.target.value)}
                            className="text-xs bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800 rounded-xl h-20"
                          />
                        </div>

                        <div>
                          <label className="text-xs font-bold text-slate-800 dark:text-slate-200 block mb-2">
                            Q4: How do you currently communicate with partner organizations?
                          </label>
                          <div className="grid grid-cols-2 gap-2">
                            {["Phone calls", "WhatsApp messages", "Paper slips", "Email attachments", "Separate software"].map((ch) => {
                              const selected = communicationChannels.includes(ch);
                              return (
                                <button
                                  key={ch}
                                  type="button"
                                  onClick={() => toggleArrayItem(communicationChannels, setCommunicationChannels, ch)}
                                  className={`p-2.5 rounded-xl border text-xs font-semibold text-left transition-all flex items-center justify-between cursor-pointer ${
                                    selected
                                      ? "bg-teal-50 dark:bg-teal-950 text-[#0F766E] dark:text-teal-300 border-[#0F766E]"
                                      : "bg-slate-50 dark:bg-slate-900 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800"
                                  }`}
                                >
                                  <span>{ch}</span>
                                  {selected && <Check className="w-3.5 h-3.5 text-[#0F766E]" />}
                                </button>
                              );
                            })}
                          </div>
                        </div>
                      </div>
                    )}

                    {personaCategory === "Supplier" && (
                      <div className="space-y-4">
                        <div className="space-y-1">
                          <h3 className="text-lg sm:text-xl font-black text-slate-900 dark:text-white tracking-tight">
                            Supplier Order &amp; Settlement Coordination
                          </h3>
                          <p className="text-xs text-slate-500 dark:text-slate-400">
                            Help us simplify B2B supply procurement and buyer fulfillment.
                          </p>
                        </div>

                        <div>
                          <label className="text-xs font-bold text-slate-800 dark:text-slate-200 block mb-2">
                            What healthcare products do you supply?
                          </label>
                          <div className="grid grid-cols-2 gap-2">
                            {["Medical Equipment", "Reagents & Lab Supplies", "Pharmaceuticals", "Support Services"].map((opt) => (
                              <button
                                key={opt}
                                type="button"
                                onClick={() => setOrgType(opt)}
                                className={`p-3 rounded-xl border text-xs font-bold transition-all cursor-pointer ${
                                  orgType === opt
                                    ? "bg-teal-50 dark:bg-teal-950 text-[#0F766E] dark:text-teal-300 border-[#0F766E]"
                                    : "bg-slate-50 dark:bg-slate-900 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800"
                                }`}
                              >
                                {opt}
                              </button>
                            ))}
                          </div>
                        </div>

                        <div>
                          <label className="text-xs font-bold text-slate-800 dark:text-slate-200 block mb-1.5">
                            Where does B2B order &amp; payment coordination break down today?
                          </label>
                          <Textarea
                            placeholder="e.g. Delayed purchase order verification and manual payment collection..."
                            value={referralBreakdownStory}
                            onChange={(e) => setReferralBreakdownStory(e.target.value)}
                            className="text-xs bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800 rounded-xl h-24"
                          />
                        </div>
                      </div>
                    )}

                    <div className="flex items-center gap-3 pt-4 border-t border-slate-100 dark:border-slate-800">
                      <Button
                        type="button"
                        variant="outline"
                        onClick={() => setStep(2)}
                        className="px-5 h-11 text-xs font-bold rounded-xl flex items-center justify-center gap-1 border-slate-200 dark:border-slate-800 cursor-pointer"
                      >
                        <ArrowLeft className="w-4 h-4" />
                        <span>Back</span>
                      </Button>
                      <Button
                        type="button"
                        onClick={() => setStep(4)}
                        className="flex-1 bg-[#0F766E] hover:bg-[#115E59] text-white font-extrabold h-11 text-xs sm:text-sm rounded-xl flex items-center justify-center gap-2 cursor-pointer shadow-md border-0"
                      >
                        <span>Continue to Launch Timeline</span>
                        <ArrowRight className="w-4 h-4" />
                      </Button>
                    </div>
                  </div>
                )}

                {/* ── STEP 4: LAUNCH TIMELINE & URGENCY ── */}
                {step === 4 && (
                  <div className="space-y-4 animate-in fade-in duration-200">
                    <div className="space-y-1">
                      <h3 className="text-lg sm:text-xl font-black text-slate-900 dark:text-white tracking-tight">
                        How soon do you need Curexal?
                      </h3>
                      <p className="text-xs text-slate-500 dark:text-slate-400">
                        Help us prioritize onboarding slots based on operational urgency.
                      </p>
                    </div>

                    <div className="grid grid-cols-2 gap-3">
                      {["Immediately", "Within 3 months", "Within 6 months", "Just exploring"].map((t) => {
                        const active = timeline === t;
                        return (
                          <button
                            key={t}
                            type="button"
                            onClick={() => setTimeline(t)}
                            className={`p-4 rounded-2xl border text-xs font-extrabold text-center transition-all cursor-pointer ${
                              active
                                ? "bg-teal-50 dark:bg-teal-950 text-[#0F766E] dark:text-teal-300 border-[#0F766E] shadow-sm"
                                : "bg-slate-50 dark:bg-slate-900 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800 hover:border-slate-300"
                            }`}
                          >
                            {t}
                          </button>
                        );
                      })}
                    </div>

                    <div className="flex items-center gap-3 pt-4 border-t border-slate-100 dark:border-slate-800">
                      <Button
                        type="button"
                        variant="outline"
                        onClick={() => setStep(3)}
                        className="px-5 h-11 text-xs font-bold rounded-xl flex items-center justify-center gap-1 border-slate-200 dark:border-slate-800 cursor-pointer"
                      >
                        <ArrowLeft className="w-4 h-4" />
                        <span>Back</span>
                      </Button>
                      <Button
                        type="button"
                        onClick={() => setStep(5)}
                        className="flex-1 bg-[#0F766E] hover:bg-[#115E59] text-white font-extrabold h-11 text-xs sm:text-sm rounded-xl flex items-center justify-center gap-2 cursor-pointer shadow-md border-0"
                      >
                        <span>Final Step: Participation</span>
                        <ArrowRight className="w-4 h-4" />
                      </Button>
                    </div>
                  </div>
                )}

                {/* ── STEP 5: CO-DESIGN PARTICIPATION PREFERENCE & FINAL SUBMIT ── */}
                {step === 5 && (
                  <div className="space-y-5 animate-in fade-in duration-200">
                    <div className="space-y-1">
                      <h3 className="text-lg sm:text-xl font-black text-slate-900 dark:text-white tracking-tight">
                        How would you like to participate with our product team?
                      </h3>
                      <p className="text-xs text-slate-500 dark:text-slate-400">
                        Choose your co-design involvement preference.
                      </p>
                    </div>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
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
                            className={`p-3.5 rounded-2xl border text-left transition-all cursor-pointer flex items-start gap-3 ${
                              active
                                ? "bg-teal-50 dark:bg-teal-950 text-[#0F766E] dark:text-teal-300 border-[#0F766E] shadow-sm"
                                : "bg-slate-50 dark:bg-slate-900 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800"
                            }`}
                          >
                            <div className={`w-7 h-7 rounded-xl flex items-center justify-center flex-shrink-0 mt-0.5 ${active ? "bg-[#0F766E] text-white" : "bg-slate-200 dark:bg-slate-800 text-slate-500"}`}>
                              <Icon className="w-4 h-4" />
                            </div>
                            <div className="flex-1 min-w-0">
                              <h4 className="text-xs font-bold leading-tight">{opt.label}</h4>
                              <p className="text-[10px] text-slate-500 dark:text-slate-400 leading-tight mt-0.5">{opt.desc}</p>
                            </div>
                          </button>
                        );
                      })}
                    </div>

                    <div className="flex items-center gap-3 pt-4 border-t border-slate-100 dark:border-slate-800">
                      <Button
                        type="button"
                        variant="outline"
                        onClick={() => setStep(4)}
                        className="px-5 h-12 text-xs font-bold rounded-xl flex items-center justify-center gap-1 border-slate-200 dark:border-slate-800 cursor-pointer"
                      >
                        <ArrowLeft className="w-4 h-4" />
                        <span>Back</span>
                      </Button>
                      <Button
                        type="submit"
                        disabled={loading}
                        className="flex-1 bg-[#0F766E] hover:bg-[#115E59] text-white font-extrabold h-12 text-xs sm:text-sm rounded-xl shadow-lg transition-all flex items-center justify-center gap-2 cursor-pointer border-0"
                      >
                        <Sparkles className="w-4 h-4" />
                        <span>{loading ? "Submitting Co-Design Input..." : "Complete Co-Design Registration"}</span>
                      </Button>
                    </div>

                    <p className="text-[10px] text-center text-slate-400 flex items-center justify-center gap-1 pt-1">
                      <ShieldCheck className="w-3.5 h-3.5 text-[#0F766E]" />
                      <span>Zero spam. Strict tenant data boundaries &amp; NDPR privacy controls.</span>
                    </p>
                  </div>
                )}
              </form>
            </div>
          )}
        </div>
      </main>

      <MarketingFooter />
    </div>
  );
}
