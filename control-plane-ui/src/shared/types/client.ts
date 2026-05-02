export type ClientStatus =
  | "created"
  | "running"
  | "restarting"
  | "stopped"
  | "stopping"
  | "deleting"
  | "deploying"
  | "disabled";

export interface APIClient {
  id: string;
  name: string;
  slug: string
  url: string
  description?: string;
  status: ClientStatus;
  api_service_id: string;
  activeConfigId?: string;
}
