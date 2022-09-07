from multiprocessing import Lock
from typing import List

from common.utils import Contestant, persist_winners


class Persistance:
    def __init__(self) -> None:
        self.lock = Lock()

    def persist_winners(self, con: List[Contestant]):
        self.lock.acquire()

        persist_winners(con)

        self.lock.release()
