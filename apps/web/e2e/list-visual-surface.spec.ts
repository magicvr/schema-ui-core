import { expect, test, type Locator, type Page } from "@playwright/test";

import { openSidebarGroup, signInAsAdmin } from "./sign-in";

/**
 * R6 list-surface visual guard (GOAL-009; closes GOAL-007 A-002 F-003).
 *
 * WHY THIS SPEC EXISTS
 * --------------------
 * R6 (list-page visual alignment) shipped through three correction rounds
 * (C5 → C7 → C8) and **every regression in those rounds was reported by the
 * user, not caught by a test**. The reason is structural: the R6 contract is
 * about real layout — computed heights, the adjacency of a button to its input,
 * and which controls a responsive breakpoint hides. jsdom has no layout engine
 * and stubs `matchMedia`, so the jsdom suites can only assert class strings;
 * they cannot observe the properties the contract is actually about.
 *
 * This spec asserts those properties in a real browser, on a real page, through
 * the real seed. It is deliberately *relational* rather than pixel-exact
 * (equal heights, negative overlap, "an item is really hidden") so that
 * retuning tokens or spacing does not produce noise — only breaking the
 * contract does.
 *
 * COVERAGE MAP (contract → assertion)
 *   C7 / A-003 · W13 T-03  search button paired with its keyword input
 *   C7 item 1              icons on the columns / save-view triggers
 *   C7 item 3              the view form renders in the view surface, near its trigger
 *   C8 item 1              every page-action control shares one height
 *   C8 item 2              the expand toggle wears the `control` token, not the reset style
 *   C8 item 3              the toggle exists only when the collapsed row hides a filter
 *   C5 / D-002             filter panel → page actions → list, footer inside the list surface
 *   GOAL-011               page-size control shows the effective size; 10 reaches the API
 *   GOAL-011               the jump confirm action is labelled as a jump
 */

const DESKTOP = { width: 1440, height: 900 };
/** `sm` tier: the collapsed row holds one control, so two filters overflow. */
const NARROW = { width: 700, height: 900 };

/** Signs in and lands on the roles list (which owns a search form + toolbar). */
async function openRolesList(page: Page, viewport = DESKTOP): Promise<void> {
  await page.setViewportSize(viewport);
  await signInAsAdmin(page);
  // The sidebar is `hidden lg:block`, so navigate at desktop width first and
  // resize afterwards when a narrow tier is under test.
  await openSidebarGroup(page, "manifest.nav.group.identityAccess");
  await page.getByRole("link", { name: "Roles" }).click();
  await expect(page.getByRole("heading", { name: "Roles" })).toBeVisible();
  await expect(page.locator('[data-table-surface="true"]')).toBeVisible();
}

/** Resolves a CSS custom property to the value the browser actually uses. */
async function readToken(page: Page, name: string): Promise<string> {
  return page.evaluate(
    (token) => getComputedStyle(document.documentElement).getPropertyValue(token).trim(),
    name,
  );
}

function searchForm(page: Page): Locator {
  return page.locator('[data-form-search-mode="true"]');
}

