import { Input } from '@/components/ui/input/input';
import { Button, OAuthButton } from '@/components/ui/button';
import { Link as TextLink } from '@/components/ui/link';
import AuthTabs from '../auth-tabs';
import styles from '../auth.module.scss';

export default function SignInPage() {
	return (
		<div className={styles.loginContainer}>
			<AuthTabs defaultActive={0} />

			<form className={styles.loginForm}>
				<Input placeholder='Username on e-mail' />
				<Input placeholder='Password' type='password' />
				<Button className={styles.loginFormSubmit} type='submit'>
					Log In
				</Button>
				<TextLink style={{ textAlign: 'center' }} size='small' href='/#'>
					Forgot your password?
				</TextLink>
			</form>

			<div className={styles.loginOAuth}>
				<OAuthButton service='google' href='/#' />
				<OAuthButton service='apple' href='/#' />
			</div>

			<div className={styles.loginContainerHighlight} />
		</div>
	);
}
