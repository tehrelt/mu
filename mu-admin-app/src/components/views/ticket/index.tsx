import { routes } from "@/shared/routes";
import {
  localizeTicketType,
  Ticket,
  TicketConnectService,
  TicketNewAccount,
  TicketStatusEnum,
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
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Link } from "react-router-dom";
import { UUID } from "@/components/ui/uuid";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { userService } from "@/shared/services/users.service";
import { fio } from "@/shared/lib/utils";
import { Button } from "@/components/ui/button";
import { rateService } from "@/shared/services/rates.service";
import { localizeServiceType } from "@/shared/types/rate";
import { ticketService } from "@/shared/services/tickets.service";
import { accountService } from "@/shared/services/account.service";

type Props = {
  ticket: Ticket;
};

export const TicketViewer = ({ ticket }: Props) => {
  const createdByQuery = useQuery({
    queryKey: ["ticket", { id: ticket.id }, { createdBy: ticket.createdBy }],
    queryFn: async () => await userService.find(ticket.createdBy),
  });

  const client = useQueryClient();

  const { mutate: patchStatus } = useMutation({
    mutationKey: ["ticket-status-patch", { id: ticket.id }],
    mutationFn: async (status: TicketStatusEnum) =>
      await ticketService.updateStatus(ticket.id, status),
    onSettled: () => {
      client.invalidateQueries({
        queryKey: ["ticket", { id: ticket.id }],
      });
    },
  });

  const approve = () => patchStatus("TicketStatusApproved");
  const reject = () => patchStatus("TicketStatusRejected");

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
          <CardFooter className="grid gap-y-2">
            {ticket.ticketStatus === "TicketStatusPending" ? (
              <>
                <div className="text-lg font-medium">Изменить статус?</div>
                <div className="flex gap-x-2">
                  <Button variant={"success"} onClick={approve}>
                    Одобрить
                  </Button>
                  <Button variant={"destructive"} onClick={reject}>
                    Отклонить
                  </Button>
                </div>
              </>
            ) : (
              <div>
                Заявка закрыта - <TicketStatus val={ticket.ticketStatus} />
              </div>
            )}
          </CardFooter>
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
      const account = await accountService.find(ticket.accountId);
      return { rate, account };
    },
  });

  return (
    <div>
      <div>
        {query.isLoading && !query.data ? (
          <>
            <Link to={routes.rate.detail(ticket.serviceId)}>
              <Button variant={"link"}>
                <UUID
                  uuid={ticket.serviceId}
                  className="text-muted-foreground"
                />
              </Button>
            </Link>
            <Link to={routes.accounts.detail(ticket.accountId)}>
              <Button variant={"link"}>
                <UUID
                  uuid={ticket.accountId}
                  className="text-muted-foreground"
                />
              </Button>
            </Link>
          </>
        ) : (
          <div>
            <div>
              <span className="text-xl font-medium">Счёт</span>
              <Link to={routes.accounts.detail(ticket.accountId)}>
                <Button variant={"link"}>
                  {query.data?.account.id} (адрес:{" "}
                  {query.data?.account.house.address})
                </Button>
              </Link>
            </div>
            <p className="text-xl font-medium">Услуга</p>
            <div className="">
              <p>
                Поставщик:
                <Link to={routes.rate.detail(ticket.serviceId)}>
                  <Button variant={"link"}>{query.data?.rate.name}</Button>
                </Link>
              </p>
              <p>
                Вид услуги: {localizeServiceType(query.data!.rate.serviceType)}
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
