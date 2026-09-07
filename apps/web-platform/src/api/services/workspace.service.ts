import { apiGet, apiPost } from "@curexal/api-client";
import type {
  LabSamplePayload,
  ConsultationQueuePayload,
  HospitalBedPayload,
  RadiologyScanPayload,
  PharmacyPrescriptionPayload,
} from "@curexal/contracts";

export type {
  LabSamplePayload,
  ConsultationQueuePayload,
  HospitalBedPayload,
  RadiologyScanPayload,
  PharmacyPrescriptionPayload,
};

class WorkspaceService {
  // LIMS Diagnostic Orders & Accessioning
  async getLabSamples(tenantId: string): Promise<LabSamplePayload[]> {
    try {
      const rawOrders = await apiGet<any[]>(`/lims/orders?organization_id=${tenantId}`);
      return (rawOrders || []).map((o: any) => ({
        id: o.id,
        orderNumber: o.orderNumber,
        sampleId: o.orderNumber,
        barcode: o.specimens?.[0]?.barcode || o.orderNumber,
        patientName: o.patientName || "Walk-in Patient",
        patientMrn: o.patientMrn || "PAT-REQ",
        testName: o.results?.[0]?.testName || "Diagnostic Panel",
        specimenType: o.specimens?.[0]?.specimenType || "Whole Blood (EDTA)",
        status: (o.status || "ordered").toLowerCase(),
        priority: (o.priority || "routine").toLowerCase(),
        receivedAt: new Date(o.createdAt || Date.now()).toLocaleTimeString(),
        results: (o.results || []).map((r: any) => ({
          parameter: r.parameterName,
          value: r.value,
          unit: r.unit,
          referenceRange: r.referenceRange,
          refRange: r.referenceRange,
          flag: r.flag,
        })),
      }));
    } catch {
      return [];
    }
  }

  async accessionSpecimen(tenantId: string, orderId: string, barcode: string, specimenType: string): Promise<boolean> {
    try {
      await apiPost(`/lims/specimens/accession?organization_id=${tenantId}`, { orderId, barcode, specimenType });
      return true;
    } catch {
      return false;
    }
  }

  async enterLabResults(tenantId: string, orderId: string, results: Array<{ testName: string; parameterName: string; value: string; unit: string; referenceRange: string; flag?: string }>): Promise<boolean> {
    try {
      await apiPost(`/lims/results?organization_id=${tenantId}`, { orderId, results });
      return true;
    } catch {
      return false;
    }
  }

  async authorizeLabOrder(tenantId: string, orderId: string, authorizerName: string, authorizerTitle: string, notes?: string): Promise<boolean> {
    try {
      await apiPost(`/lims/authorizations?organization_id=${tenantId}`, { orderId, authorizerName, authorizerTitle, notes });
      return true;
    } catch {
      return false;
    }
  }

  // EMR Clinical Queue
  async getConsultationQueue(tenantId: string): Promise<ConsultationQueuePayload[]> {
    try {
      const data = await apiGet<ConsultationQueuePayload[]>(`/encounters?organization_id=${tenantId}`);
      return data || [];
    } catch {
      return [];
    }
  }

  // HIS Hospital Beds
  async getHospitalBeds(tenantId: string): Promise<HospitalBedPayload[]> {
    try {
      const data = await apiGet<HospitalBedPayload[]>(`/workspace/${tenantId}/hospital/beds`);
      return data || [];
    } catch {
      return [];
    }
  }

  // RIS Radiology & PACS Modality Worklist
  async getRadiologyScans(tenantId: string): Promise<RadiologyScanPayload[]> {
    try {
      const rawOrders = await apiGet<any[]>(`/radiology/orders?organization_id=${tenantId}`);
      return (rawOrders || []).map((o: any) => ({
        id: o.id,
        accessionNumber: o.accessionNumber,
        patientName: o.patientName || "Patient",
        patientMrn: o.patientMrn || "PAT-RAD",
        modality: o.modality,
        procedureName: o.procedureName,
        status: (o.status || "scheduled").toLowerCase(),
        scheduledAt: new Date(o.createdAt || Date.now()).toLocaleTimeString(),
        radiologistName: o.report?.radiologistName,
        pacsUrl: o.studies?.[0]?.pacsSeriesUrl,
      }));
    } catch {
      return [];
    }
  }

  async recordStudyAcquisition(tenantId: string, orderId: string, studyInstanceUid: string, numberOfInstances: number, pacsSeriesUrl?: string): Promise<boolean> {
    try {
      await apiPost(`/radiology/studies?organization_id=${tenantId}`, { orderId, studyInstanceUid, numberOfInstances, pacsSeriesUrl });
      return true;
    } catch {
      return false;
    }
  }

  async signRadiologyReport(tenantId: string, orderId: string, radiologistName: string, findings: string, conclusion: string): Promise<boolean> {
    try {
      await apiPost(`/radiology/reports?organization_id=${tenantId}`, { orderId, radiologistName, findings, conclusion });
      return true;
    } catch {
      return false;
    }
  }

  // Pharmacy Dispensary Queue & Stock
  async getPharmacyPrescriptions(tenantId: string): Promise<PharmacyPrescriptionPayload[]> {
    try {
      const rawRxs = await apiGet<any[]>(`/pharmacy/prescriptions?organization_id=${tenantId}`);
      return (rawRxs || []).map((rx: any) => ({
        id: rx.id,
        rxNumber: rx.rxNumber,
        prescriptionNumber: rx.rxNumber,
        patientName: rx.patientName || "Patient",
        patientMrn: rx.patientMrn || "PAT-RX",
        prescribingDoctor: rx.prescribingDoctorName || "Consultant",
        status: (rx.status || "pending_verification").toLowerCase(),
        items: (rx.items || []).map((it: any) => ({
          medicationName: it.medicationName || it.name,
          drugName: it.medicationName || it.name,
          dosage: it.dosage,
          quantity: it.quantity,
          batchNumber: "BATCH-2026-A",
        })),
        createdAt: new Date(rx.createdAt || Date.now()).toLocaleTimeString(),
      }));
    } catch {
      return [];
    }
  }

  async getPharmacyMedications(tenantId: string): Promise<any[]> {
    try {
      const data = await apiGet<any[]>(`/pharmacy/medications?organization_id=${tenantId}`);
      return data || [];
    } catch {
      return [];
    }
  }

  async dispensePrescription(tenantId: string, prescriptionId: string, pharmacistName: string, itemsDispensed: any[]): Promise<boolean> {
    try {
      await apiPost(`/pharmacy/dispense?organization_id=${tenantId}`, { prescriptionId, pharmacistName, itemsDispensed });
      return true;
    } catch {
      return false;
    }
  }
}

export const workspaceService = new WorkspaceService();
