import { useQuery } from "@tanstack/react-query";
import { accountService } from "../services/account.service";
import { accountStore } from "../store/account-store";

export const useAccounts = () => {
  const query = useQuery({
    queryKey: ["accounts"],
    queryFn: async () => await accountService.getUserAccounts(),
  });

  const selectedAccount = accountStore((s) => s.account);
  const selectAccount = accountStore((s) => s.select);
  const clearAccount = accountStore((s) => s.clear);

  if (query.isError) clearAccount();

  if (query.isSuccess && query.data && selectedAccount) {
    const acc = query.data.accounts.find((a) => a.id === selectedAccount.id);
    if (!acc) clearAccount();
    else selectAccount(acc);
  }

  return { ...query, accounts: query.data };
};
