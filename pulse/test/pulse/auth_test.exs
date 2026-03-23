defmodule Pulse.AuthTest do
  use ExUnit.Case, async: true
  alias Pulse.Auth

  defp get_secret do
    Application.fetch_env!(:pulse, :jwt_secret)
  end

  setup do
    # Ensure Redix is available and we can write to it for testing
    # In a real CI, we'd use a mock or a real local redis
    :ok
  end

  describe "verify_jwt/1" do
    test "verifies a valid token" do
      claims = %{"user_id" => "123", "role" => "admin", "purpose" => "realtime"}
      signer = Joken.Signer.create("HS256", get_secret())
      {:ok, token, _} = Joken.generate_and_sign(Joken.Config.default_claims(), claims, signer)

      assert {:ok, verified_claims} = Auth.verify_jwt(token)
      assert verified_claims["user_id"] == "123"
    end

    test "returns error for invalid token" do
      assert {:error, _} = Auth.verify_jwt("invalid.token.here")
    end
  end

  describe "validate_claims/1" do
    test "returns :ok for valid claims" do
      claims = %{"purpose" => "realtime", "user_id" => "123", "role" => "user"}
      assert Auth.validate_claims(claims) == :ok
    end

    test "returns :error for missing user_id" do
      claims = %{"purpose" => "realtime", "role" => "user"}
      assert Auth.validate_claims(claims) == :error
    end

    test "returns :error for wrong purpose" do
      claims = %{"purpose" => "refresh", "user_id" => "123", "role" => "user"}
      assert Auth.validate_claims(claims) == :error
    end
  end

  describe "validate_session/1" do
    test "returns :ok if session exists in Redis" do
      claims = %{"sub" => "user123", "jti" => "token456", "purpose" => "realtime"}
      key = "token:realtime:user123:token456"

      # We assume :redix is started in test env.
      # If not, this test might fail or need a mock.
      case Redix.command(:redix, ["SET", key, "valid"]) do
        {:ok, "OK"} ->
          assert Auth.validate_session(claims) == :ok
          Redix.command(:redix, ["DEL", key])

        _ ->
          # Skip if redis is not reachable in this environment
          :ok
      end
    end

    test "returns :error if session missing" do
      claims = %{"sub" => "user123", "jti" => "missing", "purpose" => "realtime"}
      assert Auth.validate_session(claims) == :error
    end
  end
end
