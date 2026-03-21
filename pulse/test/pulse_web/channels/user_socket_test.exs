defmodule PulseWeb.UserSocketTest do
  use PulseWeb.ChannelCase, async: true
  alias PulseWeb.UserSocket

  defp get_secret do
    Application.fetch_env!(:pulse, :jwt_secret)
  end

  defp generate_token(claims) do
    signer = Joken.Signer.create("HS256", get_secret())
    {:ok, token, _} = Joken.generate_and_sign(Joken.Config.default_claims(), claims, signer)
    token
  end

  describe "connect/3" do
    test "authenticates with valid token and redis session" do
      claims = %{
        "user_id" => "u1",
        "role" => "admin",
        "purpose" => "auth",
        "sub" => "u1",
        "jti" => "t1"
      }

      token = generate_token(claims)
      key = "token:auth:u1:t1"

      # Setup Redis session
      case Redix.command(:redix, ["SET", key, "valid"]) do
        {:ok, "OK"} ->
          assert {:ok, socket} = connect(UserSocket, %{"token" => token})
          assert socket.assigns.user_id == "u1"
          assert socket.assigns.role == "admin"
          Redix.command(:redix, ["DEL", key])

        _ ->
          # Skip if redis is not reachable
          :ok
      end
    end

    test "rejects invalid tokens" do
      assert :error = connect(UserSocket, %{"token" => "bad_token"})
    end

    test "rejects missing tokens" do
      assert :error = connect(UserSocket, %{})
    end
  end

  describe "id/1" do
    test "returns correct socket id" do
      socket = %Phoenix.Socket{assigns: %{user_id: "user123"}}
      assert UserSocket.id(socket) == "user_socket:user123"
    end
  end
end
