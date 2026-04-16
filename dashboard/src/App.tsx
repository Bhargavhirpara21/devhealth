import { BrowserRouter, Routes, Route } from "react-router-dom";
import Dashboard from "./pages/Dashboard";
import RepoDetail from "./pages/RepoDetail";

export default function App() {
  return (
    <BrowserRouter>
      <div className="min-h-screen bg-gray-50">
        <header className="border-b border-gray-200 bg-white">
          <div className="mx-auto max-w-6xl px-6 py-4 flex items-center gap-3">
            <span className="text-xl"></span>
            <span className="text-lg font-bold text-gray-900">DevHealth</span>
            <span className="text-sm text-gray-400">
              GitHub Org Health Scanner
            </span>
          </div>
        </header>

        <main className="mx-auto max-w-6xl px-6 py-8">
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/repos/:owner/:repo" element={<RepoDetail />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  );
}
