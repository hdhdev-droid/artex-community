#!/usr/bin/env python3
"""Copy credentials from a stopped old container before recreating it.

Only Docker CLI and Python's standard library are required. Nothing is printed
from the credential files or Docker errors. The old container is never changed.
"""
import argparse
import io
import os
from pathlib import Path
import subprocess
import tarfile


def docker_file(container, path, optional=False):
    result = subprocess.run(
        ["docker", "cp", f"{container}:{path}", "-"], capture_output=True
    )
    if result.returncode:
        # Only a missing optional file is acceptable, not a Docker/permission error.
        missing = b"Could not find the file" in result.stderr or b"no such file or directory" in result.stderr
        if optional and missing:
            return None
        raise RuntimeError(f"Cannot read {path} from the stopped container")
    with tarfile.open(fileobj=io.BytesIO(result.stdout), mode="r:*") as archive:
        members = archive.getmembers()
        if len(members) != 1 or not members[0].isfile() or not 0 < members[0].size <= 8 * 1024 * 1024:
            raise RuntimeError("Credential source must be a nonempty regular file under 8 MiB")
        return archive.extractfile(members[0]).read()


def copy_credentials(output, files, check=False):
    output = Path(output).absolute()
    # Reject links before resolving or creating anything, including on reruns.
    for path in [output, *output.parents]:
        if path.is_symlink():
            raise RuntimeError("Credential destination must not contain symlinks")
    for name, content in files.items():
        path = output / name
        if path.is_symlink() or path.parent.is_symlink():
            raise RuntimeError("Credential destination must not contain symlinks")
        if path.exists() and (not path.is_file() or path.read_bytes() != content):
            raise RuntimeError("Destination contains different credentials; nothing was overwritten")
    if "chatgpt-oauth/credentials.enc" not in files and (output / "chatgpt-oauth/credentials.enc").exists():
        raise RuntimeError("Destination has an unrelated OAuth vault; use an empty destination")
    if check:
        return
    output.mkdir(parents=True, exist_ok=True, mode=0o700)
    os.chmod(output, 0o700)
    for name, content in files.items():
        path = output / name
        path.parent.mkdir(parents=True, exist_ok=True, mode=0o700)
        os.chmod(path.parent, 0o700)
        if not path.exists():
            # Exclusive create makes an interrupted run resumable, never destructive.
            fd = os.open(path, os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600)
            with os.fdopen(fd, "wb") as stream:
                stream.write(content)
                stream.flush()
                os.fsync(stream.fileno())
        os.chmod(path, 0o600)
        if path.read_bytes() != content:
            raise RuntimeError("Credential copy verification failed; retain the old container")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--container", required=True, help="Stopped old container name or ID")
    parser.add_argument("--source-dir", default="/app", help="Old authentication directory")
    parser.add_argument("--output", default="auth", help="Directory used by the new Compose auth mount")
    parser.add_argument("--check", action="store_true", help="Read and validate only; do not copy")
    args = parser.parse_args()
    state = subprocess.run(
        ["docker", "inspect", "--format", "{{.State.Status}}", args.container], capture_output=True, text=True
    )
    if state.returncode or state.stdout.strip() not in {"exited", "created"}:
        raise RuntimeError("Stop the old container first; credential refresh must be quiescent during the copy")
    output = Path(args.output).absolute()
    workspace = Path("data").resolve()
    if output.resolve().is_relative_to(workspace):
        raise RuntimeError("Authentication files must stay outside the browsable data directory")
    source = args.source_dir.rstrip("/")
    key = docker_file(args.container, f"{source}/jwt.key")
    if len(key.strip()) < 32:
        raise RuntimeError("Existing JWT key is invalid; restore the original backup")
    files = {"jwt.key": key}
    vault = docker_file(args.container, f"{source}/chatgpt-oauth/credentials.enc", optional=True)
    if vault is not None:
        files["chatgpt-oauth/credentials.enc"] = vault
    copy_credentials(output, files, args.check)
    print(f"Authentication {'check' if args.check else 'copy and verification'} complete ({len(files)} files). Old container preserved.")


if __name__ == "__main__":
    try:
        main()
    except (RuntimeError, OSError, tarfile.TarError) as error:
        raise SystemExit(str(error)) from None
