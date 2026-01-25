"use client";

import React from "react";
import { Terminal } from "lucide-react";
import { Input } from "@/components/ui/input";
import { CopyButton } from "@/components/utils/copy-button";

interface EndpointInputProps {
  connectionString: string;
  copiedLink: string | null;
  onCopy: (text: string) => void;
}

export function EndpointInput({ connectionString, copiedLink, onCopy }: EndpointInputProps) {
  return (
    <div className="flex items-center space-x-2">
      <div className="relative flex-grow">
        <Terminal className="h-5 w-5 absolute left-2 top-1/2 transform -translate-y-1/2 text-gray-500" />
        <Input
          value={connectionString}
          readOnly
          className="pl-8 truncate font-mono focus:outline-none focus-visible:ring-0 focus-visible:ring-offset-0"
        />
      </div>
      <CopyButton copiedLink={copiedLink} content={connectionString} copyToClipboard={onCopy} />
    </div>
  );
}
