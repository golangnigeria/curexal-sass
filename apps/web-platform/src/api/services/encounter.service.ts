import { apiClient } from "../client";
import type {
  ClinicalEncounter,
  ClinicalDiagnosis,
  Prescription,
  StartEncounterPayload,
  UpdateSOAPNotesPayload,
  AddDiagnosisPayload,
  CreatePrescriptionPayload,
  CompleteEncounterResponse,
} from "../contracts";

export interface LabOrderItem {
  testCode: string;
  testName: string;
  urgency: "ROUTINE" | "STAT" | "URGENT" | string;
  notes?: string;
}

export interface RadiologyOrderItem {
  modality: "XRAY" | "CT" | "MRI" | "ULTRASOUND" | string;
  studyName: string;
  urgency: "ROUTINE" | "STAT" | "URGENT" | string;
  notes?: string;
}

export interface PrescriptionOrderItem {
  medicationName: string;
  dosage: string;
  frequency: string;
  durationDays: number;
  instructions?: string;
}

export interface DispatchOrdersPayload {
  labOrders: LabOrderItem[];
  radiologyOrders: RadiologyOrderItem[];
  prescriptionList: PrescriptionOrderItem[];
}

export interface EncounterFilter {
  status?: string;
  providerId?: string;
  patientId?: string;
}

export const encounterService = {
  /**
   * List encounters with optional status and provider filters
   */
  async listEncounters(filter?: EncounterFilter): Promise<ClinicalEncounter[]> {
    const params = new URLSearchParams();
    if (filter?.status) params.append("status", filter.status);
    if (filter?.providerId) params.append("provider_id", filter.providerId);
    if (filter?.patientId) params.append("patient_id", filter.patientId);

    const res = await apiClient.get<{ success: boolean; data: ClinicalEncounter[] }>(
      `/encounters?${params.toString()}`
    );
    return res.data?.data || [];
  },

  /**
   * Start a clinical encounter (in-person or telehealth channel)
   */
  async startEncounter(payload: StartEncounterPayload): Promise<ClinicalEncounter> {
    const res = await apiClient.post<{ success: boolean; data: ClinicalEncounter }>(
      "/encounters",
      payload
    );
    return res.data?.data || (res.data as any);
  },

  /**
   * Get active encounter by ID with enriched data
   */
  async getEncounterById(encounterId: string): Promise<ClinicalEncounter> {
    const res = await apiClient.get<{ success: boolean; data: ClinicalEncounter }>(
      `/encounters/${encounterId}`
    );
    return res.data?.data || (res.data as any);
  },

  /**
   * Save SOAP clinical notes and primary diagnosis
   */
  async saveSOAPNotes(encounterId: string, payload: UpdateSOAPNotesPayload): Promise<void> {
    await apiClient.put(`/encounters/${encounterId}/soap`, payload);
  },

  /**
   * Add ICD-10 diagnosis to encounter
   */
  async addDiagnosis(encounterId: string, payload: AddDiagnosisPayload): Promise<ClinicalDiagnosis> {
    const res = await apiClient.post<{ success: boolean; data: ClinicalDiagnosis }>(
      `/encounters/${encounterId}/diagnoses`,
      payload
    );
    return res.data?.data || (res.data as any);
  },

  /**
   * List diagnoses for encounter
   */
  async listDiagnoses(encounterId: string): Promise<ClinicalDiagnosis[]> {
    const res = await apiClient.get<{ success: boolean; data: ClinicalDiagnosis[] }>(
      `/encounters/${encounterId}/diagnoses`
    );
    return res.data?.data || [];
  },

  /**
   * Create electronic prescription with items
   */
  async createPrescription(
    encounterId: string,
    payload: CreatePrescriptionPayload
  ): Promise<Prescription> {
    const res = await apiClient.post<{ success: boolean; data: Prescription }>(
      `/encounters/${encounterId}/prescriptions`,
      payload
    );
    return res.data?.data || (res.data as any);
  },

  /**
   * List prescriptions for encounter
   */
  async listPrescriptions(encounterId: string): Promise<Prescription[]> {
    const res = await apiClient.get<{ success: boolean; data: Prescription[] }>(
      `/encounters/${encounterId}/prescriptions`
    );
    return res.data?.data || [];
  },

  /**
   * Dispatch diagnostic lab tests, radiology studies, and e-prescriptions
   */
  async dispatchOrders(encounterId: string, payload: DispatchOrdersPayload): Promise<void> {
    await apiClient.post(`/encounters/${encounterId}/orders`, payload);
  },

  /**
   * Finalize consultation, advance Care Journey, and generate cashier POS invoice
   */
  async completeEncounter(
    encounterId: string,
    consultationFee?: number
  ): Promise<CompleteEncounterResponse> {
    const res = await apiClient.post<{ success: boolean; data: CompleteEncounterResponse }>(
      `/encounters/${encounterId}/complete`,
      { consultationFee: consultationFee || 5000 }
    );
    return res.data?.data || (res.data as any);
  },
};
