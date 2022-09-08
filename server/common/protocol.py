from enum import Enum
import logging
from socket import socket
from typing import List

from common.utils import Contestant

BUF_SIZE = 4096
ENDIAN = 'big'
SEPARATOR = ';'
CON_SEPARATOR = '|'
ENCODING = "utf-8"

"""
Response Structure for Winners consulting

| N~WINNERS | PAYLOAD_LENGTH | ID OF WINNERS |
|  2 bytes  |   4 bytes      |  Dynamic      |

"""


class OpCode(Enum):
    CheckClient = 1
    GetStats = 2

    def __str__(self):
        if(self.value == 1):
            return "Check Contestant"

        if(self.value == 2):
            return "Get Stats"


class Protocol:
    def __init__(self, socket: socket) -> None:
        self.socket = socket
        self.__buffer = bytearray()

    def recv_contestants(self) -> Contestant:
        buff_len = self.__recv_message(4)
        value = int.from_bytes(buff_len, ENDIAN)

        buff_parsed = self.__recv_message(value).decode(ENCODING)
        raw_contestants = buff_parsed.split(CON_SEPARATOR)

        contestants = []
        for el in raw_contestants:
            id, first_name, last_name, birth = el.split(SEPARATOR)
            val = Contestant(first_name, last_name, id, birth)

            contestants.append(val)

        return contestants

    def send_response(self, winners: List[Contestant]):

        ids = map(lambda x: x.document, winners)
        ids_buf = bytes(CON_SEPARATOR.join(ids), ENCODING)

        buf = len(ids_buf).to_bytes(4, ENDIAN)
        buf += ids_buf

        self.__send_message(buf)

    def send_partial_stats(self, winners):

        buf = b'1'
        buf += winners.to_bytes(4, ENDIAN)

        self.__send_message(buf)

    def send_definite_stats(self, winners):

        buf = b'0'
        buf += winners.to_bytes(4, ENDIAN)

        self.__send_message(buf)

    def get_next_message_type(self) -> OpCode:
        opcode = self.__recv_message(1)

        value = int.from_bytes(opcode, ENDIAN)
        return OpCode(value)

    def __recv_message(self, buf_len: int) -> List[bytes]:

        while(len(self.__buffer) < buf_len):

            buffer = self.socket.recv(BUF_SIZE)
            if (len(buffer) == 0):
                raise OSError

            self.__buffer += buffer

        val = self.__buffer[:buf_len]
        del self.__buffer[:buf_len]
        return val

    def __send_message(self, buf: bytearray):
        sent = 0
        while(sent < len(buf)):
            sent += self.socket.send(buf[sent:])

    def close(self):
        self.socket.close()
