defmodule PulseWeb.TeamChannel do
  use PulseWeb, :channel

  @impl true
  def join("team:" <> team_id_str, _params, socket) do
    team_id = String.to_integer(team_id_str)

    if socket.assigns.team_id == team_id do
      {:ok, socket}
    else
      {:error, %{reason: "unauthorized"}}
    end
  end

  # Channels in Pulse are read-only for clients.
  # Events are pushed from the Kafka consumer via the Endpoint.
end
