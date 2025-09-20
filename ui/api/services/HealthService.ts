/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Response } from '../models/Response';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class HealthService {
    /**
     * Health check endpoint
     * Checks if the API is running
     * @returns any API is running
     * @throws ApiError
     */
    public static getPing(): CancelablePromise<(Response & {
        data?: any;
        message?: any;
    })> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/ping',
        });
    }
}
