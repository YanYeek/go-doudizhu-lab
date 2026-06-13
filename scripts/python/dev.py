#!/usr/bin/env python3
"""Cross-platform development commands for the game server."""

from __future__ import annotations

import argparse
import os
import shutil
import socket
import subprocess
import sys
import tempfile
import time
from pathlib import Path
from typing import Sequence


REPO_ROOT = Path(__file__).resolve().parents[2]
SERVER_DIR = REPO_ROOT / "server"

SERVICES = {
    "etcd": 2379,
    "redis": 6379,
    "websocket gate": 3553,
}


class DevCommandError(RuntimeError):
    """Raised when a development command cannot continue."""


def run(
    command: Sequence[str],
    *,
    cwd: Path = SERVER_DIR,
    check: bool = True,
    capture_output: bool = False,
) -> subprocess.CompletedProcess[str]:
    print(f"\n> {' '.join(command)}", flush=True)
    return subprocess.run(
        command,
        cwd=cwd,
        check=check,
        text=True,
        capture_output=capture_output,
    )


def require_command(name: str) -> None:
    if shutil.which(name) is None:
        raise DevCommandError(
            f"找不到命令 {name!r}，请先安装并确保它已经加入 PATH。"
        )


def docker_is_ready() -> bool:
    result = subprocess.run(
        ["docker", "info"],
        cwd=SERVER_DIR,
        check=False,
        text=True,
        capture_output=True,
    )
    return result.returncode == 0


def start_docker_desktop() -> bool:
    if sys.platform == "win32":
        candidates = [
            Path(os.environ.get("ProgramFiles", "")) / "Docker/Docker/Docker Desktop.exe",
            Path(os.environ.get("LOCALAPPDATA", "")) / "Docker/Docker Desktop.exe",
        ]
        executable = next((path for path in candidates if path.is_file()), None)
        if executable is None:
            return False

        print(f"\nDocker 引擎未运行，正在启动 Docker Desktop：{executable}")
        subprocess.Popen(
            [str(executable)],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
            creationflags=subprocess.CREATE_NEW_PROCESS_GROUP
            | subprocess.DETACHED_PROCESS,
        )
        return True

    if sys.platform == "darwin" and shutil.which("open"):
        print("\nDocker 引擎未运行，正在启动 Docker Desktop。")
        subprocess.Popen(
            ["open", "-a", "Docker"],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
        )
        return True

    return False


def wait_for_docker(timeout: float = 120.0) -> None:
    deadline = time.monotonic() + timeout
    print("等待 Docker 引擎就绪", end="", flush=True)

    while time.monotonic() < deadline:
        if docker_is_ready():
            print("\nDocker 引擎已就绪。")
            return
        print(".", end="", flush=True)
        time.sleep(2)

    print()
    raise DevCommandError(
        f"等待 Docker 引擎启动超时（{int(timeout)} 秒）。"
        "请查看 Docker Desktop 是否需要手动确认或登录。"
    )


def require_docker(*, auto_start: bool = False) -> None:
    require_command("docker")
    if docker_is_ready():
        return

    if auto_start and start_docker_desktop():
        wait_for_docker()
        return

    raise DevCommandError(
        "Docker 引擎尚未运行。Windows 与 macOS 请启动 Docker Desktop，"
        "Linux 请启动 Docker 服务。"
    )


def port_is_open(port: int, host: str = "127.0.0.1") -> bool:
    try:
        with socket.create_connection((host, port), timeout=0.5):
            return True
    except OSError:
        return False


def wait_for_ports(ports: Sequence[int], timeout: float = 30.0) -> None:
    deadline = time.monotonic() + timeout
    pending = set(ports)

    while pending and time.monotonic() < deadline:
        pending = {port for port in pending if not port_is_open(port)}
        if pending:
            time.sleep(0.5)

    if pending:
        values = ", ".join(str(port) for port in sorted(pending))
        raise DevCommandError(f"等待端口超时：{values}")


def deps_up() -> None:
    require_docker(auto_start=True)
    run(["docker", "compose", "up", "-d"])
    wait_for_ports([SERVICES["etcd"], SERVICES["redis"]])
    print("\nRedis 与 etcd 已就绪。")


def deps_down() -> None:
    require_docker()
    run(["docker", "compose", "down"])


def server() -> None:
    require_command("go")
    print("\n启动 Gate。按 Ctrl+C 正常停止服务器。")
    run(["go", "run", "./cmd/server"])


def testclient() -> None:
    require_command("go")
    if not port_is_open(SERVICES["websocket gate"]):
        raise DevCommandError(
            "Gate 未在 127.0.0.1:3553 监听。请先在另一个终端运行 "
            "`python scripts/python/dev.py up` 启动 Gate 与依赖。"
        )
    print("\n运行 Go 测试客户端：向 Gate 发送一条 Greet 请求并打印响应。")
    run(["go", "run", "./cmd/testclient"])


def up() -> None:
    deps_up()
    server()


def test() -> None:
    require_command("go")
    run(["go", "test", "./..."])


def vet() -> None:
    require_command("go")
    run(["go", "vet", "./..."])


def build() -> None:
    require_command("go")
    with tempfile.TemporaryDirectory(prefix="doudizhu-build-") as directory:
        output = Path(directory) / "server"
        run(["go", "build", "-o", str(output), "./cmd/server"])


def check() -> None:
    test()
    vet()
    build()
    print("\n测试、静态检查与构建全部通过。")


def status() -> None:
    print("\n本地服务状态：")
    for name, port in SERVICES.items():
        state = "可连接" if port_is_open(port) else "未监听"
        print(f"  {name:<16} 127.0.0.1:{port:<5} {state}")

    if shutil.which("docker") is None:
        print("\nDocker：未安装或不在 PATH")
        return

    if not docker_is_ready():
        print("\nDocker：已安装，但引擎未运行")
        return

    run(["docker", "compose", "ps"], check=False)


def doctor() -> None:
    print(f"仓库目录：{REPO_ROOT}")
    print(f"服务器目录：{SERVER_DIR}")

    for command in ("python", "go", "docker"):
        path = shutil.which(command)
        print(f"{command:<8}：{path or '未找到'}")

    status()


def create_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        description="斗地主服务器跨平台开发助手",
    )
    subparsers = parser.add_subparsers(dest="command", required=True)

    commands = {
        "up": "启动 Redis、etcd，并在前台运行 Gate",
        "server": "只在前台运行 Gate",
        "testclient": "运行一次性 Go 测试客户端，向 Gate 发消息验证链路",
        "deps-up": "只启动 Redis 与 etcd",
        "deps-down": "停止并删除 Redis 与 etcd 容器",
        "status": "查看依赖与 Gate 端口状态",
        "test": "运行全部 Go 测试",
        "vet": "运行 Go 静态检查",
        "build": "构建服务器",
        "check": "依次运行测试、静态检查与构建",
        "doctor": "检查本机开发环境",
    }

    for name, help_text in commands.items():
        subparsers.add_parser(name, help=help_text)

    return parser


def main() -> int:
    args = create_parser().parse_args()
    commands = {
        "up": up,
        "server": server,
        "testclient": testclient,
        "deps-up": deps_up,
        "deps-down": deps_down,
        "status": status,
        "test": test,
        "vet": vet,
        "build": build,
        "check": check,
        "doctor": doctor,
    }

    try:
        commands[args.command]()
    except KeyboardInterrupt:
        print("\n已停止。")
        return 130
    except (DevCommandError, subprocess.CalledProcessError) as error:
        print(f"\n错误：{error}", file=sys.stderr)
        return 1

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
