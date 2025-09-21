/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Submission } from './Submission';
import type { Team } from './Team';
import type { User } from './User';
export type ProfileTeam = (Team & {
    score: number;
    rank?: number | null;
    members: Array<User>;
    submissions: Array<Submission>;
});

