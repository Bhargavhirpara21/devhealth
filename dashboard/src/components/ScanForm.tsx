import { useState } from "react";

interface ScanFormProps {
  onScan: (owner: string, type: "user" | "org") => void;
  isLoading: boolean;
}

export default function ScanForm({ onScan, isLoading }: ScanFormProps) {
  const [owner, setOwner] = useState("");
  const [type, setType] = useState<"user" | "org">("user");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = owner.trim();
    if (!trimmed) return;
    onScan(trimmed, type);
  };

  return (
    <form onSubmit={handleSubmit} className="flex items-end gap-3">
      <div className="flex-1">
        <label
          htmlFor="owner"
          className="block text-sm font-medium text-gray-700 mb-1"
        >
          GitHub Owner
        </label>
        <input
          id="owner"
          type="text"
          value={owner}
          onChange={(e) => setOwner(e.target.value)}
          placeholder="e.g. facebook, BhargavHirpara"
          className="w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm focus:border-blue-500 focus:ring-2 focus:ring-blue-200 outline-none transition"
          disabled={isLoading}
        />
      </div>

      <div>
        <label
          htmlFor="ownerType"
          className="block text-sm font-medium text-gray-700 mb-1"
        >
          Type
        </label>
        <select
          id="ownerType"
          value={type}
          onChange={(e) => setType(e.target.value as "user" | "org")}
          className="rounded-lg border border-gray-300 px-3 py-2.5 text-sm focus:border-blue-500 focus:ring-2 focus:ring-blue-200 outline-none transition"
          disabled={isLoading}
        >
          <option value="user">User</option>
          <option value="org">Organization</option>
        </select>
      </div>

      <button
        type="submit"
        disabled={isLoading || !owner.trim()}
        className="rounded-lg bg-blue-600 px-6 py-2.5 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition"
      >
        {isLoading ? (
          <span className="flex items-center gap-2">
            <svg
              className="animate-spin h-4 w-4"
              viewBox="0 0 24 24"
              fill="none"
            >
              <circle
                className="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                strokeWidth="4"
              />
              <path
                className="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
              />
            </svg>
            Scanning…
          </span>
        ) : (
          "Scan"
        )}
      </button>
    </form>
  );
}
