# proc-peepin
Capture process cpu and memory and send it off to influx

## Running locally
Rename .env.local to .env and update the environment variables.

Finally, run `go run main.go`.

**NOTE**: If you are running on MacOS you'll need to run `sudo go run main.go`.

## Using systemd
Update the environment variables in the systemd/proc-peepin.service file.

Run `systemctl daemon-reload` then start the timer `systemctl start proc-peepin.timer` and enable it on boot `systemctl enable proc-peepin.timer`. Verify the timer is running by running `systemctl list-timers --all`. Finally, watch the logs by running `journalctl -fu proc-peepin`.
