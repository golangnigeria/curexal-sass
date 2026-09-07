import { apiGet } from "./client";
import type { NavigationResponse, NavigationQueryParams } from "@curexal/contracts";

export { type NavigationQueryParams };

export const navigationService = {
  async getNavigation(params?: NavigationQueryParams): Promise<NavigationResponse> {
    const qs = new URLSearchParams();
    if (params?.scope) qs.set("scope", params.scope);
    if (params?.branch) qs.set("branch", params.branch);
    if (params?.host) qs.set("host", params.host);
    const query = qs.toString();
    return apiGet<NavigationResponse>(`/navigation${query ? `?${query}` : ""}`);
  },
};
