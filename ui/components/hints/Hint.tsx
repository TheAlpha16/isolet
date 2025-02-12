import type { HintType } from "@/utils/types";
import { Button } from "@/components/ui/button";
import { Lightbulb, Lock } from "lucide-react";
import { useChallengeStore } from "@/store/challengeStore";
import { showHint } from "@/components/hints/HintToastContainer";
import { useState } from "react";
import { HintUnlockConfirmation } from "@/components/hints/HintUnlockComponent";

function Hint(hint: HintType) {
	const { unlockHint } = useChallengeStore();
  	const [isConfirmationOpen, setIsConfirmationOpen] = useState(false);

	const handleUnlock = () => {
		setIsConfirmationOpen(true);
	};

	const confirmUnlock = () => {
		unlockHint(hint.chall_id, hint.hid);
		setIsConfirmationOpen(false);
	};

	return (
		<>{
			hint.unlocked ? (
			<Button
				variant={"outline"}
				size={"sm"}
				onClick={() => showHint(hint.hint)}
			><Lightbulb className="text-green-500"/></Button>
			) : (
			<Button
				variant={"secondary"}
				size={"sm"}
				onClick={handleUnlock}
			><Lock className="text-yellow-500"/></Button>
			)
		}
		<HintUnlockConfirmation
			isOpen={isConfirmationOpen}
			onClose={() => setIsConfirmationOpen(false)}
			onConfirm={confirmUnlock}
			cost={hint.cost}
		/></>
	);
};

export default Hint;
