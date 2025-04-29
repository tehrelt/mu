import { useQuery } from "@tanstack/react-query";
import React from "react";
import { useParams } from "react-router-dom";

type Props = {};

export const ProcessPaymentPage = (props: Props) => {
  const params = useParams();

  const payment = useQuery({
    queryKey: ["payment", { id: params["id"] }],
  });

  return <div>Process of {params["id"]}</div>;
};
