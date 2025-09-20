/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Submission } from './Submission';
import type { Team } from './Team';
export type ProfileTeam = (Team & {
    score: number;
    rank?: number | null;
    submissions: Array<Submission>;
});

