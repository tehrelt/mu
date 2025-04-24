import { ColumnDef } from "@tanstack/react-table";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card";
import { DataTable } from "../data-table";

type Account = {
  id: string;
  address: string;
  debt: number;
};

const columns: ColumnDef<Account>[] = [
  {
    accessorKey: "id",
    header: "ID",
  },
  {
    accessorKey: "address",
    header: "Адрес",
  },
  {
    accessorKey: "debt",
    header: "Долг",
  },
];

const sampleData: Account[] = [
  {
    id: "1",
    address: "0x1234567890abcdef",
    debt: 1000,
  },
  {
    id: "2",
    address: "0xabcdef1234567890",
    debt: 500,
  },
];

export function ConnectedAccounts() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Подключенные счета</CardTitle>
      </CardHeader>
      <CardContent>
        <DataTable columns={columns} data={sampleData} />
      </CardContent>
    </Card>
  );
}
