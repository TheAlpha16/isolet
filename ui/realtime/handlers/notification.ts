export const notificationHandlers = {
  "notification.info": (payload: any) => {
    console.log("Info notification:", payload.message, payload);
    // TODO: show a toast notification
  },
  "notification.warning": (payload: any) => {
    console.warn("Warning notification:", payload.message, payload);
    // TODO: show a toast notification
  },
  "notification.error": (payload: any) => {
    console.error("Error notification:", payload.message, payload);
    // TODO: show a toast notification
  },
};
