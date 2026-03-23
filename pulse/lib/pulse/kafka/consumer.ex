defmodule Pulse.Kafka.Consumer do
  @moduledoc """
  Kafka consumer module for Pulse.

  Uses `:brod_group_subscriber` to consume events from the `notifications` topic.
  Each event is decoded, transformed (event name derived), and broadcast to the
  appropriate channels (team-scoped or global).
  """

  alias Pulse.Kafka.Config
  alias Pulse.Event

  require Logger

  # brod_group_subscriber callbacks
  @behaviour :brod_group_subscriber

  def child_spec(opts) do
    %{
      id: __MODULE__,
      start: {__MODULE__, :start_link, [opts]}
    }
  end

  @doc """
  Starts the consumer as part of the supervision tree.
  """
  def start_link(_opts \\ []) do
    client_id = Config.client_id()
    topic = Config.topic()
    group_id = Config.group_id()

    group_config = [
      offset_commit_policy: :commit_to_kafka_v2,
      offset_commit_interval_seconds: 5
    ]

    # Consumer config
    consumer_config = [
      begin_offset: :earliest
    ]

    :brod_group_subscriber.start_link(
      client_id,
      group_id,
      [topic],
      group_config,
      consumer_config,
      __MODULE__,
      _cb_init_data = %{}
    )
  end

  # --- brod_group_subscriber callbacks ---

  @impl true
  def init(_group_id, _cb_init_data) do
    {:ok, %{}}
  end

  @impl true
  def handle_message(topic, partition, message, state) do
    Logger.debug("Received message from Kafka: topic=#{topic}, partition=#{partition}")

    case process_message(message) do
      :ok ->
        {:ok, :ack, state}

      {:error, reason} ->
        Logger.error(
          "Failed to process Kafka message on #{topic}:#{partition}. Reason: #{inspect(reason)}"
        )

        # We ack even on error to avoid blocking the partition infinitely.
        {:ok, :ack, state}
    end
  end

  # --- Internal Helpers ---

  # `:brod_group_subscriber` may deliver either a single kafka message or a kafka message set.
  defp process_message({:kafka_message_set, _topic, _partition, _offset, messages})
       when is_list(messages) do
    Enum.reduce_while(messages, :ok, fn message, :ok ->
      case decode_and_broadcast(message) do
        :ok -> {:cont, :ok}
        {:error, reason} -> {:halt, {:error, reason}}
      end
    end)
  end

  defp process_message({:kafka_message_set, _topic, _partition, _offset, {:incomplete_batch, _}}) do
    :ok
  end

  defp process_message(message), do: decode_and_broadcast(message)

  # Newer brod versions include headers in the kafka message tuple.
  defp decode_and_broadcast({:kafka_message, _offset, _key, value, _ts_type, _ts, _headers}) do
    decode_payload_and_broadcast(value)
  end

  # Keep compatibility with older brod tuple shape.
  defp decode_and_broadcast({:kafka_message, _offset, _key, value, _ts_type, _ts}) do
    decode_payload_and_broadcast(value)
  end

  defp decode_payload_and_broadcast(value) do
    with {:ok, raw} <- Jason.decode(value) do
      Logger.debug("Decoded Kafka payload: #{inspect(raw)}")

      case Pulse.Notification.from_map(raw) do
        {:ok, notification} ->
          event_name = Event.derive(notification)
          payload = Map.put(raw, "event", event_name)

          broadcast(notification, payload)

        {:error, reason} ->
          {:error, reason}
      end
    end
  end

  defp broadcast(%Pulse.Notification{team_ids: team_ids}, payload)
       when is_list(team_ids) and team_ids != [] do
    Enum.each(team_ids, fn team_id ->
      Logger.debug("Broadcasting to team:#{team_id}, event=#{payload["event"]}")
      PulseWeb.Endpoint.broadcast!("team:#{team_id}", "notification", payload)
    end)

    :ok
  end

  defp broadcast(_notification, payload) do
    # Fallback to global if no team_ids or team_ids is empty
    Logger.debug("Broadcasting to global, event=#{payload["event"]}")
    PulseWeb.Endpoint.broadcast!("global", "notification", payload)
    :ok
  end
end
