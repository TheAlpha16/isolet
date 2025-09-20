"use client";

import { ChallengeSkeleton } from "@/components/skeletons/challenge";
import { useProfileStore } from "@/store";
import { UI_ROUTES } from "@/utils/routes";
import { redirect } from "next/navigation";
import React from "react";

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const { user, team, meLoading } = useProfileStore();

  if (meLoading) {
    return <ChallengeSkeleton />;
  }

  if (!user) {
    return redirect(UI_ROUTES.auth.login);
  }

  if (!team) {
    return redirect(UI_ROUTES.onboarding.team);
  }

  return children;
}
