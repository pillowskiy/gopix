import type { User } from "@/types/users";
import { createContext, useContext } from "react";
import { createStore } from "zustand";

export interface UserProps {
  data: User | null;
  dirty: boolean;
}

interface UserState extends UserProps {
  setData: (data: User | null) => void;
}

export type UserStore = ReturnType<typeof createUserStore>;
export const createUserStore = (init?: Partial<UserProps>) => {
  const DEFAULT_PROPS: UserProps = { data: null, dirty: true };
  return createStore<UserState>()((set) => ({
    ...DEFAULT_PROPS,
    ...init,
    setData: (data: User | null) => set({ data, dirty: true }),
  }));
};

export const UserContext = createContext<UserStore | null>(null);

export const useUserStore = () => {
  const store = useContext(UserContext);

  if (!store)
    throw new Error("useUserStore must be used within a UserStoreProvider");

  return store.getState();
};
