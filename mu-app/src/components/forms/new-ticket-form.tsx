import {
  NewTicketInput,
  TicketConnectServiceInput,
  ticketConnectServiceInputSchema,
  TicketNewAccountInput,
  ticketNewAccountInputSchema,
  TicketType,
  ticketTypeEnum,
} from "@/shared/types/ticket";
import { zodResolver } from "@hookform/resolvers/zod";
import React from "react";
import { useForm } from "react-hook-form";
import { Form, FormField } from "../ui/form";
import { Label } from "../ui/label";
import { Input } from "../ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/select";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card";
import { Button } from "../ui/button";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useAccounts } from "@/shared/hooks/use-accounts";
import { ticketService } from "@/shared/services/ticket.service";
import { rateService } from "@/shared/services/rate.service";
import { Balance } from "../ui/balance";
import { useNavigate } from "react-router-dom";
import { routes } from "@/shared/routes";
import { toast } from "sonner";

export const NewTicketForm = () => {
  const [ticketType, setTicketType] = React.useState<TicketType | undefined>();
  const navigate = useNavigate();

  const { mutate } = useMutation({
    mutationKey: ["ticket", "new"],
    mutationFn: async (dto: NewTicketInput) => await ticketService.create(dto),
    onSuccess: (res) => {
      toast.success("Заявка создана", {
        description: `id: ${res.id}`,
      });
    },
    onSettled: () => {
      navigate(routes.dashboard.index);
    },
  });

  return (
    <div className="flex flex-col gap-6 items-center">
      <Card className="w-full min-w-[320px]">
        <CardHeader>
          <CardTitle className="grid gap-y-4">
            <span>Выберите вид заявки:</span>
            <Select
              onValueChange={(v) => setTicketType(v as TicketType)}
              value={ticketType}
            >
              <SelectTrigger className="w-[240px]">
                <SelectValue placeholder="Выберите вид заявки" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value={ticketTypeEnum.Enum.account}>
                  Новый счет
                </SelectItem>
                <SelectItem value={ticketTypeEnum.Enum["connect-service"]}>
                  Подключить услугу
                </SelectItem>
              </SelectContent>
            </Select>
          </CardTitle>
        </CardHeader>
        <CardContent>
          {!!ticketType && (
            <div className="grid gap-y-4">
              <p className="font-semibold leading-none">Заполните форму</p>
              {renderForm(ticketType, {
                submit: mutate,
              })}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
};

const renderForm = (type: TicketType, props: TicketFormProps) => {
  switch (type) {
    case ticketTypeEnum.Values.account:
      return <TicketNewAccountForm {...props} />;
    case ticketTypeEnum.Values["connect-service"]:
      return <TicketConnectServiceForm {...props} />;
  }
};

type TicketFormProps = {
  submit: (payload: NewTicketInput) => void | Promise<void>;
};

const TicketNewAccountForm = ({ submit }: TicketFormProps) => {
  const form = useForm<TicketNewAccountInput>({
    resolver: zodResolver(ticketNewAccountInputSchema),
  });

  const onSubmit = (data: TicketNewAccountInput) => {
    submit({
      type: "account",
      payload: data,
    });
  };

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit, console.warn)}
        className="grid gap-4"
      >
        <div>
          <FormField
            control={form.control}
            name="address"
            render={({ field }) => {
              return (
                <div className="grid gap-2">
                  <Label>Адрес</Label>
                  <div>
                    <Input {...field} />
                  </div>
                </div>
              );
            }}
          />
        </div>
        <div className="">
          <Button>Отправить</Button>
        </div>
      </form>
    </Form>
  );
};

const TicketConnectServiceForm = ({ submit }: TicketFormProps) => {
  const form = useForm<TicketConnectServiceInput>({
    resolver: zodResolver(ticketConnectServiceInputSchema),
  });

  const accountsQuery = useAccounts();
  const servicesQuery = useQuery({
    queryKey: ["services"],
    queryFn: async () => await rateService.list(),
  });

  const onSubmit = (data: TicketConnectServiceInput) => {
    submit({
      type: "connect-service",
      payload: data,
    });
  };

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit, console.warn)}
        className="grid gap-4"
      >
        <div className="grid gap-y-2">
          <FormField
            control={form.control}
            name="accountId"
            render={({ field }) => {
              return (
                <div className="grid gap-2">
                  <Label>Счет</Label>
                  <Select onValueChange={field.onChange}>
                    <SelectTrigger className="w-[380px]">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {accountsQuery.data &&
                        accountsQuery.data.accounts.map((acc) => (
                          <SelectItem value={acc.id}>
                            {acc.house.address}
                          </SelectItem>
                        ))}
                    </SelectContent>
                  </Select>
                </div>
              );
            }}
          />
          <FormField
            control={form.control}
            name="serviceId"
            render={({ field }) => {
              return (
                <div className="grid gap-2">
                  <Label>Услуга</Label>
                  <Select onValueChange={field.onChange}>
                    <SelectTrigger className="w-[380px]">
                      <SelectValue className="" />
                    </SelectTrigger>
                    <SelectContent>
                      {servicesQuery.data &&
                        servicesQuery.data.rates.map((r) => (
                          <SelectItem value={r.id}>
                            {r.name} - <Balance balance={r.rate} />
                          </SelectItem>
                        ))}
                    </SelectContent>
                  </Select>
                </div>
              );
            }}
          />
        </div>
        <div className="">
          <Button>Отправить</Button>
        </div>
      </form>
    </Form>
  );
};
