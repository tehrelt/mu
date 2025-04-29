import { routes } from "@/shared/routes";
import {
  localizeTicketType,
  Ticket,
  TicketConnectService,
  TicketHeader,
  TicketNewAccount,
} from "@/shared/types/ticket";
import {
  Breadcrumb,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "../../ui/breadcrumb";
import { TicketStatus } from "@/components/ticket-status";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Link } from "react-router-dom";
import { UUID } from "@/components/ui/uuid";
import { useQuery } from "@tanstack/react-query";
import { userService } from "@/shared/services/users.service";
import { fio } from "@/shared/lib/utils";
import { Button } from "@/components/ui/button";
import { rateService } from "@/shared/services/rates.service";
import { localizeServiceType } from "@/shared/types/rate";

type Props = {
  ticket: Ticket;
};

export const TicketViewer = ({ ticket }: Props) => {
  const createdByQuery = useQuery({
    queryKey: ["ticket", { id: ticket.id }, { createdBy: ticket.createdBy }],
    queryFn: async () => await userService.find(ticket.createdBy),
  });

  return (
    <>
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbLink href={routes.tickets.list}>Заявки</BreadcrumbLink>
          <BreadcrumbSeparator />
          <BreadcrumbPage className="flex gap-x-1">
            <span>{ticket.id}</span>
            <TicketStatus val={ticket.ticketStatus} />
          </BreadcrumbPage>
        </BreadcrumbList>
      </Breadcrumb>

      <div className="">
        <Card className="">
          <CardHeader className="flex justify-between">
            <div>
              <CardTitle>
                <p>{localizeTicketType(ticket.ticketType)}</p>
              </CardTitle>
              <CardDescription>
                <p className="text-muted-foreground">Заявка#{ticket.id}</p>
              </CardDescription>
            </div>
            <div>
              <TicketStatus val={ticket.ticketStatus} />
            </div>
          </CardHeader>
          <CardContent className="">
            <div>
              Создал:
              <Link to={routes.users.detail(ticket.createdBy)}>
                <Button variant={"link"}>
                  {createdByQuery.isLoading && !createdByQuery.data ? (
                    <UUID uuid={ticket.createdBy} />
                  ) : (
                    <span>{fio(createdByQuery.data!)}</span>
                  )}
                </Button>
              </Link>
            </div>
            {renderTicket(ticket)}
          </CardContent>
        </Card>
      </div>
    </>
  );
};

const renderTicket = (ticket: Ticket) => {
  switch (ticket.ticketType) {
    case "TicketTypeAccount":
      return <AccountTicket ticket={ticket as TicketNewAccount} />;
    case "TicketTypeConnectService":
      return <ConnectServiceTicket ticket={ticket as TicketConnectService} />;
    default:
      return null;
  }
};

const AccountTicket = ({ ticket }: { ticket: TicketNewAccount }) => {
  return (
    <div>
      <p>Адрес: {ticket.address}</p>
    </div>
  );
};

const ConnectServiceTicket = ({ ticket }: { ticket: TicketConnectService }) => {
  const query = useQuery({
    queryKey: ["ticket", { id: ticket.id }, { service: ticket.serviceId }],
    queryFn: async () => {
      const rate = await rateService.find(ticket.serviceId);
      return rate;
    },
  });

  return (
    <div>
      <div>
        <p className="text-xl font-medium">Услуга</p>
        {query.isLoading && !query.data ? (
          <Link to={routes.rate.detail(ticket.serviceId)}>
            <Button variant={"link"}>
              <UUID uuid={ticket.serviceId} className="text-muted-foreground" />
            </Button>
          </Link>
        ) : (
          <div className="">
            <p>
              Поставщик:
              <Link to={routes.rate.detail(ticket.serviceId)}>
                <Button variant={"link"}>{query.data?.name}</Button>
              </Link>
            </p>
            <p>Вид услуги: {localizeServiceType(query.data!.serviceType)}</p>
          </div>
        )}
      </div>
    </div>
  );
};
