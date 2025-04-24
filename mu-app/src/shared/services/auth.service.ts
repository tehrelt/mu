import { api } from "../../app/api";
import {
  LoginInput,
  AccessToken,
  tokenPairSchema,
  profileSchema,
  Profile,
} from "../types/auth";
import { sessionService } from "./session.service";

class AuthService {
  async login(data: LoginInput): Promise<AccessToken> {
    const response = await api.post("/auth/login", data);

    const token = tokenPairSchema.safeParse(response.data);
    if (!token.success) {
      throw new Error("Invalid response");
    }

    sessionService.set(token.data.accessToken);

    return token.data;
  }

  async logout() {
    await api.post("/auth/logout");
  }

  async profile(): Promise<Profile> {
    const response = await api.get("/auth/profile");

    const profile = profileSchema.safeParse(response.data);
    if (!profile.success) {
      throw new Error("Invalid response");
    }

    return profile.data;
  }

  async refresh() {
    const response = await api.put("/auth/refresh");

    const token = tokenPairSchema.safeParse(response.data);
    if (!token.success) {
      throw new Error("Invalid response");
    }

    sessionService.set(token.data.accessToken);

    return token.data;
  }
}

export const authService = new AuthService();
