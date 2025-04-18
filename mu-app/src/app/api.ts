import axios, { AxiosError, AxiosRequestConfig } from "axios";
import { sessionService } from "../shared/services/session.service";

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_ADDRESS,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.request.use((config) => {
  const token = sessionService.get();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (cfg) => cfg,
  (e) => {
    if (e instanceof AxiosError) {
      if (!e.response) {
        throw e;
      }

      const status = e.response.status;
      if (status === 401) {
        const config = e.config as AxiosRequestConfig & { isRetry: boolean };
        if (config.isRetry) {
          sessionService.clear();
          throw e;
        }
      }
    }
    throw e;
  },
);
