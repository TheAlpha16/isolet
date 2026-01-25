import type { Endpoint } from "../api/models/Endpoint";

function GenerateChallengeEndpoint(endpoint: Endpoint, username: string = "hacker"): string {
  let connString: string = "";
  const { protocol, hostname, port } = endpoint;

  switch (protocol) {
    case "http":
    case "https":
      connString = `${protocol}://${hostname}`;
      if (port && port !== 80 && port !== 443) {
        connString += `:${port}`;
      }
      break;

    case "ssh":
      connString = `ssh ${username}@${hostname}`;
      if (port && port !== 22) {
        connString += ` -p ${port}`;
      }
      break;

    case "nc":
      connString = `nc ${hostname}`;
      if (port) {
        connString += ` ${port}`;
      }
      break;

    default:
      // Generic format for unknown protocols
      connString = `${protocol}://${hostname}`;
      if (port) {
        connString += `:${port}`;
      }
      break;
  }

  return connString;
}

export { GenerateChallengeEndpoint };
