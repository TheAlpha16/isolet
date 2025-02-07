import { createServer } from "http";
import { socketInit } from "@/socket";
import { logger } from "@/config/logger";
import { dbConnect } from "@/database";
import { updateInstance } from "@/services/instances/update";
import { stopInstance } from "@/services/instances/stop";

async function startServer() {
    try {
        await dbConnect();
        logger.info("Connected to database!");

        const server = createServer();
        const io = socketInit(server);

        server.listen(8888, () => {
            logger.info("Attaching postgres listeners...");
            updateInstance(io);
            stopInstance(io);
            logger.info("Postgres listeners attached!");
            logger.info("Socket.IO relay server running on port 8888");
        });
    } catch (error) {
        logger.error(`Failed to start server: ${(error as Error).message}`);
        process.exit(1);
    }
}

startServer();
