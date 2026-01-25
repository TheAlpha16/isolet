/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Instance } from '../models/Instance';
import type { InstanceExtendInput } from '../models/InstanceExtendInput';
import type { InstanceStartInput } from '../models/InstanceStartInput';
import type { InstanceStopInput } from '../models/InstanceStopInput';
import type { Response } from '../models/Response';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class InstanceService {
    /**
     * List instances
     * Returns a list of instances running for the user's team
     * @returns any List of instances
     * @throws ApiError
     */
    public static getInstance(): CancelablePromise<(Response & {
        data?: Array<Instance>;
    })> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/instance',
            errors: {
                401: `Unauthorized`,
                403: `Forbidden`,
            },
        });
    }
    /**
     * Start instance
     * Starts a new instance for an on-demand challenge
     * @param requestBody
     * @returns any Instance started
     * @throws ApiError
     */
    public static postInstanceStart(
        requestBody: InstanceStartInput,
    ): CancelablePromise<(Response & {
        message?: any;
        data?: Instance;
    })> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/instance/start',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `Bad request`,
                401: `Unauthorized`,
                403: `Forbidden`,
                408: `Request timeout`,
            },
        });
    }
    /**
     * Stop instance
     * Stops a running instance
     * @param requestBody
     * @returns any Instance stopped
     * @throws ApiError
     */
    public static postInstanceStop(
        requestBody: InstanceStopInput,
    ): CancelablePromise<(Response & {
        message?: any;
        data?: any;
    })> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/instance/stop',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `Bad request`,
                401: `Unauthorized`,
                403: `Forbidden`,
                408: `Request timeout`,
            },
        });
    }
    /**
     * Extend instance
     * Extends the lifetime of a running instance
     * @param requestBody
     * @returns any Instance extended
     * @throws ApiError
     */
    public static postInstanceExtend(
        requestBody: InstanceExtendInput,
    ): CancelablePromise<(Response & {
        message?: any;
        data?: Instance;
    })> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/instance/extend',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `Bad request`,
                401: `Unauthorized`,
                403: `Forbidden`,
                408: `Request timeout`,
            },
        });
    }
}
