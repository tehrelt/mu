import { refreshTokens } from "@/actions/auth";
import { getSession } from "@/actions/session.server";
import ky, { HTTPError } from "ky";
import "server-only";

export const api = ky.create({
  prefixUrl: process.env.API_ADDRESS,
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 10000,
  hooks: {
    beforeRequest: [
      async (r, options) => {
        console.log("requesting", r.method, r.url);
        const pair = await getSession();
        // console.log("beforeRequest", pair);
        if (pair.accessToken) {
          // console.log("attaching access token");
          r.headers.set("Authorization", `Bearer ${pair.accessToken}`);
        }
        return r;
      },
    ],

    afterResponse: [
      async (request, opts, response) => {
        // console.log("response received", {
        //   url: request.url,
        //   status: response.status,
        //   body: await response.json(),
        // });
        return response;
      },
    ],
  },
});
