import { Rate } from "@/shared/types/rate";
import React from "react";

type Props = {
  rate: Rate;
};

const RateViewer = ({ rate }: Props) => {
  return <div>{rate.name}</div>;
};

export default RateViewer;
