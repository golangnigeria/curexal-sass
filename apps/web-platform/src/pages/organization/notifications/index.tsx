import React, { useState, useEffect } from "react";
import { organizationService } from "@/api/services/organization.service";
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
  Loader2,
} from "lucide-react";

export default function OrganizationNotificationsPage() {
  const [activeTab, setActiveTab] = useState<"channels" | "templates">("channels");

  // SMTP State
  const [smtpHost, setSmtpHost] = useState("smtp.sendgrid.net");
  const [smtpPort, setSmtpPort] = useState("587");
  const [smtpUser, setSmtpUser] = useState("apikey");
  const [smtpPass, setSmtpPass] = useState("");
  const [senderEmail, setSenderEmail] = useState("notifications@curexalhealth.com");
  const [senderName, setSenderName] = useState("Diagnostic Center Laboratory");

  // SMS Gateway
  const [smsProvider, setSmsProvider] = useState("termii");
  const [smsApiKey, setSmsApiKey] = useState("");
  const [smsSenderId, setSmsSenderId] = useState("CUREXAL-LAB");

  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    let isMounted = true;
    async function loadConfigs() {
      try {
        const configs = await organizationService.getNotificationConfigs();
        if (isMounted && configs) {
          if (configs.smtp) {
            setSmtpHost(configs.smtp.host || "smtp.sendgrid.net");
            setSmtpPort(String(configs.smtp.port || "587"));
            setSmtpUser(configs.smtp.user || "apikey");
            setSenderEmail(configs.smtp.fromEmail || "notifications@curexalhealth.com");
            setSenderName(configs.smtp.fromName || "Diagnostic Center Laboratory");
          }
          if (configs.sms) {
            setSmsProvider(configs.sms.provider || "termii");
            setSmsSenderId(configs.sms.senderId || "CUREXAL-LAB");
          }
        }
      } catch {
        // Fallback default
      } finally {
        if (isMounted) setIsLoading(false);
      }
    }
    loadConfigs();
    return () => {
      isMounted = false;
    };
  }, []);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSaving(true);
    try {
      await organizationService.saveNotificationConfig({
        smtp: {
          host: smtpHost,
          port: parseInt(smtpPort, 10) || 587,
          user: smtpUser,
          pass: smtpPass || undefined,
          fromEmail: senderEmail,
          fromName: senderName,
        },
        sms: {
          provider: smsProvider,
          apiKey: smsApiKey || undefined,
          senderId: smsSenderId,
        },
      });
      toast.success("Notification Channel Settings Saved!", {
        description: "Automated SMS and email dispatches will use these credentials.",
      });
    } catch (err: any) {
      toast.error("Failed to save notification settings: " + (err.message || "Network error"));
    } finally {
      setIsSaving(false);
    }
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
            {isSaving ? <Loader2 className="w-3.5 h-3.5 animate-spin" /> : <Save className="w-3.5 h-3.5" />}
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

      {activeTab === "channels" && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Custom SMTP Email */}
          <Card className="border-border shadow-sm">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-bold flex items-center gap-2">
                <Mail className="w-4 h-4 text-primary" />
                Custom SMTP Mail Server
              </CardTitle>
              <CardDescription className="text-xs">
                Deliver clinical lab reports and patient discharge summaries from your own domain.
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
                  <Label className="text-xs font-medium">Password / Secret Token</Label>
                  <Input
                    type="password"
                    placeholder="••••••••••••••••"
                    value={smtpPass}
                    onChange={(e) => setSmtpPass(e.target.value)}
                    className="text-xs h-9 font-mono"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-1.5">
                  <Label className="text-xs font-medium">Sender From Email</Label>
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
              <CardTitle className="text-sm font-bold flex items-center gap-2">
                <MessageSquare className="w-4 h-4 text-primary" />
                SMS & WhatsApp Gateway
              </CardTitle>
              <CardDescription className="text-xs">
                Send appointment reminders and critical panic test result SMS alerts directly to patient mobile lines.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-1.5">
                <Label className="text-xs font-medium">SMS Provider Gateway</Label>
                <Input value={smsProvider} onChange={(e) => setSmsProvider(e.target.value)} className="text-xs h-9 capitalize" />
              </div>

              <div className="space-y-1.5">
                <Label className="text-xs font-medium">Registered Sender ID</Label>
                <Input value={smsSenderId} onChange={(e) => setSmsSenderId(e.target.value)} className="text-xs h-9 font-mono uppercase" />
              </div>

              <div className="space-y-1.5">
                <Label className="text-xs font-medium">Gateway API Secret Key</Label>
                <Input
                  type="password"
                  placeholder="••••••••••••••••••••••••"
                  value={smsApiKey}
                  onChange={(e) => setSmsApiKey(e.target.value)}
                  className="text-xs h-9 font-mono"
                />
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {activeTab === "templates" && (
        <Card className="border-border shadow-sm">
          <CardHeader className="pb-3">
            <CardTitle className="text-sm font-bold">Automated Clinical Alert Templates</CardTitle>
            <CardDescription className="text-xs">
              System variables like <code className="text-primary">{`{{patient_name}}`}</code>, <code className="text-primary">{`{{test_name}}`}</code>, and <code className="text-primary">{`{{accession_number}}`}</code> are populated automatically at runtime.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="p-4 rounded-xl border border-border bg-card/50 space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-xs font-bold text-foreground">Laboratory Result Ready Notification</span>
                <Badge variant="outline" className="text-[10px] bg-emerald-500/10 text-emerald-500 border-emerald-500/20">SMS & Email</Badge>
              </div>
              <p className="text-xs font-mono text-muted-foreground">
                Dear {`{{patient_name}}`}, your medical laboratory results ({`{{test_name}}`}) from {`{{facility_name}}`} are now finalized and available. Access your secure PDF here: {`{{report_link}}`}
              </p>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
