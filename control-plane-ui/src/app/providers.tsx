import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import React, { useEffect } from "react";
import { useAuthStore } from "../features/auth/store";
import { api } from "../shared/api/apiClient";

const queryClient = new QueryClient();

export function Providers({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    const token = localStorage.getItem("access_token");

    if (token) {
      api
        .get("/users/me")
        .then((res) => {
          useAuthStore.getState().setUser(res.data);
        })
        .catch(() => {
          useAuthStore.getState().logout();
        });
    }
  }, []);

  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}
