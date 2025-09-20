"use client";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { CopyButton } from "@/components/utils/copy-button";
import useTeamInvite from "@/hooks/useTeamInvite";
import showToast, { ToastStatus } from "@/utils/toastHelper";
import { Loader2, RefreshCw } from "lucide-react";
import { useState } from "react";

interface TeamInviteProps {
  isOpen: boolean;
  onClose: () => void;
}

export function TeamInvite({ isOpen, onClose }: TeamInviteProps) {
  const [copiedLink, setCopiedLink] = useState<string | null>(null);
  const { inviteLink, loading, generateInvite } = useTeamInvite();

  const onGenerate = async () => {
    await generateInvite();
  };

  const copyToClipboard = (text: string) => {
    try {
      navigator.clipboard.writeText(text);
      setCopiedLink(text);
      setTimeout(() => setCopiedLink(null), 4000);
    } catch {
      showToast(ToastStatus.Failure, "failed to copy to clipboard");
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Invite</DialogTitle>
        </DialogHeader>
        <DialogDescription>Share this link with others to invite.</DialogDescription>
        <div className="flex items-center space-x-2">
          <Button onClick={onGenerate} variant="outline" size={"icon"} disabled={loading}>
            {loading ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <RefreshCw className="h-4 w-4" />
            )}
          </Button>
          <Input
            value={inviteLink}
            className="truncate font-mono focus:outline-none focus-visible:ring-0 focus-visible:ring-offset-0"
            readOnly
          />
          <CopyButton
            copiedLink={copiedLink}
            content={inviteLink}
            copyToClipboard={copyToClipboard}
          />
        </div>
      </DialogContent>
    </Dialog>
  );
}
