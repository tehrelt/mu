import { api } from "@/app/api";
import { z } from "zod";
import { houseAccountSchema } from "../types/account";

export const getUserAccountsResponse = z.object({
  accounts: z.array(houseAccountSchema),
});

class AccountService {
  async getUserAccounts(): Promise<z.infer<typeof getUserAccountsResponse>> {
    const response = await api.get("/accounts");

    return response.data;
  }
}

export const accountService = new AccountService();
