import TicketTable from "@/components/tables/tickets";
import {
  Breadcrumb,
  BreadcrumbList,
  BreadcrumbPage,
} from "@/components/ui/breadcrumb";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
} from "@/components/ui/pagination";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useDebounce } from "@/shared/hooks/use-debounce";
import { ticketService } from "@/shared/services/tickets.service";
import {
  localizeTicketStatus,
  localizeTicketType,
  ticketStatusEnum,
  TicketStatusEnum,
  ticketTypeEnum,
  TicketTypeEnum,
} from "@/shared/types/ticket";
import { useQuery } from "@tanstack/react-query";
import { ChevronLeftIcon, ChevronRightIcon, X } from "lucide-react";
import React, { useEffect } from "react";
import { useSearchParams } from "react-router-dom";

export const TicketListPage = () => {
  const [page, setPage] = React.useState(1);
  const [limit, setLimit] = React.useState(10);
  const debouncedLimit = useDebounce(limit, 500);
  const [sp, setSp] = useSearchParams();

  const [ticketType, setTicketType] = React.useState<
    TicketTypeEnum | undefined
  >();
  const [ticketStatus, setTicketStatus] = React.useState<
    TicketStatusEnum | undefined
  >();

  useEffect(() => {
    const page = sp.get("page");
    const limit = sp.get("limit");

    if (page) setPage(+page);
    if (limit) setLimit(+limit);
  }, []);

  useEffect(() => {
    setSp({
      page: page.toString(),
      limit: debouncedLimit.toString(),
    });
  }, [debouncedLimit, page, setSp]);

  const data = useQuery({
    queryKey: ["tickets", { ticketType: ticketType, status: ticketStatus }],
    queryFn: async () =>
      await ticketService.list({
        type: ticketType,
        status: ticketStatus,
      }),
  });

  return (
    <div className="space-y-4 relative">
      <div className="flex justify-between items-center sticky top-4 rounded-md shadow-md bg-background py-2 px-2 z-999 outline mx-2">
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbPage>Заявки</BreadcrumbPage>
          </BreadcrumbList>
        </Breadcrumb>

        <div className="flex gap-x-2 items-center ">
          <Input
            className="w-[100px]"
            value={limit}
            type="number"
            min={1}
            onChange={(e) => setLimit(+e.target.value)}
            placeholder="Записей на страницу"
          />
          <div className="flex">
            <Select
              value={ticketType ?? ""}
              onValueChange={(v) => setTicketType(v as TicketTypeEnum)}
            >
              <SelectTrigger className="w-[128px] rounded-r-none">
                <SelectValue placeholder={"Тип тикета"} />
              </SelectTrigger>
              <SelectContent>
                {Object.entries(ticketTypeEnum.enum).map((e) => (
                  <SelectItem value={e[0]}>
                    {localizeTicketType(e[0] as TicketTypeEnum)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Button
              variant={"outline"}
              className="rounded-l-none"
              disabled={!ticketType}
              onClick={() => setTicketType(undefined)}
            >
              <X />
            </Button>
          </div>
          <div className="flex">
            <Select
              value={ticketStatus ?? ""}
              onValueChange={(v) => setTicketStatus(v as TicketStatusEnum)}
            >
              <SelectTrigger className="w-[128px] rounded-r-none">
                <SelectValue placeholder={"Тип тикета"} />
              </SelectTrigger>
              <SelectContent>
                {Object.entries(ticketStatusEnum.enum).map((e) => (
                  <SelectItem value={e[0]}>
                    {localizeTicketStatus(e[0] as TicketStatusEnum)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Button
              variant={"outline"}
              className="rounded-l-none"
              disabled={!ticketStatus}
              onClick={() => setTicketStatus(undefined)}
            >
              <X />
            </Button>
          </div>
          <div>
            <Pagination className="flex justify-end">
              <PaginationContent>
                <PaginationItem>
                  <Button
                    variant={"outline"}
                    disabled={page === 1}
                    onClick={() => setPage((v) => v - 1)}
                  >
                    <ChevronLeftIcon />
                  </Button>
                </PaginationItem>
                <PaginationItem>
                  <PaginationLink>{page}</PaginationLink>
                </PaginationItem>
                <PaginationItem>
                  <Button
                    variant={"outline"}
                    onClick={() => setPage((v) => v + 1)}
                  >
                    <ChevronRightIcon />
                  </Button>
                </PaginationItem>
              </PaginationContent>
            </Pagination>
          </div>
        </div>
      </div>
      <div>
        {data.isError && <p>Error with fetching users: {data.error.message}</p>}
        {data.isLoading && <p>Loading users</p>}
        {data.isSuccess && data.data && (
          <TicketTable data={data.data.tickets} />
        )}
      </div>
    </div>
  );
};
