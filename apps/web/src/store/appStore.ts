import { create } from 'zustand';

interface AppState {
  unreadNotificationCount: number;
  aiDrawerOpen: boolean;
  aiQuestion: string;
  setUnreadNotificationCount: (count: number) => void;
  openAiDrawer: (question: string) => void;
  closeAiDrawer: () => void;
}

export const useAppStore = create<AppState>((set) => ({
  unreadNotificationCount: 3,
  aiDrawerOpen: false,
  aiQuestion: '',
  setUnreadNotificationCount: (unreadNotificationCount) => set({ unreadNotificationCount }),
  openAiDrawer: (aiQuestion) => set({ aiDrawerOpen: true, aiQuestion }),
  closeAiDrawer: () => set({ aiDrawerOpen: false }),
}));
