import socket
import yaml
import os

def listen():
    with open('routingtable.yaml', 'r') as f:
        routing_table = yaml.load(f, Loader=yaml.Loader)

    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

    port = os.getenv('MULTICASTER_PORT', 1140)
    server_address = ('localhost', int(port))

    print('starting on %s port %d' % server_address)
    sock.bind(server_address)

    while True:
        data = sock.recv(4096)

        if data:
            print('\n\ngot data: %s' % data)
            for dest in routing_table["nodes"]:
                dest_addr = (dest["addr"], int(dest["port"]))
                print('forwarding to %s:%d' % (dest_addr[0], dest_addr[1]))
                sock.sendto(data, dest_addr)

if __name__ == "__main__":
    listen()
