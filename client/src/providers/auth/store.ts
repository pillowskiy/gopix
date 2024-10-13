import { getMe } from "@/shared/actions/users";
import type { User } from "@/types/users";
import { createContext, useContext } from "react";
import { createStore } from "zustand";

export interface UserProps {
  data: User | null;
  dirty: boolean;
}

export interface UserState extends UserProps {
  resolve: () => Promise<void>;
}

export type UserStore = ReturnType<typeof createUserStore>;
export const createUserStore = (init?: Partial<UserProps>) => {
  const DEFAULT_PROPS: UserProps = { data: null, dirty: true };
  return createStore<UserState>()((set) => ({
    ...DEFAULT_PROPS,
    ...init,
    resolve: async () => {
      const user = await getMe().catch(() => null);

      set({ data: user, dirty: false });
    },
  }));
};

export const UserContext = createContext<UserStore | null>(null);

export const useUserStore = () => {
  const store = useContext(UserContext);

  if (!store)
    throw new Error("useUserStore must be used within a UserStoreProvider");

  return store.getState();
};
