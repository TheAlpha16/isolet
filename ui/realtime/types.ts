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

/**
 * Notification severity constants - matches backend NotificationSeverity
 */
export const Severity = {
  INFO: "info",
  WARNING: "warning",
  SUCCESS: "success",
} as const;

export type SeverityType = (typeof Severity)[keyof typeof Severity];

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
  severity?: SeverityType;
  team_ids?: number[];
  type?: string;
  at?: string;
};
