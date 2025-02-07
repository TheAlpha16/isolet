import dotenv from "dotenv";

dotenv.config();

const secrets = {
    SESSION_SECRET: process.env.SESSION_SECRET || "",
}

export { secrets };