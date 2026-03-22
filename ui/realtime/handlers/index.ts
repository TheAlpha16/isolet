import { instanceHandlers } from "@/realtime/handlers/instance";
import { notificationHandlers } from "@/realtime/handlers/notification";

export const handlers: Record<string, (payload: any) => void> = {
  ...instanceHandlers,
  ...notificationHandlers,
};
