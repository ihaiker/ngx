# 语法介绍

```shell script
cat nginx.conf | ngc '.http.server | 
  select(.server_name equal "api.aginx.io" and equal "v2.aginx.io") |
  args() | args_length() | length()'
```
