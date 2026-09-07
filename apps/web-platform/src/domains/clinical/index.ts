export interface ClinicalVitals {
  bp: string;
  pulse: string;
  temp: string;
  weight?: string;
  spo2?: string;
  height?: string;
}

export interface SoapEncounterNote {
  id?: string;
  patientId: string;
  encounterId: string;
  subjective: string;
  objective: string;
  assessment: string;
  plan: string;
  diagnosesIcd10?: string[];
  orders?: {
    labTestIds?: string[];
    radiologyScanIds?: string[];
    medicationIds?: string[];
  };
  createdAt?: string;
  signedByDoctorId?: string;
}
