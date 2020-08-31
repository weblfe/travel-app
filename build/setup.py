#!/usr/bin/env python

import getopt
import os
import sys

docker_deploy_bin = "docker-compose"
default_password = os.getenv("MY_ROOT_PASSWORD")
current_dir = os.getcwd()
os.chdir(current_dir)


def init_env(password=None):
    if password == "":
        password = input("please set mysql root password ")
    os.putenv("MY_ROOT_PASSWORD", password)


def init_files():
    env_file = os.path.join(current_dir, "./.env")
    if not os.path.exists(env_file):
        os.mknod(env_file)


def start():
    init_env(default_password)
    init_files()


def docker():
    state = os.system(docker_deploy_bin + " up -d")
    if state == 0:
        print("success")
    else:
        print("failed")


def build():
    docker_cmd("stop")
    start()
    docker()


def main(argv):
    try:
        opts, args = getopt.getopt(argv[1:], "ht:", ["help", "tag="])
    except getopt.GetoptError:
        print('setup.py -t <tag> \n')
        print(' tags : build , start, stop ,restart, rm , clean \n')
        sys.exit(2)
    opt_num = len(opts)
    args_num = len(args)
    if opt_num == 0 and args_num == 0:
        print("test.py -t <tag> \n")
        print(' tags : build , start, stop ,restart, rm , clean \n')
        sys.exit()

    if opt_num == 0 and args_num == 1:
        sys.exit(docker_cmd(args[0]))

    for opt, arg in opts:
        if opt == '-h':
            print("test.py -t <tag> \n")
            print(' tags : build , start, stop ,restart, rm , clean \n')
            sys.exit()
        elif opt in ("-t", "--tag"):
            action(arg)
        else:
            action("start")


def action(opt):
    if opt == "start":
        docker()
    elif opt == "build":
        build()
    elif opt == "stop":
        docker_cmd(opt)
        print("action stop ok ")
    elif opt == "restart":
        docker_cmd(opt)
    elif opt == "rm":
        docker_cmd("stop")
        docker_cmd("rm -f")
    elif opt == "clean":
        docker_cmd("stop")
        docker_cmd("down --volumes")
    elif opt in ("ps", "status"):
        docker_cmd("ps")


def docker_cmd(cmd):
    return os.system(docker_deploy_bin + " " + cmd)


if __name__ == '__main__':
    main(sys.argv)
