# Zod jenny — v2beta1 generation vs handwritten mutation-API schemas

A snapshot comparison of the cog Zod jenny's output for the dashboard
`v2beta1` CUE schema against the handwritten Zod schemas Grafana uses
for the dashboard mutation API.

- **Generated:** [`generated/zod/src/v2beta1/schemas.gen.ts`](../../../generated/zod/src/v2beta1/schemas.gen.ts) (produced by `cog/config/zod_poc.yaml`)
- **Handwritten:** [`grafana/public/app/features/dashboard-scene/mutation-api/commands/schemas.ts`](https://github.com/grafana/grafana/blob/main/public/app/features/dashboard-scene/mutation-api/commands/schemas.ts)

## Scope of this comparison

The two files are not 1:1. The handwritten file mirrors a *subset* of
v2beta1 — the kinds reachable from mutation commands — plus
mutation-specific payload schemas that have no CUE equivalent. 
This doc only compares features that exist in both
files. Things present in one but not the other (mutation payload
schemas, hand-added defaults the CUE source doesn't declare,
hoisted-vs-inlined struct exports) are out of scope.

## What differs

### `z.lazy(() => …)` everywhere

Cog wraps every cross-reference inside the same package in
`z.lazy(() => XSchema)`. The handwritten file has zero `z.lazy` calls
because it declares schemas in topological order and references them
directly as identifiers.

```ts
// generated  – noisy, defeats some inference
hide: z.lazy(() => VariableHideSchema).optional().default("dontHide")

// handwritten – terse, full inference
hide: variableHideSchema  // declared above with .default('dontHide')
```

## Open issues

- ⏳ `z.lazy()` everywhere — do we need topological sort of objects within
  a schema?  - for now I dont see a strong reason why we do
- ⏳ Recursive cycles fail TS strict-mode — even after topo-sort, the
  layout/tab/formatter cycles need explicit `z.ZodType<…>`
  annotations to break inference. Currently produces ~55 TS7022/TS7024
  errors per v2-family file under `tsc --strict`.
- ⏳ Numeric type fidelity — `uint*` doesn't add `.nonnegative()`;
  bit widths don't add `.gte/.lte` constraints; `z.number()` accepts
  `NaN`/`Infinity`. Affects e.g. `revision: uint16` in v2beta1.
- ⏳ `index.ts` generation — bare-package imports (`../<pkg>`) only
  resolve under TS module resolution that finds an index file. Adding
  an `Index{}` jenny mirroring `typescript/index.go` would close the
  gap.

## Reproducing this output

```bash
cd cog
go run ./cmd/cli generate --config config/zod_poc.yaml
# Output: generated/zod/src/v2beta1/schemas.gen.ts
```

For the production-shaped output (per-kind-version subdirectories
written into `@grafana/schema`):

```bash
go run ./cmd/cli generate --config config/zod_grafana_dashboard.yaml
```
