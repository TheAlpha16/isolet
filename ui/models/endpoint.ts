import { Endpoint } from "@/api";
import { Notification } from "@/realtime/types";

export interface EndpointEventData {
  name: string;
  protocol: string;
  hostname: string;
  port?: number | null;
  ready: boolean;
  instance_id: number;
  team_id?: number;
}

export type EndpointNotification = Notification<EndpointEventData>;

export function parseEndpointFromNotification(notification: EndpointNotification): Endpoint | null {
  if (!notification.entity?.data) {
    return null;
  }

  const data = notification.entity.data;

  return {
    name: data.name,
    protocol: data.protocol,
    hostname: data.hostname,
    port: data.port,
    ready: data.ready,
  };
}
