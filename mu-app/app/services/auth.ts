import { api, apiWithoutInterceptors } from "~/api/api.server";
import {
  LoginRequestSchema,
  loginResponseSchema,
  LoginResponseSchema,
  TokenPair,
} from "~/schemes/login";
import { Profile, profileScheme } from "~/schemes/profile";
import { sessionService } from "./session.server";

export class AuthService {
  baseUrl: string;
  cookies?: string;

  constructor(cookies?: string | null) {
    this.baseUrl = "http://localhost:8080/api";
    if (cookies) this.cookies = cookies;
  }

  async refresh(token: string): Promise<TokenPair> {
    const response = await apiWithoutInterceptors.put(
      "/auth/refresh",
      {},
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );

    const newPair = loginResponseSchema.safeParse(response.data);
    if (!newPair.success) {
      throw new Error("Invalid response");
    }

    return newPair.data;
  }

  async login(req: LoginRequestSchema): Promise<LoginResponseSchema> {
    const response = await api.post("/auth/login", req);

    const pair = loginResponseSchema.safeParse(response.data);

    if (!pair.success) {
      throw new Error("Invalid response");
    }

    return pair.data;
  }

  async profile(): Promise<Profile> {
    const response = await api.get("/auth/profile", {
      headers: {
        Cookie: this.cookies,
      },
    });

    if (!response) {
      throw new Error("Failed to fetch profile");
    }

    console.log("profile data", response.data);

    const profile = profileScheme.safeParse(response.data);

    if (!profile.success) {
      throw new Error("Invalid profile data");
    }

    return profile.data;
  }

  async logout(): Promise<void> {
    // const response = await api.post("/auth/logout", {
    //   headers: {
    //     Cookie: this.cookies,
    //   },
    // });
    // if (!response) {
    //   throw new Error("Failed to logout");
    // }

    if (this.cookies) sessionService.clear(this.cookies);
    else console.warn("AuthService.logout()", "where cookie?");
  }
}
