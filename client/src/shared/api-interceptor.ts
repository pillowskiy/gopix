import ky from "ky";

export const apiBaseUrl = `${process.env.API_URL}/${process.env.API_VERSION}`;

export const $api = ky.create({
  prefixUrl: apiBaseUrl,
  mode: "cors",
  credentials: "include",
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
  },
});
