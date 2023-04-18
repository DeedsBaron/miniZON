#!/bin/bash

MESSAGE=$(cat $1)

curl -X POST \
  -H 'Content-Type: application/json' \
  -d "{\"chat_id\": \"$TELEGRAM_CHAT_ID\", \"text\": \"$MESSAGE\"}" \
  https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/sendMessage