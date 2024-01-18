# websocket server 2.1.2

* websocket server
  * 統一 websocket manager 管理
  * client 連線傳送值timeout 1s

## 啟動方式

#### 1. 使用docker-compose

沒有mysql container可以選擇此方式

1. 確定docker 有mysql image

   `docker pull mysql`
2. 修改docker-compose.yaml
   **Windows:** line:51 註解linux命令 取消line:52的註解
   **Linux:** 取消line:52註解 line:51註解windows命令
   line:51:

   ```
   dockerfile: deploy/api/linux/Dockerfile
   ```

   line52:

   ```
   dockerfile: deploy/api/windows/Dockerfile
   ```
3. docker-compose up

   `docker-compose up`

#### 2. 使用DockerFile啟動

有固定的mysql container可選擇此方式

1. ##### Windows系統

   1. 創建需要的image
      api image:
      `docker build -t websocket_server:2.1.2 -f deploy/api/windows/Dockerfile .`
   2. 創建並啟動container

      1. run api container

         `docker run --name websocket_server -p 5488:5488 -e DB_HOST=host.docker.internal -v ${PWD}/docker/log:/app/log websocket_server:2.1.2`
2. ##### Linux系統

   1. 創建需要的image
      api image:
      `docker build -t websocket_server:2.1.2 -f deploy/api/linux/Dockerfile .`

   2. 創建並啟動container

      1. run api container
         `docker run --name websocket_server -p 5488:5488 --network="host" -v ${PWD}/docker/log:/app/log websocket_server:2.1.2 -e TZ=Asia/Taipei`
         `docker run --name websocket_server -p 5488:5488 -e TZ=Asia/Taipei -e INFLUXDB_HOST=192.168.1.11 -e REDIS_HOST=192.168.1.11 -v ${PWD}/docker/log:/app/log websocket_server:2.1.2`

# Log File

用docker啟動的程式log file在  docker/log/
