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

  # Enable LiveDashboard in development
  if Application.compile_env(:pulse, :dev_routes) do
    # If you want to use the LiveDashboard in production, you should put
    # it behind authentication and allow only admins to access it.
    # If your application does not have an admins-only section yet,
    # you can use Plug.BasicAuth to set up some basic authentication
    # as long as you are also using SSL (which you should anyway).
    import Phoenix.LiveDashboard.Router

    scope "/dev" do
      pipe_through(:dashboard)

      live_dashboard("/dashboard", metrics: PulseWeb.Telemetry)
    end
  end
end
