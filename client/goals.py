import RPi.GPIO as GPIO
import requests
import argparse
import asyncio

# Argparsing
parser = argparse.ArgumentParser(
    description='Setting up goal counter for kicker table')
parser.add_argument('-backend-url', type=str,
                    default='https://localhost', help='URL of Pinguin-backend')
parser.add_argument('-kid', type=str, default='0', help='ID of kicker table')


BACKEND_URL = "https://localhost"
KICKER_ID = 0

GOAL_ONE_PIN = 10
GOAL_TWO_PIN = 12

GOAL_ONE_ID = 0
GOAL_TWO_ID = 1

PINS = [GOAL_ONE_PIN, GOAL_TWO_PIN]

DEBOUNCING = 500


def button_callback(channel):
    if channel == GOAL_ONE_PIN:
        goal_id = GOAL_ONE_ID
    else:
        goal_id = GOAL_TWO_ID
    url = f"{BACKEND_URL}/e/tor/{KICKER_ID}/{goal_id}"

    resp = requests.post(url)

    if not resp.ok:
        print("Error: " + resp.status_code)
        print(resp.text)


def setup():
    GPIO.setmode(GPIO.BOARD)
    for pin in PINS:
        GPIO.setup(pin, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
        GPIO.add_event_detect(
            pin, GPIO.RISING, callback=button_callback, bouncetime=DEBOUNCING)

if __name__ == "__main__":
    args = parser.parse_args()
    BACKEND_URL = args.backend_url
    KICKER_ID = args.kid
    setup()

    loop = asyncio.new_event_loop()
    loop.run_forever()
