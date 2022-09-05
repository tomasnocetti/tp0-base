from select import select
import socket
import logging
import signal

from common.protocol import OpCode, Protocol
from common.utils import is_winner


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self.running = True
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
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
                client_sock = self.__accept_new_connection()
                protocol = Protocol(client_sock)
                self.__handle_client_connection(protocol)
            except OSError:
                # When closing socket connection error is thrown, skip handling.
                pass

    def __exit_gracefully(self, *args):
        self.running = False
        logging.info(
            'Closing socket connection')
        self._server_socket.close()

    def __handle_client_connection(self, protocol: Protocol):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            code = protocol.get_next_message_type()
            logging.info(
                f'New client message received from connection: {code}')

            if (code == OpCode.CheckClient):

                contestant = protocol.recv_contestant()
                winner = is_winner(contestant)

                protocol.send_response(int(winner))

        except OSError:
            logging.info("Error while reading socket")
        finally:
            protocol.close()

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
        return c
