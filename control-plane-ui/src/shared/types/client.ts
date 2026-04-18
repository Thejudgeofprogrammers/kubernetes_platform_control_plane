export type ClientStatus =
  | "created"
  | "running"
  | "restarting"
  | "stopped"
  | "deleting"
  | "deploying";

export interface APIClient {
  id: string;
  name: string;
  description?: string;
  status: ClientStatus;
  api_service_id: string;
  activeConfigId?: string;
}
