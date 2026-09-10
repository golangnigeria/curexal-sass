import { apiClient } from "../client";
import type {
  Appointment,
  AppointmentFilter,
  CreateAppointmentPayload,
  ProviderProfile,
} from "../contracts";

export const operationsService = {
  /**
   * Schedule a new appointment
   */
  async createAppointment(payload: CreateAppointmentPayload): Promise<Appointment> {
    const res = await apiClient.post<{ success: boolean; data: Appointment }>(
      "/appointments",
      payload
    );
    return res.data?.data || res.data;
  },

  /**
   * List appointments with optional filters
   */
  async listAppointments(filter?: AppointmentFilter): Promise<Appointment[]> {
    const res = await apiClient.get<{ data: Appointment[] }>("/appointments", {
      params: filter,
    });
    if (Array.isArray(res.data?.data)) return res.data.data;
    if (Array.isArray(res.data)) return res.data;
    return [];
  },

  /**
   * Retrieve single appointment
   */
  async getAppointmentById(id: string): Promise<Appointment> {
    const res = await apiClient.get<{ data: Appointment }>(`/appointments/${id}`);
    return res.data?.data || res.data;
  },

  /**
   * Update appointment status (e.g. CANCELLED)
   */
  async updateAppointmentStatus(id: string, status: string, cancellationReason?: string): Promise<void> {
    await apiClient.put(`/appointments/${id}/status`, {
      status,
      cancellationReason,
    });
  },

  /**
   * List provider profiles
   */
  async listProviderProfiles(status?: string, specialty?: string): Promise<ProviderProfile[]> {
    const res = await apiClient.get<{ data: ProviderProfile[] }>("/providers/profiles", {
      params: { status, specialty },
    });
    if (Array.isArray(res.data?.data)) return res.data.data;
    if (Array.isArray(res.data)) return res.data;
    return [];
  },

  /**
   * Update provider on-duty status
   */
  async updateProviderStatus(id: string, status: string, reason?: string): Promise<void> {
    await apiClient.put(`/providers/profiles/${id}/status`, {
      status,
      reason,
    });
  },
};
