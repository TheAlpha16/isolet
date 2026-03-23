defmodule PulseWeb.TeamChannelTest do
  use PulseWeb.ChannelCase, async: true
  alias PulseWeb.TeamChannel

  setup do
    {:ok, _, socket} =
      PulseWeb.UserSocket
      |> socket("user_socket:10", %{user_id: 10, team_id: 9, role: "user"})
      |> subscribe_and_join(TeamChannel, "team:9")

    %{socket: socket}
  end

  test "join returns ok when team_id matches", %{socket: socket} do
    assert socket.topic == "team:9"
  end

  test "join returns error when team_id doesn't match" do
    assert {:error, %{reason: "unauthorized"}} =
             PulseWeb.UserSocket
             |> socket("user_socket:10", %{user_id: 10, team_id: 8, role: "user"})
             |> subscribe_and_join(TeamChannel, "team:9")
  end

  test "broadcasts are received by the client", %{socket: _socket} do
    payload = %{"event" => "user:create", "data" => "test"}
    PulseWeb.Endpoint.broadcast!("team:9", "notification", payload)
    assert_push "notification", ^payload
  end

  test "client cannot push to the channel", %{socket: socket} do
    ref = push(socket, "ping", %{})
    assert_reply ref, :error, %{reason: "read_only"}
  end
end
