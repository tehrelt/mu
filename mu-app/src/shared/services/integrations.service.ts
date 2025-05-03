import { api } from "@/app/api";
import { z } from "zod";

class IntegrationService {
  async getOtpCode() {
    const response = await api.post("/integrations/link-telegram");

    return z
      .object({
        otp: z.string().length(6),
        userId: z.string().uuid(),
      })
      .parse(response.data);
  }
}
export const integrationService = new IntegrationService();
