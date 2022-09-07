from select import select
import socket
import logging
import signal

from common.protocol import Protocol
from common.client_connection import Client
from common.persistance import Persistance


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self.running = True
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.clients = []
        self.persistance = Persistance()
        signal.signal(signal.SIGTERM, self.__exit_gracefully)
        signal.signal(signal.SIGINT, self.__exit_gracefully)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        while self.running:
            try:
                client_sock, addr = self.__accept_new_connection()

                self.__collect_garbage()

                protocol = Protocol(client_sock)
                self.clients.append(Client(addr, protocol, self.persistance))
            except OSError:
                # When closing socket connection error is thrown, skip handling.
                pass

        for client in self.clients:
            client.finish()

        self.__collect_garbage()

    def __exit_gracefully(self, *args):
        self.running = False
        logging.info(
            'Closing server socket connection')
        self._server_socket.close()

    def __collect_garbage(self):
        for client in self.clients:
            if (client.is_closed()):
                client.join()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info("Proceed to accept new connections")
        c, addr = self._server_socket.accept()
        logging.info('Got connection from {}'.format(addr))
        return c, addr
