/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Challenge } from '../models/Challenge';
import type { Hint } from '../models/Hint';
import type { Response } from '../models/Response';
import type { SubmitFlagInput } from '../models/SubmitFlagInput';
import type { SubmitFlagOutput } from '../models/SubmitFlagOutput';
import type { UnlockHintInput } from '../models/UnlockHintInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ChallengeService {
    /**
     * List challenges
     * Returns a list of challenges
     * @returns any List of challenges
     * @throws ApiError
     */
    public static getChallenge(): CancelablePromise<(Response & {
        data?: Array<Challenge>;
    })> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/challenge',
            errors: {
                401: `Unauthorized`,
                403: `Forbidden`,
            },
        });
    }
    /**
     * Submit flag
     * Submits a flag for a challenge
     * @param requestBody
     * @returns any Flag submitted
     * @throws ApiError
     */
    public static postChallengeSubmit(
        requestBody: SubmitFlagInput,
    ): CancelablePromise<(Response & {
        data?: SubmitFlagOutput;
    })> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/challenge/submit',
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
     * Unlock hint
     * Unlocks a hint for a challenge
     * @param requestBody
     * @returns any Hint unlocked
     * @throws ApiError
     */
    public static postChallengeHintUnlock(
        requestBody: UnlockHintInput,
    ): CancelablePromise<(Response & {
        data?: Hint;
    })> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/challenge/hint/unlock',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `Bad request`,
                401: `Unauthorized`,
                403: `Forbidden`,
            },
        });
    }
}
