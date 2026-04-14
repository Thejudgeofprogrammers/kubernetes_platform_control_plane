import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { getClients } from "../shared/api/clientApi";
import { useCreateClient } from "../features/client/hooks/useCreateClient";
import Button from "../shared/ui/Button";
import Card from "../shared/ui/Card";
import ClientCard from "../widgets/ClientCard";
import { SkeletonClients } from "../shared/ui/SceletonClients";
import Input from "../shared/ui/Input";
import type { APIClient } from "../shared/types/client";
import { useAPIServices } from "../features/api-service/hooks/useAPIServices";
import Select from "../shared/ui/Select";
import Layout from "../shared/ui/Layout";
import { Link } from "react-router-dom";


export default function ClientsPage() {
  const { data, isLoading } = useQuery<APIClient[]>({
    queryKey: ["clients"],
    queryFn: getClients,
  });

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [apiServiceId, setApiServiceId] = useState("");
  const { data: services = [] } = useAPIServices();

  const createMutation = useCreateClient();

  const handleCreate = async () => {
    try {
      await createMutation.mutateAsync({
        name,
        description,
        api_service_id: apiServiceId,
      });
    } catch (err) {
      alert(err);
    }
  };

  if (isLoading) return <SkeletonClients />;

  return (
    <Layout>
    <>
      <Card>
        <h3>Create client</h3>


        <Input
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Name"
        />
        <div style={{ marginTop: "12px"}}></div>
        <Input
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="Description"
        />
        <div style={{ marginTop: "12px"}}></div>
        <Select
          value={apiServiceId}
          onChange={setApiServiceId}
          options={services.map((s) => ({
            value: s.id,
            label: s.name,
          }))}
        />

        <div style={{ marginTop: "12px"}}>
          <Button
            onClick={handleCreate}
            disabled={!name || !apiServiceId || createMutation.isPending}
          >
            {createMutation.isPending ? "Creating..." : "Create"}
          </Button>

          <span style={{ marginTop: "12px", marginLeft: "8px" }}>
            <Link to="/api-services" style={{ textDecoration: "none" }}>
              <Button variant="secondary">
                API Services
              </Button>
            </Link>
          </span>
        </div>
      </Card>

      {data?.map((client) => (
        <ClientCard key={client.id} client={client} />
      ))}
    </>
    </Layout>
  );
}