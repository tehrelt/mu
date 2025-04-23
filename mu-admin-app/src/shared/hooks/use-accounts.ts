import { useQuery } from "@tanstack/react-query";
import { accountService } from "../services/account.service";

export const useAccounts = () => {
  const query = useQuery({
    queryKey: ["accounts"],
    queryFn: async () => await accountService.getUserAccounts(),
  });

  return { ...query, accounts: query.data };
};
