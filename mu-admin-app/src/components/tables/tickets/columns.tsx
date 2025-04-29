import { TicketStatus } from "@/components/ticket-status";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { UUID } from "@/components/ui/uuid";
import { routes } from "@/shared/routes";
import {
  localizeTicketType,
  TicketHeader,
  TicketStatusEnum,
  TicketTypeEnum,
} from "@/shared/types/ticket";
import { ColumnDef } from "@tanstack/react-table";
import { Link } from "react-router-dom";

export const ticketColumns: ColumnDef<TicketHeader>[] = [
  {
    size: 20,
    maxSize: 30,
    enableResizing: true,
    accessorKey: "id",
    header: "ID",
    cell: ({ getValue }) => {
      const id = getValue() as string;
      return (
        <Link to={routes.tickets.detail(id)}>
          <Button variant={"link"}>
            <UUID uuid={id} length={4} />
          </Button>
        </Link>
      );
    },
  },
  {
    accessorKey: "ticketType",
    header: "Вид заявки",
    cell: ({ getValue }) => {
      const type = getValue() as TicketTypeEnum;
      return <Badge>{localizeTicketType(type)}</Badge>;
    },
  },
  {
    accessorKey: "ticketStatus",
    header: "Статус",
    cell: ({ getValue }) => {
      const val = getValue() as TicketStatusEnum;
      return <TicketStatus val={val} />;
    },
  },
];
