defmodule Pulse.Kafka.Config do
  @moduledoc """
  Helper for Kafka-related configuration.
  All values are pulled from the Application environment (config).
  """

  def brokers, do: Application.fetch_env!(:pulse, :kafka_brokers)
  def topic, do: Application.fetch_env!(:pulse, :kafka_topic)
  def group_id, do: Application.fetch_env!(:pulse, :kafka_group_id)
  def client_id, do: Application.fetch_env!(:pulse, :kafka_client_id)

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
