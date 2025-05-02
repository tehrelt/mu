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
import MasonryLayout from "@/components/ui/masonry-layout";
import { PaymentStatusBadge } from "@/components/ui/payment-status-badge";
import { LogsChart } from "@/components/views/cabinet-logs-chart";
import { AccountBalanceCard } from "@/components/widgets/account-balance-card";
import { PendingPaymentsAlert } from "@/components/widgets/pending-payments-alert";
import { datef } from "@/shared/lib/utils";
import { routes } from "@/shared/routes";
import { accountService } from "@/shared/services/account.service";
import { accountStore } from "@/shared/store/account-store";
import { Payment, PaymentStatus } from "@/shared/types/payment";
import { useQuery } from "@tanstack/react-query";
import { ColumnDef } from "@tanstack/react-table";
import { Link } from "react-router-dom";

export const Dashboard = () => {
  const accountId = accountStore((s) => s.account!.id);

  const limit = 10;
  const paymentsQuery = useQuery({
    queryKey: ["account", "payments"],
    queryFn: async () => await accountService.payments(accountId, { limit }),
  });
  const pendingPaymentsQuery = useQuery({
    queryKey: ["account", "payments", { status: "pending" }],
    queryFn: async () =>
      await accountService.payments(accountId, { limit, status: "pending" }),
  });

  const logsQuery = useQuery({
    queryKey: ["account", "logs"],
    queryFn: async () => await accountService.logs(accountId),
  });

  return (
    <MasonryLayout>
      {pendingPaymentsQuery.data &&
        pendingPaymentsQuery.data.payments.length != 0 && (
          <PendingPaymentsAlert />
        )}
      <AccountBalanceCard id={accountId} />
      {paymentsQuery.data && (
        <Operations payments={paymentsQuery.data.payments} limit={limit} />
      )}
      {logsQuery.data && <LogsChart services={logsQuery.data.items} />}
    </MasonryLayout>
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

const Operations = ({
  payments,
  limit,
}: {
  payments: Payment[];
  limit: number;
}) => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>
          Последние {Math.min(limit, payments.length)} платежа
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="max-w-full overflow-x-auto">
          <DataTable columns={cols} data={payments} />
        </div>
      </CardContent>
      <CardFooter>
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
  );
};
