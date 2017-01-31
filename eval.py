import argparse
import simpy
import docker
import jinja2
import os

TARGET_DOCKER_COMPOSE_FILE="docker-compose.cooked.yml"


def render(tpl_path, context):
    path, filename = os.path.split(tpl_path)
    return jinja2.Environment(
        loader=jinja2.FileSystemLoader(path or './')
    ).get_template(filename).render(context)



context={}
context={"peer_count":5,
         "te_count":20,
         "cp_count": 10,
         "chaincode_id":"fb35eb0071ba47fdcf7ce76359e9f9e4c847a74d8bf0187c1a79ce3468322d378530ca9fc6aefb9c7b074c3d423dfb664142e0c0a3d905a173f4fce253f67372",
         "te_percent":100,
         "te_percent_price": 100,
         "cp_percent":100

         }


with open(TARGET_DOCKER_COMPOSE_FILE,"w") as f:
    f.write(render("./docker-compose.yml.tpl",context))

os.system("docker-compose -f %s up -d "%TARGET_DOCKER_COMPOSE_FILE)
raw_input()
os.system("docker-compose -f %s down"%TARGET_DOCKER_COMPOSE_FILE)