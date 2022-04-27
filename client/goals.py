import RPi.GPIO as GPIO
import requests
import argparse
import asyncio

# Argparsing
parser = argparse.ArgumentParser(
    description='Setting up goal counter for kicker table')
parser.add_argument('--backend-url', type=str,
                    default='http://localhost', help='URL of Penguin')
parser.add_argument('--kid', type=str, default='1', help='ID of kicker table')


BACKEND_URL = "http://localhost"
KICKER_ID = 0

GOAL_BLACK_LASER = 16
GOAL_WHITE_LASER = 18

LASERS = [GOAL_BLACK_LASER, GOAL_WHITE_LASER]

GOAL_BLACK_PIN = 10
GOAL_WHITE_PIN = 12

GOAL_BLACK_ID = 1
GOAL_WHITE_ID = 2

PINS = [GOAL_BLACK_PIN, GOAL_WHITE_PIN]

DEBOUNCING = 500


def button_callback(channel):
    if channel == GOAL_BLACK_PIN:
        goal_id = GOAL_BLACK_ID
    else:
        goal_id = GOAL_WHITE_ID
    url = f"{BACKEND_URL}/e/tor/{KICKER_ID}/{goal_id}"

    resp = requests.post(url)
    print(url)

    if not resp.ok:
        print("Error: {resp.status_code}")
        print(resp.text)


def setup():
    GPIO.setmode(GPIO.BOARD)
    GPIO.setwarnings(False)

    for laser in LASERS:
        GPIO.setup(laser, GPIO.OUT)
        GPIO.output(laser, GPIO.HIGH)

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
    try:
        loop.run_forever()
    except KeyboardInterrupt:
        GPIO.cleanup()
