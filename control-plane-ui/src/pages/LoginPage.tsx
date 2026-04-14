import { useState } from "react";
import { api } from "../shared/api/apiClient";
import { useAuthStore } from "../features/auth/store";
import { useNavigate } from "react-router-dom";
import Input from "../shared/ui/Input";
import Button from "../shared/ui/Button";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [code, setCode] = useState("");
  const [step, setStep] = useState<"email" | "code">("email");

  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const requestCode = async () => {
    try {
      await api.post("/auth/request-code", { email });
      setStep("code");
    } catch (err) {
      alert(err);
    }
  };

  const handleLogin = async () => {
    setLoading(true);
    try {
      const res = await api.post("/auth/verify-code", { email, code });

      useAuthStore.getState().setToken(res.data.access_token);
      localStorage.setItem("refresh_token", res.data.refresh_token);
      
      const me = await api.get("/users/me");
      useAuthStore.getState().setUser(me.data);

      navigate("/");
  } catch (err) {
      alert(err);
  } finally {
      setLoading(false);
  }
  };

  return (
    <div
      style={{
        height: "100vh",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        background: "linear-gradient(to bottom, #f9fafb, #ffffff)",
        fontFamily: "Arial",
      }}
    >
      <div
        style={{
          width: "360px",
          padding: "32px",
          borderRadius: "16px",
          background: "white",
          boxShadow: "0 10px 30px rgba(0,0,0,0.08)",
        }}
      >
        <h2 style={{ marginBottom: "6px" }}>
          🔐 Sign in
        </h2>

        <p style={{ color: "#777", fontSize: "14px", marginBottom: "24px" }}>
          {step === "email"
            ? "Enter your email to receive a login code"
            : "Enter the code sent to your email"}
        </p>

        <div style={{ display: "flex", gap: "6px", marginBottom: "20px" }}>
          <div
            style={{
              flex: 1,
              height: "4px",
              borderRadius: "4px",
              background: step === "email" ? "#6366f1" : "#6366f1",
            }}
          />
          <div
            style={{
              flex: 1,
              height: "4px",
              borderRadius: "4px",
              background: step === "code" ? "#6366f1" : "#e5e7eb",
            }}
          />
        </div>

        {step === "email" && (
          <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
          <Input
            placeholder="you@example.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter" && email && !loading) {
                requestCode();
              }
            }}
          />

            <Button
              onClick={requestCode}
              disabled={!email || loading}
            >
              {loading ? "Sending..." : "Send Code"}
            </Button>
          </div>
        )}

        {step === "code" && (
          <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
          <Input
            placeholder="Enter code"
            value={code}
            onChange={(e) => setCode(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter" && code && !loading) {
                handleLogin();
              }
            }}
          />

            <Button
              onClick={handleLogin}
              disabled={!code || loading}
            >
              {loading ? "Verifying..." : "Verify & Login"}
            </Button>

            <button
              onClick={() => setStep("email")}
              style={{
                marginTop: "6px",
                background: "none",
                border: "none",
                color: "#6366f1",
                cursor: "pointer",
                fontSize: "13px",
              }}
            >
              ← Change email
            </button>
          </div>
        )}
      </div>
    </div>
  );
}