/**
 * Telemetry and Clinical Audit Logger
 */

export interface TelemetryEvent {
  action: string;
  category: "CLINICAL" | "SECURITY" | "FINANCIAL" | "SYSTEM";
  entityId?: string;
  metadata?: Record<string, any>;
  timestamp?: string;
}

export function logTelemetry(event: TelemetryEvent): void {
  if (typeof window !== "undefined" && (window as any).__DEV__) {
    // console.debug(`[Telemetry] ${event.category}:${event.action}`, event.metadata);
  }
}
