#!/usr/bin/env python3

import getopt
import os
import sys
import subprocess
from string import Template

image_version = os.getenv("API_VERSION")
image_name = os.getenv("IMAGE_NAME")
docker_file = os.getenv("DOCKERFILE")
docker_bin = "docker"
docker_deploy_bin = "docker-compose"
default_password = os.getenv("MY_ROOT_PASSWORD")
deploy_server_name = os.getenv("DEPLOY_SERVER_NAME") or "travel"
deploy_server_dir = os.getenv("DEPLOY_SERVER_DIR") or "/data/dockers/images"
deploy_worker_dir = os.getenv("DEPLOY_WORKER_DIR") or "/data/dockers/travel-app-docker"
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
    if password == "" or len(password) <= 7:
        while len(password) <= 7:
            print("please input mysql root password less 8 word!\n")
            password = input("please set mysql root password ")
        fs = open("build.log", "w+")
        fs.write("mysql_password : " + password + "\n")
        fs.close()
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
    docker_cmd(" up -d --build")


def help_menu():
    print('setup.py [command] -t <tag>  -f <dockerfile> -i <image> -v <version> -d <image>\n')
    print(' tags : build , start, stop ,restart, rm , clean ')
    print(' [command] docker-compose\'s command ')
    print(' v,version build image tag version ')
    print(' i,image build image tag name ')
    print(' f,file  build docker image dockerfile ')
    print(' d,deploy docker image to prod server')


def main(argv):
    try:
        opts, args = getopt.getopt(argv[1:], "-h:-t:-v:-i:-f:-d:", ["help", "tag=", "version=", "image=", "file=",
                                                                    "deploy="])
    except getopt.GetoptError:
        print("msg:" + getopt.GetoptError.msg, "code:" + getopt.GetoptError.opt)
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
        elif opt in ("-d", "--deploy"):
            version = arg
            action_opt = "deploy"
        else:
            return action("start")
    if action_opt != "" and action_opt == "deploy":
        return deploy(version)
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


def format_regex(value):
    return value.replace('/', '\\/').replace('.', '\\.')


def deploy(image_ver):
    if image_ver == "":
        print("镜像名缺失")
        return 1
    # # AAA:8.2,8.2表示镜像版本号
    # docker save -o tar名称.tar AAA:8.2 BBB:5.6
    versions = image_ver.split(":")
    image_tmp = versions[1]
    #  docker images -f 'reference=weblinuxgame/travel-app:v3.36' -q
    result = subprocess.run([docker_bin, "images", "-f", 'reference={}'.format(image_ver), '-q'], capture_output=True)
    if result is not None and result.stdout is not None:
        image_id = str(result.stdout.decode('utf-8')).replace("\n", "")
    else:
        return 0
    image_file = '{}.tar'.format(image_tmp)
    remote_image_file = '{}/{}'.format(deploy_server_dir, image_file)
    export_image_cmd = '{} save -o {} {}'.format(docker_bin, image_file, image_id)
    upload_image_cmd = 'scp {} {}:{}'.format(image_file, deploy_server_name, deploy_server_dir)
    image_import_cmd = 'ssh {} docker import {} {}'.format(deploy_server_name, remote_image_file, image_ver)
    clear_remote_cmd = 'ssh {} rm -rf {}'.format(deploy_server_name, remote_image_file)
    update_service_ver_cmd = ('ssh {} "sed -i.bck  \'s/weblinuxgame\\/travel-app:.*/{}/g\' {}/docker-compose.yml"'.
                              format(deploy_server_name, format_regex(image_ver), deploy_worker_dir))
    restart_service_cmd = ('ssh {} "cd {} ; docker-compose rm -f api ; docker-compose up -d"'.
                           format(deploy_server_name, deploy_worker_dir))
    # 导出镜像
    print(export_image_cmd)
    if os.system(export_image_cmd) == 1:
        print("导出镜像失败")
        return 1
    else:
        print("导出镜像成功")
    # 上传镜像
    print(upload_image_cmd)
    if os.system(upload_image_cmd) == 1:
        print("上传镜像失败")
        return 1
    # 镜像导入
    print(image_import_cmd)
    if os.system(image_import_cmd) == 1:
        print("镜像导入失败")
        return 1
    # 镜像清理
    print(clear_remote_cmd)
    if os.system(clear_remote_cmd) == 1:
        print("镜像清理失败")
        return 1
    # 更新docker-compose.yml
    print(update_service_ver_cmd)
    if os.system(update_service_ver_cmd) == 1:
        print("更新docker-compose失败")
        return 1
    # 服务重启
    print(restart_service_cmd)
    if os.system(restart_service_cmd) == 1:
        print("服务重启失败")
        return 1
    # 清理本地导出镜像文件
    os.unlink(image_file)
    print("---清理本地导出镜像文件---")
    print("---服务部署完成---")
    return 0


if __name__ == '__main__':
    main(sys.argv)
