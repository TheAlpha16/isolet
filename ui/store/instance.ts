import { create } from "zustand";
import { Instance } from "@/api";
import { InstanceService } from "@/services/instance";

interface InstanceStore {
  loading: boolean;
  instanceIdMap: Record<number, Instance>;
  challengeInstanceMap: Record<number, number>;

  fetchInstances: () => Promise<void>;
}

export const useInstanceStore = create<InstanceStore>((set, get) => ({
  loading: false,
  instanceIdMap: {},
  challengeInstanceMap: {},

  fetchInstances: async () => {
    set({ loading: true });
    try {
      const instances = await InstanceService.getInstances();
      if (!instances) {
        return;
      }
      const instanceIdMap: Record<number, Instance> = {};
      const challengeInstanceMap: Record<number, number> = {};

      for (const instance of instances) {
        instanceIdMap[instance.id] = instance;
        challengeInstanceMap[instance.challenge_id] = instance.id;
      }

      set({ instanceIdMap, challengeInstanceMap });
    } catch {
    } finally {
      set({ loading: false });
    }
  },
}));
