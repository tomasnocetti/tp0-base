from enum import Enum
from socket import socket
from typing import List

from common.utils import Contestant

BUF_SIZE = 4096
ENDIAN = 'big'
SEPARATOR = ';'
ENCODING = "utf-8"

"""
Response Structure for Winners consulting

| N~WINNERS | PAYLOAD_LENGTH | ID OF WINNERS |
|  2 bytes  |   4 bytes      |  Dynamic      |

"""


class OpCode(Enum):
    CheckClient = 1

    def __str__(self):
        if(self.value == 1):
            return "Check Contestant"


class Protocol:
    def __init__(self, socket: socket) -> None:
        self.socket = socket
        self.__buffer = bytearray()

    def recv_contestant(self) -> Contestant:
        buff_len = self.__recv_message(4)
        value = int.from_bytes(buff_len, ENDIAN)

        buff_parsed = self.__recv_message(value).decode(ENCODING)
        id, first_name, last_name, birth = buff_parsed.split(SEPARATOR)

        val = Contestant(first_name, last_name, id, birth)

        return val

    def send_response(self, winners: int):

        buf_len = winners.to_bytes(2, ENDIAN)
        self.__send_message(buf_len)

    def get_next_message_type(self) -> OpCode:
        opcode = self.__recv_message(1)

        value = int.from_bytes(opcode, ENDIAN)
        return OpCode(value)

    def __recv_message(self, buf_len: int) -> List[bytes]:

        while(len(self.__buffer) < buf_len):

            buffer = self.socket.recv(BUF_SIZE)
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
