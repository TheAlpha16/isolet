import { EventInfo } from "@/api";
import { fromUnixSeconds } from "@/utils/helpers";

export interface EventInfoUI extends Omit<EventInfo, "event_start" | "event_end"> {
  event_start: Date;
  event_end: Date;
}

export function toEventInfoUI(eventInfo: EventInfo): EventInfoUI {
  return {
    ...eventInfo,
    event_start: fromUnixSeconds(eventInfo.event_start),
    event_end: fromUnixSeconds(eventInfo.event_end),
  };
}
