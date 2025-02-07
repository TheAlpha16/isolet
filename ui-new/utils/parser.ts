function GenerateChallengeEndpoint(
    method: string,
    domain: string,
    port: number,
    username: string = "hacker",
): string {
    let connString: string = "";

    switch (method) {
        case "http":
            if (port === 80) {
                connString = `http://${domain}`;
            } else if (port === 443) {
                connString = `https://${domain}`;
            } else {
                connString = `http://${domain}:${port}`;
            }
            break;

        case "ssh":
            let user: string = username;

            if (port === 22) {
                connString = `ssh ${user}@${domain}`;
            } else {
                connString = `ssh ${user}@${domain} -p ${port}`;
            }
            break;

        case "nc":
            connString = `nc ${domain} ${port}`;
            break;
    }

    return connString;
}

export { GenerateChallengeEndpoint };
