// Package search 封装 Elasticsearch 索引写入与查询，支撑帖子 / 用户 / 圈子 / 话题全文搜索。
package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

// Client ES 封装，按 indexPrefix 管理各索引。
type Client struct {
	es          *elasticsearch.Client
	indexPrefix string
}

func NewClient(es *elasticsearch.Client, indexPrefix string) *Client {
	if indexPrefix == "" {
		indexPrefix = "feedora"
	}
	return &Client{es: es, indexPrefix: indexPrefix}
}

// 各索引名。
func (c *Client) PostIndex() string   { return c.indexPrefix + "_posts" }
func (c *Client) UserIndex() string   { return c.indexPrefix + "_users" }
func (c *Client) CircleIndex() string { return c.indexPrefix + "_circles" }
func (c *Client) TopicIndex() string  { return c.indexPrefix + "_topics" }

// EnsureIndices 确保各索引存在（不存在则创建，使用默认动态映射）。
func (c *Client) EnsureIndices(ctx context.Context) error {
	for _, idx := range []string{c.PostIndex(), c.UserIndex(), c.CircleIndex(), c.TopicIndex()} {
		exists, err := c.indexExists(ctx, idx)
		if err != nil {
			return err
		}
		if !exists {
			if err := c.createIndex(ctx, idx); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Client) indexExists(ctx context.Context, index string) (bool, error) {
	res, err := c.es.Indices.Exists([]string{index}, c.es.Indices.Exists.WithContext(ctx))
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	return res.StatusCode == 200, nil
}

func (c *Client) createIndex(ctx context.Context, index string) error {
	res, err := c.es.Indices.Create(index, c.es.Indices.Create.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("创建索引 %s 失败: %s", index, res.String())
	}
	return nil
}

// IndexDoc 写入 / 更新一条文档。
func (c *Client) IndexDoc(ctx context.Context, index, id string, doc any) error {
	b, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	res, err := c.es.Index(index, bytes.NewReader(b),
		c.es.Index.WithDocumentID(id),
		c.es.Index.WithContext(ctx),
		c.es.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("写入文档失败: %s", res.String())
	}
	return nil
}

// DeleteDoc 删除一条文档（忽略不存在）。
func (c *Client) DeleteDoc(ctx context.Context, index, id string) error {
	res, err := c.es.Delete(index, id, c.es.Delete.WithContext(ctx), c.es.Delete.WithRefresh("true"))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

// SearchIDs 在指定索引上按 query DSL 搜索，返回命中文档的 _id 列表与总数。
func (c *Client) SearchIDs(ctx context.Context, index string, query map[string]any, from, size int) ([]int64, int64, error) {
	query["from"] = from
	query["size"] = size
	query["_source"] = false
	b, _ := json.Marshal(query)
	res, err := c.es.Search(
		c.es.Search.WithContext(ctx),
		c.es.Search.WithIndex(index),
		c.es.Search.WithBody(bytes.NewReader(b)),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, 0, fmt.Errorf("搜索失败: %s", res.String())
	}
	var r struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				ID string `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	ids := make([]int64, 0, len(r.Hits.Hits))
	for _, h := range r.Hits.Hits {
		if id, err := strconv.ParseInt(h.ID, 10, 64); err == nil {
			ids = append(ids, id)
		}
	}
	return ids, r.Hits.Total.Value, nil
}

// MatchQuery 构建 multi_match + filter 的查询 DSL。
func MatchQuery(keyword string, fields []string, filters map[string]any) map[string]any {
	must := []any{}
	if strings.TrimSpace(keyword) != "" {
		must = append(must, map[string]any{
			"multi_match": map[string]any{
				"query":  keyword,
				"fields": fields,
			},
		})
	} else {
		must = append(must, map[string]any{"match_all": map[string]any{}})
	}
	filterClauses := []any{}
	for k, v := range filters {
		filterClauses = append(filterClauses, map[string]any{"term": map[string]any{k: v}})
	}
	return map[string]any{
		"query": map[string]any{
			"bool": map[string]any{
				"must":   must,
				"filter": filterClauses,
			},
		},
	}
}
