import { useEffect } from "react";
import { startRealtime } from "@/realtime/client";
import { User } from "@/api/models/User";

export function useRealtimeClient(teamId: number | undefined, user: User | null) {
  useEffect(() => {
    if (!user) {
      return;
    }

    const cleanup = startRealtime(teamId);

    return () => {
      if (cleanup) {
        cleanup();
      }
    };
  }, [teamId, user]);
}
