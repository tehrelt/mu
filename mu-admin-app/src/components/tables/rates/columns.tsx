import { Balance } from "@/components/ui/balance";
import { Button } from "@/components/ui/button";
import { UUID } from "@/components/ui/uuid";
import { routes } from "@/shared/routes";
import { localizeServiceType, Rate, ServiceType } from "@/shared/types/rate";
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
    header: "Наименование организации",
  },
  {
    accessorKey: "serviceType",
    header: "Вид услуги",
    cell: ({ getValue }) => {
      const serviceType = getValue() as ServiceType;
      return <span>{localizeServiceType(serviceType)}</span>;
    },
  },
  {
    accessorKey: "measureUnit",
    header: "Ед. измерения",
  },
  {
    accessorKey: "rate",
    header: "Цена",
    cell: ({ getValue }) => {
      const price = getValue() as number;
      return <Balance balance={price} />;
    },
  },
];
