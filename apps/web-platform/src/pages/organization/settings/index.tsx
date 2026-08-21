import React, { useState } from "react";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
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
} from "lucide-react";

export default function OrganizationSettingsPage() {
  const { data: bootstrap, refetch } = useBootstrap();

  const org = bootstrap?.organization;
  const workspace = bootstrap?.workspace;

  const [orgName, setOrgName] = useState(org?.name || "Curexal Health Network");
  const [legalName, setLegalName] = useState("Curexal Healthcare Services Ltd.");
  const [taxId, setTaxId] = useState("TIN-9948201-001");
  const [currency, setCurrency] = useState(workspace?.currency || "NGN");
  const [timezone, setTimezone] = useState(workspace?.timezone || "Africa/Lagos");
  const [primaryPhone, setPrimaryPhone] = useState("+234 803 123 4567");
  const [supportEmail, setSupportEmail] = useState("care@curexalhealth.com");
  const [headquartersAddress, setHeadquartersAddress] = useState("Plot 14, Commercial District, Victoria Island, Lagos");

  const [isSaving, setIsSaving] = useState(false);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSaving(true);
    await new Promise((r) => setTimeout(r, 600));
    setIsSaving(false);
    toast.success("Organization Profile Updated!", {
      description: "Legal, regional, and financial preferences saved.",
    });
    if (refetch) refetch();
  };

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
          <Save className="w-3.5 h-3.5" />
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
              <Label className="text-xs font-medium">Corporate Tax Identification Number (TIN / VAT)</Label>
              <Input value={taxId} onChange={(e) => setTaxId(e.target.value)} className="text-xs h-9 font-mono" />
            </div>

            <div className="space-y-1.5">
              <Label className="text-xs font-medium">Headquarters Address</Label>
              <Input value={headquartersAddress} onChange={(e) => setHeadquartersAddress(e.target.value)} className="text-xs h-9" />
            </div>
          </CardContent>
        </Card>

        {/* Financial & Regional Defaults */}
        <Card className="border-border shadow-sm">
          <CardHeader className="pb-3">
            <CardTitle className="text-sm font-bold flex items-center gap-2">
              <Globe className="w-4 h-4 text-primary" />
              Regional & Financial Defaults
            </CardTitle>
            <CardDescription className="text-xs">
              Default currency for cashier billing and timezone for clinical timestamps.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-2 gap-3">
              <div className="space-y-1.5">
                <Label className="text-xs font-medium">Primary Currency</Label>
                <select
                  value={currency}
                  onChange={(e) => setCurrency(e.target.value)}
                  className="w-full h-9 rounded-md border border-input bg-background px-3 text-xs focus:outline-none focus:ring-1 focus:ring-primary font-mono"
                >
                  <option value="NGN">NGN (Nigerian Naira - ₦)</option>
                  <option value="USD">USD (US Dollar - $)</option>
                  <option value="GBP">GBP (British Pound - £)</option>
                  <option value="EUR">EUR (Euro - €)</option>
                  <option value="GHS">GHS (Ghana Cedi - ₵)</option>
                  <option value="KES">KES (Kenyan Shilling - KSh)</option>
                </select>
              </div>

              <div className="space-y-1.5">
                <Label className="text-xs font-medium">Default Timezone</Label>
                <select
                  value={timezone}
                  onChange={(e) => setTimezone(e.target.value)}
                  className="w-full h-9 rounded-md border border-input bg-background px-3 text-xs focus:outline-none focus:ring-1 focus:ring-primary font-mono"
                >
                  <option value="Africa/Lagos">Africa/Lagos (GMT+1)</option>
                  <option value="Africa/Accra">Africa/Accra (GMT+0)</option>
                  <option value="Africa/Nairobi">Africa/Nairobi (GMT+3)</option>
                  <option value="Africa/Johannesburg">Africa/Johannesburg (GMT+2)</option>
                  <option value="Europe/London">Europe/London (GMT+0/BST)</option>
                </select>
              </div>
            </div>

            <div className="space-y-1.5">
              <Label className="text-xs font-medium">Support Contact Email</Label>
              <Input value={supportEmail} onChange={(e) => setSupportEmail(e.target.value)} className="text-xs h-9" />
            </div>

            <div className="space-y-1.5">
              <Label className="text-xs font-medium">Central Reception Phone</Label>
              <Input value={primaryPhone} onChange={(e) => setPrimaryPhone(e.target.value)} className="text-xs h-9" />
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
