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
import useLogin from "@/hooks/useLogin";
import { useProfileStore } from "@/store";
import { UI_ROUTES } from "@/utils/routes";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";

export default function Login() {
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const { loading, login } = useLogin();
  const router = useRouter();
  const { fetchMe } = useProfileStore();

  async function onSubmit(event: React.FormEvent) {
    event.preventDefault();
    if (loading) return;

    const result = await login(identifier, password);
    if (result) {
      await fetchMe();
      router.replace(UI_ROUTES.home);
    }
  }

  return (
    <div className="container flex flex-col items-center justify-center h-full">
      <Card className="w-[350px]">
        <form onSubmit={onSubmit}>
          <CardHeader className="space-y-1">
            <CardTitle className="text-2xl">Sign in</CardTitle>
            <CardDescription>Enter your email/username and password</CardDescription>
          </CardHeader>
          <CardContent className="grid gap-4">
            <div className="grid gap-2">
              <Label htmlFor="identifier">Email/Username</Label>
              <Input
                id="identifier"
                type="text"
                placeholder="thealpha16@isolet.dev"
                name="identifier"
                autoComplete="email"
                onChange={(event) => setIdentifier(event.target.value)}
                required
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="password">Password</Label>
              <PasswordInput
                id="password"
                name="password"
                placeholder="password"
                autoComplete="current-password"
                onChange={(event) => setPassword(event.target.value)}
                required
              />
            </div>
            <div className="text-sm text-right">
              <Link href="/forgot-password" className="text-primary hover:underline">
                Forgot password?
              </Link>
            </div>
          </CardContent>
          <CardFooter>
            <LoadingButton type="submit" className="w-full" loading={loading}>
              Sign In
            </LoadingButton>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}
