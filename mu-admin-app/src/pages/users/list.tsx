import { UserTable } from "@/components/tables/users";
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
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { useDebounce } from "@/shared/hooks/use-debounce";
import { userService } from "@/shared/services/users.service";
import { useQuery } from "@tanstack/react-query";
import { ChevronLeftIcon, ChevronRightIcon } from "lucide-react";
import React from "react";

export const UserListPage = () => {
  const [search, setSearch] = React.useState("");
  const [page, setPage] = React.useState(1);
  const [limit, setLimit] = React.useState(10);
  const debouncedLimit = useDebounce(limit, 500);
  const debouncedSearch = useDebounce(search, 500);

  const data = useQuery({
    queryKey: [
      "users",
      { query: debouncedSearch, limit: debouncedLimit, page },
    ],
    queryFn: async () =>
      await userService.list({ query: debouncedSearch, page, limit }),
  });

  const onSearchChange = (val: string) => setSearch(val);

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center sticky">
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbPage>Пользователи</BreadcrumbPage>
          </BreadcrumbList>
        </Breadcrumb>

        <div className="flex gap-x-2">
          <Input
            className="w-[100px]"
            value={limit}
            type="number"
            min={1}
            onChange={(e) => setLimit(+e.target.value)}
            placeholder="Записей на страницу"
          />
          <div className="w-[256px]">
            <Input
              placeholder="Поиск"
              value={search}
              onChange={(e) => onSearchChange(e.target.value)}
            />
          </div>
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
      <div>
        {data.isError && <p>Error with fetching users: {data.error.message}</p>}
        {data.isLoading && <p>Loading users</p>}
        {data.isSuccess && data.data && <UserTable data={data.data.users} />}
      </div>
    </div>
  );
};
