import { api } from "./apiClient";

export interface CreateClientPayload {
  name: string;
  description?: string;
  api_service_id: string;
}

export const createClient = (data: CreateClientPayload) => {
  return api.post("/clients", data);
};

export const getClients = async () => {
  const res = await api.get("/clients");
  return res.data.items;
};
