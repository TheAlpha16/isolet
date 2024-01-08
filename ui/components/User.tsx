'use client'

import { useContext, createContext } from 'react'

export interface stateVars {
	loggedin: boolean
	setLoggedin: any
	respHook: boolean
}

export const Context = createContext<stateVars>({
	loggedin: false,
	setLoggedin: false,
	respHook: false
})


export default function User() {
	return useContext(Context)
}
