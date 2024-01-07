'use client'
import { useState, useEffect, useRef } from "react"

interface Props {
    deadline: number
	level: number
	classes: string
}

function Timer(props: Props) {
    const [deadline, setDeadLine] = useState(props.deadline)

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
			setTimeLeft(calculateTimeLeft(deadline));
		}, 1000)

		return () => clearTimeout(timer)
	})

    useEffect(() => {
        setDeadLine(props.deadline)
    }, [props.deadline])

	return (
		<div data-level={ props.level } className={`${props.classes}`}>
			<div data-level={ props.level }>{`${('0' + timeLeft.hours).slice(-2)}h ${('0' + timeLeft.minutes).slice(-2)}m ${('0' + timeLeft.seconds).slice(-2)}s`}</div>
		</div>
	)
}

export default Timer