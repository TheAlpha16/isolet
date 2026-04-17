FROM elixir:1.16-slim AS builder

ARG EXTRA_CA_CERT=""

# install build dependencies (apt is HTTP, so no cert needed for this step)
RUN apt-get update \
  && apt-get install -y --no-install-recommends build-essential git ca-certificates \
  && rm -rf /var/lib/apt/lists/*

# inject extra CA cert before any HTTPS calls (mix deps.get, hex, etc.)
RUN if [ -n "$EXTRA_CA_CERT" ]; then \
        echo "$EXTRA_CA_CERT" > /usr/local/share/ca-certificates/extra.crt && \
        update-ca-certificates; \
    fi

# prepare build dir
WORKDIR /app

# Under QEMU (cross-platform builds), Erlang's multi-scheduler can cause
# module-loading race conditions. +S 1 forces single-threaded compilation.
ENV ELIXIR_ERL_OPTIONS="+S 1"

# rebar3 is needed to compile erlang deps (brod); hex is already installed in the base image
RUN mix local.rebar --force

# set build ENV
ENV MIX_ENV="prod"

# install mix dependencies
COPY mix.exs mix.lock ./
RUN mix deps.get --only $MIX_ENV
RUN mkdir config

# copy compile-time config files before we compile dependencies
# to ensure any relevant config change will trigger the dependencies
# to be re-compiled.
COPY config/config.exs config/${MIX_ENV}.exs config/
RUN mix deps.compile

COPY lib lib

# Compile the release
RUN mix compile

# Changes to config/runtime.exs don't require recompiling the code
COPY config/runtime.exs config/

COPY rel rel
RUN mix release

# start a new build stage so that the final image will only contain
# the compiled release and other runtime necessities
FROM debian:bookworm-slim AS final

RUN apt-get update \
  && apt-get install -y --no-install-recommends libstdc++6 openssl libncurses5 locales ca-certificates \
  && rm -rf /var/lib/apt/lists/*

# Set the locale
RUN sed -i '/en_US.UTF-8/s/^# //g' /etc/locale.gen \
  && locale-gen

ENV LANG=en_US.UTF-8
ENV LANGUAGE=en_US:en
ENV LC_ALL=en_US.UTF-8

WORKDIR "/app"
RUN chown nobody /app

# set runner ENV
ENV MIX_ENV="prod"

# Only copy the final release from the build stage
COPY --from=builder --chown=nobody:root /app/_build/prod/rel/pulse ./

USER nobody

CMD ["/app/bin/server"]
