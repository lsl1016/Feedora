import { expect, test } from '@playwright/test';

test('web runs completely in mock mode without backend', async ({ page }) => {
  await page.goto('http://127.0.0.1:5173/');
  await expect(page.getByText('分享你的技术经验、项目复盘或学习笔记...')).toBeVisible();
  await expect(page.getByText(/Feedora Desktop V1：从社区到本地 Knowledge IDE/).first()).toBeVisible();
  await expect(page.getByText('推荐圈子')).toBeVisible();

  await page.goto('http://127.0.0.1:5173/workspace');
  await expect(page.locator('.workspace-content').getByRole('heading', { name: '工作空间' })).toBeVisible();
  await expect(page.getByText('Feedora Desktop V1 架构拆解').first()).toBeVisible();

  await page.goto('http://127.0.0.1:5173/topics');
  await expect(page.getByRole('heading', { name: '话题广场' })).toBeVisible();
  await expect(page.getByText(/AI 工程化/).first()).toBeVisible();

  await page.goto('http://127.0.0.1:5173/notifications');
  await expect(page.getByText('Shirley 评论了你的帖子')).toBeVisible();

  await page.goto('http://127.0.0.1:5173/search?keyword=Feedora');
  await expect(page.getByText(/Feedora Desktop V1：从社区到本地 Knowledge IDE/).first()).toBeVisible();
});


test('web mock write interactions work without backend', async ({ page }) => {
  await page.goto('http://127.0.0.1:5173/login');
  await page.getByRole('button', { name: '登录' }).click();
  await expect(page).toHaveURL('http://127.0.0.1:5173/');

  const firstPost = page.locator('.post-card').first();
  await expect(firstPost).toContainText('Feedora Desktop V1');

  const likeButton = firstPost.getByRole('button', { name: /点赞 684/ });
  await likeButton.click();
  await expect(firstPost.getByRole('button', { name: /点赞 685/ })).toBeVisible();

  const favoriteButton = firstPost.getByRole('button', { name: /收藏 438/ });
  await favoriteButton.click();
  await expect(firstPost.getByRole('button', { name: /收藏 437/ })).toBeVisible();

  await firstPost.locator('.post-click-zone').click();
  await expect(page).toHaveURL(/\/posts\/101/);
  await page.getByPlaceholder('说点什么...').fill('Mock E2E 评论：桌面端和 Web 的共享 Client Core 很清晰。');
  await page.getByRole('button', { name: '发表评论' }).click();
  await expect(page.getByText('Mock E2E 评论：桌面端和 Web 的共享 Client Core 很清晰。')).toBeVisible();

  await page.goto('http://127.0.0.1:5173/workspace/notes/create');
  await page.getByLabel('标题').fill('E2E 新建知识笔记');
  await page.getByLabel('正文').fill('这是通过 Web Mock Runtime 保存的一篇测试笔记，用于验证无后端情况下的写入交互。');
  await page.getByRole('button', { name: '保存' }).click();
  await expect(page).toHaveURL(/\/workspace\/notes\/\d+/);
  await expect(page.getByDisplayValue('E2E 新建知识笔记')).toBeVisible();
});
