defmodule PulseWeb.Router do
  use PulseWeb, :router

  pipeline :api do
    plug(:accepts, ["json"])
  end

  pipeline :dashboard do
    plug(Plug.Session,
      store: :cookie,
      key: "_pulse_key",
      signing_salt: "z451qUF9",
      same_site: "Lax"
    )

    plug(:fetch_session)
    plug(:protect_from_forgery)
  end

  scope "/api", PulseWeb do
    pipe_through(:api)

    get("/health", HealthController, :index)
  end

  if Application.compile_env(:pulse, :dev_routes) do
    import Phoenix.LiveDashboard.Router

    scope "/dev" do
      pipe_through(:dashboard)

      live_dashboard("/dashboard", metrics: PulseWeb.Telemetry)
    end
  end
end
