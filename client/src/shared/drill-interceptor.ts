import { cookies } from "next/headers";
import { $api } from "./api-interceptor";

// Use it when u want call api auth required endpoint on server
export const $sapi = $api.extend({
  hooks: {
    beforeRequest: [
      (req) => {
        req.headers.set("Cookie", cookies().toString());
      },
    ],
  },
});
