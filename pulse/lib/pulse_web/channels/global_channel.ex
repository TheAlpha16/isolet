defmodule PulseWeb.GlobalChannel do
  use PulseWeb, :channel

  @impl true
  def join("global", _params, socket) do
    # Open to all authenticated users (who have already passed UserSocket.connect)
    {:ok, socket}
  end

  # Channels in Pulse are read-only for clients.
  # Events are pushed from the Kafka consumer via the Endpoint.
end
