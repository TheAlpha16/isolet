"use client";

import React from "react";
import { Badge } from "@/components/ui/badge";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";

interface InstanceStatusBadgeProps {
  isActive: boolean;
  hasInstance: boolean;
}

function getStatusColor(hasInstance: boolean, isActive: boolean): string {
  if (!hasInstance) return "text-gray-500";
  if (isActive) return "text-green-500";
  return "text-red-500";
}

export function InstanceStatusBadge({ isActive, hasInstance }: InstanceStatusBadgeProps) {
  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <Badge variant="outline" className={getStatusColor(hasInstance, isActive)}>
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
  );
}
