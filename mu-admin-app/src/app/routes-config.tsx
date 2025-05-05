import { DashboardLayout } from "@/layouts/dashboard";
import { Index } from "@/pages";
import { AccountDetailPage } from "@/pages/accounts/detail";
import { LoginPage } from "@/pages/login";
import { RateCreatePage } from "@/pages/rates/create";
import RateDetailsPage from "@/pages/rates/detail";
import { RateListPage } from "@/pages/rates/list";
import { RateNotFoundPage } from "@/pages/rates/not-found";
import { TicketListPage } from "@/pages/tickets/list";
import { UserListPage } from "@/pages/users/list";
import { routes } from "@/shared/routes";
import { Route, Routes } from "react-router-dom";
import { ProtectedRoute } from "../layouts/protected-route";
import { UserDetailsPage } from "@/pages/users/detail";
import { TicketDetailsPage } from "@/pages/tickets/detail";

export const RoutesConfig = () => {
  return (
    <Routes>
      <Route element={<ProtectedRoute />}>
        <Route element={<DashboardLayout />}>
          <Route index element={<Index />} />

          <Route path={routes.rate.index}>
            <Route path={routes.rate.create} element={<RateCreatePage />} />
            <Route path={routes.rate.list} element={<RateListPage />} />
            <Route path={routes.rate.detail()} element={<RateDetailsPage />} />
            <Route path={"*"} element={<RateNotFoundPage />} />
          </Route>

          <Route path={routes.users.index}>
            <Route path={routes.users.list} element={<UserListPage />} />
            <Route path={routes.users.detail()} element={<UserDetailsPage />} />
          </Route>

          <Route path={routes.tickets.index}>
            <Route path={routes.tickets.list} element={<TicketListPage />} />
            <Route
              path={routes.tickets.detail()}
              element={<TicketDetailsPage />}
            />
          </Route>

          <Route path={routes.accounts.index}>
            <Route
              path={routes.accounts.detail()}
              element={<AccountDetailPage />}
            />
          </Route>
        </Route>
      </Route>

      <Route path={routes.login} element={<LoginPage />} />
      <Route path={routes.requestAccess} />
    </Routes>
  );
};
