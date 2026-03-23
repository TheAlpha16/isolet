defmodule PulseWeb.GlobalChannel do
  use PulseWeb, :channel
  require Logger

  @impl true
  def join("global", _payload, socket) do
    # Open to all authenticated users (who have already passed UserSocket.connect)
    Logger.debug("User #{socket.assigns.user_id} joined global channel")
    {:ok, socket}
  end

  @impl true
  def handle_in(_event, _payload, socket) do
    {:reply, {:error, %{reason: "read_only"}}, socket}
  end

  # Channels in Pulse are read-only for clients.
  # Events are pushed from the Kafka consumer via the Endpoint.
end
