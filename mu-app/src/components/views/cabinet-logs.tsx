import { datef } from "@/shared/lib/utils";
import { CabinetLog } from "@/shared/types/cabinet";
import { ColumnDef } from "@tanstack/react-table";
import { DataTable } from "../data-table";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card";

const cols = (unit?: string): ColumnDef<CabinetLog>[] => [
  {
    accessorKey: "createdAt",
    header: "Когда?",
    cell: ({ getValue }) => {
      const date = getValue() as Date;
      return datef(date, "DD/MM/YYYY hh:mm:ss");
    },
  },
  {
    accessorKey: "consumed",
    header: `Сколько потреблено ${unit ? `(${unit})` : ""}`,
    cell: ({ getValue }) => {
      const consumed = getValue() as number;
      return <p>{consumed}</p>;
    },
  },
];

type Props = {
  logs: CabinetLog[];
  unit?: string;
};

export const ConsumptionHistory = ({ logs, unit }: Props) => {
  return (
    <div>
      <Card>
        <CardHeader>
          <CardTitle>История потрбления</CardTitle>
        </CardHeader>
        <CardContent>
          <DataTable columns={cols(unit)} data={logs} />
        </CardContent>
      </Card>
    </div>
  );
};
