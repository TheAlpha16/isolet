/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { CreateTeamInput } from '../models/CreateTeamInput';
import type { GenerateInviteOutput } from '../models/GenerateInviteOutput';
import type { JoinTeamInput } from '../models/JoinTeamInput';
import type { Response } from '../models/Response';
import type { Session } from '../models/Session';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class TeamService {
    /**
     * Create team
     * Creates a new team
     * @param requestBody
     * @returns any Team created
     * @throws ApiError
     */
    public static postTeamCreate(
        requestBody: CreateTeamInput,
    ): CancelablePromise<(Response & {
        message?: any;
        data?: Session;
    })> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/team/create',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `Bad request`,
                401: `Unauthorized`,
                403: `Forbidden`,
            },
        });
    }
    /**
     * Join team
     * Joins an existing team
     * @param requestBody
     * @returns any Team joined
     * @throws ApiError
     */
    public static postTeamJoin(
        requestBody: JoinTeamInput,
    ): CancelablePromise<(Response & {
        message?: any;
        data?: Session;
    })> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/team/join',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `Bad request`,
                401: `Unauthorized`,
                403: `Forbidden`,
            },
        });
    }
    /**
     * Generate invite
     * Generates a team invite link
     * @returns any Invite generated
     * @throws ApiError
     */
    public static postTeamInvite(): CancelablePromise<(Response & {
        message?: any;
        data?: GenerateInviteOutput;
    })> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/team/invite',
            errors: {
                401: `Unauthorized`,
                403: `Forbidden`,
            },
        });
    }
    /**
     * Accept invite
     * Accepts a team invite
     * @param token
     * @returns any Team joined
     * @throws ApiError
     */
    public static getTeamInvite(
        token: string,
    ): CancelablePromise<(Response & {
        message?: any;
        data?: Session;
    })> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/team/invite',
            query: {
                'token': token,
            },
            errors: {
                400: `Bad request`,
                401: `Unauthorized`,
                403: `Forbidden`,
            },
        });
    }
}
