import { TeamService } from "@/services";
import showToast, { ToastStatus } from "@/utils/toastHelper";
import { useState } from "react";

export enum ActionType {
  JOIN = "join",
  CREATE = "create",
}

export default function useTeamOnboard() {
  const [loading, setLoading] = useState(false);

  const teamOnboard = async (teamName: string, password: string, action: ActionType) => {
    // normalize
    teamName = teamName.trim();
    password = password.trim();

    if (!teamName || !password) {
      showToast(ToastStatus.Failure, "team name and password are required");
      return false;
    }

    let apiService;
    switch (action) {
      case ActionType.JOIN:
        apiService = TeamService.joinTeam;
      case ActionType.CREATE:
        apiService = TeamService.createTeam;
        break;
      default:
        showToast(ToastStatus.Failure, "invalid action!");
        return false;
    }

    setLoading(true);
    try {
      await apiService({ team_name: teamName, password });
      return true;
    } catch (err: any) {
      return false;
    } finally {
      setLoading(false);
    }
  };

  return { teamOnboard, loading };
}
