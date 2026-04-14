import type { APIClientConfig } from "../types/config";
import type { APIClientHealth, HealthStatus } from "../types/health";

interface OldConfig {
    id: string;
    version: string;
    auth_type: string;
    auth_ref?: string;
    timeout_ms: number;
    retry_count: number;
    retry_backoff: number;
    headers: Record<string, string>
}

interface OldHealth {
  ClientID: string;
  Status: HealthStatus;
  Message: string;
  LastCheck: string;
}


export const mapConfig = (c: OldConfig): APIClientConfig => ({
    id: c.id,
    version: c.version,
    authType: c.auth_type,
    authRef: c.auth_ref,
    timeoutMs: c.timeout_ms,
    retryCount: c.retry_count,
    retryBackoff: c.retry_backoff,
    headers: c.headers ?? {},
});

export const mapHealth = (h: OldHealth): APIClientHealth => ({
  clientId: h.ClientID,
  status: h.Status,
  message: h.Message,
  lastCheck: h.LastCheck,
});