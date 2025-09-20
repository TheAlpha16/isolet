"use client";

import { Button } from "@/components/ui/button";
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
import useLogin from "@/hooks/useLogin";
import { Eye, EyeClosed, Loader2 } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";

export default function Login() {
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const { loading, login } = useLogin();
  const router = useRouter();
  const [showPasswd, setShowPasswd] = useState(false);

  async function onSubmit(event: React.SyntheticEvent) {
    event.preventDefault();
    const result = await login(identifier, password);
    if (result) {
      router.replace("/");
    }
  }

  return (
    <div className="container flex flex-col items-center justify-center h-full">
      <Card className="w-[350px]">
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
              placeholder="titan@titancrew"
              name="identifier"
              autoComplete="email"
              onChange={(event) => {
                setIdentifier(event.target.value);
              }}
              required
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="password">Password</Label>
            <div className="relative">
              <Input
                id="password"
                type={showPasswd ? "text" : "password"}
                placeholder="password"
                name="password"
                autoComplete="current-password"
                onChange={(event) => {
                  setPassword(event.target.value);
                }}
                className="pr-10"
                required
              />
              <Button
                variant={"ghost"}
                size="icon"
                className="absolute inset-y-0 right-0"
                onClick={() => setShowPasswd(!showPasswd)}
              >
                {showPasswd ? <Eye className="h-5 w-5" /> : <EyeClosed className="h-5 w-5" />}
              </Button>
            </div>
          </div>
          <div className="text-sm text-right">
            <Link href="/forgot-password" className="text-primary hover:underline">
              Forgot password?
            </Link>
          </div>
        </CardContent>
        <CardFooter>
          <Button className="w-full" onClick={onSubmit} disabled={loading}>
            {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Sign In
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}
