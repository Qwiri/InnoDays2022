import time

import RPi.GPIO as GPIO
from mfrc522 import SimpleMFRC522
import requests

GPIO.setmode(GPIO.BOARD)
GPIO.setwarnings(False)


def send_request(s_id, c_id) -> requests.Response:
    request_url = "https://{0}/e/rfid/{1}/{2}/{3}".format("localhost", "0", s_id, c_id)
    print(f"POST-request to {request_url}")
    r = requests.post(request_url)
    return r


class RfidReader:
    def __init__(self):
        self.last_id = 0
        self.scanner_ids = [111, 222]
        self.reader = [SimpleMFRC522(0), SimpleMFRC522(1)]

    def get_next_id(self) -> (int, int):
        c_id = s_id = 0
        while c_id == 0 or c_id is None:
            for index, s_id in enumerate(self.scanner_ids):
                c_id = self.reader[index].read_id_no_block()
                start_time = time.time()
                while c_id is None and time.time() - start_time < 0.2:
                    c_id = self.reader[index].read_id_no_block()
                if c_id is not None:
                    print(f"Got id {c_id} from scanner {s_id}")
                    break
        return s_id, c_id


if __name__ == "__main__":
    try:
        rfid_reader = RfidReader()
        while True:
            print("RFID-Reader ready...")
            scanner_id, card_id = rfid_reader.get_next_id()
            try:
                resp = send_request(scanner_id, card_id)
                if not resp.ok:
                    print("Request failed\n")
                else:
                    print(f"Response: {resp.text}\n")
            except requests.exceptions.ConnectionError:
                print("Failed to establish a connection!\n")
            time.sleep(1)
    except KeyboardInterrupt:
        GPIO.cleanup()
        raise
