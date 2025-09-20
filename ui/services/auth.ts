"use client";

import { AuthService as ApiAuthService } from "../api/services/AuthService";
import { BaseService } from "./base";
import type {
  ForgotPasswordInput,
  LoginInput,
  RegisterInput,
  ResetPasswordInput,
  Session,
} from "../api";

export class AuthService extends BaseService {
  /**
   * Register a new user
   * Creates a new user account
   */
  public static async register(requestBody: RegisterInput): Promise<Session | undefined> {
    return this.handleResponse<Session>(ApiAuthService.postAuthRegister(requestBody));
  }

  /**
   * Verify email
   * Verifies the user's email using a token
   */
  public static async verifyEmail(token: string): Promise<any | undefined> {
    return this.handleResponse<any>(ApiAuthService.postAuthVerify(token));
  }

  /**
   * Login
   * Authenticate a user
   */
  public static async login(requestBody: LoginInput): Promise<Session | undefined> {
    return this.handleResponse<Session>(ApiAuthService.postAuthLogin(requestBody));
  }

  /**
   * Forgot password
   * Sends a password reset email
   */
  public static async forgotPassword(requestBody: ForgotPasswordInput): Promise<any | undefined> {
    return this.handleResponse<any>(ApiAuthService.postAuthForgotPassword(requestBody));
  }

  /**
   * Reset password
   * Resets the user's password using a token
   */
  public static async resetPassword(requestBody: ResetPasswordInput): Promise<any | undefined> {
    return this.handleResponse<any>(ApiAuthService.postAuthResetPassword(requestBody));
  }

  /**
   * Logout
   * Logs out the user
   */
  public static async logout(): Promise<void> {
    return this.handleVoidResponse(ApiAuthService.getAuthLogout());
  }
}
