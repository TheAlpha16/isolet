import { create } from "zustand";

interface MetadataStore {
    metadataLoaded: boolean;
    ctfName: string;
    eventStart: string;
    eventEnd: string;
    postEvent: boolean;
    teamLen: number;
    fetchMetadata: () => void;
}

export const useMetadataStore = create<MetadataStore>((set) => ({
    metadataLoaded: false,
    ctfName: "isolet",
    eventStart: "",
    eventEnd: "",
    postEvent: false,
    teamLen: 0,

    fetchMetadata: async () => {
        const res = await fetch("/auth/metadata", {
            method: "GET",
        });

        if (res.ok) {
            const data = await res.json();
            set({
                ctfName: data.ctf_name,
                eventStart: data.event_start,
                eventEnd: data.event_end,
                postEvent: data.post_event,
                teamLen: data.team_len,
            });
        }
    },
}));
