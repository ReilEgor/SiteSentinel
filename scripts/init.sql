CREATE TABLE IF NOT EXISTS sites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url TEXT NOT NULL UNIQUE,
    interval_seconds INTEGER NOT NULL DEFAULT 30,
    is_up BOOLEAN DEFAULT true,
    last_check TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS check_results (
    id BIGSERIAL PRIMARY KEY,
    site_id UUID NOT NULL,
    status_code INTEGER,
    response_time_ms INTEGER,
    is_up BOOLEAN NOT NULL,
    checked_at TIMESTAMP WITH TIME ZONE NOT NULL,
    
    CONSTRAINT fk_site FOREIGN KEY (site_id) REFERENCES sites(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_sites_last_check ON sites (last_check);