export type HealthStatus = "healthy" | "degraded" | "unhealthy" | "unknown";

export interface APIClientHealth {
  clientId: string;
  status: HealthStatus;
  lastCheck: string;
  message: string;
}
