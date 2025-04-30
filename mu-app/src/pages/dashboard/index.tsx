import { DataTable } from "@/components/data-table";
import { Balance } from "@/components/ui/balance";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { PaymentStatusBadge } from "@/components/ui/payment-status-badge";
import { AccountBalanceCard } from "@/components/widgets/account-balance-card";
import { useAccount } from "@/shared/hooks/use-account";
import { datef } from "@/shared/lib/utils";
import { routes } from "@/shared/routes";
import { accountService } from "@/shared/services/account.service";
import { accountStore } from "@/shared/store/account-store";
import { Payment, PaymentStatus } from "@/shared/types/payment";
import { useQuery } from "@tanstack/react-query";
import { ColumnDef } from "@tanstack/react-table";
import { Rows } from "lucide-react";
import { Link } from "react-router-dom";

export const Dashboard = () => {
  const accountId = accountStore((s) => s.account!.id);

  const limit = 3;

  const paymentsQuery = useQuery({
    queryKey: ["payments"],
    queryFn: async () => await accountService.payments(accountId, { limit }),
  });

  return (
    <div className="flex gap-6">
      <div>
        <AccountBalanceCard id={accountId} />
      </div>
      <div>
        <Card className="">
          <CardHeader>
            <CardTitle>
              Последние {Math.min(limit, paymentsQuery.data?.payments.length)}{" "}
              платежа
            </CardTitle>
          </CardHeader>
          <CardContent className="">
            <div>
              {paymentsQuery.data && (
                <DataTable
                  data={paymentsQuery.data.payments}
                  columns={cols}
                  className="h-[200px]"
                />
              )}
            </div>
          </CardContent>
          <CardFooter className="flex w-full">
            <Link
              to={routes.dashboard.account.transactionHistory}
              className="w-full"
            >
              <Button className="w-full" variant={"outline"}>
                Показать полную историю транзакции
              </Button>
            </Link>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
};

const cols: ColumnDef<Payment>[] = [
  {
    accessorKey: "status",
    header: "Статус",
    cell: ({ getValue }) => {
      const status = getValue() as PaymentStatus;
      return (
        <div className="flex justify-center">
          <PaymentStatusBadge status={status} />
        </div>
      );
    },
  },
  {
    accessorKey: "amount",
    header: "Сумма",
    cell: ({ getValue }) => {
      const balance = getValue() as number;
      return <Balance balance={balance} />;
    },
  },
  {
    accessorKey: "createdAt",
    header: "Создан",
    cell: ({ getValue }) => {
      const date = getValue() as Date;
      return datef(date, "DD/MM/YYYY");
    },
  },
  {
    accessorKey: "id",
    header: "",
    cell: ({ getValue, row }) => {
      const status = row.original.status;
      if (status != "pending") {
        return null;
      }

      const id = getValue() as string;
      return (
        <Link to={routes.billing.process(id)} target="_blank">
          <Button variant={"outline"}>Перейти к оплате</Button>
        </Link>
      );
    },
  },
];
