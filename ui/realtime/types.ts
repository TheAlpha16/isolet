export interface RealtimeMessage {
  event: string;
  [key: string]: any;
}

export const EntityType = {
  INSTANCE: "instance",
  ENDPOINT: "endpoint",
  NOTIFICATION: "notification",
} as const;

export type EntityTypeName = (typeof EntityType)[keyof typeof EntityType];

export const Action = {
  CREATED: "created",
  UPDATED: "updated",
  DELETED: "deleted",
} as const;

export type ActionType = (typeof Action)[keyof typeof Action];

export type NotificationEntity<T = any> = {
  name: EntityTypeName;
  id: number;
  data: T;
};

export type Notification<T = any> = {
  id: string;
  event: string;
  entity?: NotificationEntity<T>;
  action?: ActionType;
  message?: string;
  severity?: string;
  team_ids?: number[];
  type?: string;
  at?: string;
};
