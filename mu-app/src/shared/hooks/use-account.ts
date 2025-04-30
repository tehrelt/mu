import { useQuery } from "@tanstack/react-query";
import { accountService } from "../services/account.service";

export const useAccount = (id: string) => {
  const query = useQuery({
    queryKey: ["account", id],
    queryFn: async () => await accountService.find(id),
  });

  return { ...query, account: query.data };
};
