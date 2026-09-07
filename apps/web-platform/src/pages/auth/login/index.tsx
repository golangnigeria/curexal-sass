import React, { useState, useEffect } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { motion, AnimatePresence } from "framer-motion";
import { authClient } from "@/lib/auth-client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  ShieldCheck,
  ArrowRight,
  Lock,
  Mail,
  Loader2,
  KeyRound,
  Sparkles,
  CheckCircle2,
  Eye,
  EyeOff,
  Building2,
  Check,
  AlertCircle,
  RefreshCw,
  Crown,
  Stethoscope,
  Activity,
  Microscope,
  Info,
  ChevronRight,
  Zap,
} from "lucide-react";
import { CurexalLogoSymbol } from "@/components/brand/curexal-logo";
import { isPlatformHost, isPatientHost, getCanonicalPlatformUrl, getCanonicalOrgUrl } from "@/lib/url-builder";

interface Day1Persona {
  id: string;
  name: string;
  roleTitle: string;
  email: string;
  category: "Clinic Owner" | "Doctor" | "Nurse" | "HIPAA Test" | "Platform Admin";
  badge: string;
  badgeClass: string;
  description: string;
  expectedFlow: string;
  icon: React.ComponentType<{ className?: string }>;
}

const DAY1_PERSONAS: Day1Persona[] = [
  {
    id: "owner",
    name: "Curexal Clinic Owner",
    roleTitle: "Owner & Managing Director",
    email: "owner@curexal.com",
    category: "Clinic Owner",
    badge: "Multi-Branch (3 Sites)",
    badgeClass: "bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/20",
    description: "Full clinical and administrative governance across all 3 facility branches.",
    expectedFlow: "Prompts deterministic Branch Resolution Modal -> enter workspace with Topbar Switcher",
    icon: Crown,
  },
  {
    id: "dr-emeka",
    name: "Dr. Emeka Okonkwo",
    roleTitle: "Chief Medical Officer / Attending",
    email: "dr.emeka@curexal.com",
    category: "Doctor",
    badge: "Multi-Branch Doctor",
    badgeClass: "bg-teal-500/10 text-teal-600 dark:text-teal-400 border-teal-500/20",
    description: "Attending doctor assigned to Main Campus and Ikeja Specialty Center.",
    expectedFlow: "Prompts Branch Resolution Modal -> direct entrance to Outpatient Clinical EMR",
    icon: Stethoscope,
  },
  {
    id: "dr-sarah",
    name: "Dr. Sarah Alabi",
    roleTitle: "Consultant Specialist",
    email: "dr.sarah@curexal.com",
    category: "Doctor",
    badge: "Single Branch (HO-01)",
    badgeClass: "bg-blue-500/10 text-blue-600 dark:text-blue-400 border-blue-500/20",
    description: "Consultant assigned exclusively to Main Campus (Victoria Island).",
    expectedFlow: "Direct bypass without modal -> immediately connects Clinical Consultation Room",
    icon: Activity,
  },
  {
    id: "nurse-chioma",
    name: "Nurse Chioma Eze",
    roleTitle: "Lead Triage Nurse",
    email: "nurse.chioma@curexal.com",
    category: "Nurse",
    badge: "Single Branch (HO-01)",
    badgeClass: "bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 border-indigo-500/20",
    description: "Triage lead assigned exclusively to Main Campus.",
    expectedFlow: "Direct bypass -> immediately connects Care Desk & Patient Vitals MPI",
    icon: ShieldCheck,
  },
  {
    id: "unassigned",
    name: "Provisional Care Staff",
    roleTitle: "Unassigned Clinical Staff",
    email: "unassigned.staff@curexal.com",
    category: "HIPAA Test",
    badge: "Zero Branches (HIPAA)",
    badgeClass: "bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/20",
    description: "Verified account without active clinic branch provisioning.",
    expectedFlow: "Zero-Trust Enforcement: HTTP 403 UNASSIGNED_FACILITY_BRANCH & remediation modal",
    icon: AlertCircle,
  },
  {
    id: "admin",
    name: "Platform Super Admin",
    roleTitle: "Global Systems & Audit",
    email: "admin@curexal.com",
    category: "Platform Admin",
    badge: "Global Governance",
    badgeClass: "bg-purple-500/10 text-purple-600 dark:text-purple-400 border-purple-500/20",
    description: "Master administrative access with cryptographic audit ledger viewing.",
    expectedFlow: "Direct entrance to Platform Console & Feature #24 Immutable Audit Ledger",
    icon: Sparkles,
  },
];

