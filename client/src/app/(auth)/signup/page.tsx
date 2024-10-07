'use client';

import { Field, Label } from '@headlessui/react';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { OAuthButton } from '@/components/oauth';
import AuthTabs from '../auth-tabs';
import styles from '../auth.module.scss';
import { signup } from './actions';
import { useFormState } from 'react-dom';
import { InputWithError } from '@/components/ui/input';
import { useCallback, useEffect } from 'react';

export default function SignUpPage() {
	const [state, action, isPending] = useFormState(signup, null);

	useEffect(() => {
		if (!state?.success && !state?.errors && state?.message) {
			alert(state.message);
		}
	}, [state]);

	const fieldErrorState = useCallback((cur: typeof state, key: string) => {
		return cur?.success === false ? cur?.errors?.[key] : '';
	}, []);

	return (
		<div className={styles.authContainer}>
			<AuthTabs defaultActive={1} />

			<form action={action} className={styles.authForm}>
				<InputWithError
					error={fieldErrorState(state, 'username')}
					name='username'
					placeholder='Username'
				/>

				<InputWithError error={fieldErrorState(state, 'email')} name='email' placeholder='Email' />

				<InputWithError
					error={fieldErrorState(state, 'password')}
					name='password'
					placeholder='Password'
					type='password'
				/>

				<InputWithError
					error={fieldErrorState(state, 'passwordConfirmation')}
					name='passwordConfirmation'
					placeholder='Confirm your password'
					type='password'
				/>

				<Field
					data-error={!state?.success && state?.errors?.termsConditions ? '' : undefined}
					className={styles.authFormTerms}
				>
					<Checkbox name='termsConditions' size='small' />
					<Label className={styles.authFormTermsText}>
						I agree to the Terms of Service and Privacy Policy
					</Label>
				</Field>

				<Button className={styles.authFormSubmit} type='submit' disabled={isPending}>
					{isPending ? 'Signing up...' : 'Sign up'}
				</Button>
			</form>

			<div className={styles.authContainerOAuth}>
				<OAuthButton prefix='Log in with' service='google' href='/oauth/google' />
				<OAuthButton prefix='Log in with' service='apple' href='/oauth/apple' />
			</div>

			<div className={styles.authContainerHighlight} />
		</div>
	);
}
