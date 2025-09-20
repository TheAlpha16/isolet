"use client";

import { ChallengeService as ApiChallengeService } from "../api/services/ChallengeService";
import { BaseService } from "./base";
import type { Challenge, Hint, SubmitFlagInput, SubmitFlagOutput, UnlockHintInput } from "../api";

export class ChallengeService extends BaseService {
  /**
   * List challenges
   * Returns a list of challenges
   */
  public static async getChallenges(): Promise<Challenge[] | undefined> {
    return this.handleResponse<Challenge[]>(ApiChallengeService.getChallenge());
  }

  /**
   * Submit flag
   * Submits a flag for a challenge
   */
  public static async submitFlag(
    requestBody: SubmitFlagInput
  ): Promise<SubmitFlagOutput | undefined> {
    return this.handleResponse<SubmitFlagOutput>(
      ApiChallengeService.postChallengeSubmit(requestBody)
    );
  }

  /**
   * Unlock hint
   * Unlocks a hint for a challenge
   */
  public static async unlockHint(requestBody: UnlockHintInput): Promise<Hint | undefined> {
    return this.handleResponse<Hint>(ApiChallengeService.postChallengeHintUnlock(requestBody));
  }
}
