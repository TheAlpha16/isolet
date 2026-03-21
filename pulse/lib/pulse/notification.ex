defmodule Pulse.Notification do
  @moduledoc """
  Represents a canonical Notification fact consumed from Kafka.

  These are emitted by Herald and carry domain events to Pulse
  for realtime delivery to connected clients.
  """

  @type t :: %__MODULE__{
          type: String.t(),
          id: String.t(),
          at: String.t(),
          entity: map() | nil,
          action: String.t() | nil,
          message: String.t() | nil,
          severity: String.t() | nil,
          team_ids: [integer()]
        }

  defstruct [:type, :id, :at, :entity, :action, :message, :severity, team_ids: []]

  @doc """
  Decodes a raw JSON-decoded map into a `Pulse.Notification` struct.

  Returns `{:ok, notification}` or `{:error, reason}`.
  """
  @spec from_map(map()) :: {:ok, t()} | {:error, term()}
  def from_map(%{"type" => type, "id" => id, "at" => at} = map) do
    notification = %__MODULE__{
      type: type,
      id: id,
      at: at,
      entity: map["entity"],
      action: map["action"],
      message: map["message"],
      severity: map["severity"],
      team_ids: map["team_ids"] || []
    }

    {:ok, notification}
  end

  def from_map(_), do: {:error, :invalid_notification}
end
