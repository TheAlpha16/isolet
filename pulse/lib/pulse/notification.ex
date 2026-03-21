defmodule Pulse.Notification do
  @moduledoc """
  Represents a generic Notification fact consumed from Kafka.
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
  Decodes a JSON map into a `Pulse.Notification` struct.
  Performs only minimal validation required for system stability.
  """
  @spec from_map(map()) :: {:ok, t()} | {:error, term()}
  def from_map(%{"type" => "Notification", "id" => id, "at" => at} = map)
      when is_binary(id) and id != "" do
    team_ids = map["team_ids"] || []

    if is_list(team_ids) and Enum.all?(team_ids, &is_integer/1) do
      notification = %__MODULE__{
        type: "Notification",
        id: id,
        at: at,
        entity: map["entity"],
        action: map["action"],
        message: map["message"],
        severity: normalize_severity(map["severity"], map["message"]),
        team_ids: team_ids
      }

      {:ok, notification}
    else
      {:error, :invalid_team_ids}
    end
  end

  def from_map(%{"type" => "Notification"}), do: {:error, :invalid_id}
  def from_map(_), do: {:error, :invalid_notification_type}

  defp normalize_severity(nil, message) when not is_nil(message), do: "info"
  defp normalize_severity("", message) when not is_nil(message), do: "info"
  defp normalize_severity(severity, _), do: severity
end
