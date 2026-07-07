package ossx

// 阶段二占位：MinIO 存储实现。
// 阶段一使用 LocalStorage 模拟 OSS；接入 MinIO 时在此实现 Storage 接口，
// 通过配置 oss.type=minio 切换，无需改动业务代码。
