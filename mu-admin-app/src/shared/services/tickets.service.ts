import { api } from "@/app/api";
import { RateCreate, rateListSchema, rateSchema } from "../types/rate";
import { z } from "zod";
import {
  ticketHeaderSchema,
  TicketStatusEnum,
  TicketTypeEnum,
} from "../types/ticket";

type TicketListRequest = {
  status?: TicketStatusEnum;
  type?: TicketTypeEnum;
};

const listResponseSchema = z.object({
  tickets: z.array(ticketHeaderSchema),
});

class TicketService {
  // async create(data: RateCreate) {
  //   const response = await api.post("/rates", data);
  //   return response.data;
  // }

  // async find(id: string) {
  //   const response = await api.get("/rates/" + id);
  //   const parsed = rateSchema.parse(response.data);
  //   return parsed;
  // }

  async list(params?: TicketListRequest) {
    const response = await api.get("/tickets", {
      params,
    });
    const parsed = listResponseSchema.parse(response.data);
    return parsed;
  }
}

export const ticketService = new TicketService();
