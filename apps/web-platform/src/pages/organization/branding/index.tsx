import React, { useState } from "react";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { useBrandTheme } from "@/lib/theme/brand-theme-provider";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import {
  Palette,
  Image as ImageIcon,
  Type,
  Globe,
  FileCheck,
  CheckCircle2,
  Sparkles,
  Shield,
  RotateCcw,
  Save,
  Eye,
  Sliders,
  Award,
} from "lucide-react";

// Curated clinical brand color presets
const colorPresets = [
  { name: "Precision Teal (LIS)", hex: "#0d9488", desc: "Diagnostic & Laboratory" },
  { name: "Medical Cyan (EMR)", hex: "#0284c7", desc: "Clinical & Outpatient" },
  { name: "Corporate Indigo", hex: "#4f46e5", desc: "Multi-facility Enterprise" },
  { name: "Dispensary Mint", hex: "#059669", desc: "Pharmacy & Inventory" },
  { name: "Royal Violet (HIS)", hex: "#7c3aed", desc: "Inpatient Hospital" },
  { name: "Radiology Amber", hex: "#d97706", desc: "Imaging & PACS" },
  { name: "Ruby Rose", hex: "#e11d48", desc: "Emergency & Cardiology" },
  { name: "Curexal Slate", hex: "#0f172a", desc: "Executive Dark" },
];

const fontPresets = [
  { name: "Plus Jakarta Sans", family: "Plus Jakarta Sans", style: "Modern Clinical" },
  { name: "Outfit", family: "Outfit", style: "Geometric Enterprise" },
  { name: "Inter", family: "Inter", style: "Clean Technical" },
  { name: "Roboto", family: "Roboto", style: "Universal Neutral" },
];

