defmodule PulseWeb.Router do
  use PulseWeb, :router

  pipeline :api do
    plug(:accepts, ["json"])
  end

  scope "/api", PulseWeb do
    pipe_through(:api)

    get("/health", HealthController, :index)
  end

  if Application.compile_env(:pulse, :dev_routes) do
    import Phoenix.LiveDashboard.Router

    scope "/dev" do
      pipe_through([:fetch_session, :protect_from_forgery])

      live_dashboard("/dashboard", metrics: PulseWeb.Telemetry)
    end
  end
end
