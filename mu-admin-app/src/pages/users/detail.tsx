import { userService } from "@/shared/services/users.service";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { UserNotFound } from "../tickets/not-found";
import { UserViewer } from "@/components/views/user";

export const UserDetailsPage = () => {
  const { id } = useParams();

  const query = useQuery({
    queryKey: ["user", { id }],
    queryFn: async () => userService.find(id!),
  });

  if (query.isLoading) return <>Loading...</>;
  if (query.isError) return <UserNotFound id={id} />;
  if (!query.data) return <p>Query data is null</p>;

  return <UserViewer user={query.data} />;
};
