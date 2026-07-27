# Changelog

All notable changes to AstraMind are documented in this file.

The project follows [Semantic Versioning](https://semver.org/).
---
---
## [v0.9.2] - 2026-07-24

### Highlights

Started as four related changes (a licensing decision, `.docx` import, single-fact precision, and closing a web/CLI reliability gap), then grew substantially after real manual testing on real Ollama surfaced three further correctness issues in `/kb ask` - each found, root-caused, and fixed before release rather than shipped and patched later. The `scripts/`/`tests/` project structure was also reorganized for clarity during this cycle.

### Changed

- **Relicensed from AGPL-3.0 to Apache License 2.0.** AGPL's network-copyleft terms are a common, explicit reason corporate legal teams block adoption of a dependency outright. Apache 2.0 removes that friction, better matching AstraMind's goal of being a widely adoptable open platform. See `README.md`'s License section for further detail.
- **Reorganized `scripts/` and `tests/` for clarity.** Every runnable action now lives in `scripts/` (build, test, coverage, regression, and behavioral checks); `tests/` holds only fixture data (`tests/fixtures/`), generated output (`tests/output/`), and planning docs (`tests/docs/`) - nothing runnable. Two genuinely duplicate scripts were removed, and one retired script (`run_all.sh`, superseded by `scripts/regression.sh`) was deleted rather than kept as dead weight. `internal/testutil` was also renamed to `internal/utilityforunittest` - the old name was ambiguous about whether it meant "utility for unit tests" (correct) or "tests for a util package" (no such package exists).

### Added

- **`.docx` (modern Word) import support**, with zero new dependencies - a `.docx` file is a ZIP archive of XML, so Go's standard library (`archive/zip` + `encoding/xml`) is sufficient to extract real paragraph text directly, with no external library and no CGO.
- `extractDocxText`, wired into `ImportDocument` alongside the existing `.txt`/`.md` handling
- Regression tests run against a real python-docx-generated file (`testdata/sample_contract.docx`), not a hand-built XML fixture - matching the standard set earlier in this project's history of testing against real files rather than idealized ones (see the CRLF chunking fix in v0.9.1)
- **Precise single-line extraction for single-fact `/kb ask` questions**, refining the whole-chunk return that v0.9.1 shipped. Within the chunk already selected as most relevant, `findPreciseLine` looks for a single paragraph containing an unambiguous highest count of the question's content words (stopwords stripped) - a deterministic, explainable, literal word-overlap count, not a semantic technique. This succeeds on exactly the case that defeated an earlier embedding-similarity-based approach: asking for a Zoom meeting ID no longer risks matching "Not meeting on 16 February" (which shares only one content word with the question) over the actual "Meeting ID 795 777 3585" (which shares two) - a clear, checkable margin instead of a similarity score that measured data showed could rank the wrong sentence higher. Falls back to the existing whole-chunk return on any tie or zero-overlap case, so this can only ever narrow an already-correct answer, never produce a wrong one.
- `ai.ChatRequest`/embedding requests now accept an optional model override (`OPENAI_EMBEDDING_MODEL` env var). Previously the embedding model was silently hardcoded to `text-embedding-3-small` regardless of what chat model (`OPENAI_MODEL`) was configured - meaning switching providers appeared to do nothing to retrieval quality, with no visibility that a separate, unconfigured, separately-billed model was in use the whole time. Default behavior is unchanged if the new variable is unset.

### Fixed

- **The web UI (`/api/ask`) was silently still using the old free-form-LLM-only answer path**, never updated when v0.9.1 redesigned `/kb ask`'s reliability model (deterministic enumeration/extraction, see v0.9.1 above). This wasn't caught by any test, because `webserver.go`'s ask handler and the CLI's ask handler were two entirely separate, hand-written implementations with no shared code and nothing to flag the drift between them - the web UI's users were getting the older, less reliable free-form LLM behavior (including internal prompt-formatting artifacts like `[Source 1 of 2, Document: ...]` leaking into visible answers) with no indication anything was different from the CLI. Fixed by extracting a single `chat.Service.Ask` method that both the CLI and the web UI now call - the same class of drift is now structurally impossible, since there is only one implementation of the routing logic to be wrong.
- **"What is the [plural noun]" questions were misrouted to single-fact extraction instead of enumeration.** `"what is the timings of sanskrit classes?"` matched the `"what is the"` single-fact pattern and returned only 4 of 9 real entries from one chunk, dropping the rest silently - found via real manual testing, not a synthetic case. Fixed with a narrowly-scoped override: specifically for the `"what is the"` pattern (the one pattern in the list genuinely ambiguous between single-fact and enumeration), if the word immediately following "the" looks plural, the question is treated as enumeration instead. Verified this does not regress the pattern's real, working case (`"what is the zoom meeting id"` - "zoom" is not plural - still correctly routes to single-fact).
- **`/kb ask` could confidently answer completely out-of-scope questions using unrelated content.** `SemanticSearch` has no minimum relevance threshold - it always returns its top-ranked chunks regardless of whether any of them are actually related to the question, so a query like "what is the capital of delhi" against a Sanskrit class schedule returned that schedule as if it answered the question. Fixed with `HasAnyContentWordOverlap`, a relevance gate applied before either the enumeration or single-fact path: if the question shares zero content words with anything retrieved, the answer is "no relevant knowledge found" rather than a confidently-presented irrelevant chunk. Deliberately built as literal word-overlap rather than a numeric embedding-similarity cutoff - real measured data from earlier work (v0.9.1) already showed cosine similarity does not reliably separate related from unrelated content for this embedding model, and guessing a new threshold here would risk repeating that exact mistake.
- **A real false positive was found and fixed in the relevance gate above, live, before release**: the word "id" (2 characters) matched as a literal substring inside "confidentiality", letting an unrelated service contract pass the gate for the question "what is the zoom meeting id" even though the contract has no Zoom-related content at all. Fixed by excluding words of 2 characters or fewer from the gate's overlap check - short words collide as substrings inside unrelated longer words far too easily to be a meaningful relevance signal.
- **Enumeration answers could bundle content from unrelated documents when the total knowledge base chunk count was at or below the retrieval limit (5).** Confirmed on real Ollama embeddings, not a mock-provider artifact: a KB containing both `Sanskrit1.txt` and an unrelated service contract answered a Sanskrit-specific enumeration question with both documents included, because `SemanticSearch` returns every chunk it has whenever the total is within its limit - the relevance gate above only ever confirmed "something in here is relevant," not "everything in here is relevant." Fixed with `FilterRelevantChunks`, applied only to whole chunks (never to individual items within a chunk - that was tried in an earlier design and correctly rejected, since it silently dropped real, related entries like "Tuesday Chanting" from a Sanskrit-timings answer just because the entry itself didn't repeat the word "Sanskrit"). A chunk with zero content-word overlap is excluded entirely; every chunk that survives is still returned in full, so no entry within a genuinely relevant chunk can be dropped.

### Known Limitations

- **PDF import is not yet implemented.** Requires a genuine pure-Go PDF text-extraction library (no CGO, no external binary, to stay consistent with this project's architecture). Extraction quality will vary by how the source PDF was produced - clean text-based PDFs extract well, complex multi-column layouts may not, and scanned/image-only PDFs have no extractable text at all without OCR, which remains out of scope.
- **Legacy `.doc` (pre-2007 binary Word format) is not planned for pure-Go support.** Mature pure-Go extractors for this format don't exist; the ecosystem's best tools (antiword, LibreOffice headless conversion) are external binaries, which conflicts with this project's pure-Go constraint. Current recommendation: convert legacy `.doc` files to `.docx` before import (a one-click "Save As" in Word) rather than accept a lower-confidence best-effort extractor for a format explicitly meant to reach "100% output" reliability.

### Testing - Coverage Closed in Previously-0% Code Paths

Three specific, previously-uncovered risk areas were closed with real tests, not synthetic smoke tests:

- **`internal/infrastructure/ai/ollama_embedding_test.go`** (`TestOllamaProvider_Embed`, 4 table-driven subtests: valid response, empty embedding array, HTTP 500, malformed JSON) - `Embed()` coverage 0% → 90.9%. Uses `httptest.Server`, not a mocked client field - `OllamaProvider` has no injectable client field to mock into, confirmed by reading the real struct before writing the test.
- **`internal/infrastructure/ai/openai_embedding_test.go`** (`TestOpenAIProvider_Embed`, 4 subtests, plus explicit `Authorization` header assertion) - `Embed()` coverage 0% → 90.9%; package overall 58.4% → 76.0% at this point, 78.9% after the provider-fallback tests below.
- **`internal/infrastructure/ai/provider_manager_fallback_test.go`** (5 tests) - documents a real, verified, previously-untested behavioral asymmetry: `Chat()` permanently switches to the fallback provider on failure (`TestProviderManager_Chat_FailoverIsPermanent` confirms the switch persists across a second call), but `Stream()`/`Embed()` have **no fallback of their own** (`TestProviderManager_Stream_RuntimeFailureHasNoFallback`, `TestProviderManager_Embed_RuntimeFailureHasNoFallback`) - a genuine runtime failure on either propagates directly with no recovery attempt, confirmed by reading the real implementation, not assumed.
- **`internal/features/kb/json_storage_delete_test.go`** (6 tests) - `DeleteDocument`/`DeleteChunks` coverage 50% → 100%. Documents a real, confirmed asymmetry: `DeleteDocument` on a missing file returns `ErrDocumentNotFound` (an error); `DeleteChunks` on a missing file returns `nil` (success) - same underlying condition, different contract. Also covers the genuine (non-"not found") OS error path for both, using a portable technique (a non-empty directory in place of the expected file) rather than `chmod`, since Unix permission bits aren't reliably testable on Windows.
- **`internal/features/kb/filter_relevant_chunks_test.go`** (`FilterRelevantChunks`, 4 tests) - reproduces the real Issue-1 failure case exactly (`TestFilterRelevantChunks_Issue1RealFailureCase`), and includes a dedicated regression guard (`TestFilterRelevantChunks_DoesNotDropRelatedEntriesWithinAChunk`) proving the fix operates at whole-chunk granularity only, never dropping individual related entries the way an earlier, rejected item-level filtering design did.

Combined effect: `internal/infrastructure/ai` 58.4% → 78.9%; `internal/features/kb` 84.7% → 86.4%; overall project coverage 48.7% → 53.2%.

### Cross-Platform Verification - Real Bugs Found Testing on Real Machines

This release was tested by actually running every script on Linux, MINGW64 (git-bash on Windows), and native Windows `cmd.exe` - not just reading the code. That surfaced a genuinely long list of real, previously-undetected bugs:

- **`check_knowledge_base.bat` had a stale path from before the `scripts/`/`tests/` reorganization** - `cd /d "%~dp0..\.."` (two directory levels up) instead of one, and `tests\integration\commands\kb.txt` instead of the current `tests\fixtures\commands\kb.txt`. Root cause of a real, reproduced `"Could not find astramind.exe in D:\Working"` failure (missing the `\AstraMind` folder entirely) during a `regression.bat` run. The `.sh` equivalent had been correctly updated during the reorganization; the `.bat` was simply missed.
- **A `cmd.exe` parser crash in `check_rag_behavior.bat`**: literal parentheses in printed text (`"Web UI (/api/ask) Smoke Test"`), sitting inside an already-open parenthesized `if` block, caused `cmd.exe`'s paren-counting parser to misread where the block ended - manifesting as `"Smoke was unexpected at this time."` Root-caused via a single targeted debug line rather than repeated guessing, then fixed by escaping every literal parenthesis in the file (`^(` / `^)`), not just the two that crashed.
- **Log files silently missing their own trailer content, found and fixed three separate times in three different scripts** (`check_rag_behavior.sh`/`.bat`, `manual_walkthrough.sh`, then `regression.sh`/`.bat`) - each only piped its main command's output to its log file, so the "Expected results" checklist, sanity-scan results, and final summary reached the terminal but never the saved file. Fixed in `.sh` via `exec > >(tee "$LOGFILE") 2>&1` placed before all output, not just wrapping one command; verified by diffing the saved log against the real terminal transcript byte-for-byte (zero differences). Fixed in `.bat` via paired screen+file `echo` statements, since batch has no native `tee`.
- **`regression.sh`/`.bat` printed a hardcoded `"PASS"` for every pipeline step, unconditionally** - the text never reflected a real, checked outcome; it was only ever "safe" because `set -e`/`if errorlevel` guaranteed an early, silent exit before the summary could print on a real failure, meaning a failure produced no summary at all. Replaced with real per-step status tracking (`BUILD_STATUS`/`TEST_STATUS`/`COVERAGE_STATUS`/`KBRAG_STATUS`, each driven by the step's actual exit code); later steps are marked `SKIPPED`, not silently assumed passing, if an earlier one failed. Verified two ways: a real successful run, and a deliberately broken build (bad Go syntax committed on purpose, reverted after) confirming `Build: FAIL` and every downstream step correctly `SKIPPED`, with exit code `1` and `"AstraMind is NOT READY."`
- **`manual_walkthrough.sh` crashed when passed `--web`** - its argument parser (`BIN="${1:-}"`) had no flag recognition at all and silently misread `--web` as a binary path override, then tried to execute a program literally named `--web`. Fixed with the same argument-loop pattern already used in `check_rag_behavior.sh`. (The script's web smoke test already runs unconditionally on every invocation regardless of this flag - `--web` is now a documented no-op here, not a mode selector.)
- **Non-descriptive, inconsistent log filenames.** `check_rag_behavior.sh`/`.bat` previously wrote to a generic `test_output_<value>.log`, indistinguishable at a glance from any other script's output; the `.bat` version's "timestamp" was actually `%RANDOM%%RANDOM%` (two pseudo-random numbers, no date/time information despite the variable name). Every log filename now matches the exact script and platform that produced it (e.g. `check_rag_behavior_win_20260726_152805.log`, `regression_mingw64_nt-10.0-19045_20260726_153127.log`), with a real PowerShell-sourced timestamp replacing the fake one in every `.bat` script.
- **`regression.sh` had no log file of its own at all** - only `check_rag_behavior.sh`'s narrower, correctly-scoped sub-log existed; the outer script's own `[1/4]`-`[4/4]` step output and final summary were never saved anywhere, only ever shown on the terminal. Fixed with the same `exec > >(tee ...)` pattern, producing a full-run log distinct from (and encompassing) the nested sub-log. `regression.bat`'s equivalent required a different technique - batch has no native `tee` for a multi-process run like this - implemented via a self-relaunch through PowerShell's `Tee-Object`.
- **A stale `.gitignore`** referenced the pre-reorganization `tests/coverage/*` path, meaning `tests/output/coverage/coverage.html`, `.txt`, and `package_coverage.txt` were being silently tracked and re-committed on every regression run. Fixed, and `reports/*.xml` (this release's new machine-readable reports) added preemptively to avoid the identical mistake.
- **Removed genuinely dead code**: `config.Config` (an unused struct in `internal/infrastructure/config`) was never instantiated anywhere in the real application - confirmed by search before removal, not assumed. The package's real, actively-used constants (`Version`, `MaxMessages`, `HistoryFile`) and their tests were kept.

### Added - Cross-Platform Task Runner & Machine-Readable Reporting

- **`cmd/dev`**, a Go-native, cross-platform task runner (`go run ./cmd/dev -run=build|test|coverage|junit|regression-report|regression|clean`) - eliminates the exact class of bug listed above (a `.sh`/`.bat` pair silently diverging) by making build/test/coverage logic exist in exactly one place, callable identically from bash and `cmd.exe`, with zero new toolchain dependency on either platform (Go is already required). `scripts/build.sh`/`.bat`, `test.sh`/`.bat`, and `coverage.sh`/`.bat` are now thin wrappers into it.
- **`reports/junit.xml`** (`go run ./cmd/dev -run=junit`) - parses `go test -json` output directly into real JUnit XML, entirely in Go standard library (`encoding/json` + `encoding/xml`); no external tool (`go-junit-report`, `gotestsum`) required, consistent with this project's pure-Go principle. Verified against the real test suite (171 tests, 12 packages) and against a deliberately broken test, including a genuine cross-package cascading failure (breaking a shared `internal/features/kb` function correctly surfaced a failure in `internal/engine` too) captured correctly.
- **`reports/regression.xml`** (`go run ./cmd/dev -run=regression-report`) - the 4 regression pipeline steps themselves, represented as JUnit test cases using the real captured status from `regression.sh`/`.bat`, with proper `<skipped/>` (not just `<failure/>`) semantics for steps that never ran.
- **CI workflow rebuilt** (`.github/workflows/go.yml`): switched from duplicating `go fmt`/`vet`/`build`/`test` directly to calling `cmd/dev`; fixed the same stale `tests/coverage/` path found in `.gitignore`; wired in `mikepenz/action-junit-report` so results render directly on each commit/PR instead of only a raw log; **added a Windows job** - the previous workflow was Linux-only, meaning it would have caught none of the real `.bat` bugs listed above. First real run: 171/171 tests passing on both Linux and Windows.

---


## [v0.9.1] - 2026-07-24

### Highlights

This release closes out the v0.9.1 validation branch: a hardware check on `gemma2:9b` (issue #55) that, in the course of validating it, surfaced a real chunking bug, a real prompt-construction gap, and a hard limit on how far prompt engineering alone can make free-form LLM enumeration reliable. Rather than continuing to tune prompt wording and generation temperature, `/kb ask` was split into two paths: a deterministic, zero-LLM-call extraction path for questions with a discrete, verifiable answer, and the existing free-form RAG path, now used only as a fallback. Every fix below was verified against a real, previously-used document (`Sanskrit1.txt`, CRLF line endings included) through the actual CLI and real Ollama embeddings - not synthetic test fixtures alone.

### Added

#### Deterministic Query Answering

- `kb.IsEnumerationQuery` - classifies a question as enumeration-style ("what are all the X", "what are the X timings") versus a single specific fact, used to route `/kb ask`
- `kb.ExtractItems` / `kb.BuildListAnswer` - deterministic enumeration path: retrieves chunks as usual, returns each retrieved chunk's content as one item, formats as a list. No LLM call, so no possibility of a matching entry being silently dropped during generation
- `kb.ExtractiveAnswer` - deterministic single-fact path: ranks chunks by embedding similarity to the question, returns the single best-matching chunk's content verbatim. No LLM paraphrase step, so no possibility of a generative model misstating a date, fee, or other fact while rewording it
- `ai.ChatRequest.Temperature` (optional, `*float64`) - configurable generation temperature, used at 0.4 for the RAG fallback path

#### Testing

- Content-fidelity and determinism checks in `tests/integration/manual_testing.sh`: imports a fixture with known facts, asks the same question multiple times, and greps the transcript for every expected fact rather than eyeballing the output
- `TestChunkHandlesCRLFLineEndings`, `TestChunkRespectsParagraphBoundaries` - regression tests for the chunking fixes below, using real CRLF content
- `TestExtractItems_RealSanskritDocument`, `TestExtractiveAnswer_SelectsCorrectChunk` - run against the real `Sanskrit1.txt` bytes, not a synthetic fixture

### Fixed

- **Chunker corrupted real documents mid-word.** Byte-offset chunking had no word or paragraph awareness and could split a chunk boundary through the middle of a word (e.g. "Youth group" -> "uth group"). Fixed with paragraph-aware splitting on blank-line boundaries, falling back to the old sliding-window split only for a single paragraph too large to fit in one chunk.
- **The paragraph-aware fix above silently did not apply to real-world CRLF files.** The real test document uses `\r\n` line endings; a CRLF blank line (`\r\n\r\n`) never contains the substring `\n\n`, so the paragraph splitter found zero boundaries in it and fell straight back to the byte-offset splitter - reproducing the exact bug it was meant to fix, invisibly, because every test fixture at the time used LF-only content. Fixed by normalizing `\r\n`/`\r` to `\n` before splitting.
- **`/kb ask` silently omitted valid entries on enumeration questions, even with a complete and correct prompt.** Root-caused through a full pipeline audit (disk read -> chunk -> embed -> prompt -> generation): retrieval, chunking, and prompt content were all confirmed correct and complete in isolation. The remaining variance was in the LLM's free-form generation step - a model reliably following the question's literal keywords over weaker instruction wording, especially at low temperature. Not fixable through prompt wording, question rewording, or temperature tuning alone (several approaches were tried and measured; see Known Limitations). Resolved architecturally: see "Deterministic Query Answering" above.
- `BuildSemanticPrompt` had regressed to its pre-fix form (a single generic trailing instruction instead of explicit per-source enumeration) due to an unrelated file being committed under its name in a prior commit. Restored.

### Known Limitations

- **Single-fact and enumeration answers return whole chunks verbatim, which can be verbose when a chunk bundles multiple unrelated entries.** Two windowing approaches were tried and abandoned before landing on whole-chunk return:
  - A fixed-size window (N sentences before/after the matched sentence) sometimes bled into a genuinely unrelated neighboring entry.
  - A dynamic, embedding-similarity-threshold window (expand while neighboring sentences stay semantically similar to the match) was tried next, on the theory that a real topic boundary would show up as a similarity drop. Live testing against real Ollama embeddings disproved this for short-sentence prose: an unrelated sentence ("Not meeting on 16 February", about a different class) scored *higher* similarity to the anchor ("Meeting ID 795 777 3585") than genuinely related sentences did (the actual Zoom URL and password), apparently because both happened to share the literal word "meeting" despite unrelated meaning. No threshold could separate on-topic from off-topic content using this signal - this was a real, measured finding, not a miscalibration. Whole-chunk return was kept as the safer default: chunking already guarantees no entry is corrupted or split mid-content, so the worst failure mode is extra (correct) text, never wrong or missing text.
- `ExtractiveAnswer` re-embeds every sentence in every retrieved chunk on every single-fact question, with no caching. Fine at the scale tested (single-digit chunk counts); will need embedding caching at import time before this scales to a larger knowledge base.
- The free-form LLM RAG path (`BuildSemanticPrompt` + `Chat`) is now only used as a fallback when no embedder is configured. It is not covered by the completeness/precision guarantees of the deterministic paths above.
- `internal/features/kb/query_expansion.go`'s `ExpandQuery` function (question-rewording via appended instructions) is superseded by the deterministic extraction path and is no longer called anywhere in the codebase. `IsEnumerationQuery` from the same file is still in active use as the query router. Cleanup (removing the now-dead `ExpandQuery` code) is planned but not yet done.

### Tested

- go fmt
- go vet
- go build
- go test -v ./...
- tests/integration/run_all.sh
- tests/integration/manual_testing.sh (extended with content-fidelity and determinism scans)
- Full build and test suite re-verified against a fresh pull of the actual pushed branch from GitHub, independent of local working-directory state

### Verified

- Chunking fix confirmed against the real `Sanskrit1.txt` file's actual bytes (CRLF line endings, real diacritics) - zero corruption, all 9 real entries intact
- Enumeration path (`/kb ask what are the Sanskrit class timings`) confirmed live via the CLI against real Ollama embeddings: all 9 entries present, correctly cited across 3 sources
- Single-fact path (`/kb ask what is the meeting id`) confirmed live via the CLI against real Ollama embeddings: correct chunk selected, single source cited, correct answer present, no fabrication
- Hardware validation (issue #55) closed: `gemma2:9b` produces correct output on the target hardware (Intel i5-4210U, 16GB RAM, no GPU); brief UI stutter observed only under simultaneous heavy multitasking, not disqualifying for sequential use

---

## [v0.9.0] - 2026-07-21

### Highlights

This release turns AstraMind from a keyword-search Knowledge Base into a working Retrieval-Augmented Generation assistant, and closes out a substantial architectural cleanup identified in a full review of the v0.8.0 codebase. It also adds a local, browser-based interface for non-technical users, alongside the existing CLI. A known model-capability limitation was found and documented through controlled testing - see Known Limitations below before treating this as demo-ready for high-stakes use.

### Added

#### Semantic Search & RAG

- `ai.EmbeddingProvider` interface (Ollama, OpenAI, and a deterministic mock for tests), matching the existing `StreamingProvider` pattern
- Embedding generation wired into `/kb import`
- `kb.CosineSimilarity` and `Repository.SemanticSearch`, ranking chunks by embedding similarity rather than keyword count
- `/kb ssearch <text>` - semantic search, kept as a separate command alongside keyword `/kb search`
- `kb.BuildSemanticPrompt` and `/kb ask <question>` - completes the RAG loop (import -> chunk -> embed -> retrieve -> **answer**), citing sources with every response
- Local web UI (`--web` flag): embedded HTTP server and single-file browser interface, exposing import/list/ask over a JSON API, reusing the same backend and provider configuration as the CLI (online or offline)

#### Architecture

- `Dependencies` struct extended with `HistoryService`, `SessionService`, `ExportService`, `SearchService` - every feature service now constructed once at startup and shared, instead of built ad-hoc per command
- `storage.FileHistoryStore` - configurable-directory session storage, mirroring the existing `kb.JSONStorage` pattern
- `history.Store` and `session.Store` interfaces, enabling isolated, test-injectable storage
- `--script` mode routes through the full command dispatcher, matching interactive mode

#### Testing

- First unit tests for `history` and `session` packages (previously untested)
- Semantic search, RAG prompt, and context-window regression test suites
- `tests/integration/manual_testing.sh` - full manual walkthrough including a live keyword-vs-semantic comparison and a `--web` API smoke test

### Improved

- `search` decoupled from `storage`, now depends on `history` (`search -> history -> storage`)
- `session` decoupled from duplicated storage logic, now delegates to `history` for Save/Load/Delete/List
- `README.md` and `docs/roadmap.md` updated to reflect delivered semantic search and RAG

### Fixed

- Silent embedding failure during `/kb import`, caused by two missing provider files - now surfaced with visible import feedback instead of failing invisibly
- Test pollution of the real `data/sessions` folder - storage location is now injectable, and tests use isolated temp directories
- Ollama RAG answers truncating mid-generation: no context window was ever specified, so Ollama fell back to its own small default; requests now set `num_ctx: 8192`
- `SemanticSearch` had no result limit - every embedded chunk in the entire knowledge base was returned on every question with no cap, which is harmless with a few test documents but would stuff the whole knowledge base into every prompt at real scale; capped to the top 5 most relevant chunks
- Broken `.gitignore` line from a prior append with no trailing newline
- Compiled binary (`astramind`) no longer tracked in git

### Removed

- `chat/dispatcher.go` and `chat/script.go` - dead code superseded by dispatcher-based `--script` routing
- `kb.Service` - unused wrapper with zero callers anywhere in the codebase

### Known Limitations

- **`gemma3:1b` (the smallest practical local Ollama model) can fabricate facts and omit valid entries on exhaustive multi-item extraction, even with correct source text directly in its context.** Proven via a controlled experiment: an identical prompt, identical document, and identical question given to `gemma3:1b` versus a ~25B-parameter model (`google/gemma-4-26b-a4b-it` via OpenRouter) - the larger model correctly listed all 5 matching entries with zero errors; the smaller model returned only 2 of 5, and in an earlier run fabricated a duplicate entry with an invented date. This is a model-capability ceiling, not a defect in AstraMind's retrieval, prompt construction, or context handling - all of which were independently verified as correct during this investigation.
- Practical consequence: `/kb ask` and the web UI's answer generation should not be treated as reliable for high-stakes use (e.g. legal research) on `gemma3:1b` or comparably small local models. A real minimum local-model size/hardware floor has not yet been established - testing a mid-size local model (e.g. `gemma2:9b`) against realistic hardware is a recommended next step before any customer-facing demo of the RAG feature.
- `/about` and `/config` still report `v0.8.0` - the version constant has not yet been bumped for this release.

### Tested

- go fmt
- go vet
- go build
- go test -v ./...
- tests/integration/run_all.sh
- tests/integration/manual_testing.sh (including the `--web` API smoke test)

### Verified

- Semantic search and RAG proven end-to-end on real Ollama + `nomic-embed-text`, including a real, previously-unseen document (not a synthetic test fixture)
- No test-created sessions or documents polluting the real `data/` folder after the storage isolation and test-cleanup fixes
- Model-capability limitation isolated via a controlled comparison (see Known Limitations), confirming the issue is not in AstraMind's code

---

## [v0.8.0] - 2026-07-13

### Highlights

This release introduces AstraMind's first Knowledge Base implementation, providing document import, persistent storage, automatic chunking, keyword search, knowledge management commands, and the architectural foundation for Retrieval-Augmented Generation (RAG).

### Added

#### Knowledge Base

- Built-in Knowledge Base framework
- Text document import
- Markdown document import
- Automatic document chunking
- Persistent document storage
- Persistent chunk storage
- Repository abstraction
- Keyword search engine
- Prompt builder for future RAG support
- Knowledge Base statistics
- Knowledge Base management API

#### CLI Commands

- `/kb import <file>`
- `/kb list`
- `/kb search <text>`
- `/kb remove <id>`
- `/kb clear`
- `/kb stats`

#### Testing

- Knowledge Base unit tests
- Repository tests
- Search tests
- Prompt builder tests
- Management tests
- CLI command tests
- End-to-end Knowledge Base workflow validation

### Improved

- Chat service architecture
- Dependency injection for future extensibility
- Storage abstraction
- Repository design
- Internal package organization
- Command dispatch architecture
- Overall extensibility for future RAG integration

### Tested

- go fmt
- go vet
- go build
- go test -v ./...

### Verified

- Document import
- Persistent Knowledge Base storage
- Chunk generation
- Keyword search
- Knowledge Base statistics
- Document removal
- Knowledge Base clearing
- Complete `/kb` command workflow

---

## [v0.7.0] 

### Added

- Native Ollama provider
- Local LLM support through Ollama
- Runtime provider selection
- Streaming support for Ollama
- Streaming integration tests
- HTTP error handling tests
- Invalid JSON handling tests
- Connection failure tests

### Improved

- Provider factory architecture
- Provider manager
- Streaming abstraction
- Request builder reuse
- Test coverage

### Supported Providers

- Mock AI
- OpenAI
- OpenRouter
- Ollama

### Tested

- go fmt
- go vet
- go build
- go test -v ./...

### Verified

- Local Ollama installation
- Gemma3:1b model
- Conversation history
- Session persistence
- Provider switching
- Streaming implementation

---
## v0.6.0 

### Added

- Added `/search` command to search the current conversation.
- Added `/searchall` command to search across all saved sessions.
- Added case-insensitive conversation search.
- Added grouped search results by session.
- Added reusable search renderer.
- Added search models for current and multi-session search.

### Improved

- Improved search result presentation.
- Improved session grouping for multi-session search.

### Fixed

- Fixed conversation persistence after successful AI responses.
- Fixed search across sessions by ensuring conversations are saved immediately.
- Improved session consistency when switching between conversations.

### Tests

- Added unit tests for conversation search.
- Added integration tests for multi-session search.

---
# v0.5.0 

## Highlights

This release introduces real-time streaming AI responses, OpenRouter integration, configurable OpenAI-compatible endpoints, improved provider architecture, and a comprehensive integration testing framework. AstraMind now supports streaming and non-streaming providers through a unified interface while significantly improving maintainability and test coverage.

## Added

### Streaming Responses
- Real-time token-by-token AI streaming.
- Streaming renderer for terminal output.
- Streaming provider interface.
- Mock streaming provider.
- Automatic fallback for non-streaming providers.

### OpenAI-Compatible Providers
- OpenRouter integration.
- Configurable API base URL.
- Configurable AI model selection.
- Shared HTTP request builder.
- Improved provider abstraction.

### Integration Testing
- HTTP integration tests using `httptest.Server`.
- Streaming integration tests.
- HTTP error integration tests.
- Invalid JSON response tests.
- End-to-end provider validation without external services.

### Configuration
- `OPENAI_BASE_URL` environment variable.
- Support for OpenAI-compatible APIs.
- Runtime provider configuration improvements.

## Improved

### AI Provider Framework
- Reduced duplicated request construction.
- Cleaner provider implementation.
- Improved request handling.
- Better error propagation.
- More maintainable streaming architecture.

### Testing
- Expanded automated test suite.
- Increased integration test coverage.
- Improved provider validation.
- Enhanced streaming validation.

### Developer Experience
- Cleaner internal architecture.
- Better separation of provider responsibilities.
- Improved code maintainability.
- Simplified future provider integrations.

## Fixed

- Streaming response handling.
- Provider request construction.
- HTTP error handling consistency.
- OpenAI-compatible endpoint support.
- Runtime streaming stability.

## Testing

Validated successfully with:

- Unit tests
- Integration tests
- Streaming integration tests
- HTTP error integration tests
- Mock provider tests
- Mock streaming provider tests
- Provider manager tests
- Renderer tests
- Storage tests
- Runtime validation against OpenRouter

---

# v0.4.1 

## Highlights

This release introduces a major architectural milestone for AstraMind, including a modular AI provider framework, conversation export capabilities, automated testing, and a production-grade GitHub Actions CI pipeline.

## Added

### Export System
- TXT conversation export
- Markdown conversation export
- Automatic export directory creation

### AI Provider Framework
- AI provider abstraction
- Provider factory
- Provider manager
- Mock AI provider
- OpenAI provider
- Automatic provider failover

### Session Management
- Session creation
- Session loading
- Session deletion
- Session listing
- Active session tracking

### Testing
- Storage integration tests
- Mock provider tests
- Regression test suite
- Coverage generation

### CI/CD
- GitHub Actions workflow
- Automatic formatting (`go fmt`)
- Static analysis (`go vet`)
- Automated builds
- Automated unit tests
- Coverage reporting
- Coverage artifact publishing
- Workflow concurrency
- Workflow timeout
- Least-privilege workflow permissions

## Improved

- Modular project architecture
- Export subsystem
- Storage layer
- Test infrastructure
- Repository documentation
- Project roadmap
- README documentation

## Fixed

- Export failures when the export directory did not exist
- GitHub Actions workflow stability
- Cross-platform CI compatibility

---

# v0.4.0

## Added

- Multi-session support
- Session management
- Improved API error handling

---

# v0.3.0

## Added

- Persistent conversation history
- Statistics command
- Configuration command
- About command

---

# v0.2.0

## Added

- Conversation memory
- History support
- Chat commands

---

# v0.1.0

## Initial Release

- AstraMind CLI chatbot
- OpenAI integration
- Environment configuration