defmodule Pulse.NotificationTest do
  use ExUnit.Case, async: true
  alias Pulse.Notification

  @valid_base %{
    "type" => "Notification",
    "id" => "note123",
    "at" => "2026-03-21T23:13:18Z"
  }

  describe "from_map/1" do
    test "decodes valid notification with arbitrary message and entity" do
      map =
        Map.merge(@valid_base, %{
          "message" => "any message",
          "severity" => "emergency",
          "entity" => "custom_system",
          "action" => "restarted"
        })

      assert {:ok, %Notification{message: "any message", severity: "emergency"}} =
               Notification.from_map(map)
    end

    test "handles nested entity maps without validation" do
      entity = %{"id" => 999, "name" => "new_entity_type", "extra" => "data"}
      map = Map.merge(@valid_base, %{"entity" => entity, "action" => "upgraded"})

      assert {:ok, %Notification{entity: ^entity, action: "upgraded"}} =
               Notification.from_map(map)
    end

    test "defaults severity to info if message present but severity empty" do
      map = Map.merge(@valid_base, %{"message" => "hello", "severity" => ""})
      assert {:ok, %Notification{severity: "info"}} = Notification.from_map(map)
    end

    test "rejects missing id" do
      map = Map.delete(@valid_base, "id")
      assert {:error, :invalid_id} = Notification.from_map(map)
    end

    test "rejects invalid team ids format" do
      map = Map.merge(@valid_base, %{"message" => "hi", "team_ids" => "not_a_list"})
      assert {:error, :invalid_team_ids} = Notification.from_map(map)
    end

    test "rejects non-integer team ids" do
      map = Map.merge(@valid_base, %{"message" => "hi", "team_ids" => ["1"]})
      assert {:error, :invalid_team_ids} = Notification.from_map(map)
    end
  end
end
