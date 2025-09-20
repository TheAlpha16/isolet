"use client";

import type { ProfileMe, ProfileTeam } from "../api";
import { ProfileService as ApiProfileService } from "../api/services/ProfileService";
import { BaseService } from "./base";

export class ProfileService extends BaseService {
  /**
   * User profile
   * Returns the user's profile
   */
  public static async getUserProfile(): Promise<ProfileMe | undefined> {
    return this.handleSilentResponse<ProfileMe>(ApiProfileService.getProfileMe());
  }

  /**
   * Team profile
   * Returns the team's profile
   */
  public static async getTeamProfile(): Promise<ProfileTeam | undefined> {
    return this.handleSilentResponse<ProfileTeam>(ApiProfileService.getProfileTeam());
  }
}
