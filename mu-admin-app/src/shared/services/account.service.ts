import { api } from "@/app/api";
import { z } from "zod";
import { houseAccountSchema, paymentSchema } from "../types/account";
import { accountInfoSchema } from "../types/user";
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
  async find(id: string) {
    const response = await api.get(`/accounts/${id}`);
    return accountInfoSchema.parse(response.data);
  }

  async payments(accountId: string) {
    const response = await api.get(`/accounts/${accountId}/payments`);
    return accountPaymentsResponseSchema.parse(response.data);
    // return response.data;
  }

  async services(accountId: string) {
    const response = await api.get(`/accounts/${accountId}/services`);
    return accountServicesResponseSchema.parse(response.data);
  }
}

export const accountService = new AccountService();
