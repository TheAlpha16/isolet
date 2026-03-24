defmodule Pulse.Application do
  @moduledoc false

  use Application

  @impl true
  def start(_type, _args) do
    redis_url = Application.fetch_env!(:pulse, :redis_url)
    kafka_brokers = Pulse.Kafka.Config.brokers()
    kafka_client_config = Pulse.Kafka.Config.client_config()

    kafka_client_id = Pulse.Kafka.Config.client_id()

    pubsub_config =
      case Application.get_env(:pulse, :pubsub_adapter, :redis) do
        :redis ->
          node_name =
            System.get_env("HOSTNAME") ||
              System.get_env("POD_NAME") ||
              "pulse_#{:erlang.phash2(Node.self())}"

          [
            name: Pulse.PubSub,
            adapter: Phoenix.PubSub.Redis,
            url: redis_url,
            node_name: node_name
          ]

        :local ->
          [name: Pulse.PubSub]
      end

    children = [
      PulseWeb.Telemetry,
      {DNSCluster, query: Application.get_env(:pulse, :dns_cluster_query) || :ignore},
      # Redis — must start before PubSub and endpoint
      {Redix, {redis_url, [name: :redix]}},
      # Redis-backed PubSub for cluster-wide event distribution across pods
      {Phoenix.PubSub, pubsub_config},
      # Kafka client — must use map-based child spec for :brod_client
      %{
        id: kafka_client_id,
        start: {:brod_client, :start_link, [kafka_brokers, kafka_client_id, kafka_client_config]}
      },
      # Kafka consumer
      {Pulse.Kafka.Consumer, []},
      # Start to serve requests, typically the last entry
      PulseWeb.Endpoint
    ]

    opts = [strategy: :one_for_one, name: Pulse.Supervisor]
    Supervisor.start_link(children, opts)
  end

  # Tell Phoenix to update the endpoint configuration
  # whenever the application is updated.
  @impl true
  def config_change(changed, _new, removed) do
    PulseWeb.Endpoint.config_change(changed, removed)
    :ok
  end
end
