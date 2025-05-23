import { api } from "@/app/api";
import { userAccountsSchema, userListSchema, userSchema } from "../types/user";

type UserListRequest = {
  query?: string;
  page?: number;
  limit?: number;
};

class UserService {
  async find(id: string) {
    const response = await api.get("/users/" + id);
    const parsed = userSchema.parse(response.data);
    return parsed;
  }

  async list(params?: UserListRequest) {
    const response = await api.get("/users", {
      params,
    });
    const parsed = userListSchema.parse(response.data);
    return parsed;
  }

  async accountsOfUser(id: string) {
    const response = await api.get(`/users/${id}/accounts`);
    const parsed = userAccountsSchema.parse(response.data);
    return parsed;
  }
}

export const userService = new UserService();
