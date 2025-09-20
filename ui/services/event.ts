"use client";

import { EventService as ApiEventService } from "../api/services/EventService";
import { BaseService } from "./base";
import type { EventInfo } from "../api";

export class EventService extends BaseService {
  /**
   * Event info
   * Returns information about the event
   */
  public static async getEventInfo(): Promise<EventInfo | undefined> {
    return this.handleResponse<EventInfo>(ApiEventService.getEventInfo());
  }
}
