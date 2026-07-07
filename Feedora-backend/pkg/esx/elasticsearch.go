// Package esx 封装 Elasticsearch 客户端初始化。
// 用于阶段二的全文搜索、搜索联想与聚合查询。
package esx

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/feedora/backend/pkg/config"
)

// New 根据配置创建 ES 客户端并校验连通性。
func New(cfg config.ESConfig) (*elasticsearch.Client, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	})
	if err != nil {
		return nil, err
	}
	res, err := client.Info()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return client, nil
}
