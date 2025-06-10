import { Route, Routes } from "react-router-dom";
import { MainLayout } from "../layouts/main";
import { Index } from "../pages";
import { SignInPage } from "@/pages/sign-in";
import { routes } from "@/shared/routes";
import { Dashboard } from "@/pages/dashboard";
import { AccountCheck, DashboardLayout } from "@/layouts/dashboard";
import { SignUpPage } from "@/pages/sign-up";
import { NewTicketPage } from "@/pages/dashboard/new-ticket";
import { IntegrationsSettingsPage } from "@/pages/settings/integrations";
import { AddFundsPage } from "@/pages/dashboard/add-funds";
import { ProcessPaymentPage } from "@/pages/billing/process";
import NotFoundPage from "@/pages/not-found-page";
import ServiceDashboard from "@/pages/dashboard/services/dashboard";
import { ProtectedRoute } from "@/layouts/protected-route";

export const RoutesConfig = () => {
  return (
    <Routes>
      <Route element={<MainLayout />}>
        <Route index element={<Index />} />
        <Route path={routes.unmatched} element={<NotFoundPage />} />
      </Route>

      <Route path={routes.dashboard.index} element={<DashboardLayout />}>
        <Route element={<ProtectedRoute />}>
          <Route
            path={routes.dashboard.newTicket}
            element={<NewTicketPage />}
          />
          <Route element={<AccountCheck />}>
            <Route index element={<Dashboard />} />
            <Route
              path={routes.dashboard.addFunds}
              element={<AddFundsPage />}
            />

            <Route
              path={routes.dashboard.cabinet.dashboard()}
              element={<ServiceDashboard />}
            />
          </Route>
          <Route
            path={routes.dashboard.settings.integrations}
            element={<IntegrationsSettingsPage />}
          />
        </Route>
      </Route>

      <Route path={routes.billing.process()} element={<ProcessPaymentPage />} />
      <Route path={routes.signIn} element={<SignInPage />} />
      <Route path={routes.signUp} element={<SignUpPage />} />
    </Routes>
  );
};
