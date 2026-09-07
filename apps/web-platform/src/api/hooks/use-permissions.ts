import { authClient } from "@/lib/auth-client";
import { useBootstrap } from "./use-bootstrap";
import { can as canEvaluate, hasRole as hasRoleEvaluate } from "@curexal/permissions";
import type { UserRoleResponse } from "@/api/contracts";

export function usePermissions() {
  const { data: session } = authClient.useSession();
  const { data: bootstrap } = useBootstrap();

  const user = session?.user as (Partial<UserRoleResponse> & { role?: string; organizationRole?: string }) | undefined;

  const rawRole = (user?.role || bootstrap?.organization?.role || "").toLowerCase();
  const userPermissions = user?.permissions || [];

  const can = (permission: string): boolean => {
    if (user?.isPlatformAdmin || rawRole === "super_admin") return true;
    if (userPermissions.includes("*") || userPermissions.includes(permission)) return true;
    return canEvaluate(permission, user);
  };

  const hasRole = (role: string): boolean => {
    return hasRoleEvaluate(role, user);
  };

  const isClinician = Boolean(
    can("hms.consultation.create") ||
    rawRole === "doctor" ||
    rawRole === "clinician" ||
    rawRole === "specialist"
  );

  const isNurse = Boolean(
    !isClinician && (
      can("hms.vitals.record") ||
      rawRole === "nurse" ||
      rawRole === "triage_nurse"
    )
  );

  const isScientist = Boolean(
    can("lis.result.create") ||
    rawRole === "scientist" ||
    rawRole === "senior_scientist" ||
    rawRole === "lab_technician"
  );

  const isExecutive = Boolean(
    rawRole === "owner" ||
    rawRole === "org_admin" ||
    rawRole === "org_regional_manager" ||
    rawRole === "org_quality_manager" ||
    rawRole === "org_finance_manager" ||
    rawRole === "org_hr_manager" ||
    user?.isPlatformAdmin === true
  );

  const isManager = Boolean(
    isExecutive ||
    rawRole === "branch_manager" ||
    rawRole === "branch_admin" ||
    rawRole === "lab_manager" ||
    rawRole === "pharmacy_manager"
  );

  return {
    user,
    bootstrap,
    can,
    hasRole,
    isClinician,
    isNurse,
    isScientist,
    isExecutive,
    isManager,
    role: rawRole,
    permissions: userPermissions,
  };
}