test.describe("R6 list-surface visual contract", () => {
  test("keeps the keyword search control, page-action heights and list layout intact", async ({
    page,
  }) => {
    await openRolesList(page, DESKTOP);

    // ── C7 / A-003 · W13 T-03: the search submit button is ONE unit with its
    // keyword input. C5 moved it into the filter action cell and broke this;
    // jsdom could only see the class names, the browser sees the geometry.
    const form = searchForm(page);
    const keyword = form.locator('input[type="text"]').first();
    const submit = form.locator('button[type="submit"]').first();
    await expect(keyword).toBeVisible();
    await expect(submit).toBeVisible();

    const pairing = await page.evaluate(() => {
      const root = document.querySelector('[data-form-search-mode="true"]');
      const input = root?.querySelector('input[type="text"]') ?? null;
      const button = root?.querySelector('button[type="submit"]') ?? null;
      if (input === null || button === null) {
        return null;
      }
      const i = input.getBoundingClientRect();
      const b = button.getBoundingClientRect();
      return {
        sameCell:
          button.closest("[data-filter-item]") === input.closest("[data-filter-item]"),
        inActionCell: button.closest('[data-filter-actions="true"]') !== null,
        horizontalGap: Math.round(b.left - i.right),
        verticalOffset: Math.round(b.top - i.top),
        inputClass: input.className,
        buttonClass: button.className,
      };
    });
    expect(pairing).not.toBeNull();
    expect(pairing!.sameCell, "submit button must live in the keyword input's own cell").toBe(
      true,
    );
    expect(pairing!.inActionCell, "submit button must not sit in the filter action cell").toBe(
      false,
    );
    // Overlap by ~1px (`-ml-px`) or butt exactly; never a gap.
    expect(pairing!.horizontalGap, "button must be attached to the input").toBeLessThanOrEqual(0);
    expect(Math.abs(pairing!.verticalOffset), "pair must sit on one line").toBeLessThanOrEqual(1);
    expect(pairing!.buttonClass).toContain("rounded-l-none");
    expect(pairing!.inputClass).toContain("rounded-r-none");

    // ── C8 item 1: one control height across the whole page-actions row.
    const heights = await page.evaluate(() => {
      const row = document.querySelector("[data-list-page-actions]");
      if (row === null) return null;
      return {
        children: Array.from(row.children).map((child) =>
          Math.round(child.getBoundingClientRect().height),
        ),
        toolbarButtons: Array.from(row.querySelectorAll("button")).map((button) =>
          Math.round(button.getBoundingClientRect().height),
        ),
        columnsTrigger: (() => {
          const el = document.querySelector('[data-saved-view-columns-trigger="true"]');
          return el === null ? null : Math.round(el.getBoundingClientRect().height);
        })(),
      };
    });
    expect(heights).not.toBeNull();
    const allHeights = [...heights!.children, ...heights!.toolbarButtons];
    expect(allHeights.length).toBeGreaterThan(1);
    expect(
      new Set(allHeights).size,
      `page-action controls must share one height (got ${allHeights.join(", ")})`,
    ).toBe(1);
    expect(heights!.columnsTrigger).toBe(allHeights[0]);

    // ── C7 item 1: the reference page pairs these triggers with a glyph.
    await expect(page.locator('[data-saved-view-columns-trigger="true"] svg')).toHaveCount(1);
    await expect(page.locator('[data-saved-view-action="save"] svg')).toHaveCount(1);

    // ── C5 / D-002: filter panel → page actions → list surface, and the
    // pagination footer lives INSIDE the list surface (not floating below it).
    const layout = await page.evaluate(() => {
      const panel = document.querySelector('[data-list-filter-panel="true"]');
      const actions = document.querySelector("[data-list-page-actions]");
      const surface = document.querySelector('[data-table-surface="true"]');
      const footer = document.querySelector('[data-table-footer="true"]');
      const pagination = document.querySelector('[data-pagination-footer="true"]');
      const follows = (a: Element | null, b: Element | null) =>
        a !== null && b !== null && (a.compareDocumentPosition(b) & Node.DOCUMENT_POSITION_FOLLOWING) !== 0;
      const above = (a: Element | null, b: Element | null) =>
        a !== null && b !== null && a.getBoundingClientRect().bottom <= b.getBoundingClientRect().top + 1;
      return {
        panelBeforeActions: follows(panel, actions) && above(panel, actions),
        actionsBeforeSurface: follows(actions, surface) && above(actions, surface),
        footerInsideSurface: surface !== null && footer !== null && surface.contains(footer),
        hasPagination: pagination !== null,
      };
    });
    expect(layout.panelBeforeActions, "filter panel must precede page actions").toBe(true);
    expect(layout.actionsBeforeSurface, "page actions must precede the list").toBe(true);
    expect(layout.footerInsideSurface, "list footer must render inside the list surface").toBe(true);
    expect(layout.hasPagination, "a valid list response always shows pagination").toBe(true);

    // ── C8 item 3: two filters + the action cell fit the lg row, so nothing is
    // hidden and the toggle must be absent (a single-row set has nothing to expand).
    await expect(page.locator('[data-filter-toggle="true"]')).toHaveCount(0);
    const hiddenAtDesktop = await page.evaluate(
      () =>
        Array.from(document.querySelectorAll("[data-filter-item]")).filter(
          (el) => getComputedStyle(el).display === "none",
        ).length,
    );
    expect(hiddenAtDesktop, "no filter may be hidden while the toggle is absent").toBe(0);

    // ── C7 item 3: the view form opens in the view surface, beside its trigger
    // — never in the page-actions row far from the button that opened it.
    await page.locator('[data-saved-view-action="save"]').click();
    await expect(page.locator("#roles-table-saved-view-name")).toBeVisible();
    const viewForm = await page.evaluate(() => {
      const surface = document.querySelector('[data-saved-views-surface="true"]');
      const management = document.querySelector('[data-saved-view-management="true"]');
      const pageActions = document.querySelector("[data-list-page-actions]");
      const save = document.querySelector('[data-saved-view-action="save"]');
      const box = (el: Element | null) => {
        if (el === null) return null;
        const r = el.getBoundingClientRect();
        return { top: Math.round(r.top), bottom: Math.round(r.bottom) };
      };
      return {
        inSurface: surface !== null && management !== null && surface.contains(management),
        inPageActions: pageActions !== null && management !== null && pageActions.contains(management),
        save: box(save),
        management: box(management),
      };
    });
    expect(viewForm.inSurface, "view form must render inside the view surface").toBe(true);
    expect(viewForm.inPageActions, "view form must not render in the page-actions row").toBe(false);
    // "Close to its trigger" as a rule, not a magic number: the form starts
    // within one control-height of the button that opened it.
    expect(viewForm.management!.top - viewForm.save!.bottom).toBeLessThanOrEqual(40);
  });

  test("shows the expand toggle only when a filter is really hidden, wearing the control token", async ({
    page,
  }) => {
    await openRolesList(page, DESKTOP);

    // Nothing hidden at desktop → no toggle.
    await expect(page.locator('[data-filter-toggle="true"]')).toHaveCount(0);

    // ── C8 item 3: narrow the viewport until the collapsed row hides a filter.
    await page.setViewportSize(NARROW);
    const toggle = page.locator('[data-filter-toggle="true"]');
    await expect(toggle).toHaveCount(1);
    await expect(toggle).toBeVisible();

    const narrow = await page.evaluate(() => {
      const items = Array.from(document.querySelectorAll("[data-filter-item]"));
      return {
        hidden: items.filter((el) => getComputedStyle(el).display === "none").length,
        visible: items.filter((el) => getComputedStyle(el).display !== "none").length,
      };
    });
    expect(narrow.hidden, "narrow tier must actually hide a filter").toBeGreaterThan(0);
    expect(narrow.visible, "narrow tier still shows the first control").toBeGreaterThan(0);

    // ── C8 item 2: the toggle is a `control` chip, deliberately distinct from
    // the plain bordered reset action beside it. Asserting "equals the token"
    // (not a literal colour) keeps this stable across token retuning.
    const controlToken = await readToken(page, "--control");
    const controlForeground = await readToken(page, "--control-foreground");
    expect(controlToken, "--control must be declared").not.toBe("");
    expect(controlForeground, "--control-foreground must be declared").not.toBe("");

    const surfaces = await page.evaluate(() => {
      const actions = document.querySelector('[data-filter-actions="true"]');
      const toggleEl = document.querySelector('[data-filter-toggle="true"]');
      const reset = actions === null
        ? null
        : Array.from(actions.querySelectorAll("button")).find(
            (b) => (b as HTMLButtonElement).type !== "submit" && b !== toggleEl,
          ) ?? null;
      const read = (el: Element | null) =>
        el === null ? null : getComputedStyle(el).backgroundColor;
      return {
        toggle: read(toggleEl),
        reset: read(reset),
        resetIsToggle: reset !== null && reset === toggleEl,
      };
    });
    expect(surfaces.resetIsToggle, "the reset action must be a distinct control").toBe(false);
    expect(surfaces.toggle, "the toggle must paint a background").not.toBeNull();
    expect(surfaces.toggle, "the toggle must use the --control token").toBe(controlToken);
    expect(
      surfaces.toggle,
      "the toggle must not read as the plain reset action",
    ).not.toBe(surfaces.reset);

    // ── The dark theme overrides the token (C8 defines a recessed dark value);
    // asserting the override EXISTS keeps the token two-layer contract honest.
    const darkControl = await page.evaluate(() => {
      document.documentElement.classList.add("dark");
      const value = getComputedStyle(document.documentElement)
        .getPropertyValue("--control")
        .trim();
      document.documentElement.classList.remove("dark");
      return value;
    });
    expect(darkControl, "--control must have a dark override").not.toBe("");
    expect(darkControl, "--control dark value must differ from the light value").not.toBe(
      controlToken,
    );

    // ── Page actions keep their uniform height at the narrow tier too.
    const headerHeights = await page.evaluate(() => {
      const trigger = document.querySelector('[data-saved-view-columns-trigger="true"]');
      return trigger === null ? null : Math.round(trigger.getBoundingClientRect().height);
    });
    expect(headerHeights).not.toBeNull();
  });

  /**
   * GOAL-011: the pager's page-size control must show the size that is actually
   * in effect, and every other option must reach the wire. The defect this
   * covers (client default 10 vs server default 20, with the parameter omitted
   * whenever the two matched) was invisible to jsdom tests of the day because
   * their mock modelled the wrong contract — and it was the user who saw it in
   * the browser. So assert it here, against the real API: the first list request
   * carries no `pageSize` (the server's own 20 applies), the control displays
   * 20, and choosing 10 sends `pageSize=10`.
   */
  test("shows the effective page size and makes 10 reach the API", async ({ page }) => {
    const listRequests: string[] = [];
    page.on("request", (request) => {
      const url = request.url();
      if (url.includes("/api/roles") && !url.includes("pageSize=100")) {
        listRequests.push(url);
      }
    });

    await openRolesList(page, DESKTOP);

    const sizeSelect = page.locator("select[data-pagination-page-size]");
    await expect(sizeSelect).toHaveValue("20");
    expect(
      listRequests.some((url) => new URL(url).searchParams.has("pageSize")),
      "the default page size must be omitted so the API's own default applies",
    ).toBe(false);

    const listRequest = page.waitForRequest(
      (request) => request.url().includes("/api/roles") && request.url().includes("pageSize=10"),
    );
    await sizeSelect.selectOption("10");
    await listRequest;
    await expect(sizeSelect).toHaveValue("10");

    // The confirm action of the go-to-page form reads as a jump, not a search
    // (zh "跳转" / en "Go"), while the form itself stays "go to page".
    const jump = page.locator('[data-pagination-jump="true"] button[type="submit"]');
    await expect(jump).toHaveText(/^(跳转|Go)$/);
    await expect(page.locator('[data-pagination-jump="true"] label')).toHaveText(
      /^(跳至页|Go to page)$/,
    );
  });
});
