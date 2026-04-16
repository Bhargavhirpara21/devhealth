import { useState } from "react";
import { triggerScan, listRepos, getSummary } from "../api/client";
import type { RepoHealth, Summary } from "../types";
import ScanForm from "../components/ScanForm";
import ScoreCard from "../components/ScoreCard";
import ScoreDistChart from "../components/ScoreDistChart";
import PassRateChart from "../components/PassRateChart";
import RepoTable from "../components/RepoTable";

export default function Dashboard() {
  const [repos, setRepos] = useState<RepoHealth[]>([]);
  const [summary, setSummary] = useState<Summary | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [lastOwner, setLastOwner] = useState<string | null>(null);

  const handleScan = async (owner: string, type: "user" | "org") => {
    setIsLoading(true);
    setError(null);

    try {
      await triggerScan({ owner, type });

      const [repoData, summaryData] = await Promise.all([
        listRepos(owner),
        getSummary(owner),
      ]);

      setRepos(repoData);
      setSummary(summaryData);
      setLastOwner(owner);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "An unexpected error occurred";
      setError(message);
      setRepos([]);
      setSummary(null);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">
          DevHealth Dashboard
        </h1>
        <p className="mt-1 text-sm text-gray-500">
          Scan GitHub repositories for health, security, and best-practice
          compliance.
        </p>
      </div>

      <ScanForm onScan={handleScan} isLoading={isLoading} />

      {error && (
        <div className="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
          {error}
        </div>
      )}

      {summary && (
        <>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            <ScoreCard
              label="Average Score"
              value={summary.average_score}
              subtitle={`across ${summary.total_repos} repo(s)`}
            />
            <ScoreCard
              label="Total Repos"
              value={summary.total_repos}
              color="blue"
            />
            <ScoreCard
              label="Excellent"
              value={summary.score_distribution.excellent}
              color="green"
              subtitle="score ≥ 90"
            />
            <ScoreCard
              label="Critical"
              value={summary.score_distribution.critical}
              color="red"
              subtitle="score < 40"
            />
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            <ScoreDistChart distribution={summary.score_distribution} />
            <PassRateChart passRates={summary.check_pass_rates} />
          </div>
        </>
      )}

      {lastOwner && <RepoTable repos={repos} />}
    </div>
  );
}
