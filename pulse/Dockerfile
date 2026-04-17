FROM elixir:1.15-alpine AS build

ARG EXTRA_CA_CERT=""

RUN apk add --no-cache build-base git ca-certificates && \
    if [ -n "$EXTRA_CA_CERT" ]; then \
        echo "$EXTRA_CA_CERT" > /usr/local/share/ca-certificates/extra.crt && \
        update-ca-certificates; \
    fi

WORKDIR /app

ENV MIX_ENV=prod

RUN mix local.hex --force && mix local.rebar --force

COPY mix.exs mix.lock ./
RUN mix deps.get --only prod

COPY config config
RUN mix deps.compile

COPY lib lib
COPY priv priv

RUN mix compile
RUN mix release

FROM alpine:3.18 AS runner

RUN apk add --no-cache openssl libgcc libstdc++ ncurses-libs

WORKDIR /app

RUN addgroup -S -g 1001 elixir && adduser -S -u 1001 -G elixir elixir
RUN chown elixir:elixir /app

COPY --from=build --chown=elixir:elixir /app/_build/prod/rel/pulse ./

USER elixir

EXPOSE 4000

ENV PHX_SERVER=true

ENTRYPOINT ["/app/bin/pulse"]
CMD ["start"]
