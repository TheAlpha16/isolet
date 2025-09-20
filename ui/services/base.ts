"use client";

import { Response } from "../api/models/Response";
import { CancelablePromise } from "../api/core/CancelablePromise";
import showToast, { ToastStatus } from "../utils/toastHelper";

export class BaseService {
  /**
   * Handles API response with toast notifications
   * @param apiCall - The API call that returns a CancelablePromise
   * @returns Promise that resolves to the data or undefined
   */
  protected static async handleResponse<TData>(
    apiCall: CancelablePromise<Response & { data?: TData; message?: any }>
  ): Promise<TData | undefined> {
    try {
      const response = await apiCall;

      if (response.status === Response.status.SUCCESS) {
        // Show success toast if message exists and is not empty
        if (response.message && response.message !== "") {
          showToast(ToastStatus.Success, response.message);
        }

        // Return data if it exists
        return response.data;
      } else if (response.status === Response.status.ERROR) {
        // Show error toast with the error message
        showToast(ToastStatus.Failure, response.message);
        return undefined;
      }
    } catch (error: any) {
      // Handle API errors (like 400, 401, etc.)
      if (error?.body?.message) {
        showToast(ToastStatus.Failure, error.body.message);
      } else if (error?.message) {
        showToast(ToastStatus.Failure, error.message);
      } else {
        showToast(ToastStatus.Failure, "An unexpected error occurred");
      }
      return undefined;
    }
  }

  /**
   * Handles API calls that don't return a Response object (like logout)
   * @param apiCall - The API call that returns a CancelablePromise
   * @returns Promise that resolves to void
   */
  protected static async handleVoidResponse(apiCall: CancelablePromise<void>): Promise<void> {
    try {
      await apiCall;
    } catch (error: any) {
      // Handle API errors (like 400, 401, etc.)
      if (error?.body?.message) {
        showToast(ToastStatus.Failure, error.body.message);
      } else if (error?.message) {
        showToast(ToastStatus.Failure, error.message);
      } else {
        showToast(ToastStatus.Failure, "An unexpected error occurred");
      }
    }
  }
}
