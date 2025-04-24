import { DataTable } from "@/components/data-table";
import { Rate } from "@/shared/types/rate";
import React from "react";
import { rateColumns } from "./columns";

type Props = {
  data?: Rate[];
};

const sampleData: Rate[] = [
  {
    id: "1",
    name: "Rate 1",
    measureUnit: "USD",
    rate: 1.0,
  },
  {
    id: "2",
    name: "Rate 2",
    measureUnit: "EUR",
    rate: 0.8,
  },
];

const RateTable = ({ data }: Props) => {
  return (
    <div>
      <DataTable data={data ?? sampleData} columns={rateColumns} />
    </div>
  );
};

export default RateTable;
