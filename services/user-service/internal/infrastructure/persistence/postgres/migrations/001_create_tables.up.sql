SET search_path TO imcs_users;

CREATE TABLE user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    country_code VARCHAR(2) NOT NULL,
    phone VARCHAR(20),
    mobile VARCHAR(20),
    job_title VARCHAR(100),
    department VARCHAR(100),
    avatar_url VARCHAR(500),
    timezone VARCHAR(50),
    locale VARCHAR(10) DEFAULT 'en',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    preference_key VARCHAR(100) NOT NULL,
    preference_value TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, preference_key)
);

CREATE INDEX idx_user_profiles_user_id ON user_profiles(user_id);
CREATE INDEX idx_user_profiles_country ON user_profiles(country_code);
CREATE INDEX idx_user_preferences_user ON user_preferences(user_id);
