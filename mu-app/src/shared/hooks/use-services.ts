import { useQuery } from "@tanstack/react-query";
import { accountStore } from "../store/account-store";
import { accountService } from "../services/account.service";

export const useConnectedServices = () => {
  const accountId = accountStore((s) => s.account?.id);

  return useQuery({
    queryKey: ["account", accountId, "services"],
    queryFn: async () => await accountService.services(accountId!),
  });
};
