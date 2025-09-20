/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type Response = {
    status: Response.status;
    message: string;
    data?: Record<string, any> | null;
};
export namespace Response {
    export enum status {
        SUCCESS = 'success',
        ERROR = 'error',
    }
}

