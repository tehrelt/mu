import { api } from "@/app/api";
import { cabinetSchema } from "../types/cabinet";

class CabinetService {
  async find(id: string) {
    const res = await api.get("/cabinets/" + id);
    return cabinetSchema.parse(res.data);
  }

  async consume(id: string, value: number) {
    await api.post("/cabinets/" + id + "/consume", { consumed: value });
  }
}

export const cabinetService = new CabinetService();
