import { CabinetViewer } from "@/components/views/cabinet";
import { cabinetService } from "@/shared/services/cabinet.service";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";

const ServiceDashboard = () => {
  const cabinetId = useParams().id!;

  const cabinetQuery = useQuery({
    queryKey: ["cabinet", cabinetId],
    queryFn: async () => await cabinetService.find(cabinetId),
  });

  if (cabinetQuery.isLoading) {
    return <div>Loading...</div>;
  }

  return (
    <CabinetViewer
      // @ts-ignore
      cabinet={cabinetQuery.data}
    />
  );
};

export default ServiceDashboard;
