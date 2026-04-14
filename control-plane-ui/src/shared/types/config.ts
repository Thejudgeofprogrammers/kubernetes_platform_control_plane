export type AuthType = "none" | "api_key" | "bearer";

// export interface APIClientConfig {
//   id: string;
//   clientId: string;
//   version: string;
//   authType: AuthType;
//   authRef: string;
//   timeoutMs: number;
//   retryCount: number;
//   retryBackoff: number;
//   headers: Record<string, string>;
//   createdAt: string;
//   createdBy: string;
// }

export type APIClientConfig = {
  id: string;
  version: string;
  authType: string;
  authRef?: string;
  timeoutMs: number;
  retryCount: number;
  retryBackoff: number;
  headers: Record<string, string>;
};