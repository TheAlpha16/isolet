import { useInstanceStore } from "@/store/instance";
import { EndpointNotification } from "@/models/endpoint";

export function handleEndpointEvent(payload: EndpointNotification) {
  console.log(`Endpoint ${payload.action}:`, payload);
  useInstanceStore.getState().handleEndpointEvent(payload);
}

export const endpointHandlers = {
  "endpoint.created": handleEndpointEvent,
  "endpoint.updated": handleEndpointEvent,
  "endpoint.deleted": handleEndpointEvent,
};