function resolveAuthorizedDestination(sessionData: any, returnTo?: string | null, branchContextSelected?: boolean): string {
  const effectiveBootstrap = sessionData?.bootstrap;
  const user = sessionData?.user;

  const isPlatformStaff = Boolean(
    effectiveBootstrap?.platform?.isStaff === true ||
    user?.isPlatformAdmin === true ||
    user?.platformRole === "super_admin" ||
    user?.role === "super_admin"
  );

  const isOrgAuthorized = Boolean(
    isPlatformStaff ||
    user?.role === "owner" ||
    user?.role === "org_admin" ||
    user?.role === "org_regional_manager" ||
    user?.role === "org_quality_manager" ||
    user?.role === "org_finance_manager" ||
    user?.role === "org_hr_manager" ||
    effectiveBootstrap?.organization?.role === "owner" ||
    effectiveBootstrap?.organization?.role === "org_admin" ||
    effectiveBootstrap?.contexts?.current === "organization"
  );

  const hasBranchSelected = Boolean(
    branchContextSelected ||
    sessionData?.activeBranch?.id ||
    sessionData?.activeBranch?.slug ||
    effectiveBootstrap?.branch?.id ||
    effectiveBootstrap?.branch?.slug ||
    user?.activeBranchId
  );

  const isWorkspaceAuthorized = Boolean(
    hasBranchSelected ||
    Boolean(effectiveBootstrap?.workspace?.id) ||
    Boolean(effectiveBootstrap?.availableBranches?.length) ||
    effectiveBootstrap?.contexts?.current === "workspace" ||
    Boolean(user?.activeTenantId) ||
    Boolean(user?.workspaceId) ||
    !isOrgAuthorized
  );

  // Validate returnTo if present (prevent open redirect & privilege escalation)
  if (returnTo && returnTo.startsWith("/") && !returnTo.startsWith("//")) {
    if (returnTo.startsWith("/platform/")) {
      if (isPlatformStaff) return returnTo;
    } else if (returnTo.startsWith("/organization/")) {
      if (isOrgAuthorized) return returnTo;
    } else if (returnTo.startsWith("/workspace/")) {
      const bSlug =
        sessionData?.activeBranch?.slug ||
        effectiveBootstrap?.branch?.slug ||
        effectiveBootstrap?.branch?.code?.toLowerCase() ||
        effectiveBootstrap?.workspace?.slug ||
        "curexal-clinic";
      return returnTo.replace("/workspace/", `/${bSlug}/`);
    } else if (isWorkspaceAuthorized) {
      return returnTo;
    }
  }

  // Canonical destination resolution
  if (isPlatformStaff) return "/platform/dashboard";

  // If a branch was explicitly selected or caller is clinic staff:
  if (hasBranchSelected || !isOrgAuthorized) {
    const bSlug =
      sessionData?.activeBranch?.slug ||
      sessionData?.activeBranch?.code?.toLowerCase() ||
      effectiveBootstrap?.branch?.slug ||
      effectiveBootstrap?.branch?.code?.toLowerCase() ||
      effectiveBootstrap?.availableBranches?.[0]?.slug ||
      effectiveBootstrap?.availableBranches?.[0]?.code?.toLowerCase() ||
      effectiveBootstrap?.workspace?.slug ||
      "curexal-clinic";

    const userRoleStr = (user?.role || effectiveBootstrap?.identity?.role || "").toLowerCase();
    let mod = "dashboard";
    if (userRoleStr.includes("doc") || userRoleStr.includes("clin") || userRoleStr.includes("phys")) {
      mod = "clinical";
    } else if (userRoleStr.includes("nurse") || userRoleStr.includes("triage")) {
      mod = "care-desk";
    } else if (userRoleStr.includes("recep") || userRoleStr.includes("front")) {
      mod = "reception";
    } else if (userRoleStr.includes("cash") || userRoleStr.includes("acc") || userRoleStr.includes("bill")) {
      mod = "billing";
    } else {
      mod = "dashboard";
    }
    return `/${bSlug}/${mod}`;
  }

  if (isOrgAuthorized) return "/organization/dashboard";

  return "/curexal-clinic/dashboard";
}

