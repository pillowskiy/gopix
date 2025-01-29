import { getNotificationsStats } from "@/shared/actions/notifications";
import { NotificationStats } from "@/types/notifications";
import { createContext, useContext } from "react";
import { createStore } from "zustand";

export interface NotificationStatsProps {
  stats: NotificationStats | null;
  dirty: boolean;
}

export interface NotificationStatsState extends NotificationStatsProps {
  resolve: () => Promise<void>;
}

export type NotificationStatsStore = ReturnType<
  typeof createNotificationStatsStore
>;
export const createNotificationStatsStore = (
  init?: Partial<NotificationStatsProps>,
) => {
  const DEFAULT_PROPS: NotificationStatsProps = { stats: null, dirty: true };
  return createStore<NotificationStatsState>()((set) => ({
    ...DEFAULT_PROPS,
    ...init,
    resolve: async () => {
      const stats = await getNotificationsStats().catch(() => null);

      set({ stats, dirty: false });
    },
  }));
};

export const NotificationStatsContext =
  createContext<NotificationStatsStore | null>(null);

export const useNotificationStats = () => {
  const store = useContext(NotificationStatsContext);

  if (!store)
    throw new Error(
      "useNotificationStats must be used within a NotificationStatsContext",
    );

  return store.getState();
};
