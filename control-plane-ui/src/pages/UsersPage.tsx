import { useEffect, useState } from "react";
import { usersApi } from "../shared/api/userAPI";
import { CreateUserForm } from "./CreateUserFrom";
import Button from "../shared/ui/Button";
import "../shared/styles/users.css";
import SelectRole from "../shared/ui/SelectForRole";
import Layout from "../shared/ui/Layout";

type User = {
  id: string;
  email: string;
  full_name: string;
  role: "viewer" | "operator" | "owner";
};

export default function UsersPage() {
  const [users, setUsers] = useState<User[]>([]);

    const load = async () => {
        try {
            const data = await usersApi.getUsers();
            setUsers(data);
        } catch (e) {
            console.error(e);
        }
    };

    useEffect(() => {
        let mounted = true;

        const init = async () => {
            try {
                const data = await usersApi.getUsers();
            if (mounted) setUsers(data);
            } catch (e) {
                console.error(e);
            }
        };

        init();

        return () => {
            mounted = false;
        };
    }, []);

    const handleDelete = async (id: string) => {
        await usersApi.deleteUser(id);
        setUsers(prev => prev.filter(u => u.id !== id));
    };

    const handleRoleChange = async (
        id: string,
        role: User["role"]
    ) => {
        await usersApi.updateRole(id, role);

        setUsers(prev =>
            prev.map(u => (u.id === id ? { ...u, role } : u))
        );
    };

  return (
    <Layout>
        <>
            <div className="users-page">
            <h1 className="title">Users</h1>

            <div className="card">
                <CreateUserForm onCreated={load} />
            </div>

            <div className="card">
                <table className="users-table">
                <thead>
                    <tr>
                    <th>Email</th>
                    <th>Name</th>
                    <th>Role</th>
                    <th></th>
                    </tr>
                </thead>

                <tbody>
                    {users.map(u => (
                    <tr key={u.id}>
                        <td>{u.email}</td>
                        <td>{u.full_name}</td>

                        <td>
                            {u.role === "owner" ? (
                                <span className="role-owner">Owner</span>
                            ) : (
                                <SelectRole
                                value={u.role}
                                onChange={(val) => handleRoleChange(u.id, val as User["role"])}
                                options={[
                                    { value: "viewer", label: "Viewer" },
                                    { value: "operator", label: "Operator" },
                                ]}
                                />
                            )}
                        </td>

                        <td>
                        <Button
                            onClick={() => handleDelete(u.id)}
                            variant="danger"
                        >
                            Delete
                        </Button>
                        </td>
                    </tr>
                    ))}
                </tbody>
                </table>
            </div>
            </div>
        </>
    </Layout>
  );
}
