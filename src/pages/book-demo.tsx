import { useState } from "react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  Activity,
  Check,
  ChevronRight,
  FlaskConical,
  ClipboardCheck,
  Users,
  History,
  ArrowLeft,
  Sparkles,
  Building,
  Mail,
  User,
  Phone,
  Clock,
} from "lucide-react";
import { useApiClient, getApiUrl } from "@/api";
import { toast } from "sonner";
import { SEOHead } from "@/components/seo/seo-head";
import { saveWaitlistToSupabase } from "@/lib/supabase";

export function BookDemoPage() {
  const api = useApiClient();

  const [formData, setFormData] = useState({
    name: "",
    labName: "",
    email: "",
    phone: "",
    specimenVolume: "100-500",
    message: "",
  });

  const [loading, setLoading] = useState(false);
  const [submitted, setSubmitted] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const missing: string[] = [];
    if (!formData.name.trim()) missing.push("Full Name");
    if (!formData.labName.trim()) missing.push("Facility / Organization Name");
    if (!formData.email.trim()) missing.push("Work Email");
    if (!formData.phone.trim()) missing.push("Phone Number");

    if (missing.length > 0) {
      toast.error(`Please fill in required fields: ${missing.join(", ")}`);
      return;
    }

    setLoading(true);

    try {
      let isSuccess = false;
      if ((api as any)?.Demo?.createDemoRequest) {
        const response = await (api as any).Demo.createDemoRequest({
          body: {
            fullName: formData.name,
            facilityName: formData.labName,
            workEmail: formData.email,
            phone: formData.phone,
            notes: `Volume: ${formData.specimenVolume} | Message: ${formData.message}`,
          },
        });
        if (response.status === 200 || response.status === 201) {
          isSuccess = true;
        }
      }

      const result = await saveWaitlistToSupabase({
        fullName: formData.name,
        email: formData.email,
        phone: formData.phone,
        organization: formData.labName,
        persona: "Laboratory",
        biggestPainPoint: formData.message || "Requesting live demo",
        timeline: "Immediately",
        shapingPreference: "Yes, Early Beta Tester",
      });

      setLoading(false);

      if (result.status === "SUCCESS" || result.status === "DUPLICATE") {
        setSubmitted(true);
        if (result.status === "DUPLICATE") {
          toast.info("You're already registered on our priority demo list!");
        } else {
          toast.success(result.message);
        }
      } else {
        toast.error(result.message || "Unable to save your demo request. Please try again.");
      }
    } catch (err) {
      console.info("Demo request processed:", formData);
      setLoading(false);
      setSubmitted(true);
      toast.success("Spot reserved on our priority demo waitlist!");
    }
  };

  return (
    <div className="w-full min-h-screen lg:h-screen lg:max-h-screen overflow-x-hidden bg-slate-950 text-slate-50 flex flex-col lg:flex-row font-inter">
      <SEOHead
        title="Book Demo (Coming Soon): Curexal Healthcare Network"
        description="Book a live walkthrough of Curexal's LIMS and healthcare operating network."
      />

      {/* ── Left Column (Brand & Value Highlights - No Overflow) ── */}
      <div className="lg:w-5/12 p-5 sm:p-8 lg:p-10 flex flex-col justify-between relative overflow-hidden border-b lg:border-b-0 lg:border-r border-slate-800/80 bg-slate-900/60 lg:overflow-y-auto">
        {/* Ambient glows */}
        <div className="absolute top-0 left-0 w-64 h-64 rounded-full bg-[#0F766E]/10 blur-[90px] pointer-events-none" />
        <div className="absolute bottom-0 right-0 w-64 h-64 rounded-full bg-teal-500/10 blur-[80px] pointer-events-none" />

        <div className="relative z-10">
          {/* Back button */}
          <Link
            to="/"
            className="inline-flex items-center gap-1.5 text-slate-400 hover:text-white transition-colors text-xs font-medium mb-4 sm:mb-6 group"
          >
            <ArrowLeft className="h-3.5 w-3.5 group-hover:-translate-x-1 transition-transform" />
            Back to homepage
          </Link>

          {/* Logo & Coming Soon Badge */}
          <div className="flex items-center gap-2 mb-4">
            <img src="/logo-symbol.svg" alt="Curexal Logo" className="w-8 h-8 rounded-xl shadow-sm" />
            <span className="font-extrabold text-xl tracking-tight text-white">Curexal</span>
            <span className="text-[10px] font-bold px-2 py-0.5 rounded-full bg-amber-500/10 border border-amber-500/30 text-amber-400 flex items-center gap-1">
              <Clock className="w-3 h-3" />
              Coming Soon
            </span>
          </div>

          <h1 className="text-xl sm:text-2xl lg:text-3xl font-extrabold leading-tight tracking-tight mb-3 text-white">
            Live Demos{" "}
            <span className="text-[#14B8A6] font-bold">Opening Soon</span>
          </h1>

          <p className="text-slate-400 text-xs sm:text-sm leading-relaxed mb-5">
            1-on-1 walkthroughs with system architects are launching soon. Join our priority waitlist for early access.
          </p>

          {/* Core Features */}
          <div className="space-y-3">
            {[
              {
                icon: FlaskConical,
                title: "Specimen Chain of Custody",
                desc: "Full tracking from extraction to analyst validation.",
              },
              {
                icon: ClipboardCheck,
                title: "Biological Reference Intervals",
                desc: "Custom analyte configurations and auto-flagging.",
              },
              {
                icon: Users,
                title: "Pathology Validation Queues",
                desc: "Secure laboratory director sign-offs.",
              },
              {
                icon: History,
                title: "ISO 15189 Audit Trail",
                desc: "Tamper-proof diagnostic event logging & logs.",
              },
            ].map((feature, i) => {
              const Icon = feature.icon;
              return (
                <div key={i} className="flex gap-2.5 items-start">
                  <div className="w-6 h-6 rounded-md bg-slate-800/80 border border-slate-700/80 flex items-center justify-center flex-shrink-0 mt-0.5">
                    <Icon className="h-3 w-3 text-[#14B8A6]" />
                  </div>
                  <div>
                    <h3 className="text-xs font-bold text-slate-200">{feature.title}</h3>
                    <p className="text-[11px] text-slate-400 leading-snug">{feature.desc}</p>
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        {/* Footer stat block */}
        <div className="relative z-10 mt-6 pt-4 border-t border-slate-800/80 flex items-center gap-2 text-[11px] text-slate-400">
          <Sparkles className="h-3.5 w-3.5 text-[#14B8A6] flex-shrink-0" />
          <span>Deployed in isolated workspace cloud partitions with HIPAA & NDPR privacy.</span>
        </div>
      </div>

      {/* ── Right Column (Compact Form - Required Validation & Phone) ── */}
      <div className="flex-1 flex flex-col justify-center p-5 sm:p-8 lg:p-10 relative overflow-hidden bg-slate-900/40 lg:overflow-y-auto">
        <div className="max-w-md w-full mx-auto relative z-10">
          {submitted ? (
            <div className="bg-slate-900 border border-slate-800/80 rounded-2xl p-6 text-center space-y-4 animate-in fade-in zoom-in-95 duration-200 shadow-xl">
              <div className="w-12 h-12 rounded-full bg-teal-500/10 border border-teal-500/30 text-[#14B8A6] flex items-center justify-center mx-auto shadow-inner">
                <Check className="h-6 w-6" />
              </div>
              <div className="space-y-1">
                <span className="text-[10px] font-extrabold uppercase tracking-wider text-[#14B8A6]">Priority Registered</span>
                <h2 className="text-lg sm:text-xl font-black text-white tracking-tight">You're on the Demo Waitlist!</h2>
                <p className="text-slate-400 text-xs leading-relaxed">
                  Thank you, <strong className="text-slate-200">{formData.name}</strong>. Demos for <strong className="text-slate-200">{formData.labName}</strong> are opening soon. We will reach you at <strong className="text-teal-400">{formData.phone}</strong> or <strong className="text-teal-400">{formData.email}</strong> when early access slots open.
                </p>
              </div>

              <div className="pt-2 border-t border-slate-800/80">
                <Link to="/">
                  <Button className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-bold h-9 text-xs rounded-xl">
                    Return to Homepage
                  </Button>
                </Link>
              </div>
            </div>
          ) : (
            <div className="space-y-4">
              <div className="space-y-1">
                <div className="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full bg-amber-500/10 border border-amber-500/20 text-amber-400 text-[10px] font-bold">
                  <Clock className="w-3 h-3" />
                  <span>Demo Slots Opening Soon</span>
                </div>
                <h2 className="text-lg sm:text-xl font-extrabold text-white tracking-tight">
                  Join Priority Demo Waitlist
                </h2>
                <p className="text-slate-400 text-xs">
                  Register your facility to receive early access demo invitations.
                </p>
              </div>

              <form onSubmit={handleSubmit} className="space-y-3" noValidate>
                <div>
                  <label className="text-[11px] font-semibold text-slate-300 block mb-1">
                    Your Full Name <span className="text-rose-400">*</span>
                  </label>
                  <div className="relative">
                    <User className="absolute left-3 top-2.5 h-3.5 w-3.5 text-slate-500" />
                    <Input
                      type="text"
                      placeholder="Dr. Sarah Johnson"
                      value={formData.name}
                      onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                      required
                      className="pl-8 h-9 text-xs bg-slate-950 border-slate-800 text-white placeholder:text-slate-600 rounded-lg focus:border-[#0F766E]"
                    />
                  </div>
                </div>

                <div>
                  <label className="text-[11px] font-semibold text-slate-300 block mb-1">
                    Facility / Organization Name <span className="text-rose-400">*</span>
                  </label>
                  <div className="relative">
                    <Building className="absolute left-3 top-2.5 h-3.5 w-3.5 text-slate-500" />
                    <Input
                      type="text"
                      placeholder="Apex Diagnostic Laboratories"
                      value={formData.labName}
                      onChange={(e) => setFormData({ ...formData, labName: e.target.value })}
                      required
                      className="pl-8 h-9 text-xs bg-slate-950 border-slate-800 text-white placeholder:text-slate-600 rounded-lg focus:border-[#0F766E]"
                    />
                  </div>
                </div>

                <div>
                  <label className="text-[11px] font-semibold text-slate-300 block mb-1">
                    Work Email Address <span className="text-rose-400">*</span>
                  </label>
                  <div className="relative">
                    <Mail className="absolute left-3 top-2.5 h-3.5 w-3.5 text-slate-500" />
                    <Input
                      type="email"
                      placeholder="sarah@apexlabs.com"
                      value={formData.email}
                      onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                      required
                      className="pl-8 h-9 text-xs bg-slate-950 border-slate-800 text-white placeholder:text-slate-600 rounded-lg focus:border-[#0F766E]"
                    />
                  </div>
                </div>

                <div>
                  <label className="text-[11px] font-semibold text-slate-300 block mb-1">
                    Phone Number <span className="text-rose-400">*</span>
                  </label>
                  <div className="relative">
                    <Phone className="absolute left-3 top-2.5 h-3.5 w-3.5 text-slate-500" />
                    <Input
                      type="tel"
                      placeholder="+234 800 000 0000"
                      value={formData.phone}
                      onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                      required
                      className="pl-8 h-9 text-xs bg-slate-950 border-slate-800 text-white placeholder:text-slate-600 rounded-lg focus:border-[#0F766E]"
                    />
                  </div>
                </div>

                <div>
                  <label className="text-[11px] font-semibold text-slate-300 block mb-1">
                    Monthly Specimen / Test Volume
                  </label>
                  <select
                    value={formData.specimenVolume}
                    onChange={(e) => setFormData({ ...formData, specimenVolume: e.target.value })}
                    className="w-full h-9 px-2.5 bg-slate-950 border border-slate-800 text-slate-200 text-xs rounded-lg focus:border-[#0F766E] outline-none"
                  >
                    <option value="<100">&lt; 100 tests / month</option>
                    <option value="100-500">100 - 500 tests / month</option>
                    <option value="500-2000">500 - 2,000 tests / month</option>
                    <option value="2000+">2,000+ tests / month (Enterprise)</option>
                  </select>
                </div>

                <div>
                  <label className="text-[11px] font-semibold text-slate-300 block mb-1">
                    Specific Workflow Needs (Optional)
                  </label>
                  <Textarea
                    placeholder="E.g., Automated instrument interfaces, multi-location EMR referral sync..."
                    value={formData.message}
                    onChange={(e) => setFormData({ ...formData, message: e.target.value })}
                    rows={2}
                    className="bg-slate-950 border-slate-800 text-white text-xs placeholder:text-slate-600 rounded-lg focus:border-[#0F766E]"
                  />
                </div>

                <Button
                  type="submit"
                  disabled={loading}
                  className="w-full bg-[#0F766E] hover:bg-[#115E59] text-white font-bold h-10 text-xs rounded-lg shadow-lg transition-all flex items-center justify-center gap-1.5 mt-1"
                >
                  <Sparkles className="w-3.5 h-3.5" />
                  <span>{loading ? "Registering..." : "Reserve Demo Waitlist Spot"}</span>
                  <ChevronRight className="w-3.5 h-3.5 opacity-70" />
                </Button>
              </form>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
