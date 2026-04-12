CREATE TABLE IF NOT EXISTS journals (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  current_mileage INTEGER NOT NULL DEFAULT 0 CHECK (current_mileage >= 0 ),
  mileage_last_update_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);

CREATE TABLE IF NOT EXISTS records (
  id INTEGER PRIMARY KEY,
  journal_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  sku TEXT NOT NULL DEFAULT "",
  replaced_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
  next_date INTEGER DEFAULT NULL,
  mileage INTEGER NOT NULL DEFAULT 0 CHECK(mileage >= 0),
  next_mileage INTEGER DEFAULT NULL,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
  FOREIGN KEY (journal_id) REFERENCES journals(id) ON DELETE CASCADE
);
