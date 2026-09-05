import { expect, test } from "@playwright/test";

import { signInAsAdmin } from "./sign-in";

const appProfile = (process.env.APP_PROFILE || "mvp").trim().toLowerCase();

test.use({ viewport: { width: 1440, height: 900 } });

function jsonResponse(body: unknown, status = 200) {
  return {
    status,
    contentType: "application/json",
    body: JSON.stringify(body),
  };
}

function operatorStatus() {
  return {
    configured: true,
    token_set: true,
    secret_set: true,
    mode: "webhook",
    webhook_public_base_url: "",
    connection_state: "running",
    receiver: "webhook",
    business_occupied: false,
    bot_id: 42,
    bot_username: "fixture_bot",
    captured_messages_count: 240,
  };
}

function operatorSessions() {
  return Array.from({ length: 80 }, (_, index) => ({
    chatId: `chat-${index + 1}`,
    chatType: "private",
    title: `Contact ${index + 1}`,
    username: `contact_${index + 1}`,
    lastMessageAt: `2026-09-05T00:${String(index % 60).padStart(2, "0")}:00Z`,
  }));
}

function operatorMessages() {
  return Array.from({ length: 120 }, (_, index) => ({
    chatId: "chat-1",
    direction: index % 3 === 0 ? "outbound" : "inbound",
    status: index % 3 === 0 ? "sent" : "received",
    occurredAt: new Date(Date.parse("2026-09-05T01:00:00Z") + index * 60_000).toISOString(),
    ...(index % 3 === 0 ? { requestId: `operator-${index + 1}` } : { updateId: String(index + 1) }),
    text: `Message ${index + 1} ${"with enough content to wrap inside the transcript without widening the page. ".repeat(3)}`,
  })).reverse();
}

type TelegramFixtureOptions = {
  sessions?: Array<{
    chatId: string;
    chatType: string;
    title?: string;
    username?: string;
    lastMessageAt?: string;
  }>;
  messagesByChatId?: Record<string, Array<Record<string, unknown>>>;
};

function groupAndChannelOperatorSessions() {
  return [
    {
      chatId: "chat-group",
      chatType: "group",
      title: "Support group",
      username: "support",
      lastMessageAt: "2026-09-05T02:02:00Z",
    },
    {
      chatId: "chat-channel",
      chatType: "channel",
      title: "Announcements",
      username: "announcements",
      lastMessageAt: "2026-09-05T03:02:00Z",
    },
  ];
}

function groupAndChannelOperatorMessages(): Record<string, Array<Record<string, unknown>>> {
  const groupMessages = [
    {
      chatId: "chat-group",
      direction: "inbound",
      status: "received",
      occurredAt: "2026-09-05T02:00:00Z",
      updateId: "group-1",
      senderUsername: "alice",
      text: "Message from Alice",
    },
    {
      chatId: "chat-group",
      direction: "inbound",
      status: "received",
      occurredAt: "2026-09-05T02:01:00Z",
      updateId: "group-2",
      text: "Message without a sender username",
    },
    {
      chatId: "chat-group",
      direction: "outbound",
      status: "sent",
      occurredAt: "2026-09-05T02:02:00Z",
      requestId: "group-3",
      text: "Bot response",
    },
  ].reverse();
  const channelMessages = [
    {
      chatId: "chat-channel",
      direction: "inbound",
      status: "received",
      occurredAt: "2026-09-05T03:00:00Z",
      updateId: "channel-1",
      senderUsername: "channel_editor",
      text: "Message from the channel editor",
    },
    {
      chatId: "chat-channel",
      direction: "inbound",
      status: "received",
      occurredAt: "2026-09-05T03:01:00Z",
      updateId: "channel-2",
      text: "Channel message without a sender username",
    },
    {
      chatId: "chat-channel",
      direction: "outbound",
      status: "sent",
      occurredAt: "2026-09-05T03:02:00Z",
      requestId: "channel-3",
      text: "Bot channel response",
    },
  ].reverse();
  return { "chat-group": groupMessages, "chat-channel": channelMessages };
}

