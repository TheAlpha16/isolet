defmodule PulseWeb.TeamChannel do
  use PulseWeb, :channel

  @impl true
  def join("team:" <> team_id, _payload, socket) do
    if socket.assigns.team_id == team_id do
      {:ok, socket}
    else
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
