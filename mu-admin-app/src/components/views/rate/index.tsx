import { Balance } from "@/components/ui/balance";
import { routes } from "@/shared/routes";
import { localizeServiceType, Rate } from "@/shared/types/rate";
import { ChevronDownIcon, PencilIcon, TrashIcon } from "lucide-react";
import {
  Breadcrumb,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "../../ui/breadcrumb";
import { Button } from "../../ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../../ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "../../ui/dropdown-menu";
import { ConnectedAccounts } from "./connected-accounts";

type Props = {
  rate: Rate;
};

const RateViewer = ({ rate }: Props) => (
  <>
    <Breadcrumb>
      <BreadcrumbList>
        <BreadcrumbLink href={routes.rate.list}>Тарифы</BreadcrumbLink>
        <BreadcrumbSeparator />
        <BreadcrumbPage>
          <DropdownMenu>
            <DropdownMenuTrigger className="flex items-center gap-x-1">
              {rate.name}
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
            <div>{rate.name}</div>
          </CardTitle>
          <CardDescription>
            <Button variant={"outline"}>Обновить цену</Button>
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="text-muted-foreground text-sm">
            <span>ID тарифа: {rate.id}</span>
          </div>
          <p>Тип услуги: {localizeServiceType(rate.serviceType)}</p>
          <p>Единица измерения: {rate.measureUnit}</p>
          <p className="">
            Цена: <Balance balance={rate.rate} />
          </p>
        </CardContent>
      </Card>

      <ConnectedAccounts />
    </div>
  </>
);

export default RateViewer;
