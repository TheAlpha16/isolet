FROM elixir:1.16-slim AS build

ARG EXTRA_CA_CERT=""

# apt uses HTTP — no TLS cert needed; install ca-certificates so update-ca-certificates is available
RUN apt-get update -y && apt-get install -y ca-certificates build-essential git \
    && apt-get clean && rm -f /var/lib/apt/lists/*_*

# Inject extra CA cert before any HTTPS calls (mix deps.get, mix local.hex, etc.)
RUN if [ -n "$EXTRA_CA_CERT" ]; then \
        echo "$EXTRA_CA_CERT" > /usr/local/share/ca-certificates/extra.crt && \
        update-ca-certificates; \
    fi

WORKDIR /app

RUN mix local.hex --force && mix local.rebar --force

ENV MIX_ENV=prod

COPY mix.exs mix.lock ./
RUN mix deps.get --only $MIX_ENV

RUN mkdir config
COPY config/config.exs config/prod.exs config/
RUN mix deps.compile

COPY lib lib

RUN mix compile

COPY config/runtime.exs config/

RUN mix release

FROM debian:bookworm-slim

RUN apt-get update -y && \
    apt-get install -y libstdc++6 openssl libncurses5 locales ca-certificates \
    && apt-get clean && rm -f /var/lib/apt/lists/*_*

RUN sed -i '/en_US.UTF-8/s/^# //g' /etc/locale.gen && locale-gen

ENV LANG=en_US.UTF-8
ENV LANGUAGE=en_US:en
ENV LC_ALL=en_US.UTF-8

WORKDIR /app
RUN chown nobody /app

ENV MIX_ENV=prod

COPY --from=build --chown=nobody:root /app/_build/prod/rel/pulse ./

USER nobody

EXPOSE 4000

ENV PHX_SERVER=true

CMD ["/app/bin/pulse", "start"]
