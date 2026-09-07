import { useEffect, useRef, useState, useCallback } from "react";

export interface RealtimeMessage<T = any> {
  topic: string;
  event: string;
  payload: T;
  timestamp: string;
}

export interface RealtimeClientConfig {
  wsUrl?: string;
  enableDebug?: boolean;
}

export type RealtimeEventHandler<T = any> = (payload: T) => void;

class RealtimeBus {
  private listeners: Map<string, Set<RealtimeEventHandler>> = new Map();
  private ws: WebSocket | null = null;
  private config: RealtimeClientConfig = {};
  private reconnectTimeout: any = null;

  init(config: RealtimeClientConfig) {
    this.config = config;
    this.connect();
  }

  private connect() {
    if (typeof window === "undefined") return;

    // 1. Determine connection URL (WebSocket broker)
    const wsUrl = this.config.wsUrl;
    if (!wsUrl) return;

    try {
      this.ws = new WebSocket(wsUrl);

      this.ws.onopen = () => {
        if (this.config.enableDebug) {
          console.log("[Realtime] Connected to realtime stream:", wsUrl);
        }
      };

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          const channel = data.topic || data.channel || "";
          const eventType = data.event || "";
          const payload = data.payload || data;

          const key = `${channel}:${eventType}`;
          const handlers = this.listeners.get(key);
          if (handlers) {
            handlers.forEach((h) => h(payload));
          }

          // Channel-wide handlers
          const channelHandlers = this.listeners.get(channel);
          if (channelHandlers) {
            channelHandlers.forEach((h) => h({ event: eventType, payload }));
          }
        } catch {
          // Non-JSON message
        }
      };

      this.ws.onclose = () => {
        if (this.config.enableDebug) {
          console.log("[Realtime] Stream disconnected. Reconnecting in 5s...");
        }
        clearTimeout(this.reconnectTimeout);
        this.reconnectTimeout = setTimeout(() => this.connect(), 5000);
      };

      this.ws.onerror = (err) => {
        if (this.config.enableDebug) {
          console.warn("[Realtime] Stream error:", err);
        }
      };
    } catch (err) {
      if (this.config.enableDebug) {
        console.warn("[Realtime] Failed to initialize WebSocket:", err);
      }
    }
  }

  subscribe<T = any>(channel: string, event: string, handler: RealtimeEventHandler<T>): () => void {
    const key = event ? `${channel}:${event}` : channel;
    if (!this.listeners.has(key)) {
      this.listeners.set(key, new Set());
    }
    this.listeners.get(key)!.add(handler);

    return () => {
      const handlers = this.listeners.get(key);
      if (handlers) {
        handlers.delete(handler);
        if (handlers.size === 0) {
          this.listeners.delete(key);
        }
      }
    };
  }

  publishLocal<T = any>(channel: string, event: string, payload: T) {
    const key = `${channel}:${event}`;
    const handlers = this.listeners.get(key);
    if (handlers) {
      handlers.forEach((h) => h(payload));
    }
  }
}

export const realtimeBus = new RealtimeBus();

/**
 * Vendor-agnostic React hook for subscribing to a specific realtime channel and event.
 */
export function useRealtimeEvent<T = any>(
  channel: string | null | undefined,
  event: string,
  handler: RealtimeEventHandler<T>
) {
  const handlerRef = useRef(handler);
  handlerRef.current = handler;

  useEffect(() => {
    if (!channel) return;
    const unsub = realtimeBus.subscribe<T>(channel, event, (payload) => {
      handlerRef.current(payload);
    });
    return unsub;
  }, [channel, event]);
}

/**
 * Realtime hook for Outpatient Queue updates in a clinic branch.
 */
export function useLiveQueue(
  organizationId: string | undefined,
  branchId: string | undefined,
  onQueueUpdate: (eventData: any) => void
) {
  const channel = organizationId ? `clinic:${organizationId}:${branchId || "all"}:queue` : null;
  useRealtimeEvent(channel, "queue.updated", onQueueUpdate);
  useRealtimeEvent(channel, "appointment.checked_in", onQueueUpdate);
  useRealtimeEvent(channel, "appointment.triage_completed", onQueueUpdate);
  useRealtimeEvent(channel, "appointment.consultation_started", onQueueUpdate);
}

/**
 * Realtime hook for Doctor Consultation canvas vitals & triage alerts.
 */
export function useEncounterVitals(
  encounterId: string | undefined,
  onVitalsRecorded: (vitals: any) => void
) {
  const channel = encounterId ? `encounter:${encounterId}` : null;
  useRealtimeEvent(channel, "vitals.recorded", onVitalsRecorded);
}
