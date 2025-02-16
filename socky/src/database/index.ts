import pg from "pg";
import { logger } from "../config/logger";
import { dbCreds } from "../config/db";

const { Pool } = pg;

let pool: pg.Pool;

async function dbConnect(retries: number = 5, delay: number = 100) {
    try {
        if (!pool) {
            pool = new Pool({
                user: dbCreds.POSTGRES_USER,
                host: dbCreds.POSTGRES_HOST,
                database: dbCreds.POSTGRES_DATABASE,
                password: dbCreds.POSTGRES_PASSWORD,
                port: 5432,
                max: 10,
                idleTimeoutMillis: 30000
            });

            pool.on("connect", () => {
                logger.info("New client connected to database");
            });

            pool.on("error", (err: Error) => {
                logger.error(`Database Error: ${err.message}`);
            });
        }

    } catch (error) {
        if (retries <= 0) {
            logger.error(`Failed to connect to database: ${(error as Error).message}`);
            throw error
        };
        await new Promise((res) => setTimeout(res, delay));
        return dbConnect(retries - 1, delay * 2);
    }
}

async function fetchInstances(teamid: number): Promise<any[]> {
    const client = await pool.connect();
    try {
        const
            res = await client.query("SELECT flags.chall_id, flags.password, flags.port, flags.hostname, flags.deadline, challenges.deployment FROM flags JOIN challenges ON challenges.chall_id = flags.chall_id WHERE teamid = $1", [teamid]);
        return res.rows;
    }
    catch (error) {
        logger.error(`Failed to fetch instances: ${(error as Error).message}`);
        return [];
    }
    finally {
        client.release();
    }
}

export { pool, dbConnect, fetchInstances };
