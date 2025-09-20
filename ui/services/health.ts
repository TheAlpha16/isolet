"use client";

import { HealthService as ApiHealthService } from "../api/services/HealthService";
import { BaseService } from "./base";

export class HealthService extends BaseService {
  /**
   * Health check endpoint
   * Checks if the API is running
   */
  public static async ping(): Promise<any | undefined> {
    return this.handleResponse<any>(ApiHealthService.getPing());
  }
}
