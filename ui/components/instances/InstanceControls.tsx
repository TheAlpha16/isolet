"use client";

import React from "react";
import { Button } from "@/components/ui/button";
import { Loader2, Play, StopCircle, RefreshCw } from "lucide-react";

interface InstanceControlsProps {
  isActive: boolean;
  loading: boolean;
  timeLeft: number;
  onStart: () => void;
  onStop: () => void;
  onExtend: () => void;
}

function formatTime(seconds: number): string {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  return `${minutes}:${remainingSeconds.toString().padStart(2, "0")}`;
}

export function InstanceControls({
  isActive,
  loading,
  timeLeft,
  onStart,
  onStop,
  onExtend,
}: InstanceControlsProps) {
  return (
    <div className="flex space-x-2 items-center">
      <Button onClick={isActive ? onStop : onStart} disabled={loading} variant="outline" size="sm">
        {loading ? (
          <Loader2 className="animate-spin text-yellow-500" />
        ) : isActive ? (
          <StopCircle className="text-red-500" />
        ) : (
          <Play className="text-green-500" />
        )}
      </Button>

      {isActive && (
        <Button size="sm" variant="outline" onClick={onExtend} disabled={loading}>
          <RefreshCw className="h-4 w-4 mr-2" />
          Extend
        </Button>
      )}

      {isActive && (
        <span className="text-sm font-mono border h-9 p-2 rounded-lg">{formatTime(timeLeft)}</span>
      )}
    </div>
  );
}
