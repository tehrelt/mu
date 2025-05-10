import { Balance } from "@/components/ui/balance";
import { routes } from "@/shared/routes";
import { localizeServiceType, Rate } from "@/shared/types/rate";
import {
  Breadcrumb,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "../../ui/breadcrumb";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../../ui/card";

type Props = {
  rate: Rate;
};

const RateViewer = ({ rate }: Props) => (
  <>
    <Breadcrumb>
      <BreadcrumbList>
        <BreadcrumbLink href={routes.rate.list}>Тарифы</BreadcrumbLink>
        <BreadcrumbSeparator />
        <BreadcrumbPage>{rate.name}</BreadcrumbPage>
      </BreadcrumbList>
    </Breadcrumb>

    <div className="space-y-2">
      <Card>
        <CardHeader>
          <CardTitle className="">
            <div>Услуга: {rate.name}</div>
          </CardTitle>
          <CardDescription>
            {/* <Button variant={"outline"}>Обновить цену</Button> */}
            <div className="text-muted-foreground text-sm">
              <span>ID тарифа: {rate.id}</span>
            </div>
          </CardDescription>
        </CardHeader>
        <CardContent>
          <p>Тип услуги: {localizeServiceType(rate.serviceType)}</p>
          <p>Единица измерения: {rate.measureUnit}</p>
          <p className="">
            Цена: <Balance balance={rate.rate} />
          </p>
        </CardContent>
      </Card>

      {/* <ConnectedAccounts /> */}
    </div>
  </>
);

export default RateViewer;
