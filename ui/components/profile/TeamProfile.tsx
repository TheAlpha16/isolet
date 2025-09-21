"use client";

import { Team, User } from "@/api";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Stars, Trophy, UsersRound } from "lucide-react";

interface TeamProfileProps {
  team: Team;
  members?: User[];
  score: number;
  rank: number | null;
}

export function TeamProfile({ team, members, score, rank }: TeamProfileProps) {
  return (
    <Card className="flex flex-col">
      <CardHeader className="flex flex-row items-center gap-4">
        <div className="flex items-center gap-2">
          <Avatar className="w-12 h-12">
            <AvatarFallback className="text-lg">
              {team.name.slice(0, 2).toUpperCase()}
            </AvatarFallback>
          </Avatar>
          <div style={{ marginTop: "0px" }}>
            <CardTitle className="text-2xl">{team.name}</CardTitle>
          </div>
        </div>
      </CardHeader>
      <CardContent className="sm:flex w-full">
        <div className="flex flex-wrap justify-between items-center w-full gap-4">
          <div className="flex items-center gap-2">
            <UsersRound className="h-6 w-6 text-green-500" />
            <div>
              <p className="text-sm font-medium">Members</p>
              <p className="text-3xl font-bold">{members?.length}</p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Stars className="h-6 w-6 text-blue-500" />
            <div>
              <p className="text-sm font-medium">Score</p>
              <p className="text-3xl font-bold">{score}</p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Trophy className="h-6 w-6 text-yellow-500" />
            <div>
              <p className="text-sm font-medium">Rank</p>
              <p className="text-3xl font-bold">{rank ?? "-"}</p>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
