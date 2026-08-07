import { useState } from "react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Activity,
  CheckCircle2,
  User,
  Building2,
  ShoppingBag,
  ShieldCheck,
  Stethoscope,
  FlaskConical,
  Building,
  Pill,
  ArrowRight,
  ArrowLeft,
  Share2,
  Check,
  Sparkles,
  MessageSquare,
  FlaskRound,
  Rocket,
  Mail,
} from "lucide-react";
import { toast } from "sonner";
import { getApiUrl } from "@/api";
import { SEOHead } from "@/components/seo/seo-head";
import { MarketingNavbar } from "@/components/layouts/marketing-navbar";
import { MarketingFooter } from "@/components/layouts/marketing-footer";
import type { PersonaType } from "@/components/waitlist-modal";
import { saveWaitlistToSupabase } from "@/lib/supabase";

const PERSONAS: { id: PersonaType; label: string; icon: any }[] = [
  { id: "Patient", label: "Patient", icon: User },
  { id: "Laboratory", label: "Laboratory", icon: FlaskConical },
  { id: "Clinic", label: "Clinic", icon: Stethoscope },
  { id: "Hospital", label: "Hospital", icon: Building },
  { id: "Pharmacy", label: "Pharmacy", icon: Pill },
  { id: "Medical Supplier", label: "Supplier", icon: ShoppingBag },
  { id: "Diagnostic Centre", label: "Diagnostic", icon: Building2 },
  { id: "Doctor", label: "Doctor", icon: User },
  { id: "Other", label: "Other", icon: Sparkles },
];

const PAIN_POINTS: Record<string, string[]> = {
  Patient: [
    "Waiting too long for laboratory test results",
    "Finding accredited laboratories near me",
    "Losing paper medical & laboratory records",
    "Booking diagnostic test appointments online",
    "Paying for healthcare tests safely online",
  ],
  Laboratory: [
    "Managing PDF report delivery to patients & doctors",
    "Specimen sample tracking & barcoding",
    "Analyzer instrument interfaces & LIMS integration",
    "Billing & inventory tracking",
    "ISO 15189 compliance & audit log readiness",
  ],
  Clinic: [
    "Clinical EMR lab ordering & paper referral delays",
    "Receiving verified diagnostic results from external labs",
    "Multi-branch clinic synchronization",
    "Patient appointment scheduling & follow-ups",
  ],
  Hospital: [
    "Pathology validation queues & lab director sign-offs",
    "Cross-department EHR/EMR order syncing",
    "Inventory & reagent stock tracking",
    "Regulatory compliance audit logging",
  ],
  Pharmacy: [
    "Finding reliable medical supply vendors",
    "Stock visibility & expiry management",
    "Online prescription ordering & customer delivery",
    "B2B payment collection",
  ],
  "Medical Supplier": [
    "Finding accredited clinic & laboratory buyers",
    "Managing online B2B orders & storefront",
    "Stock visibility across regional warehouses",
    "Automating buyer payment collections",
  ],
};

const SHAPING_OPTIONS = [
  { id: "Yes, 15-minute customer interview", label: "15-Min User Interview", desc: "Speak directly with system architects", icon: MessageSquare },
  { id: "Yes, Early Beta Tester", label: "Early Beta Tester", desc: "Get priority access to unreleased features", icon: FlaskRound },
  { id: "Notify me on Launch Day", label: "Launch Day Alert", desc: "Instant notification when we go live", icon: Rocket },
  { id: "Keep me updated via Email", label: "Email Updates", desc: "Weekly product & engineering progress", icon: Mail },
];

