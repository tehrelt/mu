import { routes } from "@/shared/routes";
import { AccountInfo } from "@/shared/types/user";
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
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { UUID } from "@/components/ui/uuid";
import AccountPayments from "./payments";

type Props = {
  account: AccountInfo;
};

export const AccountViewer = ({ account }: Props) => {
  return (
    <>
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbLink href={routes.users.list}>Счета</BreadcrumbLink>
          <BreadcrumbSeparator />
          <BreadcrumbPage>{account.house.address}</BreadcrumbPage>
        </BreadcrumbList>
      </Breadcrumb>

      <div className="flex space-x-4">
        <div>
          <AccountCard acc={account} />
        </div>
        <div>
          <AccountPayments accId={account.id} />
        </div>
      </div>
    </>
  );
};

function AccountCard({ acc }: { acc: AccountInfo }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="">{acc.house.address}</CardTitle>
        <CardDescription>
          <p className="text-sm text-muted-foreground">id: {acc.house.id}</p>
        </CardDescription>
      </CardHeader>
      <CardContent>
        <p>
          Владелец:
          <Link to={routes.users.detail(acc.userId)}>
            <Button variant={"link"}>
              <UUID uuid={acc.userId} />
            </Button>
          </Link>
        </p>
      </CardContent>
    </Card>
  );
}
