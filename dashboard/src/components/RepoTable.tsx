import { useNavigate } from "react-router-dom";
import type { RepoHealth } from "../types";
import { useState } from "react";

interface RepoTableProps {
  repos: RepoHealth[];
}

type SortKey = "full_name" | "score" | "scanned_at";

function scoreBadge(score: number) {
  if (score >= 90)
    return "bg-emerald-100 text-emerald-800";
  if (score >= 70)
    return "bg-blue-100 text-blue-800";
  if (score >= 40)
    return "bg-amber-100 text-amber-800";
  return "bg-red-100 text-red-800";
}

export default function RepoTable({ repos }: RepoTableProps) {
  const navigate = useNavigate();
  const [sortKey, setSortKey] = useState<SortKey>("score");
  const [sortAsc, setSortAsc] = useState(true);
  const [filter, setFilter] = useState("");

  const handleSort = (key: SortKey) => {
    if (sortKey === key) {
      setSortAsc(!sortAsc);
    } else {
      setSortKey(key);
      setSortAsc(true);
    }
  };

  const filtered = repos.filter((r) =>
    r.full_name.toLowerCase().includes(filter.toLowerCase()),
  );

  const sorted = [...filtered].sort((a, b) => {
    let cmp = 0;
    if (sortKey === "full_name") cmp = a.full_name.localeCompare(b.full_name);
    else if (sortKey === "score") cmp = a.score - b.score;
    else cmp = a.scanned_at.localeCompare(b.scanned_at);
    return sortAsc ? cmp : -cmp;
  });

  const sortIcon = (key: SortKey) => {
    if (sortKey !== key) return "↕";
    return sortAsc ? "↑" : "↓";
  };

  return (
    <div className="rounded-xl border border-gray-200 bg-white shadow-sm overflow-hidden">
      <div className="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
        <h3 className="text-sm font-medium text-gray-500">
          Repositories ({filtered.length})
        </h3>
        <input
          type="text"
          placeholder="Filter by name…"
          value={filter}
          onChange={(e) => setFilter(e.target.value)}
          className="rounded-lg border border-gray-300 px-3 py-1.5 text-sm focus:border-blue-500 focus:ring-2 focus:ring-blue-200 outline-none transition w-56"
        />
      </div>

      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-gray-100 text-left text-gray-500">
            <th
              className="px-6 py-3 font-medium cursor-pointer hover:text-gray-800 select-none"
              onClick={() => handleSort("full_name")}
            >
              Repository {sortIcon("full_name")}
            </th>
            <th
              className="px-6 py-3 font-medium cursor-pointer hover:text-gray-800 select-none"
              onClick={() => handleSort("score")}
            >
              Score {sortIcon("score")}
            </th>
            <th className="px-6 py-3 font-medium">Checks</th>
            <th
              className="px-6 py-3 font-medium cursor-pointer hover:text-gray-800 select-none"
              onClick={() => handleSort("scanned_at")}
            >
              Last Scan {sortIcon("scanned_at")}
            </th>
            <th className="px-6 py-3 font-medium" />
          </tr>
        </thead>
        <tbody>
          {sorted.map((repo) => {
            const passed = repo.checks.filter((c) => c.passed).length;
            const total = repo.checks.length;

            return (
              <tr
                key={repo.id}
                className="border-b border-gray-50 hover:bg-gray-50 cursor-pointer transition"
                onClick={() =>
                  navigate(`/repos/${repo.owner}/${repo.repo}`)
                }
              >
                <td className="px-6 py-4 font-medium text-gray-900">
                  {repo.full_name}
                </td>
                <td className="px-6 py-4">
                  <span
                    className={`inline-block rounded-full px-3 py-1 text-xs font-semibold ${scoreBadge(repo.score)}`}
                  >
                    {repo.score}
                  </span>
                </td>
                <td className="px-6 py-4 text-gray-600">
                  {passed}/{total} passed
                </td>
                <td className="px-6 py-4 text-gray-400">
                  {new Date(repo.scanned_at).toLocaleDateString()}
                </td>
                <td className="px-6 py-4 text-gray-400">→</td>
              </tr>
            );
          })}

          {sorted.length === 0 && (
            <tr>
              <td
                colSpan={5}
                className="px-6 py-12 text-center text-gray-400"
              >
                {filter
                  ? "No repos matching your filter."
                  : "No repos scanned yet. Enter a GitHub owner and click Scan."}
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}
