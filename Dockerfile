FROM common-env:3.0

WORKDIR /app

COPY config/config.toml /app/
COPY bin/* /app/
