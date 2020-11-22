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
   pp=pprint.PrettyPrinter(indent=4)

   #Create new Vnet
   Socket = client.get_model('Socket')

   socket_obj = Socket(name="test socket",socket_type="tcp")
   socket_obj = client.default.post_socket(payload=socket_obj).response().result
   socket_id = socket_obj['socket_id']
   print("created new socket " + socket_id)
   pp.pprint(socket_obj)

   #create tunnel for socket
   Tunnel = client.get_model('SocketTunnel')
   tunnel_obj = Tunnel()
   tunnel_result = client.default.post_tunnel(socket_id=socket_id,payload=tunnel_obj).response().result
   tunnel_id = tunnel_result['tunnel_id']
   print("created new tunnel " + tunnel_id)
   pp.pprint(tunnel_result)


   #now clean up
   print("cleaning up, deleting created objects")
   delete_tunnel_result = client.default.delete_tunnel(socket_id=socket_id,tunnel_id=tunnel_id).response()
   http_response = delete_tunnel_result.incoming_response
   assert http_response.status_code == 204

   delete_socket_reponse = client.default.delete_socket(socket_id=socket_id).response()
   http_response = delete_socket_reponse.incoming_response
   assert http_response.status_code == 204


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

    token = get_token("atoonk@gmail.com","bla")["token"]

    http_client = RequestsClient()
    http_client.set_api_key(
     api_host, token,
     param_name='x-access-token', param_in='header'
    )

    client = SwaggerClient.from_url(
        'https://'+api_host+'/swagger.json',
        http_client=http_client,
    )

    print("Starting tests")
    test_api()
