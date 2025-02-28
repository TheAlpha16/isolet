import { createServer, IncomingMessage, ServerResponse } from "http";
import { socketInit } from "./socket";
import { logger } from "./config/logger";
import { connectToDB, db } from "./database";
import { startListeners } from "./services/instances/handler";

async function startServer() {
    try {
        await connectToDB();

        const server = createServer((req: IncomingMessage, res: ServerResponse) => {
            if (req.method === "GET" && req.url === "/ping") {
                if (!db) {
                    res.writeHead(500, { "Content-Type": "text/plain" });
                    res.end("Database connection not established");
                    return;
                }

                res.writeHead(200, { "Content-Type": "text/plain" });
                res.end("pong");
            }
        });

        const io = socketInit(server);

        server.listen(8888, () => {
            logger.info("Attaching postgres listeners...");
            startListeners(io);
            logger.info("Socket.IO relay server running on port 8888");
        });
    } catch (error) {
        logger.error(`Failed to start server: ${(error as Error).message}`);
        process.exit(1);
    }
}

startServer();
