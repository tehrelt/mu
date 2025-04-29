import { api } from "@/app/api";
import { z } from "zod";
import { paymentSchema } from "../types/payment";

const createPaymentInput = z.object({
  accountId: z.string().uuid(),
  amount: z.number(),
});

type CreatePaymentInput = z.infer<typeof createPaymentInput>;

const createPaymentOutput = z.object({
  id: z.string().uuid(),
});

class BillingService {
  async create(input: CreatePaymentInput) {
    const res = await api.post("/billing/", input);
    return createPaymentOutput.parse(res.data);
  }

  async find(id: string) {
    const res = await api.get(`/billing/${id}`);
    return paymentSchema.parse(res.data);
  }
}

export const billingService = new BillingService();
