import { DataTable } from "@/components/data-table";
import { Rate } from "@/shared/types/rate";
import { usersColumns } from "./columns";
import { User } from "@/shared/types/user";

type Props = {
  data: User[];
};

export const UserTable = ({ data }: Props) => {
  return (
    <div>
      <DataTable data={data} columns={usersColumns} />
    </div>
  );
};
