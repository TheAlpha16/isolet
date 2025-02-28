import { Server } from "socket.io";
import { ListenerConnection } from "../../database/database";

function createListener(io: Server, channel: string, event: string) {
    const listener = new ListenerConnection(io, channel, event);
    return listener.connect();
}

export const startListeners = async (io: Server) => {
    await createListener(io, "notify_instance_stop", "instanceStop");
    await createListener(io, "notify_instance_update", "instanceUpdate");
};
