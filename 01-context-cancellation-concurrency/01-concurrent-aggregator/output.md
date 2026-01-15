📚 What I Learned — Kata 01: The Fail-Fast Data Aggregator
🧠 Conceptual Takeaways
Concept	Description
Context Propagation	Used context.WithTimeout(ctx, duration) to ensure time-bounded operations, enabling request-level timeouts and cancellation.
Fail-Fast Concurrency	Leveraged golang.org/x/sync/errgroup to run multiple service calls concurrently, cancelling all on first error — crucial for real-time systems.
Functional Options Pattern	Created a flexible, clean constructor using options like WithTimeout() and WithLogger() — eliminates messy parameter lists.
Structured Logging (slog)	Integrated modern structured logging to emit rich, machine-readable logs — helpful for observability and debugging.
Concurrency-Safe State Sharing	Used sync.Mutex to protect shared data (map[string]string) — avoids data races while aggregating results.
Clean Abstraction Design	Wrapped behavior in an Aggregator struct, exposing a single Aggregate(ctx, id) method — keeps the API minimal and testable.
🌍 Real-World Relevance
🏦 Use Case: Dashboard Backend for a Banking App

Imagine a user opens their banking dashboard:

You fetch Profile Info (e.g., name, KYC status) from Service A

And Recent Transactions or Balance Summary from Service B

Constraints:

You want the total dashboard response in under 1s

If either service fails or is too slow, you must cancel the other

Logs should show where and why things failed, with user ID context

This kata mimics that exact flow:

👨‍💼 Parallel service calls for performance

⏱ Global timeout using context

🧨 Fast failure with cancellation to save backend resources

📊 Structured logs for traceability in real systems

🧪 Testing Strategy Covered

Timeouts: Validated that the aggregator respects the configured timeout

Cancellation: Ensured that one failing service cancels the other

Mocking Services: Simulated delays and errors for robust fail-fast logic

🧾 TL;DR

This kata taught me how to write idiomatic, production-ready Go using:

context for lifecycle control

errgroup for parallel execution with cancellation

slog for structured observability

Functional options for flexible design

Concurrency-safe result collection with minimal race conditions