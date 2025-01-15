import { ChallengeType, useChallengeStore } from "@/store/challengeStore";
import FormButton from "@/components/extras/buttons";
import { useState } from "react";

function ChallInfo(challenge: ChallengeType) {
    return (
        <div className="flex flex-col gap-2">
            <div className="text-2xl font-bold">{challenge.type}</div>
            <div className="text-2xl font-bold">{challenge.name}</div>
            <div className="text-lg">{challenge.prompt}</div>
            <div className="text-lg">Points: {challenge.points}</div>
            <div className="text-lg">Solves: {challenge.solves}</div>
            <div className="text-lg">Author: {challenge.author}</div>
            <div className="text-lg">Tags: {challenge.tags.join(", ")}</div>
            <div className="text-lg">Links: 
                {challenge.links.map((link) => {
                    return (
                        <a key={link} href={link} target="_blank" rel="noreferrer">
                            {link}
                        </a>
                    );
                }
            )}
            </div>
        </div>
    );
};

function FlagInput({ chall_id }: { chall_id: number }) { 
    const [flag, setFlag] = useState("");
    const { submitFlag } = useChallengeStore();

    const flagSubmit = async () => {
        await submitFlag(chall_id, flag);
    }

    return (
        <form
            onSubmit={(event) => {
                event.preventDefault();
                flagSubmit();
            }}>
            <input
                id={ `${chall_id}-input` }
                type="text"
                name="flag"
                placeholder="flag"
                value={flag}
                onChange={(event) => {
                    setFlag(event.target.value);
                }}
                required
            ></input>
            <FormButton type="submit">
                Submit
            </FormButton>
        </form>
    )
};

export function StaticChallenge(challenge: ChallengeType) {
    return (
        <div className="flex flex-col gap-4">
            <ChallInfo {...challenge} />
            <FlagInput chall_id={challenge.chall_id} />
        </div>
    );
}