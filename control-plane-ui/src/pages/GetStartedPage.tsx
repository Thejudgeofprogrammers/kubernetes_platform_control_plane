import { useNavigate } from "react-router-dom";

export default function GetStartedPage() {
  const navigate = useNavigate();

  return (
    <div
      style={{
        minHeight: "100vh",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        fontFamily: "Arial",
        padding: "40px 20px",
        background: "linear-gradient(to bottom, #f9fafb, #ffffff)",
      }}
    >
      <h1
        style={{ fontSize: "42px", marginBottom: "16px", textAlign: "center" }}
      >
        🚀 Control Plane
      </h1>

      <p
        style={{
          fontSize: "18px",
          color: "#555",
          maxWidth: "600px",
          textAlign: "center",
          marginBottom: "30px",
        }}
      >
        Управляй API-клиентами, конфигурациями и деплоем через Kubernetes —
        централизованно, безопасно и без лишней боли.
      </p>

      <button
        onClick={() => navigate("/clients")}
        style={{
          padding: "12px 28px",
          fontSize: "16px",
          borderRadius: "8px",
          border: "none",
          background: "#6366f1",
          color: "white",
          cursor: "pointer",
          marginBottom: "50px",
        }}
      >
        Get Started
      </button>

      <div
        style={{
          display: "grid",
          gridTemplateColumns: "repeat(auto-fit, minmax(240px, 1fr))",
          gap: "20px",
          maxWidth: "900px",
          width: "100%",
        }}
      >
        <FeatureCard
          title="⚙️ Управление конфигами"
          desc="Создавай версии конфигураций с retry, timeout и авторизацией."
        />

        <FeatureCard
          title="🚀 Деплой в Kubernetes"
          desc="Автоматическое развёртывание API-клиентов через Deployment."
        />

        <FeatureCard
          title="🔄 Управление состоянием"
          desc="Start, Restart, Deploy — всё через UI без ручных команд."
        />

        <FeatureCard
          title="❤️ Health Monitoring"
          desc="Отслеживай состояние клиентов в реальном времени."
        />
      </div>

      <div
        style={{ marginTop: "60px", textAlign: "center", maxWidth: "700px" }}
      >
        <h2 style={{ marginBottom: "16px" }}>Как это работает</h2>

        <p style={{ color: "#666", lineHeight: "1.6" }}>
          Ты создаёшь API-client → настраиваешь конфигурацию → Control Plane
          разворачивает его в Kubernetes → весь трафик маршрутизируется через
          Ingress.
        </p>
      </div>
    </div>
  );
}

function FeatureCard({ title, desc }: { title: string; desc: string }) {
  return (
    <div
      style={{
        padding: "20px",
        borderRadius: "12px",
        border: "1px solid #eee",
        background: "white",
        boxShadow: "0 4px 12px rgba(0,0,0,0.04)",
      }}
    >
      <h3 style={{ marginBottom: "8px" }}>{title}</h3>
      <p style={{ color: "#666", fontSize: "14px" }}>{desc}</p>
    </div>
  );
}
