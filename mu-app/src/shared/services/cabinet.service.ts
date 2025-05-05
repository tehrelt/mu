import { api } from "@/app/api";
import { cabinetLog, cabinetSchema } from "../types/cabinet";
import { z } from "zod";

class CabinetService {
  async find(id: string) {
    const res = await api.get("/cabinets/" + id);
    return cabinetSchema.parse(res.data);
  }

  async consume(id: string, value: number) {
    await api.post("/cabinets/" + id + "/consume", { consumed: value });
  }

  async logs(id: string, f?: Partial<{ limit: number }>) {
    const res = await api.get(`/cabinets/${id}/logs`, { params: f });
    return z
      .object({
        logs: z.array(cabinetLog),
        total: z.number(),
      })
      .parse(res.data);
  }
}

export const cabinetService = new CabinetService();
