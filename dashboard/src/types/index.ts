export type CheckName =
  | "branch_protection"
  | "secret_scanning"
  | "dependabot_alerts"
  | "cicd_pipeline"
  | "readme"
  | "license"
  | "codeowners"
  | "vulnerabilities"
  | "stale_repo";

export interface CheckResult {
  name: CheckName;
  passed: boolean;
  details: string;
  severity: "critical" | "high" | "medium" | "low";
  recommendation: string;
}

export interface RepoHealth {
  id: number;
  owner: string;
  repo: string;
  full_name: string;
  url: string;
  score: number;
  checks: CheckResult[];
  scanned_at: string;
  default_branch: string;
  last_commit_at: string;
  is_archived: boolean;
  is_fork: boolean;
}

export interface ScanRequest {
  owner: string;
  type: "org" | "user";
}

export interface ScanResponse {
  message: string;
  repos_found: number;
  scan_id: string;
}

export interface ScoreBuckets {
  critical: number;
  warning: number;
  good: number;
  excellent: number;
}

export interface Summary {
  owner: string;
  total_repos: number;
  average_score: number;
  score_distribution: ScoreBuckets;
  check_pass_rates: Record<CheckName, number>;
  last_scan_at: string;
}

export interface ErrorResponse {
  error: string;
  details?: string;
}
