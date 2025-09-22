"use client";

import { FormSkeleton } from "@/components/skeletons/form";
import { Skeleton } from "@/components/ui/skeleton";
import { useProfileStore } from "@/store";
import { UI_ROUTES } from "@/utils/routes";
import { redirect } from "next/navigation";
import React, { Suspense } from "react";

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const { user, meLoading, hydrated } = useProfileStore();

  if (meLoading || !hydrated) {
    return <FormSkeleton />;
  }

  if (user) {
    return redirect(UI_ROUTES.home);
  }

  return (
    <Suspense
      fallback={
        <div className="w-full h-full flex items-center justify-center">
          <Skeleton className="w-[350px] h-[350px]" />
        </div>
      }
    >
      {children}
    </Suspense>
  );
}
