defmodule PulseWeb.TeamChannel do
  use PulseWeb, :channel
  require Logger

  @impl true
  def join("team:" <> team_id, _payload, socket) do
    if to_string(socket.assigns.team_id) == team_id do
      Logger.debug("User #{socket.assigns.user_id} joined channel team:#{team_id}")
      {:ok, socket}
    else
      Logger.debug(
        "User #{socket.assigns.user_id} refused join team:#{team_id} (assigns.team_id=#{socket.assigns.team_id})"
      )

      {:error, %{reason: "unauthorized"}}
    end
  end

  @impl true
  def handle_in(_event, _payload, socket) do
    {:reply, {:error, %{reason: "read_only"}}, socket}
  end

  # Channels in Pulse are read-only for clients.
  # Events are pushed from the Kafka consumer via the Endpoint.
end
