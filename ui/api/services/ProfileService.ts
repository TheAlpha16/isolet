/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ProfileMe } from '../models/ProfileMe';
import type { ProfileTeam } from '../models/ProfileTeam';
import type { Response } from '../models/Response';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ProfileService {
    /**
     * User profile
     * Returns the user's profile
     * @returns any User profile
     * @throws ApiError
     */
    public static getProfileMe(): CancelablePromise<(Response & {
        data?: ProfileMe;
    })> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/profile/me',
            errors: {
                401: `Unauthorized`,
            },
        });
    }
    /**
     * Team profile
     * Returns the team's profile
     * @returns any Team profile
     * @throws ApiError
     */
    public static getProfileTeam(): CancelablePromise<(Response & {
        data?: ProfileTeam;
    })> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/profile/team',
            errors: {
                401: `Unauthorized`,
                403: `Forbidden`,
            },
        });
    }
}
