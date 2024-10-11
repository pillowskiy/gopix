import { getMe } from "@/shared/users";

export const GET = async () => {
  const user = await getMe();
  return new Response(JSON.stringify(user), { status: 200 });
};
