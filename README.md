# Тестовое задание
1. Решение внутри рпеозитория
2. При увеличение количества записей свыше одного миллиона записей необходимо будет сделать рейт лимиты для каждого пользователя, ограничить максимальный размер ответа(сделано), использовать пагинацию(сделано). Возможно использовать кэши на выдачу списка, а также шардирование по айди пользователя в перспективе.
3. 
  ```sql
  CREATE TABLE users (
      id SERIAL PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      email VARCHAR(255) UNIQUE NOT NULL,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
  );

  CREATE TABLE products (
      id SERIAL PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      description TEXT,
      price NUMERIC(10, 2) CHECK (price > 0) NOT NULL,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );

  CREATE TABLE orders (
      id SERIAL PRIMARY KEY,
      user_id INT NOT NULL,
      total_amount NUMERIC(10, 2) CHECK (total_amount >= 0) NOT NULL,
      order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      status VARCHAR(50) DEFAULT 'pending',
      FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
  );

  CREATE TABLE order_items (
      id SERIAL PRIMARY KEY,
      order_id INT NOT NULL,
      product_id INT NOT NULL,
      quantity INT CHECK (quantity > 0) NOT NULL,
      price NUMERIC(10, 2) CHECK (price >= 0) NOT NULL,
      FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
      FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
  );

  -- Скорее всего придётся ускорить поиск по email
  CREATE INDEX idx_users_email ON users (email);

  -- Понадобится доставать заказы по айди пользователя и дате
  CREATE INDEX idx_orders_user_id_order_date ON orders (user_id, order_date);
  CREATE INDEX idx_orders_user_id ON orders (user_id);
  CREATE INDEX idx_orders_order_date ON orders (order_date);


  CREATE INDEX idx_order_items_order_id ON order_items (order_id);
  CREATE INDEX idx_order_items_product_id ON order_items (product_id);
  ```
  Для обеспечение целостности используются индексы юникальности, ограничения и 3НФ(насколько я понимаю уточнить форму пока что нельзя, так недостаточно связей)
  
4. 
```sql
SELECT u.id AS user_id, u.name, u.email, o.id AS order_id, o.total_amount, o.order_date
FROM users u
JOIN orders o
  ON u.id = o.user_id
WHERE o.order_date = (
    SELECT MAX(order_date)
    FROM orders
    WHERE user_id = u.id
);
```

5. Решение внутри репозитория. 

6. Решение внутри репозитория. Часть про развёртывание на сервере не сделана так как развернуть внутри PaaS это одно, внутри DigitalOcean, Kubernetes и прочего совсем разные вещи.

7. Компонент лежит внутри web/src/UsersList.tsx

8. Если количество пользователей которое надо будет отрендерит будет очень большим то в соответствиии с документацией React использую react-window или react-virtualized, то есть виртуализацию списков

9. Для интеграции с системой 1С внутри 1с должен быть поднят компонент веб-сервера, который уже может отдавать данные в любом формате, задача интеграции Go и 1C заключается в вызове REST API через Golang http клиента, тут уже необходимо будует соблюдать идемпотентность, ретраи, но всё точно также как и при интеграции например с юкассой, по крайней мере в тривиальных задачах.

10. На стороне 1С необходимо включить дополнительные отчеты и обработки, после чего создать обработчики, которые будет сервить специальный объект веб-сервиса, после чего уже можно обеспечивать дуплекную связь с апишкой написанной на голанге. Часто, интеграция 1c требует передачу файлов, CSV или XML, но Http предоставляет хороший функционал для передачи даже достаточно больших файлов.

11. Сам по себе nginx не может спасти от высокой нагрузки, так как всегда есть одна точка входа, которая должна обеспечиваться сторонним сервисом для обеспечения высокой доступности, который уже будет использоваться в качестве внутреннего балансировщика нагрузки на наши сервисы, при настройке необходимо будет установить таймауты для реверс прокси, а также создать резервацию для ендпоинтов, поставить для них вес, количетсво неудачных попыток для признания нерабочей и таймату на каждый и определиться с методом балансировки. Nginx также имеет встроенное HA решение для поддержки нескольких нод, тем самым получится добиться отказоустойчивости.

12. Взаимодействие между севрисами может быть организованно через синхронные каналы(unary gRPC), в данном слуаче придётся обеспокоиться обнаружением сервисом, в случае Docker Swarm можно обойтись встроенными средствми днс маршаллинга, а в случае Kubernetes можно использовать встроенные объекты Service, но из-за ограничений kube-proxy такой вариант лишает нас многих преимуществ кубера, поэтому лучше будет рассмотреть использование Service Mesh(Istio). Либо через асинхронные каналы(Kafka), здесь уже нет проблемы обнаружение, но зато чаще возникают проблемы консистености данных, которые придётся решать с помощью настроки кафки и некоторыми паттернами.

13. Использовал бы managed kuberenets, настроил бы его для того чтобы трафик с 1С мог проходить дуплексно в кластер, после чего спроектировал микросервисы исходя из предполагаемых данных, которые будут нужны, возможно использовал бы CockroachDB, если бы понял что будет много межсервисных взаимодействий. Скорее всего это были бы сервисы (Gateway, Интеграция 1C, Сервис кэширования данных из 1C, Сервис обработки даннхы из 1c(возможно разделить на подмодули в соответствии с 1С)), небходимо было бы также обеспечить хорошую телеметрию через Grafana, Prometheus, Tempo, и прочее. Также будет много асинхронных взаимодействий, поэтому придётся обеспечить хорошую шину данных, например Kafka. И как следствие основные паттерны рапреденных систем, например сервсиы оркестровых сага, CQRS сервисы, ретраи, батчи для передачи данных.