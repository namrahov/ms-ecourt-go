#!/bin/bash

if [ -f /vault/secrets/secrets ]
then
  echo "Loading values from secrets file"
  source /vault/secrets/secrets

  sudo rm /vault/secrets/secrets
fi

sleep 5

exec ./ms-ecourt
