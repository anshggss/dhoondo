import { useMemo, useState } from "react";
import {
  ArrowLeft,
  ArrowUpRight,
  Check,
  Globe2,
  LoaderCircle,
  Radio,
  Search,
  X,
} from "lucide-react";
import "./index.css";

type DispatchState = "idle" | "sending" | "success" | "error";

const API_URL =
  import.meta.env?.BUN_PUBLIC_API_URL ??
  `${window.location.protocol}//${window.location.hostname}:8080`;

const starterUrls = [
  "https://example.com",
  "https://news.ycombinator.com",
  "https://developer.mozilla.org",
];

export function App() {
  const [draft, setDraft] = useState(starterUrls.join("\n"));
  const [state, setState] = useState<DispatchState>("idle");
  const [message, setMessage] = useState("Ready for a new crawl.");
  const [publishedCount, setPublishedCount] = useState(0);

  const urls = useMemo(
    () =>
      [...new Set(draft.split(/\n|,/).map((url) => url.trim()).filter(Boolean))],
    [draft],
  );

  async function publishWebsites(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (urls.length === 0) {
      setState("error");
      setMessage("Add at least one URL before dispatching.");
      return;
    }

    setState("sending");
    setMessage("Opening the publishing channel…");

    try {
      const response = await fetch(`${API_URL}/publish`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ websites: urls }),
      });

      const result = await response.json().catch(() => null);
      if (!response.ok) {
        throw new Error(result?.message ?? "The server rejected this batch.");
      }

      setPublishedCount(result?.published ?? urls.length);
      setState("success");
      setMessage("Batch accepted. The crawler has the next move.");
    } catch (error) {
      setState("error");
      setMessage(
        error instanceof Error
          ? error.message
          : "Could not reach the publishing server.",
      );
    }
  }

  if (window.location.pathname === "/search") {
    return <SearchPage />;
  }

  return (
    <main className="app-shell">
      <section className="workbench">
        <header className="topline">
          <div className="breadcrumb">
            <Radio size={15} strokeWidth={2.5} />
            <span>Dispatch console</span>
          </div>
          <a className="server-link" href={`${API_URL}/`} target="_blank" rel="noreferrer">
            Server status <ArrowUpRight size={14} />
          </a>
        </header>

        <div className="content-grid">
          <section className="intro">
            <p className="kicker">A small map of the open web</p>
            <h1>
              Put the web
              <br />
              <em>on the map.</em>
            </h1>
            <p className="intro-copy">
              Send a handful of starting points into the crawler. Dhoondo will
              follow the trail, collect the language, and make it searchable.
            </p>
            <div className="signal-line">
              <span />
              <span />
              <span />
              <span />
              <span />
            </div>
          </section>

          <aside className="dispatch-note">
            <div className="note-icon">
              <Globe2 size={22} />
            </div>
            <p className="note-label">How it works</p>
            <p className="note-copy">
              One URL becomes a thread. Threads become a living index.
            </p>
          </aside>

          <form className="dispatch-panel" onSubmit={publishWebsites}>
            <div className="panel-heading">
              <div>
                <p className="kicker">New dispatch</p>
                <h2>Choose your starting points</h2>
              </div>
              <span className="count-badge">{urls.length.toString().padStart(2, "0")}</span>
            </div>
            <label htmlFor="website-list">Websites to crawl</label>
            <textarea
              id="website-list"
              value={draft}
              onChange={(event) => {
                setDraft(event.target.value);
                if (state !== "idle") {
                  setState("idle");
                  setMessage("Ready for a new crawl.");
                }
              }}
              placeholder="https://example.com"
              spellCheck={false}
            />
            <div className="form-footer">
              <span>One URL per line, or separate with commas.</span>
              <button type="submit" disabled={state === "sending"}>
                {state === "sending" ? (
                  <>
                    <LoaderCircle className="spin" size={16} /> Dispatching
                  </>
                ) : (
                  <>
                    Dispatch batch <ArrowUpRight size={16} />
                  </>
                )}
              </button>
            </div>
          </form>

          <a className="search-index-link" href="/search">
            <span>
              <Search size={16} />
              Search the indexed web
            </span>
            <ArrowUpRight size={16} />
          </a>

          <section className={`dispatch-result dispatch-result--${state}`} aria-live="polite">
            <div className="result-symbol">
              {state === "success" ? (
                <Check size={18} />
              ) : state === "error" ? (
                <X size={18} />
              ) : (
                <span className="result-symbol__idle" />
              )}
            </div>
            <div>
              <p className="result-title">
                {state === "success"
                  ? `${publishedCount} ${publishedCount === 1 ? "website" : "websites"} published`
                  : state === "error"
                    ? "Dispatch paused"
                    : "Awaiting dispatch"}
              </p>
              <p className="result-copy">{message}</p>
            </div>
          </section>
        </div>

        <footer className="page-footer">
          <span>DOHOONDO / WEB INDEXER</span>
          <span>Built for curious trails</span>
        </footer>
      </section>
    </main>
  );
}

