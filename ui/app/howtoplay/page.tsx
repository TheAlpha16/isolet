'use client'

export default function HowToPlay() {
	return (
		<div className="flex min-h-screen flex-col items-center p-4 font-mono gap-5">
			<div className="font-bold text-3xl">
				HOW TO PLAY?
			</div>
			<div className="flex flex-col sm:flex-row sm:justify-between gap-2">
				<img src={`/static/assets/info.png`} alt={`info image`} className="sm:w-2/5"/>
				<div className="flex items-center grow justify-center">
					Read the description about the challenge
				</div>
			</div>
			<div className="flex flex-col sm:flex-row-reverse sm:justify-between gap-2">
				<img src={`/static/assets/start.png`} alt={`info image`} className="sm:w-2/5"/>
				<div className="flex items-center grow justify-center">
					Click on Start button to spawn a personal instance
				</div>
			</div>
			<div className="flex flex-col sm:flex-row sm:justify-between gap-2">
				<img src={`/static/assets/creds.png`} alt={`info image`} className="sm:w-2/5"/>
				<div className="flex items-center grow justify-center">
					Click on the button next to password to copy
				</div>
			</div>
			<div className="flex flex-col sm:flex-row-reverse sm:justify-between gap-2">
				<img src={`/static/assets/submit.png`} alt={`info image`} className="sm:w-2/5"/>
				<div className="flex items-center grow justify-center">
					Submit the flag after pwning the machine
				</div>
			</div>
			<div className="flex flex-col sm:flex-row sm:justify-between gap-2">
				<img src={`/static/assets/stop.png`} alt={`info image`} className="sm:w-2/5"/>
				<div className="flex items-center grow justify-center">
					Be sure to stop the instance after working
				</div>
			</div>
		</div>
	)
}
