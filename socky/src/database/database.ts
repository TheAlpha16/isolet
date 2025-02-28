import { logger } from "../config/logger";
import { db } from "../database";
import { Server } from "socket.io";
import pgPromise from "pg-promise";

class ListenerConnection {
    private connection: pgPromise.IConnected<{}, any> | null;
    private channel: string;
    private event: string;
    private io: Server;

    constructor(io: Server, channel: string, event: string) {
        this.connection = null;
        this.channel = channel;
        this.event = event;
        this.io = io;
    }

    async connect(retries: number = 5, delay: number = 1000) {
        logger.info(`Connecting to database for ${this.channel}...`);
        
        return new Promise((resolve, reject) => {
            setTimeout(() => {
                db.connect({ direct: true, onLost: this.onConnectionLost })
                    .then(obj => {
                        this.connection = obj;
                        resolve(obj);
                        return this.setupListener(obj.client);
                    })
                    .catch((err: any) => {
                        logger.error(`Failed to connect to database: ${err.message}`);
                        if (retries <= 0) {
                            logger.error("Exiting...");
                            process.exit(1);
                        } else {
                            this.connect(retries - 1, delay * 2)
                                .then(resolve)
                                .catch(reject);
                        }
                    });
            }, delay);
        });
    }

    private onConnectionLost(err: any, e: any) {
        logger.error(`Connection lost on ${this.channel}: `, err);
        this.connection = null;
        this.removeListener(e.client);

        this.connect()
            .then(() => {
                logger.info(`Reconnected to database for ${this.channel}`);
            })
            .catch(() => {
                logger.error(`Failed to reconnect to database on ${this.channel}`);
                process.exit(1);
            });
    }

    private setupListener(client: any) {
        logger.info(`Listening to channel ${this.channel}`);
        client.on("notification", this.onNotification.bind(this));

        return this.connection?.none('LISTEN ${channel:name}', { channel: this.channel })
            .catch((err: any) => {
                logger.error(`Failed to listen to channel ${this.channel}: ${err.message}`);
            });
    }

    private removeListener(client: any) {
        logger.info(`Removing listener for ${this.channel}`);
        client.removeListener('notification', this.onNotification);
    }

    private onNotification(msg: any) {
        try {
            if (!msg.payload) {
                logger.error(`Received empty payload from ${this.channel}`);
                return;
            }
            const payload = JSON.parse(msg.payload);
            this.io.to(`team-${payload.teamid}`).emit(this.event, payload);
        } catch (error) {
            logger.error(`Error parsing ${this.channel} payload: ${(error as Error).message}`);
        }
    }
}

export { ListenerConnection };