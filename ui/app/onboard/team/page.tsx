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
import useTeamOnboard, { ActionType } from "@/hooks/useTeamOnboard";
import { useProfileStore } from "@/store";
import { UI_ROUTES } from "@/utils/routes";
import { useRouter } from "next/navigation";
import { useState } from "react";

export default function TeamInit() {
  const [teamname, setTeamName] = useState("");
  const [password, setPassword] = useState("");
  const { loading, teamOnboard } = useTeamOnboard();
  const router = useRouter();
  const { fetchMe } = useProfileStore();

  async function onSubmit(action: ActionType) {
    const res = await teamOnboard(teamname, password, action);
    if (res) {
      await fetchMe();
      router.replace(UI_ROUTES.home);
    }
  }

  return (
    <div className="container flex flex-col items-center justify-center h-full">
      <Card className="w-[350px]">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl">Team</CardTitle>
          <CardDescription>Please create/join a team to continue</CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4">
          <div className="grid gap-2">
            <Label htmlFor="teamname">Team Name</Label>
            <Input
              id="teamname"
              type="text"
              placeholder="teamname"
              name="teamname"
              autoComplete="off"
              onChange={(event) => setTeamName(event.target.value)}
              required
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="password">Password</Label>
            <PasswordInput
              id="password"
              name="password"
              placeholder="password"
              autoComplete="off"
              onChange={(event) => setPassword(event.target.value)}
              required
            />
          </div>
        </CardContent>
        <CardFooter>
          <div className="flex gap-2 w-full">
            <LoadingButton
              className="w-full"
              loading={loading}
              onClick={(event) => {
                event.preventDefault();
                onSubmit(ActionType.JOIN);
              }}
            >
              Join
            </LoadingButton>
            <LoadingButton
              className="w-full"
              variant="secondary"
              loading={loading}
              onClick={(event) => {
                event.preventDefault();
                onSubmit(ActionType.CREATE);
              }}
            >
              Create
            </LoadingButton>
          </div>
        </CardFooter>
      </Card>
    </div>
  );
}
