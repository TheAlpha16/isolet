import { Instance } from "@/api";
import { Notification } from "@/realtime/types";

export interface InstanceEventData {
  id: number;
  challenge_id: number;
  team_id: number;
  expires_at: number;
}

export type InstanceNotification = Notification<InstanceEventData>;

export function parseInstanceFromNotification(notification: InstanceNotification): Instance | null {
  if (!notification.entity?.data) {
    return null;
  }

  const data = notification.entity.data;

  return {
    id: data.id,
    challenge_id: data.challenge_id,
    team_id: data.team_id,
    expires_at: data.expires_at,
    endpoints: [], // Endpoints will be populated separately or via updates
  };
}
