import { DataTable } from "@/components/data-table";
import { rateColumns } from "./columns";
import { TicketHeader } from "@/shared/types/ticket";

type Props = {
  data: TicketHeader[];
};

const TicketTable = ({ data }: Props) => {
  return (
    <div>
      <DataTable data={data} columns={rateColumns} />
    </div>
  );
};

export default TicketTable;
