import { TeamService } from "@/services";
import { useState } from "react";

export default function useTeamInvite() {
  const [inviteLink, setInviteLink] = useState("");
  const [loading, setLoading] = useState(false);

  const generateInvite = async () => {
    setLoading(true);

    try {
      const res = await TeamService.generateInvite();
      if (res) {
        setInviteLink(res.invite_link);
      }
    } catch {
    } finally {
      setLoading(false);
    }
  };

  return { generateInvite, inviteLink, loading };
}
