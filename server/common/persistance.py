from multiprocessing import Lock, Value
from typing import List

from common.utils import Contestant, persist_winners


class Persistance:
    def __init__(self) -> None:
        self.lock = Lock()
        self.total_winners = Value('i', 0)
        self.total_agencies_pending = Value('i', 0)

    def persist_winners(self, con: List[Contestant]):
        self.lock.acquire()

        persist_winners(con)

        self.lock.release()

    def increment_agencies_pending(self):
        self.total_agencies_pending.value += 1

    def decrement_agencies_pending(self):
        self.total_agencies_pending.value -= 1

    def add_winners_to_stats(self, size):
        self.total_winners.value += size
