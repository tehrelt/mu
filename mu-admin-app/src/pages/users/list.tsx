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
} from "@/components/ui/pagination";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useDebounce } from "@/shared/hooks/use-debounce";
import { userService } from "@/shared/services/users.service";
import { useQuery } from "@tanstack/react-query";
import { ChevronLeftIcon, ChevronRightIcon } from "lucide-react";
import React, { useEffect } from "react";
import { useSearchParams } from "react-router-dom";

export const UserListPage = () => {
  const [search, setSearch] = React.useState("");
  const [page, setPage] = React.useState(1);
  const [limit, setLimit] = React.useState(10);
  const debouncedLimit = useDebounce(limit, 500);
  const debouncedSearch = useDebounce(search, 500);
  const [sp, setSp] = useSearchParams();

  const [fieldName, setFieldName] = React.useState("email");

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
      query: debouncedSearch,
    });
  }, [debouncedLimit, page, setSp, debouncedSearch]);

  const data = useQuery({
    queryKey: [
      "users",
      { limit: debouncedLimit, page, [fieldName]: debouncedSearch },
    ],
    queryFn: async () =>
      await userService.list({ page, limit, [fieldName]: debouncedSearch }),
  });

  const onSearchChange = (val: string) => setSearch(val);

  return (
    <div className="space-y-4 relative">
      <div className="flex justify-between items-center sticky top-4 rounded-md shadow-md bg-background py-2 px-2 z-999 outline mx-2">
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbPage>Пользователи</BreadcrumbPage>
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
            <Input
              className="w-[256px] rounded-r-none"
              placeholder={`Поиск ${fieldName != "" && `по ${fieldName}`}`}
              value={search}
              onChange={(e) => onSearchChange(e.target.value)}
            />
            <Select value={fieldName} onValueChange={(v) => setFieldName(v)}>
              <SelectTrigger className="w-[128px] rounded-l-none">
                <SelectValue placeholder={"Поле..."} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="email">Почта</SelectItem>
                <SelectItem value="lastName">Фамилия</SelectItem>
              </SelectContent>
            </Select>
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
        {data.isSuccess && data.data && <UserTable data={data.data.users} />}
      </div>
    </div>
  );
};