async function installTelegramFixtures(
  page: import("@playwright/test").Page,
  options: TelegramFixtureOptions = {},
): Promise<void> {
  const sessions = options.sessions ?? operatorSessions();
  const messagesByChatId: Record<string, Array<Record<string, unknown>>> = options.messagesByChatId
    ?? { "chat-1": operatorMessages() };
  await page.route("**/api/channel/telegram/**", async (route) => {
    const url = new URL(route.request().url());
    if (url.pathname === "/api/channel/telegram/settings") {
      await route.fulfill(jsonResponse(operatorStatus()));
      return;
    }
    if (url.pathname === "/api/channel/telegram/operator/sessions") {
      const items = sessions;
      await route.fulfill(jsonResponse({ items, total: items.length, page: 1, pageSize: 100 }));
      return;
    }
    const match = url.pathname.match(/^\/api\/channel\/telegram\/operator\/sessions\/([^/]+)\/(messages|capability)$/);
    if (match !== null) {
      const chatId = decodeURIComponent(match[1]);
      if (!(chatId in messagesByChatId)) {
        await route.continue();
        return;
      }
      if (match[2] === "capability") {
        await route.fulfill(jsonResponse({ chatId, canSend: true }));
        return;
      }
      if (route.request().method() === "POST") {
        await route.fulfill(jsonResponse({ status: "sent" }, 201));
        return;
      }
      const items = messagesByChatId[chatId];
      await route.fulfill(jsonResponse({ items, total: items.length, page: 1, pageSize: 100 }));
      return;
    }
    await route.continue();
  });
}

