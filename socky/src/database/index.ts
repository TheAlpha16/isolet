import { logger } from "../config/logger";
import { dbCreds } from "../config/db";
import pgPromise from "pg-promise";

const pgp = pgPromise();
let db: pgPromise.IDatabase<any>;

async function connectToDB(retries: number = 5, delay: number = 1000) {
    try {
        db = pgp({
            user: dbCreds.POSTGRES_USER,
            host: dbCreds.POSTGRES_HOST,
            database: dbCreds.POSTGRES_DATABASE,
            password: dbCreds.POSTGRES_PASSWORD,
            port: 5432,
            max: 10,
            idleTimeoutMillis: 30000
        });

        await db.connect();
        logger.info("Connected to database!");
    } catch (error) {
        logger.error(`Failed to connect to database: ${(error as Error).message}`);

        if (retries <= 0) {
            process.exit(1);
        };

        await new Promise((res) => setTimeout(res, delay));
        return connectToDB(retries - 1, delay * 2);
    }
}

async function fetchInstances(teamid: number) {

    try {
        return await db.any("SELECT flags.chall_id, flags.password, flags.port, flags.hostname, flags.deadline, challenges.deployment FROM flags JOIN challenges ON challenges.chall_id = flags.chall_id WHERE teamid = $1", [teamid]);
    } catch (error) {
        logger.error(`Failed to fetch instances: ${(error as Error).message}`);
        return [];
    }
}

export { db, connectToDB, fetchInstances };
