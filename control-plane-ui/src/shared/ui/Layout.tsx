import { Link, useNavigate } from "react-router-dom";
import { useAuthStore } from "../../features/auth/store";
import "./layout.css";

export default function Layout({ children }: { children: React.ReactNode }) {
  const navigate = useNavigate();
  const { user, logout } = useAuthStore();

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <div className="layout">
      <div className="navbar">
        <div className="nav-left">
          <Link to="/" className="nav-link">Main</Link>
          <Link to="/clients" className="nav-link">Clients</Link>
          <Link to="/api-services" className="nav-link">API Services</Link>
          <Link to="/users" className="nav-link">Users</Link>
        </div>

        <div className="nav-right">
          <span className="user-email">{user?.email}</span>
          <button className="logout-btn" onClick={handleLogout}>
            Logout
          </button>
        </div>
      </div>

      <div className="content">
        {children}
      </div>
    </div>
  );
}
