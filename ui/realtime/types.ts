export interface RealtimeMessage {
  event: string;
  [key: string]: any;
}

export type NotificationEntity<T = any> = {
  name: string;
  id: number;
  data: T;
};

export type Notification<T = any> = {
  id: string;
  event: string;
  entity?: NotificationEntity<T>;
  action?: string;
  message?: string;
  severity?: string;
  team_ids?: number[];
  type?: string;
  at?: string;
};
