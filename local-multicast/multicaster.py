import socket
import yaml

def listen():
    with open('routingtable.yaml', 'r') as f:
        routing_table = yaml.load(f, Loader=yaml.Loader)

    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

    server_address = ('localhost', 1140)

    print('starting on %s port %d' % server_address)
    sock.bind(server_address)

    while True:
        data = sock.recv(4096)

        if data:
            for dest in routing_table:
                sock.sendto(data, (dest["address"], dest["port"]))

if __name__ == "__main__":
    print('starting')
    listen()