import { Server } from "socket.io";
import jwt from "jsonwebtoken";
import * as cookie from "cookie";
import { logger } from "../config/logger";
import { secrets } from "../config/general";
import { fetchInstances } from "../database";

export const socketInit = (server: any) => {
    const io = new Server(server, {
        path: "/ws",
    });

    io.use((socket, next) => {
        try {
            const cookies = cookie.parse(socket.handshake.headers.cookie || "");
            const token = cookies.token;

            if (!token) throw new Error("missing token");

            const decoded = jwt.verify(token, secrets.SESSION_SECRET) as jwt.JwtPayload;
            (socket as any).user = decoded;

            if (decoded.teamid) {
                socket.join(`team-${decoded.teamid}`);
            } else {
                throw new Error("invalid/expired token");
            }
            next();
        } catch (error) {
            logger.error(`Authentication failed: ${(error as Error).message}`);
            next(new Error("Authentication error"));
        }
    });

    io.on("connection", async (socket) => {
        io.socketsJoin(`team-${(socket as any).user.teamid}`);

        const instances = await fetchInstances((socket as any).user.teamid);
        socket.emit("instances", instances);
    });

    return io;
};
