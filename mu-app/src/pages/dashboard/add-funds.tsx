import { DataTable } from "@/components/data-table";
import { Balance } from "@/components/ui/balance";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import MasonryLayout from "@/components/ui/masonry-layout";
import { PaymentStatusBadge } from "@/components/ui/payment-status-badge";
import { datef } from "@/shared/lib/utils";
import { routes } from "@/shared/routes";
import { accountService } from "@/shared/services/account.service";
import { billingService } from "@/shared/services/billing.service";
import { accountStore } from "@/shared/store/account-store";
import {
  Payment,
  PaymentStatus,
  paymentStatusSchema,
} from "@/shared/types/payment";
import { useMutation, useQuery } from "@tanstack/react-query";
import { ColumnDef } from "@tanstack/react-table";
import React from "react";
import { Link, useNavigate } from "react-router-dom";
import { toast } from "sonner";

type Props = {};

const presets = [100, 250, 500, 1000];

export function AddFundsPage() {
  const [amount, setAmount] = React.useState(0);
  const account = accountStore((s) => s.account);
  const navigate = useNavigate();

  const { mutate } = useMutation({
    mutationKey: ["payment-create", { amount, accountId: account!.id }],
    mutationFn: (amount: number) =>
      billingService.create({
        amount: amount,
        accountId: account!.id,
      }),
    onSuccess: (res) => {
      toast.success("Платеж создан");
      navigate(routes.billing.process(res.id));
    },
    onSettled: () => {
      setAmount(0);
    },
  });

  const paymentsQuery = useQuery({
    queryKey: [
      "account",
      account?.id,
      "payments",
      { status: paymentStatusSchema.enum.pending },
    ],
    queryFn: async () =>
      await accountService.payments(account!.id, {
        status: paymentStatusSchema.enum.pending,
      }),
  });

  return (
    <MasonryLayout>
      <Card className="">
        <CardHeader>
          <CardTitle>Пополнить счёт</CardTitle>
          <CardDescription>
            Пополните на готовую сумму или введите вручную
          </CardDescription>
        </CardHeader>
        <CardContent className="grid">
          <div className="flex flex-col items-center gap-4">
            <div className="flex gap-x-2">
              {presets.map((preset) => (
                <Button
                  variant={"outline"}
                  key={preset}
                  onClick={() => setAmount(preset)}
                >
                  <Balance balance={preset} />
                </Button>
              ))}
            </div>
            <div className="text-muted-foreground flex">или</div>
            <Input
              value={amount}
              onChange={(e) => setAmount(+e.target.value)}
              type="number"
            />
          </div>
        </CardContent>
        <CardFooter className="flex justify-end">
          <CardAction>
            <Button onClick={() => mutate(amount)}>Пополнить</Button>
          </CardAction>
        </CardFooter>
      </Card>

      {paymentsQuery.data && (
        <PendingPayments payments={paymentsQuery.data?.payments} />
      )}
    </MasonryLayout>
  );
}

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
    header: "Действия",
    cell: ({ getValue }) => {
      const id = getValue() as string;
      return (
        <Link to={routes.billing.process(id)} target="_blank">
          <Button variant={"outline"}>Перейти к оплате</Button>
        </Link>
      );
    },
  },
];

const PendingPayments = ({ payments }: { payments: Payment[] }) => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Платежи в ожидании оплаты</CardTitle>
      </CardHeader>
      <CardContent>
        <DataTable data={payments} columns={cols} />
      </CardContent>
    </Card>
  );
};
