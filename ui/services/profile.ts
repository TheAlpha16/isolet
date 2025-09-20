"use client";

import { ProfileService as ApiProfileService } from "../api/services/ProfileService";
import { BaseService } from "./base";
import type { ProfileMe, ProfileTeam } from "../api";

export class ProfileService extends BaseService {
  /**
   * User profile
   * Returns the user's profile
   */
  public static async getUserProfile(): Promise<ProfileMe | undefined> {
    return this.handleResponse<ProfileMe>(ApiProfileService.getProfileMe());
  }

  /**
   * Team profile
   * Returns the team's profile
   */
  public static async getTeamProfile(): Promise<ProfileTeam | undefined> {
    return this.handleResponse<ProfileTeam>(ApiProfileService.getProfileTeam());
  }
}
