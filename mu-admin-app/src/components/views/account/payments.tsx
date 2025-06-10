import { DataTable } from "@/components/data-table";
import { Badge } from "@/components/ui/badge";
import { Balance } from "@/components/ui/balance";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { UUID } from "@/components/ui/uuid";
import { datef } from "@/shared/lib/utils";
import { accountService } from "@/shared/services/account.service";
import { Payment, PaymentStatus } from "@/shared/types/account";
import { useQuery } from "@tanstack/react-query";
import { ColumnDef } from "@tanstack/react-table";

type Props = {
  accId: string;
};

const columns: ColumnDef<Payment>[] = [
  {
    size: 20,
    maxSize: 30,
    enableResizing: true,
    accessorKey: "id",
    header: "ID",
    cell: ({ getValue }) => {
      const id = getValue() as string;
      return <UUID uuid={id} length={4} />;
    },
  },
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
    accessorKey: "updatedAt",
    header: "Обновлён",
    cell: ({ getValue }) => {
      const val = getValue();
      if (!val) return;

      const date = val as Date;
      return datef(date);
    },
  },
];

const AccountPayments = ({ accId }: Props) => {
  const query = useQuery({
    queryKey: ["account", { accId }, "payments"],
    queryFn: async () => await accountService.payments(accId),
  });

  return (
    <div>
      <Card>
        <CardHeader>
          <CardTitle>Платежи</CardTitle>
        </CardHeader>
        <CardContent>
          {query.data && (
            <DataTable data={query.data.payments} columns={columns} />
          )}
        </CardContent>
        <CardFooter>
          <CardDescription className="space-x-2 text-xs text-center">
            <span>
              <Badge variant="pending" className="p-1" /> - Ожидание платежа
            </span>
            <span>
              <Badge variant="success" className="p-1" /> - Платёж принят
            </span>
            <span>
              <Badge variant="destructive" className="p-1" /> - Платёж отклонён
            </span>
          </CardDescription>
        </CardFooter>
      </Card>
    </div>
  );
};

const PaymentStatusBadge = ({ status }: { status: PaymentStatus }) => {
  const variant =
    status === "pending"
      ? "pending"
      : status === "success"
      ? "success"
      : "destructive";

  const tooltip =
    status === "pending"
      ? "Ожидание платежа"
      : status === "success"
      ? "Платёж принят"
      : "Платёж отклонен";

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger>
          <Badge className="p-1" variant={variant} />
        </TooltipTrigger>
        <TooltipContent>{tooltip}</TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
};

export default AccountPayments;
