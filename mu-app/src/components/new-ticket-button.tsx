import { Button } from "./ui/button";
import { Link } from "react-router-dom";
import { routes } from "@/shared/routes";
import { PlusCircle } from "lucide-react";

export const NewTicketButton = () => {
  return (
    <Link to={routes.dashboard.newTicket} className="">
      <Button className="w-full">
        <PlusCircle />
        Новая заявка
      </Button>
    </Link>
  );
};
