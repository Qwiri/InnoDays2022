#!/usr/bin/env bash

# Start client sripts
screen -UmdS rfid-scan python3 rfid.py
screen -UmdS goal-1 python3 goals.py --backend-url "http://localhost" --kid "1"