import { apiClient } from "@curexal/api-client";
import type {
  CanonicalPatient,
  DuplicateEvaluationRequest,
  DuplicateEvaluationResponse,
  PatientListFilter,
  PatientListResponse,
  PortalAuthResponse,
  RegisterCanonicalPatientPayload,
  SendPortalOTPPayload,
  SetPortalPINPayload,
  VerifyPortalOTPPayload,
} from "@curexal/contracts";

export interface RegisterPatientResult {
  patient: CanonicalPatient;
  token?: string;
  accessToken?: string;
}

export const patientService = {
  /**
   * Evaluate identity signals against existing patients to detect duplicates before registration
   */
  async resolveDuplicates(payload: DuplicateEvaluationRequest): Promise<DuplicateEvaluationResponse> {
    const res = await apiClient.post<{ success: boolean; data: DuplicateEvaluationResponse }>(
      "/patients/resolve",
      payload
    );
    return res.data?.data || res.data;
  },

  /**
   * Register a new canonical patient with MPI duplicate verification & instant session establishment
   */
  async registerPatient(payload: RegisterCanonicalPatientPayload): Promise<RegisterPatientResult> {
    const res = await apiClient.post<{
      success: boolean;
      data: CanonicalPatient;
      token?: string;
      accessToken?: string;
    }>("/patients", payload);
    return {
      patient: res.data?.data || (res.data as any),
      token: res.data?.token || res.data?.accessToken,
      accessToken: res.data?.accessToken || res.data?.token,
    };
  },

  /**
   * List and search patients within the current tenant / branch
   */
  async listPatients(filter?: PatientListFilter): Promise<PatientListResponse> {
    const res = await apiClient.get<PatientListResponse>("/patients", {
      params: filter,
    });
    return res.data;
  },

  /**
   * Retrieve a patient's complete 360 profile
   */
  async getPatientById(patientId: string): Promise<CanonicalPatient> {
    const res = await apiClient.get<{ success: boolean; data: CanonicalPatient }>(
      `/patients/${patientId}`
    );
    return res.data?.data || res.data;
  },

  /**
   * Send a passwordless OTP for patient portal login
   */
  async sendPortalOTP(payload: SendPortalOTPPayload): Promise<{ identifier: string; debugCode?: string }> {
    const res = await apiClient.post<{ success: boolean; data: { identifier: string; debugCode?: string } }>(
      "/portal/auth/send-otp",
      payload
    );
    return res.data?.data || res.data;
  },

  /**
   * Verify passwordless OTP and establish authenticated portal identity session
   */
  async verifyPortalOTP(payload: VerifyPortalOTPPayload): Promise<PortalAuthResponse> {
    const res = await apiClient.post<PortalAuthResponse>(
      "/portal/auth/verify-otp",
      payload
    );
    return res.data;
  },

  /**
   * Set a 4-6 digit quick login PIN for a patient portal account
   */
  async setPortalPIN(patientId: string, payload: SetPortalPINPayload): Promise<void> {
    await apiClient.post(`/portal/auth/patients/${patientId}/pin`, payload);
  },
};
