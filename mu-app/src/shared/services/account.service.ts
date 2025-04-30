import { api } from "@/app/api";
import { z } from "zod";
import { houseAccountSchema } from "../types/account";
import { paymentSchema, PaymentStatus } from "../types/payment";
import { rateSchema } from "../types/rate";

export const getUserAccountsResponse = z.object({
  accounts: z.array(houseAccountSchema),
});

const accountPaymentsResponseSchema = z.object({
  payments: z.array(paymentSchema),
});

const accountServicesResponseSchema = z.object({
  services: z.array(rateSchema),
});

class AccountService {
  async getUserAccounts(): Promise<z.infer<typeof getUserAccountsResponse>> {
    const response = await api.get("/accounts");
    return getUserAccountsResponse.parse(response.data);
  }

  async find(id: string) {
    const response = await api.get(`/accounts/${id}`);
    return houseAccountSchema.parse(response.data);
  }

  async payments(
    id: string,
    f?: Partial<{ status: PaymentStatus; limit: number }>
  ) {
    const res = await api.get(`/accounts/${id}/payments`, { params: f });
    const data = accountPaymentsResponseSchema.parse(res.data);
    return data;
  }

  async services(accountId: string) {
    const response = await api.get(`/accounts/${accountId}/services`);
    return accountServicesResponseSchema.parse(response.data);
  }
}

export const accountService = new AccountService();
