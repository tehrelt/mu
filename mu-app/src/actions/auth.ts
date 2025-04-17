"use server";

import { LoginSchema, TokenSchema } from "@/schemes/login";
import { createSession, destroySession, getSession } from "./session.server";
import { authService } from "@/service/auth.service";

type LoginOpts = {
  callbackUrl?: string;
};

export async function login(data: LoginSchema, opts?: LoginOpts) {
  "use server";
  try {
    const response = await authService.login(data);
    await createSession(response);
  } catch (e) {
    console.log(e);
    throw new Error("Login failed");
  }
}

export async function logout() {
  "use server";
  await destroySession();
}

export async function refreshTokens() {
  "use server";
  const currentSession = await getSession();

  if (!currentSession.refreshToken) {
    return;
  }

  try {
    const newTokens = await authService.refresh(currentSession.refreshToken);
    await createSession(newTokens);
    return newTokens;
  } catch (_) {
    await destroySession();
  }
}
