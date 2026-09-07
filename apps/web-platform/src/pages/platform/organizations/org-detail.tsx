import React, { useState } from "react";
import { useParams, Link } from "react-router-dom";
import {
  Building2,
  ArrowLeft,
  ShieldCheck,
  CheckCircle2,
  XCircle,
  FileText,
  Layers,
  Settings,
  Plus,
  Trash2,
  Clock,
  ExternalLink,
  Upload,
  AlertTriangle,
  Mail,
  Eye,
  Download,
  Pencil,
} from "lucide-react";
import { DocumentPreviewViewer } from "@/components/document/document-preview-viewer";
import {
  useOrganization,
  useOrganizationSettings,
  useOrganizationDocuments,
  useUpdateOrganization,
  useUpdateOrganizationSettings,
  useApproveOrganization,
  useRejectOrganization,
  useReviewDocument,
  useResendOwnerInvite,
} from "@/api/hooks/use-organizations";
import {
  useOrganizationCapabilities,
  useOrganizationEntitlements,
  useGrantCapability,
  useRevokeCapability,
  useStartTrialCapability,
  useCapabilityCatalog,
} from "@/api/hooks/use-marketplace";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { formatDate, formatFileSize } from "@/lib/utils";
import { toast } from "sonner";

export default function OrganizationDetailPage() {
  const { id } = useParams<{ id: string }>();
  const orgId = id || "";

  const { data: org, isLoading: orgLoading, refetch: refetchOrg } = useOrganization(orgId);
  const { data: settings } = useOrganizationSettings(orgId);
  const { data: documents, refetch: refetchDocs } = useOrganizationDocuments(orgId);
  const { data: capsData, refetch: refetchCaps } = useOrganizationCapabilities(orgId);
  const { data: entitlements, refetch: refetchEntitlements } = useOrganizationEntitlements(orgId);
  const { data: catalog } = useCapabilityCatalog();

  const updateOrgMutation = useUpdateOrganization(orgId);
  const updateSettingsMutation = useUpdateOrganizationSettings(orgId);
  const approveMutation = useApproveOrganization(orgId);
  const rejectMutation = useRejectOrganization(orgId);
  const grantCapMutation = useGrantCapability(orgId);
  const revokeCapMutation = useRevokeCapability(orgId);
  const trialCapMutation = useStartTrialCapability(orgId);
  const resendInviteMutation = useResendOwnerInvite(orgId);

  // Edit Org State
  const [editName, setEditName] = useState("");
  const [editPlan, setEditPlan] = useState("");
  const [editCustomDomain, setEditCustomDomain] = useState("");

  // Edit Corporate Profile & Tax Metadata State
  const [openEditMetadata, setOpenEditMetadata] = useState(false);
  const [editCacNumber, setEditCacNumber] = useState("");
  const [editTinNumber, setEditTinNumber] = useState("");
  const [editBusinessType, setEditBusinessType] = useState("");
  const [editSupportEmail, setEditSupportEmail] = useState("");
  const [editSupportPhone, setEditSupportPhone] = useState("");
  const [editAddress, setEditAddress] = useState("");
  const [editCurrency, setEditCurrency] = useState("NGN");

  const handleOpenEditMetadata = () => {
    setEditCacNumber(settings?.cacNumber || org?.registrationNumber || "");
    setEditTinNumber(settings?.tinNumber || settings?.taxNumber || org?.taxId || "");
    setEditBusinessType(settings?.businessType || "Healthcare Provider");
    setEditSupportEmail(settings?.supportEmail || org?.email || "");
    setEditSupportPhone(settings?.supportPhone || org?.phone || "");
    setEditAddress(settings?.address || org?.address || "");
    setEditCurrency(settings?.currency || "NGN");
    setOpenEditMetadata(true);
  };

  const handleUpdateMetadata = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await updateSettingsMutation.mutateAsync({
        cacNumber: editCacNumber,
        tinNumber: editTinNumber,
        taxNumber: editTinNumber,
        businessType: editBusinessType,
        supportEmail: editSupportEmail,
        supportPhone: editSupportPhone,
        address: editAddress,
        currency: editCurrency,
      });
      toast.success("Corporate profile & tax metadata updated successfully!");
      setOpenEditMetadata(false);
      refetchOrg();
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Failed to update corporate metadata");
    }
  };

  // Reject Modal State
  const [rejectReason, setRejectReason] = useState("");
  const [openReject, setOpenReject] = useState(false);

  // Grant Capability Modal State
  const [openGrant, setOpenGrant] = useState(false);
  const [selectedCapCode, setSelectedCapCode] = useState("");
  const [grantSource, setGrantSource] = useState("addon");

  // Review Doc Modal State
  const [reviewDocId, setReviewDocId] = useState<string | null>(null);
  const [docRejectReason, setDocRejectReason] = useState("");
  const [docReviewStatus, setDocReviewStatus] = useState<"approved" | "rejected">("approved");
  const reviewDocMutation = useReviewDocument(reviewDocId || "", orgId);

  // Preview Modal State
  const [previewDoc, setPreviewDoc] = useState<any | null>(null);
  const [isPreviewOpen, setIsPreviewOpen] = useState(false);

  const handleOpenPreview = (doc: any) => {
    setPreviewDoc(doc);
    setIsPreviewOpen(true);
  };

  const handleUpdateOrg = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await updateOrgMutation.mutateAsync({
        name: editName || org?.name,
        plan: editPlan || org?.plan,
        customDomain: editCustomDomain || org?.customDomain,
      });
      toast.success("Organization details updated successfully!");
      refetchOrg();
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Failed to update organization");
    }
  };

  const handleApprove = async () => {
    try {
      await approveMutation.mutateAsync();
      toast.success("Organization approved and verified successfully!");
      refetchOrg();
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Approval failed");
    }
  };

  const handleReject = async () => {
    if (!rejectReason) {
      toast.error("Please provide a rejection reason.");
      return;
    }
    try {
      await rejectMutation.mutateAsync(rejectReason);
      toast.success("Organization verification rejected.");
      setOpenReject(false);
      refetchOrg();
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Rejection failed");
    }
  };

  const handleResendInvite = async () => {
    try {
      await resendInviteMutation.mutateAsync();
      toast.success("Owner invitation verification code resent successfully!");
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Failed to resend owner invitation");
    }
  };

  const handleGrantCapability = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedCapCode) return;
    try {
      await grantCapMutation.mutateAsync({
        capabilityCode: selectedCapCode,
        source: grantSource,
      });
      toast.success(`Granted capability '${selectedCapCode}'!`);
      setOpenGrant(false);
      refetchCaps();
      refetchEntitlements();
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Failed to grant capability");
    }
  };

  const handleRevokeCapability = async (capCode: string) => {
    if (!confirm(`Are you sure you want to revoke '${capCode}' entitlement?`)) return;
    try {
      await revokeCapMutation.mutateAsync(capCode);
      toast.success(`Revoked capability '${capCode}'`);
      refetchCaps();
      refetchEntitlements();
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Failed to revoke capability");
    }
  };

  const handleReviewDocSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!reviewDocId) return;
    try {
      await reviewDocMutation.mutateAsync({
        status: docReviewStatus,
        rejectionReason: docReviewStatus === "rejected" ? docRejectReason : undefined,
      });
      toast.success(`Document marked as ${docReviewStatus}`);
      setReviewDocId(null);
      refetchDocs();
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Failed to review document");
    }
  };

  if (orgLoading) {
    return (
      <div className="py-16 text-center text-sm text-muted-foreground">
        Loading organization details...
      </div>
    );
  }

  if (!org) {
    return (
      <div className="py-16 text-center space-y-3">
        <h2 className="text-lg font-bold">Organization Not Found</h2>
        <Button asChild variant="outline" size="sm">
          <Link to="/platform/organizations">Back to Organizations</Link>
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Back Link & Header */}
      <div>
        <Link
          to="/platform/organizations"
          className="inline-flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground mb-3"
        >
          <ArrowLeft className="h-3.5 w-3.5" /> Back to Organizations Directory
        </Link>

        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10 text-primary border border-primary/20">
              <Building2 className="h-6 w-6" />
            </div>
            <div>
              <div className="flex items-center gap-2">
                <h1 className="text-2xl font-bold tracking-tight text-foreground">
                  {org.name}
                </h1>
                <Badge variant="outline" className="capitalize text-xs font-mono">
                  {org.plan || "Smart"}
                </Badge>
                {org.status === "active" ? (
                  <Badge variant="outline" className="border-emerald-500/30 text-emerald-600 bg-emerald-500/5 text-xs">
                    <CheckCircle2 className="mr-1 h-3 w-3 inline" /> Active
                  </Badge>
                ) : org.status === "pending_verification" ? (
                  <Badge variant="outline" className="border-amber-500/30 text-amber-600 bg-amber-500/5 text-xs">
                    <Clock className="mr-1 h-3 w-3 inline" /> Pending Verification
                  </Badge>
                ) : (
                  <Badge variant="outline" className="border-destructive/30 text-destructive bg-destructive/5 text-xs">
                    {org.status}
                  </Badge>
                )}
              </div>
              <p className="text-xs text-muted-foreground font-mono mt-0.5">
                slug: {org.slug} • id: {org.id}
              </p>
            </div>
          </div>

          {/* Verification Actions */}
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={handleResendInvite}
              disabled={resendInviteMutation.isPending}
              className="text-xs h-9 gap-1.5 border-primary/30 text-primary hover:bg-primary/5"
            >
              <Mail className="h-3.5 w-3.5" />
              {resendInviteMutation.isPending ? "Resending..." : "Resend Invitation"}
            </Button>

            {org.status !== "active" && (
              <Button
                size="sm"
                onClick={handleApprove}
                disabled={approveMutation.isPending}
                className="bg-emerald-600 hover:bg-emerald-700 text-white text-xs h-9 gap-1.5 shadow-sm"
              >
                <CheckCircle2 className="h-3.5 w-3.5" />
                Approve & Activate
              </Button>
            )}

            <Dialog open={openReject} onOpenChange={setOpenReject}>
              <DialogTrigger asChild>
                <Button
                  variant="outline"
                  size="sm"
                  className="border-destructive/30 text-destructive hover:bg-destructive/10 text-xs h-9 gap-1.5"
                >
                  <XCircle className="h-3.5 w-3.5" />
                  Reject Organization
                </Button>
              </DialogTrigger>
              <DialogContent className="sm:max-w-md">
                <DialogHeader>
                  <DialogTitle className="text-base">Reject Organization</DialogTitle>
                  <DialogDescription className="text-xs">
                    Reject onboarding for {org.name}. Please specify the administrative reason.
                  </DialogDescription>
                </DialogHeader>
                <div className="space-y-2 py-2">
                  <Label htmlFor="rejReason" className="text-xs">Rejection Reason</Label>
                  <Input
                    id="rejReason"
                    placeholder="e.g. Incomplete CAC corporate filing documents"
                    value={rejectReason}
                    onChange={(e) => setRejectReason(e.target.value)}
                    className="text-xs"
                  />
                </div>
                <DialogFooter>
                  <Button variant="outline" size="sm" onClick={() => setOpenReject(false)} className="text-xs">
                    Cancel
                  </Button>
                  <Button
                    size="sm"
                    variant="destructive"
                    onClick={handleReject}
                    disabled={rejectMutation.isPending}
                    className="text-xs"
                  >
                    Confirm Rejection
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>
          </div>
        </div>
      </div>

      {/* Tabs Navigation */}
      <Tabs defaultValue="overview" className="space-y-4">
        <TabsList className="bg-muted/60 p-1 border border-border">
          <TabsTrigger value="overview" className="text-xs">Overview & Settings</TabsTrigger>
          <TabsTrigger value="capabilities" className="text-xs">Capabilities & Entitlements</TabsTrigger>
          <TabsTrigger value="compliance" className="text-xs">Compliance Documents</TabsTrigger>
        </TabsList>

        {/* TAB 1: OVERVIEW */}
        <TabsContent value="overview" className="space-y-4">
          <div className="grid gap-6 md:grid-cols-2">
            <Card className="card-enterprise">
              <CardHeader>
                <CardTitle className="text-sm font-semibold">General Details</CardTitle>
                <CardDescription className="text-xs">
                  Update primary parameters and subscription tier for this organization.
                </CardDescription>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleUpdateOrg} className="space-y-3 text-xs">
                  <div className="space-y-1">
                    <Label htmlFor="orgName">Organization Name</Label>
                    <Input
                      id="orgName"
                      defaultValue={org.name}
                      onChange={(e) => setEditName(e.target.value)}
                      className="text-xs"
                    />
                  </div>

                  <div className="space-y-1">
                    <Label htmlFor="orgPlan">Subscription Tier</Label>
                    <Select
                      defaultValue={org.plan || "smart"}
                      onValueChange={setEditPlan}
                    >
                      <SelectTrigger id="orgPlan" className="text-xs">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="smart">Smart (Starter)</SelectItem>
                        <SelectItem value="optimize">Optimize</SelectItem>
                        <SelectItem value="pro">Pro</SelectItem>
                        <SelectItem value="enterprise">Enterprise</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>

                  <div className="space-y-1">
                    <Label htmlFor="domain">Custom Domain</Label>
                    <Input
                      id="domain"
                      placeholder="portal.apex.org"
                      defaultValue={org.customDomain || ""}
                      onChange={(e) => setEditCustomDomain(e.target.value)}
                      className="text-xs font-mono"
                    />
                  </div>

                  <Button
                    type="submit"
                    size="sm"
                    disabled={updateOrgMutation.isPending}
                    className="bg-primary text-primary-foreground text-xs mt-2"
                  >
                    {updateOrgMutation.isPending ? "Saving..." : "Save Changes"}
                  </Button>
                </form>
              </CardContent>
            </Card>

            <Card className="card-enterprise">
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle className="text-sm font-semibold">Corporate Profile & Tax Metadata</CardTitle>
                    <CardDescription className="text-xs">
                      Business credentials, legal registration, and tax filings.
                    </CardDescription>
                  </div>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={handleOpenEditMetadata}
                    className="text-xs h-7 gap-1.5 border-primary/30 text-primary hover:bg-primary/5 shadow-xs"
                  >
                    <Pencil className="w-3 h-3" />
                    Edit Metadata
                  </Button>
                </div>
              </CardHeader>
              <CardContent className="space-y-2.5 text-xs">
                <div className="flex justify-between border-b border-border/50 pb-2">
                  <span className="text-muted-foreground">CAC Registration:</span>
                  <span className="font-mono font-medium">{settings?.cacNumber || org?.registrationNumber || "—"}</span>
                </div>
                <div className="flex justify-between border-b border-border/50 pb-2">
                  <span className="text-muted-foreground">TIN / Tax ID:</span>
                  <span className="font-mono font-medium">{settings?.tinNumber || settings?.taxNumber || org?.taxId || "—"}</span>
                </div>
                <div className="flex justify-between border-b border-border/50 pb-2">
                  <span className="text-muted-foreground">Business Type:</span>
                  <span className="font-medium capitalize">{settings?.businessType || "Healthcare Provider"}</span>
                </div>
                <div className="flex justify-between border-b border-border/50 pb-2">
                  <span className="text-muted-foreground">Support Email:</span>
                  <span className="font-medium">{settings?.supportEmail || org?.email || "—"}</span>
                </div>
                <div className="flex justify-between border-b border-border/50 pb-2">
                  <span className="text-muted-foreground">Support Phone:</span>
                  <span className="font-medium">{settings?.supportPhone || org?.phone || "—"}</span>
                </div>
                <div className="flex justify-between border-b border-border/50 pb-2">
                  <span className="text-muted-foreground">Headquarters Address:</span>
                  <span className="font-medium text-right max-w-[200px] truncate" title={settings?.address || org?.address || ""}>
                    {settings?.address || org?.address || "—"}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Base Currency:</span>
                  <span className="font-mono font-medium">{settings?.currency || "NGN"}</span>
                </div>
              </CardContent>
            </Card>

            {/* Edit Corporate Profile & Tax Metadata Dialog */}
            <Dialog open={openEditMetadata} onOpenChange={setOpenEditMetadata}>
              <DialogContent className="sm:max-w-lg max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                  <DialogTitle className="text-base font-bold flex items-center gap-2">
                    <Building2 className="w-4 h-4 text-primary" />
                    Edit Corporate Profile & Tax Metadata
                  </DialogTitle>
                  <DialogDescription className="text-xs">
                    Update legal entity registration, tax identification numbers, and contact channels for {org.name}.
                  </DialogDescription>
                </DialogHeader>

                <form onSubmit={handleUpdateMetadata} className="space-y-3.5 py-2 text-xs">
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                    <div className="space-y-1">
                      <Label htmlFor="cacNumber" className="text-xs font-semibold">
                        CAC Registration Number
                      </Label>
                      <Input
                        id="cacNumber"
                        placeholder="e.g. RC-1928374 or BN-982374"
                        value={editCacNumber}
                        onChange={(e) => setEditCacNumber(e.target.value)}
                        className="text-xs font-mono"
                      />
                    </div>

                    <div className="space-y-1">
                      <Label htmlFor="tinNumber" className="text-xs font-semibold">
                        TIN / Tax Identification ID
                      </Label>
                      <Input
                        id="tinNumber"
                        placeholder="e.g. 23891048-0001"
                        value={editTinNumber}
                        onChange={(e) => setEditTinNumber(e.target.value)}
                        className="text-xs font-mono"
                      />
                    </div>
                  </div>

                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                    <div className="space-y-1">
                      <Label htmlFor="businessType" className="text-xs font-semibold">
                        Business Entity Type
                      </Label>
                      <Select value={editBusinessType} onValueChange={setEditBusinessType}>
                        <SelectTrigger id="businessType" className="text-xs">
                          <SelectValue placeholder="Select business type" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="Hospital & Medical Center">Hospital & Medical Center</SelectItem>
                          <SelectItem value="Diagnostic & Laboratory Center">Diagnostic & Laboratory Center</SelectItem>
                          <SelectItem value="Pharmacy Network">Pharmacy Network</SelectItem>
                          <SelectItem value="Specialist Clinic">Specialist Clinic</SelectItem>
                          <SelectItem value="Healthcare Provider">Healthcare Provider</SelectItem>
                          <SelectItem value="Research & Pathology Institute">Research & Pathology Institute</SelectItem>
                          <SelectItem value="Public Health Facility">Public Health Facility</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>

                    <div className="space-y-1">
                      <Label htmlFor="currency" className="text-xs font-semibold">
                        Operating Currency
                      </Label>
                      <Select value={editCurrency} onValueChange={setEditCurrency}>
                        <SelectTrigger id="currency" className="text-xs">
                          <SelectValue placeholder="Select currency" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="NGN">NGN - Nigerian Naira (₦)</SelectItem>
                          <SelectItem value="USD">USD - US Dollar ($)</SelectItem>
                          <SelectItem value="GBP">GBP - British Pound (£)</SelectItem>
                          <SelectItem value="EUR">EUR - Euro (€)</SelectItem>
                          <SelectItem value="GHS">GHS - Ghanaian Cedi (₵)</SelectItem>
                          <SelectItem value="KES">KES - Kenyan Shilling (KSh)</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                  </div>

                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                    <div className="space-y-1">
                      <Label htmlFor="supportEmail" className="text-xs font-semibold">
                        Official Support Email
                      </Label>
                      <Input
                        id="supportEmail"
                        type="email"
                        placeholder="support@organization.com"
                        value={editSupportEmail}
                        onChange={(e) => setEditSupportEmail(e.target.value)}
                        className="text-xs"
                      />
                    </div>

                    <div className="space-y-1">
                      <Label htmlFor="supportPhone" className="text-xs font-semibold">
                        Official Support Phone
                      </Label>
                      <Input
                        id="supportPhone"
                        type="tel"
                        placeholder="+234 800 000 0000"
                        value={editSupportPhone}
                        onChange={(e) => setEditSupportPhone(e.target.value)}
                        className="text-xs"
                      />
                    </div>
                  </div>

                  <div className="space-y-1">
                    <Label htmlFor="address" className="text-xs font-semibold">
                      Headquarters Physical Address
                    </Label>
                    <Textarea
                      id="address"
                      placeholder="e.g. Plot 12, Medical District Boulevard, Victoria Island, Lagos"
                      value={editAddress}
                      onChange={(e) => setEditAddress(e.target.value)}
                      rows={2}
                      className="text-xs resize-none"
                    />
                  </div>

                  <DialogFooter className="pt-2 gap-2">
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() => setOpenEditMetadata(false)}
                      className="text-xs"
                    >
                      Cancel
                    </Button>
                    <Button
                      type="submit"
                      size="sm"
                      disabled={updateSettingsMutation.isPending}
                      className="bg-primary text-primary-foreground text-xs shadow-xs"
                    >
                      {updateSettingsMutation.isPending ? "Saving..." : "Save Corporate Metadata"}
                    </Button>
                  </DialogFooter>
                </form>
              </DialogContent>
            </Dialog>
          </div>
        </TabsContent>

        {/* TAB 2: CAPABILITIES & ENTITLEMENTS */}
        <TabsContent value="capabilities" className="space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-sm font-semibold">Effective Capabilities</h3>
              <p className="text-xs text-muted-foreground">
                Runtime capabilities enabled through subscription tier and commercial add-ons.
              </p>
            </div>

            <Dialog open={openGrant} onOpenChange={setOpenGrant}>
              <DialogTrigger asChild>
                <Button size="sm" className="bg-primary text-primary-foreground text-xs h-8 gap-1.5 shadow-sm">
                  <Plus className="h-3.5 w-3.5" />
                  Grant Capability Add-On
                </Button>
              </DialogTrigger>
              <DialogContent className="sm:max-w-md">
                <DialogHeader>
                  <DialogTitle className="text-base">Grant Capability Add-On</DialogTitle>
                  <DialogDescription className="text-xs">
                    Assign a capability add-on or administrative override to this organization.
                  </DialogDescription>
                </DialogHeader>
                <form onSubmit={handleGrantCapability} className="space-y-3 py-2 text-xs">
                  <div className="space-y-1">
                    <Label htmlFor="capSelect">Select Capability</Label>
                    <Select value={selectedCapCode} onValueChange={setSelectedCapCode}>
                      <SelectTrigger id="capSelect" className="text-xs">
                        <SelectValue placeholder="Choose capability code" />
                      </SelectTrigger>
                      <SelectContent>
                        {(catalog || []).map((c) => {
                          const capCode = c.code || c.capability?.code || "";
                          const capName = c.name || c.capability?.name || capCode;
                          return (
                            <SelectItem key={capCode} value={capCode}>
                              {capName} ({capCode})
                            </SelectItem>
                          );
                        })}
                      </SelectContent>
                    </Select>
                  </div>

                  <div className="space-y-1">
                    <Label htmlFor="sourceSelect">Entitlement Source</Label>
                    <Select value={grantSource} onValueChange={setGrantSource}>
                      <SelectTrigger id="sourceSelect" className="text-xs">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="addon">Commercial Add-On</SelectItem>
                        <SelectItem value="override">Platform Administrative Override</SelectItem>
                        <SelectItem value="trial">Promotional Trial</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>

                  <DialogFooter className="pt-3">
                    <Button type="button" variant="outline" size="sm" onClick={() => setOpenGrant(false)} className="text-xs">
                      Cancel
                    </Button>
                    <Button
                      type="submit"
                      size="sm"
                      disabled={grantCapMutation.isPending}
                      className="bg-primary text-primary-foreground text-xs"
                    >
                      {grantCapMutation.isPending ? "Granting..." : "Grant Entitlement"}
                    </Button>
                  </DialogFooter>
                </form>
              </DialogContent>
            </Dialog>
          </div>

          <div className="grid gap-3 sm:grid-cols-2 md:grid-cols-3">
            {(capsData?.capabilities || []).map((cap) => (
              <div
                key={cap}
                className="flex items-center justify-between p-3 rounded-lg border border-border bg-card shadow-sm text-xs"
              >
                <div className="flex items-center gap-2">
                  <CheckCircle2 className="h-4 w-4 text-emerald-500 shrink-0" />
                  <span className="font-mono font-medium text-foreground">{cap}</span>
                </div>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => handleRevokeCapability(cap)}
                  className="h-7 w-7 text-muted-foreground hover:text-destructive"
                >
                  <Trash2 className="h-3.5 w-3.5" />
                </Button>
              </div>
            ))}
          </div>
        </TabsContent>

        {/* TAB 3: COMPLIANCE DOCUMENTS */}
        <TabsContent value="compliance" className="space-y-4">
          <Card className="card-enterprise overflow-hidden">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-semibold">Verification & Legal Documents</CardTitle>
              <CardDescription className="text-xs">
                CAC certificates, medical licenses, and compliance proofs uploaded by the organization.
              </CardDescription>
            </CardHeader>
            <CardContent className="p-0">
              <table className="w-full text-left text-xs">
                <thead className="border-y border-border bg-muted/40 font-semibold uppercase tracking-wider text-muted-foreground">
                  <tr>
                    <th className="py-2.5 px-4">Document Type</th>
                    <th className="py-2.5 px-4">File Name</th>
                    <th className="py-2.5 px-4">Status</th>
                    <th className="py-2.5 px-4">Uploaded At</th>
                    <th className="py-2.5 px-4 text-right">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {!documents || documents.length === 0 ? (
                    <tr>
                      <td colSpan={5} className="py-8 text-center text-muted-foreground">
                        No compliance documents uploaded yet.
                      </td>
                    </tr>
                  ) : (
                    documents.map((doc) => {
                      const displayFilename = doc.originalFilename || doc.fileName || "Document File";
                      const displaySize = doc.fileSizeBytes ?? doc.fileSize ?? 0;
                      return (
                        <tr key={doc.id} className="hover:bg-muted/20">
                          <td className="py-3 px-4 font-semibold capitalize">
                            {doc.documentType.replace(/_/g, " ")}
                          </td>
                          <td className="py-3 px-4 text-muted-foreground font-mono">
                            <span className="truncate max-w-[220px] inline-block align-bottom font-medium text-foreground">{displayFilename}</span>{" "}
                            <span className="text-[11px] text-muted-foreground">({formatFileSize(displaySize)})</span>
                          </td>
                          <td className="py-3 px-4">
                            {doc.status === "approved" ? (
                              <Badge variant="outline" className="border-emerald-500/30 text-emerald-600 bg-emerald-500/5 text-[10px]">
                                Approved
                              </Badge>
                            ) : doc.status === "rejected" ? (
                              <Badge variant="outline" className="border-destructive/30 text-destructive bg-destructive/5 text-[10px]">
                                Rejected
                              </Badge>
                            ) : (
                              <Badge variant="outline" className="border-amber-500/30 text-amber-600 bg-amber-500/5 text-[10px]">
                                Pending Review
                              </Badge>
                            )}
                          </td>
                          <td className="py-3 px-4 text-muted-foreground">
                            {formatDate(doc.uploadedAt || doc.createdAt)}
                          </td>
                        <td className="py-3 px-4 text-right">
                          <div className="flex items-center justify-end gap-1.5">
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => handleOpenPreview(doc)}
                              className="h-7 text-xs gap-1 text-primary hover:bg-primary/10"
                              title="Preview document in-app"
                            >
                              <Eye className="w-3.5 h-3.5" />
                              Preview
                            </Button>
                            <Button
                              variant="ghost"
                              size="icon"
                              asChild
                              className="h-7 w-7 text-muted-foreground hover:text-foreground"
                              title="Download original file"
                            >
                              <a
                                href={`/api/v1/organizations/${orgId}/documents/${doc.id}/download`}
                                download={doc.originalFilename || doc.fileName || "document"}
                              >
                                <Download className="w-3.5 h-3.5" />
                              </a>
                            </Button>
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => setReviewDocId(doc.id)}
                              className="h-7 text-xs"
                            >
                              Review
                            </Button>
                          </div>
                        </td>
                      </tr>
                    );
                  })
                )}
                </tbody>
              </table>
            </CardContent>
          </Card>

          {/* Document Review Modal */}
          <Dialog open={!!reviewDocId} onOpenChange={(open) => !open && setReviewDocId(null)}>
            <DialogContent className="sm:max-w-md">
              <DialogHeader>
                <DialogTitle className="text-base font-semibold">Review Compliance Document</DialogTitle>
                <DialogDescription className="text-xs">
                  Set verification decision for uploaded certificate.
                </DialogDescription>
              </DialogHeader>
              <form onSubmit={handleReviewDocSubmit} className="space-y-3 py-2 text-xs">
                <div className="space-y-1">
                  <Label htmlFor="docStatus">Decision</Label>
                  <Select
                    value={docReviewStatus}
                    onValueChange={(val: any) => setDocReviewStatus(val)}
                  >
                    <SelectTrigger id="docStatus" className="text-xs">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="approved">Approve Document</SelectItem>
                      <SelectItem value="rejected">Reject Document</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                {docReviewStatus === "rejected" && (
                  <div className="space-y-1">
                    <Label htmlFor="docRejReason">Rejection Reason</Label>
                    <Input
                      id="docRejReason"
                      placeholder="e.g. Expired medical practicing license"
                      value={docRejectReason}
                      onChange={(e) => setDocRejectReason(e.target.value)}
                      required
                      className="text-xs"
                    />
                  </div>
                )}

                <DialogFooter className="pt-3">
                  <Button type="button" variant="outline" size="sm" onClick={() => setReviewDocId(null)} className="text-xs">
                    Cancel
                  </Button>
                  <Button type="submit" size="sm" className="bg-primary text-primary-foreground text-xs">
                    Save Review
                  </Button>
                </DialogFooter>
              </form>
            </DialogContent>
          </Dialog>

          {/* In-App Document Previewer */}
          <DocumentPreviewViewer
            open={isPreviewOpen}
            onClose={() => {
              setIsPreviewOpen(false);
              setPreviewDoc(null);
            }}
            documentId={previewDoc?.id}
            organizationId={orgId}
            title={previewDoc?.documentType?.replace(/_/g, " ")}
            fileName={previewDoc?.originalFilename || previewDoc?.fileName}
            mimeType={previewDoc?.mimeType}
            fileSizeBytes={previewDoc?.fileSizeBytes || previewDoc?.fileSize}
            version={previewDoc?.version}
          />
        </TabsContent>
      </Tabs>
    </div>
  );
}
