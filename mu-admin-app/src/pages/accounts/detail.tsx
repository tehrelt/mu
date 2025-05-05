import { AccountViewer } from "@/components/views/account";
import { accountService } from "@/shared/services/account.service";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { UserNotFound } from "./not-found";

export const AccountDetailPage = () => {
  const { id } = useParams();

  const query = useQuery({
    queryKey: ["account", { id }],
    queryFn: async () => accountService.find(id!),
  });

  if (query.isLoading) return <>Loading...</>;
  if (query.isError) return <UserNotFound id={id} />;
  if (!query.data) return <p>Query data is null</p>;

  return <AccountViewer account={query.data} />;
};
