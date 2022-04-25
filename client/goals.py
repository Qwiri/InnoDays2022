import RPi.GPIO as GPIO
import requests
import argparse

#Argparsing
parser = argparse.ArgumentParser(description='Setting up goal counter for kicker table')
parser.add_argument('-host', type=str, default='localhost')
parser.add_argument('-kid', type=str, default='0')


BACKEND_URL = "localhost"
KICKER_ID = 0

GOAL_ONE_PIN = 10
GOAL_TWO_PIN = 12

GOAL_ONE_ID = 0
GOAL_TWO_ID = 1

PINS = [GOAL_ONE_PIN, GOAL_TWO_PIN]

DEBOUNCING = 1000

def button_callback(channel):
    if channel == GOAL_ONE_PIN:
        goal_id = GOAL_ONE_ID
    else:
        goal_id = GOAL_TWO_ID
    print(f"{BACKEND_URL}/e/tor/{KICKER_ID}/{goal_id}")

def setup():
    GPIO.setmode(GPIO.BOARD)
    for pin in PINS:
        GPIO.setup(pin, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
        GPIO.add_event_detect(pin,GPIO.RISING,callback=button_callback, bouncetime=DEBOUNCING)


args = parser.parse_args()
BACKEND_URL = args.host
KICKER_ID = args.kid
setup()
# perma loop; optimizable I guess
while True:
    continue