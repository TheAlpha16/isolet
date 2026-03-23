import { create } from "zustand";
import { Instance, Endpoint } from "@/api";
import { InstanceService } from "@/services/instance";
import { InstanceNotification, parseInstanceFromNotification } from "@/models/instance";
import { EndpointNotification, parseEndpointFromNotification } from "@/models/endpoint";
import { Action } from "@/realtime/types";

interface InstanceStore {
  loading: boolean;
  instanceIdMap: Record<number, Instance>;
  challengeInstanceMap: Record<number, number>;

  fetchInstances: () => Promise<void>;
  startInstance: (challenge_id: number) => Promise<void>;
  stopInstance: (challenge_id: number) => Promise<void>;
  extendInstance: (challenge_id: number) => Promise<void>;

  handleInstanceEvent: (notification: InstanceNotification) => void;
  handleEndpointEvent: (notification: EndpointNotification) => void;
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

  handleInstanceEvent: (notification: InstanceNotification) => {
    const action = notification.action;

    switch (action) {
      case Action.CREATED: {
        const instance = parseInstanceFromNotification(notification);
        if (!instance) return;

        set((state) => {
          const existingInstance = state.instanceIdMap[instance.id];

          const mergedInstance = existingInstance
            ? {
                ...instance,
                endpoints:
                  existingInstance.endpoints.length > 0
                    ? existingInstance.endpoints
                    : instance.endpoints,
              }
            : instance;

          return {
            instanceIdMap: {
              ...state.instanceIdMap,
              [instance.id]: mergedInstance,
            },
            challengeInstanceMap: {
              ...state.challengeInstanceMap,
              [instance.challenge_id]: instance.id,
            },
          };
        });
        break;
      }

      case Action.UPDATED: {
        const instance = parseInstanceFromNotification(notification);
        if (!instance) return;

        set((state) => {
          const existingInstance = state.instanceIdMap[instance.id];
          if (!existingInstance) return state;

          return {
            instanceIdMap: {
              ...state.instanceIdMap,
              [instance.id]: {
                ...instance,
                endpoints: existingInstance.endpoints,
              },
            },
          };
        });
        break;
      }

      case Action.DELETED: {
        const instanceId = notification.entity?.id;
        if (!instanceId) return;

        set((state) => {
          const instance = state.instanceIdMap[instanceId];
          if (!instance) return state;

          const newInstanceIdMap = { ...state.instanceIdMap };
          const newChallengeInstanceMap = { ...state.challengeInstanceMap };

          delete newInstanceIdMap[instanceId];
          delete newChallengeInstanceMap[instance.challenge_id];

          return {
            instanceIdMap: newInstanceIdMap,
            challengeInstanceMap: newChallengeInstanceMap,
          };
        });
        break;
      }
    }
  },

  handleEndpointEvent: (notification: EndpointNotification) => {
    const endpoint = parseEndpointFromNotification(notification);
    const instanceId = notification.entity?.data?.instance_id;
    const action = notification.action;

    if (!endpoint || !instanceId) {
      return;
    }

    set((state) => {
      let instance = state.instanceIdMap[instanceId];

      // If instance doesn't exist yet (endpoint event arrived first), create a partial instance
      if (!instance) {
        const entityData = notification.entity?.data;
        if (!entityData) {
          return state;
        }

        instance = {
          id: instanceId,
          team_id: entityData.team_id,
          challenge_id: 0,
          endpoints: [],
        };
      }

      const existingEndpointIndex = instance.endpoints.findIndex((e) => e.name === endpoint.name);

      let updatedEndpoints: Endpoint[];

      if (action === Action.DELETED) {
        updatedEndpoints = instance.endpoints.filter((e) => e.name !== endpoint.name);
      } else if (existingEndpointIndex >= 0) {
        updatedEndpoints = [...instance.endpoints];
        updatedEndpoints[existingEndpointIndex] = endpoint;
      } else {
        updatedEndpoints = [...instance.endpoints, endpoint];
      }

      return {
        instanceIdMap: {
          ...state.instanceIdMap,
          [instanceId]: {
            ...instance,
            endpoints: updatedEndpoints,
          },
        },
      };
    });
  },
}));
