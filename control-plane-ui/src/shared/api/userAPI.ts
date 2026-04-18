import { api } from "../../shared/api/apiClient";

export const usersApi = {
  getUsers: () =>
    api.get("/users").then(res => res.data.items),

  deleteUser: (id: string) =>
    api.delete(`/users/${id}`),

  updateRole: (id: string, role: string) =>
    api.patch(`/users/${id}/role`, { role }),

  createUser: (data: { email: string; full_name: string }) =>
    api.post("/auth/register", data),
};
