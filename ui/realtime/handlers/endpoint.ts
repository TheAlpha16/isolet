import { useInstanceStore } from "@/store/instance";
import { EndpointNotification } from "@/models/endpoint";

export const endpointHandlers = {
  "endpoint.updated": (payload: EndpointNotification) => {
    console.log("Endpoint updated:", payload);
    useInstanceStore.getState().handleEndpointUpdated(payload);
  },
  "endpoint.ready": (payload: EndpointNotification) => {
    console.log("Endpoint ready:", payload);
    useInstanceStore.getState().handleEndpointUpdated(payload);
  },
};
