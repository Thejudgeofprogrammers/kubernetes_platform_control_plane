import type { APIClient } from "../shared/types/client";
import { api } from "../shared/api/apiClient";
import { useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import Button from "../shared/ui/Button";
import { useRole } from "../features/auth/store";

interface Props {
  client: APIClient;
}

export default function ClientCard({ client }: Props) {
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const { isOwner } = useRole();

  const restart = async () => {
    await api.post(`/clients/${client.id}/restart`);
    queryClient.invalidateQueries({ queryKey: ["clients"] });
  };
  const canDelete = client.status !== "deleting";
  const remove = async () => {
    if (!confirm(`Delete client "${client.name}"?`)) return;

    await api.post(`/clients/${client.id}/delete`);

    queryClient.invalidateQueries({ queryKey: ["clients"] });
  };

  const getStatusColor = () => {
    switch (client.status) {
      case "running":
        return "green";
      case "restarting":
        return "orange";
      case "stopped":
        return "gray";
      case "deleting":
        return "red";
      default:
        return "black";
    }
  };

  return (
    <div
      style={{
        border: "1px solid #ccc",
        padding: "12px",
        marginBottom: "10px",
        borderRadius: "8px",
      }}
    >
      <h3
        style={{ cursor: "pointer" }}
        onClick={() => navigate(`/clients/${client.id}`)}
      >
        {client.name}
      </h3>

      <p style={{ color: getStatusColor() }}>Status: {client.status}</p>

      <div style={{ display: "flex", gap: "8px" }}>
        <Button disabled={client.status !== "running"} onClick={restart}>
          Restart
        </Button>
        {isOwner && (
          <Button onClick={remove} variant="danger" disabled={!canDelete}>
            Delete
          </Button>
        )}
      </div>
    </div>
  );
}
