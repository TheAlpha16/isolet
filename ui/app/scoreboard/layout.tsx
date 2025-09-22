"use client";

import { ScoreboardSkeleton } from "@/components/skeletons/scoreboard";
import { useProfileStore } from "@/store";
import { UI_ROUTES } from "@/utils/routes";
import { redirect } from "next/navigation";
import React from "react";

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const { user, team, meLoading, hydrated } = useProfileStore();

  if (meLoading || !hydrated) {
    return <ScoreboardSkeleton />;
  }

  if (!user) {
    return redirect(UI_ROUTES.auth.login);
  }

  if (!team) {
    return redirect(UI_ROUTES.onboard.team);
  }

  return children;
}
