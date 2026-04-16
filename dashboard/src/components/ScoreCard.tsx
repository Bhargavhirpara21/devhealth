interface ScoreCardProps {
  label: string;
  value: number | string;
  subtitle?: string;
  color?: "blue" | "green" | "yellow" | "red";
}

const colorMap = {
  blue: "text-blue-600",
  green: "text-emerald-600",
  yellow: "text-amber-500",
  red: "text-red-600",
};

function getScoreColor(score: number): "green" | "yellow" | "red" {
  if (score >= 70) return "green";
  if (score >= 40) return "yellow";
  return "red";
}

export default function ScoreCard({
  label,
  value,
  subtitle,
  color,
}: ScoreCardProps) {
  const resolvedColor =
    color ?? (typeof value === "number" ? getScoreColor(value) : "blue");

  return (
    <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
      <p className="text-sm font-medium text-gray-500">{label}</p>
      <p className={`mt-2 text-4xl font-bold ${colorMap[resolvedColor]}`}>
        {typeof value === "number" ? Math.round(value) : value}
      </p>
      {subtitle && <p className="mt-1 text-sm text-gray-400">{subtitle}</p>}
    </div>
  );
}
