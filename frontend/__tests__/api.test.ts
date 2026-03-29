import { describe, it, expect, vi, beforeEach } from "vitest";

// Mock fetch globally
const mockFetch = vi.fn();
vi.stubGlobal("fetch", mockFetch);

// Import after mocking
import { api } from "@/lib/api";

describe("api client", () => {
  beforeEach(() => {
    mockFetch.mockClear();
  });

  it("sends GET requests with credentials", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => ({ data: "test" }),
    });

    const result = await api.get("/api/test");
    expect(result).toEqual({ data: "test" });
    expect(mockFetch).toHaveBeenCalledWith(
      "/api/test",
      expect.objectContaining({
        credentials: "include",
      }),
    );
  });

  it("sends POST requests with JSON body", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => ({ ok: true }),
    });

    await api.post("/api/items", { name: "Widget" });
    expect(mockFetch).toHaveBeenCalledWith(
      "/api/items",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ name: "Widget" }),
      }),
    );
  });

  it("sends PUT requests", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => ({ ok: true }),
    });

    await api.put("/api/items/123", { name: "Updated" });
    expect(mockFetch).toHaveBeenCalledWith(
      "/api/items/123",
      expect.objectContaining({
        method: "PUT",
      }),
    );
  });

  it("sends DELETE requests", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => ({ ok: true }),
    });

    await api.delete("/api/items/123");
    expect(mockFetch).toHaveBeenCalledWith(
      "/api/items/123",
      expect.objectContaining({
        method: "DELETE",
      }),
    );
  });

  it("throws on non-OK responses", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 500,
      json: async () => ({ error: "Server error" }),
    });

    await expect(api.get("/api/fail")).rejects.toThrow("Server error");
  });
});
