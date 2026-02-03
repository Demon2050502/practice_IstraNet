CREATE TABLE roles (
  id          BIGSERIAL PRIMARY KEY,
  code        TEXT NOT NULL UNIQUE,     -- 'user', 'operator', 'admin'
  name        TEXT NOT NULL
);

CREATE TABLE users (
  id            BIGSERIAL PRIMARY KEY,
  email         TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  full_name     TEXT NOT NULL,
  is_active     BOOLEAN NOT NULL DEFAULT TRUE,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_roles (
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
  PRIMARY KEY (user_id, role_id)
);

INSERT INTO roles(code, name) VALUES
('user','Пользователь'),
('operator','Оператор'),
('admin','Администратор');


CREATE TABLE application_statuses (
  id        SMALLSERIAL PRIMARY KEY,
  code      TEXT NOT NULL UNIQUE,     -- 'new', 'in_progress', 'waiting', 'resolved', 'closed'
  name      TEXT NOT NULL,
  is_final  BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE application_priorities (
  id    SMALLSERIAL PRIMARY KEY,
  code  TEXT NOT NULL UNIQUE,         -- 'low', 'normal', 'high', 'critical'
  name  TEXT NOT NULL,
  weight SMALLINT NOT NULL            -- для сортировки/очереди
);

CREATE TABLE application_categories (
  id        BIGSERIAL PRIMARY KEY,
  name      TEXT NOT NULL,
  parent_id BIGINT NULL REFERENCES application_categories(id) ON DELETE SET NULL
);

CREATE TABLE applications (
  id              BIGSERIAL PRIMARY KEY,
  title           TEXT NOT NULL,
  description     TEXT NOT NULL,

  status_id       SMALLINT NOT NULL REFERENCES application_statuses(id),
  priority_id     SMALLINT NOT NULL REFERENCES application_priorities(id),
  category_id     BIGINT NULL REFERENCES application_categories(id) ON DELETE SET NULL,

  created_by      BIGINT NOT NULL REFERENCES users(id),                 -- кто создал
  assigned_to     BIGINT NULL REFERENCES users(id) ON DELETE SET NULL,  -- исполнитель (оператор)

  contact_phone   TEXT NULL,
  contact_address TEXT NULL,

  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  closed_at       TIMESTAMPTZ NULL
);

CREATE INDEX idx_applications_created_by   ON applications(created_by);
CREATE INDEX idx_applications_assigned_to  ON applications(assigned_to);
CREATE INDEX idx_applications_status       ON applications(status_id);
CREATE INDEX idx_applications_priority     ON applications(priority_id);
CREATE INDEX idx_applications_created_at   ON applications(created_at);

INSERT INTO application_statuses(code, name, is_final) VALUES
('new','Новая',FALSE),
('in_progress','В работе',FALSE),
('waiting','Ожидание',FALSE),
('resolved','Решена',TRUE),
('closed','Закрыта',TRUE);

INSERT INTO application_priorities(code, name, weight) VALUES
('low','Низкий',10),
('normal','Обычный',20),
('high','Высокий',30),
('critical','Критический',40);


CREATE TABLE application_history (
  id          BIGSERIAL PRIMARY KEY,
  application_id   BIGINT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
  actor_id    BIGINT NOT NULL REFERENCES users(id),

  action      TEXT NOT NULL,  -- 'create', 'status_change', 'assign', 'comment', 'edit'
  field       TEXT NULL,      -- например 'status_id', 'assigned_to', 'priority_id'
  old_value   TEXT NULL,
  new_value   TEXT NULL,

  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE application_comments (
  id         BIGSERIAL PRIMARY KEY,
  application_id  BIGINT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
  author_id  BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  body       TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
