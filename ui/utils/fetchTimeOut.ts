function fetchTimeout(url: string, ms: number, signal: AbortSignal, options = {}) {
    const controller = new AbortController()
    // const promise = fetch(url, { signal: controller.signal, ...options })

    if (signal) signal.addEventListener("abort", () => controller.abort(), true)
    
    // const timeout = setTimeout(() => controller.abort(), ms)
    // return promise.finally(() => clearTimeout(timeout))

    return fetch(url, { signal: controller.signal, ...options })
}

export default fetchTimeout;
