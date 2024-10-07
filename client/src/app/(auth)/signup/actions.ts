"use server";

import { $api } from "@/shared/api-interceptor";
import type { SafeActionState } from "@/types/common/react";
import type { User } from "@/types/users";
import { drillCookies } from "@/utils/common/next";
import { FetchError } from "@/utils/fetch";
import { prettifyZodErrors } from "@/utils/validation";
import z from "zod";

const signupSchema = z
  .object({
    username: z
      .string({ message: "Name is required" })
      .min(1, "Name should be at least 1 character")
      .max(48, "Name should be at most 48 characters"),
    email: z.string().email("Please enter a valid email address"),
    password: z
      .string({ message: "Password is required" })
      .min(6, "Password should be at least 6 characters")
      .max(48, "Password should be at most 48 characters"),
    passwordConfirmation: z.string({
      message: "Password confirmation is required",
    }),
    termsConditions: z
      .string()
      .refine((value) => value === "on", { path: ["termsConditions"] }),
  })
  .refine((data) => data.password === data.passwordConfirmation, {
    message: "Passwords don't match",
    path: ["passwordConfirmation"],
  });

export async function signup(
  _: SafeActionState<User>,
  formData: FormData,
): Promise<SafeActionState<User>> {
  const parsedInput = signupSchema.safeParse({
    username: formData.get("username"),
    email: formData.get("email"),
    password: formData.get("password"),
    passwordConfirmation: formData.get("passwordConfirmation"),
    termsConditions: formData.get("termsConditions"),
  });

  if (!parsedInput.success) {
    return {
      message: "Invalid input",
      errors: prettifyZodErrors(parsedInput.error),
      success: false,
    };
  }

  const { data } = parsedInput;
  const res = await $api.post("auth/register", {
    json: {
      username: data.username,
      email: data.email,
      password: data.password,
    },
  });

  try {
    const body = await res.json<User>();
    drillCookies(res);
    return { success: true, data: body };
  } catch (err) {
    const fetchErr = await FetchError.fromKy(err);
    return { success: false, message: fetchErr.error };
  }
}
