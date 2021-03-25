# hcotf-cli

CLI that syncs a Peloton Cycle stack with the "hardCORE On The Floor" google sheet

# How do I use it

- Download the latest program for either
  [Windows](https://github.com/caesarshift/hcotf-cli/releases/download/v0.2/windows.zip), [Linux](https://github.com/caesarshift/hcotf-cli/releases/download/v0.2/linux.zip), or [Mac OS](https://github.com/caesarshift/hcotf-cli/releases/download/v0.2/darwin.zip)
- Extract the hcotf-cli program from the zipped folder
- Run the program from the command line (see below if you've never done this before) using the following format:

`hcotf-cli -username YOURUSERNAME -password YOURPASSWORD`

Where `YOURUSERNAME` is your Peloton username or email address, and `YOURPASSWORD` is your Peloton
password.

# Windows - How to install/run a program from the "command line"

- If you've never used a command program before, install it on your Desktop, so you can easily find
  it.
- To open cmd.exe ("command line), open File Explorer and navigate to the Desktop
- Type "cmd" into the File Explorer navigator and hit enter
- A cmd.exe window opens with the Desktop as the current working directory

![](hcotf-cli-windows-help.gif)

# MacOS - How to run a program from the "command line"

- Open the app named "Terminal" which will open a command window
- Navigate to the folder that contains hcotfl-cli by using `cd` (e.g., `cd ~/Downloads/darwin`)

# Linux - How to run a program from the "command line"

- You installed your OS from the command line. You got this down already.

# Load a specific date

If you want to load an alternate date's classes, use the `date` parameter. Example:

`hcotf-cli -username YOURUSERNAME -password YOURPASSWORD -date 3/23/2021`

NOTE: the specified date must match the sheet date _exactly_. For example, `-date 3/23/2021` will
work. `-date 03/23/2021` will not (because the google sheet has the date listed as `3/23/2021`).

# What it does

The program will login to Peloton with the provided username and password.

Assuming successful authentication to Peloton, the program will then attempt to pull today's
schedule from the public "hardCORE On The Floor" google sheet. If it can find today's schedule, it
will then CLEAR your existing Peloton stack and load today's schedule into your Peloton stack.

# What it can't do

- Determine which class to take. If the "hardCORE On The Floor" google sheet has two classes that
are listed as OR, it will load both classes into your stack.

# Why did you make this?

My wife initially joined the hcotf facebook group, and I saw her having to daily locate the classes.
When I asked her if everyone was having to do this, she said "yes", and then showed me the long list
of people asking for a better way. This is my attempt at "a better way."
