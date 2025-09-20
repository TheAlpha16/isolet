import { EventService } from "@/api";
import { EventInfoUI, toEventInfoUI } from "@/models/event";
import { create } from "zustand";

interface EventStore {
  loaded: boolean;
  info: EventInfoUI;
  fetchInfo: () => void;
}

export const useEventStore = create<EventStore>((set) => ({
  loaded: false,
  fetching: false,
  info: {
    name: "isolet",
    event_start: new Date(0),
    event_end: new Date(0),
    post_event: false,
    team_length: 0,
  },

  fetchInfo: async () => {
    const res = await EventService.getEventInfo();
    set({ info: toEventInfoUI(res.data!), loaded: true });
  },
}));
