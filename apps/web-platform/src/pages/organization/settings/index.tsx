import React, { useState, useEffect } from "react";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { organizationService } from "@/api/services/organization.service";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import {
  Building2,
  Save,
  Globe,
  Clock,
  Shield,
  FileText,
  DollarSign,
  Phone,
  Mail,
  Loader2,
} from "lucide-react";

export default function OrganizationSettingsPage() {
  const { data: bootstrap, refetch } = useBootstrap();

  const org = bootstrap?.organization;
  const workspace = bootstrap?.workspace;

  const [orgName, setOrgName] = useState("");
  const [legalName, setLegalName] = useState("");
  const [taxId, setTaxId] = useState("");
  const [currency, setCurrency] = useState("NGN");
  const [timezone, setTimezone] = useState("Africa/Lagos");
  const [primaryPhone, setPrimaryPhone] = useState("");
  const [supportEmail, setSupportEmail] = useState("");
  const [headquartersAddress, setHeadquartersAddress] = useState("");
  const [version, setVersion] = useState<number>(0);

  const [isLoadingProfile, setIsLoadingProfile] = useState(true);
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    let isMounted = true;
    async function loadProfile() {
      try {
        const profile = await organizationService.getProfile();
        if (isMounted && profile) {
          const profileSettings = (typeof profile.settings === "object" && profile.settings !== null) ? (profile.settings as Record<string, any>) : {};
          setOrgName(profile.name || org?.name || "");
          setLegalName(profile.registrationNumber || profile.legalName || profileSettings.cacNumber || org?.name || "");
          setTaxId(profile.taxId || profileSettings.tinNumber || profileSettings.taxNumber || "");
          setCurrency(profileSettings.currency || profile.currency || workspace?.currency || "NGN");
          setTimezone(profileSettings.timezone || profile.timezone || workspace?.timezone || "Africa/Lagos");
          setPrimaryPhone(profile.phone || profileSettings.supportPhone || "");
          setSupportEmail(profile.email || profileSettings.supportEmail || "");
          setHeadquartersAddress(profile.address || profileSettings.address || "");
          setVersion(profile.version || 0);
        } else if (isMounted) {
          setOrgName(org?.name || "");
          setCurrency(workspace?.currency || "NGN");
          setTimezone(workspace?.timezone || "Africa/Lagos");
        }
      } catch (err) {
        console.error("Failed to load organization profile:", err);
        if (isMounted) {
          setOrgName(org?.name || "");
          setCurrency(workspace?.currency || "NGN");
          setTimezone(workspace?.timezone || "Africa/Lagos");
        }
      } finally {
        if (isMounted) setIsLoadingProfile(false);
      }
    }
    loadProfile();
    return () => {
      isMounted = false;
    };
  }, [org?.name, workspace?.currency, workspace?.timezone]);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSaving(true);
    try {
      const updated = await organizationService.updateProfile({
        name: orgName,
        legalName,
        registrationNumber: legalName,
        taxId,
        currency,
        timezone,
        phone: primaryPhone,
        email: supportEmail,
        address: headquartersAddress,
        version,
      });
      if (updated?.version) {
        setVersion(updated.version);
      }
      toast.success("Organization Profile Updated!", {
        description: "Legal, regional, and financial preferences saved.",
      });
      if (refetch) refetch();
    } catch (err: any) {
      const errMsg = err.response?.data?.message || err.message || "Failed to update organization profile";
      toast.error("Failed to save settings: " + errMsg);
    } finally {
      setIsSaving(false);
    }
  };

  if (isLoadingProfile) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-6">
        <div>
          <div className="flex items-center gap-2.5 mb-1">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              Organization Profile & Regional Preferences
            </h1>
            <Badge variant="outline" className="border-primary/40 text-primary bg-primary/10 text-[10px] font-mono">
              Corporate Governance
            </Badge>
          </div>
          <p className="text-xs text-muted-foreground">
            Configure legal entity details, primary operating currencies, regional timezones, and support channels.
          </p>
        </div>

        <Button
          type="button"
          size="sm"
          onClick={handleSave}
          disabled={isSaving}
          className="text-xs h-9 gap-1.5 bg-primary text-primary-foreground shadow"
        >
          {isSaving ? <Loader2 className="w-3.5 h-3.5 animate-spin" /> : <Save className="w-3.5 h-3.5" />}
          {isSaving ? "Saving..." : "Save Settings"}
        </Button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Legal & Corporate Details */}
        <Card className="border-border shadow-sm">
          <CardHeader className="pb-3">
            <CardTitle className="text-sm font-bold flex items-center gap-2">
              <Building2 className="w-4 h-4 text-primary" />
              Corporate & Legal Entity
            </CardTitle>
            <CardDescription className="text-xs">
              Official company name and tax identification for invoicing and regulatory compliance.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-1.5">
              <Label className="text-xs font-medium">Display Organization Name</Label>
              <Input value={orgName} onChange={(e) => setOrgName(e.target.value)} className="text-xs h-9" />
            </div>

            <div className="space-y-1.5">
              <Label className="text-xs font-medium">Registered Legal Entity Name</Label>
              <Input value={legalName} onChange={(e) => setLegalName(e.target.value)} className="text-xs h-9" />
            </div>

            <div className="space-y-1.5">
              <Label className="text-xs font-medium">Tax Identification Number (TIN / CAC)</Label>
              <Input value={taxId} onChange={(e) => setTaxId(e.target.value)} className="text-xs h-9" />
            </div>

            <div className="space-y-1.5">
              <Label className="text-xs font-medium">Corporate Headquarters Address</Label>
              <Input
                value={headquartersAddress}
                onChange={(e) => setHeadquartersAddress(e.target.value)}
                className="text-xs h-9"
              />
            </div>
          </CardContent>
        </Card>

        {/* Regional & Financial Preferences */}
        <Card className="border-border shadow-sm">
          <CardHeader className="pb-3">
            <CardTitle className="text-sm font-bold flex items-center gap-2">
              <Globe className="w-4 h-4 text-primary" />
              Regional & Financial Preferences
            </CardTitle>
            <CardDescription className="text-xs">
              Set default ledger currency, clinical timezone, and emergency support dispatches.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <Label className="text-xs font-medium">Primary Ledger Currency</Label>
                <Input value={currency} onChange={(e) => setCurrency(e.target.value)} className="text-xs h-9 font-mono" />
              </div>

              <div className="space-y-1.5">
                <Label className="text-xs font-medium">Operational Timezone</Label>
                <Input value={timezone} onChange={(e) => setTimezone(e.target.value)} className="text-xs h-9 font-mono" />
              </div>
            </div>

            <div className="space-y-1.5">
              <Label className="text-xs font-medium">Emergency Care Support Phone</Label>
              <Input value={primaryPhone} onChange={(e) => setPrimaryPhone(e.target.value)} className="text-xs h-9" />
            </div>

            <div className="space-y-1.5">
              <Label className="text-xs font-medium">Support & Inquiries Email</Label>
              <Input value={supportEmail} onChange={(e) => setSupportEmail(e.target.value)} className="text-xs h-9" />
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
