import { showHint } from "@/components/hints/HintToastContainer";
import { HintUnlockConfirmation } from "@/components/hints/HintUnlockComponent";
import { Button } from "@/components/ui/button";
import { HintUI } from "@/models/hint";
import { useChallengeStore } from "@/store";
import { Lightbulb, Lock } from "lucide-react";
import { useState } from "react";

export default function Hint(hint: HintUI) {
  const { unlockHint } = useChallengeStore();
  const [isConfirmationOpen, setIsConfirmationOpen] = useState(false);

  const handleUnlock = () => {
    setIsConfirmationOpen(true);
  };

  const confirmUnlock = () => {
    unlockHint(hint.challenge_id, hint.id);
    setIsConfirmationOpen(false);
  };

  return (
    <>
      {hint.unlocked ? (
        <Button variant={"outline"} size={"sm"} onClick={() => showHint(hint.text)}>
          <Lightbulb className="text-green-500" />
        </Button>
      ) : (
        <Button variant={"secondary"} size={"sm"} onClick={handleUnlock}>
          <Lock className="text-yellow-500" />
        </Button>
      )}
      <HintUnlockConfirmation
        isOpen={isConfirmationOpen}
        onClose={() => setIsConfirmationOpen(false)}
        onConfirm={confirmUnlock}
        cost={hint.cost}
      />
    </>
  );
}
