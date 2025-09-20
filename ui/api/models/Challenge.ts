/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Category } from './Category';
import type { Hint } from './Hint';
export type Challenge = {
    id?: number;
    name?: string;
    prompt?: string;
    category?: Category;
    type?: Challenge.type;
    points?: number;
    files?: Array<string>;
    hints?: Array<Hint>;
    author?: string;
    tags?: Array<string>;
    links?: Array<string>;
    max_attempts?: number;
    total_solves?: number;
    solved?: boolean;
    attempt_count?: number;
};
export namespace Challenge {
    export enum type {
        STATIC = 'static',
        DYNAMIC = 'dynamic',
        ON_DEMAND = 'on-demand',
    }
}

