defmodule PulseWeb do
  @moduledoc """
  The entrypoint for defining your web interface.

  This can be used in your application as:

      use PulseWeb, :controller
      use PulseWeb, :channel
      use PulseWeb, :router
  """

  def static_paths, do: ~w()

  def router do
    quote do
      use Phoenix.Router, helpers: false

      import Plug.Conn
      import Phoenix.Controller
    end
  end

  def channel do
    quote do
      use Phoenix.Channel, log_join: false, log_handle_in: false
    end
  end

  def controller do
    quote do
      use Phoenix.Controller, formats: [:json]

      import Plug.Conn

      unquote(verified_routes())
    end
  end

  def verified_routes do
    quote do
      use Phoenix.VerifiedRoutes,
        endpoint: PulseWeb.Endpoint,
        router: PulseWeb.Router,
        statics: PulseWeb.static_paths()
    end
  end

  defmacro __using__(which) when is_atom(which) do
    apply(__MODULE__, which, [])
  end
end
