defmodule Pulse.EventTest do
  use ExUnit.Case, async: true
  alias Pulse.Event
  alias Pulse.Notification

  describe "derive/1" do
    test "returns notification.<severity> for any severity" do
      notification = %Notification{severity: "critical", action: "shipped", entity: "order"}
      assert Event.derive(notification) == "notification.critical"

      notification = %Notification{severity: "debug", action: "ping", entity: "system"}
      assert Event.derive(notification) == "notification.debug"
    end

    test "returns <entity>.<action> when severity is missing" do
      notification = %Notification{severity: nil, action: "delivered", entity: "parcel"}
      assert Event.derive(notification) == "parcel.delivered"

      notification = %Notification{severity: "", action: "failed", entity: "payment"}
      assert Event.derive(notification) == "payment.failed"
    end

    test "handles nested entity maps" do
      notification = %Notification{entity: %{"name" => "sensor_a"}, action: "triggered"}
      assert Event.derive(notification) == "sensor_a.triggered"
    end

    test "handles missing action or entity with defaults" do
      notification = %Notification{entity: "device_1"}
      assert Event.derive(notification) == "device_1.unknown"

      notification = %Notification{action: "scan"}
      assert Event.derive(notification) == "unknown.scan"

      notification = %Notification{}
      assert Event.derive(notification) == "unknown.unknown"
    end
  end
end
