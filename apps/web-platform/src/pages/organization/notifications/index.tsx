import React, { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import {
  Bell,
  Mail,
  MessageSquare,
  Send,
  Save,
  CheckCircle2,
  Lock,
  Sparkles,
} from "lucide-react";

export default function OrganizationNotificationsPage() {
  const [activeTab, setActiveTab] = useState<"channels" | "templates">("channels");

  // SMTP State
  const [smtpHost, setSmtpHost] = useState("smtp.sendgrid.net");
  const [smtpPort, setSmtpPort] = useState("587");
  const [smtpUser, setSmtpUser] = useState("apikey");
  const [smtpPass, setSmtpPass] = useState("••••••••••••••••");
  const [senderEmail, setSenderEmail] = useState("notifications@curexalhealth.com");
  const [senderName, setSenderName] = useState("Diagnostic Center Laboratory");

  // SMS Gateway
  const [smsProvider, setSmsProvider] = useState("termii");
  const [smsApiKey, setSmsApiKey] = useState("••••••••••••••••••••••••");
  const [smsSenderId, setSmsSenderId] = useState("CUREXAL-LAB");

  const [isSaving, setIsSaving] = useState(false);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSaving(true);
    await new Promise((r) => setTimeout(r, 500));
    setIsSaving(false);
    toast.success("Notification Channel Settings Saved!", {
      description: "Automated SMS and email dispatches will use these credentials.",
    });
  };

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-6">
        <div>
          <div className="flex items-center gap-2.5 mb-1">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              Notification Channels & Templates
            </h1>
            <Badge variant="outline" className="border-primary/40 text-primary bg-primary/10 text-[10px] font-mono">
              AEAD Encrypted
            </Badge>
          </div>
          <p className="text-xs text-muted-foreground">
            Configure custom SMTP email servers, SMS gateways, and automated patient test result alerts.
          </p>
        </div>

        <div className="flex items-center gap-2">
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
      </div>

      {/* Tabs */}
      <div className="flex items-center gap-2 border-b border-border pb-2">
        <button
          type="button"
          onClick={() => setActiveTab("channels")}
          className={`text-xs font-semibold px-3 py-1.5 rounded-lg transition-all ${
            activeTab === "channels" ? "bg-primary text-primary-foreground shadow-sm" : "text-muted-foreground hover:text-foreground"
          }`}
        >
          Delivery Channels (SMTP / SMS)
        </button>
        <button
          type="button"
          onClick={() => setActiveTab("templates")}
          className={`text-xs font-semibold px-3 py-1.5 rounded-lg transition-all ${
            activeTab === "templates" ? "bg-primary text-primary-foreground shadow-sm" : "text-muted-foreground hover:text-foreground"
          }`}
        >
          Message Templates
        </button>
      </div>

      {activeTab === "channels" ? (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          {/* Custom SMTP Email Provider */}
          <Card className="border-border shadow-sm">
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between">
                <CardTitle className="text-sm font-bold flex items-center gap-2">
                  <Mail className="w-4 h-4 text-primary" />
                  Custom SMTP Email Provider
                </CardTitle>
                <Badge className="text-[9px] bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/30">
                  Verified
                </Badge>
              </div>
              <CardDescription className="text-xs">
                Deliver patient test results and invoices from your own custom email domain.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-3 gap-3">
                <div className="col-span-2 space-y-1.5">
                  <Label className="text-xs font-medium">SMTP Server Host</Label>
                  <Input value={smtpHost} onChange={(e) => setSmtpHost(e.target.value)} className="text-xs h-9 font-mono" />
                </div>
                <div className="space-y-1.5">
                  <Label className="text-xs font-medium">Port</Label>
                  <Input value={smtpPort} onChange={(e) => setSmtpPort(e.target.value)} className="text-xs h-9 font-mono" />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-1.5">
                  <Label className="text-xs font-medium">Username / API Key</Label>
                  <Input value={smtpUser} onChange={(e) => setSmtpUser(e.target.value)} className="text-xs h-9 font-mono" />
                </div>
                <div className="space-y-1.5">
                  <Label className="text-xs font-medium">Password / Secret</Label>
                  <Input type="password" value={smtpPass} onChange={(e) => setSmtpPass(e.target.value)} className="text-xs h-9 font-mono" />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-1.5">
                  <Label className="text-xs font-medium">Sender Email</Label>
                  <Input value={senderEmail} onChange={(e) => setSenderEmail(e.target.value)} className="text-xs h-9" />
                </div>
                <div className="space-y-1.5">
                  <Label className="text-xs font-medium">Sender Display Name</Label>
                  <Input value={senderName} onChange={(e) => setSenderName(e.target.value)} className="text-xs h-9" />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* SMS Notification Gateway */}
          <Card className="border-border shadow-sm">
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between">
                <CardTitle className="text-sm font-bold flex items-center gap-2">
                  <MessageSquare className="w-4 h-4 text-primary" />
                  SMS Notification Gateway
                </CardTitle>
                <Badge className="text-[9px] bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/30">
                  Connected
                </Badge>
              </div>
              <CardDescription className="text-xs">
                Dispatch automated SMS alerts when diagnostic results are ready for download.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-1.5">
                <Label className="text-xs font-medium">SMS Provider</Label>
                <select
                  value={smsProvider}
                  onChange={(e) => setSmsProvider(e.target.value)}
                  className="w-full h-9 rounded-md border border-input bg-background px-3 text-xs focus:outline-none focus:ring-1 focus:ring-primary"
                >
                  <option value="termii">Termii (African Telcos & DND Bypassing)</option>
                  <option value="twilio">Twilio (Global Telephony)</option>
                  <option value="africastalking">Africa's Talking</option>
                </select>
              </div>

              <div className="space-y-1.5">
                <Label className="text-xs font-medium">Provider API Key / Secret</Label>
                <Input type="password" value={smsApiKey} onChange={(e) => setSmsApiKey(e.target.value)} className="text-xs h-9 font-mono" />
              </div>

              <div className="space-y-1.5">
                <Label className="text-xs font-medium">Registered Sender ID (Alphanumeric)</Label>
                <Input value={smsSenderId} onChange={(e) => setSmsSenderId(e.target.value)} className="text-xs h-9 font-mono uppercase" />
              </div>
            </CardContent>
          </Card>
        </div>
      ) : (
        /* Message Templates List */
        <div className="space-y-4">
          {[
            {
              key: "lab.results_ready",
              title: "Laboratory Diagnostic Results Ready",
              channel: "SMS & Email",
              preview: "Dear {{patient_name}}, your {{test_name}} results from {{org_name}} are ready. Download: {{result_link}}",
            },
            {
              key: "radiology.scan_scheduled",
              title: "Radiology Scan Appointment Confirmation",
              channel: "SMS",
              preview: "Hello {{patient_name}}, your {{modality}} appointment is scheduled on {{appointment_time}} at {{branch_name}}.",
            },
            {
              key: "billing.invoice_receipt",
              title: "Patient Receipt & Payment Confirmation",
              channel: "Email",
              preview: "Payment received: {{currency}} {{amount_paid}} for Invoice {{invoice_number}}. Thank you for choosing {{org_name}}.",
            },
          ].map((t) => (
            <Card key={t.key} className="border-border shadow-sm p-4 flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
              <div className="space-y-1 max-w-xl">
                <div className="flex items-center gap-2">
                  <span className="text-xs font-bold text-foreground">{t.title}</span>
                  <Badge variant="outline" className="text-[9px] font-mono">{t.channel}</Badge>
                </div>
                <p className="text-[11px] font-mono text-muted-foreground bg-secondary/30 p-2 rounded-lg border border-border/50">
                  {t.preview}
                </p>
              </div>
              <Button size="sm" variant="outline" className="text-xs h-8 shrink-0">
                Edit Template
              </Button>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
