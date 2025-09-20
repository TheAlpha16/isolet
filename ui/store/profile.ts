import { Team, User } from "@/api";
import { ProfileService } from "@/services";
import { create } from "zustand";

interface ProfileStore {
  user: User | null;
  team: Team | null;
  setUser: (user: ProfileStore["user"]) => void;
  setTeam: (team: ProfileStore["team"]) => void;
  fetchMe: () => void;
  logout: () => void;
}

export const useProfileStore = create<ProfileStore>((set) => ({
  user: null,
  team: null,

  setUser: (user) => set({ user }),

  setTeam: (team) => set({ team }),

  fetchMe: async () => {
    const profileMe = await ProfileService.getUserProfile();
    set({ user: profileMe?.user, team: profileMe?.team! });
  },

  logout: async () => {
    set({ user: null, team: null });
  },
}));
