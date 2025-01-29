"use client";

import { useEffect, useRef } from "react";
import {
  NotificationStatsContext,
  type NotificationStatsProps,
  type NotificationStatsState,
  createNotificationStatsStore,
} from "./store";

export default function NotificationStatsProvider({
  children,
  ...state
}: React.PropsWithChildren<Partial<NotificationStatsProps>>) {
  const store = useRef(createNotificationStatsStore(state)).current;

  // biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
  useEffect(() => {
    const cleanState = (state: NotificationStatsState) => {
      if (state.dirty) state.resolve();
    };

    cleanState(store.getState());
    return store.subscribe(cleanState);
  }, []);

  return (
    <NotificationStatsContext.Provider value={store}>
      {children}
    </NotificationStatsContext.Provider>
  );
}
