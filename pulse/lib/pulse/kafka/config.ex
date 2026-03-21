defmodule Pulse.Kafka.Config do
  @moduledoc """
  Helper for Kafka-related configuration.
  """

  @doc """
  Parses a comma-separated list of brokers into the format required by :brod.
  Example: "localhost:9092,broker2:9092" -> [{"localhost", 9092}, {"broker2", 9092}]
  """
  def brokers do
    Application.get_env(:pulse, :kafka_brokers)
    |> parse_brokers()
  end

  defp parse_brokers(brokers) when is_list(brokers), do: brokers
  defp parse_brokers(nil), do: [{"localhost", 9092}]

  defp parse_brokers(brokers_str) when is_binary(brokers_str) do
    brokers_str
    |> String.split(",")
    |> Enum.map(fn s ->
      [host, port] = String.split(s, ":")
      {host, String.to_integer(port)}
    end)
  end

  def topic, do: Application.get_env(:pulse, :kafka_topic, "notifications")
  def group_id, do: Application.get_env(:pulse, :kafka_group_id, "pulse")

  @doc """
  Returns the client config for :brod.
  """
  def client_config do
    [
      auto_start_producers: false,
      reconnect_cool_down_seconds: 5
    ]
  end
end
