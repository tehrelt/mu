import RateViewer from "@/components/views/rate";
import { userService } from "@/shared/services/users.service";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { UserNotFound } from "./not-found";
import { UserViewer } from "@/components/views/user";

type Props = {};

export const UserDetailsPage = (props: Props) => {
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
