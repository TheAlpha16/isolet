"use client";

import { ProfileService as ApiProfileService } from "../api/services/ProfileService";
import { BaseService } from "./base";
import type { ProfileMe, ProfileTeam } from "../api";

export class ProfileService extends BaseService {
  /**
   * User profile
   * Returns the user's profile
   */
  public static async getUserProfile(): Promise<ProfileMe> {
    const profileMe = await this.handleResponse<ProfileMe>(ApiProfileService.getProfileMe());
    if (!profileMe) {
      throw new Error("failed to fetch user profile");
    }
    return profileMe;
  }

  /**
   * Team profile
   * Returns the team's profile
   */
  public static async getTeamProfile(): Promise<ProfileTeam | undefined> {
    const profileTeam = await this.handleResponse<ProfileTeam>(ApiProfileService.getProfileTeam());
    if (!profileTeam) {
      throw new Error("failed to fetch team profile");
    }
    return profileTeam;
  }
}
