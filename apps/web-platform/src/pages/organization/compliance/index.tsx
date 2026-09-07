import React, { useState } from "react";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { useOrganizationDocuments } from "@/api/hooks/use-organizations";
import { organizationService } from "@/api/services/organization.service";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
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
import { formatDate, formatFileSize } from "@curexal/utils";
import { toast } from "sonner";
import { DocumentPreviewViewer } from "@curexal/documents";
import {
  ShieldCheck,
  FileText,
  Upload,
  CheckCircle2,
  XCircle,
  Clock,
  ExternalLink,
  Plus,
  AlertTriangle,
  Send,
  Loader2,
  FileCheck,
  Building2,
  Award,
  Sparkles,
  Info,
  Eye,
  Download,
} from "lucide-react";

const regulatoryDocumentCatalog = [
  {
    type: "registration_certificate",
    name: "Certificate of Incorporation / CAC Registration",
    authority: "Corporate Affairs Commission (CAC) / Registrar General",
    desc: "Official legal entity registration document proving incorporation as a licensed business entity.",
    required: true,
  },
  {
    type: "operating_license",
    name: "State Health Facility Operating License",
    authority: "Ministry of Health / HEFAMAA / State Health Inspectorate",
    desc: "Annual operating permit authorizing healthcare delivery at the registered physical facility premises.",
    required: true,
  },
  {
    type: "medical_license",
    name: "Medical Superintendent / Practicing License",
    authority: "MDCN / MLSCN / PCN / Radiographers Registration Board",
    desc: "Valid annual practicing certificate of the chief medical officer, lab director, or superintendent pharmacist.",
    required: true,
  },
  {
    type: "facility_license",
    name: "Clinical Premises & Inpatient Facility Permit",
    authority: "State Environmental & Health Protection Agency",
    desc: "Premises inspection certificate for clinical wards, emergency rooms, or outpatient consulting suites.",
    required: false,
  },
  {
    type: "accreditation_certificate",
    name: "ISO 15189 / Radiation Safety Accreditation",
    authority: "NNRA (Nuclear Regulatory) / MLSCN / International Accreditation",
    desc: "Tamper-evident accreditation for diagnostic pathology laboratories and radiology X-ray/CT modalities.",
    required: false,
  },
];

