"use client";

import { Challenge } from "@/api";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { useInstanceStore } from "@/store";
import { Check, Flag, Trophy } from "lucide-react";
import { useMemo } from "react";

interface ChallengeCardProps {
  challenge: Challenge;
  onClick: () => void;
}

export function ChallengeCard({ challenge, onClick }: ChallengeCardProps) {
  const { instanceIdMap, challengeInstanceMap } = useInstanceStore();

  const isRunning = useMemo(() => {
    if (challenge.type !== Challenge.type.ON_DEMAND) return false;

    const instanceId = challengeInstanceMap[challenge.id];
    if (!instanceId) return false;

    const instance = instanceIdMap[instanceId];
    return instance?.expires_at ? instance.expires_at * 1000 > Date.now() : false;
  }, [challenge.id, challenge.type, challengeInstanceMap, instanceIdMap]);

  return (
    <Card
      className={`hover:shadow-lg dark:hover:shadow-zinc-900 transition-shadow w-[300px] relative`}
    >
      {isRunning && (
        <span className="absolute -top-1.5 -right-1.5 flex size-3 z-10">
          <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-green-400 opacity-75"></span>
          <span className="relative inline-flex size-3 rounded-full bg-green-500"></span>
        </span>
      )}
      <CardHeader>
        <CardTitle className="flex justify-between items-center">
          <span className="max-w-[75%] overflow-hidden text-ellipsis whitespace-nowrap leading-tight">
            {challenge.name}
          </span>
          <Badge variant={challenge.solved ? "secondary" : "default"} className="ml-2 min-w-fit">
            {challenge.points} pts
          </Badge>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-muted-foreground mb-2 truncate">{challenge.prompt}</p>
        <div className="flex flex-wrap gap-1">
          {challenge.tags.map((tag) => (
            <Badge key={tag} variant="outline">
              {tag}
            </Badge>
          ))}
        </div>
      </CardContent>
      <CardFooter className="flex justify-between items-center">
        <div className="flex items-center text-sm text-muted-foreground">
          <Trophy className="w-4 h-4 mr-1" />
          <span>{challenge.total_solves} solves</span>
        </div>
        <Button
          onClick={onClick}
          size="sm"
          variant={challenge.solved ? "secondary" : "default"}
          className="transition-all active:scale-95"
        >
          {challenge.solved ? (
            <>
              <Check className="w-4 h-4 mr-2 text-green-600 dark:text-green-500" />
              Solved
            </>
          ) : (
            <>
              <Flag className="w-4 h-4 mr-2" />
              Solve
            </>
          )}
        </Button>
      </CardFooter>
    </Card>
  );
}
