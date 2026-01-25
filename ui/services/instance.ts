"use client";

import type { Instance, InstanceStartInput, InstanceStopInput, InstanceExtendInput } from "@/api";
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

  /**
   * Start instance
   * Starts a new instance for a challenge
   */
  public static async startInstance(challenge_id: number): Promise<Instance | undefined> {
    const requestBody: InstanceStartInput = { challenge_id };
    return this.handleResponse<Instance>(ApiInstanceService.postInstanceStart(requestBody));
  }

  /**
   * Stop instance
   * Stops a running instance
   */
  public static async stopInstance(instance_id: number): Promise<void> {
    const requestBody: InstanceStopInput = { instance_id };
    await this.handleResponse(ApiInstanceService.postInstanceStop(requestBody));
  }

  /**
   * Extend instance
   * Extends the lifetime of a running instance
   */
  public static async extendInstance(instance_id: number): Promise<Instance | undefined> {
    const requestBody: InstanceExtendInput = { instance_id };
    return this.handleResponse<Instance>(ApiInstanceService.postInstanceExtend(requestBody));
  }
}
