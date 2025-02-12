import { pool } from "../../database";
import { logger } from "../../config/logger";
import { Server } from "socket.io";

export const stopInstance = async (io: Server) => {
    pool.connect((err, client) => {
        if (err) {
            logger.error(`Failed to connect to DB for instance stops: ${err.message}`);
            return;
        }

        if (!client) {
            logger.error("Failed to connect to DB for instance stops: No client returned");
            return;
        }

        client.query("LISTEN notify_instance_stop");

        client.on("notification", (msg) => {
            try {
                if (!msg.payload) {
                    logger.error("Received empty payload from DB");
                    return;
                }
                const payload = JSON.parse(msg.payload);
                io.to(`team-${payload.teamid}`).emit("instanceStop", payload);
            } catch (error) {
                const err = error as Error;
                logger.error(`Error parsing instance stop: ${err.message}`);
            }
        });

        logger.info("Listening for instance stops...");
    });
};
