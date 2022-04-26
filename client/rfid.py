import time
from mfrc522 import SimpleMFRC522
import requests


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
        c_id = s_id = None
        while c_id is None:
            for index, s_id in enumerate(self.scanner_ids):  # iterate over each reader

                # read id from reader until time expiration
                c_id = self.reader[index].read_id_no_block()
                start_time = time.time()
                while c_id is None and time.time() - start_time < 0.2:
                    c_id = self.reader[index].read_id_no_block()

                if c_id is not None and c_id != self.last_id:
                    print(f"Got id {c_id} from scanner {s_id}")  # new id
                    break
                elif c_id == self.last_id:
                    c_id = None  # same id as last read

        self.last_id = c_id
        return s_id, c_id


if __name__ == "__main__":
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
        time.sleep(0.2)
