# PTCG Trader
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