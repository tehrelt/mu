import { api } from "@/api";
import { LoginSchema, tokenSchema, TokenSchema } from "@/schemes/login";
import { Profile, profileSchema } from "@/schemes/profile";

class AuthService {
  async login(req: LoginSchema): Promise<TokenSchema> {
    const response = await api
      .post<TokenSchema>("auth/login", {
        json: req,
      })
      .json();

    const pair = tokenSchema.safeParse(response);
    if (!pair.success) {
      console.log("invalid response", response);
      throw new Error("invalid response");
    }

    return pair.data;
  }

  async profile(): Promise<Profile | undefined | null> {
    try {
      const response = await api.get<Profile>("auth/profile").json();

      const profile = profileSchema.safeParse(response);
      if (!profile.success) {
        console.log("invalid response", response);
        throw new Error("invalid response");
      }

      return profile.data;
    } catch (e) {
      // console.log("error fetchin profile", e);
    }
  }

  async refresh(token: string): Promise<TokenSchema> {
    const response = await api
      .put<TokenSchema>("auth/refresh", {
        headers: { Authorization: `Bearer ${token}` },
      })
      .json();

    const pair = tokenSchema.safeParse(response);
    if (!pair.success) {
      console.log("invalid response", response);
      throw new Error("invalid response");
    }

    return pair.data;
  }

  async logout() {
    try {
      await api.post("auth/logout");
    } catch (e) {
      console.log(e);
    }
  }
}

export const authService = new AuthService();
