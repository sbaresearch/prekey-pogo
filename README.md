# Prekey Pogo Artifacts
This artifact contains source code demonstrating how to leverage an open-source WhatsApp client (emulating a companion session) to interact with WhatsApp's internal API.
The software allows retrieval of a session directory and the corresponding prekey material for any arbitrary phone number.

## Security, privacy, and ethical concerns
We want to highlight that all findings (including our corresponding WOOT paper) were reported via Meta's bug bounty program (ticket #10212619137590341) in March 2025.
The ticket was closed as a duplicate and Meta neither followed up on our questions, nor request an embargo for the public release of our findings.
To the best of our knowledge, most attacks presented in the paper remain unfixed and should still work.
Therefore, this client does not cover all attacks presented in the paper to limit the potential for abuse. 
We removed any functionality that could be considered as offensive (e.g., prekey depletion via rapid and iterative prekey retrieval and DoS via prekey clogging by overloading the server with concurrent requests) and just provide a PoC 
which can be used to retrieve prekey bundles manually.

Executing the artifact requires an official WhatsApp account.
Using the PoC code with your official WhatsApp account should not lead to any negative consequences, such as a blocked account.
Although we did not experience any blocked accounts throughout our entire study, we nevertheless recommend using a test account on a burner phone just to be sure as we of course cannot provide any guarantees/legal advice.

## Setup
Execution in host OS (requires golang):
```shell
git clone https://github.com/sbaresearch/prekey-pogo
./setup.sh
go run .
```

Alternatively, containerized version (via podman):
```shell
git clone https://github.com/sbaresearch/prekey-pogo
podman build -t prekey-pogo .
podman run -it -v ./session:/app/session prekey-pogo:latest
```

## Program Execution
After successful execution, the software prints a QR-code that can be used to establish a WhatsApp companion session (similar to WhatsApp Web).
To pair the program with a WhatsApp account, the QR code needs to be scanned from the main device (i.e., the official WhatsApp application on Android or iOS).

### Available commands
```
Enter command (write help to list available commands): help
Available commands:
	(h)elp     -- Show this help message
	(t)arget   -- Update the current target number
	(d)evices  -- Display existing sessions (main and companion devices) for the target number
	(p)rekey   -- Retrieve a prekey bundle for the target number (main device only)
	(c)ombine  -- Retrieve prekey bundles for all existing sessions (main and companion devices) of the target number
	(e)xit     -- Exit the program
```
