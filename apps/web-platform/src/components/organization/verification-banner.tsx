import React, { useState } from "react";
import { Link, useLocation } from "react-router-dom";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { useOrganizationDocuments } from "@/api/hooks/use-organizations";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  AlertTriangle,
  Clock,
  CheckCircle2,
  XCircle,
  ArrowRight,
  ShieldCheck,
  FileCheck,
  ChevronDown,
  ChevronUp,
  X,
} from "lucide-react";

const MANDATORY_DOCUMENT_TYPES = [
  "registration_certificate",
  "operating_license",
  "medical_license",
];

export function VerificationBanner() {
  const location = useLocation();
  const { data: bootstrap } = useBootstrap();
  const [isExpanded, setIsExpanded] = useState(false);
  const [isDismissed, setIsDismissed] = useState(false);

  // Only display in Organization context (hide on /platform/* console)
  const isPlatform = location.pathname.startsWith("/platform");
  const org = bootstrap?.organization;
  const orgId = org?.id || "";
  const status = (org?.status || "").toLowerCase();
  const setupState = (org?.setupState || org?.setup_state || "").toUpperCase();

  const { data: documents } = useOrganizationDocuments(orgId);

  // 1. If no organization loaded, on platform console, or organization is active / verified / approved, DO NOT SHOW
  if (
    !org ||
    !orgId ||
    isPlatform ||
    status === "active" ||
    status === "approved" ||
    status === "verified" ||
    setupState === "VERIFIED"
  ) {
    return null;
  }

  // 2. Check if all mandatory regulatory documents are approved
  const docsList = documents || [];
  const approvedMandatoryCount = MANDATORY_DOCUMENT_TYPES.filter((reqType) =>
    docsList.some(
      (d: any) =>
        (d.documentType === reqType || d.document_type === reqType) &&
        (d.status === "approved" || d.status === "active")
    )
  ).length;

  // If all mandatory documents are verified & approved, DO NOT SHOW
  if (approvedMandatoryCount >= MANDATORY_DOCUMENT_TYPES.length) {
    return null;
  }

  // 3. Persistent dismissal per org and status within user session
  const dismissalKey = `curexal_dismissed_banner_${orgId}_${status}_${approvedMandatoryCount}`;
  const isSessionDismissed =
    typeof window !== "undefined" &&
    window.sessionStorage.getItem(dismissalKey) === "true";

  if (isDismissed || isSessionDismissed) {
    return null;
  }

  const handleDismiss = () => {
    setIsDismissed(true);
    if (typeof window !== "undefined") {
      window.sessionStorage.setItem(dismissalKey, "true");
    }
  };

  const isCompliancePage = location.pathname === "/organization/compliance";

  // Check if documents are under review
  const hasPendingReviewDocs = docsList.some(
    (d: any) =>
      d.status === "pending" ||
      d.status === "review_pending" ||
      d.status === "under_review"
  );
  const hasRejectedDocs = docsList.some((d: any) => d.status === "rejected");

  // 1. Rejected State
  if (status === "rejected" || hasRejectedDocs) {
    if (!isExpanded) {
      return (
        <div className="relative overflow-hidden bg-rose-500/10 border-b border-rose-500/25 px-4 py-1.5 text-xs text-foreground transition-all duration-200">
          <div className="max-w-7xl mx-auto flex items-center justify-between gap-2">
            <div className="flex items-center gap-2 min-w-0">
              <XCircle className="w-3.5 h-3.5 text-rose-500 shrink-0" />
              <span className="font-semibold text-rose-600 dark:text-rose-400 truncate">
                Accreditation Verification Rejected
              </span>
              <Badge variant="outline" className="hidden sm:inline-flex border-rose-500/40 text-rose-600 dark:text-rose-400 bg-rose-500/10 text-[9px] uppercase font-mono py-0 px-1.5">
                Action Required
              </Badge>
            </div>
            <div className="flex items-center gap-2 shrink-0">
              <button
                onClick={() => setIsExpanded(true)}
                className="text-[11px] font-medium text-rose-600 dark:text-rose-400 hover:underline flex items-center gap-1 cursor-pointer"
              >
                <span>Details</span>
                <ChevronDown className="w-3 h-3" />
              </button>
              {!isCompliancePage && (
                <Button asChild size="sm" className="text-xs h-6 px-2 gap-1 bg-rose-600 hover:bg-rose-700 text-white shadow-xs">
                  <Link to="/organization/compliance">
                    <FileCheck className="w-3 h-3" />
                    Fix Docs
                  </Link>
                </Button>
              )}
              <button
                onClick={handleDismiss}
                className="p-1 rounded-md text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                title="Dismiss notification"
              >
                <X className="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>
      );
    }

    return (
      <div className="relative overflow-hidden bg-gradient-to-r from-rose-500/15 via-rose-500/10 to-card border-b border-rose-500/30 px-4 py-3 text-xs text-foreground transition-all duration-200">
        <div className="max-w-7xl mx-auto flex flex-col sm:flex-row sm:items-center justify-between gap-3">
          <div className="flex items-start gap-3">
            <div className="p-1.5 rounded-lg bg-rose-500/20 text-rose-500 shrink-0 mt-0.5 sm:mt-0">
              <XCircle className="w-4 h-4" />
            </div>
            <div>
              <div className="flex items-center gap-2">
                <span className="font-bold text-rose-600 dark:text-rose-400">
                  Accreditation Verification Rejected
                </span>
                <Badge variant="outline" className="border-rose-500/40 text-rose-600 dark:text-rose-400 bg-rose-500/10 text-[9px] uppercase font-mono">
                  Action Required
                </Badge>
              </div>
              <p className="text-muted-foreground mt-0.5 text-[11px]">
                Your compliance documents were rejected during regulatory audit. Please review the inspector's notes and upload revised copies to activate your organization.
              </p>
            </div>
          </div>

          <div className="flex items-center gap-2 shrink-0">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setIsExpanded(false)}
              className="text-xs h-7 gap-1 text-muted-foreground hover:text-foreground"
            >
              <ChevronUp className="w-3.5 h-3.5" />
              Hide
            </Button>
            {!isCompliancePage && (
              <Button asChild size="sm" className="text-xs h-7 gap-1.5 bg-rose-600 hover:bg-rose-700 text-white shadow-sm">
                <Link to="/organization/compliance">
                  <FileCheck className="w-3.5 h-3.5" />
                  Fix Documents
                  <ArrowRight className="w-3 h-3" />
                </Link>
              </Button>
            )}
            <button
              onClick={handleDismiss}
              className="p-1 rounded-md text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
              title="Dismiss notification"
            >
              <X className="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
      </div>
    );
  }

  // 2. Under Review State
  if (
    status === "review_pending" ||
    status === "under_review" ||
    status === "submitted" ||
    hasPendingReviewDocs
  ) {
    if (!isExpanded) {
      return (
        <div className="relative overflow-hidden bg-blue-500/10 border-b border-blue-500/25 px-4 py-1.5 text-xs text-foreground transition-all duration-200">
          <div className="max-w-7xl mx-auto flex items-center justify-between gap-2">
            <div className="flex items-center gap-2 min-w-0">
              <Clock className="w-3.5 h-3.5 text-blue-500 shrink-0 animate-pulse" />
              <span className="font-semibold text-blue-600 dark:text-blue-400 truncate">
                Verification In Review
              </span>
              <Badge variant="outline" className="hidden sm:inline-flex border-blue-500/40 text-blue-600 dark:text-blue-400 bg-blue-500/10 text-[9px] uppercase font-mono py-0 px-1.5">
                Pending Audit ({approvedMandatoryCount}/3)
              </Badge>
            </div>
            <div className="flex items-center gap-2 shrink-0">
              <button
                onClick={() => setIsExpanded(true)}
                className="text-[11px] font-medium text-blue-600 dark:text-blue-400 hover:underline flex items-center gap-1 cursor-pointer"
              >
                <span>Details</span>
                <ChevronDown className="w-3 h-3" />
              </button>
              {!isCompliancePage && (
                <Button asChild variant="outline" size="sm" className="text-xs h-6 px-2 gap-1 border-blue-500/40 text-blue-600 dark:text-blue-400 hover:bg-blue-500/10">
                  <Link to="/organization/compliance">
                    <FileCheck className="w-3 h-3" />
                    Status
                  </Link>
                </Button>
              )}
              <button
                onClick={handleDismiss}
                className="p-1 rounded-md text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                title="Dismiss notification"
              >
                <X className="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>
      );
    }

    return (
      <div className="relative overflow-hidden bg-gradient-to-r from-blue-500/15 via-blue-500/10 to-card border-b border-blue-500/30 px-4 py-3 text-xs text-foreground transition-all duration-200">
        <div className="max-w-7xl mx-auto flex flex-col sm:flex-row sm:items-center justify-between gap-3">
          <div className="flex items-start gap-3">
            <div className="p-1.5 rounded-lg bg-blue-500/20 text-blue-500 shrink-0 mt-0.5 sm:mt-0">
              <Clock className="w-4 h-4 animate-pulse" />
            </div>
            <div>
              <div className="flex items-center gap-2">
                <span className="font-bold text-blue-600 dark:text-blue-400">
                  Verification In Review
                </span>
                <Badge variant="outline" className="border-blue-500/40 text-blue-600 dark:text-blue-400 bg-blue-500/10 text-[9px] uppercase font-mono">
                  Pending Audit ({approvedMandatoryCount}/3 Approved)
                </Badge>
              </div>
              <p className="text-muted-foreground mt-0.5 text-[11px]">
                Your compliance credentials have been submitted and are currently undergoing regulatory review by the Platform Compliance Team. Estimated turnaround: 24–48 hours.
              </p>
            </div>
          </div>

          <div className="flex items-center gap-2 shrink-0">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setIsExpanded(false)}
              className="text-xs h-7 gap-1 text-muted-foreground hover:text-foreground"
            >
              <ChevronUp className="w-3.5 h-3.5" />
              Hide
            </Button>
            {!isCompliancePage && (
              <Button asChild variant="outline" size="sm" className="text-xs h-7 gap-1.5 border-blue-500/40 text-blue-600 dark:text-blue-400 hover:bg-blue-500/10">
                <Link to="/organization/compliance">
                  <FileCheck className="w-3.5 h-3.5" />
                  View Status
                </Link>
              </Button>
            )}
            <button
              onClick={handleDismiss}
              className="p-1 rounded-md text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
              title="Dismiss notification"
            >
              <X className="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
      </div>
    );
  }

  // 3. Pending / Draft / Unverified Initial State
  if (!isExpanded) {
    return (
      <div className="relative overflow-hidden bg-amber-500/10 border-b border-amber-500/25 px-4 py-1.5 text-xs text-foreground transition-all duration-200">
        <div className="max-w-7xl mx-auto flex items-center justify-between gap-2">
          <div className="flex items-center gap-2 min-w-0">
            <AlertTriangle className="w-3.5 h-3.5 text-amber-500 shrink-0" />
            <span className="font-semibold text-amber-600 dark:text-amber-400 truncate">
              Regulatory Verification Required
            </span>
            <Badge variant="outline" className="hidden sm:inline-flex border-amber-500/40 text-amber-600 dark:text-amber-400 bg-amber-500/10 text-[9px] uppercase font-mono font-bold py-0 px-1.5">
              {approvedMandatoryCount}/3 Licenses Uploaded
            </Badge>
          </div>
          <div className="flex items-center gap-2 shrink-0">
            <button
              onClick={() => setIsExpanded(true)}
              className="text-[11px] font-medium text-amber-600 dark:text-amber-400 hover:underline flex items-center gap-1 cursor-pointer"
            >
              <span>Details</span>
              <ChevronDown className="w-3 h-3" />
            </button>
            {!isCompliancePage && (
              <Button asChild size="sm" className="text-xs h-6 px-2 gap-1 bg-amber-600 hover:bg-amber-700 text-white shadow-xs">
                <Link to="/organization/compliance">
                  <ShieldCheck className="w-3 h-3" />
                  Verify
                </Link>
              </Button>
            )}
            <button
              onClick={handleDismiss}
              className="p-1 rounded-md text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
              title="Dismiss notification"
            >
              <X className="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="relative overflow-hidden bg-gradient-to-r from-amber-500/15 via-amber-500/10 to-card border-b border-amber-500/30 px-4 py-3 text-xs text-foreground transition-all duration-200">
      <div className="max-w-7xl mx-auto flex flex-col sm:flex-row sm:items-center justify-between gap-3">
        <div className="flex items-start gap-3">
          <div className="p-1.5 rounded-lg bg-amber-500/20 text-amber-500 shrink-0 mt-0.5 sm:mt-0">
            <AlertTriangle className="w-4 h-4" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <span className="font-bold text-amber-600 dark:text-amber-400">
                Regulatory Verification Required
              </span>
              <Badge variant="outline" className="border-amber-500/40 text-amber-600 dark:text-amber-400 bg-amber-500/10 text-[9px] uppercase font-mono font-bold">
                {approvedMandatoryCount}/3 Licenses Uploaded
              </Badge>
            </div>
            <p className="text-muted-foreground mt-0.5 text-[11px]">
              To comply with health authority regulations and unlock operational clinical & diagnostic workspaces, please upload your operating licenses and legal registration certificates.
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2 shrink-0">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setIsExpanded(false)}
            className="text-xs h-7 gap-1 text-muted-foreground hover:text-foreground"
          >
            <ChevronUp className="w-3.5 h-3.5" />
            Hide
          </Button>
          {!isCompliancePage && (
            <Button asChild size="sm" className="text-xs h-7 gap-1.5 bg-amber-600 hover:bg-amber-700 text-white shadow-sm">
              <Link to="/organization/compliance">
                <ShieldCheck className="w-3.5 h-3.5" />
                Complete Verification
                <ArrowRight className="w-3 h-3" />
              </Link>
            </Button>
          )}
          <button
            onClick={handleDismiss}
            className="p-1 rounded-md text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
            title="Dismiss notification"
          >
            <X className="w-3.5 h-3.5" />
          </button>
        </div>
      </div>
    </div>
  );
}

