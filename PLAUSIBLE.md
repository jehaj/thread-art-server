I have been trying to setup [plausible.io](https://plausible.io/) for self-host with podman - quite unsuccessfully. Until now. It took more than 4 hours, so this is for future reference. Do not waste more time :-)

The `docker-compose.yaml` should look like
```
version: "3.3"
services:
  mail:
    image: bytemark/smtp
    restart: always

  plausible_db:
    # supported versions are 12, 13, and 14
    image: postgres:14-alpine
    restart: always
    volumes:
      - db-data:/var/lib/postgresql/data
    environment:
      - POSTGRES_PASSWORD=postgres

  plausible_events_db:
    image: clickhouse/clickhouse-server:22.6-alpine
    restart: always
    volumes:
      - event-data:/var/lib/clickhouse
      - ./clickhouse/clickhouse-config.xml:/etc/clickhouse-server/config.d/logging.xml:Z
      - ./clickhouse/clickhouse-user-config.xml:/etc/clickhouse-server/users.d/logging.xml:Z
    ulimits:
      nofile:
        soft: 262144
        hard: 262144

  plausible:
    image: plausible/analytics:latest
    restart: always
    command: sh -c "sleep 10 && /entrypoint.sh db createdb && /entrypoint.sh db migrate && /entrypoint.sh run"
    depends_on:
      - plausible_db
      - plausible_events_db
      - mail
    ports:
      - 8000:8000
    env_file:
      - plausible-conf.env

  caddy:
    image: caddy:2.6
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
      - "443:443/udp"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile:Z
      - caddy_data:/data
      - caddy_config:/config

volumes:
  db-data:
    driver: local
  event-data:
    driver: local
  geoip:
    driver: local
```

The modification consists of changing :ro (read-only) to :Z. It has something to do with SELinux, which fedora uses. Not making this change will result in permissions denied errors. More information may be found at [blog.christophersmart.com](https://blog.christophersmart.com/2021/01/31/podman-volumes-and-selinux/). `caddy` has been added as an reverse_proxy to enable SSL down the road.

The `Caddyfile` should look like
```
:80

reverse_proxy {
  to plausible:8000
}
```

I believe the Caddyfile could also look like
```
stats.domain.com {
  reverse_proxy {
    to plausible:8000
  }
}
```

If `plausible` exists in another `.yaml` than `caddy`, then they need to be on the same network. Consult [baeldung.com](https://www.baeldung.com/ops/docker-compose-communication) for more information.

Good luck!
