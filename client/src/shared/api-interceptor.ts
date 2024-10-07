import "server-only";

import ky from "ky";
import { cookies } from "next/headers";

export const apiBaseUrl = `${process.env.API_URL}/${process.env.API_VERSION}`;

export const $api = ky.create({
  prefixUrl: apiBaseUrl,
  mode: "cors",
  credentials: "include",
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
  },
  hooks: {
    beforeRequest: [
      (req) => {
        req.headers.set("Cookie", cookies().toString());
      },
    ],
  },
});
