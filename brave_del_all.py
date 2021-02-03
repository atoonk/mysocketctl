#!/usr/bin/env python3

import json
import os
import requests
import sys
import jwt
import pprint

from bravado.client import SwaggerClient
from bravado.requests_client import RequestsClient

api_host = "api.mysocket.io"


def test_api():
    pp = pprint.PrettyPrinter(indent=4)

    Socket = client.get_model("Socket")

    all_sockets = client.default.get_sockets().response().result
    for socket in all_sockets:
        # pp.pprint(socket)
        print(f"{socket['socket_id']} {socket['dnsname']} => {socket.name}")

    # selectivly delete sockets, based on name for example
    if socket.name == "string":
        reponse = client.default.delete_socket(socket_id=socket.socket_id).response()
        assert reponse.incoming_response.status_code == 204


def get_token(user_email, user_pass):
    params = {"email": user_email, "password": user_pass}
    token = requests.post(
        "https://" + api_host + "/login",
        data=json.dumps(params),
        headers={"accept": "application/json", "Content-Type": "application/json"},
    )
    if token.status_code == 401:
        print("Login failed")
        sys.exit(1)
    if token.status_code != 200:
        print(token.status_code, token.text)
        sys.exit(1)
    return token.json()


if __name__ == "__main__":

    token = get_token("xxx@gmail.com", "xxx")["token"]

    http_client = RequestsClient()
    http_client.set_api_key(
        api_host, token, param_name="x-access-token", param_in="header"
    )

    client = SwaggerClient.from_url(
        "https://" + api_host + "/swagger.json",
        http_client=http_client,
    )

    print("Starting tests")
    test_api()
