import { Route, Routes } from "react-router-dom";
import { routes } from "@/shared/routes";
import { ProtectedRoute } from "../layouts/protected-route";
import { Index } from "@/pages";
import { LoginPage } from "@/pages/login";
import { DashboardLayout } from "@/layouts/dashboard";
import { RateCreatePage } from "@/pages/rates/create";

export const RoutesConfig = () => {
  return (
    <Routes>
      <Route element={<ProtectedRoute />}>
        <Route element={<DashboardLayout />}>
          <Route index element={<Index />} />

          <Route path={routes.rate.create} element={<RateCreatePage />} />
        </Route>
      </Route>

      <Route path={routes.login} element={<LoginPage />} />
      <Route path={routes.requestAccess} />
    </Routes>
  );
};
