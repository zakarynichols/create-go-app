import test, { expect } from "@playwright/test";

test("Server responds with a school", async ({ page }) => {
  await page.goto("/"); // (*NOTE*: localhost:7777 -> hostname:7777 -> node:7777) 'node' is the name of the service in docker-compose

  // Initial state
  await expect(page.getByText("School ID:", { exact: true })).toBeVisible();
  await expect(page.getByText("School Name:", { exact: true })).toBeVisible();
  await expect(page.getByText("Location:", { exact: true })).toBeVisible();
  await expect(page.getByText("Type:", { exact: true })).toBeVisible();

  // Request school with id: 1
  await page.getByLabel("School ID").click();
  await page.getByLabel("School ID").fill("1");
  await page.getByRole("button", { name: "Get School" }).click();

  // await expect(page.getByText("failed")).toBeVisible(); // -> test if failed

  // New state from server
  await expect(page.getByText("School ID: 1")).toBeVisible();
  await expect(page.getByText("School Name: University of XYZ")).toBeVisible();
  await expect(page.getByText("Location: City A")).toBeVisible();
  await expect(page.getByText("Type: public")).toBeVisible();
});
