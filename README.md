# hcotf-cli

cli that syncs a Peloton Cycle stack with the hcotf google sheet of classes

# How do I use it

- Download the latest program for either
  [Windows](https://github.com/caesarshift/hcotf-cli/releases/download/v0.1/windows.zip), [Linux](https://github.com/caesarshift/hcotf-cli/releases/download/v0.1/linux.zip), or [Mac OS](https://github.com/caesarshift/hcotf-cli/releases/download/v0.1/darwin.zip)
- Extract the hcotf-cli program from the zipped folder
- Run the porgram using the following format:

`hcotf-cli -username <YOURUSERNAME> -password <YOURPASSWORD>`

Where `YOURUSERNAME` is your Peloton username or email address, and `YOURPASSWORD` is your Peloton
password.

# What it does

The program will login to Peloton with the provided username and password.

Assuming successful authentication to Peloton, the program will then attempt to pull today's schedule from the public
hardCORE on the floor google sheet. If it can find today's schedule, it will then CLEAR your
existing Peloton stack and load today's schedule into your Peloton stack.

# What it can't do

- Determine which class to take. If the hcotf google sheet has classes that are listed as OR, it will load both of
  them into your stack

# Why did you make this?

My wife initially joined the hcotf facebook group, and I saw her having to daily locate the classes.
When I asked her if everyone was having to do this, she said "yes", and then showed me the long list
of people asking for a better way. This is my attempt at "a better way."
