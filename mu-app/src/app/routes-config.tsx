import { Route, Routes } from "react-router-dom";
import { MainLayout } from "../layouts/main";
import { Index } from "../pages";
import { SignInPage } from "@/pages/sign-in";
import { routes } from "@/shared/routes";
import { Dashboard } from "@/pages/dashboard";
import { DashboardLayout } from "@/layouts/dashboard";
import { SignUpPage } from "@/pages/sign-up";
import { NewTicketPage } from "@/pages/dashboard/new-ticket";
import { IntegrationsSettingsPage } from "@/pages/settings/integrations";
import { AddFundsPage } from "@/pages/dashboard/add-funds";
import { ProtectedRoute } from "@/layouts/protected-route";
import { ProcessPaymentPage } from "@/pages/billing/process";

export const RoutesConfig = () => {
  return (
    <Routes>
      <Route element={<MainLayout />}>
        <Route index element={<Index />} />
      </Route>

      <Route path={routes.dashboard.index} element={<DashboardLayout />}>
        <Route index element={<Dashboard />} />
        <Route path={routes.dashboard.newTicket} element={<NewTicketPage />} />
        <Route
          path={routes.dashboard.settings.integrations}
          element={<IntegrationsSettingsPage />}
        />
        <Route path={routes.dashboard.addFunds} element={<AddFundsPage />} />
      </Route>

      <Route path={routes.billing.process()} element={<ProcessPaymentPage />} />
      <Route path={routes.signIn} element={<SignInPage />} />
      <Route path={routes.signUp} element={<SignUpPage />} />
    </Routes>
  );
};
