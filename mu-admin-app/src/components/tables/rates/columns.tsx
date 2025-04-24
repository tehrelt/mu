import { Button } from "@/components/ui/button";
import { UUID } from "@/components/ui/uuid";
import { routes } from "@/shared/routes";
import { Rate } from "@/shared/types/rate";
import { ColumnDef } from "@tanstack/react-table";
import { Link } from "react-router-dom";

export const rateColumns: ColumnDef<Rate>[] = [
  {
    size: 20,
    maxSize: 30,
    enableResizing: true,
    accessorKey: "id",
    header: "ID",
    cell: ({ getValue }) => {
      const uuid = getValue() as string;
      return (
        <Link to={routes.rate.detail(uuid)}>
          <Button variant={"link"}>
            <UUID uuid={uuid} length={4} />
          </Button>
        </Link>
      );
    },
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
