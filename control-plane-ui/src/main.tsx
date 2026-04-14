import { createRoot } from "react-dom/client";
import { RouterProvider } from "react-router-dom";
import { router } from "./app/router";
import { Providers } from "./app/providers";
import "./shared/styles/global.css";
import "./shared/styles/layout.css";
import "./shared/styles/components.css";

createRoot(document.getElementById('root')!).render(
  <Providers>
    <RouterProvider router={router} />
  </Providers>,
)
