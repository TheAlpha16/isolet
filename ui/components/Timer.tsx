'use client'
import { useState, useEffect, useContext } from "react"
import { challContext } from "./Contexts"

interface Props {
	level: number
	classes: string
	setTimeLeft: any
}

function Timer(props: Props) {
	const { deadline } = useContext(challContext)

	const calculateTimeLeft = (deadline: number) => {
		let difference = deadline - Date.now()
		let timeLeft = {
			hours: 0,
			minutes: 0,
			seconds: 0
		}
	
		if (difference > 0) {
			timeLeft = {
				hours: Math.floor((difference / (1000 * 60 * 60))),
				minutes: Math.floor((difference / 1000 / 60) % 60),
				seconds: Math.floor((difference / 1000) % 60)
			}
		}
		return timeLeft
	}

	const [timeLeft, setTimeLeft] = useState(calculateTimeLeft(deadline))

	useEffect(() => {
		const timer = setTimeout(() => {
			setTimeLeft(calculateTimeLeft(deadline))
		}, 1000)

		if (timeLeft.hours === 0 && timeLeft.minutes === 0 && timeLeft.seconds === 0) {
			props.setTimeLeft(true)
		}

		return () => clearTimeout(timer)
	})

	return (
		<div data-level={ props.level } className={`${props.classes}`}>
			<div data-level={ props.level }>{`${('0' + timeLeft.hours).slice(-2)}h ${('0' + timeLeft.minutes).slice(-2)}m ${('0' + timeLeft.seconds).slice(-2)}s`}</div>
		</div>
	)
}

export default Timer