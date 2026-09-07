import { apiGet, apiPost } from "@curexal/api-client";
import type {
  CareJourneyMilestone,
  CareRequest,
  CareRequestFilter,
  CareRequestListResponse,
} from "@curexal/contracts";

export interface CreateCareRequestPayload {
  patientId?: string;
  serviceType: "GENERAL_CONSULTATION" | "SPECIALIST" | "LAB_TEST" | "REFILL" | "TELEHEALTH" | string;
  preferredMode: "IN_PERSON" | "VIDEO" | "AUDIO" | "ASYNC_CHAT" | string;
  urgency?: "ROUTINE" | "URGENT" | "EMERGENCY" | string;
  chiefComplaint: string;
  symptoms?: string[];
  preferredTimeWindow?: Record<string, any>;
}

export const orchestrationService = {
  /**
   * Submit a new Care Request (from patient portal or staff desk)
   */
  async createCareRequest(payload: CreateCareRequestPayload): Promise<CareRequest> {
    const res = await apiPost<{ success: boolean; data: CareRequest }>(
      "/care-requests",
      payload
    );
    return res?.data || (res as any);
  },

  /**
   * List Care Requests for staff workspace queue
   */
  async listCareRequests(filter?: CareRequestFilter): Promise<CareRequestListResponse> {
    return await apiGet<CareRequestListResponse>("/care-requests", {
      params: filter,
    });
  },

  /**
   * Retrieve single Care Request by ID
   */
  async getCareRequestById(requestId: string): Promise<CareRequest> {
    const res = await apiGet<{ success: boolean; data: CareRequest }>(
      `/care-requests/${requestId}`
    );
    return res?.data || (res as any);
  },

  /**
   * Retrieve Care Requests for the authenticated patient portal user
   */
  async getMyCareRequests(): Promise<CareRequest[]> {
    const res = await apiGet<{ success: boolean; data: CareRequest[] }>(
      "/portal/my-care-requests"
    );
    if (Array.isArray(res?.data)) return res.data;
    if (Array.isArray(res)) return res;
    return [];
  },

  /**
   * Retrieve live Curexal Care Journey milestones for the authenticated patient
   */
  async getMyCareJourney(): Promise<CareJourneyMilestone[]> {
    const res = await apiGet<{ success: boolean; data: CareJourneyMilestone[] }>(
      "/portal/my-care-journey"
    );
    if (Array.isArray(res?.data)) return res.data;
    if (Array.isArray(res)) return res;
    return [];
  },

  /**
   * Submit clinical triage vitals & calculate urgency acuity
   */
  async submitTriage(requestId: string, payload: any): Promise<any> {
    const res = await apiPost<{ success: boolean; data: any }>(
      `/care-requests/${requestId}/triage`,
      payload
    );
    return res?.data || (res as any);
  },

  /**
   * Calculate ranked qualified providers for a care request
   */
  async matchProviders(requestId: string): Promise<any> {
    const res = await apiGet<{ success: boolean; data: any }>(
      `/care-requests/${requestId}/match-provider`
    );
    return res?.data || (res as any);
  },

  /**
   * Assign a doctor to a care request and advance Care Journey
   */
  async assignProvider(requestId: string, providerId: string): Promise<void> {
    await apiPost(`/care-requests/${requestId}/assign-provider`, {
      providerId,
    });
  },
};
