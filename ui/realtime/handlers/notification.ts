import { Notification, Severity } from "@/realtime/types";
import showToast, { ToastStatus } from "@/utils/toastHelper";

function getSeverityToastStatus(severity?: string): ToastStatus {
  switch (severity) {
    case Severity.SUCCESS:
      return ToastStatus.Success;
    case Severity.WARNING:
      return ToastStatus.Warning;
    case Severity.INFO:
      return ToastStatus.Info;
    default:
      return ToastStatus.Info;
  }
}

export function handleNotificationEvent(payload: Notification) {
  console.log("Notification received:", payload);

  if (payload.message) {
    const toastStatus = getSeverityToastStatus(payload.severity);
    showToast(toastStatus, payload.message);
  }
}

export const notificationHandlers = {
  "notification.info": handleNotificationEvent,
  "notification.warning": handleNotificationEvent,
  "notification.error": handleNotificationEvent,
  "notification.success": handleNotificationEvent,
};
