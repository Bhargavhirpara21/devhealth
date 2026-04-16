import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  Cell,
} from "recharts";
import type { ScoreBuckets } from "../types";

interface ScoreDistChartProps {
  distribution: ScoreBuckets;
}

const COLORS = {
  critical: "#ef4444",
  warning: "#f59e0b",
  good: "#3b82f6",
  excellent: "#10b981",
};

export default function ScoreDistChart({
  distribution,
}: ScoreDistChartProps) {
  const data = [
    { name: "Critical (0–39)", value: distribution.critical, key: "critical" },
    { name: "Warning (40–69)", value: distribution.warning, key: "warning" },
    { name: "Good (70–89)", value: distribution.good, key: "good" },
    {
      name: "Excellent (90–100)",
      value: distribution.excellent,
      key: "excellent",
    },
  ];

  return (
    <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
      <h3 className="text-sm font-medium text-gray-500 mb-4">
        Score Distribution
      </h3>
      <ResponsiveContainer width="100%" height={220}>
        <BarChart data={data} layout="vertical" margin={{ left: 20 }}>
          <XAxis type="number" allowDecimals={false} />
          <YAxis
            type="category"
            dataKey="name"
            width={130}
            tick={{ fontSize: 12 }}
          />
          <Tooltip />
          <Bar dataKey="value" radius={[0, 6, 6, 0]} barSize={28}>
            {data.map((entry) => (
              <Cell
                key={entry.key}
                fill={COLORS[entry.key as keyof typeof COLORS]}
              />
            ))}
          </Bar>
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
