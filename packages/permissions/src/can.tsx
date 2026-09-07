import React from "react";
import { can, hasRole, hasCapability } from "./evaluator";
import type { UserRoleResponse, BootstrapContractResponse } from "@curexal/contracts";

export interface CanProps {
  permission?: string;
  role?: string;
  capability?: string;
  user?: Partial<UserRoleResponse> | null;
  bootstrap?: Partial<BootstrapContractResponse> | null;
  children: React.ReactNode;
  fallback?: React.ReactNode;
}

/**
 * Declarative UX permission/role gate component.
 * Note: UX only. Backend enforces authoritative authorization.
 */
export function Can({
  permission,
  role,
  capability,
  user,
  bootstrap,
  children,
  fallback = null,
}: CanProps) {
  if (permission && !can(permission, user)) {
    return <>{fallback}</>;
  }

  if (role && !hasRole(role, user)) {
    return <>{fallback}</>;
  }

  if (capability && !hasCapability(capability, bootstrap)) {
    return <>{fallback}</>;
  }

  return <>{children}</>;
}

export interface PermissionGateProps {
  permission: string;
  user?: Partial<UserRoleResponse> | null;
  children: React.ReactNode;
  fallback?: React.ReactNode;
}

export function PermissionGate({
  permission,
  user,
  children,
  fallback = null,
}: PermissionGateProps) {
  if (!can(permission, user)) {
    return <>{fallback}</>;
  }
  return <>{children}</>;
}
