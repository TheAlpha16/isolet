defmodule PulseWeb.TeamChannelTest do
  use PulseWeb.ChannelCase, async: true
  alias PulseWeb.TeamChannel

  setup do
    {:ok, _, socket} =
      PulseWeb.UserSocket
      |> socket("user_id", %{user_id: "u1", team_id: "team_a", role: "user"})
      |> subscribe_and_join(TeamChannel, "team:team_a")

    %{socket: socket}
  end

  test "join returns ok when team_id matches", %{socket: socket} do
    assert socket.topic == "team:team_a"
  end

  test "join returns error when team_id doesn't match" do
    assert {:error, %{reason: "unauthorized"}} =
             PulseWeb.UserSocket
             |> socket("user_id", %{user_id: "u1", team_id: "team_b", role: "user"})
             |> subscribe_and_join(TeamChannel, "team:team_a")
  end

  test "broadcasts are received by the client", %{socket: _socket} do
    payload = %{"event" => "user:create", "data" => "test"}
    PulseWeb.Endpoint.broadcast!("team:team_a", "notification", payload)
    assert_push "notification", ^payload
  end

  test "client cannot push to the channel", %{socket: socket} do
    ref = push(socket, "ping", %{})
    assert_reply ref, :error, %{reason: "read_only"}
  end
end
