export const API_ROUTES = {
  ping: "/api/v1/ping",

  auth: {
    register: "/api/v1/auth/register",
    verify: "/api/v1/auth/verify",
    login: "/api/v1/auth/login",
    forgotPassword: "/api/v1/auth/forgot-password",
    resetPassword: "/api/v1/auth/reset-password",
    logout: "/api/v1/auth/logout",
  },

  team: {
    create: "/api/v1/team/create",
    join: "/api/v1/team/join",
    invite: {
      generate: "/api/v1/team/invite",
      accept: "/api/v1/team/invite", // GET with token query param
    },
  },

  event: {
    info: "/api/v1/event/info",
  },

  profile: {
    me: "/api/v1/profile/me",
    team: "/api/v1/profile/team",
  },

  challenge: {
    list: "/api/v1/challenge",
    submitFlag: "/api/v1/challenge/submit",
    unlockHint: "/api/v1/challenge/hint/unlock",
  },

  score: {
    scoreboard: "/api/v1/score",
    graph: "/api/v1/score/graph",
  },
} as const;
