import { api } from "./apiClient";
import type { APIService } from "../types/apiService";

export const getAPIServices = async (): Promise<APIService[]> => {
  const res = await api.get("/api-services");
  return res.data ?? [];
};

export const createAPIService = (data: {
  name: string;
  base_url: string;
  protocol: string;
}) => {
  return api.post("/api-services", data);
};

export const deleteAPIService = (id: string) => {
  return api.delete(`/api-services/${id}`);
};