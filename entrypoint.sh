#!/usr/bin/env sh

/go/bin/api -dsn $GREENLIGHT_DB_DSN -smtp-host $SMTP_HOST -smtp-port $SMTP_PORT -smtp-username $MAIL_USERNAME -smtp-password $MAIL_PASSOWRD -smtp-sender "$MAIL_SENDER" -cors-trusted-origins "$CORS_ALLOWED"