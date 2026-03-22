import { handlers } from "@/realtime/handlers";
import { RealtimeMessage } from "@/realtime/types";

export function dispatch(message: RealtimeMessage) {
  const handler = handlers[message.event];

  if (handler) {
    handler(message);
  } else {
    console.warn("Unhandled event:", message.event, message);
  }
}
