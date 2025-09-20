/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Response } from '../models/Response';
import type { Scoreboard } from '../models/Scoreboard';
import type { ScoreGraph } from '../models/ScoreGraph';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ScoreService {
    /**
     * Scoreboard
     * Returns the scoreboard
     * @param page
     * @param pageSize
     * @returns any Scoreboard
     * @throws ApiError
     */
    public static getScore(
        page: number = 1,
        pageSize: number = 50,
    ): CancelablePromise<(Response & {
        data?: Scoreboard;
    })> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/score',
            query: {
                'page': page,
                'page_size': pageSize,
            },
            errors: {
                401: `Unauthorized`,
                403: `Forbidden`,
            },
        });
    }
    /**
     * Score graph
     * Returns the score graph
     * @returns any Score graph
     * @throws ApiError
     */
    public static getScoreGraph(): CancelablePromise<(Response & {
        data?: ScoreGraph;
    })> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/score/graph',
            errors: {
                401: `Unauthorized`,
                403: `Forbidden`,
            },
        });
    }
}
