import axios, { AxiosError } from "axios";
import { Cog } from "lucide-react";
import { AuthService } from "~/services/auth";
import { sessionService } from "~/services/session.server";

export const api = axios.create({
  baseURL: process.env.API_ADDR,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

export const apiWithoutInterceptors = axios.create({
  baseURL: process.env.API_ADDR,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.request.use(async (config) => {
  const cookies = String(config.headers.get("cookie"));
  const pair = await sessionService.retrieve(cookies);
  if (pair) {
    if (config.headers) {
      config.headers.Authorization = `Bearer ${pair.accessToken}`;
    }
  }
  return config;
});

api.interceptors.response.use(
  (config) => config,
  async (error: AxiosError) => {
    const code = error.response?.status;
    if (code == 401) {
      if (error.config) {
        const config = error.config as typeof error.config & {
          isRetry: boolean;
        };
        if (!config.isRetry) {
          const cookie = String(error.config.headers.get("cookie"));
          const pair = await sessionService.retrieve(cookie);
          if (pair?.refreshToken) {
            config.isRetry = true;

            const auther = new AuthService(cookie);
            const newPair = await auther.refresh(pair.refreshToken);
            const header = await sessionService.set(newPair);
            config.headers.set("cookie", header);

            return api(config);
          }
        }
      }
    }
  }
);
