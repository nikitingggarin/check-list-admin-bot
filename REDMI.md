# README - Telegram бот для управления чек-листами

## 🚀 Быстрый старт

### 1. Клонирование репозитория
git clone https://github.com/nikitingggarin/check-list-admin-bot.git


### 2. Настройка переменных окружения
Создайте файл `.env` в корне проекта:
TELEGRAM_BOT_TOKEN=ваш_токен_бота
SUPABASE_URL=ваш_url_supabase
SUPABASE_KEY=ваш_ключ_supabase

#### 🔧 Получение TELEGRAM_BOT_TOKEN:
1. Откройте Telegram и найдите @BotFather
2. Начните диалог и отправьте команду /newbot
3. Следуйте инструкциям, получите токен вида: 1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
4. Скопируйте этот токен в `.env` файл

#### 🗄️ Получение SUPABASE_URL и SUPABASE_KEY:
1. Зарегистрируйтесь на supabase.com
2. Создайте новый проект
3. После создания проекта перейдите в Settings → API
4. Скопируйте Project URL → SUPABASE_URL и anon public ключ → SUPABASE_KEY
5. Вставьте значения в `.env` файл

### 3. Настройка базы данных
Выполните SQL скрипт в Supabase SQL Editor:

-- Пользователи
CREATE TABLE public.users (
  id bigserial NOT NULL,
  telegram_id bigint NOT NULL,
  username text NULL,
  full_name text NULL,
  role public.user_role NOT NULL DEFAULT 'user'::user_role,
  created_at timestamp with time zone NULL DEFAULT now(),
  CONSTRAINT users_pkey PRIMARY KEY (id),
  CONSTRAINT users_telegram_id_key UNIQUE (telegram_id)
) TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON public.users USING btree (telegram_id) TABLESPACE pg_default;

-- Чеклисты
CREATE TABLE public.checklists (
  id bigserial NOT NULL,
  name text NOT NULL,
  user_id bigint NOT NULL,
  created_at timestamp with time zone NULL DEFAULT now(),
  status public.checklist_status NOT NULL DEFAULT 'draft'::checklist_status,
  CONSTRAINT checklists_pkey PRIMARY KEY (id),
  CONSTRAINT checklists_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
) TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS idx_checklists_user_id ON public.checklists USING btree (user_id) TABLESPACE pg_default;
CREATE INDEX IF NOT EXISTS idx_checklists_status ON public.checklists USING btree (status) TABLESPACE pg_default;

-- Блоки вопросов
CREATE TABLE public.question_blocks (
  id bigserial NOT NULL,
  name text NOT NULL,
  description text NULL,
  checklist_id bigint NOT NULL,
  created_at timestamp with time zone NULL DEFAULT now(),
  CONSTRAINT question_blocks_pkey PRIMARY KEY (id),
  CONSTRAINT question_blocks_checklist_id_fkey FOREIGN KEY (checklist_id) REFERENCES checklists (id) ON DELETE CASCADE
) TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS idx_question_blocks_checklist_id ON public.question_blocks USING btree (checklist_id) TABLESPACE pg_default;

-- Вопросы
CREATE TABLE public.questions (
  id bigserial NOT NULL,
  text text NOT NULL,
  category public.question_category NOT NULL DEFAULT 'compliance'::question_category,
  checklist_id bigint NOT NULL,
  created_at timestamp with time zone NULL DEFAULT now(),
  updated_at timestamp with time zone NULL DEFAULT now(),
  CONSTRAINT questions_pkey PRIMARY KEY (id),
  CONSTRAINT questions_checklist_id_fkey FOREIGN KEY (checklist_id) REFERENCES checklists (id) ON DELETE CASCADE
) TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS idx_questions_category ON public.questions USING btree (category) TABLESPACE pg_default;
CREATE INDEX IF NOT EXISTS idx_questions_checklist_id ON public.questions USING btree (checklist_id) TABLESPACE pg_default;

-- Варианты ответов
CREATE TABLE public.answer_options (
  id bigserial NOT NULL,
  question_id bigint NOT NULL,
  text text NOT NULL,
  is_correct boolean NULL DEFAULT false,
  created_at timestamp with time zone NULL DEFAULT now(),
  CONSTRAINT answer_options_pkey PRIMARY KEY (id),
  CONSTRAINT answer_options_question_id_fkey FOREIGN KEY (question_id) REFERENCES questions (id) ON DELETE CASCADE
) TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS idx_answer_options_question_id ON public.answer_options USING btree (question_id) TABLESPACE pg_default;

-- Шаблоны чеклистов
CREATE TABLE public.checklist_templates (
  id bigserial NOT NULL,
  checklist_id bigint NOT NULL,
  question_id bigint NOT NULL,
  block_id bigint NULL,
  created_at timestamp with time zone NULL DEFAULT now(),
  CONSTRAINT checklist_templates_pkey PRIMARY KEY (id),
  CONSTRAINT checklist_templates_checklist_id_fkey FOREIGN KEY (checklist_id) REFERENCES checklists (id) ON DELETE CASCADE,
  CONSTRAINT checklist_templates_question_id_fkey FOREIGN KEY (question_id) REFERENCES questions (id) ON DELETE CASCADE,
  CONSTRAINT checklist_templates_block_id_fkey FOREIGN KEY (block_id) REFERENCES question_blocks (id) ON DELETE CASCADE
) TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS idx_checklist_templates_checklist_id ON public.checklist_templates USING btree (checklist_id) TABLESPACE pg_default;
CREATE INDEX IF NOT EXISTS idx_checklist_templates_question_id ON public.checklist_templates USING btree (question_id) TABLESPACE pg_default;
CREATE INDEX IF NOT EXISTS idx_checklist_templates_block_id ON public.checklist_templates USING btree (block_id) TABLESPACE pg_default;

### 4. Добавление администратора
После создания таблиц добавьте себя как администратора:
INSERT INTO public.users (telegram_id, username, full_name, role) 
VALUES (ваш_telegram_id, 'ваш_username', 'Ваше Имя', 'admin');

Чтобы узнать ваш Telegram ID: откройте Telegram и найдите @userinfobot

### 5. Установка зависимостей
go mod tidy

### 6. Запуск бота
go run cmd/bot/main.go

Если все настроено правильно, вы увидите:
🤖 <имя_бота> ЗАПУЩЕН
==========================================
🚀 Бот запущен и ожидает сообщений...
==========================================