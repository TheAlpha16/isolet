"use client";

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { LoadingButton } from "@/components/ui/loading-button";
import { PasswordInput } from "@/components/ui/password-input";
import useResetPassword from "@/hooks/useResetPassword";
import showToast, { ToastStatus } from "@/utils/toastHelper";
import { useSearchParams } from "next/navigation";
import React, { useState } from "react";

export default function ResetPassword() {
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const { loading, resetPassword } = useResetPassword();
  const searchParams = useSearchParams();
  const token = searchParams.get("token");

  async function onSubmit(event: React.FormEvent) {
    event.preventDefault();
    if (loading) return;

    if (!token) {
      showToast(ToastStatus.Failure, "missing token");
      return;
    }

    await resetPassword(password, confirm, token);
  }

  return (
    <div className="container flex flex-col items-center justify-center h-full">
      <Card className="w-[350px]">
        <form onSubmit={onSubmit}>
          <CardHeader className="space-y-1">
            <CardTitle className="text-2xl">Reset Password</CardTitle>
            <CardDescription>Enter your new password</CardDescription>
          </CardHeader>
          <CardContent className="grid gap-4">
            <div className="grid gap-2">
              <Label htmlFor="password">Password</Label>
              <PasswordInput
                id="password"
                name="password"
                placeholder="password"
                autoComplete="new-password"
                onChange={(event) => setPassword(event.target.value)}
                required
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="confirm-password">Confirm</Label>
              <PasswordInput
                id="confirm-password"
                name="confirm-password"
                placeholder="confirm password"
                autoComplete="new-password"
                onChange={(event) => setConfirm(event.target.value)}
                required
              />
            </div>
          </CardContent>
          <CardFooter>
            <LoadingButton type="submit" className="w-full" loading={loading}>
              Reset
            </LoadingButton>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}
