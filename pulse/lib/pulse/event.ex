defmodule Pulse.Event do
  @moduledoc """
  Derives a semantic event name from a raw Notification map.

  The derivation rules (from PLAN.md):

    - Entity events:      `entity.name + "." + action`  e.g. `"instance.created"`
    - Severity events:    `"notification." + severity`   e.g. `"notification.info"`
  """

  @doc """
  Derives a semantic event name from a raw (string-keyed) notification map.

  ## Examples

      iex> Pulse.Event.derive(%{"entity" => %{"name" => "instance"}, "action" => "created"})
      "instance.created"

      iex> Pulse.Event.derive(%{"severity" => "warning"})
      "notification.warning"

  """
  @spec derive(map()) :: String.t()
  def derive(%{"entity" => %{"name" => name}, "action" => action})
      when is_binary(name) and is_binary(action) do
    name <> "." <> action
  end

  def derive(%{"severity" => severity}) when is_binary(severity) do
    "notification." <> severity
  end
end
