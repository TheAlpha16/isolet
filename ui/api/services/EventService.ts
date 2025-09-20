/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { EventInfo } from '../models/EventInfo';
import type { Response } from '../models/Response';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class EventService {
    /**
     * Event info
     * Returns information about the event
     * @returns any Event information
     * @throws ApiError
     */
    public static getEventInfo(): CancelablePromise<(Response & {
        data?: EventInfo;
    })> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/event/info',
        });
    }
}
