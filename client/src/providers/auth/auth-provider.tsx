'use client';

import { useEffect, useRef } from 'react';
import { UserContext, type UserProps, type UserState, createUserStore } from './store';

export default function AuthStoreProvider({
	children,
	...state
}: React.PropsWithChildren<Partial<UserProps>>) {
	const store = useRef(createUserStore(state)).current;

	// biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
	useEffect(() => {
		const cleanState = (state: UserState) => {
			if (state.dirty) state.resolve();
		};

		cleanState(store.getState());
		return store.subscribe(cleanState);
	}, []);

	return <UserContext.Provider value={store}>{children}</UserContext.Provider>;
}
