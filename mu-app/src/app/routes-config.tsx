import { Route, Routes } from "react-router-dom";
import { MainLayout } from "../layouts/main";
import { Index } from "../pages";
import { SignInPage } from "@/pages/sign-in";
import { routes } from "@/shared/routes";
import { Dashboard } from "@/pages/dashboard";
import { DashboardLayout } from "@/layouts/dashboard";
import { SignUpPage } from "@/pages/sign-up";
import { NewTicketPage } from "@/pages/dashboard/new-ticket";

export const RoutesConfig = () => {
  return (
    <Routes>
      <Route element={<MainLayout />}>
        <Route index element={<Index />} />
      </Route>

      <Route path={routes.dashboard} element={<DashboardLayout />}>
        <Route index element={<Dashboard />} />
        <Route path={routes.newTicket} element={<NewTicketPage />} />
      </Route>

      <Route path={routes.signIn} element={<SignInPage />} />
      <Route path={routes.signUp} element={<SignUpPage />} />
    </Routes>
  );
};
