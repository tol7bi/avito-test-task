
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('employee', 'moderator'))
);

CREATE TABLE IF NOT EXISTS pvz (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    city TEXT NOT NULL CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань'))
);

CREATE TABLE IF NOT EXISTS receptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    pvz_id UUID REFERENCES pvz(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('in_progress', 'close'))
);

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    type TEXT NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь')),
    reception_id UUID REFERENCES receptions(id) ON DELETE CASCADE
);
