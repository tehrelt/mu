import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { UserNotFound } from "../tickets/not-found";
import { TicketViewer } from "@/components/views/ticket";
import { ticketService } from "@/shared/services/tickets.service";

type Props = {};

export const TicketDetailsPage = (props: Props) => {
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
