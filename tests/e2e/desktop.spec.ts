import { expect, test } from '@playwright/test';

test('desktop mock workflow works across core tabs', async ({ page }) => {
  await page.goto('http://127.0.0.1:5174/#/community');
  await expect(page.getByRole('heading', { name: '社区首页' })).toBeVisible();
  await expect(page.getByText('Feedora Desktop V1：从社区到本地 Knowledge IDE 的架构设计')).toBeVisible();

  const firstFeed = page.locator('.feed-card').first();
  await firstFeed.getByRole('button', { name: '保存为笔记' }).click();
  await expect(page.getByText('已保存到本地笔记')).toBeVisible();

  await firstFeed.getByRole('button', { name: '添加到 Context' }).click();
  await expect(page.locator('.ai-context-card')).toContainText('帖子: Feedora Desktop V1');

  await page.locator('.ai-panel-input textarea').fill('请总结当前帖子');
  await page.locator('.ai-panel-input').getByRole('button', { name: '发送' }).click();
  await expect(page.locator('.ai-chat-scroll')).toContainText('请总结当前帖子');
  await expect(page.locator('.ai-chat-scroll')).toContainText('当前已使用');

  await page.getByRole('link', { name: '工作台' }).click();
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
  await expect(page.getByText('最近编辑')).toBeVisible();

  await page.goto('http://127.0.0.1:5174/#/workspace/notes');
  await expect(page.getByText('笔记编辑器').or(page.getByText('Feedora Desktop V1 架构拆解')).first()).toBeVisible();
  await expect(page.getByText('Feedora Desktop V1：从社区到本地 Knowledge IDE 的架构设计')).toBeVisible();

  await page.goto('http://127.0.0.1:5174/#/repository');
  await expect(page.getByText('本地仓库阅读、索引与 AI 辅助理解')).toBeVisible();
  await expect(page.locator('.code-preview')).toContainText('package main');

  await page.goto('http://127.0.0.1:5174/#/search');
  await expect(page.getByRole('heading', { name: '统一搜索' })).toBeVisible();
  await expect(page.getByText('Wails Bridge 设计要点')).toBeVisible();

  await page.goto('http://127.0.0.1:5174/#/notifications');
  await expect(page.getByRole('heading', { name: '通知中心' })).toBeVisible();
  await page.getByRole('button', { name: '全部标记已读' }).click();
  await expect(page.locator('.notification-card.unread')).toHaveCount(0);

  await page.goto('http://127.0.0.1:5174/#/settings');
  await expect(page.getByRole('heading', { name: '设置' })).toBeVisible();
  await page.getByRole('button', { name: '保存全部设置' }).click();
  await expect(page.getByText('设置已保存到 Mock Desktop State')).toBeVisible();
});
