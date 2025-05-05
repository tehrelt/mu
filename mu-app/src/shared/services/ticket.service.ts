import { api } from "@/app/api";
import { NewTicketInput, TicketType } from "../types/ticket";

type CreateResponse = {
  id: string;
};

class TicketService {
  async create(input: NewTicketInput) {
    const res = await api.post(newTicketEndpoint(input.type), input.payload);
    return res.data as CreateResponse;
  }
}

const newTicketEndpoint = (t: TicketType) => {
  switch (t) {
    case "account":
      return "/tickets/new-account";
    case "connect-service":
      return "/tickets/connect-service";
    default:
      return "";
  }
};

export const ticketService = new TicketService();