export default function OrganizationBrandingPage() {
  const { data: bootstrap, refetch } = useBootstrap();
  const { branding } = useBrandTheme();

  const orgName = bootstrap?.organization?.name || "Healthcare Organization";
  const orgPlan = bootstrap?.organization?.subscription || "pro";

  // Form State
  const [logoUrl, setLogoUrl] = useState(branding?.logoUrl || "");
  const [faviconUrl, setFaviconUrl] = useState(branding?.faviconUrl || "");
  const [primaryColor, setPrimaryColor] = useState(branding?.primaryColor || "#0d9488");
  const [secondaryColor, setSecondaryColor] = useState(branding?.secondaryColor || "#0f172a");
  const [fontFamily, setFontFamily] = useState(branding?.fontFamily || "Plus Jakarta Sans");
  const [borderRadius, setBorderRadius] = useState(branding?.borderRadius || "0.5rem");
  const [customDomain, setCustomDomain] = useState(branding?.customDomain || "");
  const [headerText, setHeaderText] = useState("ISO 15189 Accredited Diagnostic Facility");
  const [footerDisclaimer, setFooterDisclaimer] = useState("Medical Laboratory & Imaging results are electronically validated and tamper-evident.");
  
  const [isSaving, setIsSaving] = useState(false);
  const [previewTab, setPreviewTab] = useState<"ui" | "report">("ui");

  // Save changes handler
  const handleSaveBranding = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSaving(true);

    try {
      // In production: PUT /api/v1/organization/branding
      // Live apply to current DOM document
      const root = document.documentElement;
      root.style.setProperty("--brand-primary", primaryColor);
      root.style.setProperty("--font-sans", `'${fontFamily}', sans-serif`);
      root.style.setProperty("--radius", borderRadius);

      // Simulate API latency
      await new Promise((resolve) => setTimeout(resolve, 600));

      toast.success("Organization Branding Updated Successfully!", {
        description: "Your custom brand colors, typography, and assets are now live.",
      });

      if (refetch) refetch();
    } catch (err: any) {
      toast.error("Failed to update branding: " + (err.message || "Network error"));
    } finally {
      setIsSaving(false);
    }
  };

  const handleResetDefaults = () => {
    setPrimaryColor("#0284c7");
    setSecondaryColor("#0f172a");
    setFontFamily("Outfit");
    setBorderRadius("0.5rem");
    setLogoUrl("");
    toast.info("Reset to Curexal default palette");
  };

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-6">
        <div>
          <div className="flex items-center gap-2.5 mb-1">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              Branding & Multi-Tenant Customization
            </h1>
            <Badge variant="outline" className="border-primary/40 text-primary bg-primary/10 uppercase tracking-wider text-[10px] font-mono font-bold">
              {orgPlan} Plan
            </Badge>
          </div>
          <p className="text-xs text-muted-foreground">
            Personalize your Organization Console, Facility Workspaces, and Patient Diagnostic Reports.
          </p>
        </div>

        <div className="flex items-center gap-2.5">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={handleResetDefaults}
            className="text-xs h-9 gap-1.5"
          >
            <RotateCcw className="w-3.5 h-3.5" />
            Reset Defaults
          </Button>
          <Button
            type="button"
            size="sm"
            onClick={handleSaveBranding}
            disabled={isSaving}
            style={{ backgroundColor: primaryColor }}
            className="text-xs h-9 gap-1.5 text-white shadow-md hover:opacity-90 transition-opacity"
          >
            {isSaving ? (
              <div className="w-3.5 h-3.5 border-2 border-white border-t-transparent animate-spin rounded-full" />
            ) : (
              <Save className="w-3.5 h-3.5" />
            )}
            Save & Publish Brand
          </Button>
        </div>
      </div>

      {/* Main Grid: Controls on Left (5 cols), Live Preview on Right (7 cols) */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
        {/* Left Column: Brand Configuration Controls */}
        <div className="lg:col-span-5 space-y-6">
          {/* Logo & Visual Identity */}
          <Card className="border-border shadow-sm">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-semibold flex items-center gap-2">
                <ImageIcon className="w-4 h-4 text-primary" />
                Logo & Visual Assets
              </CardTitle>
              <CardDescription className="text-xs">
                Provide your official logo and favicon for sidebars, topbars, and reports.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-1.5">
                <Label htmlFor="logoUrl" className="text-xs font-medium">Corporate Logo URL</Label>
                <Input
                  id="logoUrl"
                  placeholder="https://example.com/logo.png"
                  value={logoUrl}
                  onChange={(e) => setLogoUrl(e.target.value)}
                  className="text-xs h-9"
                />
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="faviconUrl" className="text-xs font-medium">Favicon URL (.ico / .png)</Label>
                <Input
                  id="faviconUrl"
                  placeholder="https://example.com/favicon.ico"
                  value={faviconUrl}
                  onChange={(e) => setFaviconUrl(e.target.value)}
                  className="text-xs h-9"
                />
              </div>
            </CardContent>
          </Card>

          {/* Color Palette & Swatches */}
          <Card className="border-border shadow-sm">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-semibold flex items-center gap-2">
                <Palette className="w-4 h-4 text-primary" />
                Primary Brand Color
              </CardTitle>
              <CardDescription className="text-xs">
                Applied to primary buttons, active navigation markers, and analytics highlights.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Presets Grid */}
              <div className="grid grid-cols-4 gap-2">
                {colorPresets.map((preset) => (
                  <button
                    key={preset.hex}
                    type="button"
                    onClick={() => setPrimaryColor(preset.hex)}
                    className="group relative flex flex-col items-center gap-1 p-2 rounded-lg border border-border hover:border-primary/50 transition-all text-left bg-card hover:bg-secondary/40"
                    title={`${preset.name} - ${preset.desc}`}
                  >
                    <div
                      className="w-6 h-6 rounded-full border border-black/10 shadow-sm flex items-center justify-center transition-transform group-hover:scale-110"
                      style={{ backgroundColor: preset.hex }}
                    >
                      {primaryColor === preset.hex && (
                        <CheckCircle2 className="w-3.5 h-3.5 text-white" />
                      )}
                    </div>
                    <span className="text-[10px] font-medium text-foreground truncate w-full text-center">
                      {preset.name.split(" ")[0]}
                    </span>
                  </button>
                ))}
              </div>

              {/* Custom Hex Input */}
              <div className="flex items-center gap-3 pt-2">
                <div
                  className="w-9 h-9 rounded-lg border border-border shrink-0 shadow-inner"
                  style={{ backgroundColor: primaryColor }}
                />
                <div className="flex-1">
                  <Label htmlFor="customColor" className="text-[11px] font-medium text-muted-foreground">
                    Custom Hex Code
                  </Label>
                  <Input
                    id="customColor"
                    value={primaryColor}
                    onChange={(e) => setPrimaryColor(e.target.value)}
                    className="text-xs h-8 font-mono"
                    placeholder="#0d9488"
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Typography & Geometry */}
          <Card className="border-border shadow-sm">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-semibold flex items-center gap-2">
                <Type className="w-4 h-4 text-primary" />
                Typography & Interface Geometry
              </CardTitle>
              <CardDescription className="text-xs">
                Select your preferred brand typography and UI component corner roundness.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-1.5">
                <Label className="text-xs font-medium">Font Family</Label>
                <div className="grid grid-cols-2 gap-2">
                  {fontPresets.map((f) => (
                    <button
                      key={f.family}
                      type="button"
                      onClick={() => setFontFamily(f.family)}
                      className={`p-2.5 rounded-lg border text-left transition-all ${
                        fontFamily === f.family
                          ? "border-primary bg-primary/5 shadow-sm"
                          : "border-border hover:border-border/80 bg-card"
                      }`}
                    >
                      <p className="text-xs font-semibold text-foreground">{f.name}</p>
                      <p className="text-[10px] text-muted-foreground">{f.style}</p>
                    </button>
                  ))}
                </div>
              </div>

              <div className="space-y-1.5 pt-2">
                <Label className="text-xs font-medium">Corner Roundness</Label>
                <div className="grid grid-cols-4 gap-2">
                  {[
                    { label: "Sharp", val: "0.25rem" },
                    { label: "Modern", val: "0.5rem" },
                    { label: "Smooth", val: "0.75rem" },
                    { label: "Curved", val: "1rem" },
                  ].map((r) => (
                    <button
                      key={r.val}
                      type="button"
                      onClick={() => setBorderRadius(r.val)}
                      className={`py-1.5 px-2 rounded-lg border text-center text-xs transition-all ${
                        borderRadius === r.val
                          ? "border-primary bg-primary/5 font-semibold text-primary"
                          : "border-border hover:border-border/80 text-muted-foreground"
                      }`}
                    >
                      {r.label}
                    </button>
                  ))}
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Whitelabel Custom Domain (Enterprise) */}
          <Card className="border-border shadow-sm">
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between">
                <CardTitle className="text-sm font-semibold flex items-center gap-2">
                  <Globe className="w-4 h-4 text-primary" />
                  Whitelabel Custom Domain
                </CardTitle>
                <Badge variant="secondary" className="text-[9px] uppercase font-mono">
                  Enterprise
                </Badge>
              </div>
              <CardDescription className="text-xs">
                Serve Curexal directly from your clinic’s custom domain with automatic SSL.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-3">
              <Input
                placeholder="app.yourhospital.org"
                value={customDomain}
                onChange={(e) => setCustomDomain(e.target.value)}
                className="text-xs h-9 font-mono"
              />
              <p className="text-[11px] text-muted-foreground flex items-center gap-1.5">
                <Shield className="w-3.5 h-3.5 text-emerald-500" />
                CNAME records automatically provision Let's Encrypt SSL certificates.
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Right Column: Live Real-Time Interactive Preview (7 cols) */}
        <div className="lg:col-span-7 space-y-4">
          <div className="flex items-center justify-between bg-card border border-border p-2 rounded-xl">
            <span className="text-xs font-semibold text-foreground px-2 flex items-center gap-1.5">
              <Eye className="w-3.5 h-3.5 text-primary" />
              Live Real-Time Brand Preview
            </span>
            <div className="flex items-center gap-1 bg-secondary/50 p-0.5 rounded-lg">
              <button
                type="button"
                onClick={() => setPreviewTab("ui")}
                className={`text-xs px-3 py-1 rounded-md font-medium transition-all ${
                  previewTab === "ui" ? "bg-card text-foreground shadow-sm" : "text-muted-foreground"
                }`}
              >
                Console UI
              </button>
              <button
                type="button"
                onClick={() => setPreviewTab("report")}
                className={`text-xs px-3 py-1 rounded-md font-medium transition-all ${
                  previewTab === "report" ? "bg-card text-foreground shadow-sm" : "text-muted-foreground"
                }`}
              >
                Diagnostic PDF Report
              </button>
            </div>
          </div>

          {previewTab === "ui" ? (
            /* Live Simulated Console Card */
            <div
              className="rounded-2xl border border-border bg-card shadow-lg p-6 space-y-6 transition-all"
              style={{ fontFamily, borderRadius }}
            >
              {/* Simulated Topbar */}
              <div className="flex items-center justify-between border-b border-border pb-4">
                <div className="flex items-center gap-3">
                  <div
                    className="w-8 h-8 rounded-xl flex items-center justify-center text-white shadow-sm font-bold text-xs"
                    style={{ backgroundColor: primaryColor, borderRadius }}
                  >
                    {logoUrl ? (
                      <img src={logoUrl} alt="logo" className="w-full h-full object-contain rounded-xl" />
                    ) : (
                      orgName.slice(0, 2).toUpperCase()
                    )}
                  </div>
                  <div>
                    <h3 className="text-sm font-bold text-foreground leading-tight">{orgName}</h3>
                    <p className="text-[10px] text-muted-foreground">Executive Healthcare Console</p>
                  </div>
                </div>
                <Badge
                  style={{ backgroundColor: `${primaryColor}20`, color: primaryColor, borderColor: `${primaryColor}40` }}
                  className="text-[10px] uppercase font-mono font-semibold"
                >
                  Active Org HQ
                </Badge>
              </div>

              {/* Simulated Navigation Bar */}
              <div className="flex items-center gap-2 overflow-x-auto pb-1 text-xs">
                {[
                  { name: "Executive Dashboard", active: true },
                  { name: "Branch Facilities", active: false },
                  { name: "Medical Lab (LIS)", active: false },
                  { name: "Billing POS", active: false },
                ].map((tab, i) => (
                  <div
                    key={tab.name}
                    className={`px-3 py-1.5 text-xs font-medium cursor-pointer transition-all ${
                      tab.active
                        ? "text-white shadow-sm font-semibold"
                        : "text-muted-foreground hover:text-foreground bg-secondary/30"
                    }`}
                    style={{
                      backgroundColor: tab.active ? primaryColor : undefined,
                      borderRadius,
                    }}
                  >
                    {tab.name}
                  </div>
                ))}
              </div>

              {/* Simulated Stat Cards */}
              <div className="grid grid-cols-2 gap-4">
                <div className="p-4 rounded-xl border border-border bg-secondary/10 space-y-1">
                  <span className="text-[11px] text-muted-foreground">Daily Diagnostic Tests</span>
                  <div className="text-2xl font-bold text-foreground">1,248</div>
                  <span className="text-[10px] font-medium text-emerald-600 dark:text-emerald-400 flex items-center gap-1">
                    ↑ +12.4% vs last week
                  </span>
                </div>
                <div className="p-4 rounded-xl border border-border bg-secondary/10 space-y-1">
                  <span className="text-[11px] text-muted-foreground">Consolidated Revenue</span>
                  <div className="text-2xl font-bold text-foreground">₦4,850,000</div>
                  <span className="text-[10px] font-medium text-emerald-600 dark:text-emerald-400 flex items-center gap-1">
                    ↑ +8.1% vs last week
                  </span>
                </div>
              </div>

              {/* Simulated Action Button */}
              <div className="flex items-center justify-between p-4 rounded-xl bg-secondary/20 border border-border">
                <div className="space-y-0.5">
                  <p className="text-xs font-semibold text-foreground">Patient Accessioning Register</p>
                  <p className="text-[11px] text-muted-foreground">14 specimen batches pending review</p>
                </div>
                <button
                  type="button"
                  style={{ backgroundColor: primaryColor, borderRadius }}
                  className="px-3.5 py-1.5 text-xs font-medium text-white shadow hover:opacity-90 transition-opacity"
                >
                  Open Worklist
                </button>
              </div>
            </div>
          ) : (
            /* Live Branded PDF Report Preview */
            <div
              className="rounded-2xl border border-border bg-white text-slate-900 shadow-lg p-8 space-y-6 transition-all"
              style={{ fontFamily, borderRadius }}
            >
              {/* Report Header */}
              <div className="flex items-start justify-between border-b-2 pb-4" style={{ borderColor: primaryColor }}>
                <div className="flex items-center gap-3">
                  <div
                    className="w-12 h-12 rounded-xl flex items-center justify-center text-white font-bold text-base shadow-sm"
                    style={{ backgroundColor: primaryColor }}
                  >
                    {logoUrl ? <img src={logoUrl} alt="logo" className="w-full h-full object-contain" /> : orgName.slice(0, 2).toUpperCase()}
                  </div>
                  <div>
                    <h2 className="text-base font-extrabold text-slate-900 leading-tight uppercase">{orgName}</h2>
                    <p className="text-[11px] text-slate-600 font-medium">{headerText}</p>
                    <p className="text-[10px] text-slate-500">Accredited Pathology & Diagnostic Center</p>
                  </div>
                </div>
                <div className="text-right text-[11px] text-slate-600 space-y-0.5">
                  <p className="font-bold text-slate-900">LABORATORY TEST REPORT</p>
                  <p>Sample ID: <span className="font-mono font-semibold">LAB-2026-0841</span></p>
                  <p>Date: {new Date().toLocaleDateString()}</p>
                </div>
              </div>

              {/* Patient Info Bar */}
              <div className="grid grid-cols-3 gap-3 p-3 rounded-lg bg-slate-50 border border-slate-200 text-xs">
                <div>
                  <span className="text-[10px] text-slate-500 uppercase">Patient Name</span>
                  <p className="font-bold text-slate-900">Patrick Dimkpa</p>
                </div>
                <div>
                  <span className="text-[10px] text-slate-500 uppercase">Age / Gender</span>
                  <p className="font-bold text-slate-900">38 Yrs / Male</p>
                </div>
                <div>
                  <span className="text-[10px] text-slate-500 uppercase">Referring Physician</span>
                  <p className="font-bold text-slate-900">Dr. A. Adebayo (Consultant)</p>
                </div>
              </div>

              {/* Diagnostic Results Table */}
              <table className="w-full text-xs text-left">
                <thead>
                  <tr className="border-b border-slate-200 text-[10px] uppercase font-bold text-slate-500">
                    <th className="py-2">Test Investigation</th>
                    <th className="py-2">Observed Result</th>
                    <th className="py-2">Biological Reference</th>
                    <th className="py-2">Flag</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100 text-slate-800">
                  <tr>
                    <td className="py-2 font-medium">Hemoglobin (Hb)</td>
                    <td className="py-2 font-bold text-slate-900">14.8 g/dL</td>
                    <td className="py-2 text-slate-500">13.0 - 17.5 g/dL</td>
                    <td className="py-2"><span className="text-[10px] px-1.5 py-0.5 rounded bg-emerald-100 text-emerald-800 font-bold">NORMAL</span></td>
                  </tr>
                  <tr>
                    <td className="py-2 font-medium">Fasting Blood Glucose</td>
                    <td className="py-2 font-bold text-amber-700">114 mg/dL</td>
                    <td className="py-2 text-slate-500">70 - 99 mg/dL</td>
                    <td className="py-2"><span className="text-[10px] px-1.5 py-0.5 rounded bg-amber-100 text-amber-800 font-bold">ELEVATED</span></td>
                  </tr>
                  <tr>
                    <td className="py-2 font-medium">Total Cholesterol</td>
                    <td className="py-2 font-bold text-slate-900">182 mg/dL</td>
                    <td className="py-2 text-slate-500">&lt; 200 mg/dL</td>
                    <td className="py-2"><span className="text-[10px] px-1.5 py-0.5 rounded bg-emerald-100 text-emerald-800 font-bold">NORMAL</span></td>
                  </tr>
                </tbody>
              </table>

              {/* Report Footer */}
              <div className="border-t border-slate-200 pt-4 flex items-center justify-between text-[10px] text-slate-500">
                <p className="max-w-md">{footerDisclaimer}</p>
                <div className="flex items-center gap-1 font-semibold text-slate-700">
                  <Award className="w-3.5 h-3.5" style={{ color: primaryColor }} />
                  Verified by Chief Pathologist
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
