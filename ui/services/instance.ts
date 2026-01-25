"use client";

import type { Instance } from "@/api";
import { InstanceService as ApiInstanceService } from "../api/services/InstanceService";
import { BaseService } from "./base";

export class InstanceService extends BaseService {
  /**
   * List instances
   * Returns list of running instances for the current team
   */
  public static async getInstances(): Promise<Instance[] | undefined> {
    return this.handleResponse<Instance[]>(ApiInstanceService.getInstance());
  }
}
