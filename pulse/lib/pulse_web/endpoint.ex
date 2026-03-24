defmodule PulseWeb.Endpoint do
  use Phoenix.Endpoint, otp_app: :pulse

  @session_options [
    store: :cookie,
    key: "_pulse_key",
    signing_salt: "z451qUF9",
    same_site: "Lax"
  ]

  socket("/live", Phoenix.LiveView.Socket,
    websocket: [connect_info: [session: @session_options]],
    longpoll: [connect_info: [session: @session_options]]
  )

  socket("/socket", PulseWeb.UserSocket,
    websocket: [log: false],
    longpoll: false
  )

  if code_reloading? do
    plug(Phoenix.CodeReloader)
  end

  plug(Phoenix.LiveDashboard.RequestLogger,
    param_key: "request_logger",
    cookie_key: "request_logger"
  )

  plug(Plug.RequestId)
  plug(Plug.Telemetry, event_prefix: [:phoenix, :endpoint])

  plug(Plug.Parsers,
    parsers: [:json],
    pass: ["*/*"],
    json_decoder: Phoenix.json_library()
  )

  plug(PulseWeb.Router)
end
