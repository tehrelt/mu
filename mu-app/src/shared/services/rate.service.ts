import { api } from "@/app/api";
import { TicketType } from "../types/ticket";
import { rateListSchema } from "../types/rate";

export type RateFilters = {
  type?: string[];
};

class RateService {
  async list(f?: RateFilters) {
    const res = await api.get("/rates", {
      params: f,
    });

    const data = rateListSchema.parse(res.data);

    return data;
  }
}

export const rateService = new RateService();
