import { rateService } from "@/shared/services/rates.service";
import { useQuery } from "@tanstack/react-query";
import React from "react";
import { useParams } from "react-router-dom";
import { RateNotFoundPage } from "./not-found";
import RateViewer from "@/components/views/rate";

type Props = {};

const RateDetailsPage = (props: Props) => {
  const { id } = useParams();

  const query = useQuery({
    queryKey: ["rate", { id }],
    queryFn: async () => rateService.find(id!),
  });

  if (query.isLoading) return <>Loading...</>;
  if (query.isError) return <RateNotFoundPage id={id} />;
  if (!query.data) return <p>Query data is null</p>;

  return <RateViewer rate={query.data} />;
};

export default RateDetailsPage;
