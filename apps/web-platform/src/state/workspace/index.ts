import { useLocation, useParams } from "react-router-dom";
import { resolveEntitledWorkspaces } from "@/lib/permissions/entitlements";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { authClient } from "@/lib/auth-client";

export function useActiveWorkspace() {
  const location = useLocation();
  const { branchSlug } = useParams<{ branchSlug: string }>();
  const { data: bootstrap } = useBootstrap(branchSlug);
  const { data: session } = authClient.useSession();

  const pathParts = location.pathname.split("/").filter(Boolean);
  const activeModuleCode = pathParts[1] || "dashboard";

  const entitledWorkspaces = resolveEntitledWorkspaces(bootstrap, session?.user);
  const activeWorkspace = entitledWorkspaces.find((w) => w.path === activeModuleCode) || null;

  return {
    branchSlug: branchSlug || pathParts[0] || "main",
    activeModuleCode,
    activeWorkspace,
    entitledWorkspaces,
  };
}
