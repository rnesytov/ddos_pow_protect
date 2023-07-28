import os
import sys
import socket
import logging
import pow
import time

logger = logging.getLogger("client")
logging.basicConfig(stream=sys.stdout, level=logging.INFO)


def main() -> None:
    # Setup params
    host = os.getenv("HOST", "0.0.0.0")
    port = os.getenv("PORT", "8000")
    try:
        port = int(port)
    except ValueError:
        logger.error(f"Invalid port: {port}")
        sys.exit(1)

    # Connect to the server
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect((host, port))
        logger.info(f"Connected to {host}:{port}")

        # Read the challenge and difficulty
        data = s.makefile("rb").readline()
        parts = data.decode().strip().split(";")
        if len(parts) != 2:
            logger.error("Invalid challenge")
            return
        challenge, difficulty = parts
        try:
            difficulty = int(difficulty)
        except ValueError:
            logger.error("Invalid difficulty")
            return
        logger.info(
            f"Starting to solve the challenge: {challenge} with difficulty: {difficulty}"
        )
        now = time.time()

        # Solve the proof of work
        nonce = pow.generate_proof_of_work(challenge, difficulty)
        logger.info(f"Found nonce={nonce} in {time.time() - now} seconds")
        s.sendall(f"{nonce}\n".encode())

        # Read the response
        data = s.makefile("rb").readline()
        logger.info(data.decode().strip())


if __name__ == "__main__":
    main()
