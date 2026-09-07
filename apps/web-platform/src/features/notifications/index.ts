export interface NotificationConfig {
  smtp?: {
    host: string;
    port: number;
    user: string;
    fromEmail: string;
    fromName: string;
  };
  sms?: {
    provider: string;
    senderId: string;
  };
}
