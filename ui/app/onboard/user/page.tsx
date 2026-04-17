"use client";

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { LoadingButton } from "@/components/ui/loading-button";
import { PasswordInput } from "@/components/ui/password-input";
import useUserOnboard from "@/hooks/useUserOnboard";
import { useProfileStore } from "@/store";
import { UI_ROUTES } from "@/utils/routes";
import { useRouter } from "next/navigation";
import { useState } from "react";

export default function Register() {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const { loading, register } = useUserOnboard();
  const { fetchMe } = useProfileStore();
  const router = useRouter();

  async function onSubmit(event: React.FormEvent) {
    event.preventDefault();
    if (loading) return;

    const res = await register(username, email, password, confirm);
    if (res) {
      await fetchMe();
      router.replace(UI_ROUTES.onboard.team);
    }
  }

  return (
    <div className="container flex flex-col items-center justify-center h-full">
      <Card className="w-[350px]">
        <form onSubmit={onSubmit}>
          <CardHeader className="space-y-1">
            <CardTitle className="text-2xl">Register</CardTitle>
            <CardDescription>Sign up for a new account</CardDescription>
          </CardHeader>
          <CardContent className="grid gap-4">
            <div className="grid gap-2">
              <Label htmlFor="username">Username</Label>
              <Input
                id="username"
                type="text"
                placeholder="username"
                autoComplete="username"
                onChange={(event) => setUsername(event.target.value)}
                required
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="thealpha16@isolet.dev"
                autoComplete="email"
                onChange={(event) => setEmail(event.target.value)}
                required
              />
            </div>
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
              Register
            </LoadingButton>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}
