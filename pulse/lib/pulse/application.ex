defmodule Pulse.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application

  @impl true
  def start(_type, _args) do
    redis_url = Application.fetch_env!(:pulse, :redis_url)
    kafka_brokers = Pulse.Kafka.Config.brokers()
    kafka_client_config = Pulse.Kafka.Config.client_config()

    kafka_client_id = Pulse.Kafka.Config.client_id()

    children = [
      PulseWeb.Telemetry,
      {DNSCluster, query: Application.get_env(:pulse, :dns_cluster_query) || :ignore},
      {Phoenix.PubSub, name: Pulse.PubSub},
      # Redis — must start before the endpoint so session validation is available
      {Redix, {redis_url, [name: :redix]}},
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

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
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
