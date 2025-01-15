import { HintType } from "@/store/challengeStore";
import { TooltipProvider, Tooltip, TooltipContent, TooltipTrigger } from "@radix-ui/react-tooltip";
import { Button } from "@/components/ui/button";
import { Lightbulb, Lock } from "lucide-react";
import { useChallengeStore } from "@/store/challengeStore";

function Hint(hint: HintType) {
	// const { unlockHint } = useChallengeStore();

	const unlockHint = (hid: number) => {
		console.log(`Unlocking hint ${hid}`);
		// hint.unlocked = true;
	}

	return (
		<>{
			hint.unlocked ? (
			<Button
				variant={"outline"}
				size={"sm"}
				onClick={() => alert(hint.hint)}
			><Lightbulb className="text-green-500"/></Button>
			) : (
			<Button
				variant={"secondary"}
				size={"sm"}
				onClick={() => unlockHint(hint.hid)}
			><Lock className="text-yellow-500"/></Button>
			)
		}</>
	);
};

export default Hint;
