import os
import re
import subprocess
import argparse

def load_env(path):
    env = {}
    with open(path, "r") as f:
        for line in f:
            line = line.strip()
            if not line or line.startswith("#"):
                continue
            key, value = line.split("=", 1)
            env[key.strip()] = value.strip()
    return env

swag_path = os.path.expanduser("~/go/bin/swag")

result = subprocess.run(
    [
        swag_path, 
        "init", 
        "-g", "cmd/api/main.go", 
        "--parseInternal",
        "-o", "docs/"
    ], 
    capture_output=True, 
    text=True
)

file_path = "docs/docs.go"

with open(file_path, "r") as f:
    content = f.read()

content = re.sub(r'\s*LeftDelim:\s*".*?"\s*,?\s*', '', content)
content = re.sub(r'\s*RightDelim:\s*".*?"\s*,?\s*', '\n', content)

hostname = subprocess.run([
        "hostname", "-I"
    ],
    capture_output=True,
    text=True
).stdout.strip().split()[0]

parser = argparse.ArgumentParser()
parser.add_argument("--env_file", required=True)
args = parser.parse_args()

env_path = args.env_file
env = load_env(env_path)

port = env.get("PORT", "8000")

content = re.sub(
    r'Host:\s*"[^"]+"',
    f'Host: \t"{hostname}:{port}"',
    content
)

with open(file_path, "w") as f:
    f.write(content)
    f.close()


print(result.stdout)
print(result.stderr)