export default function OrganizationCompliancePage() {
  const { data: bootstrap, refetch: refetchBootstrap } = useBootstrap();
  const orgId = bootstrap?.organization?.id || "";
  const orgStatus = (bootstrap?.organization?.status || "pending").toLowerCase();
  const orgName = bootstrap?.organization?.name || "Healthcare Organization";

  const { data: documents, isLoading, refetch: refetchDocs } = useOrganizationDocuments(orgId);

  // Upload Modal State
  const [isUploadOpen, setIsUploadOpen] = useState(false);
  const [selectedDocType, setSelectedDocType] = useState("registration_certificate");
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isUploading, setIsUploading] = useState(false);

  // Submit for Review State
  const [isSubmittingReview, setIsSubmittingReview] = useState(false);

  // In-App Document Previewer State
  const [previewDoc, setPreviewDoc] = useState<any | null>(null);
  const [isPreviewOpen, setIsPreviewOpen] = useState(false);

  const handleOpenPreview = (doc: any) => {
    setPreviewDoc(doc);
    setIsPreviewOpen(true);
  };

  const docsList = documents || [];
  const approvedDocs = docsList.filter((d: any) => d.status === "approved" || d.status === "active");
  const pendingDocs = docsList.filter((d: any) => d.status === "pending" || d.status === "review_pending" || d.status === "under_review");
  const rejectedDocs = docsList.filter((d: any) => d.status === "rejected");

  const requiredCatalog = regulatoryDocumentCatalog.filter((r) => r.required);
  const requiredCount = requiredCatalog.length;
  const approvedRequiredCount = requiredCatalog.filter((r) =>
    docsList.some((d: any) => d.documentType === r.type && (d.status === "approved" || d.status === "active"))
  ).length;

  const isVerified = orgStatus === "active" || orgStatus === "approved";
  const isUnderReview = orgStatus === "review_pending" || orgStatus === "under_review";
  const isRejected = orgStatus === "rejected";

  const getDocTypeName = (typeCode: string) => {
    const found = regulatoryDocumentCatalog.find((c) => c.type === typeCode);
    if (found) return found.name;
    return (typeCode || "License").replace(/_/g, " ");
  };

  const getCleanPreviewUrl = (url?: string) => {
    if (!url) return "";
    if (url.startsWith("http://localhost:8080")) {
      return url.replace("http://localhost:8080", "");
    }
    return url;
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      if (file.size > 15 * 1024 * 1024) {
        toast.error("File exceeds 15MB size limit.");
        return;
      }
      setSelectedFile(file);
    }
  };

  const handleUpload = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedFile) {
      toast.error("Please select a valid document file.");
      return;
    }

    setIsUploading(true);
    try {
      const formData = new FormData();
      formData.append("document_type", selectedDocType);
      formData.append("file", selectedFile);

      const res = await fetch(`/api/v1/organizations/${orgId}/documents`, {
        method: "POST",
        credentials: "include",
        body: formData,
      });

      if (!res.ok) {
        const err = await res.json().catch(() => ({}));
        throw new Error(err.message || "Failed to upload document");
      }

      toast.success("Regulatory Document Uploaded Successfully!", {
        description: "Your document has been registered in the compliance ledger and queued for verification.",
      });

      setIsUploadOpen(false);
      setSelectedFile(null);
      refetchDocs();
      if (refetchBootstrap) refetchBootstrap();
    } catch (err: any) {
      toast.error(err.message || "Upload failed");
    } finally {
      setIsUploading(false);
    }
  };

  const handleSubmitReview = async () => {
    if (docsList.length === 0) {
      toast.error("Please upload at least one required regulatory document before submitting for review.");
      return;
    }

    setIsSubmittingReview(true);
    try {
      const res = await fetch(`/api/v1/organization/setup/submit-review`, {
        method: "POST",
        credentials: "include",
      });

      if (!res.ok) {
        const err = await res.json().catch(() => ({}));
        throw new Error(err.message || "Failed to submit for review");
      }

      toast.success("Organization Submitted for Compliance Review!", {
        description: "Platform compliance inspectors have been notified. Your status is now under review.",
      });

      if (refetchBootstrap) refetchBootstrap();
      refetchDocs();
    } catch (err: any) {
      toast.error(err.message || "Submission failed");
    } finally {
      setIsSubmittingReview(false);
    }
  };

  const openUploadForType = (typeCode: string) => {
    setSelectedDocType(typeCode);
    setSelectedFile(null);
    setIsUploadOpen(true);
  };

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Header & Main Status */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-6">
        <div>
          <div className="flex items-center gap-2.5 mb-1">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              Healthcare Regulatory Compliance & Accreditation
            </h1>
            {isVerified ? (
              <Badge variant="outline" className="border-emerald-500/40 text-emerald-600 dark:text-emerald-400 bg-emerald-500/10 text-xs font-mono font-bold">
                <CheckCircle2 className="w-3.5 h-3.5 mr-1 inline" /> Verified Healthcare Provider
              </Badge>
            ) : isUnderReview ? (
              <Badge variant="outline" className="border-blue-500/40 text-blue-600 dark:text-blue-400 bg-blue-500/10 text-xs font-mono font-bold">
                <Clock className="w-3.5 h-3.5 mr-1 inline animate-pulse" /> Review In Progress
              </Badge>
            ) : isRejected ? (
              <Badge variant="outline" className="border-rose-500/40 text-rose-600 dark:text-rose-400 bg-rose-500/10 text-xs font-mono font-bold">
                <XCircle className="w-3.5 h-3.5 mr-1 inline" /> Revision Requested
              </Badge>
            ) : (
              <Badge variant="outline" className="border-amber-500/40 text-amber-600 dark:text-amber-400 bg-amber-500/10 text-xs font-mono font-bold">
                <AlertTriangle className="w-3.5 h-3.5 mr-1 inline" /> Verification Pending
              </Badge>
            )}
          </div>
          <p className="text-xs text-muted-foreground">
            Manage your legal entity incorporation, operating permits, clinical licenses, and health authority accreditation files.
          </p>
        </div>

        <div className="flex items-center gap-2.5">
          <Dialog open={isUploadOpen} onOpenChange={setIsUploadOpen}>
            <DialogTrigger asChild>
              <Button size="sm" className="text-xs h-9 gap-1.5 bg-primary text-primary-foreground shadow">
                <Upload className="w-3.5 h-3.5" />
                Upload Regulatory Document
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[480px]">
              <form onSubmit={handleUpload}>
                <DialogHeader>
                  <DialogTitle className="text-base font-bold text-foreground">
                    Upload Regulatory Compliance Document
                  </DialogTitle>
                  <DialogDescription className="text-xs">
                    Upload official government-issued licenses, permits, or registration certificates (PDF, JPEG, PNG max 15MB).
                  </DialogDescription>
                </DialogHeader>

                <div className="space-y-4 py-4 text-xs">
                  <div className="space-y-1.5">
                    <Label className="text-xs font-medium">Document Category</Label>
                    <Select value={selectedDocType} onValueChange={setSelectedDocType}>
                      <SelectTrigger className="text-xs h-9">
                        <SelectValue placeholder="Select Document Type" />
                      </SelectTrigger>
                      <SelectContent>
                        {regulatoryDocumentCatalog.map((cat) => (
                          <SelectItem key={cat.type} value={cat.type} className="text-xs">
                            {cat.name} {cat.required ? "(Required)" : "(Optional)"}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>

                  <div className="space-y-1.5">
                    <Label className="text-xs font-medium">Select Document File</Label>
                    <div className="border-2 border-dashed border-border rounded-xl p-5 text-center hover:border-primary/50 transition-colors bg-card/50">
                      <Input
                        type="file"
                        accept=".pdf,.png,.jpg,.jpeg"
                        onChange={handleFileChange}
                        className="hidden"
                        id="compliance-file-input"
                      />
                      <label htmlFor="compliance-file-input" className="cursor-pointer block space-y-2">
                        <div className="w-10 h-10 rounded-full bg-primary/10 text-primary flex items-center justify-center mx-auto">
                          <Upload className="w-5 h-5" />
                        </div>
                        <div>
                          <p className="text-xs font-bold text-foreground">
                            {selectedFile ? selectedFile.name : "Click to browse or drag and drop"}
                          </p>
                          <p className="text-[11px] text-muted-foreground mt-0.5">
                            {selectedFile
                              ? formatFileSize(selectedFile.size)
                              : "PDF, PNG, JPG up to 15MB (Tamper-evident SHA-256 verified)"}
                          </p>
                        </div>
                      </label>
                    </div>
                  </div>
                </div>

                <DialogFooter>
                  <Button type="button" variant="outline" size="sm" onClick={() => setIsUploadOpen(false)} disabled={isUploading} className="text-xs">
                    Cancel
                  </Button>
                  <Button type="submit" size="sm" disabled={!selectedFile || isUploading} className="text-xs gap-1.5 bg-primary text-primary-foreground shadow">
                    {isUploading ? <Loader2 className="w-3.5 h-3.5 animate-spin" /> : <Upload className="w-3.5 h-3.5" />}
                    {isUploading ? "Uploading..." : "Save & Register"}
                  </Button>
                </DialogFooter>
              </form>
            </DialogContent>
          </Dialog>

          {!isVerified && !isUnderReview && (
            <Button
              size="sm"
              onClick={handleSubmitReview}
              disabled={isSubmittingReview || docsList.length === 0}
              className="text-xs h-9 gap-1.5 bg-emerald-600 hover:bg-emerald-700 text-white shadow"
            >
              {isSubmittingReview ? (
                <Loader2 className="w-3.5 h-3.5 animate-spin" />
              ) : (
                <Send className="w-3.5 h-3.5" />
              )}
              Submit for Verification
            </Button>
          )}
        </div>
      </div>

      {/* Compliance Overview Progress Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card className="border-border shadow-sm bg-card">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-medium">Compliance Readiness Score</CardDescription>
            <CardTitle className="text-2xl font-bold font-mono text-foreground flex items-center justify-between">
              <span>{Math.round((approvedRequiredCount / requiredCount) * 100)}%</span>
              <Award className="w-5 h-5 text-primary" />
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="w-full bg-muted rounded-full h-2 overflow-hidden mb-2">
              <div
                className="bg-primary h-2 rounded-full transition-all duration-500"
                style={{ width: `${(approvedRequiredCount / requiredCount) * 100}%` }}
              />
            </div>
            <p className="text-[11px] text-muted-foreground">
              {approvedRequiredCount} of {requiredCount} mandatory regulatory licenses verified
            </p>
          </CardContent>
        </Card>

        <Card className="border-border shadow-sm bg-card">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-medium">Total Registered Documents</CardDescription>
            <CardTitle className="text-2xl font-bold font-mono text-foreground flex items-center justify-between">
              <span>{docsList.length} Files</span>
              <FileText className="w-5 h-5 text-blue-500" />
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-[11px] text-muted-foreground">
              {approvedDocs.length} Approved • {pendingDocs.length} In Review • {rejectedDocs.length} Rejected
            </p>
          </CardContent>
        </Card>

        <Card className="border-border shadow-sm bg-card">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-medium">Accreditation Level</CardDescription>
            <CardTitle className="text-2xl font-bold font-mono text-foreground flex items-center justify-between">
              <span>{isVerified ? "Level 3" : "Level 1 (Draft)"}</span>
              <ShieldCheck className="w-5 h-5 text-emerald-500" />
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-[11px] text-muted-foreground">
              {isVerified ? "Full Clinical & Diagnostic Clearance" : "Complete document review to upgrade clearance"}
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Required Regulatory Licenses Checklist */}
      <div className="space-y-4">
        <div>
          <h3 className="text-base font-bold text-foreground flex items-center gap-2">
            <ShieldCheck className="w-4 h-4 text-primary" />
            Mandatory Healthcare Licenses & Permits
          </h3>
          <p className="text-xs text-muted-foreground">
            Curexal requires all healthcare providers to upload verifiable documentation for patient safety and regulatory compliance.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {regulatoryDocumentCatalog.map((cat) => {
            const matchingDocs = docsList.filter((d: any) => d.documentType === cat.type);
            matchingDocs.sort((a: any, b: any) => (b.version || 1) - (a.version || 1));
            const uploaded = matchingDocs[0];
            const docStatus = (uploaded?.status || "").toLowerCase();
            const isDocApproved = docStatus === "approved" || docStatus === "active" || docStatus === "verified";
            const isDocPending = docStatus === "pending" || docStatus === "review_pending" || docStatus === "under_review";
            const isDocRejected = docStatus === "rejected";

            return (
              <Card
                key={cat.type}
                className={`border transition-all flex flex-col justify-between ${
                  isDocApproved
                    ? "border-emerald-500/30 bg-emerald-500/5 shadow-sm"
                    : isDocRejected
                    ? "border-rose-500/30 bg-rose-500/5 shadow-sm"
                    : isDocPending
                    ? "border-blue-500/30 bg-blue-500/5 shadow-sm"
                    : "border-border bg-card shadow-sm hover:border-border/80"
                }`}
              >
                <CardHeader className="pb-3">
                  <div className="flex items-start justify-between gap-2">
                    <div>
                      <div className="flex items-center gap-2">
                        <span className="text-xs font-bold text-foreground">{cat.name}</span>
                        {cat.required && (
                          <Badge variant="outline" className="text-[9px] font-mono uppercase bg-primary/10 text-primary border-primary/20">
                            Mandatory
                          </Badge>
                        )}
                      </div>
                      <p className="text-[11px] text-primary/80 font-medium mt-0.5">{cat.authority}</p>
                    </div>

                    {isDocApproved ? (
                      <Badge variant="outline" className="text-[10px] bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/30">
                        <CheckCircle2 className="w-3 h-3 mr-1 inline" /> Verified
                      </Badge>
                    ) : isDocPending ? (
                      <Badge variant="outline" className="text-[10px] bg-blue-500/10 text-blue-600 dark:text-blue-400 border-blue-500/30">
                        <Clock className="w-3 h-3 mr-1 inline" /> In Review
                      </Badge>
                    ) : isDocRejected ? (
                      <Badge variant="outline" className="text-[10px] bg-rose-500/10 text-rose-600 dark:text-rose-400 border-rose-500/30">
                        <XCircle className="w-3 h-3 mr-1 inline" /> Rejected
                      </Badge>
                    ) : (
                      <Badge variant="outline" className="text-[10px] border-border text-muted-foreground">
                        Not Uploaded
                      </Badge>
                    )}
                  </div>
                  <CardDescription className="text-xs mt-2 text-muted-foreground">
                    {cat.desc}
                  </CardDescription>
                </CardHeader>

                <CardContent className="pt-0">
                  <div className="flex items-center justify-between border-t border-border/50 pt-3">
                    {uploaded ? (
                      <div className="text-[11px] text-muted-foreground flex items-center gap-1.5 font-mono truncate max-w-[220px]">
                        <FileText className="w-3.5 h-3.5 text-primary shrink-0" />
                        <span className="truncate">{uploaded.originalFilename}</span>
                        <span className="text-[10px] opacity-70">(v{uploaded.version || 1})</span>
                      </div>
                    ) : (
                      <span className="text-[11px] text-muted-foreground italic">No file registered</span>
                    )}

                    <div className="flex items-center gap-1.5">
                      {uploaded && (
                        <Button
                          size="sm"
                          variant="ghost"
                          onClick={() => handleOpenPreview(uploaded)}
                          className="text-xs h-7 gap-1 text-primary hover:bg-primary/10"
                          title="Preview document in-app"
                        >
                          <Eye className="w-3 h-3" />
                          Preview
                        </Button>
                      )}
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => openUploadForType(cat.type)}
                        className="text-xs h-7 gap-1 hover:border-primary hover:text-primary"
                      >
                        <Upload className="w-3 h-3" />
                        {uploaded ? "Update (v" + ((uploaded.version || 1) + 1) + ")" : "Upload"}
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            );
          })}
        </div>
      </div>

      {/* Document History & Audit Repository */}
      <div className="space-y-4 pt-4 border-t border-border">
        <div>
          <h3 className="text-base font-bold text-foreground flex items-center gap-2">
            <FileCheck className="w-4 h-4 text-primary" />
            Registered Document Audit Vault
          </h3>
          <p className="text-xs text-muted-foreground">
            All documents are stored with AES-256 encryption and tamper-evident SHA-256 cryptographic hashes.
          </p>
        </div>

        {docsList.length === 0 ? (
          <Card className="border-border shadow-sm p-8 text-center text-muted-foreground text-xs bg-card">
            <FileText className="w-8 h-8 mx-auto text-muted-foreground/50 mb-2" />
            No regulatory documents have been uploaded yet. Click "Upload Regulatory Document" above to get started.
          </Card>
        ) : (
          <div className="border border-border rounded-xl overflow-hidden shadow-sm bg-card">
            <table className="w-full text-left text-xs">
              <thead className="bg-muted/50 border-b border-border text-muted-foreground">
                <tr>
                  <th className="p-3 font-medium">Document Type</th>
                  <th className="p-3 font-medium">File Name</th>
                  <th className="p-3 font-medium">Version</th>
                  <th className="p-3 font-medium">Size</th>
                  <th className="p-3 font-medium">Date Uploaded</th>
                  <th className="p-3 font-medium">Verification Status</th>
                  <th className="p-3 font-medium text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {docsList.map((doc: any) => {
                  return (
                    <tr key={doc.id} className="hover:bg-muted/20 transition-colors">
                      <td className="p-3 font-bold text-foreground">
                        {getDocTypeName(doc.documentType)}
                      </td>
                      <td className="p-3 text-muted-foreground font-mono truncate max-w-[200px]">
                        {doc.originalFilename || "Document File"}
                      </td>
                      <td className="p-3 font-mono text-muted-foreground">v{doc.version || 1}</td>
                      <td className="p-3 font-mono text-muted-foreground">{formatFileSize(doc.fileSizeBytes || 0)}</td>
                      <td className="p-3 text-muted-foreground">{formatDate(doc.uploadedAt || doc.createdAt)}</td>
                      <td className="p-3">
                        {doc.status === "approved" || doc.status === "active" ? (
                          <Badge variant="outline" className="text-[10px] bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/30">
                            <CheckCircle2 className="w-3 h-3 mr-1 inline" /> Approved
                          </Badge>
                        ) : doc.status === "rejected" ? (
                          <div>
                            <Badge variant="outline" className="text-[10px] bg-rose-500/10 text-rose-600 dark:text-rose-400 border-rose-500/30">
                              <XCircle className="w-3 h-3 mr-1 inline" /> Rejected
                            </Badge>
                            {doc.rejectionReason && (
                              <p className="text-[10px] text-rose-500 mt-1 max-w-[240px] leading-tight">
                                Note: {doc.rejectionReason}
                              </p>
                            )}
                          </div>
                        ) : (
                          <Badge variant="outline" className="text-[10px] bg-blue-500/10 text-blue-600 dark:text-blue-400 border-blue-500/30">
                            <Clock className="w-3 h-3 mr-1 inline" /> Pending Review
                          </Badge>
                        )}
                      </td>
                      <td className="p-3 text-right">
                        <div className="flex items-center justify-end gap-1.5">
                          <Button
                            size="sm"
                            variant="ghost"
                            onClick={() => handleOpenPreview(doc)}
                            className="h-7 text-xs gap-1.5 text-primary hover:text-primary hover:bg-primary/10 font-medium"
                            title="Preview in-app without downloading"
                          >
                            <Eye className="w-3.5 h-3.5" />
                            Preview
                          </Button>
                          <Button
                            size="icon"
                            variant="ghost"
                            asChild
                            className="h-7 w-7 text-muted-foreground hover:text-foreground"
                            title="Download original document"
                          >
                            <a
                              href={`/api/v1/organizations/${orgId}/documents/${doc.id}/download`}
                              download={doc.originalFilename || "document"}
                            >
                              <Download className="w-3.5 h-3.5" />
                            </a>
                          </Button>
                        </div>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* In-App PDF.js & Image Previewer Modal */}
      <DocumentPreviewViewer
        open={isPreviewOpen}
        onClose={() => {
          setIsPreviewOpen(false);
          setPreviewDoc(null);
        }}
        documentId={previewDoc?.id}
        organizationId={orgId}
        title={getDocTypeName(previewDoc?.documentType)}
        fileName={previewDoc?.originalFilename}
        mimeType={previewDoc?.mimeType}
        fileSizeBytes={previewDoc?.fileSizeBytes}
        version={previewDoc?.version}
      />
    </div>
  );
}
