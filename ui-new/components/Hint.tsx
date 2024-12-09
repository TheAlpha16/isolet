import { HintType } from "@/store/challengeStore";

function Hint(hint: HintType) {
    return (
        <div className="flex gap-2 items-center">
            <div className="text-lg font-bold">{hint.cost}</div>
            <div className="text-lg">{hint.hint}</div>
        </div>
    );
};

export default Hint;
