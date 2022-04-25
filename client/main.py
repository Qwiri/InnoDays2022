import RPi.GPIO as GPIO
from time import sleep
import sys
from mfrc522 import SimpleMFRC522
import requests

GPIO.setmode(GPIO.BCM)
GPIO.setwarnings(False)

if __name__ == "__main__":
    reader = SimpleMFRC522()

    try:
        last_id = 0
        while True:
            print("RFID-Reader ready...")
            id, text = reader.read()
            if last_id != id:
                print("ID: %s\nText: %s" % (id, text))
                resp = requests.post(f"{0}/e/rfid/{1}/{2}/{3}".format("localhost", "0", "0", id))
                if not resp.ok:
                    print("Request failed")
            else:
                print("Same ID as last time")
            sleep(5)
    except KeyboardInterrupt:
        GPIO.cleanup()
        raise

