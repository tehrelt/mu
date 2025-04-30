import { api } from "@/app/api";
import { cabinetSchema } from "../types/cabinet";

class CabinetService {
  async find(id: string) {
    const res = await api.get("/cabinets/" + id);
    return cabinetSchema.parse(res.data);
  }
}

export const cabinetService = new CabinetService();
