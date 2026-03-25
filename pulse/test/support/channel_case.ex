defmodule PulseWeb.ChannelCase do
  @moduledoc """
  This module defines the test case to be used by
  channel tests.
  """

  use ExUnit.CaseTemplate

  using do
    quote do
      import Phoenix.ChannelTest
      import PulseWeb.ChannelCase

      @endpoint PulseWeb.Endpoint
    end
  end

  setup _tags do
    :ok
  end
end
