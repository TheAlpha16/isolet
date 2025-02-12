"use client"

import React from 'react';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { Button } from "@/components/ui/button";
import { Copy, Check } from 'lucide-react';

interface CopyButtonProps {
    copiedLink: string | null;
    content: string;
    copyToClipboard: (text: string) => void;
}

export function CopyButton({ copiedLink, content, copyToClipboard }: CopyButtonProps) {

    return (
        <TooltipProvider>
            <Tooltip>
                <TooltipTrigger asChild>
                    <Button
                        variant="outline"
                        size="icon"
                        onClick={() => copyToClipboard(content)}
                    >
                        {copiedLink === content ? (
                            <Check className="h-4 w-4 text-green-500" />
                        ) : (
                            <Copy className="h-4 w-4" />
                        )}
                    </Button>
                </TooltipTrigger>
                <TooltipContent>
                    <p>
                        {copiedLink === content ? "Copied!" : "Copy to clipboard"}
                    </p>
                </TooltipContent>
            </Tooltip>
        </TooltipProvider>
    );
}
