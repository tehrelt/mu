import { createCookieSessionStorage } from "@remix-run/node";
import { TokenPair } from "~/schemes/login";

const tokensKey = "session";
const accessKey = "accessToken";
const refreshKey = "refreshToken";

const tokensCookie = createCookieSessionStorage({
  cookie: {
    name: tokensKey,
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    sameSite: "lax",
    path: "/",
    maxAge: 60 * 60 * 24 * 7,
    secrets: [process.env.SESSION_SECRET ?? "secret"],
  },
});

class SessionService {
  async set(pair: TokenPair) {
    const session = await tokensCookie.getSession();
    session.set(accessKey, pair.accessToken);
    session.set(refreshKey, pair.refreshToken);
    return await tokensCookie.commitSession(session);
  }

  async retrieve(cookieHeader: string): Promise<TokenPair | undefined> {
    console.log("cookieheader", cookieHeader);
    const session = await tokensCookie.getSession(cookieHeader);
    if (!session.data) return undefined;
    return session.data as TokenPair;
  }

  async clear(): Promise<string> {
    const session = await tokensCookie.getSession();
    return await tokensCookie.destroySession(session);
  }
}

export const sessionService = new SessionService();
