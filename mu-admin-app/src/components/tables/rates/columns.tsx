import { Rate } from "@/shared/types/rate";
import { ColumnDef } from "@tanstack/react-table";

export const rateColumns: ColumnDef<Rate>[] = [
  {
    accessorKey: "id",
    header: "ID",
  },
  {
    accessorKey: "name",
    header: "Name",
  },
  {
    accessorKey: "measureUnit",
    header: "Measure Unit",
  },
  {
    accessorKey: "rate",
    header: "Rate",
  },
];
