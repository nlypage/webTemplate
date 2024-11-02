# Web Template

Simple web application template for Go.

## Libraries used:
- Fiber
- Gorm (PGSQL)
- Viper
- Zap

## OpenAPI Docs
```shell
# 1. Install OpenAPI generator
# Don't forget to add this to your path
go install github.com/swaggo/swag/cmd/swag@latest

# 2. Generate docs
# add this to your build configuration to regenerate it automatically
swag init
```