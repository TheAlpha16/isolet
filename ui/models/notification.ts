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
