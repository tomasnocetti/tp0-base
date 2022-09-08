
import logging
from multiprocessing import Process, Value

from common.protocol import OpCode, Protocol
from common.utils import is_winner
from common.persistance import Persistance


class Client:
    def __init__(self, addr, protocol: Protocol, persistance: Persistance) -> None:
        self.addr = addr
        self.protocol = protocol
        self.running = Value('i', 1)
        self.persistance = persistance
        self.process = Process(target=self.run)
        self.process.start()

    def run(self):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            while(self.running.value):
                logging.debug(
                    f'[SERVER - CLIENT {self.addr}] Waiting Message')

                code = self.protocol.get_next_message_type()

                logging.debug(
                    f'[SERVER - CLIENT {self.addr}] New client message received from connection: {code}')

                if (code == OpCode.CheckClient):
                    self.persistance.increment_agencies_pending()
                    self.check_contestants()
                    self.persistance.decrement_agencies_pending()

                if (code == OpCode.GetStats):
                    self.get_stats()

        except OSError:
            pass
        finally:
            self.finish()

    def check_contestants(self):

        logging.debug(
            f'[SERVER - CLIENT {self.addr}] Parsing')

        contestant = self.protocol.recv_contestants()

        logging.debug(
            f'[SERVER - CLIENT {self.addr}] Contestants total: {len(contestant)}')

        winners = []
        for el in contestant:
            if not self.running.value:
                return

            if is_winner(el):
                winners.append(el)

        logging.debug(
            f'[SERVER - CLIENT {self.addr}] Persisting {len(winners)} winners')
        self.persistance.persist_winners(winners)
        self.persistance.add_winners_to_stats(len(winners))

        logging.debug(
            f'[SERVER - CLIENT {self.addr}] Responding to client with {len(winners)} winners')
        self.protocol.send_response(winners)

    def get_stats(self):

        if self.persistance.pending_agencies():
            logging.debug(
                f'[SERVER - CLIENT {self.addr}] Sending pending Stats!')
            self.protocol.send_partial_stats(self.persistance.get_partial())
        else:
            logging.debug(
                f'[SERVER - CLIENT {self.addr}] Sending definite Stats!')
            self.protocol.send_definite_stats(self.persistance.get_partial())

    def finish(self):
        if self.running.value:
            logging.debug(f'[SERVER - CLIENT {self.addr}] Closing connection')
            self.protocol.close()

        self.running.value = 0

    def is_closed(self):
        return not self.running.value

    def join(self):
        logging.debug(f'[SERVER - CLIENT {self.addr}] Joining Process')

        self.process.join()
