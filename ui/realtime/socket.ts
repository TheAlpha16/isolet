import { Socket } from "phoenix";
import { getCookie } from "@/utils/helpers";

export const socket = new Socket("/socket", {
  params: { token: getCookie("realtime_token") },
});
