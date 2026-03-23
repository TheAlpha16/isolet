import { useInstanceStore } from "@/store/instance";
import { InstanceNotification } from "@/models/instance";

export function handleInstanceEvent(payload: InstanceNotification) {
  console.log(`Instance ${payload.action}:`, payload);
  useInstanceStore.getState().handleInstanceEvent(payload);
}

export const instanceHandlers = {
  "instance.created": handleInstanceEvent,
  "instance.updated": handleInstanceEvent,
  "instance.deleted": handleInstanceEvent,
};
