# PTCG Trader
[![Build Status](https://travis-ci.com/XiaoXiaoSN/ptcg_trader.svg?branch=master)](https://travis-ci.com/XiaoXiaoSN/ptcg_trader)

Please design and implement a backend system for an online trading platform of Pokémon Trading Card Game.
- This online trading platform trades 4 kinds of cards only: Pikachu, Bulbasaur, Charmander, and Squirtle.
- The price of cards is between 1.00 USD and 10.00 USD.
- Users on this platform are called traders.
- **There are 10K traders.**
- Traders own unlimited USD and cards.
- Traders can send orders to the platform when they want to buy or sell cards at certain prices.
- A trader can only buy or sell 1 card in 1 order.
- Traders can only buy cards using USD or sell cards for USD.
- Orders are first come first serve.
- There are 2 situations to make a trade:
    - When a buy order is sent to the platform, there exists an uncompleted sell order, whose price is the lowest one among all uncompleted sell orders and less than or equal to the price of the buy order. Then, a trade is made at the price of the sell order. Both buy and sell orders are completed. Otherwise, the buy order is uncompleted.
    - When a sell order is sent to the platform, there exists an uncompleted buy order, whose price is the highest one among all uncompleted buy orders and greater than or equal to the price of the sell order. Then, a trade is made at the price of the buy order. Both buy and sell orders are completed. Otherwise, the sell order is uncompleted.
- Traders can view the status of their latest 50 orders.
- Traders can view the latest 50 trades on each kind of cards.
- **If the sequence of orders is fixed, the results must be the same no matter how many times you execute the sequence.**

## Basic Requirements:
- RESTful API
- Relational database (PostgreSQL, MySQL, ...)
- Containerize
- Testing
- Gracefully shutdown
## Advanced Requirements:
- Multithreading
- Maximize performance of finishing 1M orders
- OpenAPI (Swagger)
- Set up configs using environment variables
- View logs on visualization dashboard (Kibana, Grafana, ...)
- Microservice
- Message queue (Apache Kafka, Apache Pulsar, ...)
- gRPC
- GraphQL
- Docker Compose
- Kubernetes
- Cloud computing platforms (AWS, Azure, GCP, ...) 
- NoSQL
- CI/CD
- User authentication and authorization
- High availability
- ...
Please commit code to your GitHub account.
You must complete basic requirements within 2 weeks. You could ask for 2 more weeks to complete some advance requirements further after finishing basic requirements if you want


## 專案實作使用 TechStack

專案提供 Restful API 服務，使用 Go 搭配 echo 框架實現
- [echo](https://github.com/labstack/echo) 提供高併發 HTTP 服務
- [mikunalpha/goas](https://github.com/mikunalpha/goas) 提供 Code as Documents，產生 Swagger API 文件
- [uber-go/fx](https://github.com/uber-go/fx) 依賴注入工具、生命週期管理，有效管理注入的組件在關閉時能夠 Gracefully shutdown
- [PostgreSQL](https://www.postgresql.org/) 關聯式資料庫提供 ACID Transaction 操作
- [pressly/goose](https://github.com/pressly/goose) 資料庫版本 Migration
- [NATS Streaming](https://github.com/nats-io/nats-streaming-server) 提供持久化、高吞吐量的 Message Queue 服務，其實本來是使用 kafka ，但是後來想說他的 Throughput 更高就跳槽了。沒想到之後在打算切 Partition 時遇到了[阻礙](https://github.com/nats-io/nats-streaming-server/issues/524)
- [Redis](https://redis.io/) 最初設計分散式鎖的選擇，後來發現可以改用 Database row lock 減少程式頻繁詢問鎖的狀態
- Docker 容器化，提高服務環境的一致性

啟動本地服務
- [Make](https://www.gnu.org/software/make/) 讀取 Makefile 提供新進入專案成員快速掌握服務啟動方法
- [Docker Compose](https://docs.docker.com/compose/) 編排容器、網路關聯，快速在本地啟動 NGINX load balance 服務


專案基礎設施部署
- [AWS](https://aws.amazon.com/tw/) 雲端運算服務，部署 EC2 instance Kubernates
- [Terraform](https://www.terraform.io/) IaC 工具，協助以程式碼管理 AWS 資源
- [kOps](https://github.com/kubernetes/kops) 搭配 Terraform 管理 kubernetes 叢集，提供高可用服務的基礎設施
- [Helm](https://helm.sh/) 以套件管理模式來管理 kubernates，整合複雜的 k8s yaml 檔案們
- [PrometheusOperater](https://github.com/prometheus-operator/prometheus-operator) 操作 Prometheus 蒐集服務時序狀態，搭配 AlertManager Slack 告警，並提供 Grafana 監控介面
- [cert-manager](https://cert-manager.io/) 搭配 Route53 完成 ACME 的挑戰，只要部署 Issuer 就會自動簽署、續約 TLS 憑證，提供安全的 HTTPS 服務

日誌蒐集
- [Fluent Bit](https://fluentbit.io/) 輕量型的日誌蒐集系統，部署於 k8s 內持續抓取 container 的輸出
- [Loki](https://grafana.com/oss/loki/) 輕量型 log 查詢系統，和 Grafana Dashboard 整合方便 Trace 服務狀態


持續整合與部署
- [Travis CI](https://travis-ci.org/) 跟 Github 整合提供持續整合、測試程式碼，然後他是免費ㄉ
- [Github Actions](https://github.com/features/actions) 持續部署。自動建置 docker 容器、部署 Helm release，然後他也是免費ㄉ


服務架構包含：
![](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://github.com/XiaoXiaoSN/ptcg_trader/blob/master/documents/architecture.puml)
