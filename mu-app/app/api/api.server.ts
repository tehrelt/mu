import axios from "axios";
import { Cog } from "lucide-react";
import { sessionService } from "~/services/session.server";

export const api = axios.create({
  baseURL: process.env.API_ADDR,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.request.use(async (config) => {
  try {
    const cookies = config.headers.get("cookie");
    const pair = await sessionService.retrieve(cookies);
    console.log("pair of tokens", pair);
    if (pair) {
      if (config.headers) {
        config.headers.Authorization = `Bearer ${pair.accessToken}`;
      }
    }
  } catch (e) {
    console.error(e);
    throw e;
  }
  return config;
});

api.interceptors.response.use(
  (config) => config,
  async (error) => {
    // throw error;
  }
);
