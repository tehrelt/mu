import { api } from "@/app/api";
import { z } from "zod";
import {
  Ticket,
  ticketConnectService,
  ticketHeaderSchema,
  ticketNewAccount,
  TicketStatusEnum,
  TicketTypeEnum,
} from "../types/ticket";

type TicketListRequest = {
  status: TicketStatusEnum;
  type: TicketTypeEnum;
};

const listResponseSchema = z.object({
  tickets: z.array(ticketHeaderSchema),
});

class TicketService {
  async find(id: string): Promise<Ticket> {
    const response = await api.get("/tickets/" + id);

    const header = ticketHeaderSchema.parse(response.data);
    if (header.ticketType === "TicketTypeAccount") {
      const data = ticketNewAccount.parse(response.data);
      return data;
    } else if (header.ticketType === "TicketTypeConnectService") {
      const data = ticketConnectService.parse(response.data);
      return data;
    }

    throw new Error("invalid ticket type");
  }

  async list(params?: Partial<TicketListRequest>) {
    const response = await api.get("/tickets", {
      params,
    });
    const parsed = listResponseSchema.parse(response.data);
    return parsed;
  }

  async updateStatus(id: string, status: TicketStatusEnum) {
    await api.patch(`/tickets/${id}`, {
      status,
    });
  }
}

export const ticketService = new TicketService();
