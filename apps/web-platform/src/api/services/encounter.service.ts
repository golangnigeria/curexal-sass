import { apiClient } from "../client";

export interface Encounter {
  id: string;
  tenantId: string;
  careRequestId?: string;
  patientId: string;
  providerId: string;
  encounterType: "OUTPATIENT" | "EMERGENCY" | "TELEHEALTH" | "INPATIENT_ROUND" | string;
  mode: "IN_PERSON" | "VIDEO" | "AUDIO" | "ASYNC_CHAT" | string;
  status: "WAITING" | "IN_PROGRESS" | "ON_HOLD" | "COMPLETED" | "CANCELLED" | string;
  chiefComplaint?: string;
  subjective?: string;
  objective?: string;
  assessment?: string;
  plan?: string;
  primaryDiagnosisCode?: string;
  primaryDiagnosisName?: string;
  secondaryDiagnoses?: string[];
  startedAt: string;
  completedAt?: string;
  createdAt: string;
  updatedAt: string;

  // Enriched patient info
  patientName?: string;
  mrn?: string;
  gender?: string;
  ageYears?: number;
}

export interface StartEncounterPayload {
  careRequestId?: string;
  patientId: string;
  providerId: string;
  encounterType: "OUTPATIENT" | "EMERGENCY" | "TELEHEALTH" | "INPATIENT_ROUND" | string;
  mode: "IN_PERSON" | "VIDEO" | "AUDIO" | "ASYNC_CHAT" | string;
  chiefComplaint?: string;
}

export interface UpdateSOAPNotesPayload {
  chiefComplaint?: string;
  subjective?: string;
  objective?: string;
  assessment?: string;
  plan?: string;
  primaryDiagnosisCode?: string;
  primaryDiagnosisName?: string;
  secondaryDiagnoses?: string[];
}

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

export const encounterService = {
  /**
   * Start a clinical encounter
   */
  async startEncounter(payload: StartEncounterPayload): Promise<Encounter> {
    const res = await apiClient.post<{ success: boolean; data: Encounter }>(
      "/encounters/start",
      payload
    );
    return res.data?.data || res.data;
  },

  /**
   * Get active encounter by ID
   */
  async getEncounterById(encounterId: string): Promise<Encounter> {
    const res = await apiClient.get<{ success: boolean; data: Encounter }>(
      `/encounters/${encounterId}`
    );
    return res.data?.data || res.data;
  },

  /**
   * Save SOAP clinical notes and primary diagnosis
   */
  async saveSOAPNotes(encounterId: string, payload: UpdateSOAPNotesPayload): Promise<void> {
    await apiClient.put(`/encounters/${encounterId}/notes`, payload);
  },

  /**
   * Dispatch diagnostic lab tests, radiology studies, and e-prescriptions
   */
  async dispatchOrders(encounterId: string, payload: DispatchOrdersPayload): Promise<void> {
    await apiClient.post(`/encounters/${encounterId}/orders`, payload);
  },

  /**
   * Finalize consultation and advance Care Journey to settlement
   */
  async completeEncounter(encounterId: string): Promise<void> {
    await apiClient.post(`/encounters/${encounterId}/complete`);
  },
};
