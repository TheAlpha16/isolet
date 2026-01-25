"use client";

import React from "react";
import type { Endpoint } from "@/api";
import { EndpointUtils } from "@/lib/endpoint";
import { EndpointInput } from "./EndpointInput";

interface EndpointListProps {
  endpoints: Endpoint[];
  copiedLink: string | null;
  onCopy: (text: string) => void;
}

export function EndpointList({ endpoints, copiedLink, onCopy }: EndpointListProps) {
  if (endpoints.length === 0) return null;

  return (
    <div className="flex flex-col space-y-2">
      {endpoints.map((endpoint, index) => {
        const connectionString = EndpointUtils.toConnectionString(endpoint);
        return (
          <EndpointInput
            key={index}
            connectionString={connectionString}
            copiedLink={copiedLink}
            onCopy={onCopy}
          />
        );
      })}
    </div>
  );
}
