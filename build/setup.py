#!/usr/bin/env python3

import getopt
import os
import sys
from string import Template

image_version = os.getenv("API_VERSION")
image_name = os.getenv("IMAGE_NAME")
docker_file = os.getenv("DOCKERFILE")
docker_deploy_bin = "docker-compose"
default_password = os.getenv("MY_ROOT_PASSWORD")
current_dir = os.getcwd()
os.chdir(current_dir)

'''
    获取 git 用户 | 目录名 自动填充镜像名
'''


def get_default_name():
    user = os.popen("git config --get user.name").read()
    if user is None or user == 0:
        user = os.path.basename(os.path.dirname(current_dir))
    app = os.path.basename(current_dir)
    user = user.strip("\n")
    if user != "":
        return user + "/" + app
    return app


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


def help_menu():
    print('setup.py [command] -t <tag>  -f <dockerfile> -i <image> -v <version> \n')
    print(' tags : build , start, stop ,restart, rm , clean ')
    print(' [command] docker-compose\'s command ')
    print(' v,version build image tag version ')
    print(' i,image build image tag name ')
    print(' f,file  build docker image dockerfile ')


def main(argv):
    try:
        opts, args = getopt.getopt(argv[1:], "ht:v:i:f", ["help", "tag=", "version=", "image=", "file="])
    except getopt.GetoptError:
        help_menu()
        sys.exit(2)
    opt_num = len(opts)
    args_num = len(args)
    if opt_num == 0 and args_num == 0:
        help_menu()
        sys.exit()

    if opt_num == 0 and args_num == 1:
        sys.exit(docker_cmd(args[0]))
    action_opt = version = file = name = ""
    for opt, arg in opts:
        if opt == '-h':
            help_menu()
            sys.exit()
        elif opt in ("-t", "--tag"):
            action_opt = arg
        elif opt in ("-v", "--version"):
            version = arg
        elif opt in ("-f", "--file"):
            file = arg
        elif opt in ("-i", "--image"):
            name = arg
        else:
            return action("start")
    if action_opt != "" and action_opt != "image":
        action(action_opt)
        return 0
    if action_opt == "image":
        return image(file or docker_file, version or image_version, name or image_name)


def action(opt):
    if opt == "start":
        docker()
    elif opt == "image":
        image(docker_file, image_version, image_name)
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


def image(dockerfile_name, version, name):
    tpl = Template('docker build -f ${dockerfile} -t ${name}:${version} .')
    if dockerfile_name == "" or dockerfile_name is None:
        dockerfile_name = input("please set dockerfile name :(default Dockerfile)")
        if dockerfile_name == "" or dockerfile_name == "\n":
            dockerfile_name = "Dockerfile"
    if version == "" or version is None:
        version = input("please set image version : (default 1.0.0)")
        if version == "" or version == "\n":
            version = "1.0.0"
    if name == "" or name is None:
        default_name = get_default_name()
        name = input("please set image name : (default " + default_name + " )")
        if name == "" or name == "\n":
            name = default_name
    cmd = tpl.substitute(dockerfile=dockerfile_name, version=version, name=name)
    print("command : `" + cmd + "`")
    ok = input("that command is ok? you want to exec  Y(es)/N(o) \n")
    ok = ok.strip("\n")
    # 是否确定执行
    if ok in ("y", "Y", "Yes", "yes", "1", "\n", ""):
        return os.system(cmd)
    else:
        return 1


def docker_cmd(cmd):
    return os.system(docker_deploy_bin + " " + cmd)


if __name__ == '__main__':
    main(sys.argv)
