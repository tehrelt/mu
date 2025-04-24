import RateTable from "@/components/tables/rates";
import {
  Breadcrumb,
  BreadcrumbList,
  BreadcrumbPage,
} from "@/components/ui/breadcrumb";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useDebounce } from "@/shared/hooks/use-debounce";
import { routes } from "@/shared/routes";
import { rateService } from "@/shared/services/rates.service";
import { useQuery } from "@tanstack/react-query";
import debounce from "lodash.debounce";
import React from "react";
import { Link } from "react-router-dom";

export const RateListPage = () => {
  const [search, setSearch] = React.useState("");
  const debouncedSearch = useDebounce(search, 500);

  const data = useQuery({
    queryKey: ["rates", { query: debouncedSearch }],
    queryFn: async () => await rateService.list({ query: debouncedSearch }),
  });

  const onSearchChange = (val: string) => setSearch(val);

  return (
    <>
      <div className="flex justify-between items-center">
        {/* <h1 className="text-2xl font-bold">Тарифы</h1> */}

        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbPage>Тарифы</BreadcrumbPage>
          </BreadcrumbList>
        </Breadcrumb>

        <div className="flex gap-x-2">
          <div className="w-[256px]">
            <Input
              placeholder="Поиск"
              value={search}
              onChange={(e) => onSearchChange(e.target.value)}
            />
          </div>
          <Link to={routes.rate.create}>
            <Button>Создать тариф</Button>
          </Link>
        </div>
      </div>
      {data.isSuccess && data.data && <RateTable data={data.data.rates} />}
    </>
  );
};
