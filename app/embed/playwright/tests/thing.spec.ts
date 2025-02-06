import test, { expect } from "@playwright/test";

test("Server responds with a thing", async ({ page }) => {
  await page.goto("/"); // (*NOTE*: localhost:7777 -> hostname:7777 -> node:7777) 'node' is the name of the service in docker-compose

  // Initial state
  await expect(page.getByText("thing ID:", { exact: true })).toBeVisible();

  // Request thing with id: 1
  await page.getByLabel("thing ID").click();
  await page.getByLabel("thing ID").fill("1");
  await page.getByRole("button", { name: "Get thing" }).click();

  // await expect(page.getByText("failed")).toBeVisible(); // -> test if failed

  // New state from server
  await expect(page.getByText("thing ID: 1")).toBeVisible();
});
