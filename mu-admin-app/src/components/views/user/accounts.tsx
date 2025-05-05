import { DataTable } from "@/components/data-table";
import { Balance } from "@/components/ui/balance";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { UUID } from "@/components/ui/uuid";
import { datef } from "@/shared/lib/utils";
import { routes } from "@/shared/routes";
import { userService } from "@/shared/services/users.service";
import { AccountInfo } from "@/shared/types/user";
import { useQuery } from "@tanstack/react-query";
import { ColumnDef } from "@tanstack/react-table";
import React from "react";
import { Link } from "react-router-dom";

type Props = {
  userId: string;
};

const columns: ColumnDef<AccountInfo>[] = [
  {
    size: 20,
    maxSize: 30,
    enableResizing: true,
    accessorKey: "id",
    header: "ID",
    cell: ({ getValue }) => {
      const id = getValue() as string;
      return (
        <Link to={routes.accounts.detail(id)}>
          <Button variant={"link"}>
            <UUID uuid={id} length={4} />
          </Button>
        </Link>
      );
    },
  },
  {
    accessorKey: "house.address",
    header: "Адрес",
  },
  {
    accessorKey: "balance",
    header: "Баланс",
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

export const UserAccounts = ({ userId }: Props) => {
  const query = useQuery({
    queryKey: ["user", { userId }, "accounts"],
    queryFn: async () => await userService.accountsOfUser(userId),
  });

  return (
    <div>
      <Card>
        <CardHeader>
          <CardTitle>Счета пользователя</CardTitle>
        </CardHeader>
        <CardContent>
          {query.data && (
            <DataTable data={query.data.accounts} columns={columns} />
          )}
        </CardContent>
      </Card>
    </div>
  );
};
