"use client";

import React, { useMemo } from "react";
import { useInstanceStore } from "@/store/instance";
import { useCopyToClipboard } from "@/components/instances/useCopyToClipboard";
import { useInstanceTimer } from "@/components/instances/useInstanceTimer";
import { InstanceControls } from "@/components/instances/InstanceControls";
import { InstanceStatusBadge } from "@/components/instances/InstanceStatusBadge";
import { EndpointList } from "@/components/instances/EndpointList";

interface InstanceCardProps {
  challenge_id: number;
}

export function InstanceCard({ challenge_id }: InstanceCardProps) {
  const {
    instanceIdMap,
    loading,
    challengeInstanceMap,
    startInstance,
    stopInstance,
    extendInstance,
  } = useInstanceStore();

  const instance_id = challengeInstanceMap[challenge_id];
  const instance = instance_id ? instanceIdMap[instance_id] : undefined;

  const isActive = instance?.expires_at ? instance.expires_at * 1000 > Date.now() : false;
  const timeLeft = useInstanceTimer(instance);
  const { copiedLink, copyToClipboard } = useCopyToClipboard();

  const endpoints = useMemo(() => {
    if (!instance || !instance.endpoints || instance.endpoints.length === 0) return [];
    return instance.endpoints;
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

  return (
    <div className="flex flex-col bg-card p-3 rounded-lg shadow-sm border space-y-2">
      <div className="flex items-center space-x-4 justify-between">
        <InstanceControls
          isActive={isActive}
          loading={loading}
          timeLeft={timeLeft}
          onStart={handleStart}
          onStop={handleStop}
          onExtend={handleExtend}
        />
        <InstanceStatusBadge isActive={isActive} hasInstance={!!instance} />
      </div>

      {isActive && (
        <EndpointList endpoints={endpoints} copiedLink={copiedLink} onCopy={copyToClipboard} />
      )}
    </div>
  );
}
