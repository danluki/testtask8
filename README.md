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