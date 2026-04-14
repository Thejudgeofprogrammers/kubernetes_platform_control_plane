import { create } from "zustand";

type User = {
  id: string;
  email: string;
  full_name: string;
  role: "owner" | "operator" | "viewer";
};

interface AuthState {
  accessToken: string | null;
  user: User | null;
  setToken: (token: string) => void;
  setUser: (user: User) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  accessToken: localStorage.getItem("access_token"),
  user: null,

  setToken: (token) => {
    localStorage.setItem("access_token", token);
    set({ accessToken: token });
  },

  setUser: (user) => {
    set({ user });
  },

  logout: () => {
    localStorage.removeItem("access_token");
    set({ accessToken: null, user: null });
  },
}));


export const useRole = () => {
  const user = useAuthStore((s) => s.user);

  return {
    isOwner: user?.role === "owner",
    isOperator: user?.role === "operator",
    isViewer: user?.role === "viewer",
  };
};
