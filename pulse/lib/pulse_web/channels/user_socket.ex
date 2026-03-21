defmodule PulseWeb.UserSocket do
  use Phoenix.Socket

  alias Pulse.Auth

  ## Channels
  channel "team:*", PulseWeb.TeamChannel
  channel "global", PulseWeb.GlobalChannel

  @impl true
  def connect(%{"token" => token}, socket, _connect_info) do
    with {:ok, claims} <- Auth.verify_jwt(token),
         :ok <- Auth.validate_claims(claims),
         :ok <- Auth.validate_session(claims) do
      {:ok,
       socket
       |> assign(:user_id, claims["user_id"])
       |> assign(:team_id, claims["team_id"])
       |> assign(:jti, claims["jti"])
       |> assign(:role, claims["role"])}
    else
      _ -> :error
    end
  end

  def connect(_params, _socket, _connect_info), do: :error

  @impl true
  def id(socket), do: "user_socket:#{socket.assigns.user_id}"
end
