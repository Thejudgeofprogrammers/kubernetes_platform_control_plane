import { api } from "./apiClient";
import type { APIService } from "../types/apiService";

interface RequestUpdate {
  name: string;
  base_url: string;
  protocol: "http" | "https" | "grpc";
}

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

export const getById = (id: string) => {
  api.get(`/api-services/${id}`).then((res) => res.data);
};

export const update = (id: string, data: RequestUpdate) => {
  api.put(`/api-services/${id}`, data);
};
