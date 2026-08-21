import { authClient } from "@/lib/auth-client";

export interface LabSamplePayload {
  id: string;
  sampleId: string;
  patientName: string;
  patientMrn: string;
  testName: string;
  testCode: string;
  status: "collected" | "accessioned" | "in_analysis" | "results_pending" | "authorized";
  priority: "routine" | "stat" | "urgent";
  collectedAt: string;
  results?: Array<{ parameter: string; value: string; unit: string; referenceRange: string; flag?: "NORMAL" | "HIGH" | "LOW" | "PANIC" }>;
}

export interface ConsultationQueuePayload {
  id: string;
  queueNumber: number;
  patientName: string;
  patientMrn: string;
  triageCategory: "emergency" | "urgent" | "routine";
  vitals?: { bp: string; pulse: string; temp: string; weight: string; spo2: string };
  status: "waiting" | "in_consultation" | "completed";
  waitingSince: string;
}

export interface HospitalBedPayload {
  id: string;
  wardName: string;
  bedNumber: string;
  status: "occupied" | "available" | "cleaning" | "maintenance";
  patientName?: string;
  patientMrn?: string;
  admissionDate?: string;
  attendingPhysician?: string;
}

export interface RadiologyScanPayload {
  id: string;
  accessionNumber: string;
  patientName: string;
  patientMrn: string;
  modality: "X-RAY" | "ULTRASOUND" | "CT" | "MRI" | "MAMMOGRAPHY";
  procedureName: string;
  status: "scheduled" | "in_progress" | "images_acquired" | "reported";
  scheduledAt: string;
  radiologistName?: string;
}

export interface PharmacyPrescriptionPayload {
  id: string;
  prescriptionNumber: string;
  patientName: string;
  prescribingDoctor: string;
  status: "pending_dispense" | "dispensed" | "cancelled";
  items: Array<{ drugName: string; dosage: string; quantity: number; batchNumber?: string; expiryDate?: string }>;
  createdAt: string;
}

class WorkspaceService {
  private async getCsrfHeader(): Promise<Record<string, string>> {
    const csrfToken = authClient.getCsrfToken?.() || "";
    return csrfToken ? { "X-CSRF-Token": csrfToken } : {};
  }

  // LIMS Accessioning
  async getLabSamples(tenantId: string): Promise<LabSamplePayload[]> {
    const res = await fetch(`/api/v1/workspace/${tenantId}/laboratory/samples`, { credentials: "include" });
    if (!res.ok) return [];
    return res.json();
  }

  // EMR Clinical Queue
  async getConsultationQueue(tenantId: string): Promise<ConsultationQueuePayload[]> {
    const res = await fetch(`/api/v1/workspace/${tenantId}/clinical/queue`, { credentials: "include" });
    if (!res.ok) return [];
    return res.json();
  }

  // HIS Hospital Beds
  async getHospitalBeds(tenantId: string): Promise<HospitalBedPayload[]> {
    const res = await fetch(`/api/v1/workspace/${tenantId}/hospital/beds`, { credentials: "include" });
    if (!res.ok) return [];
    return res.json();
  }

  // RIS Modality Queue
  async getRadiologyScans(tenantId: string): Promise<RadiologyScanPayload[]> {
    const res = await fetch(`/api/v1/workspace/${tenantId}/radiology/scans`, { credentials: "include" });
    if (!res.ok) return [];
    return res.json();
  }

  // Pharmacy Dispensary Queue
  async getPharmacyPrescriptions(tenantId: string): Promise<PharmacyPrescriptionPayload[]> {
    const res = await fetch(`/api/v1/workspace/${tenantId}/pharmacy/prescriptions`, { credentials: "include" });
    if (!res.ok) return [];
    return res.json();
  }
}

export const workspaceService = new WorkspaceService();
