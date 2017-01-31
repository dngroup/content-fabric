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
         "te_count_per_peer":5,
         "cp_count_per_peer": 5}

with open(TARGET_DOCKER_COMPOSE_FILE,"w") as f:
    f.write(render("./docker-compose.yml.tpl",context))

os.system("docker-compose -f %s up "%TARGET_DOCKER_COMPOSE_FILE)
raw_input()
os.system("docker-compose -f %s down"%TARGET_DOCKER_COMPOSE_FILE)