test("Telegram operator keeps document/main fixed while sessions and messages scroll internally", async ({ page }) => {
  test.skip(appProfile !== "custom", "requires APP_PROFILE=custom so channel.telegram is enabled");

  await signInAsAdmin(page);
  await installTelegramFixtures(page);
  // Stay inside the authenticated SPA route so the browser measurement covers
  // the same navigation path an operator uses and does not require a second
  // token restoration between the dashboard and the inner page.
  await page.getByRole("link", { name: "Telegram channel" }).click();
  await expect(page.locator("#page-title")).toHaveText("Telegram channel");
  await page.getByRole("button", { name: "Open operator conversations" }).click();

  await expect(page.locator('[data-telegram-operator-page="true"]')).toBeVisible();
  await expect(page.locator("[data-telegram-session='chat-1']")).toBeVisible();
  await expect(page.locator("[data-telegram-message]")).toHaveCount(120);
  await expect(page.locator("[data-telegram-message]").first()).toContainText("Message 1");
  await expect(page.locator("[data-telegram-message]").last()).toContainText("Message 120");
  await expect(page.locator("[data-telegram-sender]").first()).toHaveText("Bot");
  await expect(page.locator("[data-telegram-sender]").last()).toHaveText("Contact 1");
  await expect(page.locator("[data-telegram-shortcut-reverse]")).toBeVisible();
  await expect(page.locator("[data-telegram-shortcut-hint]")).toContainText("Enter sends");

  const metrics = await page.evaluate(() => {
    const required = (selector: string): HTMLElement => {
      const element = document.querySelector<HTMLElement>(selector);
      if (element === null) throw new Error(`missing ${selector}`);
      return element;
    };
    const scrollingElement = document.scrollingElement ?? document.documentElement;
    const root = required('[data-shell="admin"]');
    const body = required('[data-shell-region="body"]');
    const main = required('[data-shell-region="main"]');
    const pageRegion = required('[data-shell-region="page"]');
    const sessions = required("[data-telegram-sessions]");
    const messages = required("[data-telegram-message-list]");
    const operator = required("[data-telegram-operator]");
    const composer = required("[data-telegram-composer]");
    const composerRect = composer.getBoundingClientRect();
    const operatorRect = operator.getBoundingClientRect();
    return {
      viewportHeight: window.innerHeight,
      document: { clientHeight: scrollingElement.clientHeight, scrollHeight: scrollingElement.scrollHeight, scrollTop: scrollingElement.scrollTop },
      body: { clientHeight: body.clientHeight, scrollHeight: body.scrollHeight, scrollTop: body.scrollTop },
      root: { clientHeight: root.clientHeight, scrollHeight: root.scrollHeight },
      main: { clientHeight: main.clientHeight, scrollHeight: main.scrollHeight, scrollTop: main.scrollTop, overflowY: getComputedStyle(main).overflowY },
      page: { clientHeight: pageRegion.clientHeight, scrollHeight: pageRegion.scrollHeight },
      sessions: { clientHeight: sessions.clientHeight, scrollHeight: sessions.scrollHeight, overflowY: getComputedStyle(sessions).overflowY },
      messages: { clientHeight: messages.clientHeight, scrollHeight: messages.scrollHeight, overflowY: getComputedStyle(messages).overflowY },
      messagesAtBottom: messages.scrollHeight - messages.scrollTop - messages.clientHeight <= 1,
      composerWithinOperator: composerRect.bottom <= operatorRect.bottom + 1,
      windowScrollY: window.scrollY,
    };
  });

  expect(metrics.document.scrollHeight).toBeLessThanOrEqual(metrics.document.clientHeight + 1);
  expect(metrics.body.scrollHeight).toBeLessThanOrEqual(metrics.body.clientHeight + 1);
  expect(metrics.root.scrollHeight).toBeLessThanOrEqual(metrics.root.clientHeight + 1);
  expect(metrics.main.scrollHeight).toBeLessThanOrEqual(metrics.main.clientHeight + 1);
  expect(metrics.main.overflowY).toBe("hidden");
  expect(metrics.sessions.scrollHeight).toBeGreaterThan(metrics.sessions.clientHeight);
  expect(metrics.sessions.overflowY).toBe("auto");
  expect(metrics.messages.scrollHeight).toBeGreaterThan(metrics.messages.clientHeight);
  expect(metrics.messages.overflowY).toBe("auto");
  expect(metrics.messagesAtBottom).toBe(true);
  expect(metrics.composerWithinOperator).toBe(true);

  const textarea = page.locator("[data-telegram-composer] textarea");
  await textarea.fill("browser Ctrl+Enter newline");
  await textarea.press("Control+Enter");
  await expect(textarea).toHaveValue("browser Ctrl+Enter newline\n");

  const enterSendRequest = page.waitForRequest((request) =>
    request.url().endsWith("/api/channel/telegram/operator/sessions/chat-1/messages") && request.method() === "POST");
  await textarea.fill("browser Enter send");
  await textarea.press("Enter");
  const enterRequest = await enterSendRequest;
  expect(JSON.parse(enterRequest.postData() ?? "{}")).toEqual(expect.objectContaining({ text: "browser Enter send" }));
  await expect(textarea).toHaveValue("");

  const reverseShortcuts = page.locator("[data-telegram-shortcut-reverse]");
  await reverseShortcuts.check();
  await expect(page.locator("[data-telegram-shortcut-hint]")).toContainText("Ctrl+Enter sends");
  await textarea.fill("browser reversed Enter newline");
  await textarea.press("Enter");
  await expect(textarea).toHaveValue("browser reversed Enter newline\n");

  const ctrlEnterSendRequest = page.waitForRequest((request) =>
    request.url().endsWith("/api/channel/telegram/operator/sessions/chat-1/messages") && request.method() === "POST");
  await textarea.fill("browser reversed Ctrl+Enter send");
  await textarea.press("Control+Enter");
  const ctrlEnterRequest = await ctrlEnterSendRequest;
  expect(JSON.parse(ctrlEnterRequest.postData() ?? "{}")).toEqual(expect.objectContaining({ text: "browser reversed Ctrl+Enter send" }));
  await expect(textarea).toHaveValue("");

  const messageList = page.locator("[data-telegram-message-list]");
  const positionBeforeRefresh = await messageList.evaluate((element) => {
    element.scrollTop = Math.floor(element.scrollHeight / 2);
    element.dispatchEvent(new Event("scroll", { bubbles: true }));
    return element.scrollTop;
  });
  await expect.poll(() => messageList.evaluate((element) => element.scrollTop)).toBe(positionBeforeRefresh);
  const timelineRefreshResponse = page.waitForResponse((response) =>
    response.request().method() === "GET"
    && response.url().includes("/api/channel/telegram/operator/sessions/chat-1/messages"));
  await page.getByRole("button", { name: "Refresh" }).click();
  await timelineRefreshResponse;
  await expect.poll(() => messageList.evaluate((element) => element.scrollTop)).toBe(positionBeforeRefresh);

  const beforeScroll = await page.evaluate(() => ({
    windowScrollY: window.scrollY,
    documentScrollTop: (document.scrollingElement ?? document.documentElement).scrollTop,
  }));
  await messageList.evaluate((element) => {
    element.scrollTop = element.scrollHeight;
  });
  await expect.poll(() => messageList.evaluate((element) => element.scrollTop)).toBeGreaterThan(0);
  const afterMessageScroll = await page.evaluate(() => ({
    windowScrollY: window.scrollY,
    documentScrollTop: (document.scrollingElement ?? document.documentElement).scrollTop,
  }));
  expect(afterMessageScroll).toEqual(beforeScroll);

  const sessions = page.locator("[data-telegram-sessions]");
  await sessions.evaluate((element) => {
    element.scrollTop = element.scrollHeight;
  });
  await expect.poll(() => sessions.evaluate((element) => element.scrollTop)).toBeGreaterThan(0);
  const afterSessionScroll = await page.evaluate(() => ({
    windowScrollY: window.scrollY,
    documentScrollTop: (document.scrollingElement ?? document.documentElement).scrollTop,
  }));
  expect(afterSessionScroll).toEqual(beforeScroll);
});

