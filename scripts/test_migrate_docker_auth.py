import importlib.util
from pathlib import Path
import tempfile
import unittest

spec = importlib.util.spec_from_file_location("migration", Path(__file__).with_name("migrate-docker-auth.py"))
migration = importlib.util.module_from_spec(spec)
spec.loader.exec_module(migration)


class MigrationTests(unittest.TestCase):
    def test_pair_copy_check_and_resume(self):
        with tempfile.TemporaryDirectory() as root:
            target = Path(root) / "auth"
            files = {"jwt.key": b"synthetic-test-key-not-a-real-secret", "chatgpt-oauth/credentials.enc": b"fake-vault"}
            migration.copy_credentials(target, files, check=True)
            self.assertFalse(target.exists())
            migration.copy_credentials(target, files)
            migration.copy_credentials(target, files)
            self.assertTrue(all((target / name).read_bytes() == content for name, content in files.items()))

    def test_mismatch_never_overwrites_or_partially_copies(self):
        with tempfile.TemporaryDirectory() as root:
            target = Path(root)
            (target / "chatgpt-oauth").mkdir()
            vault = target / "chatgpt-oauth/credentials.enc"
            vault.write_bytes(b"existing-vault")
            with self.assertRaises(RuntimeError):
                migration.copy_credentials(target, {"jwt.key": b"new-key", "chatgpt-oauth/credentials.enc": b"other-vault"})
            self.assertFalse((target / "jwt.key").exists())
            self.assertTrue(vault.read_bytes() == b"existing-vault")

    def test_unrelated_vault_is_rejected(self):
        with tempfile.TemporaryDirectory() as root:
            target = Path(root)
            (target / "chatgpt-oauth").mkdir()
            (target / "chatgpt-oauth/credentials.enc").write_bytes(b"existing-vault")
            with self.assertRaises(RuntimeError):
                migration.copy_credentials(target, {"jwt.key": b"new-key"})
            self.assertFalse((target / "jwt.key").exists())


if __name__ == "__main__":
    unittest.main()
