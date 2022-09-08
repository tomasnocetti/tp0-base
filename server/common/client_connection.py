
import logging
from multiprocessing import Process, Value, Queue

from common.protocol import OpCode, Protocol
from common.utils import is_winner
from common.persistance import Persistance


class Client:
    def __init__(self, addr, protocol: Protocol, persistance: Persistance) -> None:
        self.addr = addr
        self.protocol = protocol
        self.running = Queue()
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
            while(self.running.empty()):
                logging.debug(
                    f'[SERVER - CLIENT {self.addr}] Waiting Message')

                code = self.protocol.get_next_message_type()

                logging.debug(
                    f'[SERVER - CLIENT {self.addr}] New client message received from connection: {code}')

                if (code == OpCode.CheckClient):

                    logging.debug(
                        f'[SERVER - CLIENT {self.addr}] Parsing')

                    contestant = self.protocol.recv_contestants()

                    logging.debug(
                        f'[SERVER - CLIENT {self.addr}] Contestants total: {len(contestant)}')

                    winners = []
                    for el in contestant:
                        if not self.running.empty():
                            return

                        if is_winner(el):
                            winners.append(el)

                    logging.debug(
                        f'[SERVER - CLIENT {self.addr}] Persisting {len(winners)} winners')
                    self.persistance.persist_winners(winners)

                    logging.debug(
                        f'[SERVER - CLIENT {self.addr}] Responding to client with {len(winners)} winners')
                    self.protocol.send_response(winners)

        except OSError:
            logging.info("Error while reading socket")
        finally:
            self.finish()

    def finish(self):
        if self.running.empty():
            logging.debug(f'[SERVER - CLIENT {self.addr}] Closing connection')
            self.protocol.close()

        self.running.put(1)

    def is_closed(self):
        return not self.running.empty()

    def join(self):
        logging.debug(f'[SERVER - CLIENT {self.addr}] Joining Process')

        self.process.join()
