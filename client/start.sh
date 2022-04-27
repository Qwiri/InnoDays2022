#!/usr/bin/env bash

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '#' | sed 's/\r$//' | awk '/=/ {print $1}' )
fi

# Start client scripts
screen -UmdS rfid-scan python3 rfid.py
screen -UmdS goal-1 python3 goals.py --backend-url ${BACKEND_URL} --kid ${KICKER_ID}

# Start chrome
chromium-browser --start-fullscreen ${BACKEND_FRONTEND_URL} &