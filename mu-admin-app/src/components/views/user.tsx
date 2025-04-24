import { fio } from "@/shared/lib/utils";
import { routes } from "@/shared/routes";
import { User } from "@/shared/types/user";
import { ChevronDownIcon, PencilIcon, TrashIcon } from "lucide-react";
import {
  Breadcrumb,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "../ui/breadcrumb";
import { Button } from "../ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";

type Props = {
  user: User;
};

export const UserViewer = ({ user }: Props) => {
  return (
    <>
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbLink href={routes.rate.list}>Тарифы</BreadcrumbLink>
          <BreadcrumbSeparator />
          <BreadcrumbPage>
            <DropdownMenu>
              <DropdownMenuTrigger className="flex items-center gap-x-1">
                {fio(user)}
                <ChevronDownIcon />
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem className="flex items-center gap-x-1">
                  <PencilIcon />
                  Редактировать
                </DropdownMenuItem>
                <DropdownMenuItem className="flex items-center gap-x-1">
                  <TrashIcon />
                  Удалить
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </BreadcrumbPage>
        </BreadcrumbList>
      </Breadcrumb>

      <div className="space-y-2">
        <Card>
          <CardHeader>
            <CardTitle className="">
              <div>{fio(user, { full: true })}</div>
            </CardTitle>
            <CardDescription></CardDescription>
          </CardHeader>
          <CardContent>
            <div>
              <p>
                <span className="text-lg font-bold">Паспорт: </span>
                <span>
                  {user.passportSeries}#{user.passportNumber}
                </span>
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </>
  );
};
