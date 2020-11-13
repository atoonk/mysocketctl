#!/usr/bin/env python3
import click


@click.group()
def cli():
    pass


# Import sub commands from commands
from commands.account import account
from commands.login import login
from commands.connect import connect
from commands.socket import socket
from commands.tunnel import tunnel

cli.add_command(account)
cli.add_command(login)
cli.add_command(connect)
cli.add_command(socket)
cli.add_command(tunnel)



if __name__ == "__main__":
    cli()
