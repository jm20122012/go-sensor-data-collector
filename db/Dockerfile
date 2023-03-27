FROM postgres:16.3-alpine3.20

EXPOSE 5432

ARG DB_USER
ARG DB_PASSWORD
ARG DB_NAME

ENV POSTGRES_USER=$POSTGRES_USER
ENV POSTGRES_PASSWORD=$POSTGRES_PASSWORD
ENV POSTGRES_DB=$POSTGRES_DB

WORKDIR /db

# Copy the database schema and seed data
COPY ./migrations/ ./migrations
COPY ./seed.sql ./seed.sql
COPY ./custom-entrypoint.sh /usr/local/bin/custom-entrypoint.sh

# Grant execute permissions to the custom entrypoint script
RUN chmod 775 /usr/local/bin/custom-entrypoint.sh

# Install the Atlas CLI and other dependencies
RUN apk update && apk add --no-cache curl bash coreutils

# Install Atlas
RUN yes | curl -sSL https://atlasgo.sh | sh

# Set the default user to postgres
USER postgres

ENTRYPOINT [ "custom-entrypoint.sh" ]