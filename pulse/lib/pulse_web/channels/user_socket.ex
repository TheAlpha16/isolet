defmodule PulseWeb.UserSocket do
  use Phoenix.Socket, log: false

  alias Pulse.Auth
  require Logger

  ## Channels
  channel "team:*", PulseWeb.TeamChannel
  channel "global", PulseWeb.GlobalChannel

  @impl true
  def connect(params, socket, _connect_info) do
    token = params["token"]

    with {:ok, t} when not is_nil(t) <- {:ok, token},
         {:ok, claims} <- Auth.verify_jwt(t),
         :ok <- Auth.validate_claims(claims),
         :ok <- Auth.validate_session(claims) do
      Logger.debug(
        "UserSocket connected: user_id=#{claims["user_id"]}, team_id=#{claims["team_id"]}"
      )

      {:ok,
       socket
       |> assign(:user_id, claims["user_id"])
       |> assign(:team_id, claims["team_id"])
       |> assign(:jti, claims["jti"])
       |> assign(:role, claims["role"])}
    else
      _ ->
        Logger.debug("UserSocket connection failed")
        :error
    end
  end

  @impl true
  def id(socket), do: "user_socket:#{socket.assigns.user_id}"
end
