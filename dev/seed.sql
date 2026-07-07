-- tabgraph デモ用 ECサイトスキーマ
-- ER図・全文検索・FK推定のテストに最適化

-- ユーザー
CREATE TABLE users (
    id          SERIAL PRIMARY KEY,
    email       VARCHAR(255) NOT NULL UNIQUE,
    username    VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(255),
    password_hash VARCHAR(255) NOT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE users IS 'サービス登録ユーザー';
COMMENT ON COLUMN users.email IS 'ログインに使うメールアドレス。小文字で正規化して保存';
COMMENT ON COLUMN users.username IS 'ユニークなユーザー名（表示用）';
COMMENT ON COLUMN users.password_hash IS 'bcrypt でハッシュ化したパスワード';
COMMENT ON COLUMN users.is_active IS 'false = 退会済み or BAN';

-- カテゴリ（自己参照で階層構造）
CREATE TABLE categories (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    slug        VARCHAR(100) NOT NULL UNIQUE,
    parent_id   INT REFERENCES categories(id),
    sort_order  INT NOT NULL DEFAULT 0
);
COMMENT ON TABLE categories IS '商品カテゴリ（階層構造）';
COMMENT ON COLUMN categories.slug IS 'URLフレンドリーな識別子';
COMMENT ON COLUMN categories.parent_id IS 'NULL = ルートカテゴリ';

-- 商品
CREATE TABLE products (
    id          SERIAL PRIMARY KEY,
    sku         VARCHAR(100) NOT NULL UNIQUE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    category_id INT NOT NULL REFERENCES categories(id),
    price       NUMERIC(10,2) NOT NULL,
    stock       INT NOT NULL DEFAULT 0,
    is_published BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE products IS '販売商品マスタ';
COMMENT ON COLUMN products.sku IS '在庫管理用の商品コード（Stock Keeping Unit）';
COMMENT ON COLUMN products.price IS '税抜き価格（円）';
COMMENT ON COLUMN products.stock IS '現在の在庫数。0以下は販売停止';
COMMENT ON COLUMN products.is_published IS 'true = ストアに公開中';

-- 商品画像
CREATE TABLE product_images (
    id          SERIAL PRIMARY KEY,
    product_id  INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    url         VARCHAR(500) NOT NULL,
    alt_text    VARCHAR(255),
    sort_order  INT NOT NULL DEFAULT 0,
    is_primary  BOOLEAN NOT NULL DEFAULT false
);
COMMENT ON TABLE product_images IS '商品に紐づく画像（複数可）';
COMMENT ON COLUMN product_images.is_primary IS 'true = 一覧ページで表示するメイン画像';

-- 注文
CREATE TABLE orders (
    id          SERIAL PRIMARY KEY,
    user_id     INT NOT NULL REFERENCES users(id),
    status      VARCHAR(50) NOT NULL DEFAULT 'pending',
    total_price NUMERIC(10,2) NOT NULL,
    note        TEXT,
    ordered_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    shipped_at  TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ
);
COMMENT ON TABLE orders IS '注文ヘッダー。明細は order_items を参照';
COMMENT ON COLUMN orders.status IS 'pending / paid / shipped / delivered / cancelled';
COMMENT ON COLUMN orders.total_price IS '送料・税込みの最終金額';

-- 注文明細
CREATE TABLE order_items (
    id          SERIAL PRIMARY KEY,
    order_id    INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id  INT NOT NULL REFERENCES products(id),
    quantity    INT NOT NULL,
    unit_price  NUMERIC(10,2) NOT NULL
);
COMMENT ON TABLE order_items IS '注文1件あたりの商品明細';
COMMENT ON COLUMN order_items.unit_price IS '購入時点の単価（後から変更されても影響なし）';

-- レビュー
CREATE TABLE reviews (
    id          SERIAL PRIMARY KEY,
    product_id  INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    user_id     INT NOT NULL REFERENCES users(id),
    rating      SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    title       VARCHAR(255),
    body        TEXT,
    is_approved BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(product_id, user_id)
);
COMMENT ON TABLE reviews IS '商品レビュー（ユーザーが1商品につき1件投稿可）';
COMMENT ON COLUMN reviews.rating IS '1〜5の星評価';
COMMENT ON COLUMN reviews.is_approved IS 'モデレーション済みのレビューのみ公開';

-- 住所（配送先）
CREATE TABLE addresses (
    id          SERIAL PRIMARY KEY,
    user_id     INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    label       VARCHAR(100),
    postal_code VARCHAR(10) NOT NULL,
    prefecture  VARCHAR(50) NOT NULL,
    city        VARCHAR(100) NOT NULL,
    line1       VARCHAR(255) NOT NULL,
    line2       VARCHAR(255),
    is_default  BOOLEAN NOT NULL DEFAULT false
);
COMMENT ON TABLE addresses IS 'ユーザーの配送先住所帳';
COMMENT ON COLUMN addresses.label IS '「自宅」「会社」など任意のラベル';
COMMENT ON COLUMN addresses.is_default IS 'true = 注文時のデフォルト配送先';

-- タグ
CREATE TABLE tags (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE
);
COMMENT ON TABLE tags IS '商品に付けるフリータグ';

-- 商品↔タグ 中間テーブル
CREATE TABLE product_tags (
    product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    tag_id     INT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, tag_id)
);
COMMENT ON TABLE product_tags IS '商品とタグの多対多関連';

-- クーポン
CREATE TABLE coupons (
    id              SERIAL PRIMARY KEY,
    code            VARCHAR(50) NOT NULL UNIQUE,
    discount_type   VARCHAR(20) NOT NULL,
    discount_value  NUMERIC(10,2) NOT NULL,
    min_order_price NUMERIC(10,2),
    usage_limit     INT,
    used_count      INT NOT NULL DEFAULT 0,
    expires_at      TIMESTAMPTZ
);
COMMENT ON TABLE coupons IS '割引クーポンマスタ';
COMMENT ON COLUMN coupons.discount_type IS '"percent" = 割引率（%）, "fixed" = 定額値引き（円）';
COMMENT ON COLUMN coupons.min_order_price IS 'このクーポンを使うための最低注文金額';
COMMENT ON COLUMN coupons.usage_limit IS 'NULL = 使用回数無制限';

-- 注文↔クーポン
CREATE TABLE order_coupons (
    order_id  INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    coupon_id INT NOT NULL REFERENCES coupons(id),
    PRIMARY KEY (order_id, coupon_id)
);
COMMENT ON TABLE order_coupons IS '注文で使用されたクーポンの記録';

-- インデックス
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_orders_user ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_order_items_order ON order_items(order_id);
CREATE INDEX idx_reviews_product ON reviews(product_id);

-- ==================== サンプルデータ ====================

-- カテゴリ
INSERT INTO categories (name, slug, parent_id, sort_order) VALUES
    ('家電', 'electronics', NULL, 1),
    ('ファッション', 'fashion', NULL, 2),
    ('本・メディア', 'books', NULL, 3),
    ('スマートフォン', 'smartphones', 1, 1),
    ('ノートPC', 'laptops', 1, 2),
    ('メンズ', 'mens', 2, 1),
    ('レディース', 'womens', 2, 2);

-- タグ
INSERT INTO tags (name, slug) VALUES
    ('新着', 'new'),
    ('セール', 'sale'),
    ('人気', 'popular'),
    ('限定', 'limited'),
    ('エコ', 'eco');

-- ユーザー
INSERT INTO users (email, username, display_name, password_hash) VALUES
    ('alice@example.com', 'alice', 'Alice Tanaka', '$2b$12$dummy_hash_alice'),
    ('bob@example.com', 'bob', 'Bob Sato', '$2b$12$dummy_hash_bob'),
    ('charlie@example.com', 'charlie', 'Charlie Suzuki', '$2b$12$dummy_hash_charlie'),
    ('diana@example.com', 'diana', 'Diana Yamamoto', '$2b$12$dummy_hash_diana');

-- 商品
INSERT INTO products (sku, name, description, category_id, price, stock, is_published) VALUES
    ('PHONE-001', 'ProMax X15', '最新フラッグシップスマートフォン。6.7インチ有機EL、5000mAhバッテリー。', 4, 129800.00, 50, true),
    ('PHONE-002', 'AquaPro 5G', 'IPX8防水対応の5G対応スマートフォン。コスパ最強モデル。', 4, 59800.00, 120, true),
    ('LAPTOP-001', 'UltraBook Air M3', '薄型軽量ノートPC。M3チップ搭載で長時間バッテリー。', 5, 168000.00, 30, true),
    ('LAPTOP-002', 'DevBook Pro', '開発者向けハイスペックノートPC。RAM 64GB、SSD 2TB。', 5, 298000.00, 10, true),
    ('MENS-001', 'オーガニックコットンTシャツ', '肌に優しいオーガニックコットン100%。', 6, 3800.00, 200, true),
    ('BOOK-001', 'Go言語による並行処理', 'Go言語の並行処理パターンを体系的に解説。', 3, 3520.00, 80, true);

-- 商品タグ
INSERT INTO product_tags (product_id, tag_id) VALUES
    (1, 1), (1, 3),  -- ProMax X15: 新着・人気
    (2, 2), (2, 3),  -- AquaPro: セール・人気
    (3, 1),          -- UltraBook: 新着
    (5, 5);          -- Tシャツ: エコ

-- 商品画像
INSERT INTO product_images (product_id, url, alt_text, sort_order, is_primary) VALUES
    (1, 'https://example.com/images/phone-001-main.jpg', 'ProMax X15 正面', 0, true),
    (1, 'https://example.com/images/phone-001-back.jpg', 'ProMax X15 背面', 1, false),
    (2, 'https://example.com/images/phone-002-main.jpg', 'AquaPro 5G', 0, true),
    (3, 'https://example.com/images/laptop-001-main.jpg', 'UltraBook Air', 0, true);

-- 住所
INSERT INTO addresses (user_id, label, postal_code, prefecture, city, line1, is_default) VALUES
    (1, '自宅', '150-0001', '東京都', '渋谷区', '神南1-2-3', true),
    (2, '自宅', '530-0001', '大阪府', '大阪市北区', '梅田4-5-6', true),
    (3, '会社', '220-0012', '神奈川県', '横浜市西区', 'みなとみらい7-8-9', true);

-- 注文
INSERT INTO orders (user_id, status, total_price, ordered_at, shipped_at) VALUES
    (1, 'delivered', 129800.00, NOW() - INTERVAL '10 days', NOW() - INTERVAL '7 days'),
    (1, 'paid',      59800.00, NOW() - INTERVAL '2 days', NULL),
    (2, 'delivered', 171800.00, NOW() - INTERVAL '15 days', NOW() - INTERVAL '12 days'),
    (3, 'pending',   3520.00, NOW() - INTERVAL '1 hour', NULL);

-- 注文明細
INSERT INTO order_items (order_id, product_id, quantity, unit_price) VALUES
    (1, 1, 1, 129800.00),
    (2, 2, 1, 59800.00),
    (3, 3, 1, 168000.00),
    (3, 6, 1, 3520.00),  -- 同一注文で2商品
    (4, 6, 1, 3520.00);

-- クーポン
INSERT INTO coupons (code, discount_type, discount_value, min_order_price, usage_limit, used_count) VALUES
    ('WELCOME10', 'percent', 10.00, 5000.00, 1000, 42),
    ('SUMMER3000', 'fixed', 3000.00, 20000.00, 500, 128),
    ('VIP20', 'percent', 20.00, 50000.00, 100, 7);

-- レビュー
INSERT INTO reviews (product_id, user_id, rating, title, body, is_approved) VALUES
    (1, 2, 5, '最高のスマホ', 'カメラの画質が素晴らしい。バッテリーも一日余裕で持ちます。', true),
    (1, 3, 4, 'ほぼ満足', '性能は申し分ないが価格が少し高い。', true),
    (3, 1, 5, '仕事がはかどる', '軽くてバッテリーが長持ち。M3の性能は本物。', true),
    (6, 1, 4, 'わかりやすい', 'Go初心者でも理解しやすい構成。サンプルコードが実用的。', true);
