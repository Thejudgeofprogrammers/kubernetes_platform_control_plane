import { useState } from "react";
import { useAPIServices } from "../features/api-service/hooks/useAPIServices";
import { useCreateAPIService } from "../features/api-service/hooks/useCreateAPIService";
import { useDeleteAPIService } from "../features/api-service/hooks/useDeleteAPIService";
import Card from "../shared/ui/Card";
import Input from "../shared/ui/Input";
import Button from "../shared/ui/Button";
import Layout from "../shared/ui/Layout";
import Select from "../shared/ui/Select";
import { Link } from "react-router-dom";

export default function APIServicesPage() {
  const { data: services = [], isLoading } = useAPIServices();
  const createMutation = useCreateAPIService();
  const deleteMutation = useDeleteAPIService();

  const [name, setName] = useState("");
  const [baseURL, setBaseURL] = useState("");
  const [protocol, setProtocol] = useState("");

  const handleCreate = async () => {
    createMutation.mutate({
        name,
        base_url: baseURL,
        protocol,
    });
    setName("");
    setBaseURL("");
  };

  if (isLoading) return <div>Loading...</div>;

  return (
    <Layout>
    <>
      <Card>
        <h3>Create API Service</h3>

        <Input
          placeholder="Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <div style={{ marginTop: "12px"}}></div>
        <Input
          placeholder="Base URL"
          value={baseURL}
          onChange={(e) => setBaseURL(e.target.value)}
        />
        <div style={{ marginTop: "12px"}}></div>
        <Select
            value={protocol}
            onChange={setProtocol}
            options={[
                { value: "http", label: "http" },
                { value: "https", label: "https" },
            ]}
        />

        <div style={{ marginTop: "12px"}}>
            <Button onClick={handleCreate} disabled={!name || !baseURL}>
                Create
            </Button>
        </div>
      </Card>

      {services.map((s) => (
        <Card key={s.id}>
          <h3>
            <Link to={`/api-services/${s.id}`}>{s.name}</Link>
          </h3>
          <p>BaseURL: {s.base_url}</p>
          <p>Protocol: {s.protocol}</p>

          <Button
            variant="danger"
            onClick={() => deleteMutation.mutate(s.id)}
          >
            Delete
          </Button>
        </Card>
      ))}
    </>
    </Layout>
  );
}
