import type { ReactNode } from "react";
import { Navigate } from "react-router-dom";
import { useAuthStore } from "../features/auth/store";

export function ProtectedRoute({ children }: { children: ReactNode }) {
  const token = useAuthStore((s) => s.accessToken);

  if (!token) {
    return <Navigate to="/login" />;
  }

  return children;
}
