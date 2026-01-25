import { create } from "zustand";
import { Instance } from "@/api";
import { InstanceService } from "@/services/instance";

interface InstanceStore {
  loading: boolean;
  instanceIdMap: Record<number, Instance>;
  challengeInstanceMap: Record<number, number>;

  fetchInstances: () => Promise<void>;
  startInstance: (challenge_id: number) => Promise<void>;
  stopInstance: (challenge_id: number) => Promise<void>;
  extendInstance: (challenge_id: number) => Promise<void>;
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

  startInstance: async (challenge_id: number) => {
    set({ loading: true });
    try {
      const instance = await InstanceService.startInstance(challenge_id);
      if (instance) {
        set((state) => ({
          instanceIdMap: {
            ...state.instanceIdMap,
            [instance.id]: instance,
          },
          challengeInstanceMap: {
            ...state.challengeInstanceMap,
            [challenge_id]: instance.id,
          },
        }));
      }
    } catch {
    } finally {
      set({ loading: false });
    }
  },

  stopInstance: async (challenge_id: number) => {
    const { challengeInstanceMap, instanceIdMap } = get();
    const instance_id = challengeInstanceMap[challenge_id];
    if (!instance_id) {
      return;
    }

    set({ loading: true });
    try {
      await InstanceService.stopInstance(instance_id);

      // Remove instance from both maps
      const newInstanceIdMap = { ...instanceIdMap };
      const newChallengeInstanceMap = { ...challengeInstanceMap };
      delete newInstanceIdMap[instance_id];
      delete newChallengeInstanceMap[challenge_id];

      set({
        instanceIdMap: newInstanceIdMap,
        challengeInstanceMap: newChallengeInstanceMap,
      });
    } catch {
    } finally {
      set({ loading: false });
    }
  },

  extendInstance: async (challenge_id: number) => {
    const { challengeInstanceMap } = get();
    const instance_id = challengeInstanceMap[challenge_id];
    if (!instance_id) {
      return;
    }

    set({ loading: true });
    try {
      const instance = await InstanceService.extendInstance(instance_id);
      if (instance) {
        set((state) => ({
          instanceIdMap: {
            ...state.instanceIdMap,
            [instance.id]: instance,
          },
        }));
      }
    } catch {
    } finally {
      set({ loading: false });
    }
  },
}));
