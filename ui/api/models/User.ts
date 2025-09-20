/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type User = {
    id: number;
    username: string;
    email: string;
    role: User.role;
};
export namespace User {
    export enum role {
        ADMIN = 'admin',
        AUTHOR = 'author',
        CAPTAIN = 'captain',
        PLAYER = 'player',
    }
}

