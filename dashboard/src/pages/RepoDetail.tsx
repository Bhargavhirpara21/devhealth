import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { getRepo } from "../api/client";
import type { RepoHealth } from "../types";
import CheckList from "../components/CheckList";

function scoreBadge(score: number) {
  if (score >= 90) return "bg-emerald-100 text-emerald-800";
  if (score >= 70) return "bg-blue-100 text-blue-800";
  if (score >= 40) return "bg-amber-100 text-amber-800";
  return "bg-red-100 text-red-800";
}

export default function RepoDetail() {
  const { owner, repo } = useParams<{ owner: string; repo: string }>();
  const navigate = useNavigate();
  const [data, setData] = useState<RepoHealth | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (!owner || !repo) return;

    setIsLoading(true);
    getRepo(owner, repo)
      .then(setData)
      .catch((err) => setError(err.message))
      .finally(() => setIsLoading(false));
  }, [owner, repo]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-20 text-gray-400">
        Loading…
      </div>
    );
  }

  if (error) {
    return (
      <div className="space-y-4">
        <button
          onClick={() => navigate("/")}
          className="text-sm text-blue-600 hover:text-blue-800 transition"
        >
          ← Back to Dashboard
        </button>
        <div className="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
          {error}
        </div>
      </div>
    );
  }

  if (!data) return null;

  const passed = data.checks.filter((c) => c.passed).length;
  const total = data.checks.length;

  return (
    <div className="space-y-6">
      <button
        onClick={() => navigate("/")}
        className="text-sm text-blue-600 hover:text-blue-800 transition"
      >
        ← Back to Dashboard
      </button>

      <div className="flex items-start justify-between flex-wrap gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">{data.full_name}</h1>
          <div className="mt-2 flex items-center gap-3 text-sm text-gray-500">
            <span>Branch: {data.default_branch}</span>
            <span>·</span>
            <span>
              Last commit:{" "}
              {new Date(data.last_commit_at).toLocaleDateString()}
            </span>
            <span>·</span>
            <a
              href={data.url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-600 hover:text-blue-800 transition"
            >
              View on GitHub ↗
            </a>
          </div>
        </div>

        <div className="text-right">
          <span
            className={`inline-block rounded-full px-4 py-2 text-2xl font-bold ${scoreBadge(data.score)}`}
          >
            {data.score}/100
          </span>
          <p className="mt-1 text-sm text-gray-400">
            {passed}/{total} checks passed
          </p>
        </div>
      </div>

      <CheckList checks={data.checks} />

      <p className="text-xs text-gray-400">
        Scanned at {new Date(data.scanned_at).toLocaleString()}
      </p>
    </div>
  );
}