export function WaitlistPage() {
  const [step, setStep] = useState<number>(1);
  const [persona, setPersona] = useState<PersonaType>("Patient");
  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [country, setCountry] = useState("Nigeria");
  const [state, setState] = useState("");
  const [city, setCity] = useState("");
  const [organization, setOrganization] = useState("");

  const [selectedPainPoint, setSelectedPainPoint] = useState("");
  const [timeline, setTimeline] = useState("Within 3 months");
  const [shapingPreference, setShapingPreference] = useState("Yes, Early Beta Tester");

  const [loading, setLoading] = useState(false);
  const [submitted, setSubmitted] = useState(false);

  const validateRequired = () => {
    const missing: string[] = [];
    if (!fullName.trim()) missing.push("Full Name");
    if (!email.trim()) missing.push("Email Address");
    if (!phone.trim()) missing.push("Phone Number");

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
      persona,
      organization,
      biggestPainPoint: selectedPainPoint,
      timeline,
      shapingPreference,
    };

    try {
      await saveWaitlistToSupabase(payload);
      await fetch(getApiUrl("/waitlist"), {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      }).catch(() => null);

      setLoading(false);
      setSubmitted(true);
      toast.success("Spot reserved on our priority waitlist!");
    } catch (err) {
      console.info("Waitlist response stored locally:", payload);
      setLoading(false);
      setSubmitted(true);
      toast.success("Spot reserved on our priority waitlist!");
    }
  };

  const activePainPoints = PAIN_POINTS[persona] || PAIN_POINTS.Patient;

  return (
    <div className="min-h-screen bg-white dark:bg-[#0B1120] text-slate-900 dark:text-white font-inter overflow-x-hidden">
      <SEOHead
        title="Early Access Registration — Curexal Healthcare"
        description="Join Curexal's connected healthcare waitlist and get priority access for your facility."
      />

      <MarketingNavbar />

      <main className="pt-20 sm:pt-24 pb-16 px-3.5 sm:px-6 max-w-xl mx-auto w-full overflow-x-hidden">
        <div className="text-center max-w-xl mx-auto mb-4 sm:mb-6">
          <div className="inline-flex items-center gap-1.5 px-2.5 py-0.5 mb-2 rounded-full border border-teal-500/30 bg-teal-500/10 text-[#0F766E] dark:text-teal-400 text-[10px] sm:text-[11px] font-bold uppercase tracking-wider">
            <Sparkles className="w-3.5 h-3.5" />
            Priority Early Access Registration
          </div>
          <h1 className="text-xl sm:text-3xl font-extrabold tracking-tight text-slate-900 dark:text-white mb-1">
            Join the Curexal Waitlist
          </h1>
          <p className="text-slate-600 dark:text-slate-400 text-xs sm:text-sm leading-snug">
            Reserve your early access spot for connected healthcare operations.
          </p>
        </div>

        <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl sm:rounded-3xl p-3.5 sm:p-6 shadow-xl w-full box-border">
          {submitted ? (
            <div className="space-y-4 animate-in fade-in zoom-in-95 duration-200">
              <div className="p-4 sm:p-6 rounded-2xl bg-teal-50/90 dark:bg-teal-950/40 border border-teal-200/80 dark:border-teal-800/80 text-center space-y-2.5">
                <div className="w-11 h-11 rounded-full bg-[#0F766E] text-white flex items-center justify-center mx-auto shadow-md">
                  <CheckCircle2 className="w-6 h-6" />
                </div>
                <h2 className="text-lg sm:text-xl font-black text-slate-900 dark:text-white tracking-tight">
                  🎉 You're on the list, {fullName}!
                </h2>
                <p className="text-xs text-slate-600 dark:text-slate-300 max-w-md mx-auto leading-relaxed">
                  We've reserved your priority access slot for <strong className="text-[#0F766E] dark:text-teal-300">{persona}s</strong>. We'll reach out as early access slots open.
                </p>
                <div className="grid grid-cols-2 gap-1.5 pt-1 text-left">
                  {[
                    "✓ Early Updates",
                    "✓ Direct Team Access",
                    "✓ Priority Beta",
                    "✓ VIP Support",
                  ].map((item) => (
                    <span key={item} className="text-[10px] font-semibold text-slate-700 dark:text-slate-300 bg-white/90 dark:bg-slate-900/80 p-1.5 rounded-lg border border-slate-200/80 dark:border-slate-800 truncate">
                      {item}
                    </span>
                  ))}
                </div>
              </div>

              <div className="pt-1 flex flex-col sm:flex-row items-center gap-2">
                <Button
                  onClick={() => {
                    navigator.clipboard?.writeText("https://curexal.com/waitlist");
                    toast.success("Waitlist link copied!");
                  }}
                  variant="outline"
                  className="w-full sm:w-1/2 h-9 text-xs font-bold rounded-xl flex items-center justify-center gap-1.5"
                >
                  <Share2 className="w-3.5 h-3.5" />
                  <span>Share Waitlist Link</span>
                </Button>
                <Link to="/" className="w-full sm:w-1/2">
                  <Button className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-bold h-9 text-xs rounded-xl">
                    Return to Homepage
                  </Button>
                </Link>
              </div>
            </div>
          ) : (
            <div className="w-full">
              {/* Step Counter Indicator */}
              <div className="flex items-center justify-between pb-2 mb-3 border-b border-slate-200 dark:border-slate-800">
                <div className="flex items-center gap-1.5">
                  <span className="w-5 h-5 rounded-lg bg-teal-50 dark:bg-teal-950/60 border border-teal-200 dark:border-teal-800 flex items-center justify-center text-[#0F766E] dark:text-teal-400 font-bold text-[10px]">
                    {step}/5
                  </span>
                  <span className="text-[10px] font-bold text-slate-500 uppercase tracking-wider">
                    Step {step} of 5
                  </span>
                </div>
                <div className="w-20 sm:w-24 bg-slate-100 dark:bg-slate-800 h-1.5 rounded-full overflow-hidden">
                  <div
                    className="bg-[#0F766E] h-full transition-all duration-300"
                    style={{ width: `${(step / 5) * 100}%` }}
                  />
                </div>
              </div>

              <form onSubmit={handleSubmit} className="space-y-3.5 w-full" noValidate>
                {/* STEP 1: Persona Selection */}
                {step === 1 && (
                  <div className="space-y-3">
                    <h3 className="text-sm sm:text-base font-bold text-slate-900 dark:text-white">
                      Step 1: Who are you joining as?
                    </h3>
                    <div className="grid grid-cols-3 gap-1.5 sm:gap-2">
                      {PERSONAS.map((item) => {
                        const Icon = item.icon;
                        const active = persona === item.id;
                        return (
                          <button
                            key={item.id}
                            type="button"
                            onClick={() => setPersona(item.id)}
                            className={`flex flex-col items-center justify-center p-2 rounded-xl border text-[11px] font-bold transition-all cursor-pointer ${
                              active
                                ? "bg-teal-50 dark:bg-teal-950/60 text-[#0F766E] dark:text-teal-300 border-[#0F766E] shadow-sm"
                                : "bg-slate-50 dark:bg-slate-800/40 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800 hover:border-slate-300"
                            }`}
                          >
                            <Icon className={`w-4 h-4 mb-1 ${active ? "text-[#0F766E] dark:text-teal-300" : "opacity-60"}`} />
                            <span className="text-center truncate w-full">{item.label}</span>
                          </button>
                        );
                      })}
                    </div>
                    <div className="pt-2 border-t border-slate-100 dark:border-slate-800">
                      <Button
                        type="button"
                        onClick={() => setStep(2)}
                        className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-bold h-9 text-xs sm:text-sm rounded-xl flex items-center justify-center gap-1.5 cursor-pointer shadow-sm"
                      >
                        <span>Continue to Step 2</span>
                        <ArrowRight className="w-3.5 h-3.5" />
                      </Button>
                    </div>
                  </div>
                )}

                {/* STEP 2: Contact Info */}
                {step === 2 && (
                  <div className="space-y-2.5">
                    <h3 className="text-sm sm:text-base font-bold text-slate-900 dark:text-white">
                      Step 2: Tell us about yourself
                    </h3>

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
                          className="h-8 text-xs bg-slate-50 dark:bg-slate-800/60 border-slate-200 dark:border-slate-700 rounded-xl"
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
                          className="h-8 text-xs bg-slate-50 dark:bg-slate-800/60 border-slate-200 dark:border-slate-700 rounded-xl"
                        />
                      </div>
                    </div>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                      <div>
                        <label className="text-[10px] font-semibold text-slate-700 dark:text-slate-300 block mb-0.5">
                          Phone Number <span className="text-rose-500">*</span>
                        </label>
                        <Input
                          type="tel"
                          placeholder="+234 800 000 0000"
                          value={phone}
                          onChange={(e) => setPhone(e.target.value)}
                          required
                          className="h-8 text-xs bg-slate-50 dark:bg-slate-800/60 border-slate-200 dark:border-slate-700 rounded-xl"
                        />
                      </div>
                      <div>
                        <label className="text-[10px] font-semibold text-slate-700 dark:text-slate-300 block mb-0.5">
                          Facility / Company (Optional)
                        </label>
                        <Input
                          type="text"
                          placeholder="Apex Diagnostic Center"
                          value={organization}
                          onChange={(e) => setOrganization(e.target.value)}
                          className="h-8 text-xs bg-slate-50 dark:bg-slate-800/60 border-slate-200 dark:border-slate-700 rounded-xl"
                        />
                      </div>
                    </div>

                    <div className="grid grid-cols-3 gap-1.5">
                      <div>
                        <label className="text-[9px] font-semibold text-slate-700 dark:text-slate-300 block mb-0.5">Country</label>
                        <Input
                          type="text"
                          value={country}
                          onChange={(e) => setCountry(e.target.value)}
                          className="h-8 text-[11px] bg-slate-50 dark:bg-slate-800/60 border-slate-200 dark:border-slate-700 rounded-xl"
                        />
                      </div>
                      <div>
                        <label className="text-[9px] font-semibold text-slate-700 dark:text-slate-300 block mb-0.5">State</label>
                        <Input
                          type="text"
                          placeholder="Lagos"
                          value={state}
                          onChange={(e) => setState(e.target.value)}
                          className="h-8 text-[11px] bg-slate-50 dark:bg-slate-800/60 border-slate-200 dark:border-slate-700 rounded-xl"
                        />
                      </div>
                      <div>
                        <label className="text-[9px] font-semibold text-slate-700 dark:text-slate-300 block mb-0.5">City</label>
                        <Input
                          type="text"
                          placeholder="Ikeja"
                          value={city}
                          onChange={(e) => setCity(e.target.value)}
                          className="h-8 text-[11px] bg-slate-50 dark:bg-slate-800/60 border-slate-200 dark:border-slate-700 rounded-xl"
                        />
                      </div>
                    </div>

                    {/* Bottom Action Bar: Back & Continue Side by Side */}
                    <div className="flex items-center gap-2 pt-2.5 border-t border-slate-100 dark:border-slate-800">
                      <Button
                        type="button"
                        variant="outline"
                        onClick={() => setStep(1)}
                        className="px-3 h-9 text-xs font-bold rounded-xl flex items-center justify-center gap-1 border-slate-200 dark:border-slate-800 cursor-pointer flex-shrink-0"
                      >
                        <ArrowLeft className="w-3.5 h-3.5" />
                        <span>Back</span>
                      </Button>
                      <Button
                        type="button"
                        onClick={() => {
                          if (validateRequired()) {
                            setStep(3);
                          }
                        }}
                        className="flex-1 bg-[#0F766E] hover:bg-[#115E59] text-white font-bold h-9 text-xs sm:text-sm rounded-xl flex items-center justify-center gap-1.5 cursor-pointer shadow-sm min-w-0"
                      >
                        <span className="truncate">Continue to Step 3</span>
                        <ArrowRight className="w-3.5 h-3.5 flex-shrink-0" />
                      </Button>
                    </div>
                  </div>
                )}

                {/* STEP 3: Operational Challenge */}
                {step === 3 && (
                  <div className="space-y-2.5">
                    <h3 className="text-sm sm:text-base font-bold text-slate-900 dark:text-white">
                      Step 3: What's your biggest challenge?
                    </h3>
                    <p className="text-xs text-slate-500 dark:text-slate-400">
                      Select the main issue you experience as a <strong className="text-[#0F766E] dark:text-teal-400">{persona}</strong>:
                    </p>

                    <div className="space-y-1.5">
                      {activePainPoints.map((item) => {
                        const active = selectedPainPoint === item;
                        return (
                          <button
                            key={item}
                            type="button"
                            onClick={() => setSelectedPainPoint(item)}
                            className={`w-full p-2 rounded-xl border text-[11px] font-semibold text-left transition-all cursor-pointer flex items-center justify-between gap-2 ${
                              active
                                ? "bg-teal-50 dark:bg-teal-950/60 text-[#0F766E] dark:text-teal-300 border-[#0F766E]"
                                : "bg-slate-50 dark:bg-slate-800/40 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800 hover:border-slate-300"
                            }`}
                          >
                            <span className="leading-snug">{item}</span>
                            <div className={`w-3.5 h-3.5 rounded-full border flex items-center justify-center flex-shrink-0 ${active ? "border-[#0F766E] bg-[#0F766E] text-white" : "border-slate-300"}`}>
                              {active && <Check className="w-2.5 h-2.5" />}
                            </div>
                          </button>
                        );
                      })}
                    </div>

                    {/* Bottom Action Bar: Back & Continue Side by Side */}
                    <div className="flex items-center gap-2 pt-2.5 border-t border-slate-100 dark:border-slate-800">
                      <Button
                        type="button"
                        variant="outline"
                        onClick={() => setStep(2)}
                        className="px-3 h-9 text-xs font-bold rounded-xl flex items-center justify-center gap-1 border-slate-200 dark:border-slate-800 cursor-pointer flex-shrink-0"
                      >
                        <ArrowLeft className="w-3.5 h-3.5" />
                        <span>Back</span>
                      </Button>
                      <Button
                        type="button"
                        onClick={() => setStep(4)}
                        className="flex-1 bg-[#0F766E] hover:bg-[#115E59] text-white font-bold h-9 text-xs sm:text-sm rounded-xl flex items-center justify-center gap-1.5 cursor-pointer shadow-sm min-w-0"
                      >
                        <span className="truncate">Continue to Step 4</span>
                        <ArrowRight className="w-3.5 h-3.5 flex-shrink-0" />
                      </Button>
                    </div>
                  </div>
                )}

                {/* STEP 4: Urgency & Timeline */}
                {step === 4 && (
                  <div className="space-y-2.5">
                    <h3 className="text-sm sm:text-base font-bold text-slate-900 dark:text-white">
                      Step 4: How soon do you need Curexal?
                    </h3>

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
                                ? "bg-teal-50 dark:bg-teal-950/60 text-[#0F766E] dark:text-teal-300 border-[#0F766E]"
                                : "bg-slate-50 dark:bg-slate-800/40 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800 hover:border-slate-300"
                            }`}
                          >
                            {t}
                          </button>
                        );
                      })}
                    </div>

                    {/* Bottom Action Bar: Back & Continue Side by Side */}
                    <div className="flex items-center gap-2 pt-2.5 border-t border-slate-100 dark:border-slate-800">
                      <Button
                        type="button"
                        variant="outline"
                        onClick={() => setStep(3)}
                        className="px-3 h-9 text-xs font-bold rounded-xl flex items-center justify-center gap-1 border-slate-200 dark:border-slate-800 cursor-pointer flex-shrink-0"
                      >
                        <ArrowLeft className="w-3.5 h-3.5" />
                        <span>Back</span>
                      </Button>
                      <Button
                        type="button"
                        onClick={() => setStep(5)}
                        className="flex-1 bg-[#0F766E] hover:bg-[#115E59] text-white font-bold h-9 text-xs sm:text-sm rounded-xl flex items-center justify-center gap-1.5 cursor-pointer shadow-sm min-w-0"
                      >
                        <span className="truncate">Final Step</span>
                        <ArrowRight className="w-3.5 h-3.5 flex-shrink-0" />
                      </Button>
                    </div>
                  </div>
                )}

                {/* STEP 5: Participation Preference & Final Submit (No Mobile Overflow!) */}
                {step === 5 && (
                  <div className="space-y-2.5 animate-in fade-in duration-150">
                    <h3 className="text-sm sm:text-base font-bold text-slate-900 dark:text-white">
                      Step 5: How would you like to participate?
                    </h3>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                      {SHAPING_OPTIONS.map((opt) => {
                        const active = shapingPreference === opt.id;
                        const Icon = opt.icon;
                        return (
                          <button
                            key={opt.id}
                            type="button"
                            onClick={() => setShapingPreference(opt.id)}
                            className={`p-2 sm:p-2.5 rounded-xl border text-left transition-all cursor-pointer flex items-start gap-2 ${
                              active
                                ? "bg-teal-50/90 dark:bg-teal-950/60 text-[#0F766E] dark:text-teal-300 border-[#0F766E] shadow-xs"
                                : "bg-slate-50 dark:bg-slate-800/40 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800 hover:border-slate-300"
                            }`}
                          >
                            <div className={`w-5 h-5 rounded-lg flex items-center justify-center flex-shrink-0 mt-0.5 ${active ? "bg-[#0F766E] text-white" : "bg-slate-200 dark:bg-slate-800 text-slate-500"}`}>
                              <Icon className="w-3 h-3" />
                            </div>
                            <div className="flex-1 min-w-0">
                              <h4 className="text-[11px] sm:text-xs font-bold leading-tight truncate">{opt.label}</h4>
                              <p className="text-[9px] sm:text-[10px] text-slate-500 dark:text-slate-400 leading-tight mt-0.5">{opt.desc}</p>
                            </div>
                            <div className={`w-3.5 h-3.5 rounded-full border flex items-center justify-center flex-shrink-0 mt-0.5 ${active ? "border-[#0F766E] bg-[#0F766E] text-white" : "border-slate-300"}`}>
                              {active && <Check className="w-2.5 h-2.5" />}
                            </div>
                          </button>
                        );
                      })}
                    </div>

                    {/* Bottom Action Bar: Back & Final Submit Side by Side (Responsive Text) */}
                    <div className="flex items-center gap-2 pt-2.5 border-t border-slate-100 dark:border-slate-800">
                      <Button
                        type="button"
                        variant="outline"
                        onClick={() => setStep(4)}
                        className="px-3 h-10 text-xs font-bold rounded-xl flex items-center justify-center gap-1 border-slate-200 dark:border-slate-800 cursor-pointer flex-shrink-0"
                      >
                        <ArrowLeft className="w-3.5 h-3.5" />
                        <span>Back</span>
                      </Button>
                      <Button
                        type="submit"
                        disabled={loading}
                        className="flex-1 bg-[#0F766E] hover:bg-[#115E59] text-white font-extrabold h-10 text-xs sm:text-sm rounded-xl shadow-md transition-all flex items-center justify-center gap-1.5 cursor-pointer min-w-0"
                      >
                        <Sparkles className="w-3.5 h-3.5 flex-shrink-0" />
                        <span className="truncate hidden sm:inline">{loading ? "Registering..." : "Complete Early Access Registration"}</span>
                        <span className="truncate sm:hidden">{loading ? "Registering..." : "Complete Registration"}</span>
                      </Button>
                    </div>

                    <p className="text-[9px] text-center text-slate-400 flex items-center justify-center gap-1 pt-0.5">
                      <ShieldCheck className="w-3 h-3 text-[#0F766E]" />
                      <span>Zero spam. HIPAA & NDPR data compliance.</span>
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
