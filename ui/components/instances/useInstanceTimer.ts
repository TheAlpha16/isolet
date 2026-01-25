import { useState, useEffect } from "react";
import type { Instance } from "@/api";

export function useInstanceTimer(instance: Instance | undefined) {
  const [timeLeft, setTimeLeft] = useState(0);

  useEffect(() => {
    if (instance?.expires_at) {
      const expiresAtMs = instance.expires_at * 1000;
      const active = expiresAtMs > Date.now();
      if (active) {
        setTimeLeft(Math.max(0, Math.floor((expiresAtMs - Date.now()) / 1000)));
      } else {
        setTimeLeft(0);
      }
    } else {
      setTimeLeft(0);
    }

    let interval: NodeJS.Timeout | undefined;
    if (instance?.expires_at) {
      const expiresAtMs = instance.expires_at * 1000;
      if (expiresAtMs > Date.now()) {
        interval = setInterval(() => {
          const timeRemaining = Math.max(0, Math.floor((expiresAtMs - Date.now()) / 1000));
          setTimeLeft(timeRemaining);

          if (timeRemaining === 0 && interval) {
            clearInterval(interval);
          }
        }, 1000);
      }
    }

    return () => {
      if (interval) clearInterval(interval);
    };
  }, [instance]);

  return timeLeft;
}
