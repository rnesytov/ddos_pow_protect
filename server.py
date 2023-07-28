import asyncio
import random
import logging
import os
import pow
import string
import sys

logger = logging.getLogger("server")
logging.basicConfig(stream=sys.stdout, level=logging.INFO)

SOLVE_TIMEOUT = 10  # seconds


def get_challenge(N=32) -> str:
    return "".join(
        random.choice(string.ascii_uppercase + string.digits) for _ in range(N)
    )


def client_handler(quotes: list[str], difficulty: int) -> callable:
    async def handler(reader: asyncio.StreamReader, writer: asyncio.StreamWriter):
        # Generate a random challenge
        challenge = get_challenge()
        writer.write(f"{challenge};{difficulty}\n".encode())
        await writer.drain()

        # Read the client's response
        try:
            async with asyncio.timeout(SOLVE_TIMEOUT):
                response = await reader.readline()
        except asyncio.TimeoutError:
            writer.write("Timeout\n".encode())
            await writer.drain()
            writer.close()
            return
        response = response.decode().strip()
        try:
            nonce = int(response)
        except ValueError:
            writer.write("Invalid nonce\n".encode())
            await writer.drain()
            writer.close()
            return

        # Verify the proof of work
        if pow.verify_proof_of_work(challenge, nonce, difficulty):
            # Send a random quote from the collection
            quote = random.choice(quotes)
            writer.write(f"{quote}\n".encode())
        else:
            writer.write("Invalid proof of work\n".encode())

        await writer.drain()
        writer.close()

    return handler


async def main() -> None:
    # Setup params
    host = os.getenv("HOST", "0.0.0.0")
    port = os.getenv("PORT", "8000")
    raw_difficulty = os.getenv("DIFFICULTY", "3")
    try:
        difficulty = int(raw_difficulty)
    except ValueError:
        logger.error(f"Invalid difficulty: {raw_difficulty}")
        sys.exit(1)

    quotes_path = os.getenv("QUOTES_PATH", "quotes.txt")
    try:
        with open(quotes_path) as f:
            quotes = f.read().splitlines()
    except FileNotFoundError:
        logger.error(f"Quotes file not found: {quotes_path}")
        sys.exit(1)

    # Start server
    server = await asyncio.start_server(client_handler(quotes, difficulty), host, port)
    logger.info(f"Serving on {host}:{port}")
    async with server:
        await server.serve_forever()


if __name__ == "__main__":
    asyncio.run(main())
