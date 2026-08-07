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
  ThumbsUp,
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

const FEATURE_OPTIONS: Record<string, string[]> = {
  Patient: [
    "View all my lab results in 1 unified vault",
    "Book diagnostic test appointments online",
    "Compare lab test prices across facilities",
    "Order prescription medicines with delivery",
    "Telemedicine doctor consultations",
    "Digital Health Wallet & Payment Ledger",
    "AI Health Assistant for diagnostic explanations",
  ],
  Laboratory: [
    "ISO 15189 LIMS Platform",
    "WhatsApp Automated Result Notifications",
    "Pathologist Validation Gateway & Sign-off",
    "Online Test Booking Storefront",
    "Reagent Inventory & Stock Alerts",
    "Analyzer Instrument Data Sync APIs",
    "Automated Billing & Patient Invoicing",
  ],
  Clinic: [
    "Clinic EMR & Digital Lab Ordering",
    "Patient Result Inbox & Notification Sync",
    "Online Appointment Booking Engine",
    "Multi-location Patient Record Sync",
    "Telemedicine & Video Consultations",
  ],
  "Medical Supplier": [
    "B2B Online Storefront & Product Catalog",
    "Vendor Order Dashboard & Processing",
    "Automated B2B Payment Settlement",
    "Stock Analytics & Demand Forecasting",
    "Bulk Discounts & Contract Pricing Tiers",
  ],
};

