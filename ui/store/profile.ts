import { Team, User } from "@/api";
import { AuthService, ProfileService } from "@/services";
import { create } from "zustand";

interface ProfileStore {
  meLoading: boolean;
  teamLoading: boolean;
  user: User | null;
  team: Team | null;
  setUser: (user: ProfileStore["user"]) => void;
  setTeam: (team: ProfileStore["team"]) => void;
  fetchMe: () => Promise<void>;
  logout: () => Promise<void>;
}

export const useProfileStore = create<ProfileStore>((set) => ({
  meLoading: false,
  teamLoading: false,
  user: null,
  team: null,

  setUser: (user) => set({ user }),

  setTeam: (team) => set({ team }),

  fetchMe: async () => {
    set({ meLoading: true });
    try {
      const profileMe = await ProfileService.getUserProfile();
      set({ user: profileMe.user, team: profileMe.team ?? null });
    } finally {
      set({ meLoading: false });
    }
  },

  logout: async () => {
    await AuthService.logout();
    set({ user: null, team: null });
  },
}));
