"use client";

import { FormSkeleton } from "@/components/skeletons/form";
import { useProfileStore } from "@/store";
import { UI_ROUTES } from "@/utils/routes";
import { redirect } from "next/navigation";
import React from "react";

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const { user, team, meLoading } = useProfileStore();

  if (meLoading) {
    return <FormSkeleton />;
  }

  if (!user) {
    return redirect(UI_ROUTES.auth.login);
  }

  if (team) {
    return redirect(UI_ROUTES.challenge.home);
  }

  return children;
}
