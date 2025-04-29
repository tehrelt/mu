import { TicketStatus } from "@/components/ticket-status";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { UUID } from "@/components/ui/uuid";
import { routes } from "@/shared/routes";
import { localizeServiceType } from "@/shared/types/rate";
import {
  localizeTicketStatus,
  localizeTicketType,
  TicketHeader,
  TicketStatusEnum,
  TicketTypeEnum,
} from "@/shared/types/ticket";
import { ColumnDef } from "@tanstack/react-table";
import { Link } from "react-router-dom";

export const rateColumns: ColumnDef<TicketHeader>[] = [
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
    accessorKey: "ticket_status",
    header: "Статус",
    cell: ({ getValue }) => {
      const val = getValue() as TicketStatusEnum;
      return <TicketStatus val={val} />;
    },
  },
  {
    accessorKey: "ticket_type",
    header: "type",
    cell: ({ getValue }) => {
      const type = getValue() as TicketTypeEnum;
      return <Badge>{localizeTicketType(type)}</Badge>;
    },
  },
];
