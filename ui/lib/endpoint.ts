import type { Endpoint } from "@/api";

/**
 * Utility functions for working with Endpoint objects
 */
export const EndpointUtils = {
  /**
   * Generates a connection string from an Endpoint object
   * @param endpoint - The endpoint object containing protocol, hostname, and optional port
   * @param username - Username for SSH connections (default: "hacker")
   * @returns A formatted connection string
   */
  toConnectionString(endpoint: Endpoint, username: string = "hacker"): string {
    let connString = "";
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
  },

  /**
   * Checks if an endpoint is ready for connections
   * @param endpoint - The endpoint to check
   * @returns True if the endpoint is ready
   */
  isReady(endpoint: Endpoint): boolean {
    return endpoint.ready;
  },

  /**
   * Gets the display name for an endpoint
   * @param endpoint - The endpoint
   * @returns A human-readable name for the endpoint
   */
  getDisplayName(endpoint: Endpoint): string {
    return endpoint.name || endpoint.protocol;
  },
};
