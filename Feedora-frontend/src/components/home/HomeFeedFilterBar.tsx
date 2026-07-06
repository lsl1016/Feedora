import { Button, Card, Select, Space } from 'antd';
import { useEffect, useState } from 'react';
import { api } from '../../api';
import type { Circle, ContentTag, Topic, TimeRange } from '../../types';

export type HomeFeedType = 'recommend' | 'latest' | 'following' | 'friend' | 'hot';

export interface HomeFeedFilterValue {
  feedType: HomeFeedType;
  timeRange: TimeRange;
  tagId?: number;
  circleId?: number;
  topicId?: number;
}

interface Props {
  value: HomeFeedFilterValue;
  onChange: (value: HomeFeedFilterValue) => void;
}

const defaultValue: HomeFeedFilterValue = { feedType: 'recommend', timeRange: 'all' };

export function HomeFeedFilterBar({ value, onChange }: Props) {
  const [tags, setTags] = useState<ContentTag[]>([]);
  const [circles, setCircles] = useState<Circle[]>([]);
  const [topics, setTopics] = useState<Topic[]>([]);

  useEffect(() => {
    api.getTags().then(setTags);
    api.getCircles({ page: 1, pageSize: 50, scope: 'all' }).then((res) => setCircles(res.list));
    api.getTopics({ page: 1, pageSize: 50, type: 'all' }).then((res) => setTopics(res.list));
  }, []);

  const update = (patch: Partial<HomeFeedFilterValue>) => onChange({ ...value, ...patch });

  return (
    <Card className="soft-card filter-card home-feed-filter-card">
      <Space wrap>
        <Select
          value={value.feedType}
          style={{ width: 120 }}
          onChange={(feedType) => update({ feedType })}
          options={[
            { value: 'recommend', label: '推荐' },
            { value: 'latest', label: '最新' },
            { value: 'following', label: '关注' },
            { value: 'friend', label: '朋友圈' },
            { value: 'hot', label: '热门' },
          ]}
        />
        <Select
          value={value.timeRange}
          style={{ width: 130 }}
          onChange={(timeRange) => update({ timeRange })}
          options={[
            { value: 'all', label: '全部时间' },
            { value: 'today', label: '今天' },
            { value: 'week', label: '本周' },
            { value: 'month', label: '本月' },
          ]}
        />
        <Select
          allowClear
          placeholder="全部标签"
          value={value.tagId}
          style={{ width: 150 }}
          onChange={(tagId) => update({ tagId })}
          options={tags.map((tag) => ({ value: tag.tagId, label: tag.tagName }))}
        />
        <Select
          allowClear
          placeholder="全部圈子"
          value={value.circleId}
          style={{ width: 170 }}
          onChange={(circleId) => update({ circleId })}
          options={circles.map((circle) => ({ value: circle.circleId, label: circle.name }))}
        />
        <Select
          allowClear
          placeholder="全部话题"
          value={value.topicId}
          style={{ width: 170 }}
          onChange={(topicId) => update({ topicId })}
          options={topics.map((topic) => ({ value: topic.topicId, label: `#${topic.name}#` }))}
        />
        <Button onClick={() => onChange(defaultValue)}>重置</Button>
      </Space>
    </Card>
  );
}

export const defaultHomeFeedFilter = defaultValue;
