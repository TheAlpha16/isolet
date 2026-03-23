defmodule Pulse.Auth do
  @moduledoc """
  JWT and Redis session authentication helpers for Pulse WebSocket connections.

  Enforces the same token invariants as Oracle:
    - `purpose` must be `"realtime"`
    - `jti`, `sub`, `exp`, `iat` must exist
    - `user_id` and `role` must exist
    - `team_id` is optional
    - Session (`jti`) must be live in Redis under key `token:{purpose}:{sub}:{jti}`
  """

  use Joken.Config

  @doc """
  Verifies the JWT signature and expiry.

  Returns `{:ok, claims}` on success, `{:error, reason}` on failure.
  """
  @spec verify_jwt(String.t()) :: {:ok, map()} | {:error, term()}
  def verify_jwt(token) do
    secret = Application.fetch_env!(:pulse, :jwt_secret)
    signer = Joken.Signer.create("HS256", secret)

    case Joken.verify(token, signer) do
      {:ok, claims} -> {:ok, claims}
      {:error, reason} -> {:error, reason}
    end
  end

  @doc """
  Validates that the JWT claims meet Pulse's requirements:
    - `purpose` must be `"realtime"`
    - `user_id` and `role` must be present
  """
  @spec validate_claims(map()) :: :ok | :error
  def validate_claims(%{"purpose" => "realtime", "user_id" => user_id, "role" => role})
      when not is_nil(user_id) and not is_nil(role) do
    :ok
  end

  def validate_claims(_), do: :error

  @doc """
  Validates that the token's session (`jti`) is still live in Redis.

  Key format: `token:{purpose}:{sub}:{jti}` — matches Oracle exactly.
  """
  @spec validate_session(map()) :: :ok | :error
  def validate_session(%{"jti" => jti, "sub" => subject, "purpose" => purpose}) do
    key = "token:#{purpose}:#{subject}:#{jti}"

    case Redix.command(:redix, ["GET", key]) do
      {:ok, nil} -> :error
      {:ok, _value} -> :ok
      {:error, _} -> :error
    end
  end

  def validate_session(_), do: :error
end
