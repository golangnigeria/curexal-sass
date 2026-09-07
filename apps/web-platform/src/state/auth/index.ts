import { authClient } from "@/lib/auth-client";

export function useAuthSession() {
  return authClient.useSession();
}

export function getCurrentUser(sessionData?: any): any {
  return sessionData?.user || null;
}
