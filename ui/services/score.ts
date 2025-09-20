"use client";

import { ScoreService as ApiScoreService } from "../api/services/ScoreService";
import { BaseService } from "./base";
import type { Scoreboard, ScoreGraph } from "../api";

export class ScoreService extends BaseService {
  /**
   * Scoreboard
   * Returns the scoreboard
   */
  public static async getScoreboard(
    page: number = 1,
    pageSize: number = 50
  ): Promise<Scoreboard | undefined> {
    return this.handleResponse<Scoreboard>(ApiScoreService.getScore(page, pageSize));
  }

  /**
   * Score graph
   * Returns the score graph
   */
  public static async getScoreGraph(): Promise<ScoreGraph | undefined> {
    return this.handleResponse<ScoreGraph>(ApiScoreService.getScoreGraph());
  }
}
