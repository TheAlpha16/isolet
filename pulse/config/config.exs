# This file is responsible for configuring your application
# and its dependencies with the aid of the Config module.
#
# This configuration file is loaded before any dependency and
# is restricted to this project.

# General application configuration
import Config

config :pulse,
  generators: [timestamp_type: :utc_datetime],
  # Kafka defaults
  kafka_brokers: [{"localhost", 9092}],
  kafka_topic: "herald.notifications",
  kafka_group_id: "pulse",
  kafka_client_id: :pulse_kafka_client,
  # Redis defaults
  redis_url: "redis://localhost:6379"

# Configure the endpoint
config :pulse, PulseWeb.Endpoint,
  url: [host: "localhost"],
  adapter: Bandit.PhoenixAdapter,
  render_errors: [
    formats: [json: PulseWeb.ErrorJSON],
    layout: false
  ],
  pubsub_server: Pulse.PubSub,
  live_view: [signing_salt: "3wrIv7Qc"]

# Configure Elixir's Logger
config :logger, :default_formatter,
  format: "$time $metadata[$level] $message\n",
  metadata: [:request_id]

# Use Jason for JSON parsing in Phoenix
config :phoenix, :json_library, Jason

# Import environment specific config. This must remain at the bottom
# of this file so it overrides the configuration defined above.
import_config "#{config_env()}.exs"
