import asyncio

HOST, PORT = '0.0.0.0', 8888
ADDR = (HOST, PORT)


async def handle_client(reader, writer):
    req = await reader.read(1024)
    if not req:
        writer.close()
        return

    req_lines = req.split(b'\r\n')
    method, path, version = req_lines[0].decode().split()
    print(method, path, version)

    resp = b'HTTP/1.1 200 OK\r\n\r\nHello, World!'
    writer.write(resp)
    await writer.drain()
    writer.close()


async def runserver():
    server = await asyncio.start_server(handle_client, HOST, PORT)
    async with server:
        await server.serve_forever()


if __name__ == '__main__':
    print('Serving HTTP on:', ADDR)
    asyncio.run(runserver())
