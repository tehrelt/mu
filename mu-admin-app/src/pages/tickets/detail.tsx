import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { TicketViewer } from "@/components/views/ticket";
import { ticketService } from "@/shared/services/tickets.service";

export const TicketDetailsPage = () => {
  const { id } = useParams();

  const query = useQuery({
    queryKey: ["ticket", { id }],
    queryFn: async () => ticketService.find(id!),
  });

  console.log("ticket page");
  if (query.isLoading) return <>Loading...</>;
  if (query.isError) return <p>{query.error.message}</p>;
  if (!query.data) return <p>Query data is null</p>;

  return <TicketViewer ticket={query.data} />;
};
