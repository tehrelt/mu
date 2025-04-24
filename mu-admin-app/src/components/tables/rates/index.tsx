import { DataTable } from "@/components/data-table";
import { Rate } from "@/shared/types/rate";
import { rateColumns } from "./columns";

type Props = {
  data: Rate[];
};

const RateTable = ({ data }: Props) => {
  return (
    <div>
      <DataTable data={data} columns={rateColumns} />
    </div>
  );
};

export default RateTable;
