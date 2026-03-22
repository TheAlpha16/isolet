export interface RealtimeMessage {
  event: string;
  [key: string]: any;
}

export type Notification = {
  id: string;
  event: string;
  entity?: {
    name: string;
    id: number;
    data: any;
  };
  action?: string;
  message?: string;
  severity?: string;
};
