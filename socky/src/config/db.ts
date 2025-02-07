import dotenv from "dotenv";

dotenv.config();

const dbCreds = {
    POSTGRES_USER: process.env.POSTGRES_USER || "",
    POSTGRES_HOST: process.env.POSTGRES_HOST || "",
    POSTGRES_DATABASE: process.env.POSTGRES_DATABASE || "",
    POSTGRES_PASSWORD: process.env.POSTGRES_PASSWORD || "",
}

export { dbCreds };