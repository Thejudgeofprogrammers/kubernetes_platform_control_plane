import { useNavigate, useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import { api } from "../shared/api/apiClient";
import Card from "../shared/ui/Card";
import Input from "../shared/ui/Input";
import Button from "../shared/ui/Button";
import Select from "../shared/ui/Select";
import Layout from "../shared/ui/Layout";
import type { APIService } from "../shared/types/apiService";

export default function APIServiceDetailPage() {
  const { id } = useParams();

  const [service, setService] = useState<APIService | null>(null);
  const [name, setName] = useState("");
  const [baseURL, setBaseURL] = useState("");
  const [protocol, setProtocol] = useState("");

  useEffect(() => {
    api.get(`/api-services/${id}`).then((res) => {
      setService(res.data);
      setName(res.data.name);
      setBaseURL(res.data.base_url);
      setProtocol(res.data.protocol);
    });
  }, [id]);

  const navigate = useNavigate();

  const handleUpdate = async () => {
    await api.put(`/api-services/${id}`, {
      name,
      base_url: baseURL,
      protocol,
    });

    navigate("/api-services");
  };

  if (!service) return <div>Loading...</div>;

  return (
    <Layout>
      <Card>
        <h2>API Service Detail</h2>

        <Input value={name} onChange={(e) => setName(e.target.value)} />
        <div style={{ marginTop: 12 }} />

        <Input value={baseURL} onChange={(e) => setBaseURL(e.target.value)} />
        <div style={{ marginTop: 12 }} />

        <Select
          value={protocol}
          onChange={setProtocol}
          options={[
            { value: "http", label: "http" },
            { value: "https", label: "https" },
          ]}
        />

        <div style={{ marginTop: 12 }}>
          <Button onClick={handleUpdate}>Save</Button>
        </div>
      </Card>
    </Layout>
  );
}
