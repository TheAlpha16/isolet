import { ChallengeType, useChallengeStore } from "@/store/challengeStore";
import FormButton from "@/components/extras/buttons";
import React, { useState } from "react";
import Image from "next/image";

function ChallInfo(challenge: ChallengeType) {
    return (
        <div className="flex flex-col gap-2 text-lg">
            <div className="text-2xl font-bold">{challenge.name}</div>
            <div>{challenge.prompt}</div>
            <div>Points: {challenge.points}</div>
            <div>Solves: {challenge.solves}</div>
            <div>Author: {challenge.author}</div>
            <div>Tags: {challenge.tags.join(", ")}</div>
            <div>Links: 
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
            }}
            className="flex gap-2">
            <input
                id={ `${chall_id}-input` }
                type="text"
                name="flag"
                placeholder="flag"
                value={flag}
                onChange={(event) => {
                    setFlag(event.target.value);
                }}
                className="rounded-md h-10 p-2"
                required
            ></input>
            <FormButton type="submit">
                Submit
            </FormButton>
        </form>
    )
};

export function StaticChallenge({ challenge, closeChallenge = () => {}, isFocussed = false, onClick = () => {}}: {
    challenge: ChallengeType,
    closeChallenge?: () => void,
    isFocussed?: boolean,
    onClick?: () => void
}) {
    return (
    <>
        {
            isFocussed ? (
                    <div className="flex flex-col gap-4 p-4 w-2/3 rounded-xl bg-[rgba(255,255,255,0.1)] relative" onClick={(event) => {event.stopPropagation()}}>
                    <Image
                        className="svg-icon hover:cursor-pointer absolute -top-3 -right-3 bg-white rounded-full"
                        src="/close.svg"
                        alt="close"
                        width={24}
                        height={24}
                        onClick={closeChallenge}
                    ></Image>
                    <ChallInfo {...challenge} />
                    <FlagInput chall_id={challenge.chall_id} />
                </div>
            ) : (
                <div className="flex flex-col w-56 h-24 bg-[#1a1a1a] rounded-md p-4 place-items-center justify-around hover:cursor-pointer" onClick={onClick}>
                    <div className="truncate">{ challenge.name }</div>
                    <div>{ challenge.points }</div>
                </div>
            )
        }
    </>);
}