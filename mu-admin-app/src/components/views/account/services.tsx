import { DataTable } from "@/components/data-table";
import { Balance } from "@/components/ui/balance";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { accountService } from "@/shared/services/account.service";
import { localizeServiceType, Rate, ServiceType } from "@/shared/types/rate";
import { useQuery } from "@tanstack/react-query";
import { ColumnDef } from "@tanstack/react-table";

type Props = {
  accountId: string;
};

const cols: ColumnDef<Rate>[] = [
  {
    accessorKey: "name",
  },
  {
    accessorKey: "serviceType",
    cell: ({ getValue }) => {
      const val = getValue() as ServiceType;
      return localizeServiceType(val);
    },
  },
  {
    accessorKey: "rate",
    cell: ({ getValue }) => {
      const val = getValue() as number;
      return <Balance balance={val} />;
    },
  },
];

const ConnectedServices = ({ accountId }: Props) => {
  const query = useQuery({
    queryKey: ["account", { id: accountId }, "services"],
    queryFn: async () => await accountService.services(accountId),
  });

  return (
    <Card>
      <CardHeader>
        <CardTitle>Подключенные услуги</CardTitle>
      </CardHeader>
      <CardContent>
        {query.data && <DataTable data={query.data.services} columns={cols} />}
      </CardContent>
    </Card>
  );
};

export default ConnectedServices;
