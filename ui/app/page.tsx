"use client";

import Link from "next/link";

export default function Home() {
  return (
    <div className="flex flex-col h-5/6 w-full p-4">
      <div className="flex flex-col h-5/6 w-3/4 justify-center self-center gap-3">
        <div className="text-6xl font-extrabold text-center">
          Deploy your linux wargames easily
        </div>
        <div className="antialiased font-Roboto w-3/4 self-center hidden sm:flex text-center">
          Isolet: A framework packed with features like Isolated instances,
          Dynamic Pod Generation - Shape Your Wargames with Ease.
        </div>
        <Link href={`/howtoplay`} className="self-center">
          <button className="p-3 text-black bg-palette-500 rounded-md font-mono">
            Learn more
          </button>
        </Link>
      </div>
      <div className="flex self-center gap-2 w-3/5 justify-around font-mono">
        <Link
          href="https://cyberlabs.club"
          className="flex flex-col items-center gap-2"
        >
          <img
            src={`/static/assets/cyberlabs.png`}
            alt="CL logo"
            width={96}
          ></img>
          <div className="hidden sm:flex">CyberLabs</div>
        </Link>
        <Link
          href="https://github.com/TheAlpha16/isolet"
          className="flex flex-col items-center gap-2"
        >
          <img src={`/static/assets/github-mark-white.svg`} alt="CL logo"></img>
          <div>SourceCode</div>
        </Link>
      </div>
    </div>
  );
}
