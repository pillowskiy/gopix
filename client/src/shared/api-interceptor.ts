import ky from "ky";

export const apiBaseUrl = `${process.env.NEXT_PUBLIC_API_URL}/${process.env.NEXT_PUBLIC_API_VERSION}`;

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
      async (req) => {
        if (typeof window === "undefined") {
          // We can't import next/headers on the client side,
          // so we have to use a hack with dynamic imports
          const { cookies } = await import("next/headers");
          req.headers.set("Cookie", cookies().toString());
        }
      },
    ],
  },
});
