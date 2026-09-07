import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import {
  Phone,
  ShieldCheck,
  CheckCircle2,
  Sparkles,
  Loader2,
  KeyRound,
  UserPlus,
  LogIn,
  HeartPulse,
  Microscope,
  Pill,
  ArrowRight,
  BadgeCheck,
} from "lucide-react";
import {
  useSendPortalOTP,
  useVerifyPortalOTP,
  useRegisterPatient,
} from "@curexal/patient";
import { CurexalLogoSymbol, useBrandTheme } from "@curexal/design-system";

export function PatientPortalLoginPage() {
  const navigate = useNavigate();
  const { orgName, logoUrl } = useBrandTheme();

  // Mode: "LOGIN" or "REGISTER"
  const [authMode, setAuthMode] = useState<"LOGIN" | "REGISTER">("LOGIN");

  // Login Form state
  const [identifier, setIdentifier] = useState("");
  const [otpCode, setOtpCode] = useState("");
  const [loginStep, setLoginStep] = useState<"IDENTIFIER" | "OTP">("IDENTIFIER");
  const [debugOtp, setDebugOtp] = useState<string | null>(null);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  // Self-Registration Form state
  const [firstName, setFirstName] = useState("");
  const [middleName, setMiddleName] = useState("");
  const [lastName, setLastName] = useState("");
  const [gender, setGender] = useState<"MALE" | "FEMALE" | "OTHER">("MALE");
  const [dateOfBirth, setDateOfBirth] = useState("");
  const [regPhone, setRegPhone] = useState("");
  const [regEmail, setRegEmail] = useState("");
  const [bloodGroup, setBloodGroup] = useState("O+");
  const [genotype, setGenotype] = useState("AA");
  const [nin, setNin] = useState("");

  const sendOtp = useSendPortalOTP();
  const verifyOtp = useVerifyPortalOTP();
  const registerPatient = useRegisterPatient();

  const organizationDisplay = orgName || "Curexal Health Network";

  // Handle Login OTP Dispatch
  const handleSendOTP = (e: React.FormEvent) => {
    e.preventDefault();
    if (!identifier.trim()) return;
    setErrorMsg(null);

    sendOtp.mutate(
      { identifier: identifier.trim() },
      {
        onSuccess: (data: any) => {
          setLoginStep("OTP");
          if (data?.debugCode) {
            setDebugOtp(data.debugCode);
          }
        },
        onError: (err: any) => {
          setErrorMsg(
            err.response?.data?.message ||
              err.message ||
              "Failed to send verification code. Please check your phone number or email."
          );
        },
      }
    );
  };

  // Handle Login OTP Verification
  const handleVerifyOTP = (e: React.FormEvent) => {
    e.preventDefault();
    if (otpCode.length !== 6) return;
    setErrorMsg(null);

    verifyOtp.mutate(
      { identifier: identifier.trim(), code: otpCode.trim() },
      {
        onSuccess: (res: any) => {
          const token = res.accessToken || res.token || "portal_jwt_" + identifier;
          localStorage.setItem("curexal_portal_token", token);
          if (res.patient) {
            localStorage.setItem("curexal_portal_patient", JSON.stringify(res.patient));
          } else {
            localStorage.setItem(
              "curexal_portal_patient",
              JSON.stringify({
                id: res.patientId || "",
                identifier,
                mrn: res.mrn,
                name: `${res.firstName || ""} ${res.lastName || ""}`.trim() || identifier,
              })
            );
          }
          navigate("/dashboard");
        },
        onError: (err: any) => {
          setErrorMsg(
            err.response?.data?.message ||
              err.message ||
              "Invalid or expired verification code. Please request a new code."
          );
        },
      }
    );
  };

  // Handle Patient Self-Registration
  const handleSelfRegister = (e: React.FormEvent) => {
    e.preventDefault();
    if (!firstName.trim() || !lastName.trim() || !dateOfBirth || !regPhone.trim()) {
      setErrorMsg("Please fill in all required fields (Name, DOB, and Phone Number).");
      return;
    }
    setErrorMsg(null);

    registerPatient.mutate(
      {
        firstName: firstName.trim(),
        middleName: middleName.trim() || undefined,
        lastName: lastName.trim(),
        gender,
        dateOfBirth,
        phone: regPhone.trim(),
        email: regEmail.trim() || undefined,
        bloodGroup: bloodGroup || undefined,
        genotype: genotype || undefined,
        nin: nin.trim() || undefined,
        registrationChannel: "PORTAL",
      },
      {
        onSuccess: (res: any) => {
          const patient = res.patient;
          const token = res.accessToken || res.token || "portal_jwt_" + patient.id;
          localStorage.setItem("curexal_portal_token", token);
          localStorage.setItem(
            "curexal_portal_patient",
            JSON.stringify({
              id: patient.id,
              userId: patient.userId,
              mrn: patient.mrn,
              name: `${patient.firstName} ${patient.lastName}`,
              identifier: regPhone.trim(),
            })
          );
          navigate("/dashboard");
        },
        onError: (err: any) => {
          if (err.response?.data?.duplicates) {
            setErrorMsg(
              "An account with this phone number or identity already exists. Please sign in instead."
            );
          } else {
            setErrorMsg(
              err.response?.data?.message || err.message || "Registration failed. Please try again."
            );
          }
        },
      }
    );
  };

  return (
    <div className="min-h-screen bg-slate-50 dark:bg-[#0B1120] text-slate-900 dark:text-white font-sans flex flex-col selection:bg-teal-500/20 relative overflow-hidden">
      {/* Background Radial Glow */}
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_80%_60%_at_50%_-20%,rgba(15,118,110,0.12),rgba(255,255,255,0))] dark:bg-[radial-gradient(ellipse_80%_60%_at_50%_-20%,rgba(20,184,166,0.15),rgba(11,17,32,0))] pointer-events-none" />
      <div className="absolute inset-0 bg-[linear-gradient(to_right,#8080800a_1px,transparent_1px),linear-gradient(to_bottom,#8080800a_1px,transparent_1px)] bg-[size:24px_24px] pointer-events-none" />

      {/* Header */}
      <header className="fixed top-0 left-0 right-0 z-50 bg-white/90 dark:bg-[#0B1120]/90 backdrop-blur-xl border-b border-slate-200/80 dark:border-slate-800 shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 h-16 flex items-center justify-between">
          <Link to="/" className="flex items-center gap-3 group">
            <div className="p-1 rounded-xl bg-teal-50 dark:bg-teal-950/50 border border-teal-200/60 dark:border-teal-800/60 shadow-sm group-hover:scale-105 transition-transform flex items-center justify-center">
              {logoUrl ? (
                <img src={logoUrl} alt={organizationDisplay} className="w-8 h-8 object-contain rounded-lg" />
              ) : (
                <CurexalLogoSymbol className="w-8 h-8" />
              )}
            </div>
            <div className="flex flex-col">
              <span className="font-extrabold text-base tracking-tight text-slate-900 dark:text-white">
                {organizationDisplay}
              </span>
              <span className="text-[10px] font-semibold tracking-wider text-teal-600 dark:text-teal-400 uppercase">
                Patient Health Portal
              </span>
            </div>
          </Link>

          <div className="flex items-center gap-3">
            <div className="hidden md:flex items-center gap-2 px-3 py-1.5 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 dark:text-emerald-400 text-xs font-semibold">
              <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
              <span>Encrypted Care Network Active</span>
            </div>
          </div>
        </div>
      </header>

      {/* Main Container */}
      <main className="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 pt-24 pb-12 flex flex-col lg:flex-row items-center justify-between gap-12 z-10 my-auto">
        {/* Left Column: Product Showcase */}
        <div className="flex-1 w-full max-w-xl space-y-6 text-left">
          <div className="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full bg-teal-500/10 border border-teal-500/20 text-teal-700 dark:text-teal-300 text-xs font-semibold">
            <ShieldCheck className="w-4 h-4 text-teal-600 dark:text-teal-400" />
            <span>Zero-Friction Health Records • HIPAA / NDPR Encrypted</span>
          </div>

          <div className="space-y-3">
            <h1 className="text-3xl sm:text-4xl lg:text-5xl font-extrabold tracking-tight text-slate-900 dark:text-white leading-[1.15]">
              Your Unified Health Record & Care Network
            </h1>
            <p className="text-sm sm:text-base text-slate-600 dark:text-slate-300 leading-relaxed">
              Instantly access authorized diagnostic laboratory results, digital prescriptions, clinical visit summaries, and live telehealth consultations in one secure space.
            </p>
          </div>

          <div className="space-y-3.5 pt-2">
            <div className="flex items-start gap-3.5 p-3.5 rounded-2xl bg-white dark:bg-slate-900/70 border border-slate-200/80 dark:border-slate-800 shadow-sm hover:border-teal-500/40 transition-colors">
              <div className="p-2 rounded-xl bg-teal-50 dark:bg-teal-950/60 border border-teal-200/40 dark:border-teal-800/40 text-teal-600 dark:text-teal-400 shrink-0">
                <Microscope className="w-5 h-5" />
              </div>
              <div className="space-y-0.5">
                <h4 className="text-xs font-bold text-slate-900 dark:text-white">
                  Real-Time Laboratory & Imaging Results
                </h4>
                <p className="text-xs text-slate-500 dark:text-slate-400">
                  Electronic sign-off from pathologists and radiologists with in-app PDF and DICOM imaging preview.
                </p>
              </div>
            </div>

            <div className="flex items-start gap-3.5 p-3.5 rounded-2xl bg-white dark:bg-slate-900/70 border border-slate-200/80 dark:border-slate-800 shadow-sm hover:border-teal-500/40 transition-colors">
              <div className="p-2 rounded-xl bg-emerald-50 dark:bg-emerald-950/60 border border-emerald-200/40 dark:border-emerald-800/40 text-emerald-600 dark:text-emerald-400 shrink-0">
                <Pill className="w-5 h-5" />
              </div>
              <div className="space-y-0.5">
                <h4 className="text-xs font-bold text-slate-900 dark:text-white">
                  Smart Digital Prescriptions & Pharmacy Refills
                </h4>
                <p className="text-xs text-slate-500 dark:text-slate-400">
                  Verified medication orders with dosage instructions, refill tracking, and batch safety verification.
                </p>
              </div>
            </div>

            <div className="flex items-start gap-3.5 p-3.5 rounded-2xl bg-white dark:bg-slate-900/70 border border-slate-200/80 dark:border-slate-800 shadow-sm hover:border-teal-500/40 transition-colors">
              <div className="p-2 rounded-xl bg-blue-50 dark:bg-blue-950/60 border border-blue-200/40 dark:border-blue-800/40 text-blue-600 dark:text-blue-400 shrink-0">
                <HeartPulse className="w-5 h-5" />
              </div>
              <div className="space-y-0.5">
                <h4 className="text-xs font-bold text-slate-900 dark:text-white">
                  Virtual Consultations & Direct Provider Matching
                </h4>
                <p className="text-xs text-slate-500 dark:text-slate-400">
                  Intelligent doctor triage matching and encrypted end-to-end video consultations from anywhere.
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Right Column: Auth Card */}
        <div className="w-full max-w-md shrink-0">
          <div className="p-6 sm:p-8 rounded-3xl bg-white/90 dark:bg-[#111827]/90 backdrop-blur-2xl border border-slate-200/90 dark:border-slate-800 shadow-2xl shadow-teal-950/10 space-y-6 relative">
            <div className="space-y-4">
              <div className="text-center space-y-1">
                <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-teal-500/10 text-teal-600 dark:text-teal-400 text-xs font-bold">
                  <Sparkles className="w-3.5 h-3.5" />
                  <span>Patient Health Portal</span>
                </div>
                <h2 className="text-xl font-bold text-slate-900 dark:text-white">
                  {authMode === "LOGIN" ? "Sign In to Your Health Record" : "Create Patient Profile"}
                </h2>
                <p className="text-xs text-slate-500 dark:text-slate-400">
                  {authMode === "LOGIN"
                    ? "Enter your mobile phone or email to receive a secure login code"
                    : "Register in 60 seconds for instant access to your records"}
                </p>
              </div>

              {/* Tab Selector */}
              <div className="flex items-center p-1 rounded-2xl bg-slate-100 dark:bg-slate-900 border border-slate-200 dark:border-slate-800">
                <button
                  type="button"
                  onClick={() => {
                    setAuthMode("LOGIN");
                    setErrorMsg(null);
                  }}
                  className={`flex-1 flex items-center justify-center gap-1.5 py-2 rounded-xl text-xs font-semibold transition-all cursor-pointer ${
                    authMode === "LOGIN"
                      ? "bg-white dark:bg-slate-800 text-teal-700 dark:text-teal-300 shadow-sm border border-slate-200 dark:border-slate-700 font-bold"
                      : "text-slate-500 hover:text-slate-800 dark:hover:text-slate-200"
                  }`}
                >
                  <LogIn className="w-3.5 h-3.5" />
                  <span>Sign In</span>
                </button>
                <button
                  type="button"
                  onClick={() => {
                    setAuthMode("REGISTER");
                    setErrorMsg(null);
                  }}
                  className={`flex-1 flex items-center justify-center gap-1.5 py-2 rounded-xl text-xs font-semibold transition-all cursor-pointer ${
                    authMode === "REGISTER"
                      ? "bg-white dark:bg-slate-800 text-teal-700 dark:text-teal-300 shadow-sm border border-slate-200 dark:border-slate-700 font-bold"
                      : "text-slate-500 hover:text-slate-800 dark:hover:text-slate-200"
                  }`}
                >
                  <UserPlus className="w-3.5 h-3.5" />
                  <span>New Patient</span>
                </button>
              </div>
            </div>

            {errorMsg && (
              <div className="p-3.5 rounded-2xl bg-rose-500/10 border border-rose-500/20 text-rose-600 dark:text-rose-400 text-xs flex items-center gap-2">
                <div className="w-1.5 h-1.5 rounded-full bg-rose-500 shrink-0" />
                <span>{errorMsg}</span>
              </div>
            )}

            {/* OTP LOGIN */}
            {authMode === "LOGIN" && (
              <>
                {loginStep === "IDENTIFIER" ? (
                  <form onSubmit={handleSendOTP} className="space-y-4">
                    <div className="space-y-1.5">
                      <label className="block text-xs font-semibold text-slate-700 dark:text-slate-300">
                        Registered Mobile Phone or Email
                      </label>
                      <div className="relative">
                        <Phone className="w-4 h-4 text-slate-400 absolute left-3.5 top-3" />
                        <input
                          type="text"
                          required
                          value={identifier}
                          onChange={(e) => setIdentifier(e.target.value)}
                          placeholder="+234 801 234 5678 or patient@example.com"
                          className="w-full pl-10 pr-4 py-2.5 bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl text-xs font-medium text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-teal-500"
                        />
                      </div>
                    </div>

                    <button
                      type="submit"
                      disabled={sendOtp.isPending || !identifier.trim()}
                      className="w-full flex items-center justify-center gap-2 py-3 rounded-xl text-white text-xs font-bold shadow-lg shadow-teal-600/20 transition-all disabled:opacity-50 hover:brightness-105 active:scale-[0.99] bg-[#0F766E] hover:bg-[#115E59] cursor-pointer"
                    >
                      {sendOtp.isPending ? (
                        <>
                          <Loader2 className="w-4 h-4 animate-spin" />
                          <span>Sending Verification Code...</span>
                        </>
                      ) : (
                        <>
                          <span>Send Verification Code</span>
                          <ArrowRight className="w-4 h-4" />
                        </>
                      )}
                    </button>
                  </form>
                ) : (
                  <form onSubmit={handleVerifyOTP} className="space-y-4">
                    <div className="space-y-2">
                      <div className="flex items-center justify-between">
                        <label className="block text-xs font-semibold text-slate-700 dark:text-slate-300">
                          Enter 6-Digit Code
                        </label>
                        <button
                          type="button"
                          onClick={() => {
                            setLoginStep("IDENTIFIER");
                            setOtpCode("");
                          }}
                          className="text-[11px] font-bold text-teal-600 dark:text-teal-400 hover:underline cursor-pointer"
                        >
                          Change Number
                        </button>
                      </div>
                      <div className="relative">
                        <KeyRound className="w-4 h-4 text-slate-400 absolute left-3.5 top-3" />
                        <input
                          type="text"
                          required
                          maxLength={6}
                          value={otpCode}
                          onChange={(e) => setOtpCode(e.target.value.trim())}
                          placeholder="123456"
                          autoFocus
                          className="w-full pl-10 pr-4 py-2.5 bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl text-lg font-mono tracking-widest text-center text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-teal-500"
                        />
                      </div>

                      {debugOtp && (
                        <div className="mt-2 p-2.5 rounded-xl bg-teal-500/10 border border-teal-500/20 text-teal-700 dark:text-teal-300 text-[11px] flex items-center justify-between">
                          <span>Active Verification Code:</span>
                          <strong className="font-mono text-sm">{debugOtp}</strong>
                        </div>
                      )}
                    </div>

                    <button
                      type="submit"
                      disabled={otpCode.length !== 6 || verifyOtp.isPending}
                      className="w-full flex items-center justify-center gap-2 py-3 rounded-xl text-white text-xs font-bold shadow-lg shadow-teal-600/20 transition-all disabled:opacity-50 hover:brightness-105 active:scale-[0.99] bg-[#0F766E] hover:bg-[#115E59] cursor-pointer"
                    >
                      {verifyOtp.isPending ? (
                        <>
                          <Loader2 className="w-4 h-4 animate-spin" />
                          <span>Verifying Secure Code...</span>
                        </>
                      ) : (
                        <>
                          <CheckCircle2 className="w-4 h-4" />
                          <span>Verify & Enter Portal</span>
                        </>
                      )}
                    </button>
                  </form>
                )}
              </>
            )}

            {/* PATIENT SELF REGISTRATION */}
            {authMode === "REGISTER" && (
              <form onSubmit={handleSelfRegister} className="space-y-3.5">
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
                  <div>
                    <label className="block text-[11px] font-semibold text-slate-700 dark:text-slate-300 mb-1">
                      First Name <span className="text-rose-500">*</span>
                    </label>
                    <input
                      type="text"
                      required
                      value={firstName}
                      onChange={(e) => setFirstName(e.target.value)}
                      placeholder="e.g. Chinelo"
                      className="w-full px-3 py-2 bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl text-xs text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-1 focus:ring-teal-500"
                    />
                  </div>
                  <div>
                    <label className="block text-[11px] font-semibold text-slate-700 dark:text-slate-300 mb-1">
                      Last Name <span className="text-rose-500">*</span>
                    </label>
                    <input
                      type="text"
                      required
                      value={lastName}
                      onChange={(e) => setLastName(e.target.value)}
                      placeholder="e.g. Okafor"
                      className="w-full px-3 py-2 bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl text-xs text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-1 focus:ring-teal-500"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
                  <div>
                    <label className="block text-[11px] font-semibold text-slate-700 dark:text-slate-300 mb-1">
                      Gender <span className="text-rose-500">*</span>
                    </label>
                    <select
                      value={gender}
                      onChange={(e) => setGender(e.target.value as any)}
                      className="w-full px-3 py-2 bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl text-xs text-slate-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-teal-500"
                    >
                      <option value="MALE">Male</option>
                      <option value="FEMALE">Female</option>
                      <option value="OTHER">Other</option>
                    </select>
                  </div>
                  <div>
                    <label className="block text-[11px] font-semibold text-slate-700 dark:text-slate-300 mb-1">
                      Date of Birth <span className="text-rose-500">*</span>
                    </label>
                    <input
                      type="date"
                      required
                      value={dateOfBirth}
                      onChange={(e) => setDateOfBirth(e.target.value)}
                      className="w-full px-3 py-2 bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl text-xs text-slate-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-teal-500"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
                  <div>
                    <label className="block text-[11px] font-semibold text-slate-700 dark:text-slate-300 mb-1">
                      Mobile Phone <span className="text-rose-500">*</span>
                    </label>
                    <input
                      type="tel"
                      required
                      value={regPhone}
                      onChange={(e) => setRegPhone(e.target.value)}
                      placeholder="+234 801 234 5678"
                      className="w-full px-3 py-2 bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl text-xs text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-1 focus:ring-teal-500"
                    />
                  </div>
                  <div>
                    <label className="block text-[11px] font-semibold text-slate-700 dark:text-slate-300 mb-1">
                      Email (Optional)
                    </label>
                    <input
                      type="email"
                      value={regEmail}
                      onChange={(e) => setRegEmail(e.target.value)}
                      placeholder="patient@example.com"
                      className="w-full px-3 py-2 bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl text-xs text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-1 focus:ring-teal-500"
                    />
                  </div>
                </div>

                <button
                  type="submit"
                  disabled={registerPatient.isPending}
                  className="w-full flex items-center justify-center gap-2 py-3 rounded-xl text-white text-xs font-bold shadow-lg shadow-teal-600/20 transition-all disabled:opacity-50 hover:brightness-105 active:scale-[0.99] bg-[#0F766E] hover:bg-[#115E59] mt-2 cursor-pointer"
                >
                  {registerPatient.isPending ? (
                    <>
                      <Loader2 className="w-4 h-4 animate-spin" />
                      <span>Creating Health Profile...</span>
                    </>
                  ) : (
                    <>
                      <BadgeCheck className="w-4 h-4" />
                      <span>Register & Access Portal</span>
                    </>
                  )}
                </button>
              </form>
            )}
          </div>
        </div>
      </main>
    </div>
  );
}
