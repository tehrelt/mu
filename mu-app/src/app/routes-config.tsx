import { Route, Routes } from "react-router-dom";
import { MainLayout } from "../layouts/main";
import { Index } from "../pages";
import { SignInPage } from "@/pages/sign-in";

export const RoutesConfig = () => {
  return (
    <Routes>
      <Route element={<MainLayout />}>
        <Route index element={<Index />} />
      </Route>

      <Route path="/sign-in" element={<SignInPage />} />
    </Routes>
  );
};
