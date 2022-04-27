import os
import time
from mfrc522 import SimpleMFRC522
import requests
from dotenv import load_dotenv


def send_request(u, k_id, s_id, c_id) -> requests.Response:
    request_url = "{0}/e/rfid/{1}/{2}/{3}".format(u, k_id, s_id, c_id)
    print(f"POST-request to {request_url}")
    r = requests.post(request_url)
    return r


class RfidReader:
    def __init__(self):
        self.reader = [SimpleMFRC522(0), SimpleMFRC522(1)]
        self.scanner_ids = [os.getenv("SCANNER_BLACK"), os.getenv("SCANNER_WHITE")]

    def get_next_id(self) -> (int, int):
        c_id = s_id = None
        while c_id is None:
            for index, s_id in enumerate(self.scanner_ids):  # iterate over each reader

                # read id from reader until time expiration
                c_id = self.reader[index].read_id_no_block()
                start_time = time.time()
                while c_id is None and time.time() - start_time < 0.2:
                    c_id = self.reader[index].read_id_no_block()

                if c_id is not None:
                    print(f"Got id {c_id} from scanner {s_id}")  # new id
                    break

        return s_id, c_id


if __name__ == "__main__":
    load_dotenv()
    url = os.getenv("BACKEND_URL")
    kicker_id = os.getenv("KICKER_ID")
    rfid_reader = RfidReader()
    while True:
        print("RFID-Reader ready...")
        scanner_id, card_id = rfid_reader.get_next_id()
        try:
            resp = send_request(url, kicker_id, scanner_id, card_id)
            if not resp.ok:
                print(f"Request failed: {resp.text}\n")
            else:
                print(f"Response: {resp.text}\n")
        except requests.exceptions.ConnectionError:
            print("Failed to establish a connection!\n")
        time.sleep(0.2)
