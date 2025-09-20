"use client";

import { ProfilePageSkeleton } from "@/components/skeletons/profile";
import { useProfileStore } from "@/store";
import { redirect } from "next/navigation";
import React from "react";

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const { user, team, meLoading } = useProfileStore();

  if (meLoading) {
    return <ProfilePageSkeleton />;
  }

  if (!user) {
    return redirect("/login");
  }

  if (!team) {
    return redirect("/teaminit");
  }

  return children;
}
