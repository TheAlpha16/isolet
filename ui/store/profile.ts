import { Submission, Team, User } from "@/api";
import { AuthService, ProfileService } from "@/services";
import { create } from "zustand";

interface ProfileStore {
  meLoading: boolean;
  teamLoading: boolean;
  user: User | null;
  team: Team | null;
  score: number;
  rank: number | null;
  members: User[];
  submissions: Submission[];
  hydrated: boolean;
  setUser: (user: ProfileStore["user"]) => void;
  setTeam: (team: ProfileStore["team"]) => void;
  setHydrated: (hydrated: boolean) => void;
  fetchMe: () => Promise<void>;
  fetchTeam: () => Promise<void>;
  logout: () => Promise<void>;
}

export const useProfileStore = create<ProfileStore>((set) => ({
  meLoading: true, // Start with loading true to prevent hydration mismatch
  teamLoading: false,
  user: null,
  team: null,
  score: 0,
  rank: null,
  members: [],
  submissions: [],
  hydrated: false,

  setUser: (user) => set({ user }),

  setTeam: (team) => set({ team }),

  setHydrated: (hydrated) => set({ hydrated }),

  fetchMe: async () => {
    set({ meLoading: true });
    try {
      const profileMe = await ProfileService.getUserProfile();
      set({ 
        user: profileMe?.user || null, 
        team: profileMe?.team ?? null,
        hydrated: true 
      });
    } catch {
      set({ hydrated: true });
    } finally {
      set({ meLoading: false });
    }
  },

  fetchTeam: async () => {
    set({ teamLoading: true });
    try {
      const profileTeam = await ProfileService.getTeamProfile();
      if (!profileTeam) {
        return;
      }
      set({
        score: profileTeam.score,
        rank: profileTeam.rank ?? null,
        members: profileTeam.members,
        submissions: profileTeam.submissions,
      });
    } catch {
    } finally {
      set({ teamLoading: false });
    }
  },

  logout: async () => {
    await AuthService.logout();
    set({ user: null, team: null, score: 0, rank: null, members: [], submissions: [] });
  },
}));
