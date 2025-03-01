"use client";

import Image from "next/image"
import Link from "next/link"
import { Anchor } from "lucide-react"
import Countdown from 'react-countdown';

import { useMetadataStore } from "@/store/metadataStore";

import googleCloud from "@/icons/google-cloud.png"
import nvcti from "@/icons/nvcti.png"
import pearlLogo from "@/icons/pearl-logo.png"

export default function HomePage() {
  const { eventStart, eventEnd } = useMetadataStore();

  return (
    <div className="min-h-screen text-white flex flex-col items-center">
      {/* Main content */}
      <main className="flex-1 w-full max-w-6xl px-4 py-12 md:py-12">
        {/* Hero section with logo */}
        <section className="flex flex-col items-center justify-center mb-16 relative">
          {/* Cyberpunk circuit lines */}
          <div className="absolute inset-0 pointer-events-none overflow-hidden">
            <div className="absolute top-1/2 left-0 w-full h-px bg-gradient-to-r from-transparent via-blue-500/30 to-transparent"></div>
            <div className="absolute top-0 left-1/4 w-px h-full bg-gradient-to-b from-transparent via-blue-500/20 to-transparent"></div>
            <div className="absolute top-0 right-1/4 w-px h-full bg-gradient-to-b from-transparent via-blue-500/20 to-transparent"></div>
          </div>

          {/* Logo */}
          <div className="relative mb-8 w-64 h-64 md:w-80 md:h-80">
            <div className="relative w-full h-full flex items-center justify-center">
              <Image
                src={pearlLogo.src}
                alt="Pearl CTF Logo"
                width={300}
                height={300}
                className="drop-shadow-[0_0_15px_rgba(59,130,246,0.5)]"
              />
            </div>
          </div>

          {/* Tagline */}
          <h1 className="font-mono text-2xl md:text-4xl text-center mb-6 tracking-wider text-gray-800 dark:text-blue-100">
            <span className="text-blue-400">&gt;</span> DIVE INTO THE CHALLENGE
          </h1>

          <p className="font-mono text-sm md:text-base text-center max-w-2xl text-gray-600 dark:text-blue-200/70 mb-8">
            An oceanic cybersecurity competition where hackers navigate the depths of digital challenges. Explore
            vulnerabilities, decrypt secrets, and capture flags in this immersive CTF experience.
          </p>

          {/* Event status */}
          <div className="flex flex-col sm:flex-row gap-4 mt-4">
            <div className="font-mono text-xl font-bold text-center max-w-2xl text-gray-600 dark:text-blue-200/70">
              {
              Date.now() < eventStart.getTime() ? 
              <>
              <p className="mb-3">Event Starts on ${eventStart.toLocaleString("en-GB", { day: "numeric", month: "short", year: "numeric", hour: "numeric", minute: "2-digit" })}</p>
              <Countdown className="text-4xl" date={eventStart}/>
              </> : 
              Date.now() > eventEnd.getTime() ? 
              <p>Event has Ended</p> : 
              <>
              <p className="mb-3">Event will end on {eventEnd.toLocaleString("en-GB", { day: "numeric", month: "short", year: "numeric", hour: "numeric", minute: "2-digit" })}</p>
              <Countdown className="text-4xl" date={eventEnd}/>
              </>
              }
            </div>
          </div>
        </section>

        {/* Wave divider */}
        <div className="w-full h-12 relative my-16">
          <div className="absolute inset-0 flex items-center">
            <div className="h-px w-full bg-gradient-to-r from-transparent via-blue-500/40 to-transparent"></div>
          </div>
          <div className="absolute inset-0 flex items-center justify-center">
            <Anchor className="h-8 w-8 text-blue-500/70" />
          </div>
        </div>

        {/* Social Links section */}
        <section className="mb-16">
          <h2 className="font-mono text-xl md:text-2xl text-center mb-12 tracking-wider text-gray-800 dark:text-blue-300">
            <span className="text-blue-400">#</span> JOIN OUR COMMUNITY
          </h2>

          <div className="flex flex-wrap justify-center gap-8">
            {[
              {
                icon: (
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="24"
                    height="24"
                    viewBox="0 0.5 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    className="h-8 w-8"
                  >
					<path d="M18.8943 4.34399C17.5183 3.71467 16.057 3.256 14.5317 3C14.3396 3.33067 14.1263 3.77866 13.977 4.13067C12.3546 3.89599 10.7439 3.89599 9.14391 4.13067C8.99457 3.77866 8.77056 3.33067 8.58922 3C7.05325 3.256 5.59191 3.71467 4.22552 4.34399C1.46286 8.41865 0.716188 12.3973 1.08952 16.3226C2.92418 17.6559 4.69486 18.4666 6.4346 19C6.86126 18.424 7.24527 17.8053 7.57594 17.1546C6.9466 16.92 6.34927 16.632 5.77327 16.2906C5.9226 16.184 6.07194 16.0667 6.21061 15.9493C9.68793 17.5387 13.4543 17.5387 16.889 15.9493C17.0383 16.0667 17.177 16.184 17.3263 16.2906C16.7503 16.632 16.153 16.92 15.5236 17.1546C15.8543 17.8053 16.2383 18.424 16.665 19C18.4036 18.4666 20.185 17.6559 22.01 16.3226C22.4687 11.7787 21.2836 7.83202 18.8943 4.34399ZM8.05593 13.9013C7.01058 13.9013 6.15725 12.952 6.15725 11.7893C6.15725 10.6267 6.98925 9.67731 8.05593 9.67731C9.11191 9.67731 9.97588 10.6267 9.95454 11.7893C9.95454 12.952 9.11191 13.9013 8.05593 13.9013ZM15.065 13.9013C14.0196 13.9013 13.1652 12.952 13.1652 11.7893C13.1652 10.6267 13.9983 9.67731 15.065 9.67731C16.121 9.67731 16.985 10.6267 16.9636 11.7893C16.9636 12.952 16.1317 13.9013 15.065 13.9013Z" strokeLinejoin="round"/>
				  </svg>
                ),
                name: "Discord",
                handle: "PearlCTF",
                color: "bg-[#5865F2]/20 border-[#5865F2]/40 hover:bg-[#5865F2]/30",
				url: "https://discord.gg/MuUwXX9P"
              },
              {
                icon: (
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    className="h-8 w-8"
                  >
                    <rect width="20" height="20" x="2" y="2" rx="5" ry="5" />
                    <path d="M16 11.37A4 4 0 1 1 12.63 8 4 4 0 0 1 16 11.37z" />
                    <line x1="17.5" x2="17.51" y1="6.5" y2="6.5" />
                  </svg>
                ),
                name: "Instagram",
                handle: "@cyberlabs_iitism",
                color: "bg-[#E1306C]/20 border-[#E1306C]/40 hover:bg-[#E1306C]/30",
				url: "https://www.instagram.com/cyberlabs_iitism?igsh=cGMzOXlxbnV4MWY3"
              },
              {
                icon: (
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    className="h-8 w-8"
                  >
                    <path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4" />
                    <path d="M9 18c-4.51 2-5-2-7-2" />
                  </svg>
                ),
                name: "GitHub",
                handle: "Cyber Labs",
                color: "bg-gray-700/30 border-gray-500/40 hover:bg-gray-700/50",
				url: "https://github.com/Cyber-Labs"
              },
              {
                icon: (
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    className="h-8 w-8"
                  >
                    <path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z" />
                    <rect width="4" height="12" x="2" y="9" />
                    <circle cx="4" cy="4" r="2" />
                  </svg>
                ),
                name: "LinkedIn",
                handle: "CyberLabs IITISM",
                color: "bg-[#0077B5]/20 border-[#0077B5]/40 hover:bg-[#0077B5]/30",
				url: "https://in.linkedin.com/company/cyberlabs-iitism"
              },
            ].map((social, index) => (
              <Link
                key={index}
                href={social.url}
                className={`flex flex-col items-center p-6 border rounded-md transition-all ${social.color} group w-64 h-64 justify-center relative overflow-hidden`}
              >
                <div className="absolute inset-0 bg-blue-500/5 opacity-0 group-hover:opacity-100 transition-opacity">
                  <div className="absolute bottom-0 left-0 w-full h-1 bg-gradient-to-r from-blue-400/0 via-blue-400/70 to-blue-400/0"></div>
                  <div className="absolute top-0 right-0 w-1 h-full bg-gradient-to-b from-blue-400/0 via-blue-400/70 to-blue-400/0"></div>
                </div>

                <div className="text-blue-400 mb-4 transform group-hover:scale-110 transition-transform">
                  {social.icon}
                </div>
                <h3 className="font-mono text-xl mb-2 text-gray-800 dark:text-blue-200">{social.name}</h3>
                <p className="text-gray-600 dark:text-blue-200/70 text-sm">{social.handle}</p>

                <div className="mt-6 font-mono text-xs text-blue-400 opacity-0 group-hover:opacity-100 transition-opacity">
                  CONNECT &gt;
                </div>
              </Link>
            ))}
          </div>
        </section>

        {/* Sponsors section */}
        <section>
          <h2 className="font-mono text-xl md:text-2xl text-center mb-12 tracking-wider text-gray-800 dark:text-blue-300">
            <span className="text-blue-400">#</span> SPONSORS
          </h2>

          <div className="gap-4 flex justify-center flex-wrap">
            {[
				{
					name: "Google Cloud",
					icon: googleCloud
				},
				{
					name: "NVCTI",
					icon: nvcti
				}
			].map((sponsor, index) => (
              <div key={index} className="w-32 h-32 md:w-40 md:h-40 relative group">
                <div className="absolute inset-0 flex items-center justify-center p-4 flex-col gap-2">
                  <Image
                    src={sponsor.icon.src}
                    alt={sponsor.name}
                    width={80}
                    height={80}
                    className="max-w-full max-h-full"
                  />
				  <p className="font-bold text-slate-400 text-lg h-full flex items-end">{sponsor.name}</p>
                </div>
                <div className="absolute inset-0 rounded-md opacity-0 group-hover:opacity-100 transition-opacity">
                  <div className="absolute bottom-0 left-0 w-full h-1 bg-gradient-to-r from-blue-400/0 via-blue-400/70 to-blue-400/0 transform translate-y-full group-hover:translate-y-0 transition-transform"></div>
                </div>
              </div>
            ))}
          </div>
        </section>
      </main>
    </div>
  )
}

