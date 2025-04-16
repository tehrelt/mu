import { api } from "~/api/api.server";
import {
  LoginRequestSchema,
  loginResponseSchema,
  LoginResponseSchema,
} from "~/schemes/login";
import { Profile, profileScheme } from "~/schemes/profile";

export class AuthService {
  baseUrl: string;
  cookies?: string;

  constructor(cookies?: string) {
    this.baseUrl = "http://localhost:8080/api";
    this.cookies = cookies;
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
}
