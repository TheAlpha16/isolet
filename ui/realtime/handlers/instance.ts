import { useInstanceStore } from "@/store/instance";
import { InstanceNotification } from "@/models/instance";

export const instanceHandlers = {
  "instance.created": (payload: InstanceNotification) => {
    console.log("Instance created:", payload);
    useInstanceStore.getState().handleInstanceCreated(payload);
  },
  "instance.updated": (payload: InstanceNotification) => {
    console.log("Instance updated:", payload);
    useInstanceStore.getState().handleInstanceUpdated(payload);
  },
  "instance.deleted": (payload: InstanceNotification) => {
    console.log("Instance deleted:", payload);
    useInstanceStore.getState().handleInstanceDeleted(payload);
  },
};
