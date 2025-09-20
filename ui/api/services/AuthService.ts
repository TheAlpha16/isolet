/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ForgotPasswordInput } from '../models/ForgotPasswordInput';
import type { LoginInput } from '../models/LoginInput';
import type { RegisterInput } from '../models/RegisterInput';
import type { ResetPasswordInput } from '../models/ResetPasswordInput';
import type { Response } from '../models/Response';
import type { Session } from '../models/Session';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class AuthService {
    /**
     * Register a new user
     * Creates a new user account
     * @param requestBody
     * @returns any Registration successful
     * @throws ApiError
     */
    public static postAuthRegister(
        requestBody: RegisterInput,
    ): CancelablePromise<(Response & {
        message?: any;
        data?: Session;
    })> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/auth/register',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `Bad request`,
            },
        });
    }
    /**
     * Verify email
     * Verifies the user's email using a token
     * @param token
     * @returns any Email verified
     * @throws ApiError
     */
    public static postAuthVerify(
        token: string,
    ): CancelablePromise<(Response & {
        message?: any;
        data?: any;
    })> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/auth/verify',
            query: {
                'token': token,
            },
            errors: {
                400: `Bad request`,
            },
        });
    }
    /**
     * Login
     * Authenticate a user
     * @param requestBody
     * @returns any Login successful
     * @throws ApiError
     */
    public static postAuthLogin(
        requestBody: LoginInput,
    ): CancelablePromise<(Response & {
        message?: any;
        data?: Session;
    })> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/auth/login',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `Bad request`,
            },
        });
    }
    /**
     * Forgot password
     * Sends a password reset email
     * @param requestBody
     * @returns any Password reset email sent
     * @throws ApiError
     */
    public static postAuthForgotPassword(
        requestBody: ForgotPasswordInput,
    ): CancelablePromise<(Response & {
        message?: any;
        data?: any;
    })> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/auth/forgot-password',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `Bad request`,
            },
        });
    }
    /**
     * Reset password
     * Resets the user's password using a token
     * @param requestBody
     * @returns any Password reset successful
     * @throws ApiError
     */
    public static postAuthResetPassword(
        requestBody: ResetPasswordInput,
    ): CancelablePromise<(Response & {
        message?: any;
        data?: any;
    })> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/auth/reset-password',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `Bad request`,
            },
        });
    }
    /**
     * Logout
     * Logs out the user
     * @returns void
     * @throws ApiError
     */
    public static getAuthLogout(): CancelablePromise<void> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/auth/logout',
            errors: {
                401: `Unauthorized`,
            },
        });
    }
}
