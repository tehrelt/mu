"use server";

import "server-only";
import { TokenSchema } from "@/schemes/login";
import { cookies } from "next/headers";

const keys = {
  accessToken: "accessToken",
  refreshToken: "refreshToken",
};

export const createSession = async (tokenPair: TokenSchema) => {
  console.log("creating session", tokenPair);
  const exp = new Date(Date.now() + 60 * 60 * 24 * 7);
  const store = await cookies();

  store.set(keys.accessToken, tokenPair.accessToken, {
    httpOnly: true,
    secure: true,
    sameSite: "lax",
    path: "/",
    expires: exp,
  });

  store.set(keys.refreshToken, tokenPair.refreshToken, {
    httpOnly: true,
    secure: true,
    sameSite: "lax",
    path: "/",
    expires: exp,
  });
};

export const destroySession = async () => {
  const store = await cookies();

  store.delete(keys.accessToken);
  store.delete(keys.refreshToken);
};

export const getSession = async () => {
  const store = await cookies();

  const accessToken = store.get(keys.accessToken)?.value;
  const refreshToken = store.get(keys.refreshToken)?.value;

  return { accessToken, refreshToken };
};
