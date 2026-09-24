import { expect, test } from '@playwright/test';

test('web runs completely in mock mode without backend', async ({ page }) => {
  await page.goto('http://127.0.0.1:5173/');
  await expect(page.getByText('分享你的技术经验、项目复盘或学习笔记...')).toBeVisible();
  await expect(page.getByText('Feedora Desktop V1：从社区到本地 Knowledge IDE 的架构设计')).toBeVisible();
  await expect(page.getByText('推荐圈子')).toBeVisible();

  await page.goto('http://127.0.0.1:5173/workspace');
  await expect(page.getByRole('heading', { name: '工作空间' })).toBeVisible();
  await expect(page.getByText('Feedora Desktop V1 架构拆解')).toBeVisible();

  await page.goto('http://127.0.0.1:5173/topics');
  await expect(page.getByRole('heading', { name: '话题广场' })).toBeVisible();
  await expect(page.getByText('#AI 工程化#')).toBeVisible();

  await page.goto('http://127.0.0.1:5173/notifications');
  await expect(page.getByText('Shirley 评论了你的帖子')).toBeVisible();

  await page.goto('http://127.0.0.1:5173/search?keyword=Feedora');
  await expect(page.getByText('Feedora Desktop V1：从社区到本地 Knowledge IDE 的架构设计')).toBeVisible();
});
