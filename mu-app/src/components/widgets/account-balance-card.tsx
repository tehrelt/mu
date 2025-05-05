import { useAccount } from "@/shared/hooks/use-account";
import { Balance } from "../ui/balance";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card";

export const AccountBalanceCard = ({ id }: { id: string }) => {
  const accountQuery = useAccount(id);
  return (
    <Card>
      <CardHeader>
        <CardTitle>Баланс счёта</CardTitle>
      </CardHeader>
      <CardContent className="flex justify-end">
        {accountQuery.account ? (
          <Balance balance={accountQuery.account.balance} />
        ) : (
          <p>Loading...</p>
        )}
      </CardContent>
    </Card>
  );
};
