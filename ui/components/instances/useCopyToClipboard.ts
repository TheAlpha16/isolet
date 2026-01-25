"use client";

import { useState } from "react";
import showToast, { ToastStatus } from "@/utils/toastHelper";

export function useCopyToClipboard() {
  const [copiedLink, setCopiedLink] = useState<string | null>(null);

  const copyToClipboard = (text: string) => {
    try {
      navigator.clipboard.writeText(text);
      setCopiedLink(text);
      setTimeout(() => setCopiedLink(null), 4000);
    } catch {
      showToast(ToastStatus.Failure, "Failed to copy to clipboard");
    }
  };

  return { copiedLink, copyToClipboard };
}
