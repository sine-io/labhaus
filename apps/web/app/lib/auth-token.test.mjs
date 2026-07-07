import assert from "node:assert/strict";
import test from "node:test";

import { buildFrontendHeaders } from "./auth-token.mjs";

test("buildFrontendHeaders adds authorization when a token is present", () => {
  assert.deepEqual(buildFrontendHeaders("abc.def.ghi"), {
    "Content-Type": "application/json",
    Authorization: "Bearer abc.def.ghi",
  });
});

test("buildFrontendHeaders accepts an already-prefixed bearer token", () => {
  assert.deepEqual(buildFrontendHeaders("Bearer token-123"), {
    "Content-Type": "application/json",
    Authorization: "Bearer token-123",
  });
});

test("buildFrontendHeaders omits authorization when token is empty", () => {
  assert.deepEqual(buildFrontendHeaders("  "), {
    "Content-Type": "application/json",
  });
});
