"use client";

import { TeamService as ApiTeamService } from "../api/services/TeamService";
import { BaseService } from "./base";
import type { CreateTeamInput, GenerateInviteOutput, JoinTeamInput, Session } from "../api";

export class TeamService extends BaseService {
  /**
   * Create team
   * Creates a new team
   */
  public static async createTeam(requestBody: CreateTeamInput): Promise<Session | undefined> {
    return this.handleResponse<Session>(ApiTeamService.postTeamCreate(requestBody));
  }

  /**
   * Join team
   * Joins an existing team
   */
  public static async joinTeam(requestBody: JoinTeamInput): Promise<Session | undefined> {
    return this.handleResponse<Session>(ApiTeamService.postTeamJoin(requestBody));
  }

  /**
   * Generate invite
   * Generates a team invite link
   */
  public static async generateInvite(): Promise<GenerateInviteOutput | undefined> {
    return this.handleResponse<GenerateInviteOutput>(ApiTeamService.postTeamInvite());
  }

  /**
   * Accept invite
   * Accepts a team invite
   */
  public static async acceptInvite(token: string): Promise<Session | undefined> {
    return this.handleResponse<Session>(ApiTeamService.getTeamInvite(token));
  }
}
