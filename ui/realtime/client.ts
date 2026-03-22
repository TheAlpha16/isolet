import { socket } from "@/realtime/socket";
import { dispatch } from "@/realtime/dispatcher";
import type { Channel } from "phoenix";

let globalChannel: Channel | null = null;
let teamChannel: Channel | null = null;

export function startRealtime(teamId?: number) {
  if (!socket) {
    return;
  }

  if (!socket.isConnected()) {
    socket.connect();
  }

  if (!globalChannel) {
    globalChannel = socket.channel("global");
    globalChannel.on("notification", dispatch);
    globalChannel.join().receive("ok", () => console.log("Joined global channel"));
  }

  if (teamId) {
    if (teamChannel && teamChannel.topic !== `team:${teamId}`) {
      teamChannel.leave();
    }

    if (!teamChannel || teamChannel.topic !== `team:${teamId}`) {
      teamChannel = socket.channel(`team:${teamId}`);
      teamChannel.on("notification", dispatch);
      teamChannel.join().receive("ok", () => console.log(`Joined team channel:${teamId}`));
    }
  }

  return () => {
    if (teamChannel) teamChannel.leave();
    if (globalChannel) globalChannel.leave();
    socket.disconnect();

    teamChannel = null;
    globalChannel = null;
  };
}
