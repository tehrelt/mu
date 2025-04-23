import { useProfile } from "@/shared/hooks/use-profile";
import { routes } from "@/shared/routes";
import { Navigate, Outlet } from "react-router-dom";

export const ProtectedRoute = () => {
  const data = useProfile();

  if (data.isLoading) {
    return <div>Loading...</div>;
  }

  if (!data.isSuccess) {
    console.log("data.isSuccess", data.isSuccess);
    return <Navigate to={routes.login} replace />;
  }

  return <Outlet />;
};
