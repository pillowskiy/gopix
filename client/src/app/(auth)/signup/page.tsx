import { Input } from '@/components/ui/input/input';
import { Button, OAuthButton } from '@/components/ui/button';
import AuthTabs from '../auth-tabs';
import styles from '../auth.module.scss';
import { Description, Field } from '@headlessui/react';
import { Checkbox } from '@/components/ui/checkbox';

export default function SignUpPage() {
	return (
		<div className={styles.loginContainer}>
			<AuthTabs defaultActive={1} />

			<form className={styles.loginForm}>
				<Input placeholder='Username' />
				<Input placeholder='Email' />
				<Input placeholder='Password' type='password' />
				<Input placeholder='Confirm your password' type='password' />
				<Field className={styles.loginFormTerms}>
					<Checkbox size='small' />
					<Description className={styles.loginFormTermsText}>
						I agree to the Terms of Service and Privacy Policy
					</Description>
				</Field>
				<Button className={styles.loginFormSubmit} type='submit'>
					Sign up
				</Button>
			</form>

			<div className={styles.loginOAuth}>
				<OAuthButton service='google' href='/#' />
				<OAuthButton service='apple' href='/#' />
			</div>

			<div className={styles.loginContainerHighlight} />
		</div>
	);
}
