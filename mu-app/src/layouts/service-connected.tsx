import { useConnectedServices } from "@/shared/hooks/use-services";
import { routes } from "@/shared/routes";
import { Navigate, Outlet, useParams } from "react-router-dom";

export const IsServiceConnectedLayout = () => {
  const serviceId = useParams().id!;
  const servicesQuery = useConnectedServices();

  if (servicesQuery.isLoading) {
    return <p>Loading...</p>;
  }

  if (!servicesQuery.isSuccess) {
    return <Navigate to={routes.dashboard.index} />;
  }

  if (servicesQuery.data) {
    if (!servicesQuery.data.services.find((svc) => svc.id === serviceId)) {
      return <Navigate to={routes.dashboard.index} />;
    }
  }

  return <Outlet />;
};
