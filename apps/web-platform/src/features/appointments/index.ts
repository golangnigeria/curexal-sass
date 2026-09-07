import React from "react";

export interface AppointmentSlot {
  id: string;
  providerId: string;
  providerName: string;
  serviceType: string;
  startTime: string;
  endTime: string;
  isBooked: boolean;
}
