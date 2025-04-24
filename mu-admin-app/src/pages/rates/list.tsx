import RateTable from "@/components/tables/rates";
import { Button } from "@/components/ui/button";
import { routes } from "@/shared/routes";
import { rateService } from "@/shared/services/rates.service";
import { useQuery } from "@tanstack/react-query";
import { Link } from "react-router-dom";

export const RateListPage = () => {
  const data = useQuery({
    queryKey: ["rates"],
    queryFn: async () => await rateService.list(),
  });

  return (
    <div className="py-8 px-12 space-y-4">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Тарифы</h1>
        <Link to={routes.rate.create}>
          <Button>Создать тариф</Button>
        </Link>
      </div>
      {data.isSuccess && data.data && <RateTable data={data.data.rates} />}
    </div>
  );
};
