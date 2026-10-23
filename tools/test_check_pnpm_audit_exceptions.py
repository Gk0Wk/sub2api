"""审计命令失败时必须阻止安全检查通过。"""

import json
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest


CHECKER = Path(__file__).with_name("check_pnpm_audit_exceptions.py")


class AuditReportTests(unittest.TestCase):
    def run_check(
        self, report: object, exceptions: str = "version: 1\nexceptions:\n"
    ) -> subprocess.CompletedProcess[str]:
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            audit = root / "audit.json"
            audit.write_text(json.dumps(report), encoding="utf-8")
            exception_file = root / "exceptions.yml"
            exception_file.write_text(exceptions, encoding="utf-8")
            return subprocess.run(
                [
                    sys.executable,
                    str(CHECKER),
                    "--audit",
                    str(audit),
                    "--exceptions",
                    str(exception_file),
                ],
                check=False,
                capture_output=True,
                text=True,
                encoding="utf-8",
            )

    def test_should_reject_registry_error_report(self) -> None:
        report = {
            "error": {
                "code": "ERR_PNPM_AUDIT_ENDPOINT_NOT_EXISTS",
                "message": "The audit endpoint does not exist.",
            }
        }
        result = self.run_check(report)
        self.assertNotEqual(result.returncode, 0)
        self.assertNotIn("Audit exceptions validated.", result.stdout)

    def test_should_reject_error_even_with_empty_advisories(self) -> None:
        result = self.run_check({"error": {"code": "FETCH_ERROR"}, "advisories": {}})
        self.assertNotEqual(result.returncode, 0)

    def test_should_reject_missing_or_invalid_vulnerability_data(self) -> None:
        for report in [{}, {"advisories": []}, {"vulnerabilities": None}, []]:
            with self.subTest(report=report):
                self.assertNotEqual(self.run_check(report).returncode, 0)

    def test_should_accept_empty_successful_reports_in_supported_formats(self) -> None:
        for field in ["advisories", "vulnerabilities"]:
            with self.subTest(field=field):
                self.assertEqual(self.run_check({field: {}}).returncode, 0)

    def test_should_require_a_current_exception_for_high_advisories(self) -> None:
        report = {
            "advisories": {
                "1": {
                    "module_name": "example-package",
                    "severity": "high",
                    "github_advisory_id": "GHSA-example",
                    "title": "Example advisory",
                }
            }
        }
        self.assertNotEqual(self.run_check(report).returncode, 0)
        exception = (
            "version: 1\nexceptions:\n"
            "  - package: example-package\n"
            "    advisory: GHSA-example\n"
            "    severity: high\n"
            "    mitigation: Fixture-only exception\n"
            "    expires_on: 2099-01-01\n"
        )
        self.assertEqual(self.run_check(report, exception).returncode, 0)
        self.assertNotEqual(
            self.run_check(report, exception.replace("2099-01-01", "2000-01-01")).returncode,
            0,
        )


if __name__ == "__main__":
    unittest.main()