export default function LoginPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const returnTo = searchParams.get("returnTo") || searchParams.get("redirect");
  const { data: session, isPending } = authClient.useSession();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const showDevPersonas = searchParams.get("dev") === "true" || searchParams.get("demo") === "true";

  // Active persona highlight for Day 1 test verification (when dev mode is active)
  const [selectedPersonaId, setSelectedPersonaId] = useState<string | null>(null);

  // Branch Selection Dialog State
  const [isBranchSelectOpen, setIsBranchSelectOpen] = useState(false);
  const [branchSelectionData, setBranchSelectionData] = useState<{
    selectionToken: string;
    assignedBranches: Array<{
      id: string;
      name: string;
      code: string;
      slug: string;
      facilityType?: string;
      isHeadquarters?: boolean;
    }>;
    organization?: { id: string; name: string; slug: string };
  } | null>(null);
  const [selectedBranchId, setSelectedBranchId] = useState<string>("");
  const [isBranchSelectLoading, setIsBranchSelectLoading] = useState(false);
  const [branchSelectError, setBranchSelectError] = useState<string | null>(null);

  // Unassigned Facility Branch Dialog State
  const [isUnassignedModalOpen, setIsUnassignedModalOpen] = useState(false);
  const [unassignedOrgData, setUnassignedOrgData] = useState<{ id?: string; name?: string; slug?: string } | null>(null);

  // Forgot Password Dialog State
  const [isForgotOpen, setIsForgotOpen] = useState(false);
  const [forgotEmail, setForgotEmail] = useState("");
  const [isForgotLoading, setIsForgotLoading] = useState(false);
  const [forgotError, setForgotError] = useState<string | null>(null);
  const [forgotSuccess, setForgotSuccess] = useState<string | null>(null);

  // Set Password Dialog State
  const [isSetOpen, setIsSetOpen] = useState(false);
  const [setEmailVal, setSetEmailVal] = useState("");
  const [setCodeVal, setSetCodeVal] = useState("");
  const [setNewPassword, setSetNewPassword] = useState("");
  const [setConfirmPassword, setSetConfirmPassword] = useState("");
  const [showSetPassword, setShowSetPassword] = useState(false);
  const [isSetLoading, setIsSetLoading] = useState(false);
  const [setErrorMsg, setSetErrorMsg] = useState<string | null>(null);
  const [setSuccessMsg, setSetSuccessMsg] = useState<string | null>(null);

  const navigateToCanonicalDestination = (sessionPayload: any, returnToParam?: string | null, branchContextSelected?: boolean) => {
    const destination = resolveAuthorizedDestination(sessionPayload, returnToParam, branchContextSelected);
    if (destination.startsWith("/platform/") && !isPlatformHost()) {
      const canonicalPlatformUrl = getCanonicalPlatformUrl(destination);
      if (canonicalPlatformUrl !== destination) {
        window.location.replace(canonicalPlatformUrl);
        return;
      }
    }
    const orgSlug = sessionPayload?.bootstrap?.organization?.slug;
    if (destination.startsWith("/organization/") && isPlatformHost() && orgSlug) {
      const canonicalOrgUrl = getCanonicalOrgUrl(orgSlug, destination);
      if (canonicalOrgUrl !== destination) {
        window.location.replace(canonicalOrgUrl);
        return;
      }
    }
    navigate(destination, { replace: true });
  };

  // If on patient portal subdomain, redirect directly to portal login
  useEffect(() => {
    if (isPatientHost()) {
      navigate("/portal/login", { replace: true });
    }
  }, [navigate]);

  // If already logged in, redirect to authorized dashboard
  useEffect(() => {
    if (!isPending && session?.user && !isPatientHost()) {
      navigateToCanonicalDestination(session, returnTo);
    }
  }, [session, isPending, returnTo]);

  const handleSelectPersona = (persona: Day1Persona, autoSubmit = false) => {
    setSelectedPersonaId(persona.id);
    setEmail(persona.email);
    setPassword("password");
    setError(null);

    if (autoSubmit) {
      executeSignIn(persona.email, "password");
    }
  };

  const executeSignIn = async (userEmail: string, userPass: string) => {
    setIsLoading(true);
    setError(null);

    try {
      const res: any = await authClient.signIn({ email: userEmail, password: userPass });

      // 1. Cross-host redirection with single-use exchange token
      if (res?.targetUrl || res?.status === "redirect_required") {
        window.location.replace(res.targetUrl);
        return;
      }

      // 2. Branch selection required modal
      if (res?.requireBranchSelection || res?.status === "branch_selection_required") {
        setBranchSelectionData({
          selectionToken: res.selectionToken,
          assignedBranches: res.assignedBranches || [],
          organization: res.organization,
        });
        if (res.assignedBranches && res.assignedBranches.length > 0) {
          setSelectedBranchId(res.assignedBranches[0].id);
        }
        setIsBranchSelectOpen(true);
        return;
      }

      // 3. Unassigned facility branch (HIPAA guard)
      if (res?.status === "unassigned_facility_branch") {
        setUnassignedOrgData(res.organization || null);
        setIsUnassignedModalOpen(true);
        return;
      }

      // 4. Organization access required
      if (res?.status === "organization_access_required") {
        setError(res.message || "No active organization memberships found for your account.");
        return;
      }

      // 5. Directly authenticated
      if (res?.destinationPath) {
        window.location.replace(res.destinationPath);
        return;
      }

      navigateToCanonicalDestination(res?.session || res, returnTo);
    } catch (err: any) {
      const errData = err?.details || err?.response?.data || err?.data;
      const errCode = errData?.error?.code || errData?.code || err?.code;
      if (errCode === "UNASSIGNED_FACILITY_BRANCH" || (err?.message && err.message.includes("No operational branch assignment found"))) {
        setUnassignedOrgData(errData?.error?.organization || null);
        setIsUnassignedModalOpen(true);
        return;
      }
      const msg =
        errData?.error?.message ||
        errData?.message ||
        err?.message ||
        "Authentication failed. Please check your credentials.";
      setError(msg);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || !password) {
      setError("Please enter both email and password.");
      return;
    }
    await executeSignIn(email, password);
  };

  const handleBranchSelectSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!branchSelectionData?.selectionToken || !selectedBranchId) {
      setBranchSelectError("Please select a clinical location to continue.");
      return;
    }

    setIsBranchSelectLoading(true);
    setBranchSelectError(null);

    try {
      const res: any = await authClient.selectBranch({
        selectionToken: branchSelectionData.selectionToken,
        branchId: selectedBranchId,
      });

      if (res?.targetUrl || res?.status === "redirect_required") {
        window.location.replace(res.targetUrl);
        return;
      }

      setIsBranchSelectOpen(false);
      if (res?.destinationPath) {
        window.location.replace(res.destinationPath);
        return;
      }
      navigateToCanonicalDestination(res?.session || res, returnTo, true);
    } catch (err: any) {
      const msg =
        err?.response?.data?.message ||
        err?.message ||
        "Branch authorization failed. Your session may have expired.";
      setBranchSelectError(msg);
    } finally {
      setIsBranchSelectLoading(false);
    }
  };

  const handleOpenForgot = () => {
    setForgotEmail(email || "");
    setForgotError(null);
    setForgotSuccess(null);
    setIsForgotOpen(true);
  };

  const handleForgotSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!forgotEmail) {
      setForgotError("Please enter your registered email address.");
      return;
    }
    setIsForgotLoading(true);
    setForgotError(null);
    setForgotSuccess(null);

    try {
      const res: any = await authClient.forgotPassword(forgotEmail);
      const msg =
        res?.message ||
        "If this email is registered and verified, password instructions have been sent to your inbox.";
      setForgotSuccess(msg);
    } catch (err: any) {
      const msg =
        err?.response?.data?.message ||
        err?.message ||
        "Unable to send recovery email. Please check your details and try again.";
      setForgotError(msg);
    } finally {
      setIsForgotLoading(false);
    }
  };

  const handleOpenSetPassword = () => {
    setSetEmailVal(email || "");
    setSetCodeVal("");
    setSetNewPassword("");
    setSetConfirmPassword("");
    setSetErrorMsg(null);
    setSetSuccessMsg(null);
    setIsSetOpen(true);
  };

  const handleSetPasswordSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!setEmailVal || !setCodeVal || !setNewPassword) {
      setSetErrorMsg("Please fill in all required fields.");
      return;
    }
    if (setNewPassword.length < 8) {
      setSetErrorMsg("Password must be at least 8 characters long.");
      return;
    }
    if (setNewPassword !== setConfirmPassword) {
      setSetErrorMsg("Passwords do not match.");
      return;
    }

    setIsSetLoading(true);
    setSetErrorMsg(null);
    setSetSuccessMsg(null);

    try {
      const isToken = setCodeVal.length > 10;
      await authClient.setPassword({
        email: setEmailVal,
        code: isToken ? undefined : setCodeVal,
        token: isToken ? setCodeVal : undefined,
        password: setNewPassword,
      });
      setSetSuccessMsg("Password set successfully! You can now log in to your account.");
      setEmail(setEmailVal);
    } catch (err: any) {
      const msg =
        err?.response?.data?.message ||
        err?.message ||
        "Failed to set password. The code or token may be expired or invalid.";
      setSetErrorMsg(msg);
    } finally {
      setIsSetLoading(false);
    }
  };

  const activePersona = DAY1_PERSONAS.find((p) => p.id === selectedPersonaId);

  return (
    <div className="relative min-h-screen w-full flex items-center justify-center bg-background px-4 py-8 overflow-x-hidden">
      {/* Dynamic Ambient Background Elements */}
      <div className="absolute top-1/4 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[700px] h-[400px] bg-primary/15 rounded-full blur-[140px] pointer-events-none -z-10" />
      <div className="absolute bottom-10 right-10 w-[400px] h-[300px] bg-teal-500/10 rounded-full blur-[120px] pointer-events-none -z-10" />

      <motion.div
        initial={{ opacity: 0, y: 15 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.4, ease: "easeOut" }}
        className="w-full max-w-5xl space-y-8"
      >
        {/* Top Header & Brand Identity */}
        <div className="flex flex-col items-center text-center space-y-3">
          <div className="flex items-center gap-3">
            <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-gradient-to-br from-primary to-teal-800 text-white shadow-lg shadow-primary/25">
              <CurexalLogoSymbol className="w-8 h-8" />
            </div>
            <div className="text-left">
              <span className="text-xl font-bold tracking-tight text-foreground block">
                Curexal
              </span>
              <span className="text-xs uppercase tracking-widest text-primary font-semibold">
                Healthcare Operating System
              </span>
            </div>
          </div>
          <p className="text-sm text-muted-foreground max-w-lg">
            Zero-Trust Clinical Workspace, Multi-Tenant Operations & HIPAA § 164.528 Compliance Engine
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
          {/* Left Column: Enterprise Healthcare Feature & Security Showcase */}
          <div className="lg:col-span-7 space-y-6 text-left">
            <div className="space-y-3">
              <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-primary/10 border border-primary/20 text-primary text-xs font-semibold">
                <ShieldCheck className="w-3.5 h-3.5" />
                <span>Zero-Trust Architecture • HIPAA § 164.528 Verified</span>
              </div>
              <h1 className="text-3xl sm:text-4xl font-extrabold tracking-tight text-foreground leading-[1.15]">
                Enterprise Healthcare Intelligence & Multi-Facility Operations
              </h1>
              <p className="text-sm text-muted-foreground leading-relaxed max-w-xl">
                Unified outpatient clinical records, deterministic branch routing, real-time patient triage, and cryptographic audit trails engineered for modern healthcare networks.
              </p>
            </div>

            <div className="space-y-3 pt-1">
              <div className="flex items-start gap-3.5 p-3.5 rounded-2xl bg-card/60 backdrop-blur-md border border-border/80 shadow-sm hover:border-primary/40 transition-colors">
                <div className="p-2 rounded-xl bg-primary/10 border border-primary/20 text-primary shrink-0">
                  <Building2 className="w-5 h-5" />
                </div>
                <div className="space-y-0.5">
                  <h3 className="text-xs font-bold text-foreground">
                    Deterministic Multi-Facility Branch Mesh
                  </h3>
                  <p className="text-xs text-muted-foreground leading-relaxed">
                    Automated location resolution and role routing across hospital headquarters, outpatient specialty centers, and diagnostic emergency hubs.
                  </p>
                </div>
              </div>

              <div className="flex items-start gap-3.5 p-3.5 rounded-2xl bg-card/60 backdrop-blur-md border border-border/80 shadow-sm hover:border-primary/40 transition-colors">
                <div className="p-2 rounded-xl bg-teal-500/10 border border-teal-500/20 text-teal-600 dark:text-teal-400 shrink-0">
                  <Stethoscope className="w-5 h-5" />
                </div>
                <div className="space-y-0.5">
                  <h3 className="text-xs font-bold text-foreground">
                    Role-Tailored Clinical Workspaces
                  </h3>
                  <p className="text-xs text-muted-foreground leading-relaxed">
                    Dedicated interfaces for attending doctors, triage nurses, surgeons, and reception coordinators with contextual patient charting.
                  </p>
                </div>
              </div>

              <div className="flex items-start gap-3.5 p-3.5 rounded-2xl bg-card/60 backdrop-blur-md border border-border/80 shadow-sm hover:border-primary/40 transition-colors">
                <div className="p-2 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 dark:text-emerald-400 shrink-0">
                  <Lock className="w-5 h-5" />
                </div>
                <div className="space-y-0.5">
                  <h3 className="text-xs font-bold text-foreground">
                    Tamper-Evident Cryptographic Ledger
                  </h3>
                  <p className="text-xs text-muted-foreground leading-relaxed">
                    Every access, clinical note, and prescription is cryptographically chained with SHA-256 signatures for immutable regulatory compliance.
                  </p>
                </div>
              </div>
            </div>

            {/* Compliance & Security Badges */}
            <div className="pt-2 flex flex-wrap items-center gap-2">
              <span className="text-[11px] font-medium px-2.5 py-1 rounded-lg bg-muted/60 border border-border text-muted-foreground flex items-center gap-1.5">
                <CheckCircle2 className="w-3.5 h-3.5 text-primary" />
                HIPAA § 164.312 Compliant
              </span>
              <span className="text-[11px] font-medium px-2.5 py-1 rounded-lg bg-muted/60 border border-border text-muted-foreground flex items-center gap-1.5">
                <CheckCircle2 className="w-3.5 h-3.5 text-primary" />
                256-Bit TLS End-to-End
              </span>
              <span className="text-[11px] font-medium px-2.5 py-1 rounded-lg bg-muted/60 border border-border text-muted-foreground flex items-center gap-1.5">
                <CheckCircle2 className="w-3.5 h-3.5 text-primary" />
                Zero-Trust Multi-Tenancy
              </span>
            </div>
          </div>

          {/* Right Column: Credentials Login Card */}
          <div className="lg:col-span-5 space-y-4">
            <Card className="border-border shadow-xl backdrop-blur-md bg-card/85">
              <CardHeader className="space-y-1 pb-4">
                <CardTitle className="text-base font-semibold flex items-center gap-2">
                  <ShieldCheck className="h-5 w-5 text-primary" />
                  Sign In to Workspace
                </CardTitle>
                <CardDescription className="text-xs">
                  Enter your verified staff or administrative credentials.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <form onSubmit={handleSubmit} className="space-y-4">
                  {error && (
                    <Alert variant="destructive" className="py-2.5 text-xs">
                      <AlertDescription>{error}</AlertDescription>
                    </Alert>
                  )}

                  <div className="space-y-1.5">
                    <Label htmlFor="email" className="text-xs font-medium">
                      Work Email
                    </Label>
                    <div className="relative">
                      <Mail className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                      <Input
                        id="email"
                        type="email"
                        placeholder="staff@curexal.com"
                        value={email}
                        onChange={(e) => {
                          setEmail(e.target.value);
                          setSelectedPersonaId(null);
                        }}
                        required
                        className="pl-9 text-sm"
                        disabled={isLoading}
                      />
                    </div>
                  </div>

                  <div className="space-y-1.5">
                    <div className="flex items-center justify-between">
                      <Label htmlFor="password" className="text-xs font-medium">
                        Password
                      </Label>
                      <button
                        type="button"
                        onClick={handleOpenForgot}
                        className="text-xs font-medium text-primary hover:text-primary/80 transition-colors hover:underline"
                      >
                        Forgot password?
                      </button>
                    </div>
                    <div className="relative">
                      <Lock className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                      <Input
                        id="password"
                        type={showPassword ? "text" : "password"}
                        placeholder="Enter password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                        className="pl-9 pr-9 text-sm"
                        disabled={isLoading}
                      />
                      <button
                        type="button"
                        onClick={() => setShowPassword(!showPassword)}
                        className="absolute right-3 top-3 text-muted-foreground hover:text-foreground transition-colors"
                        tabIndex={-1}
                      >
                        {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                      </button>
                    </div>
                  </div>

                  <Button
                    type="submit"
                    className="w-full bg-primary hover:bg-primary/90 text-primary-foreground font-medium shadow-md shadow-primary/20 transition-all text-xs h-9.5"
                    disabled={isLoading}
                  >
                    {isLoading ? (
                      <>
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        Authenticating Credentials...
                      </>
                    ) : (
                      <>
                        Sign In to Secure Workspace
                        <ArrowRight className="ml-2 h-4 w-4" />
                      </>
                    )}
                  </Button>
                </form>

                {/* Set Initial Password Link */}
                <div className="pt-2 text-center border-t border-border/60">
                  <p className="text-xs text-muted-foreground">
                    Have an invitation code or setup token?{" "}
                    <button
                      type="button"
                      onClick={handleOpenSetPassword}
                      className="font-semibold text-primary hover:underline transition-colors"
                    >
                      Set Initial Password
                    </button>
                  </p>
                </div>
              </CardContent>
            </Card>

            <div className="text-center space-y-1">
              <p className="text-[11px] text-muted-foreground flex items-center justify-center gap-1.5">
                <ShieldCheck className="w-3.5 h-3.5 text-primary" />
                Protected by Curexal Zero-Trust RBAC &amp; Deterministic Branch Routing
              </p>
            </div>

            {/* Optional Dev/QA Testing Persona Matrix (Active ONLY when ?dev=true is appended to URL) */}
            {showDevPersonas && (
              <div className="p-3.5 rounded-xl border border-dashed border-amber-500/40 bg-amber-500/5 space-y-3">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-1.5 text-amber-600 dark:text-amber-400 font-semibold text-xs">
                    <Zap className="w-3.5 h-3.5" />
                    <span>QA / Developer Persona Switcher (?dev=true)</span>
                  </div>
                  <Badge variant="outline" className="text-[9px] font-mono border-amber-500/30 text-amber-600 dark:text-amber-400">
                    Dev Mode
                  </Badge>
                </div>
                <div className="grid grid-cols-2 gap-2">
                  {DAY1_PERSONAS.map((persona) => (
                    <button
                      key={persona.id}
                      type="button"
                      onClick={() => handleSelectPersona(persona, false)}
                      className="p-2 rounded-lg border border-border/70 hover:border-primary bg-card/80 text-left transition-all text-[11px]"
                    >
                      <div className="font-semibold text-foreground truncate">{persona.name}</div>
                      <div className="text-[10px] text-muted-foreground truncate">{persona.roleTitle}</div>
                    </button>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </motion.div>

      {/* Branch Selection Dialog (Feature #1 / Feature #3) */}
      <Dialog
        open={isBranchSelectOpen}
        onOpenChange={(open) => {
          if (!isBranchSelectLoading) setIsBranchSelectOpen(open);
        }}
      >
        <DialogContent className="sm:max-w-md border-border bg-card shadow-2xl">
          <DialogHeader>
            <div className="flex items-center gap-2 text-primary mb-1">
              <div className="p-1.5 rounded-lg bg-primary/10 text-primary">
                <Building2 className="h-5 w-5" />
              </div>
              <DialogTitle className="text-base font-semibold">
                Select Operating Facility Location
              </DialogTitle>
            </div>
            <DialogDescription className="text-xs text-muted-foreground leading-relaxed">
              {branchSelectionData?.organization?.name
                ? `You hold credentials across multiple facilities in ${branchSelectionData.organization.name}. Select your current clinical location to enter.`
                : "Select the operational clinic or facility location you are currently attending to."}
            </DialogDescription>
          </DialogHeader>

          {branchSelectError && (
            <Alert variant="destructive" className="py-2 text-xs">
              <AlertDescription>{branchSelectError}</AlertDescription>
            </Alert>
          )}

          <form onSubmit={handleBranchSelectSubmit} className="space-y-4 py-2">
            <div className="space-y-2.5 max-h-72 overflow-y-auto pr-1">
              {branchSelectionData?.assignedBranches?.map((b) => {
                const isSelected = selectedBranchId === b.id;
                const isHq = b.isHeadquarters === true || b.code === "HO-01";

                return (
                  <motion.div
                    key={b.id}
                    whileHover={{ scale: 1.01 }}
                    whileTap={{ scale: 0.99 }}
                    onClick={() => setSelectedBranchId(b.id)}
                    className={`flex items-start justify-between p-3.5 rounded-xl border cursor-pointer transition-all ${
                      isSelected
                        ? "border-primary bg-primary/5 shadow-md shadow-primary/10 ring-1 ring-primary"
                        : "border-border/80 hover:border-border hover:bg-muted/30 bg-card"
                    }`}
                  >
                    <div className="flex items-start gap-3">
                      <div
                        className={`p-2 rounded-lg shrink-0 mt-0.5 ${
                          isSelected
                            ? "bg-primary text-primary-foreground"
                            : "bg-muted text-muted-foreground"
                        }`}
                      >
                        {isHq ? <Building2 className="w-4 h-4" /> : <Stethoscope className="w-4 h-4" />}
                      </div>
                      <div className="space-y-1">
                        <div className="flex items-center gap-1.5">
                          <p
                            className={`text-xs font-semibold ${
                              isSelected ? "text-primary" : "text-foreground"
                            }`}
                          >
                            {b.name}
                          </p>
                          {isHq && (
                            <span className="text-[9px] font-semibold bg-amber-500/10 text-amber-600 dark:text-amber-400 px-1.5 py-0.2 rounded border border-amber-500/20">
                              HQ
                            </span>
                          )}
                        </div>
                        <p className="text-[11px] text-muted-foreground font-mono">
                          Code: {b.code}
                          {b.facilityType ? ` • ${b.facilityType.replace(/_/g, " ")}` : ""}
                        </p>
                      </div>
                    </div>
                    <div className="pt-0.5">
                      <div
                        className={`h-4 w-4 rounded-full border flex items-center justify-center transition-all ${
                          isSelected
                            ? "border-primary bg-primary text-primary-foreground"
                            : "border-muted-foreground/40"
                        }`}
                      >
                        {isSelected && <Check className="h-2.5 w-2.5 stroke-[3]" />}
                      </div>
                    </div>
                  </motion.div>
                );
              })}
            </div>

            <DialogFooter className="gap-2 sm:gap-0 pt-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => setIsBranchSelectOpen(false)}
                disabled={isBranchSelectLoading}
                className="text-xs"
              >
                Cancel
              </Button>
              <Button
                type="submit"
                disabled={isBranchSelectLoading || !selectedBranchId}
                className="bg-primary hover:bg-primary/90 text-primary-foreground text-xs font-medium gap-1.5"
              >
                {isBranchSelectLoading ? (
                  <>
                    <Loader2 className="h-4 w-4 animate-spin" />
                    Connecting Facility...
                  </>
                ) : (
                  <>
                    Enter Clinical Workspace
                    <ArrowRight className="h-3.5 w-3.5" />
                  </>
                )}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Unassigned Facility Branch Guidance Dialog (Zero-Trust Feature #1 / Feature #3) */}
      <Dialog open={isUnassignedModalOpen} onOpenChange={setIsUnassignedModalOpen}>
        <DialogContent className="sm:max-w-md border-border bg-card shadow-2xl">
          <DialogHeader>
            <div className="flex items-center gap-2 text-amber-500 mb-1">
              <div className="p-1.5 rounded-lg bg-amber-500/10 text-amber-600 dark:text-amber-400">
                <AlertCircle className="h-5 w-5" />
              </div>
              <DialogTitle className="text-base font-semibold">
                Facility Location Assignment Required
              </DialogTitle>
            </div>
            <DialogDescription className="text-xs text-muted-foreground leading-relaxed">
              {unassignedOrgData?.name
                ? `Your account is active at ${unassignedOrgData.name}, but pending clinical location provisioning.`
                : "Your account is verified, but has not yet been assigned to a clinic or branch location."}
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-3 py-2 text-xs">
            <Alert className="border-amber-500/30 bg-amber-500/10 text-amber-800 dark:text-amber-300 py-3">
              <AlertDescription className="leading-relaxed">
                To comply with Zero-Trust Clinical Governance and HIPAA multi-center audit trails, medical and administrative personnel must be linked to at least one active clinic location before accessing patient records.
              </AlertDescription>
            </Alert>

            <div className="rounded-xl border border-border p-3 space-y-2 bg-secondary/30">
              <p className="font-semibold text-foreground flex items-center gap-1.5">
                <Building2 className="w-3.5 h-3.5 text-primary" />
                Remediation Workflow:
              </p>
              <ul className="list-disc list-inside space-y-1 text-muted-foreground text-[11px] leading-relaxed">
                <li>Contact your Chief Medical Officer, Clinic Administrator, or IT Lead.</li>
                <li>Request that they assign you to your operating facility branch in the Staff Directory.</li>
                <li>Once assigned, click <strong>"Check Assignment Status"</strong> below to enter directly.</li>
              </ul>
            </div>
          </div>

          <DialogFooter className="gap-2 sm:gap-0 pt-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => {
                setIsUnassignedModalOpen(false);
                setPassword("");
              }}
              className="text-xs"
            >
              Sign In as Different User
            </Button>
            <Button
              type="button"
              onClick={(e) => {
                setIsUnassignedModalOpen(false);
                handleSubmit(e);
              }}
              className="bg-primary hover:bg-primary/90 text-primary-foreground text-xs font-medium gap-1.5"
            >
              <RefreshCw className="h-3.5 w-3.5" />
              Check Assignment Status
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Forgot Password Dialog */}
      <Dialog open={isForgotOpen} onOpenChange={setIsForgotOpen}>
        <DialogContent className="sm:max-w-md border-border bg-card shadow-2xl">
          <DialogHeader>
            <div className="flex items-center gap-2 text-primary mb-1">
              <KeyRound className="h-5 w-5" />
              <DialogTitle className="text-base font-semibold">
                Reset Account Password
              </DialogTitle>
            </div>
            <DialogDescription className="text-xs text-muted-foreground">
              Enter your registered administrator email to receive password recovery instructions.
            </DialogDescription>
          </DialogHeader>

          {forgotSuccess ? (
            <div className="space-y-4 py-2">
              <Alert className="border-emerald-500/30 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 py-3 text-xs">
                <CheckCircle2 className="h-4 w-4 mr-2 inline" />
                <AlertDescription>{forgotSuccess}</AlertDescription>
              </Alert>
              <DialogFooter>
                <Button
                  type="button"
                  onClick={() => setIsForgotOpen(false)}
                  className="w-full bg-primary hover:bg-primary/90 text-primary-foreground text-xs"
                >
                  Back to Sign In
                </Button>
              </DialogFooter>
            </div>
          ) : (
            <form onSubmit={handleForgotSubmit} className="space-y-4 py-2">
              {forgotError && (
                <Alert variant="destructive" className="py-2 text-xs">
                  <AlertDescription>{forgotError}</AlertDescription>
                </Alert>
              )}

              <div className="space-y-1.5">
                <Label htmlFor="forgot-email" className="text-xs font-medium">
                  Administrator Email
                </Label>
                <div className="relative">
                  <Mail className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                  <Input
                    id="forgot-email"
                    type="email"
                    placeholder="name@example.com"
                    value={forgotEmail}
                    onChange={(e) => setForgotEmail(e.target.value)}
                    required
                    className="pl-9 text-sm"
                    disabled={isForgotLoading}
                  />
                </div>
              </div>

              <DialogFooter className="gap-2 sm:gap-0">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setIsForgotOpen(false)}
                  disabled={isForgotLoading}
                  className="text-xs"
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  disabled={isForgotLoading}
                  className="bg-primary hover:bg-primary/90 text-primary-foreground text-xs"
                >
                  {isForgotLoading ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Sending Instructions...
                    </>
                  ) : (
                    "Send Reset Link"
                  )}
                </Button>
              </DialogFooter>
            </form>
          )}
        </DialogContent>
      </Dialog>

      {/* Set Initial Password Dialog */}
      <Dialog open={isSetOpen} onOpenChange={setIsSetOpen}>
        <DialogContent className="sm:max-w-md border-border bg-card shadow-2xl">
          <DialogHeader>
            <div className="flex items-center gap-2 text-primary mb-1">
              <Sparkles className="h-5 w-5" />
              <DialogTitle className="text-base font-semibold">
                Set Initial Password
              </DialogTitle>
            </div>
            <DialogDescription className="text-xs text-muted-foreground">
              Complete account setup using the 6-character code or token sent in your invitation email.
            </DialogDescription>
          </DialogHeader>

          {setSuccessMsg ? (
            <div className="space-y-4 py-2">
              <Alert className="border-emerald-500/30 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 py-3 text-xs">
                <CheckCircle2 className="h-4 w-4 mr-2 inline" />
                <AlertDescription>{setSuccessMsg}</AlertDescription>
              </Alert>
              <DialogFooter>
                <Button
                  type="button"
                  onClick={() => setIsSetOpen(false)}
                  className="w-full bg-primary hover:bg-primary/90 text-primary-foreground text-xs"
                >
                  Proceed to Sign In
                </Button>
              </DialogFooter>
            </div>
          ) : (
            <form onSubmit={handleSetPasswordSubmit} className="space-y-3.5 py-2">
              {setErrorMsg && (
                <Alert variant="destructive" className="py-2 text-xs">
                  <AlertDescription>{setErrorMsg}</AlertDescription>
                </Alert>
              )}

              <div className="space-y-1.5">
                <Label htmlFor="set-email" className="text-xs font-medium">
                  Account Email
                </Label>
                <div className="relative">
                  <Mail className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                  <Input
                    id="set-email"
                    type="email"
                    placeholder="name@example.com"
                    value={setEmailVal}
                    onChange={(e) => setSetEmailVal(e.target.value)}
                    required
                    className="pl-9 text-sm"
                    disabled={isSetLoading}
                  />
                </div>
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="set-code" className="text-xs font-medium">
                  Invitation Code or Token
                </Label>
                <div className="relative">
                  <KeyRound className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                  <Input
                    id="set-code"
                    type="text"
                    placeholder="Enter code or token"
                    value={setCodeVal}
                    onChange={(e) => setSetCodeVal(e.target.value)}
                    required
                    className="pl-9 text-sm font-mono uppercase"
                    disabled={isSetLoading}
                  />
                </div>
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="set-password" className="text-xs font-medium">
                  New Password (min. 8 characters)
                </Label>
                <div className="relative">
                  <Lock className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                  <Input
                    id="set-password"
                    type={showSetPassword ? "text" : "password"}
                    placeholder="Create strong password"
                    value={setNewPassword}
                    onChange={(e) => setSetNewPassword(e.target.value)}
                    required
                    minLength={8}
                    className="pl-9 pr-9 text-sm"
                    disabled={isSetLoading}
                  />
                  <button
                    type="button"
                    onClick={() => setShowSetPassword(!showSetPassword)}
                    className="absolute right-3 top-3 text-muted-foreground hover:text-foreground"
                  >
                    {showSetPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                  </button>
                </div>
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="set-confirm-password" className="text-xs font-medium">
                  Confirm Password
                </Label>
                <div className="relative">
                  <Lock className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                  <Input
                    id="set-confirm-password"
                    type={showSetPassword ? "text" : "password"}
                    placeholder="Re-enter password"
                    value={setConfirmPassword}
                    onChange={(e) => setSetConfirmPassword(e.target.value)}
                    required
                    minLength={8}
                    className="pl-9 text-sm"
                    disabled={isSetLoading}
                  />
                </div>
              </div>

              <DialogFooter className="gap-2 sm:gap-0 pt-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setIsSetOpen(false)}
                  disabled={isSetLoading}
                  className="text-xs"
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  disabled={isSetLoading}
                  className="bg-primary hover:bg-primary/90 text-primary-foreground text-xs"
                >
                  {isSetLoading ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Saving Password...
                    </>
                  ) : (
                    "Set Password & Activate"
                  )}
                </Button>
              </DialogFooter>
            </form>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}
