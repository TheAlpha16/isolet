import { EventService } from "@/services";
import { EventInfoUI, toEventInfoUI } from "@/models/event";
import { create } from "zustand";

const defaultInfo: EventInfoUI = {
  name: "isolet",
  event_start: new Date(0),
  event_end: new Date(0),
  post_event: false,
  team_length: 0,
};

interface EventStore {
  loaded: boolean;
  info: EventInfoUI;
  fetchInfo: () => Promise<void>;
}

export const useEventStore = create<EventStore>((set) => ({
  loaded: false,
  fetching: false,
  info: defaultInfo,
  fetchInfo: async () => {
    const infoRes = await EventService.getEventInfo();
    set({ info: infoRes ? toEventInfoUI(infoRes) : defaultInfo, loaded: true });
  },
}));
