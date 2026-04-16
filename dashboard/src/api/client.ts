import type {
  RepoHealth,
  ScanRequest,
  ScanResponse,
  Summary,
} from "../types";

const BASE_URL = "/api";

class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new ApiError(res.status, body.error || res.statusText);
  }

  return res.json();
}

export async function healthCheck(): Promise<{ status: string; time: string }> {
  return request("/health");
}

export async function triggerScan(req: ScanRequest): Promise<ScanResponse> {
  return request("/scan", {
    method: "POST",
    body: JSON.stringify(req),
  });
}

export async function listRepos(owner: string): Promise<RepoHealth[]> {
  return request(`/repos?owner=${encodeURIComponent(owner)}`);
}

export async function getRepo(
  owner: string,
  repo: string,
): Promise<RepoHealth> {
  return request(
    `/repos/${encodeURIComponent(owner)}/${encodeURIComponent(repo)}`,
  );
}

export async function getSummary(owner: string): Promise<Summary> {
  return request(`/summary?owner=${encodeURIComponent(owner)}`);
}
