import { describe, it, expect } from "vitest";
import { routeLabels } from "@/lib/routes";

describe("routeLabels", () => {
  it("contains dashboard label", () => {
    expect(routeLabels.dashboard).toBe("Dashboard");
  });

  it("contains items label", () => {
    expect(routeLabels.items).toBe("Items");
  });

  it("contains categories label", () => {
    expect(routeLabels.categories).toBe("Categories");
  });

  it("contains forms label", () => {
    expect(routeLabels.forms).toBe("Forms");
  });

  it("contains workflows label", () => {
    expect(routeLabels.workflows).toBe("Workflows");
  });

  it("contains settings label", () => {
    expect(routeLabels.settings).toBe("Settings");
  });
});
