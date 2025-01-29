"use client";

import { Button } from "@/components/ui/button";
import { Link as TextLink } from "@/components/ui/link";
import { OAuthButton } from "@/components/oauth";
import AuthTabs from "../auth-tabs";
import styles from "../auth.module.scss";
import { InputWithError } from "@/components/ui/input";
import { login } from "./actions";
import { useCallback, useEffect } from "react";
import { useFormState } from "react-dom";

export default function SignInPage() {
  const [state, action, isPending] = useFormState(login, null);

  useEffect(() => {
    if (!state?.success && !state?.errors && state?.message) {
      alert(state.message);
    }
  }, [state]);

  const fieldErrorState = useCallback((cur: typeof state, key: string) => {
    return cur?.success === false ? cur?.errors?.[key] : "";
  }, []);

  return (
    <div className={styles.authContainer}>
      <AuthTabs defaultActive={0} />

      <form action={action} className={styles.authForm}>
        <InputWithError
          error={fieldErrorState(state, "initials")}
          name="initials"
          placeholder="Username on e-mail"
        />

        <InputWithError
          error={fieldErrorState(state, "password")}
          name="password"
          placeholder="Password"
          type="password"
        />

        <Button
          className={styles.authFormSubmit}
          type="submit"
          disabled={isPending}
        >
          Log In
        </Button>
        <TextLink style={{ textAlign: "center" }} size="small" href="/#">
          Forgot your password?
        </TextLink>
      </form>

      <div className={styles.authContainerOAuth}>
        <OAuthButton
          prefix="Log in with"
          service="google"
          href="/oauth/google"
        />
        <OAuthButton prefix="Log in with" service="apple" href="/oauth/apple" />
      </div>

      <div className={styles.authContainerHighlight} />
    </div>
  );
}
