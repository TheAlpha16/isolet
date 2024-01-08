import { createContext } from "react"

interface challData {
	deadline: number;
}

export const challContext = createContext<challData>({
	deadline: 0
})
