import { useState } from "react";
import { usersApi } from "../shared/api/userAPI";
import Input from "../shared/ui/Input";
import Button from "../shared/ui/Button";

export function CreateUserForm({ onCreated }: { onCreated: () => void }) {
  const [email, setEmail] = useState("");
  const [fullName, setFullName] = useState("");

  const handleCreate = async () => {
    await usersApi.createUser({
      email,
      full_name: fullName,
    });

    setEmail("");
    setFullName("");

    onCreated();
  };

  return (
    <div style={{ marginBottom: 20 }}>
      <h3>Create User</h3>

      <Input
        placeholder="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />
      <div style={{ margin: "12px" }}></div>
      <Input
        placeholder="full name"
        value={fullName}
        onChange={(e) => setFullName(e.target.value)}
      />
      <div style={{ margin: "12px" }}></div>
      <Button onClick={handleCreate}>Create</Button>
    </div>
  );
}
