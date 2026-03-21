defmodule PulseWeb.HealthController do
  use PulseWeb, :controller

  def index(conn, _params) do
    json(conn, %{status: "still alive!"})
  end
end
