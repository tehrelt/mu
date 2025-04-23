import { rateService } from "@/shared/services/rates.service";
import { useQuery } from "@tanstack/react-query";

export const RateListPage = () => {
  const data = useQuery({
    queryKey: ["rates"],
    queryFn: async () => await rateService.list(),
  });

  return (
    <div>
      RateList
      {data.isSuccess &&
        data.data &&
        data.data.rates.map((r) => (
          <div key={r.id}>
            {r.name} - {r.measureUnit} - {r.rate}
          </div>
        ))}
    </div>
  );
};
