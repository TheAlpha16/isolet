import { instanceHandlers } from "@/realtime/handlers/instance";
import { endpointHandlers } from "@/realtime/handlers/endpoint";
import { notificationHandlers } from "@/realtime/handlers/notification";

export const handlers: Record<string, (payload: any) => void> = {
  ...instanceHandlers,
  ...endpointHandlers,
  ...notificationHandlers,
};
