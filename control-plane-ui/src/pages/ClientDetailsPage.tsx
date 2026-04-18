import { useParams } from "react-router-dom";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "../shared/api/apiClient";
import type { APIClientConfig } from "../shared/types/config";
import type { APIClient } from "../shared/types/client";
import type { APIClientHealth } from "../shared/types/health";
import { useState } from "react";
import Card from "../shared/ui/Card";
import Button from "../shared/ui/Button";
import Input from "../shared/ui/Input";
import { useAPIServices } from "../features/api-service/hooks/useAPIServices";
import Layout from "../shared/ui/Layout";
import Select from "../shared/ui/Select";
import { mapConfig, mapHealth } from "../shared/api/mappers";

export default function ClientDetailsPage() {
  const { id } = useParams<{ id: string }>();
  
  const queryClient = useQueryClient();

  const [version, setVersion] = useState("");
  const [authType, setAuthType] = useState("none");
  const [authRef, setAuthRef] = useState("");
  const [timeoutMs, setTimeoutMs] = useState(1000);
  const [retryCount, setRetryCount] = useState(3);
  const [retryBackoff, setRetryBackoff] = useState(100);

  const [headers, setHeaders] = useState<{ key: string; value: string }[]>([
    { key: "", value: "" },
  ]);

  const handleHeaderChange = (index: number, field: "key" | "value", val: string) => {
    const updated = [...headers];
    updated[index][field] = val;
    setHeaders(updated);
  };

  const addHeader = () => {
    setHeaders([...headers, { key: "", value: "" }]);
  };

  const removeHeader = (index: number) => {
    if (headers.length === 1) return;
    setHeaders(headers.filter((_, i) => i !== index));
  };

  const buildHeadersMap = () => {
    const map: Record<string, string> = {};

    headers.forEach((h) => {
      if (h.key.trim()) {
        map[h.key.trim()] = h.value.trim();
      }
    });

    return map;
  };

  const { data: client } = useQuery<APIClient>({
    queryKey: ["client", id],
    queryFn: async () => {
      const res = await api.get(`/clients/${id}`);
      return res.data;
    },
    enabled: !!id,
  });

  const { data: configs } = useQuery<APIClientConfig[]>({
    queryKey: ["configs", id],
    queryFn: async () => {
      if (!id) return [];
      const res = await api.get(`/clients/${id}/configs`);
      return res.data.map(mapConfig);
    },
    enabled: !!id,
  });

  const { data: health } = useQuery<APIClientHealth>({
    queryKey: ["health", id],
    queryFn: async () => {
      try {
        const res = await api.get(`/clients/${id}/health`);
        console.log("HEALTH RAW:", res.data);
        return mapHealth(res.data);
      } catch {
        return {
          clientId: id!,
          status: "unknown",
          message: "failed to fetch",
          lastCheck: new Date().toISOString(),
        };
      }
    },
    refetchInterval: 3000,
    enabled: !!id,
  });

  const deploy = async (configId: string) => {
    await api.post(`/clients/${id}/configs/${configId}/deploy`);

    queryClient.invalidateQueries({ queryKey: ["configs", id] });
    queryClient.invalidateQueries({ queryKey: ["client", id] });
  };
  
  const { data: services = [] } = useAPIServices();

  if (!client) return <div>Loading...</div>;
  if (!id) return <div>Invalid client ID</div>;

  const service = services.find(
    (s) => s.id === client.api_service_id
  );


  
  const createConfig = async () => {
    try {
      await api.post(`/clients/${id}/configs`, {
        version: version,
        auth_type: authType,
        auth_ref: authRef,
        timeout_ms: timeoutMs,
        retry_count: retryCount,
        retry_backoff: retryBackoff,
        headers: buildHeadersMap(),
      });

      setVersion("");
      setAuthType("none");
      setAuthRef("");
      setTimeoutMs(1000);
      setRetryCount(3);
      setRetryBackoff(100);
      setHeaders([{ key: "", value: "" }]);

      queryClient.invalidateQueries({ queryKey: ["configs", id] });
    } catch (err) {
      console.error(err);
      alert("Failed to create config");
    }
  };

  const getHealthColor = () => {
    switch (health?.status) {
      case "healthy":
        return "green";
      case "degraded":
        return "orange";
      case "unhealthy":
        return "red";
      case "unknown":
        return "gray";
      default:
        return "gray";
    }
  };

    const status = client.status;

    const isBusy =
      status === "deploying" || status === "restarting";

    const canStart =
      status === "created" || status === "stopped";

    const canRestart =
      status === "running";

    const canDeploy =
      status !== "deleting";

    const canStop = status === "running";
    const canDelete = status !== "deleting";

    const stopClient = async () => {
      await api.post(`/clients/${id}/stop`);
      queryClient.invalidateQueries({ queryKey: ["client", id] });
    };

    const deleteClient = async () => {
      if (!confirm("Delete client?")) return;

      await api.post(`/clients/${id}/delete`);
      queryClient.invalidateQueries({ queryKey: ["clients"] });
    };

    const deleteConfig = async (configId: string) => {
      if (!confirm("Delete config?")) return;

      await api.delete(`/clients/${id}/configs/${configId}/delete`);

      queryClient.invalidateQueries({ queryKey: ["configs", id] });
    };

    const updateConfig = async (configId: string) => {
      const newVersion = prompt("New version?");
      if (!newVersion) return;

      await api.put(`/clients/${id}/configs/${configId}/update`, {
        version: newVersion,
        auth_type: "none",
        auth_ref: "",
        timeout_ms: 1000,
        retry_count: 3,
        retry_backoff: 100,
        headers: {},
      });

      queryClient.invalidateQueries({ queryKey: ["configs", id] });
    };

  return (
    <Layout>
    <div
      style={{
        padding: "24px",
        maxWidth: "900px",
        margin: "0 auto",
        fontFamily: "Arial",
      }}
    >
<Card>
  <div style={{ marginBottom: "12px" }}>
    <h1 style={{ margin: 0 }}>{client.name}</h1>

    <span className={`client-status status-${client.status}`}>
      status: {client.status}
    </span>
  </div>

  <div className="api-service-block">
    <div className="label">API Service</div>

      <div style={{ marginTop: "10px"}}></div>

    {service ? (
      <div className="service-card" style={{ marginLeft: "10px"}}>
        <div className="service-name">Service: {service.name}</div>
        <div className="service-url">BaseURL: {service.base_url}</div>
      </div>
    ) : (
      <div className="service-error">
        API Service not found
      </div>
    )}
  </div>
</Card>
      
      <div style={{ marginTop: "12px"}}></div>

    <Card>
      <h2>Create Config</h2>

      <div style={{ display: "flex", flexDirection: "column", gap: "10px" }}>
        
        <Input
          placeholder="Version (v1)"
          value={version}
          onChange={(e) => setVersion(e.target.value)}
        />

        <Select
          value={authType}
          onChange={setAuthType}
          options={[
            { value: "none", label: "No Auth" },
            { value: "api_key", label: "API Key" },
            { value: "bearer", label: "Bearer" },
          ]}
        />

        {authType !== "none" && (
          <Input
            placeholder={
              authType === "api_key"
                ? "API Key"
                : "Bearer token"
            }
            value={authRef}
            onChange={(e) => setAuthRef(e.target.value)}
          />
        )}

        <Input
          type="number"
          placeholder="Timeout (ms)"
          value={timeoutMs.toString()}
          onChange={(e) => setTimeoutMs(Number(e.target.value))}
        />

        <Input
          type="number"
          placeholder="Retry count"
          value={retryCount.toString()}
          onChange={(e) => setRetryCount(Number(e.target.value))}
        />

        <Input
          type="number"
          placeholder="Retry backoff (ms)"
          value={retryBackoff.toString()}
          onChange={(e) => setRetryBackoff(Number(e.target.value))}
        />

        <div style={{ marginTop: "10px" }}>
          <p style={{ fontWeight: "bold" }}>Headers</p>

          {headers.map((h, index) => (
            <div
              key={index}
              style={{
                display: "flex",
                gap: "8px",
                marginBottom: "8px",
              }}
            >
              <Input
                placeholder="Key"
                value={h.key}
                onChange={(e) =>
                  handleHeaderChange(index, "key", e.target.value)
                }
              />

              <Input
                placeholder="Value"
                value={h.value}
                onChange={(e) =>
                  handleHeaderChange(index, "value", e.target.value)
                }
              />

              <Button
                variant="danger"
                onClick={() => removeHeader(index)}
              >
                X
              </Button>
            </div>
          ))}

          <Button variant="secondary" onClick={addHeader}>
            + Add header
          </Button>
        </div>

        <Button onClick={createConfig} disabled={!version}>
          Create Config
        </Button>
      </div>
    </Card>

    <div style={{ marginTop: "12px"}}></div>

    {configs?.map((cfg) => (
      <Card key={cfg.id}>
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "flex-start",
          }}
        >
          <div>
            <b>{cfg.version}</b>{" "}
            {client.activeConfigId === cfg.id && (
              <span style={{ color: "green" }}>(active)</span>
            )}

            <p>Auth: {cfg.authType}</p>

            {cfg.authRef && (
              <p style={{ fontSize: "12px", color: "#666" }}>
                Ref: {cfg.authRef}
              </p>
            )}

            <p>Timeout: {cfg.timeoutMs} ms</p>
            <p>Retry: {cfg.retryCount}</p>
            <p>Backoff: {cfg.retryBackoff} ms</p>
            {cfg.headers && Object.keys(cfg.headers).length > 0 && (
              <div style={{ marginTop: "8px" }}>
                <p style={{ fontWeight: "bold", margin: "4px 0" }}>
                  Headers:
                </p>

                {Object.entries(cfg.headers).map(([key, value]) => (
                  <p
                    key={key}
                    style={{
                      fontSize: "12px",
                      color: "#555",
                      margin: 0,
                    }}
                  >
                    {key}: {value}
                  </p>
                ))}
              </div>
            )}
          </div>

          <div style={{ display: "flex", gap: "8px", marginTop: "12px" }}>
            <Button
              disabled={!canDeploy || isBusy || client.activeConfigId === cfg.id}
              onClick={() => deploy(cfg.id)}
            >
              🚀 Deploy
            </Button>

            <Button
              variant="secondary"
              onClick={() => updateConfig(cfg.id)}
            >
              ✏️ Update
            </Button>

            <Button
              variant="danger"
              disabled={client.activeConfigId === cfg.id}
              onClick={() => deleteConfig(cfg.id)}
            >
              🗑 Delete
            </Button>

            <Button
              variant="secondary"
              disabled={!canStart || isBusy}
              onClick={async () => {
                await api.post(`/clients/${id}/start`);
                queryClient.invalidateQueries({ queryKey: ["client", id] });
              }}
            >
              ▶️ Start
            </Button>

            <Button
              variant="secondary"
              disabled={!canRestart || isBusy}
              onClick={async () => {
                await api.post(`/clients/${id}/restart`);
                queryClient.invalidateQueries({ queryKey: ["client", id] });
              }}
            >
              🔄 Restart
            </Button>
            <Button
              variant="secondary"
              disabled={!canStop || isBusy}
              onClick={stopClient}
            >
              ⏹ Stop
            </Button>

            <Button
              variant="danger"
              disabled={!canDelete}
              onClick={deleteClient}
            >
              🗑 Delete
            </Button>
          </div>
        </div>
      </Card>
    ))}

      <Card>
        <h2>Health</h2>

        <p style={{ color: getHealthColor() }}>
          <b>{health?.status || "loading..."}</b>
        </p>

        <p>{health?.message || ""}</p>
      </Card>
    </div>
    </Layout>
  );
}