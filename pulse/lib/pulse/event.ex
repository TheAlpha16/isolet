defmodule Pulse.Event do
  @moduledoc """
  Derives a semantic event name from a Notification.
  """

  alias Pulse.Notification

  @doc """
  Derives a semantic event name from a `Pulse.Notification` struct.
  """
  @spec derive(Notification.t()) :: String.t()
  def derive(%Notification{severity: severity}) when is_binary(severity) and severity != "" do
    "notification." <> severity
  end

  def derive(%Notification{entity: entity, action: action}) do
    entity_name = get_entity_name(entity)
    action_name = action || "unknown"
    entity_name <> "." <> action_name
  end

  defp get_entity_name(%{"name" => name}) when is_binary(name), do: name
  defp get_entity_name(name) when is_binary(name), do: name
  defp get_entity_name(_), do: "unknown"
end
