# Test task for Server Engineer
Design and implement “Word of Wisdom” tcp server.
- TCP server should be protected from DDoS attacks with the Prof of Work
(https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
- Docker file should be provided both for the server and for the client that solves the POW challenge.

# Implementation details
## Choice of the POW algorithm
There are many different algorithms for PoW. One of most simple and effective is based on the number of leading zeros in the hash of the message.
The algorithm works as follows:
- Server sends a message with a random string and a difficulty level (number of leading zeros).
- Client should find a random string that will give a hash with the required number of leading zeros.
- Client sends the found string to the server.
- Server checks the hash of the received string and the difficulty level.
- If the hash is correct, the server sends a message with a quote.

The default choice for hash function is SHA256. But this function can be easily optimized for GPU or ASIC. So, for this task I decided to use Scrypt hash function. It is memory-hard function, so it is not so easy to optimize it for GPU or ASIC.
It also gives us the ability to set the required memory and CPU usage by changing the parameters N, r, p., which allows us to additionally adjust the difficulty of the PoW.

# Usage
Run this to start server and client:
```
docker-compose up
```

After starting the server you can also test it with `telnet` or `nc` by connecting to port 8080.

# Notes
This is a PoC implementation. It is not production ready. It is not secure. It is not optimized. It is not tested.
There are many things that can be improved. For example:
- Rewrite to compiled language (C, C++, Rust, Go, etc.) for better performance.
- Rewrite protocol to use binary messages instead of text and add additional metadata to messages (version, timestamp, etc.).
- Test with different devices and different Scrypt parameters to find the best balance between performance and security.


I tested algorithm with current Scrypt parameters (N=1024, r=1, p=1) and different difficulties on my laptop (M1 Max). Here are the results:
- D=3 – avg=0.83
- D=2 – avg=0.05

So on modern devices it is not so hard to solve PoW with difficulty 2. But it is still can be hard for IoT devices.