function SearchPage() {
  const [word, setWord] = useState("");
  const [results, setResults] = useState<string[]>([]);
  const [state, setState] = useState<DispatchState>("idle");
  const [message, setMessage] = useState("Search across every website Dhoondo has indexed.");

  async function searchWebsites(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const query = word.trim();
    if (!query) {
      setState("error");
      setMessage("Enter a word to search for.");
      return;
    }

    setState("sending");
    setMessage(`Looking for “${query}”…`);
    setResults([]);

    try {
      const response = await fetch(`${API_URL}/search/`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ words: [query] }),
      });
      const result = await response.json().catch(() => null);
      if (!response.ok) {
        throw new Error(result?.message ?? "The search server rejected this request.");
      }

      setResults(result?.websites ?? []);
      setState("success");
      setMessage(
        result?.websites?.length
          ? `Found ${result.websites.length} matching website${result.websites.length === 1 ? "" : "s"}.`
          : `No indexed websites contain “${query}”.`,
      );
    } catch (error) {
      setState("error");
      setMessage(
        error instanceof Error ? error.message : "Could not reach the search server.",
      );
    }
  }

  return (
    <main className="app-shell search-shell">
      <section className="workbench">
        <header className="topline">
          <a className="breadcrumb back-link" href="/">
            <ArrowLeft size={15} />
            Back to dispatch
          </a>
          <span className="server-link">
            <Radio size={14} /> Search index
          </span>
        </header>

        <div className="search-content">
          <p className="kicker">The Dhoondo index</p>
          <h1>
            Find a thread
            <br />
            <em>worth following.</em>
          </h1>
          <p className="intro-copy">
            Search the words collected from the websites you have sent into the
            crawler.
          </p>

          <form className="search-form" onSubmit={searchWebsites}>
            <label htmlFor="search-word">Search for a word</label>
            <div className="search-input-row">
              <input
                id="search-word"
                value={word}
                onChange={(event) => {
                  setWord(event.target.value);
                  if (state !== "idle") {
                    setState("idle");
                    setMessage("Search across every website Dhoondo has indexed.");
                  }
                }}
                placeholder="e.g. architecture"
                autoComplete="off"
                autoFocus
              />
              <button type="submit" disabled={state === "sending"}>
                {state === "sending" ? (
                  <LoaderCircle className="spin" size={17} />
                ) : (
                  <Search size={17} />
                )}
                Search
              </button>
            </div>
          </form>

          <section className={`search-results search-results--${state}`} aria-live="polite">
            <div className="results-heading">
              <span>Results</span>
              {results.length > 0 && <strong>{results.length.toString().padStart(2, "0")}</strong>}
            </div>
            <p className="result-copy">{message}</p>
            {results.length > 0 && (
              <ul>
                {results.map((url) => (
                  <li key={url}>
                    <a href={url} target="_blank" rel="noreferrer">
                      <Globe2 size={16} />
                      <span>{url}</span>
                      <ArrowUpRight size={15} />
                    </a>
                  </li>
                ))}
              </ul>
            )}
          </section>
        </div>
      </section>
    </main>
  );
}

export default App;
