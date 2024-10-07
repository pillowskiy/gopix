"use server";

import { $api } from "@/shared/api-interceptor";
import type { SafeActionState } from "@/types/common/react";
import type { User } from "@/types/users";
import { drillCookies } from "@/utils/common/next";
import { FetchError } from "@/utils/fetch";
import { prettifyZodErrors } from "@/utils/validation";
import z from "zod";

const signupSchema = z.object({
  initials: z
    .string({ message: "Login is required" })
    .min(1, "Login is required")
    .max(48, "Login should be at most 48 characters"),
  password: z
    .string({ message: "Password is required" })
    .min(6, "Password should be at least 6 characters")
    .max(48, "Password should be at most 48 characters"),
});

export async function login(
  _: SafeActionState<User>,
  formData: FormData,
): Promise<SafeActionState<User>> {
  const parsedInput = signupSchema.safeParse({
    initials: formData.get("initials"),
    password: formData.get("password"),
  });

  if (!parsedInput.success) {
    return {
      message: "Invalid input",
      errors: prettifyZodErrors(parsedInput.error),
      success: false,
    };
  }

  const res = await $api.post("auth/login", { json: parsedInput.data });

  try {
    const body = await res.json<User>();
    drillCookies(res);
    return { success: true, data: body };
  } catch (err) {
    const fetchErr = await FetchError.fromKy(err);
    return { success: false, message: fetchErr.error };
  }
}
