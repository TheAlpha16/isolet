import React, { useEffect, useState, useMemo } from "react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Loader2, Play, StopCircle, RefreshCw, Terminal } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { useInstanceStore } from "@/store/instance";
import { CopyButton } from "@/components/utils/copy-button";
import { GenerateChallengeEndpoint } from "@/utils/parser";
import showToast, { ToastStatus } from "@/utils/toastHelper";

interface InstanceCardProps {
  challenge_id: number;
}

export function InstanceCard({ challenge_id }: InstanceCardProps) {
  const [timeLeft, setTimeLeft] = useState(0);
  const {
    instanceIdMap,
    loading,
    challengeInstanceMap,
    startInstance,
    stopInstance,
    extendInstance,
  } = useInstanceStore();
  const [copiedLink, setCopiedLink] = useState<string | null>(null);

  const instance_id = challengeInstanceMap[challenge_id];
  const instance = instance_id ? instanceIdMap[instance_id] : undefined;

  const isActive = instance?.expires_at ? instance.expires_at * 1000 > Date.now() : false;

  const endpoints = useMemo(() => {
    if (!instance || !instance.endpoints || instance.endpoints.length === 0) return [];
    return instance.endpoints;
  }, [instance]);

  useEffect(() => {
    if (instance?.expires_at) {
      const expiresAtMs = instance.expires_at * 1000;
      const active = expiresAtMs > Date.now();
      if (active) {
        setTimeLeft(Math.max(0, Math.floor((expiresAtMs - Date.now()) / 1000)));
      } else {
        setTimeLeft(0);
      }
    } else {
      setTimeLeft(0);
    }

    let interval: NodeJS.Timeout | undefined;
    if (instance?.expires_at) {
      const expiresAtMs = instance.expires_at * 1000;
      if (expiresAtMs > Date.now()) {
        interval = setInterval(() => {
          const timeRemaining = Math.max(0, Math.floor((expiresAtMs - Date.now()) / 1000));
          setTimeLeft(timeRemaining);

          if (timeRemaining === 0 && interval) {
            clearInterval(interval);
          }
        }, 1000);
      }
    }

    return () => {
      if (interval) clearInterval(interval);
    };
  }, [instance]);

  const handleStart = async () => {
    await startInstance(challenge_id);
  };

  const handleStop = async () => {
    await stopInstance(challenge_id);
  };

  const handleExtend = async () => {
    await extendInstance(challenge_id);
  };

  const formatTime = (seconds: number) => {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, "0")}`;
  };

  const copyToClipboard = (text: string) => {
    try {
      navigator.clipboard.writeText(text);
      setCopiedLink(text);
      setTimeout(() => setCopiedLink(null), 4000);
    } catch {
      showToast(ToastStatus.Failure, "Failed to copy to clipboard");
    }
  };

  const getStatusColor = () => {
    if (!instance) return "text-gray-500";
    if (isActive) return "text-green-500";
    return "text-red-500";
  };

  return (
    <div className="flex flex-col bg-card p-3 rounded-lg shadow-sm border space-y-2">
      <div className="flex items-center space-x-4 justify-between">
        <div className="flex space-x-2 items-center">
          <Button
            onClick={isActive ? handleStop : handleStart}
            disabled={loading}
            variant="outline"
            size="sm"
          >
            {loading ? (
              <Loader2 className="animate-spin text-yellow-500" />
            ) : isActive ? (
              <StopCircle className="text-red-500" />
            ) : (
              <Play className="text-green-500" />
            )}
          </Button>

          {isActive && (
            <Button size="sm" variant="outline" onClick={handleExtend} disabled={loading}>
              <RefreshCw className="h-4 w-4 mr-2" />
              Extend
            </Button>
          )}

          {isActive && (
            <span className="text-sm font-mono border h-9 p-2 rounded-lg">
              {formatTime(timeLeft)}
            </span>
          )}
        </div>

        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger asChild>
              <Badge variant="outline" className={getStatusColor()}>
                {isActive ? "Running" : "Stopped"}
              </Badge>
            </TooltipTrigger>
            <TooltipContent>
              <p>
                {isActive
                  ? "Access the instance using the details below"
                  : "Start the instance to access it"}
              </p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </div>

      {isActive && endpoints.length > 0 && (
        <div className="flex flex-col space-y-2">
          {endpoints.map((endpoint, index) => {
            const connectionString = GenerateChallengeEndpoint(endpoint);
            return (
              <div key={index} className="flex items-center space-x-2">
                <div className="relative flex-grow">
                  <Terminal className="h-5 w-5 absolute left-2 top-1/2 transform -translate-y-1/2 text-gray-500" />
                  <Input
                    value={connectionString}
                    readOnly
                    className="pl-8 truncate font-mono focus:outline-none focus-visible:ring-0 focus-visible:ring-offset-0"
                  />
                </div>
                <CopyButton
                  copiedLink={copiedLink}
                  content={connectionString}
                  copyToClipboard={copyToClipboard}
                />
              </div>
            );
          })}
        </div>
      )}
      {/* {instance.active && instance.password && (
        <div className="flex items-center space-x-2">
          <div className="relative flex-grow">
            <KeyRound className="h-5 w-5 absolute left-2 top-1/2 transform -translate-y-1/2 text-gray-500" />
            <Input
              value={instance.password}
              readOnly
              className="pl-8 truncate font-mono focus:outline-none focus-visible:ring-0 focus-visible:ring-offset-0"
            />
          </div>
          <CopyButton
            copiedLink={copiedLink}
            content={instance.password}
            copyToClipboard={copyToClipboard}
          />
        </div>
      )} */}
    </div>
  );
}
