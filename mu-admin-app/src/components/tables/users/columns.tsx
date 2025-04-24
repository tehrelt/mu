import { Button } from "@/components/ui/button";
import { UUID } from "@/components/ui/uuid";
import { datef } from "@/lib/utils";
import { routes } from "@/shared/routes";
import { UserSnippet } from "@/shared/types/user";
import { ColumnDef } from "@tanstack/react-table";
import { Link } from "react-router-dom";

export const usersColumns: ColumnDef<UserSnippet>[] = [
  {
    size: 20,
    maxSize: 30,
    enableResizing: true,
    accessorKey: "id",
    header: "ID",
    cell: ({ getValue }) => {
      const uuid = getValue() as string;
      return (
        <Link to={routes.users.detail(uuid)}>
          <Button variant={"link"}>
            <UUID uuid={uuid} length={4} />
          </Button>
        </Link>
      );
    },
  },
  {
    accessorKey: "lastName",
    header: "Last Name",
  },
  {
    accessorKey: "firstName",
    header: "First Name",
  },
  {
    accessorKey: "middleName",
    header: "Middle Name",
  },
  {
    accessorKey: "email",
    header: "Email",
  },
  {
    accessorKey: "phone",
    header: "Phone Number",
  },
  {
    accessorKey: "createdAt",
    header: "Created At",
    cell: ({ getValue }) => {
      const date = getValue() as Date;
      return datef(date);
    },
  },
  {
    accessorKey: "updatedAt",
    header: "Updated At",
    cell: ({ getValue }) => {
      const val = getValue();
      if (!val) return null;
      const date = val as Date;
      return datef(date);
    },
  },
];
