import type { CheckResult } from "../types";

interface CheckListProps {
  checks: CheckResult[];
}

const SEVERITY_ORDER = { critical: 0, high: 1, medium: 2, low: 3 };

const SEVERITY_BADGE: Record<string, string> = {
  critical: "bg-red-100 text-red-700",
  high: "bg-orange-100 text-orange-700",
  medium: "bg-amber-100 text-amber-700",
  low: "bg-gray-100 text-gray-600",
};

const CHECK_LABELS: Record<string, string> = {
  branch_protection: "Branch Protection",
  secret_scanning: "Secret Scanning",
  dependabot_alerts: "Dependabot Alerts",
  cicd_pipeline: "CI/CD Pipeline",
  readme: "README",
  license: "License",
  codeowners: "CODEOWNERS",
  vulnerabilities: "Vulnerabilities",
  stale_repo: "Repository Activity",
};

export default function CheckList({ checks }: CheckListProps) {
  const sorted = [...checks].sort(
    (a, b) =>
      SEVERITY_ORDER[a.severity] - SEVERITY_ORDER[b.severity] ||
      Number(a.passed) - Number(b.passed),
  );

  return (
    <div className="space-y-3">
      {sorted.map((check) => (
        <div
          key={check.name}
          className={`rounded-xl border p-4 transition ${
            check.passed
              ? "border-emerald-200 bg-emerald-50/50"
              : "border-red-200 bg-red-50/50"
          }`}
        >
          <div className="flex items-start gap-3">
            <span className="mt-0.5 text-lg">
              {check.passed ? "✅" : "❌"}
            </span>

            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2 flex-wrap">
                <span className="font-medium text-gray-900">
                  {CHECK_LABELS[check.name] ?? check.name}
                </span>
                <span
                  className={`inline-block rounded-full px-2 py-0.5 text-xs font-medium ${SEVERITY_BADGE[check.severity]}`}
                >
                  {check.severity}
                </span>
              </div>

              <p className="mt-1 text-sm text-gray-600">{check.details}</p>

              {!check.passed && check.recommendation && (
                <p className="mt-2 text-sm text-blue-700 bg-blue-50 rounded-lg px-3 py-2">
                  💡 {check.recommendation}
                </p>
              )}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
