import { getMe } from '@/shared/actions/users';
import AuthStoreProvider from './auth-provider';

export default async function AuthClientWrapper({ children }: React.PropsWithChildren) {
	const user = await getMe().catch(() => null);

	return (
		<AuthStoreProvider data={user} dirty>
			{children}
		</AuthStoreProvider>
	);
}
