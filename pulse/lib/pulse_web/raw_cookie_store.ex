defmodule PulseWeb.RawCookieStore do
  @moduledoc """
  A custom plug session store used exclusively to read raw HttpOnly cookies
  during the Phoenix WebSocket upgrade handshake.

  Because Phoenix Channels do not expose the raw `cookie` header natively
  through `connect_info`, we use this store to bypass standard serialization
  and directly extract the JWT token appended by the external Go API.
  """
  @behaviour Plug.Session.Store

  @impl true
  def init(opts), do: opts

  @impl true
  def get(_conn, cookie_value, _opts) do
    # When used in connect_info, `cookie_value` is the raw value of the cookie
    # identified by the `:key` option. We wrap it in a map as the "session".
    {cookie_value, %{"token" => cookie_value}}
  end

  @impl true
  def put(_conn, _sid, _any, _opts), do: ""

  @impl true
  def delete(_conn, _sid, _opts), do: :ok
end
