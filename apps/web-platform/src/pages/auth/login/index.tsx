import React, { useState, useEffect } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { authClient } from "@/lib/auth-client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";
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
} from "lucide-react";
import { CurexalLogoSymbol } from "@/components/brand/curexal-logo";

function resolveAuthorizedDestination(sessionData: any, returnTo?: string | null): string {
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
    Boolean(effectiveBootstrap?.organization?.id) ||
    effectiveBootstrap?.contexts?.current === "organization" ||
    user?.role === "owner" ||
    user?.role === "org_admin" ||
    user?.role === "org_regional_manager" ||
    Boolean(user?.organizationId)
  );

  const isWorkspaceAuthorized = Boolean(
    isOrgAuthorized ||
    Boolean(effectiveBootstrap?.workspace?.id) ||
    effectiveBootstrap?.contexts?.current === "workspace" ||
    Boolean(user?.activeTenantId) ||
    Boolean(user?.workspaceId)
  );

  // Validate returnTo if present (prevent open redirect & privilege escalation)
  if (returnTo && returnTo.startsWith("/") && !returnTo.startsWith("//")) {
    if (returnTo.startsWith("/platform/")) {
      if (isPlatformStaff) return returnTo;
    } else if (returnTo.startsWith("/organization/")) {
      if (isOrgAuthorized) return returnTo;
    } else if (returnTo.startsWith("/workspace/")) {
      if (isWorkspaceAuthorized) return returnTo;
    }
  }

  // Canonical destination resolution
  if (isPlatformStaff) return "/platform/dashboard";
  if (isOrgAuthorized) return "/organization/dashboard";
  if (isWorkspaceAuthorized) return "/workspace/dashboard";

  return "/organization/dashboard";
}

export default function LoginPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const returnTo = searchParams.get("returnTo") || searchParams.get("redirect");
  const { data: session, isPending } = authClient.useSession();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

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

  // If already logged in, redirect to authorized dashboard
  useEffect(() => {
    if (!isPending && session?.user) {
      const destination = resolveAuthorizedDestination(session, returnTo);
      navigate(destination, { replace: true });
    }
  }, [session, isPending, returnTo, navigate]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || !password) {
      setError("Please enter both email and password.");
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const updatedSession = await authClient.signIn({ email, password });
      const destination = resolveAuthorizedDestination(updatedSession || session, returnTo);
      navigate(destination, { replace: true });
    } catch (err: any) {
      const msg =
        err?.response?.data?.message ||
        err?.message ||
        "Authentication failed. Please check your credentials.";
      setError(msg);
    } finally {
      setIsLoading(false);
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

  return (
    <div className="relative min-h-screen w-full flex items-center justify-center bg-background px-4 dot-grid">
      <div className="w-full max-w-md space-y-6">
        {/* Brand Header */}
        <div className="flex flex-col items-center text-center space-y-2">
          <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-slate-900 shadow-lg shadow-primary/20">
            <CurexalLogoSymbol className="w-10 h-10" />
          </div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            Curexal Platform
          </h1>
          <p className="text-sm text-muted-foreground">
            Global Control Center & Super Admin Console
          </p>
        </div>

        {/* Login Card */}
        <Card className="border-border shadow-card card-enterprise">
          <CardHeader className="space-y-1 pb-4">
            <CardTitle className="text-lg font-semibold flex items-center gap-2">
              <ShieldCheck className="h-5 w-5 text-primary" />
              Sign in to Console
            </CardTitle>
            <CardDescription className="text-xs">
              Enter your platform administrative credentials to proceed.
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
                  Administrator Email
                </Label>
                <div className="relative">
                  <Mail className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                  <Input
                    id="email"
                    type="email"
                    placeholder="Enter administrator email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
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
                    type="password"
                    placeholder="Enter account password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    className="pl-9 text-sm"
                    disabled={isLoading}
                  />
                </div>
              </div>

              <Button
                type="submit"
                className="w-full bg-primary hover:bg-primary/90 text-primary-foreground font-medium shadow-sm transition-all"
                disabled={isLoading}
              >
                {isLoading ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Authenticating...
                  </>
                ) : (
                  <>
                    Sign In
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </>
                )}
              </Button>
            </form>

            {/* Set Initial Password Prompt */}
            <div className="pt-2 text-center border-t border-border/60">
              <p className="text-xs text-muted-foreground">
                Have an invitation code or token?{" "}
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

        {/* Security Notice */}
        <p className="text-center text-xs text-muted-foreground">
          Protected by Curexal Zero-Trust RBAC and Multi-Tenant Isolation.
        </p>
      </div>

      {/* Forgot Password Dialog */}
      <Dialog open={isForgotOpen} onOpenChange={setIsForgotOpen}>
        <DialogContent className="sm:max-w-md border-border bg-card">
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

      {/* Set Password Dialog */}
      <Dialog open={isSetOpen} onOpenChange={setIsSetOpen}>
        <DialogContent className="sm:max-w-md border-border bg-card">
          <DialogHeader>
            <div className="flex items-center gap-2 text-primary mb-1">
              <Sparkles className="h-5 w-5" />
              <DialogTitle className="text-base font-semibold">
                Set Initial Password
              </DialogTitle>
            </div>
            <DialogDescription className="text-xs text-muted-foreground">
              Complete your account setup using the 6-character code or token sent in your invitation email.
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
                    placeholder="Enter 6-character code or token"
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