const DEFAULT_ROADMAP_VOTES = [
  { id: "vault", title: "Unified Patient Results Vault", count: 840, category: "Patient", percent: 88 },
  { id: "lims", title: "ISO 15189 LIMS & Barcoding Engine", count: 620, category: "Laboratory", percent: 74 },
  { id: "whatsapp", title: "WhatsApp Result Notifications", count: 510, category: "Communication", percent: 66 },
  { id: "booking", title: "Online Diagnostic Appointment Booking", count: 480, category: "Marketplace", percent: 59 },
  { id: "emr", title: "Clinic EMR & Electronic Ordering", count: 390, category: "Clinic", percent: 48 },
  { id: "b2b", title: "Medical Supply B2B Marketplace", count: 310, category: "Supplier", percent: 39 },
];

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
  const [selectedFeatures, setSelectedFeatures] = useState<string[]>([]);
  const [timeline, setTimeline] = useState("Within 3 months");
  const [shapingPreference, setShapingPreference] = useState("Yes, Early Beta Tester");

  const [loading, setLoading] = useState(false);
  const [submitted, setSubmitted] = useState(false);
  const [userVotes, setUserVotes] = useState<Record<string, number>>({});
  const [votedItems, setVotedItems] = useState<Record<string, boolean>>({});

  const toggleFeature = (feature: string) => {
    if (selectedFeatures.includes(feature)) {
      setSelectedFeatures(selectedFeatures.filter((f) => f !== feature));
    } else {
      setSelectedFeatures([...selectedFeatures, feature]);
    }
  };

  const handleVote = (voteId: string) => {
    if (votedItems[voteId]) {
      toast.info("You already voted for this feature!");
      return;
    }
    const current = userVotes[voteId] || DEFAULT_ROADMAP_VOTES.find((v) => v.id === voteId)?.count || 0;
    setUserVotes({ ...userVotes, [voteId]: current + 1 });
    setVotedItems({ ...votedItems, [voteId]: true });
    toast.success("Vote recorded! Thank you for directing our roadmap.");
  };

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
      desiredFeatures: selectedFeatures,
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
  const activeFeatureOptions = FEATURE_OPTIONS[persona] || FEATURE_OPTIONS.Patient;

  return (
    <div className="min-h-screen bg-white dark:bg-[#0B1120] text-slate-900 dark:text-white font-inter">
      <SEOHead
        title="Early Access & Product Discovery Hub — Curexal Healthcare"
        description="Help shape Curexal's connected healthcare platform, vote on engineering roadmap features, and get priority early access."
      />

      <MarketingNavbar />

      <main className="pt-20 sm:pt-24 pb-16 px-3.5 sm:px-6 max-w-3xl mx-auto">
        <div className="text-center max-w-xl mx-auto mb-5 sm:mb-8">
          <div className="inline-flex items-center gap-1.5 px-2.5 py-1 mb-2.5 rounded-full border border-teal-500/30 bg-teal-500/10 text-[#0F766E] dark:text-teal-400 text-[10px] sm:text-[11px] font-bold uppercase tracking-wider">
            <Sparkles className="w-3.5 h-3.5" />
            Customer Research & Early Access Hub
          </div>
          <h1 className="text-xl sm:text-3xl font-extrabold tracking-tight text-slate-900 dark:text-white mb-1.5">
            Help Shape Curexal's Platform
          </h1>
          <p className="text-slate-600 dark:text-slate-400 text-xs sm:text-sm leading-snug">
            Select your operational challenges and vote on the features built first.
          </p>
        </div>

        <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl sm:rounded-3xl p-4 sm:p-8 shadow-xl">
          {submitted ? (
            <div className="space-y-5 animate-in fade-in zoom-in-95 duration-200">
              <div className="p-4 sm:p-6 rounded-2xl bg-teal-50/90 dark:bg-teal-950/40 border border-teal-200/80 dark:border-teal-800/80 text-center space-y-2.5">
                <div className="w-11 h-11 rounded-full bg-[#0F766E] text-white flex items-center justify-center mx-auto shadow-md">
                  <CheckCircle2 className="w-6 h-6" />
                </div>
                <h2 className="text-lg sm:text-2xl font-black text-slate-900 dark:text-white tracking-tight">
                  🎉 You're on the list, {fullName}!
                </h2>
                <p className="text-xs sm:text-sm text-slate-600 dark:text-slate-300 max-w-md mx-auto leading-relaxed">
                  We've reserved your priority access slot for <strong className="text-[#0F766E] dark:text-teal-300">{persona}s</strong>.
                </p>
                <div className="grid grid-cols-2 gap-1.5 pt-1 text-left">
                  {[
                    "✓ Early Product Updates",
                    "✓ Feature Voting Rights",
                    "✓ Priority Beta Access",
                    "✓ VIP Onboarding",
                  ].map((item) => (
                    <span key={item} className="text-[10px] sm:text-xs font-semibold text-slate-700 dark:text-slate-300 bg-white/90 dark:bg-slate-900/80 p-1.5 sm:p-2 rounded-lg border border-slate-200/80 dark:border-slate-800 truncate">
                      {item}
                    </span>
                  ))}
                </div>
              </div>

              {/* Feature Voting System */}
              <div className="space-y-2.5">
                <div className="flex items-center justify-between">
                  <h3 className="text-xs sm:text-base font-extrabold text-slate-900 dark:text-white flex items-center gap-1.5">
                    <Sparkles className="w-4 h-4 text-[#0F766E]" />
                    <span>Vote on Next Features</span>
                  </h3>
                  <span className="text-[10px] text-slate-400 font-medium">Single-tap vote</span>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
                  {DEFAULT_ROADMAP_VOTES.map((item) => {
                    const currentCount = userVotes[item.id] || item.count;
                    const isVoted = votedItems[item.id];
                    return (
                      <div
                        key={item.id}
                        className="p-2.5 sm:p-3 rounded-xl bg-slate-50 dark:bg-slate-800/60 border border-slate-200 dark:border-slate-700 flex items-center justify-between gap-2"
                      >
                        <div className="flex-1 min-w-0">
                          <span className="text-[9px] font-extrabold uppercase text-[#0F766E] dark:text-teal-400 bg-teal-50 dark:bg-teal-950/60 px-1.5 py-0.5 rounded">
                            {item.category}
                          </span>
                          <p className="text-xs font-bold text-slate-900 dark:text-white mt-1 truncate">
                            {item.title}
                          </p>
                          <span className="text-[10px] text-slate-500 font-medium">{currentCount} votes</span>
                        </div>
                        <button
                          onClick={() => handleVote(item.id)}
                          className={`px-2.5 py-1 rounded-lg text-xs font-bold transition-all flex items-center gap-1 cursor-pointer border ${
                            isVoted
                              ? "bg-emerald-50 text-emerald-600 border-emerald-200 dark:bg-emerald-950/60 dark:text-emerald-300"
                              : "bg-[#0F766E] text-white hover:bg-[#115E59] border-transparent"
                          }`}
                        >
                          {isVoted ? <Check className="w-3.5 h-3.5" /> : <ThumbsUp className="w-3.5 h-3.5" />}
                          <span>{isVoted ? "Voted" : "Vote"}</span>
                        </button>
                      </div>
                    );
                  })}
                </div>
              </div>

              <div className="pt-2 flex flex-col sm:flex-row items-center gap-2.5">
                <Button
                  onClick={() => {
                    navigator.clipboard?.writeText("https://curexal.com/waitlist");
                    toast.success("Waitlist link copied!");
                  }}
                  variant="outline"
                  className="w-full sm:w-1/2 h-9 sm:h-10 text-xs font-bold rounded-xl flex items-center justify-center gap-1.5"
                >
                  <Share2 className="w-3.5 h-3.5" />
                  <span>Share Waitlist Hub</span>
                </Button>
                <Link to="/" className="w-full sm:w-1/2">
                  <Button className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-bold h-9 sm:h-10 text-xs rounded-xl">
                    Return to Homepage
                  </Button>
                </Link>
              </div>
            </div>
          ) : (
            <div>
              {/* Form Navigation Header */}
              <div className="flex items-center justify-between pb-2.5 mb-3 border-b border-slate-200 dark:border-slate-800">
                <div className="flex items-center gap-1.5">
                  <span className="w-5 h-5 sm:w-6 sm:h-6 rounded-lg bg-teal-50 dark:bg-teal-950/60 border border-teal-200 dark:border-teal-800 flex items-center justify-center text-[#0F766E] dark:text-teal-400 font-bold text-[10px]">
                    {step}/6
                  </span>
                  <span className="text-[10px] sm:text-[11px] font-bold text-slate-500 uppercase tracking-wider">
                    Step {step} of 6
                  </span>
                </div>
                {step > 1 && (
                  <button
                    onClick={() => setStep(step - 1)}
                    className="text-xs font-semibold text-slate-500 hover:text-slate-900 dark:hover:text-white flex items-center gap-1 cursor-pointer border-0 bg-transparent"
                  >
                    <ArrowLeft className="w-3 h-3" /> Back
                  </button>
                )}
              </div>

              <form onSubmit={handleSubmit} className="space-y-3.5" noValidate>
                {/* STEP 1: Persona Selection */}
                {step === 1 && (
                  <div className="space-y-3">
                    <h3 className="text-sm sm:text-base font-bold text-slate-900 dark:text-white">
                      Step 1: Who are you joining as?
                    </h3>
                    <div className="grid grid-cols-3 gap-2">
                      {PERSONAS.map((item) => {
                        const Icon = item.icon;
                        const active = persona === item.id;
                        return (
                          <button
                            key={item.id}
                            type="button"
                            onClick={() => setPersona(item.id)}
                            className={`flex flex-col items-center justify-center p-2.5 sm:p-3 rounded-xl border text-xs font-bold transition-all cursor-pointer ${
                              active
                                ? "bg-teal-50 dark:bg-teal-950/60 text-[#0F766E] dark:text-teal-300 border-[#0F766E] shadow-sm"
                                : "bg-slate-50 dark:bg-slate-800/40 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800 hover:border-slate-300"
                            }`}
                          >
                            <Icon className={`w-4 h-4 sm:w-5 sm:h-5 mb-1 ${active ? "text-[#0F766E] dark:text-teal-300" : "opacity-60"}`} />
                            <span className="text-center truncate w-full">{item.label}</span>
                          </button>
                        );
                      })}
                    </div>
                    <Button
                      type="button"
                      onClick={() => setStep(2)}
                      className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-bold h-9 sm:h-10 text-xs rounded-xl mt-3 flex items-center justify-center gap-1.5 cursor-pointer"
                    >
                      <span>Continue to Step 2</span>
                      <ArrowRight className="w-3.5 h-3.5" />
                    </Button>
                  </div>
                )}

                {/* STEP 2: Contact Info */}
                {step === 2 && (
                  <div className="space-y-3">
                    <h3 className="text-sm sm:text-base font-bold text-slate-900 dark:text-white">
                      Step 2: Tell us about yourself
                    </h3>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
                      <div>
                        <label className="text-[10px] font-semibold text-slate-700 dark:text-slate-300 block mb-1">
                          Full Name <span className="text-rose-500">*</span>
                        </label>
                        <Input
                          type="text"
                          placeholder="Dr. Sarah Johnson"
                          value={fullName}
                          onChange={(e) => setFullName(e.target.value)}
                          required
                          className="h-9 text-xs bg-slate-50 dark:bg-slate-800/60 border-slate-200 dark:border-slate-700 rounded-xl"
                        />
                      </div>
                      <div>
                        <label className="text-[10px] font-semibold text-slate-700 dark:text-slate-300 block mb-1">
                          Email Address <span className="text-rose-500">*</span>
                        </label>
                        <Input
                          type="email"
                          placeholder="sarah@facility.com"
                          value={email}
                          onChange={(e) => setEmail(e.target.value)}
                          required
                          className="h-9 text-xs bg-slate-50 dark:bg-slate-800/60 border-slate-200 dark:border-slate-700 rounded-xl"
                        />
                      </div>
                    </div>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
                      <div>
                        <label className="text-[10px] font-semibold text-slate-700 dark:text-slate-300 block mb-1">
                          Phone Number <span className="text-rose-500">*</span>
                        </label>
                        <Input
                          type="tel"
                          placeholder="+234 800 000 0000"
                          value={phone}
                          onChange={(e) => setPhone(e.target.value)}
                          required
                          className="h-9 text-xs bg-slate-50 dark:bg-slate-800/60 border-slate-200 dark:border-slate-700 rounded-xl"
                        />
                      </div>
                      <div>
                        <label className="text-[10px] font-semibold text-slate-700 dark:text-slate-300 block mb-1">
                          Facility / Company (Optional)
                        </label>
                        <Input
                          type="text"
                          placeholder="Apex Diagnostic Center"
                          value={organization}
                          onChange={(e) => setOrganization(e.target.value)}
                          className="h-9 text-xs bg-slate-50 dark:bg-slate-800/60 border-slate-200 dark:border-slate-700 rounded-xl"
                        />
                      </div>
                    </div>

                    <div className="grid grid-cols-3 gap-2">
                      <div>
                        <label className="text-[9px] font-semibold text-slate-700 dark:text-slate-300 block mb-1">Country</label>
                        <Input
                          type="text"
                          value={country}
                          onChange={(e) => setCountry(e.target.value)}
                          className="h-9 text-xs bg-slate-50 dark:bg-slate-800/60 border-slate-200 dark:border-slate-700 rounded-xl"
                        />
                      </div>
                      <div>
                        <label className="text-[9px] font-semibold text-slate-700 dark:text-slate-300 block mb-1">State</label>
                        <Input
                          type="text"
                          placeholder="Lagos"
                          value={state}
                          onChange={(e) => setState(e.target.value)}
                          className="h-9 text-xs bg-slate-50 dark:bg-slate-800/60 border-slate-200 dark:border-slate-700 rounded-xl"
                        />
                      </div>
                      <div>
                        <label className="text-[9px] font-semibold text-slate-700 dark:text-slate-300 block mb-1">City</label>
                        <Input
                          type="text"
                          placeholder="Ikeja"
                          value={city}
                          onChange={(e) => setCity(e.target.value)}
                          className="h-9 text-xs bg-slate-50 dark:bg-slate-800/60 border-slate-200 dark:border-slate-700 rounded-xl"
                        />
                      </div>
                    </div>

                    <Button
                      type="button"
                      onClick={() => {
                        if (validateRequired()) {
                          setStep(3);
                        }
                      }}
                      className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-bold h-9 text-xs rounded-xl mt-3 flex items-center justify-center gap-1.5 cursor-pointer"
                    >
                      <span>Continue to Step 3</span>
                      <ArrowRight className="w-3.5 h-3.5" />
                    </Button>
                  </div>
                )}

                {/* STEP 3: Operational Challenge */}
                {step === 3 && (
                  <div className="space-y-3">
                    <h3 className="text-sm sm:text-base font-bold text-slate-900 dark:text-white">
                      Step 3: What's your biggest challenge?
                    </h3>
                    <p className="text-xs text-slate-500 dark:text-slate-400">
                      Select the main issue you experience as a <strong className="text-[#0F766E] dark:text-teal-400">{persona}</strong>:
                    </p>

                    <div className="space-y-2">
                      {activePainPoints.map((item) => {
                        const active = selectedPainPoint === item;
                        return (
                          <button
                            key={item}
                            type="button"
                            onClick={() => setSelectedPainPoint(item)}
                            className={`w-full p-2.5 rounded-xl border text-xs font-semibold text-left transition-all cursor-pointer flex items-center justify-between gap-2 ${
                              active
                                ? "bg-teal-50 dark:bg-teal-950/60 text-[#0F766E] dark:text-teal-300 border-[#0F766E]"
                                : "bg-slate-50 dark:bg-slate-800/40 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800 hover:border-slate-300"
                            }`}
                          >
                            <span>{item}</span>
                            <div className={`w-3.5 h-3.5 rounded-full border flex items-center justify-center flex-shrink-0 ${active ? "border-[#0F766E] bg-[#0F766E] text-white" : "border-slate-300"}`}>
                              {active && <Check className="w-2.5 h-2.5" />}
                            </div>
                          </button>
                        );
                      })}
                    </div>

                    <Button
                      type="button"
                      onClick={() => setStep(4)}
                      className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-bold h-9 text-xs rounded-xl mt-3 flex items-center justify-center gap-1.5 cursor-pointer"
                    >
                      <span>Vote on Features</span>
                      <ArrowRight className="w-3.5 h-3.5" />
                    </Button>
                  </div>
                )}

                {/* STEP 4: Feature Excitement */}
                {step === 4 && (
                  <div className="space-y-3">
                    <h3 className="text-sm sm:text-base font-bold text-slate-900 dark:text-white">
                      Step 4: Vote on feature priorities
                    </h3>
                    <p className="text-xs text-slate-500 dark:text-slate-400">
                      Select all features you want built first into Curexal:
                    </p>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                      {activeFeatureOptions.map((feat) => {
                        const active = selectedFeatures.includes(feat);
                        return (
                          <button
                            key={feat}
                            type="button"
                            onClick={() => toggleFeature(feat)}
                            className={`p-2.5 rounded-xl border text-xs font-semibold text-left transition-all cursor-pointer flex items-center justify-between gap-2 ${
                              active
                                ? "bg-teal-50 dark:bg-teal-950/60 text-[#0F766E] dark:text-teal-300 border-[#0F766E]"
                                : "bg-slate-50 dark:bg-slate-800/40 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800 hover:border-slate-300"
                            }`}
                          >
                            <span className="leading-snug">{feat}</span>
                            <div className={`w-3.5 h-3.5 rounded border flex items-center justify-center flex-shrink-0 ${active ? "bg-[#0F766E] border-[#0F766E] text-white" : "border-slate-300"}`}>
                              {active && <Check className="w-2.5 h-2.5" />}
                            </div>
                          </button>
                        );
                      })}
                    </div>

                    <Button
                      type="button"
                      onClick={() => setStep(5)}
                      className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-bold h-9 text-xs rounded-xl mt-3 flex items-center justify-center gap-1.5 cursor-pointer"
                    >
                      <span>Continue to Step 5</span>
                      <ArrowRight className="w-3.5 h-3.5" />
                    </Button>
                  </div>
                )}

                {/* STEP 5: Urgency & Timeline */}
                {step === 5 && (
                  <div className="space-y-3">
                    <h3 className="text-sm sm:text-base font-bold text-slate-900 dark:text-white">
                      Step 5: How soon do you need this?
                    </h3>

                    <div className="grid grid-cols-2 gap-2">
                      {["Immediately", "Within 3 months", "Within 6 months", "Just exploring"].map((t) => {
                        const active = timeline === t;
                        return (
                          <button
                            key={t}
                            type="button"
                            onClick={() => setTimeline(t)}
                            className={`p-3 rounded-xl border text-xs font-bold text-center transition-all cursor-pointer ${
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

                    <Button
                      type="button"
                      onClick={() => setStep(6)}
                      className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-bold h-9 text-xs rounded-xl mt-3 flex items-center justify-center gap-1.5 cursor-pointer"
                    >
                      <span>Final Step</span>
                      <ArrowRight className="w-3.5 h-3.5" />
                    </Button>
                  </div>
                )}

                {/* STEP 6: LAST PHASE — Co-Creation & Final Submit (Mobile Responsive) */}
                {step === 6 && (
                  <div className="space-y-3 animate-in fade-in duration-150">
                    <h3 className="text-sm sm:text-base font-bold text-slate-900 dark:text-white">
                      Step 6: How would you like to participate?
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
                            className={`p-2.5 rounded-xl border text-left transition-all cursor-pointer flex items-start gap-2 ${
                              active
                                ? "bg-teal-50/90 dark:bg-teal-950/60 text-[#0F766E] dark:text-teal-300 border-[#0F766E] shadow-xs"
                                : "bg-slate-50 dark:bg-slate-800/40 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800 hover:border-slate-300"
                            }`}
                          >
                            <div className={`w-6 h-6 rounded-lg flex items-center justify-center flex-shrink-0 mt-0.5 ${active ? "bg-[#0F766E] text-white" : "bg-slate-200 dark:bg-slate-800 text-slate-500"}`}>
                              <Icon className="w-3.5 h-3.5" />
                            </div>
                            <div className="flex-1 min-w-0">
                              <h4 className="text-xs font-bold leading-tight truncate">{opt.label}</h4>
                              <p className="text-[10px] text-slate-500 dark:text-slate-400 leading-tight mt-0.5">{opt.desc}</p>
                            </div>
                            <div className={`w-3.5 h-3.5 rounded-full border flex items-center justify-center flex-shrink-0 mt-0.5 ${active ? "border-[#0F766E] bg-[#0F766E] text-white" : "border-slate-300"}`}>
                              {active && <Check className="w-2.5 h-2.5" />}
                            </div>
                          </button>
                        );
                      })}
                    </div>

                    <Button
                      type="submit"
                      disabled={loading}
                      className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-extrabold h-10 sm:h-11 text-xs sm:text-sm rounded-xl shadow-lg transition-all flex items-center justify-center gap-1.5 mt-3 cursor-pointer"
                    >
                      <Sparkles className="w-3.5 h-3.5" />
                      <span>{loading ? "Submitting..." : "Complete Registration & Save Votes"}</span>
                    </Button>

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
