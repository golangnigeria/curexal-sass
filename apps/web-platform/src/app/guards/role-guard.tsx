import React from "react";
import { Navigate } from "react-router-dom";
import { authClient } from "@/lib/auth-client";

export interface RoleGuardProps {
  allowedRoles: string[];
  fallbackPath?: string;
  children: React.ReactNode;
}

export const RoleGuard: React.FC<RoleGuardProps> = ({
  allowedRoles,
  fallbackPath = "/login",
  children,
}) => {
  const { data: session } = authClient.useSession();
  const userRole = (session?.user?.role || "").toLowerCase();
  const isSuperAdmin = session?.user?.isPlatformAdmin === true;

  if (isSuperAdmin || allowedRoles.map((r) => r.toLowerCase()).includes(userRole)) {
    return <>{children}</>;
  }

  return <Navigate to={fallbackPath} replace />;
};
