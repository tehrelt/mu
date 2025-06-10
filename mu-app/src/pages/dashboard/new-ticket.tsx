import { NewTicketForm } from "@/components/forms/new-ticket-form";

export const NewTicketPage = () => {
  return (
    <div className="flex flex-col  min-h-screen ">
      <div className=" p-8 rounded-lg shadow-md">
        <NewTicketForm />
      </div>
    </div>
  );
};
