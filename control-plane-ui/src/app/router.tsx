import { createBrowserRouter } from "react-router-dom";
import { ProtectedRoute } from "./ProtectedRoute";
import LoginPage from "../pages/LoginPage"
import ClientsPage from "../pages/ClientsPage"
import ClientDetailsPage from "../pages/ClientDetailsPage";
import APIServicesPage from "../pages/APIServicesPage";
import GetStartedPage from "../pages/GetStartedPage";

export const router = createBrowserRouter([
    {
        path: "/",
        element: <GetStartedPage />
    },
    {
        path: "/login",
        element: <LoginPage />,
    },
    {
        path: "/clients",
        element: (
        <ProtectedRoute>
            <ClientsPage />
        </ProtectedRoute>
        ),
    },
    {
    path: "/clients/:id",
    element: (
            <ProtectedRoute>
                <ClientDetailsPage />
            </ProtectedRoute>
        ),
    },
    {
        path: "/api-services",
        element: (
            <ProtectedRoute>
                <APIServicesPage />
            </ProtectedRoute>
        ),
    }
]);

