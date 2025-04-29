import React from "react";
import { Button } from "./ui/button";
import { Link } from "react-router-dom";
import { routes } from "@/shared/routes";
import { PlusCircle } from "lucide-react";

type Props = {};

export const NewTicketButton = (props: Props) => {
  return (
    <Link to={routes.newTicket} className="">
      <Button className="w-full">
        <PlusCircle />
        Новая заявка
      </Button>
    </Link>
  );
};
