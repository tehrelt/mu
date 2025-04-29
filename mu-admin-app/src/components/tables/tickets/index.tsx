import { DataTable } from "@/components/data-table";
import { ticketColumns } from "./columns";
import { TicketHeader } from "@/shared/types/ticket";

type Props = {
  data: TicketHeader[];
};

const TicketTable = ({ data }: Props) => {
  return (
    <div>
      <DataTable data={data} columns={ticketColumns} />
    </div>
  );
};

export default TicketTable;
