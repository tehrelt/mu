import { api } from "@/app/api";
import { RateCreate, rateListSchema, rateSchema } from "../types/rate";

export type RateListRequest = {
  query?: string;
};

class RateService {
  async create(data: RateCreate) {
    const response = await api.post("/rates", data);
    return response.data;
  }

  async find(id: string) {
    const response = await api.get("/rates/" + id);
    const parsed = rateSchema.parse(response.data);
    return parsed;
  }

  async list(params?: RateListRequest) {
    const response = await api.get("/rates", {
      params,
    });
    const parsed = rateListSchema.parse(response.data);
    return parsed;
  }
}

export const rateService = new RateService();
