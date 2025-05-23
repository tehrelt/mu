import { useProfile } from "@/shared/hooks/use-profile";
import { routes } from "@/shared/routes";
import { Navigate, Outlet } from "react-router-dom";

export const ProtectedRoute = () => {
  const data = useProfile();

  if (data.isLoading) {
    return <p>Loading</p>;
  }

  if (!data.isLoading && !data.isSuccess) {
    return <Navigate to={routes.signIn} replace />;
  }

  return <Outlet />;
};
