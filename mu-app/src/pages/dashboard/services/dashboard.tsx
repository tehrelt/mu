import { Card } from "@/components/ui/card";
import { CabinetViewer } from "@/components/views/cabinet";
import { cabinetService } from "@/shared/services/cabinet.service";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";

type Props = {};

const ServiceDashboard = (props: Props) => {
  const cabinetId = useParams().id!;

  const cabinetQuery = useQuery({
    queryKey: ["cabinet", cabinetId],
    queryFn: async () => await cabinetService.find(cabinetId),
  });

  if (cabinetQuery.isLoading) {
    return <div>Loading...</div>;
  }

  return <CabinetViewer cabinet={cabinetQuery.data} />;
};

export default ServiceDashboard;
