import { DataTable } from "@/components/data-table";
import { usersColumns } from "./columns";
import { UserSnippet } from "@/shared/types/user";

type Props = {
  data: UserSnippet[];
};

export const UserTable = ({ data }: Props) => {
  return (
    <div>
      <DataTable data={data} columns={usersColumns} />
    </div>
  );
};
