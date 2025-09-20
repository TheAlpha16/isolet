"use client";

import { useProfileStore } from "@/store";
import { UI_ROUTES } from "@/utils/routes";
import { redirect } from "next/navigation";
import React from "react";

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const { user } = useProfileStore();

  if (user) {
    return redirect(UI_ROUTES.home);
  }

  return children;
}
