import React, { useState, useEffect } from "react";
import {
  Settings,
  Shield,
  Lock,
  Globe,
  AlertTriangle,
  CheckCircle2,
  Save,
  Key,
} from "lucide-react";
import {
  usePlatformConfig,
  useUpdatePlatformConfig,
  useSecurityPolicy,
  useUpdateSecurityPolicy,
} from "@/api/hooks/use-config";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Switch } from "@/components/ui/switch";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import type { PlatformGeneralSettings, IdentitySecurityPolicy } from "@/api/contracts";
import { toast } from "sonner";

export default function SettingsPage() {
  const { data: config, isLoading: configLoading, refetch: refetchConfig } = usePlatformConfig();
  const { data: policy, isLoading: policyLoading, refetch: refetchPolicy } = useSecurityPolicy();

  const updateConfigMutation = useUpdatePlatformConfig();
  const updatePolicyMutation = useUpdateSecurityPolicy();

  // General Settings State
  const [platformName, setPlatformName] = useState("");
  const [supportEmail, setSupportEmail] = useState("");
  const [supportPhone, setSupportPhone] = useState("");
  const [defaultCountry, setDefaultCountry] = useState("");
  const [defaultCurrency, setDefaultCurrency] = useState("");
  const [maintenanceMode, setMaintenanceMode] = useState(false);
  const [announcementBanner, setAnnouncementBanner] = useState("");

  // Security Policy State
  const [minPasswordLength, setMinPasswordLength] = useState(8);
  const [passwordRequireUppercase, setPasswordRequireUppercase] = useState(true);
  const [passwordRequireNumber, setPasswordRequireNumber] = useState(true);
  const [passwordRequireSymbol, setPasswordRequireSymbol] = useState(false);
  const [emailVerificationRequired, setEmailVerificationRequired] = useState(true);
  const [maxFailedLoginAttempts, setMaxFailedLoginAttempts] = useState(5);
  const [accountLockoutDurationMinutes, setAccountLockoutDurationMinutes] = useState(15);
  const [sessionMaxDurationHours, setSessionMaxDurationHours] = useState(24);
  const [refreshTokenDurationDays, setRefreshTokenDurationDays] = useState(30);
  const [maxActiveSessions, setMaxActiveSessions] = useState(5);

  useEffect(() => {
    if (config) {
      setPlatformName(config.platformName || "");
      setSupportEmail(config.supportEmail || "");
      setSupportPhone(config.supportPhone || "");
      setDefaultCountry(config.defaultCountry || "");
      setDefaultCurrency(config.defaultCurrency || "");
      setMaintenanceMode(config.maintenanceMode || false);
      setAnnouncementBanner(config.announcementBanner || "");
    }
  }, [config]);

  useEffect(() => {
    if (policy) {
      setMinPasswordLength(policy.minPasswordLength || 8);
      setPasswordRequireUppercase(policy.passwordRequireUppercase ?? true);
      setPasswordRequireNumber(policy.passwordRequireNumber ?? true);
      setPasswordRequireSymbol(policy.passwordRequireSymbol ?? false);
      setEmailVerificationRequired(policy.emailVerificationRequired ?? true);
      setMaxFailedLoginAttempts(policy.maxFailedLoginAttempts || 5);
      setAccountLockoutDurationMinutes(policy.accountLockoutDurationMinutes || 15);
      setSessionMaxDurationHours(policy.sessionMaxDurationHours || 24);
      setRefreshTokenDurationDays(policy.refreshTokenDurationDays || 30);
      setMaxActiveSessions(policy.maxActiveSessions || 5);
    }
  }, [policy]);

  const handleSaveConfig = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!config) return;

    try {
      await updateConfigMutation.mutateAsync({
        ...config,
        platformName,
        supportEmail,
        supportPhone,
        defaultCountry,
        defaultCurrency,
        maintenanceMode,
        announcementBanner: announcementBanner || undefined,
        version: config.version,
      });
      toast.success("Platform general configuration saved!");
      refetchConfig();
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Failed to update configuration");
    }
  };

  const handleSavePolicy = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!policy) return;

    try {
      await updatePolicyMutation.mutateAsync({
        ...policy,
        minPasswordLength,
        passwordRequireUppercase,
        passwordRequireNumber,
        passwordRequireSymbol,
        emailVerificationRequired,
        maxFailedLoginAttempts,
        accountLockoutDurationMinutes,
        sessionMaxDurationHours,
        refreshTokenDurationDays,
        maxActiveSessions,
        version: policy.version,
      });
      toast.success("Identity security governance policy saved!");
      refetchPolicy();
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Failed to update security policy");
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-foreground">
          Platform Console Settings & Security
        </h1>
        <p className="text-sm text-muted-foreground mt-1">
          Global platform parameters, operational flags, and identity governance policies.
        </p>
      </div>

      <Tabs defaultValue="general" className="space-y-4">
        <TabsList className="bg-muted/60 p-1 border border-border">
          <TabsTrigger value="general" className="text-xs gap-1.5">
            <Globe className="h-3.5 w-3.5 text-primary" /> General Platform Configuration
          </TabsTrigger>
          <TabsTrigger value="security" className="text-xs gap-1.5">
            <Shield className="h-3.5 w-3.5 text-emerald-500" /> Identity Security Policy
          </TabsTrigger>
        </TabsList>

        {/* TAB 1: GENERAL CONFIGURATION */}
        <TabsContent value="general" className="space-y-4">
          <form onSubmit={handleSaveConfig}>
            <Card className="card-enterprise">
              <CardHeader className="pb-4">
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle className="text-sm font-semibold">General Platform Parameters</CardTitle>
                    <CardDescription className="text-xs">
                      Controls public defaults and platform contact channels across all tenant instances.
                    </CardDescription>
                  </div>
                  <span className="font-mono text-xs text-muted-foreground">v{config?.version || 1}</span>
                </div>
              </CardHeader>

              <CardContent className="space-y-4 text-xs">
                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="space-y-1.5">
                    <Label htmlFor="platName">Platform Name</Label>
                    <Input
                      id="platName"
                      value={platformName}
                      onChange={(e) => setPlatformName(e.target.value)}
                      required
                      className="text-xs"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="platEmail">Global Support Email</Label>
                    <Input
                      id="platEmail"
                      type="email"
                      value={supportEmail}
                      onChange={(e) => setSupportEmail(e.target.value)}
                      required
                      className="text-xs"
                    />
                  </div>
                </div>

                <div className="grid gap-4 sm:grid-cols-3">
                  <div className="space-y-1.5">
                    <Label htmlFor="platPhone">Support Phone</Label>
                    <Input
                      id="platPhone"
                      value={supportPhone}
                      onChange={(e) => setSupportPhone(e.target.value)}
                      className="text-xs"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="defCountry">Default Country</Label>
                    <Input
                      id="defCountry"
                      value={defaultCountry}
                      onChange={(e) => setDefaultCountry(e.target.value)}
                      className="text-xs font-mono"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="defCurr">Default Currency</Label>
                    <Input
                      id="defCurr"
                      value={defaultCurrency}
                      onChange={(e) => setDefaultCurrency(e.target.value)}
                      className="text-xs font-mono"
                    />
                  </div>
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="banner">Global Announcement Banner (Optional)</Label>
                  <Input
                    id="banner"
                    placeholder="e.g. Scheduled platform maintenance on Sunday 2 AM UTC"
                    value={announcementBanner}
                    onChange={(e) => setAnnouncementBanner(e.target.value)}
                    className="text-xs"
                  />
                </div>

                <div className="flex items-center justify-between rounded-lg border border-border/80 bg-secondary/30 p-3 mt-4">
                  <div className="space-y-0.5">
                    <Label htmlFor="maintenance" className="font-semibold cursor-pointer">Platform Maintenance Mode</Label>
                    <p className="text-[11px] text-muted-foreground">
                      When enabled, non-admin users will receive 503 Service Unavailable notice.
                    </p>
                  </div>
                  <Switch
                    id="maintenance"
                    checked={maintenanceMode}
                    onCheckedChange={setMaintenanceMode}
                  />
                </div>

                <div className="pt-2 flex justify-end">
                  <Button
                    type="submit"
                    size="sm"
                    disabled={updateConfigMutation.isPending}
                    className="bg-primary text-primary-foreground text-xs gap-1.5 shadow-sm"
                  >
                    <Save className="h-3.5 w-3.5" />
                    {updateConfigMutation.isPending ? "Saving..." : "Save Platform Settings"}
                  </Button>
                </div>
              </CardContent>
            </Card>
          </form>
        </TabsContent>

        {/* TAB 2: IDENTITY SECURITY POLICY */}
        <TabsContent value="security" className="space-y-4">
          <form onSubmit={handleSavePolicy}>
            <Card className="card-enterprise">
              <CardHeader className="pb-4">
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle className="text-sm font-semibold">Identity Security & Lockout Policy</CardTitle>
                    <CardDescription className="text-xs">
                      Enforce global password strength, session lifetimes, and brute-force lockout thresholds.
                    </CardDescription>
                  </div>
                  <span className="font-mono text-xs text-muted-foreground">v{policy?.version || 1}</span>
                </div>
              </CardHeader>

              <CardContent className="space-y-4 text-xs">
                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="space-y-1.5">
                    <Label htmlFor="minPass">Minimum Password Length</Label>
                    <Input
                      id="minPass"
                      type="number"
                      min={8}
                      max={32}
                      value={minPasswordLength}
                      onChange={(e) => setMinPasswordLength(parseInt(e.target.value, 10))}
                      required
                      className="text-xs font-mono"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="maxFailed">Max Failed Login Attempts (Lockout)</Label>
                    <Input
                      id="maxFailed"
                      type="number"
                      min={1}
                      max={20}
                      value={maxFailedLoginAttempts}
                      onChange={(e) => setMaxFailedLoginAttempts(parseInt(e.target.value, 10))}
                      required
                      className="text-xs font-mono"
                    />
                  </div>
                </div>

                <div className="grid gap-4 sm:grid-cols-3">
                  <div className="space-y-1.5">
                    <Label htmlFor="lockoutMin">Lockout Duration (Minutes)</Label>
                    <Input
                      id="lockoutMin"
                      type="number"
                      min={1}
                      max={1440}
                      value={accountLockoutDurationMinutes}
                      onChange={(e) => setAccountLockoutDurationMinutes(parseInt(e.target.value, 10))}
                      required
                      className="text-xs font-mono"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="sessHours">Session Max Duration (Hours)</Label>
                    <Input
                      id="sessHours"
                      type="number"
                      min={1}
                      max={720}
                      value={sessionMaxDurationHours}
                      onChange={(e) => setSessionMaxDurationHours(parseInt(e.target.value, 10))}
                      required
                      className="text-xs font-mono"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="refreshDays">Refresh Token Duration (Days)</Label>
                    <Input
                      id="refreshDays"
                      type="number"
                      min={1}
                      max={365}
                      value={refreshTokenDurationDays}
                      onChange={(e) => setRefreshTokenDurationDays(parseInt(e.target.value, 10))}
                      required
                      className="text-xs font-mono"
                    />
                  </div>
                </div>

                {/* Password Complexity Toggles */}
                <div className="space-y-2 rounded-lg border border-border/80 bg-secondary/30 p-3 mt-4">
                  <div className="flex items-center justify-between pb-2 border-b border-border/40">
                    <Label htmlFor="reqUpper" className="cursor-pointer">Require Uppercase Character (A-Z)</Label>
                    <Switch
                      id="reqUpper"
                      checked={passwordRequireUppercase}
                      onCheckedChange={setPasswordRequireUppercase}
                    />
                  </div>
                  <div className="flex items-center justify-between pb-2 border-b border-border/40">
                    <Label htmlFor="reqNum" className="cursor-pointer">Require Numeric Digit (0-9)</Label>
                    <Switch
                      id="reqNum"
                      checked={passwordRequireNumber}
                      onCheckedChange={setPasswordRequireNumber}
                    />
                  </div>
                  <div className="flex items-center justify-between pb-2 border-b border-border/40">
                    <Label htmlFor="reqSym" className="cursor-pointer">Require Special Symbol (!@#$)</Label>
                    <Switch
                      id="reqSym"
                      checked={passwordRequireSymbol}
                      onCheckedChange={setPasswordRequireSymbol}
                    />
                  </div>
                  <div className="flex items-center justify-between">
                    <Label htmlFor="reqEmail" className="cursor-pointer">Require Email Verification Before Login</Label>
                    <Switch
                      id="reqEmail"
                      checked={emailVerificationRequired}
                      onCheckedChange={setEmailVerificationRequired}
                    />
                  </div>
                </div>

                <div className="pt-2 flex justify-end">
                  <Button
                    type="submit"
                    size="sm"
                    disabled={updatePolicyMutation.isPending}
                    className="bg-primary text-primary-foreground text-xs gap-1.5 shadow-sm"
                  >
                    <Save className="h-3.5 w-3.5" />
                    {updatePolicyMutation.isPending ? "Saving..." : "Save Security Policy"}
                  </Button>
                </div>
              </CardContent>
            </Card>
          </form>
        </TabsContent>
      </Tabs>
    </div>
  );
}
