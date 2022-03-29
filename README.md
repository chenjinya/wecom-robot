# wecom-robot
企业微信机器人小工具

## Usage

```go

c := wecom_robot.NewWeComRobot("<key>")

c.Notice("haha")


```

## Functions

### MustAppendFile
追加文件内容，配合 CSV 使用

```go

wecom_robot.MustAppendFile("test.csv", wecom_robot.ToCsvRow("id", "name", "url"))

```

### ToCsvRow
将参数转换为 CSV 逗号分割格式，里面处理了值里还有逗号的情况



## Features

- 支持文件上传并发送

```go

c.SendFile(ctx, "file_path.csv", "show_name.csv")

```