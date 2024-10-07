import "server-only";

import { cookies } from "next/headers";

export function drillCookies(res: Response) {
  const setCookieHeader = res.headers.getSetCookie();
  for (const cookie of setCookieHeader) {
    const [kv, ...pairs] = cookie.split(";");
    const [key, value] = kv.split("=");
    const cookieOptions = Object.fromEntries(
      pairs.map((p) => p.trim().split("=")),
    );

    // TODO: add unit tests, better serialize
    cookies().set(key, value, {
      maxAge: cookieOptions["Max-Age"] ? Number(cookieOptions["Max-Age"]) : 0,
      // biome-ignore lint/complexity/useLiteralKeys: <explanation>
      httpOnly: cookieOptions["HttpOnly"] === "true",
      // biome-ignore lint/complexity/useLiteralKeys: <explanation>
      path: cookieOptions["Path"] || "/",
      // biome-ignore lint/complexity/useLiteralKeys: <explanation>
      secure: cookieOptions["Secure"] === "true",
    });
  }
}
