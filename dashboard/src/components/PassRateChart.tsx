import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  Cell,
} from "recharts";
import type { CheckName } from "../types";

interface PassRateChartProps {
  passRates: Record<CheckName, number>;
}

const CHECK_LABELS: Record<CheckName, string> = {
  branch_protection: "Branch Protection",
  secret_scanning: "Secret Scanning",
  dependabot_alerts: "Dependabot",
  cicd_pipeline: "CI/CD Pipeline",
  readme: "README",
  license: "License",
  codeowners: "CODEOWNERS",
  vulnerabilities: "Vulnerabilities",
  stale_repo: "Active Repo",
};

function rateColor(rate: number): string {
  if (rate >= 80) return "#10b981";
  if (rate >= 50) return "#f59e0b";
  return "#ef4444";
}

export default function PassRateChart({ passRates }: PassRateChartProps) {
  const data = Object.entries(passRates)
    .map(([key, value]) => ({
      name: CHECK_LABELS[key as CheckName] ?? key,
      rate: Math.round(value),
    }))
    .sort((a, b) => a.rate - b.rate);

  return (
    <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
      <h3 className="text-sm font-medium text-gray-500 mb-4">
        Check Pass Rates
      </h3>
      <ResponsiveContainer width="100%" height={300}>
        <BarChart data={data} layout="vertical" margin={{ left: 30 }}>
          <XAxis type="number" domain={[0, 100]} unit="%" />
          <YAxis
            type="category"
            dataKey="name"
            width={140}
            tick={{ fontSize: 12 }}
          />
          <Tooltip formatter={(value: number) => `${value}%`} />
          <Bar dataKey="rate" radius={[0, 6, 6, 0]} barSize={24}>
            {data.map((entry, index) => (
              <Cell key={index} fill={rateColor(entry.rate)} />
            ))}
          </Bar>
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
