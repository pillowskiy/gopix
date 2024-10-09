'use client';

import { useRef } from 'react';
import { createUserStore, UserContext, type UserProps } from './store';

export default function AuthStoreProvider({
	children,
	...state
}: React.PropsWithChildren<Partial<UserProps>>) {
	const store = useRef(createUserStore(state)).current;

	return <UserContext.Provider value={store}>{children}</UserContext.Provider>;
}
