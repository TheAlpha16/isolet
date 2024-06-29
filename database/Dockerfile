FROM postgres

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=1c25d2068cdd9c3452013dc106e3be8a

COPY init.sql /docker-entrypoint-initdb.d/
EXPOSE 5432
