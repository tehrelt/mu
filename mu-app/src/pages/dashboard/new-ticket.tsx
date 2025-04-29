import { NewTicketForm } from "@/components/forms/new-ticket-form";
import React from "react";

type Props = {};

export const NewTicketPage = (props: Props) => {
  return (
    <div className="flex flex-col  min-h-screen ">
      <div className=" p-8 rounded-lg shadow-md">
        <NewTicketForm />
      </div>
    </div>
  );
};
