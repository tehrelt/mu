import { Balance } from "@/components/ui/balance";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { accountStore } from "@/shared/store/account-store";

export const Dashboard = () => {
  const account = accountStore((s) => s.account);

  return (
    <div className="flex gap-6">
      <Card className="w-[180px]">
        <CardHeader>
          <CardTitle>Баланс счёта</CardTitle>
        </CardHeader>
        <CardContent>
          <Balance balance={account!.balance} />
        </CardContent>
      </Card>
      <Card className="w-[180px]">
        <CardHeader>
          <CardTitle></CardTitle>
        </CardHeader>
        <CardContent></CardContent>
      </Card>
    </div>
  );
};