test("Telegram operator uses message senders for group and channel conversations", async ({ page }) => {
  test.skip(appProfile !== "custom", "requires APP_PROFILE=custom so channel.telegram is enabled");

  await signInAsAdmin(page);
  await installTelegramFixtures(page, {
    sessions: groupAndChannelOperatorSessions(),
    messagesByChatId: groupAndChannelOperatorMessages(),
  });
  await page.getByRole("link", { name: "Telegram channel" }).click();
  await expect(page.locator("#page-title")).toHaveText("Telegram channel");
  await page.getByRole("button", { name: "Open operator conversations" }).click();

  await expect(page.locator('[data-telegram-session="chat-group"]')).toBeVisible();
  await expect(page.locator('[data-telegram-message-list] [data-telegram-message]')).toHaveCount(3);
  const groupSenders = page.locator('[data-telegram-message-list] [data-telegram-sender]');
  await expect(groupSenders.nth(0)).toHaveText("alice");
  await expect(groupSenders.nth(1)).toHaveText("User");
  await expect(groupSenders.nth(2)).toHaveText("Bot");
  expect(await groupSenders.allTextContents()).not.toContain("Support group");

  await page.locator('[data-telegram-session="chat-channel"]').click();
  await expect(page.locator('[data-telegram-message]').filter({ hasText: "Message from the channel editor" })).toBeVisible();
  const channelSenders = page.locator('[data-telegram-message-list] [data-telegram-sender]');
  await expect(channelSenders).toHaveCount(3);
  await expect(channelSenders.nth(0)).toHaveText("channel_editor");
  await expect(channelSenders.nth(1)).toHaveText("User");
  await expect(channelSenders.nth(2)).toHaveText("Bot");
  expect(await channelSenders.allTextContents()).not.toContain("Announcements");
});

test("ordinary long pages retain page-level vertical scrolling", async ({ page }) => {
  test.skip(appProfile !== "custom", "requires APP_PROFILE=custom for the shared browser harness");

  const users = Array.from({ length: 100 }, (_, index) => ({
    id: `user-${index + 1}`,
    username: `user_${index + 1}`,
    name: `User ${index + 1}`,
    roles: ["viewer"],
    enabled: true,
    mustChangePassword: false,
    createdAt: "2026-09-05T00:00:00Z",
  }));
  await page.route("**/api/users*", async (route) => {
    await route.fulfill(jsonResponse({ items: users, total: users.length, page: 1, pageSize: users.length }));
  });

  await signInAsAdmin(page);
  await page.goto("/users");
  await expect(page.getByRole("heading", { name: "Users" })).toBeVisible();
  await expect(page.locator('[data-table-presentation="desktop-table"]')).toBeVisible();

  const metrics = await page.evaluate(() => {
    const main = document.querySelector<HTMLElement>('[data-shell-region="main"]');
    if (main === null) throw new Error("missing main");
    return {
      mode: main.dataset.shellScrollMode,
      overflowY: getComputedStyle(main).overflowY,
      canScroll: main.scrollHeight > main.clientHeight,
    };
  });
  expect(metrics.mode).toBe("page");
  expect(metrics.overflowY).toBe("auto");
  expect(metrics.canScroll).toBe(true);
});
