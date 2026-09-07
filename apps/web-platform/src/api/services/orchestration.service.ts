import { apiClient } from "../client";
import type {
  CareJourneyMilestone,
  CareRequest,
  CareRequestFilter,
  CareRequestListResponse,
} from "../contracts";

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
    const res = await apiClient.post<{ success: boolean; data: CareRequest }>(
      "/care-requests",
      payload
    );
    return res.data?.data || res.data;
  },

  /**
   * List Care Requests for staff workspace queue
   */
  async listCareRequests(filter?: CareRequestFilter): Promise<CareRequestListResponse> {
    const res = await apiClient.get<CareRequestListResponse>("/care-requests", {
      params: filter,
    });
    return res.data;
  },

  /**
   * Retrieve single Care Request by ID
   */
  async getCareRequestById(requestId: string): Promise<CareRequest> {
    const res = await apiClient.get<{ success: boolean; data: CareRequest }>(
      `/care-requests/${requestId}`
    );
    return res.data?.data || res.data;
  },

  /**
   * Retrieve Care Requests for the authenticated patient portal user
   */
  async getMyCareRequests(): Promise<CareRequest[]> {
    const res = await apiClient.get<{ success: boolean; data: CareRequest[] }>(
      "/portal/my-care-requests"
    );
    if (Array.isArray(res.data?.data)) return res.data.data;
    if (Array.isArray(res.data)) return res.data;
    return [];
  },

  /**
   * Retrieve live Curexal Care Journey milestones for the authenticated patient
   */
  async getMyCareJourney(): Promise<CareJourneyMilestone[]> {
    const res = await apiClient.get<{ success: boolean; data: CareJourneyMilestone[] }>(
      "/portal/my-care-journey"
    );
    if (Array.isArray(res.data?.data)) return res.data.data;
    if (Array.isArray(res.data)) return res.data;
    return [];
  },

  /**
   * Submit clinical triage vitals & calculate urgency acuity
   */
  async submitTriage(requestId: string, payload: any): Promise<any> {
    const res = await apiClient.post<{ success: boolean; data: any }>(
      `/care-requests/${requestId}/triage`,
      payload
    );
    return res.data?.data || res.data;
  },

  /**
   * Calculate ranked qualified providers for a care request
   */
  async matchProviders(requestId: string): Promise<any> {
    const res = await apiClient.get<{ success: boolean; data: any }>(
      `/care-requests/${requestId}/match-provider`
    );
    return res.data?.data || res.data;
  },

  /**
   * Assign a doctor to a care request and advance Care Journey
   */
  async assignProvider(requestId: string, providerId: string): Promise<void> {
    await apiClient.post(`/care-requests/${requestId}/assign-provider`, {
      providerId,
    });
  },
};
