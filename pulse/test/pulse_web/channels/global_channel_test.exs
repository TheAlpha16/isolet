defmodule PulseWeb.GlobalChannelTest do
  use PulseWeb.ChannelCase, async: true
  alias PulseWeb.GlobalChannel

  setup do
    {:ok, _, socket} =
      PulseWeb.UserSocket
      |> socket("user_id", %{user_id: "u1", role: "user"})
      |> subscribe_and_join(GlobalChannel, "global")

    %{socket: socket}
  end

  test "any authenticated user can join", %{socket: socket} do
    assert socket.topic == "global"
  end

  test "broadcasts are received by the client", %{socket: _socket} do
    payload = %{"event" => "system:broadcast", "data" => "test"}
    PulseWeb.Endpoint.broadcast!("global", "notification", payload)
    assert_push "notification", ^payload
  end

  test "client cannot push to the channel", %{socket: socket} do
    ref = push(socket, "ping", %{})
    assert_reply ref, :error, %{reason: "read_only"}
  end
end
