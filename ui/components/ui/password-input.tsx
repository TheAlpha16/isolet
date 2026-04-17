"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Eye, EyeClosed } from "lucide-react";
import { useState } from "react";

type PasswordInputProps = Omit<React.ComponentProps<typeof Input>, "type">;

export function PasswordInput({ className, ...props }: PasswordInputProps) {
  const [show, setShow] = useState(false);

  return (
    <div className="relative">
      <Input type={show ? "text" : "password"} className={`pr-10 ${className ?? ""}`} {...props} />
      <Button
        type="button"
        variant="ghost"
        size="icon"
        className="absolute inset-y-0 right-0"
        onClick={() => setShow((s) => !s)}
        tabIndex={-1}
      >
        {show ? <Eye className="h-5 w-5" /> : <EyeClosed className="h-5 w-5" />}
      </Button>
    </div>
  );
}
